package bukuinduk

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"sekolah/internal/database" // sesuaikan import path package database.DB di project
)

// kolom yang dipakai bersama utk pesertadidik (aktif) & pesertadidik_mutasi
// (skema kedua tabel identik, dikonfirmasi sebelumnya).
const pesertaDidikSelectCols = `
	nis,
	id_rombel,
	COALESCE(anak_ke,'') AS anak_ke,
	COALESCE(saudara_tiri,'') AS saudara_tiri,
	COALESCE(saudara_angkat,'') AS saudara_angkat,
	COALESCE(yatim_piatu,'') AS yatim_piatu,
	COALESCE(pernah_sakit,'') AS pernah_sakit,
	COALESCE(kebutuhan_khusus,'') AS kebutuhan_khusus,
	COALESCE(lama_belajar,'') AS lama_belajar,
	COALESCE(tempat_lahirayah,'') AS tempat_lahirayah,
	COALESCE(tanggal_lahirayah,'') AS tanggal_lahirayah,
	COALESCE(agama_ayah,'') AS agama_ayah,
	COALESCE(warga_ayah,'') AS warga_ayah,
	COALESCE(gaji_ayah,'') AS gaji_ayah,
	COALESCE(alamat_ayah,'') AS alamat_ayah,
	COALESCE(telp_ayah,'') AS telp_ayah,
	COALESCE(status_ayah,'') AS status_ayah,
	COALESCE(tempat_lahiribu,'') AS tempat_lahiribu,
	COALESCE(tanggal_lahiribu,'') AS tanggal_lahiribu,
	COALESCE(agama_ibu,'') AS agama_ibu,
	COALESCE(warga_ibu,'') AS warga_ibu,
	COALESCE(gaji_ibu,'') AS gaji_ibu,
	COALESCE(alamat_ibu,'') AS alamat_ibu,
	COALESCE(telp_ibu,'') AS telp_ibu,
	COALESCE(status_ibu,'') AS status_ibu,
	COALESCE(tempat_lahirwali,'') AS tempat_lahirwali,
	COALESCE(tanggal_lahirwali,'') AS tanggal_lahirwali,
	COALESCE(agama_wali,'') AS agama_wali,
	COALESCE(warga_wali,'') AS warga_wali,
	COALESCE(gaji_wali,'') AS gaji_wali,
	COALESCE(alamat_wali,'') AS alamat_wali,
	COALESCE(telp_wali,'') AS telp_wali,
	COALESCE(hobi_seni,'') AS hobi_seni,
	COALESCE(hobi_orkes,'') AS hobi_orkes,
	COALESCE(hobi_organisasi,'') AS hobi_organisasi,
	COALESCE(hobi_lain,'') AS hobi_lain,
	COALESCE(nisn,'') AS nisn,
	COALESCE(nik_siswa,'') AS nik_siswa,
	COALESCE(no_kk,'') AS no_kk,
	COALESCE(nama_siswa,'') AS nama_siswa,
	COALESCE(jenis_kelamin,'') AS jenis_kelamin,
	COALESCE(tempat_lahirsiswa,'') AS tempat_lahirsiswa,
	COALESCE(tanggal_lahirsiswa,'') AS tanggal_lahirsiswa,
	COALESCE(agama_siswa,'') AS agama_siswa,
	COALESCE(warga_siswa,'') AS warga_siswa,
	COALESCE(saudara_kandung,'') AS saudara_kandung,
	COALESCE(bahasa_siswa,'') AS bahasa_siswa,
	COALESCE(gol_darah,'') AS gol_darah,
	COALESCE(alamat_siswa,'') AS alamat_siswa,
	COALESCE(rt_siswa,'') AS rt_siswa,
	COALESCE(rw_siswa,'') AS rw_siswa,
	COALESCE(kel_siswa,'') AS kel_siswa,
	COALESCE(kec_siswa,'') AS kec_siswa,
	COALESCE(kopos_siswa,'') AS kopos_siswa,
	COALESCE(telprmh_siswa,'') AS telprmh_siswa,
	COALESCE(hp_siswa,'') AS hp_siswa,
	COALESCE(tinggal_siswa,'') AS tinggal_siswa,
	COALESCE(jarak_sekolah,'') AS jarak_sekolah,
	COALESCE(nama_ayah,'') AS nama_ayah,
	COALESCE(pend_ayah,'') AS pend_ayah,
	COALESCE(pek_ayah,'') AS pek_ayah,
	COALESCE(nama_ibu,'') AS nama_ibu,
	COALESCE(pend_ibu,'') AS pend_ibu,
	COALESCE(pek_ibu,'') AS pek_ibu,
	COALESCE(nama_wali,'') AS nama_wali,
	COALESCE(pend_wali,'') AS pend_wali,
	COALESCE(pek_wali,'') AS pek_wali,
	COALESCE(lulus_sd,'') AS lulus_sd,
	COALESCE(tanggal_ijazah,'') AS tanggal_ijazah,
	COALESCE(nomor_ijazah,'') AS nomor_ijazah,
	COALESCE(asal_smp,'') AS asal_smp,
	COALESCE(tanggal_masuk,'') AS tanggal_masuk,
	COALESCE(tinggi_badan,'') AS tinggi_badan,
	COALESCE(berat_badan,'') AS berat_badan,
	COALESCE(foto,'') AS foto,
	COALESCE(status_keluar,'') AS status_keluar,
	COALESCE(alasan_keluar,'') AS alasan_keluar,
	COALESCE(tanggal_keluar,'') AS tanggal_keluar,
	COALESCE(tujuan,'') AS tujuan
`

// pesertaDidikRow adalah hasil scan mentah dari pesertadidik /
// pesertadidik_mutasi (kolom identik pada kedua tabel).
type pesertaDidikRow struct {
	NIS               string `db:"nis"`
	IDRombel          int    `db:"id_rombel"`
	NISN              string `db:"nisn"`
	NIKSiswa          string `db:"nik_siswa"`
	NoKK              string `db:"no_kk"`
	NamaSiswa         string `db:"nama_siswa"`
	JenisKelamin      string `db:"jenis_kelamin"`
	TempatLahirSiswa  string `db:"tempat_lahirsiswa"`
	TanggalLahirSiswa string `db:"tanggal_lahirsiswa"`
	AgamaSiswa        string `db:"agama_siswa"`
	WargaSiswa        string `db:"warga_siswa"`
	SaudaraKandung    string `db:"saudara_kandung"`
	BahasaSiswa       string `db:"bahasa_siswa"`
	GolDarah          string `db:"gol_darah"`
	AlamatSiswa       string `db:"alamat_siswa"`
	RtSiswa           string `db:"rt_siswa"`
	RwSiswa           string `db:"rw_siswa"`
	KelSiswa          string `db:"kel_siswa"`
	KecSiswa          string `db:"kec_siswa"`
	KoposSiswa        string `db:"kopos_siswa"`
	TelprmhSiswa      string `db:"telprmh_siswa"`
	HpSiswa           string `db:"hp_siswa"`
	TinggalSiswa      string `db:"tinggal_siswa"`
	JarakSekolah      string `db:"jarak_sekolah"`
	NamaAyah          string `db:"nama_ayah"`
	PendAyah          string `db:"pend_ayah"`
	PekAyah           string `db:"pek_ayah"`
	NamaIbu           string `db:"nama_ibu"`
	PendIbu           string `db:"pend_ibu"`
	PekIbu            string `db:"pek_ibu"`
	NamaWali          string `db:"nama_wali"`
	PendWali          string `db:"pend_wali"`
	PekWali           string `db:"pek_wali"`
	LulusSd           string `db:"lulus_sd"`
	TanggalIjazah     string `db:"tanggal_ijazah"`
	NomorIjazah       string `db:"nomor_ijazah"`
	AsalSmp           string `db:"asal_smp"`
	TanggalMasuk      string `db:"tanggal_masuk"`
	TinggiBadan       string `db:"tinggi_badan"`
	BeratBadan        string `db:"berat_badan"`
	Foto              string `db:"foto"`
	StatusKeluar      string `db:"status_keluar"`
	AlasanKeluar      string `db:"alasan_keluar"`
	TanggalKeluar     string `db:"tanggal_keluar"`
	Tujuan            string `db:"tujuan"`
	AnakKe            string `db:"anak_ke"`
	SaudaraTiri       string `db:"saudara_tiri"`
	SaudaraAngkat     string `db:"saudara_angkat"`
	YatimPiatu        string `db:"yatim_piatu"`
	PernahSakit       string `db:"pernah_sakit"`
	KebutuhanKhusus   string `db:"kebutuhan_khusus"`
	LamaBelajar       string `db:"lama_belajar"`
	TempatLahirAyah   string `db:"tempat_lahirayah"`
	TanggalLahirAyah  string `db:"tanggal_lahirayah"`
	AgamaAyah         string `db:"agama_ayah"`
	WargaAyah         string `db:"warga_ayah"`
	GajiAyah          string `db:"gaji_ayah"`
	AlamatAyah        string `db:"alamat_ayah"`
	TelpAyah          string `db:"telp_ayah"`
	StatusAyah        string `db:"status_ayah"`
	TempatLahirIbu    string `db:"tempat_lahiribu"`
	TanggalLahirIbu   string `db:"tanggal_lahiribu"`
	AgamaIbu          string `db:"agama_ibu"`
	WargaIbu          string `db:"warga_ibu"`
	GajiIbu           string `db:"gaji_ibu"`
	AlamatIbu         string `db:"alamat_ibu"`
	TelpIbu           string `db:"telp_ibu"`
	StatusIbu         string `db:"status_ibu"`
	TempatLahirWali   string `db:"tempat_lahirwali"`
	TanggalLahirWali  string `db:"tanggal_lahirwali"`
	AgamaWali         string `db:"agama_wali"`
	WargaWali         string `db:"warga_wali"`
	GajiWali          string `db:"gaji_wali"`
	AlamatWali        string `db:"alamat_wali"`
	TelpWali          string `db:"telp_wali"`
	HobiSeni          string `db:"hobi_seni"`
	HobiOrkes         string `db:"hobi_orkes"`
	HobiOrganisasi    string `db:"hobi_organisasi"`
	HobiLain          string `db:"hobi_lain"`
}

// pesertaDidikArsipRow = pesertaDidikRow + kode_siswa (PK asli tabel
// pesertadidik_arsip, karena nis di situ tidak unique).
type pesertaDidikArsipRow struct {
	KodeSiswa int `db:"kode_siswa"`
	pesertaDidikRow
}

// GetAktifByNIS mengambil data siswa aktif utk dicetak Buku Induk.
func GetAktifByNIS(nis string) (*BukuIndukData, error) {
	var row pesertaDidikRow
	query := fmt.Sprintf(`
		SELECT %s FROM pesertadidik WHERE nis = ?`, pesertaDidikSelectCols)
	if err := database.DB.Get(&row, query, nis); err != nil {
		return nil, err
	}
	data := mapRowToBukuInduk(row, SumberAktif)
	data.Kelas = getNamaRombel(row.IDRombel) // tambahan
	if ekskul, errEkskul := GetEkskulSiswa(data.NIS); errEkskul == nil {
		data.EkskulList = ekskul
	}
	if namaWali, nrkWali, errWali := GetWaliKelas(row.IDRombel); errWali == nil {
		data.WaliKelasNama = namaWali
		data.WaliKelasNIPNRK = nrkWali
	}
	if err := enrichKelulusan(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

// helper baru:
func getNamaRombel(idRombel int) string {
	var nama string
	err := database.DB.Get(&nama, `SELECT COALESCE(nama_rombel,'') FROM tbl_rombel WHERE id_rombel = ?`, idRombel)
	if err != nil {
		return ""
	}
	return nama
}

// GetMutasiByNIS mengambil data siswa dari pesertadidik_mutasi.
func GetMutasiByNIS(nis string) (*BukuIndukData, error) {
	var row pesertaDidikRow
	query := fmt.Sprintf(`SELECT %s FROM pesertadidik_mutasi WHERE nis = ?`, pesertaDidikSelectCols)
	if err := database.DB.Get(&row, query, nis); err != nil {
		return nil, err
	}
	data := mapRowToBukuInduk(row, SumberMutasi)
	data.Kelas = getNamaRombel(row.IDRombel) // tambahan
	if ekskul, errEkskul := GetEkskulSiswa(data.NIS); errEkskul == nil {
		data.EkskulList = ekskul
	}
	if namaWali, nrkWali, errWali := GetWaliKelas(row.IDRombel); errWali == nil {
		data.WaliKelasNama = namaWali
		data.WaliKelasNIPNRK = nrkWali
	}
	// siswa mutasi tidak perlu dicek kelulusan (sudah pasti pindah/keluar/meninggal)
	return &data, nil
}

// GetArsipByKodeSiswa mengambil data siswa dari pesertadidik_arsip.
// Identifier-nya kode_siswa (bukan nis), krn satu NIS bisa punya
// banyak baris arsip lintas tahun ajaran.
func GetArsipByKodeSiswa(kodeSiswa int) (*BukuIndukData, error) {
	var row pesertaDidikArsipRow
	query := fmt.Sprintf(`SELECT kode_siswa, %s FROM pesertadidik_arsip WHERE kode_siswa = ?`, pesertaDidikSelectCols)
	if err := database.DB.Get(&row, query, kodeSiswa); err != nil {
		return nil, err
	}
	data := mapRowToBukuInduk(row.pesertaDidikRow, SumberArsip)
	data.Kelas = getNamaRombel(row.IDRombel) // tambahan
	if ekskul, errEkskul := GetEkskulSiswa(data.NIS); errEkskul == nil {
		data.EkskulList = ekskul
	}
	if namaWali, nrkWali, errWali := GetWaliKelas(row.IDRombel); errWali == nil {
		data.WaliKelasNama = namaWali
		data.WaliKelasNIPNRK = nrkWali
	}
	if err := enrichKelulusan(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

// mapRowToBukuInduk adalah converter inti, dipakai bersama oleh sumber
// aktif/mutasi/arsip krn skema kolomnya identik.
func mapRowToBukuInduk(row pesertaDidikRow, sumber SumberData) BukuIndukData {
	d := BukuIndukData{
		Sumber:                  sumber,
		NIS:                     row.NIS,
		NISN:                    row.NISN,
		NIK:                     row.NIKSiswa,
		NoKK:                    row.NoKK,
		NamaLengkap:             row.NamaSiswa,
		NamaPanggilan:           firstWord(row.NamaSiswa),
		JenisKelamin:            formatJenisKelamin(row.JenisKelamin),
		TempatLahir:             row.TempatLahirSiswa,
		TanggalLahir:            formatTanggalDMY(row.TanggalLahirSiswa),
		Agama:                   row.AgamaSiswa,
		Kewarganegaraan:         row.WargaSiswa,
		JumlahSaudara:           row.SaudaraKandung,
		BahasaSehariHari:        row.BahasaSiswa,
		GolonganDarah:           row.GolDarah,
		Alamat:                  formatAlamat(row),
		NoTelpHp:                formatTelpHp(row.TelprmhSiswa, row.HpSiswa),
		BertempatTinggalPada:    row.TinggalSiswa,
		JarakSekolah:            row.JarakSekolah,
		NamaAyah:                row.NamaAyah,
		PendidikanAyah:          row.PendAyah,
		PekerjaanAyah:           row.PekAyah,
		NamaIbu:                 row.NamaIbu,
		PendidikanIbu:           row.PendIbu,
		PekerjaanIbu:            row.PekIbu,
		NamaWali:                row.NamaWali,
		PendidikanWali:          row.PendWali,
		PekerjaanWali:           row.PekWali,
		AsalSekolahMasuk:        row.LulusSd,
		TanggalIjazahMasuk:      formatTanggalDMY(row.TanggalIjazah),
		NomorIjazahMasuk:        row.NomorIjazah,
		AsalSekolahPindahan:     row.AsalSmp,
		DiterimaTanggalPindahan: formatTanggalDMY(row.TanggalMasuk),
		TinggiBadan:             row.TinggiBadan,
		BeratBadan:              row.BeratBadan,
		FotoURL:                 buildFotoURL(sumber, row.Foto),
		AnakKe:                  row.AnakKe,
		SaudaraTiri:             row.SaudaraTiri,
		SaudaraAngkat:           row.SaudaraAngkat,
		YatimPiatu:              row.YatimPiatu,
		PernahSakit:             row.PernahSakit,
		KebutuhanKhusus:         row.KebutuhanKhusus,
		LamaBelajar:             row.LamaBelajar,
		TempatLahirAyah:         row.TempatLahirAyah,
		TanggalLahirAyah:        formatTanggalDMY(row.TanggalLahirAyah),
		AgamaAyah:               row.AgamaAyah,
		WargaAyah:               row.WargaAyah,
		GajiAyah:                row.GajiAyah,
		AlamatAyah:              row.AlamatAyah,
		TelpAyah:                row.TelpAyah,
		StatusAyah:              row.StatusAyah,
		TempatLahirIbu:          row.TempatLahirIbu,
		TanggalLahirIbu:         formatTanggalDMY(row.TanggalLahirIbu),
		AgamaIbu:                row.AgamaIbu,
		WargaIbu:                row.WargaIbu,
		GajiIbu:                 row.GajiIbu,
		AlamatIbu:               row.AlamatIbu,
		TelpIbu:                 row.TelpIbu,
		StatusIbu:               row.StatusIbu,
		TempatLahirWali:         row.TempatLahirWali,
		TanggalLahirWali:        formatTanggalDMY(row.TanggalLahirWali),
		AgamaWali:               row.AgamaWali,
		WargaWali:               row.WargaWali,
		GajiWali:                row.GajiWali,
		AlamatWali:              row.AlamatWali,
		TelpWali:                row.TelpWali,
		HobiSeni:                row.HobiSeni,
		HobiOrkes:               row.HobiOrkes,
		HobiOrganisasi:          row.HobiOrganisasi,
		HobiLain:                row.HobiLain,
	}

	// Bagian D: default berdasar status_keluar mentah (relevan utk
	// sumber mutasi; aktif & arsip selalu kosong sesuai temuan data).
	switch strings.ToLower(strings.TrimSpace(row.StatusKeluar)) {
	case "pindah":
		d.StatusKeluar = StatusKeluarPindah
		d.PindahKeSekolah = row.Tujuan
	case "keluar":
		d.StatusKeluar = StatusKeluarKeluar
		d.KeluarAlasan = row.AlasanKeluar
		d.KeluarTanggal = formatTanggalDMY(row.TanggalKeluar)
	case "meninggal":
		d.StatusKeluar = StatusKeluarKeluar
		d.KeluarAlasan = "Meninggal"
		d.KeluarTanggal = formatTanggalDMY(row.TanggalKeluar)
	default:
		d.StatusKeluar = StatusKeluarNone
	}

	return d
}

// enrichKelulusan cek tabel kelulusan utk NIS ybs; kalau ada baris &
// kolom ket mengandung "lulus" (tanpa "tidak"), Bagian D DITIMPA jadi
// cabang Lulus (poin 17) — bisa override hasil switch status_keluar
// di atas kalau kedua sumber sama-sama ada (kasus langka).
func enrichKelulusan(d *BukuIndukData) error {
	var row struct {
		Ket     string `db:"ket"`
		TahunAj string `db:"nama_tahunpelajaran"`
	}
	query := `
		SELECT COALESCE(k.ket,'') AS ket, COALESCE(tp.nama_tahunpelajaran,'') AS nama_tahunpelajaran
		FROM kelulusan k
		LEFT JOIN tbl_tahunpelajaran tp ON tp.id_tahunpelajaran = k.id_tahunpelajaran
		WHERE k.nis = ?
		LIMIT 1`
	err := database.DB.Get(&row, query, d.NIS)
	if err == sql.ErrNoRows {
		return nil // tidak ada record kelulusan, biarkan hasil status_keluar apa adanya
	}
	if err != nil {
		return err
	}
	ketLower := strings.ToLower(row.Ket)
	if strings.Contains(ketLower, "lulus") && !strings.Contains(ketLower, "tidak") {
		d.StatusKeluar = StatusKeluarLulus
		d.LulusTahun = row.TahunAj
	}
	return nil
}

// GetSekolahProfil mengambil profil sekolah dari tabel setting, dipakai
// utk halaman Cover.
func GetSekolahProfil() (*SekolahProfil, error) {
	var p SekolahProfil
	query := `
		SELECT
			COALESCE(nama_sekolah,'') AS nama_sekolah,
			COALESCE(alamat_sekolah,'') AS alamat_sekolah,
			COALESCE(kelurahan_sekolah,'') AS kelurahan_sekolah,
			COALESCE(kecamatan_sekolah,'') AS kecamatan_sekolah,
			COALESCE(kabupaten_sekolah,'') AS kabupaten_sekolah,
			COALESCE(provinsi_sekolah,'') AS provinsi_sekolah,
			COALESCE(nsss,'') AS nss,
			COALESCE(npsn,'') AS npsn,
			COALESCE(nsp,'') AS nsp
		FROM setting LIMIT 1`
	if err := database.DB.Get(&p, query); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetKepalaSekolahInfo mengambil nama & NIP kepala sekolah dari
// kelulusan_waktu (single-row id_waktu=1), reuse pola modul Kelulusan.
func GetKepalaSekolahInfo() (*KepalaSekolahInfo, error) {
	var k KepalaSekolahInfo
	query := `
		SELECT COALESCE(kepsek,'') AS nama_kepsek, COALESCE(nipkepsek,'') AS nip_kepsek
		FROM kelulusan_waktu WHERE id_waktu = 1`
	if err := database.DB.Get(&k, query); err != nil {
		return nil, err
	}
	return &k, nil
}

// --- helper kecil ---

func firstWord(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	return strings.Fields(s)[0]
}

func formatJenisKelamin(kode string) string {
	switch strings.ToUpper(strings.TrimSpace(kode)) {
	case "L":
		return "Laki-laki"
	case "P":
		return "Perempuan"
	default:
		return ""
	}
}

func formatAlamat(row pesertaDidikRow) string {
	var parts []string
	if row.AlamatSiswa != "" {
		parts = append(parts, row.AlamatSiswa)
	}
	if row.RtSiswa != "" || row.RwSiswa != "" {
		parts = append(parts, fmt.Sprintf("RT %s/RW %s", row.RtSiswa, row.RwSiswa))
	}
	if row.KelSiswa != "" {
		parts = append(parts, row.KelSiswa)
	}
	if row.KecSiswa != "" {
		parts = append(parts, row.KecSiswa)
	}
	if row.KoposSiswa != "" {
		parts = append(parts, row.KoposSiswa)
	}
	return strings.Join(parts, ", ")
}

func formatTelpHp(telprmh, hp string) string {
	switch {
	case telprmh != "" && hp != "":
		return fmt.Sprintf("%s / %s", telprmh, hp)
	case hp != "":
		return hp
	default:
		return telprmh
	}
}

func buildFotoURL(sumber SumberData, filename string) string {
	if filename == "" {
		return ""
	}
	folder := "fotosiswa"
	if sumber == SumberArsip {
		folder = "fotosiswaarsip" // konsisten dgn pola foto arsip pesertadidik yg sudah ada
	}
	return fmt.Sprintf("/storage/%s/%s", folder, filename)
}

// GetNamaPengguna mengambil nama akun dari tabel pengguna, dipakai utk
// blok tanda tangan "Petugas" (nama admin yg sedang login mencetak).
func GetNamaPengguna(idPengguna string) (string, error) {
	var nama string
	err := database.DB.Get(&nama, `SELECT COALESCE(nama_pengguna,'') FROM pengguna WHERE id_pengguna = ?`, idPengguna)
	if err != nil {
		return "", err
	}
	return nama, nil
}

// GetRentangNIS mengambil NIS terkecil & terbesar dari pesertadidik
// aktif utk tahun pelajaran tertentu (dipakai di halaman Daftar Isi).
func GetRentangNIS(tahunPelajaran string) (string, string, error) {
	var result struct {
		MinNIS string `db:"min_nis"`
		MaxNIS string `db:"max_nis"`
	}
	query := `
		SELECT COALESCE(MIN(pd.nis),'') AS min_nis, COALESCE(MAX(pd.nis),'') AS max_nis
		FROM pesertadidik pd
		JOIN tbl_tahunpelajaran tp ON tp.id_tahunpelajaran = pd.id_tahunpelajaran
		WHERE tp.nama_tahunpelajaran = ?`
	if err := database.DB.Get(&result, query, tahunPelajaran); err != nil {
		return "", "", err
	}
	return result.MinNIS, result.MaxNIS, nil
}

// GetNIPPegawai mengambil nip_baru dari tabel pegawai berdasarkan
// nuptk, dipakai utk blok tanda tangan "Petugas" kalau admin yg
// mencetak kebetulan juga tercatat sbg pegawai (punya NIP).
func GetNIPPegawai(nuptk string) (string, error) {
	var nip string
	err := database.DB.Get(&nip, `SELECT COALESCE(nip_baru,'') FROM pegawai WHERE nuptk = ?`, nuptk)
	if err != nil {
		return "", err
	}
	return nip, nil
}

// formatTanggalDMY mengubah tanggal dari format DB (yyyy-mm-dd) ke
// format tampilan Indonesia (dd-mm-yyyy). Kalau formatnya tak
// dikenali/kosong, tampilkan apa adanya (tidak dipaksa parse).
func formatTanggalDMY(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02", trimmed)
	if err != nil {
		return raw
	}
	return t.Format("02-01-2006")
}

// GetWaliKelas mengambil nama & NRK wali kelas suatu rombel, via JOIN
// tbl_walikelas (nuptk) -> pegawai (nama_pegawai, nrk).
func GetWaliKelas(idRombel int) (nama string, nip string, err error) {
	var row struct {
		Nama string `db:"nama_pegawai"`
		NIP  string `db:"nip_baru"`
	}
	query := `
		SELECT COALESCE(p.nama_pegawai,'') AS nama_pegawai, COALESCE(p.nip_baru,'') AS nip_baru
		FROM tbl_walikelas wk
		JOIN pegawai p ON p.nuptk = wk.nuptk
		WHERE wk.id_rombel = ?
		LIMIT 1`
	e := database.DB.Get(&row, query, idRombel)
	if e == sql.ErrNoRows {
		return "", "", nil
	}
	if e != nil {
		return "", "", e
	}
	return row.Nama, row.NIP, nil
}

// GetEkskulSiswa mengambil daftar nama ekskul yg diikuti siswa, dipakai
// utk tabel Ekstrakurikuler Bagian K (cuma nama, tanpa nilai/keterangan
// krn belum ada sumbernya).
func GetEkskulSiswa(nis string) ([]string, error) {
	var nama []string
	query := `
		SELECT e.nama_ekskul
		FROM tbl_anggota_ekskul a
		JOIN tbl_ekskul e ON e.id_ekskul = a.id_ekskul
		WHERE a.nis = ? AND e.aktif = 1
		ORDER BY e.nama_ekskul`
	if err := database.DB.Select(&nama, query, nis); err != nil {
		return nil, err
	}
	return nama, nil
}
