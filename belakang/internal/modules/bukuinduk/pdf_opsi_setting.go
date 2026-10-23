package bukuinduk

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"sekolah/internal/database"
)

// cetakOpsiRow = 1 baris singleton (id selalu 1) tabel bukuinduk_opsi_cetak - setting Opsi Cetak
// default yg diisi admin lewat menu Buku Induk > Pengaturan.
type cetakOpsiRow struct {
	Sampul        bool   `db:"sampul"`
	DaftarIsi     bool   `db:"daftar_isi"`
	Pemeriksaan   bool   `db:"pemeriksaan"`
	NomorUrut     bool   `db:"nomor_urut"`
	Petunjuk      bool   `db:"petunjuk"`
	LembarIV      bool   `db:"lembar_iv"`
	HalamanAkhir  bool   `db:"halaman_akhir"`
	Tingkat       string `db:"tingkat"` // csv "7,8,9"
	KepsekNama    string `db:"kepsek_nama_override"`
	KepsekNIP     string `db:"kepsek_nip_override"`
	WaliKelasNama string `db:"walikelas_nama_override"`
	WaliKelasNIP  string `db:"walikelas_nip_override"`
}

func parseTingkatCSV(s string) []int {
	var out []int

	for _, part := range strings.Split(s, ",") {
		if n, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			out = append(out, n)
		}
	}

	return out
}

func tingkatToCSV(list []int) string {
	parts := make([]string, 0, len(list))
	for _, n := range list {
		parts = append(parts, strconv.Itoa(n))
	}

	return strings.Join(parts, ",")
}

// GetCetakOpsiSetting membaca Opsi Cetak default. Belum pernah disimpan admin (tabel masih
// kosong) -> DefaultCetakOpsi() (cetak semua lembar, spt perilaku sblm fitur Pengaturan ada).
func GetCetakOpsiSetting() (CetakOpsi, error) {
	def := DefaultCetakOpsi()

	var row cetakOpsiRow

	err := database.DB.Get(&row, `
		SELECT sampul, daftar_isi, pemeriksaan, nomor_urut, petunjuk, lembar_iv, halaman_akhir,
			tingkat, kepsek_nama_override, kepsek_nip_override, walikelas_nama_override, walikelas_nip_override
		FROM bukuinduk_opsi_cetak WHERE id = 1 LIMIT 1
	`)
	if err != nil {
		if err == sql.ErrNoRows {
			return def, nil
		}

		return def, err
	}

	tingkat := parseTingkatCSV(row.Tingkat)
	if len(tingkat) == 0 {
		tingkat = def.Tingkat
	}

	return CetakOpsi{
		Sampul: row.Sampul, DaftarIsi: row.DaftarIsi, Pemeriksaan: row.Pemeriksaan,
		NomorUrut: row.NomorUrut, Petunjuk: row.Petunjuk, LembarIV: row.LembarIV,
		HalamanAkhir: row.HalamanAkhir, Tingkat: tingkat,
		KepsekNamaOverride: row.KepsekNama, KepsekNIPOverride: row.KepsekNIP,
		WaliKelasNamaOverride: row.WaliKelasNama, WaliKelasNIPOverride: row.WaliKelasNIP,
	}, nil
}

// SaveCetakOpsiSetting menyimpan Opsi Cetak default (1 baris singleton, id selalu 1).
func SaveCetakOpsiSetting(o CetakOpsi) error {
	_, err := database.DB.Exec(`
		INSERT INTO bukuinduk_opsi_cetak (
			id, sampul, daftar_isi, pemeriksaan, nomor_urut, petunjuk, lembar_iv, halaman_akhir,
			tingkat, kepsek_nama_override, kepsek_nip_override, walikelas_nama_override, walikelas_nip_override
		) VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			sampul = VALUES(sampul), daftar_isi = VALUES(daftar_isi), pemeriksaan = VALUES(pemeriksaan),
			nomor_urut = VALUES(nomor_urut), petunjuk = VALUES(petunjuk), lembar_iv = VALUES(lembar_iv),
			halaman_akhir = VALUES(halaman_akhir), tingkat = VALUES(tingkat),
			kepsek_nama_override = VALUES(kepsek_nama_override),
			kepsek_nip_override = VALUES(kepsek_nip_override),
			walikelas_nama_override = VALUES(walikelas_nama_override),
			walikelas_nip_override = VALUES(walikelas_nip_override)
	`,
		o.Sampul, o.DaftarIsi, o.Pemeriksaan, o.NomorUrut, o.Petunjuk, o.LembarIV, o.HalamanAkhir,
		tingkatToCSV(o.Tingkat), o.KepsekNamaOverride, o.KepsekNIPOverride,
		o.WaliKelasNamaOverride, o.WaliKelasNIPOverride,
	)

	return err
}

// GetOpsiCetakHandler — GET /bukuinduk/opsi-cetak: isi form halaman Pengaturan.
func GetOpsiCetakHandler(c *fiber.Ctx) error {
	opsi, err := GetCetakOpsiSetting()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": opsi.ToResponse()})
}

// SaveOpsiCetakHandler — PUT /bukuinduk/opsi-cetak: simpan Opsi Cetak default dari halaman
// Pengaturan. Sesudah ini, CetakSatuHandler/CetakBulkHandler yg dipanggil TANPA override lain
// (dari tombol Cetak di daftar Peserta Didik) otomatis memakai pengaturan ini.
func SaveOpsiCetakHandler(c *fiber.Ctx) error {
	var req CetakOpsiRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Body tidak valid"})
	}

	opsi := req.ToCetakOpsi()
	if !opsi.adaLembarDicetak() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Pilih minimal 1 lembar yang akan dicetak."})
	}

	if err := SaveCetakOpsiSetting(opsi); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal menyimpan pengaturan: " + err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pengaturan cetak Buku Induk berhasil disimpan."})
}
