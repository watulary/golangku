package bukuinduk

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// CetakOpsi mengatur lembar/halaman mana saja yg ikut dicetak, dan boleh override nama Kepala
// Sekolah utk 1x proses cetak ini saja (TANPA mengubah setting global). Dipakai mis. utk cetak
// arsip semester Ganjil yg baru berjalan (Kelas 7/8/9 belum semuanya terisi nilai): halaman
// depan buku (Sampul/Daftar Isi/dst, biasanya cuma dicetak 1x di awal tahun pelajaran) & tingkat
// yg datanya belum diisi bisa dilewati, supaya tidak ikut cetak halaman kosong.
type CetakOpsi struct {
	Sampul       bool
	DaftarIsi    bool
	Pemeriksaan  bool
	NomorUrut    bool
	Petunjuk     bool
	LembarIV     bool  // biodata siswa (Bagian A-E: Keterangan Siswa/Ortu/Perkembangan/dst)
	Tingkat      []int // tingkat Data Nilai yg dicetak, subset dari {7, 8, 9}; kosong = semua
	HalamanAkhir bool  // hlm penutup: Rata-rata/Ujian + Ekstrakurikuler + Keputusan Lulus/Tidak Lulus

	KepsekNamaOverride string // "" = pakai setting Kepala Sekolah global (tabel kelulusan_waktu)
	KepsekNIPOverride  string

	// override nama/NIP Wali Kelas yg ditampilkan di blok ttd halaman depan (Daftar Isi &
	// Nomor Urut) - halaman itu representasinya cuma dari SISWA PERTAMA di daftar yg dicetak,
	// jadi utk cetak lintas rombel/wali kelas beda2, admin bisa isi manual siapa yg dianggap
	// bertanggung jawab. TIDAK memengaruhi data per-siswa di halaman Data Nilai masing2 (itu
	// tetap dari data rombel siswa yg sebenarnya).
	WaliKelasNamaOverride string
	WaliKelasNIPOverride  string
}

// DefaultCetakOpsi = perilaku lama (semua lembar dicetak, tanpa override Kepala Sekolah).
func DefaultCetakOpsi() CetakOpsi {
	return CetakOpsi{
		Sampul: true, DaftarIsi: true, Pemeriksaan: true, NomorUrut: true,
		Petunjuk: true, LembarIV: true, Tingkat: []int{7, 8, 9}, HalamanAkhir: true,
	}
}

// tingkatMap mengubah Tingkat jadi lookup cepat dipakai renderDataNilaiPerKelas.
func (o CetakOpsi) tingkatMap() map[int]bool {
	m := map[int]bool{}
	for _, t := range o.Tingkat {
		if t == 7 || t == 8 || t == 9 {
			m[t] = true
		}
	}
	return m
}

// terapkanKepsek mengembalikan KepalaSekolahInfo dgn nama/NIP yg sudah dioverride (kalau diisi).
func (o CetakOpsi) terapkanKepsek(kepsek KepalaSekolahInfo) KepalaSekolahInfo {
	if o.KepsekNamaOverride != "" {
		kepsek.NamaKepsek = o.KepsekNamaOverride
	}

	if o.KepsekNIPOverride != "" {
		kepsek.NIPKepsek = o.KepsekNIPOverride
	}

	return kepsek
}

// adaLembarDicetak: kalau semua bagian non-siswa dimatikan DAN tingkat data nilai juga kosong,
// tidak ada apa pun yg akan dicetak (PDF kosong) - dicegah di handler sblm generate.
func (o CetakOpsi) adaLembarDicetak() bool {
	return o.Sampul || o.DaftarIsi || o.Pemeriksaan || o.NomorUrut || o.Petunjuk ||
		o.LembarIV || o.HalamanAkhir || len(o.tingkatMap()) > 0
}

// CetakOpsiRequest adalah bentuk opsi cetak dari luar (query string GET / JSON body POST).
// Pointer bool: nil = ikut default (kompatibel dgn caller lama yg belum kirim opsi ini sama
// sekali -> tetap cetak lengkap spt sebelum fitur ini ada).
type CetakOpsiRequest struct {
	Sampul        *bool  `json:"sampul,omitempty"`
	DaftarIsi     *bool  `json:"daftar_isi,omitempty"`
	Pemeriksaan   *bool  `json:"pemeriksaan,omitempty"`
	NomorUrut     *bool  `json:"nomor_urut,omitempty"`
	Petunjuk      *bool  `json:"petunjuk,omitempty"`
	LembarIV      *bool  `json:"lembar_iv,omitempty"`
	Tingkat       []int  `json:"tingkat,omitempty"`
	HalamanAkhir  *bool  `json:"halaman_akhir,omitempty"`
	KepsekNama    string `json:"kepsek_nama,omitempty"`
	KepsekNIP     string `json:"kepsek_nip,omitempty"`
	WaliKelasNama string `json:"walikelas_nama,omitempty"`
	WaliKelasNIP  string `json:"walikelas_nip,omitempty"`
}

// ToCetakOpsi mengubah request jadi CetakOpsi siap pakai, mengisi bagian yg tidak dikirim (nil)
// dgn DefaultCetakOpsi() (cetak semua, spt perilaku lama sblm ada Pengaturan).
func (r CetakOpsiRequest) ToCetakOpsi() CetakOpsi {
	return r.ToCetakOpsiWithBase(DefaultCetakOpsi())
}

// ToCetakOpsiWithBase sama seperti ToCetakOpsi, tapi bagian yg tidak dikirim (nil) diisi dari
// base, bukan selalu DefaultCetakOpsi() - dipakai CetakSatuHandler/CetakBulkHandler supaya kalau
// tidak ada override yg dikirim dari halaman Cetak, hasilnya mengikuti Pengaturan (menu Buku
// Induk > Pengaturan) yg sudah disimpan admin, bukan selalu "cetak semua".
func (r CetakOpsiRequest) ToCetakOpsiWithBase(def CetakOpsi) CetakOpsi {
	pick := func(p *bool, fallback bool) bool {
		if p == nil {
			return fallback
		}

		return *p
	}

	tingkat := r.Tingkat
	if tingkat == nil {
		tingkat = def.Tingkat
	}

	// string (bukan pointer) tidak bisa membedakan "tidak dikirim" dari "dikirim kosong", jadi
	// kosong = ikut base (def) - konsisten dgn cara Tampil.vue memanggil cetak (tidak pernah
	// mengirim field ini sama sekali), supaya override dari Pengaturan tetap terpakai.
	pickStr := func(sent, fallback string) string {
		sent = strings.TrimSpace(sent)
		if sent == "" {
			return fallback
		}

		return sent
	}

	return CetakOpsi{
		Sampul:                pick(r.Sampul, def.Sampul),
		DaftarIsi:             pick(r.DaftarIsi, def.DaftarIsi),
		Pemeriksaan:           pick(r.Pemeriksaan, def.Pemeriksaan),
		NomorUrut:             pick(r.NomorUrut, def.NomorUrut),
		Petunjuk:              pick(r.Petunjuk, def.Petunjuk),
		LembarIV:              pick(r.LembarIV, def.LembarIV),
		Tingkat:               tingkat,
		HalamanAkhir:          pick(r.HalamanAkhir, def.HalamanAkhir),
		KepsekNamaOverride:    pickStr(r.KepsekNama, def.KepsekNamaOverride),
		KepsekNIPOverride:     pickStr(r.KepsekNIP, def.KepsekNIPOverride),
		WaliKelasNamaOverride: pickStr(r.WaliKelasNama, def.WaliKelasNamaOverride),
		WaliKelasNIPOverride:  pickStr(r.WaliKelasNIP, def.WaliKelasNIPOverride),
	}
}

// queryBoolPtr membaca query bool dari GET (?sampul=1 / ?sampul=0), nil kalau parameternya
// tidak dikirim sama sekali (beda dgn dikirim "0") - supaya default lama tetap jalan.
func queryBoolPtr(c *fiber.Ctx, key string) *bool {
	if !c.Context().QueryArgs().Has(key) {
		return nil
	}

	v := c.Query(key)
	b := v == "1" || strings.EqualFold(v, "true")

	return &b
}

// parseCetakOpsiQuery membaca CetakOpsiRequest dari query string (dipakai CetakSatuHandler).
func parseCetakOpsiQuery(c *fiber.Ctx) CetakOpsiRequest {
	var tingkat []int
	if raw := strings.TrimSpace(c.Query("tingkat")); raw != "" {
		for _, part := range strings.Split(raw, ",") {
			if n, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
				tingkat = append(tingkat, n)
			}
		}
	}

	return CetakOpsiRequest{
		Sampul:        queryBoolPtr(c, "sampul"),
		DaftarIsi:     queryBoolPtr(c, "daftar_isi"),
		Pemeriksaan:   queryBoolPtr(c, "pemeriksaan"),
		NomorUrut:     queryBoolPtr(c, "nomor_urut"),
		Petunjuk:      queryBoolPtr(c, "petunjuk"),
		LembarIV:      queryBoolPtr(c, "lembar_iv"),
		Tingkat:       tingkat,
		HalamanAkhir:  queryBoolPtr(c, "halaman_akhir"),
		KepsekNama:    c.Query("kepsek_nama"),
		KepsekNIP:     c.Query("kepsek_nip"),
		WaliKelasNama: c.Query("walikelas_nama"),
		WaliKelasNIP:  c.Query("walikelas_nip"),
	}
}

// CetakOpsiResponse = bentuk CetakOpsi yg dikirim ke frontend (halaman Pengaturan), field
// override dikirim polos sbg string (bukan pointer) krn cuma ditampilkan di form isian.
type CetakOpsiResponse struct {
	Sampul        bool   `json:"sampul"`
	DaftarIsi     bool   `json:"daftar_isi"`
	Pemeriksaan   bool   `json:"pemeriksaan"`
	NomorUrut     bool   `json:"nomor_urut"`
	Petunjuk      bool   `json:"petunjuk"`
	LembarIV      bool   `json:"lembar_iv"`
	HalamanAkhir  bool   `json:"halaman_akhir"`
	Tingkat       []int  `json:"tingkat"`
	KepsekNama    string `json:"kepsek_nama"`
	KepsekNIP     string `json:"kepsek_nip"`
	WaliKelasNama string `json:"walikelas_nama"`
	WaliKelasNIP  string `json:"walikelas_nip"`
}

func (o CetakOpsi) ToResponse() CetakOpsiResponse {
	return CetakOpsiResponse{
		Sampul: o.Sampul, DaftarIsi: o.DaftarIsi, Pemeriksaan: o.Pemeriksaan,
		NomorUrut: o.NomorUrut, Petunjuk: o.Petunjuk, LembarIV: o.LembarIV,
		HalamanAkhir: o.HalamanAkhir, Tingkat: o.Tingkat,
		KepsekNama: o.KepsekNamaOverride, KepsekNIP: o.KepsekNIPOverride,
		WaliKelasNama: o.WaliKelasNamaOverride, WaliKelasNIP: o.WaliKelasNIPOverride,
	}
}
