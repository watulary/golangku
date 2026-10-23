package pesertadidik

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

// importHeaders urutan HARUS SAMA PERSIS dengan pesertaDidikColumns (index 1..119),
// index 0 ("No") diabaikan saat parsing.
var importHeaders = []string{
	"No",
	"ID Tahun Pelajaran", "ID Kelas", "ID Rombel",
	"Nama Siswa", "NIS", "NISN", "No. Ujian SD", "No. Simdik", "Jenis Kelamin",
	"Tempat Lahir", "Tanggal Lahir", "No. Akte", "Kota Akte", "Agama", "Kewarganegaraan",
	"Anak Ke", "Saudara Kandung", "Saudara Tiri", "Saudara Angkat", "Yatim/Piatu",
	"Bahasa Sehari-hari", "Penerima Bantuan", "No. KK", "NIK Siswa",
	"Alamat", "RT", "RW", "Kode Pos", "Kelurahan/Desa", "Kecamatan",
	"Telp Rumah", "No. HP", "Tinggal Bersama", "Jarak ke Sekolah",
	"Golongan Darah", "Riwayat Sakit", "Status Jasmani", "Kebutuhan Khusus",
	"Tinggi Badan", "Berat Badan", "Lingkar Kepala", "Moda Transportasi",
	"Lulus Dari (SD)", "Tanggal Ijazah", "Kota Ijazah", "Nomor Ijazah", "Tahun Ijazah",
	"Nomor SKHUN", "Tahun SKHUN", "Nilai Ijazah", "Nilai NEM", "Lama Belajar",
	"Jalur Masuk", "Semester", "Asal Sekolah", "Tanggal Masuk", "Email",
	"Nama Ayah", "Tempat Lahir Ayah", "Tanggal Lahir Ayah", "Agama Ayah", "Kewarganegaraan Ayah",
	"Pendidikan Ayah", "Pekerjaan Ayah", "Penghasilan Ayah", "Kategori Penghasilan Ayah",
	"NIK Ayah", "Alamat Ayah", "RT Ayah", "RW Ayah", "Kode Pos Ayah",
	"Kelurahan Ayah", "Kecamatan Ayah", "No. Telp/HP Ayah", "Status Ayah",
	"Nama Ibu", "Tempat Lahir Ibu", "Tanggal Lahir Ibu", "Agama Ibu", "Kewarganegaraan Ibu",
	"Pendidikan Ibu", "Pekerjaan Ibu", "Penghasilan Ibu", "Kategori Penghasilan Ibu",
	"NIK Ibu", "Alamat Ibu", "RT Ibu", "RW Ibu", "Kode Pos Ibu",
	"Kelurahan Ibu", "Kecamatan Ibu", "No. Telp/HP Ibu", "Status Ibu",
	"Nama Wali", "Tempat Lahir Wali", "Tanggal Lahir Wali", "Agama Wali", "Kewarganegaraan Wali",
	"Pendidikan Wali", "Pekerjaan Wali", "Penghasilan Wali", "Kategori Penghasilan Wali",
	"Alamat Wali", "RT Wali", "RW Wali", "Kode Pos Wali", "Kelurahan Wali", "Kecamatan Wali", "No. Telp/HP Wali",
	"Hobi Seni", "Hobi Olahraga", "Hobi Organisasi", "Hobi Lainnya",
	"Status Keluar", "Alasan Keluar", "Tanggal Keluar", "Nomor Surat", "Tujuan",
	"Foto",
}

func getCell(row []string, idx int) string {
	if idx < len(row) {
		return strings.TrimSpace(row[idx])
	}
	return ""
}

func isRowBlank(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// validateDigitCell: exact true = harus PERSIS length digit. exact false = MAKSIMAL length digit.
// Kosong dianggap valid (field opsional), kecuali required=true.
func validateDigitCell(value, label string, length int, exact, required bool) string {
	if value == "" {
		if required {
			return fmt.Sprintf("%s wajib diisi", label)
		}
		return ""
	}
	if !isAllDigits(value) {
		return fmt.Sprintf("%s harus angka", label)
	}
	if exact && len(value) != length {
		return fmt.Sprintf("%s harus %d digit angka", label, length)
	}
	if !exact && len(value) > length {
		return fmt.Sprintf("%s maksimal %d digit angka", label, length)
	}
	return ""
}

// TemplateHandler menghasilkan file template import berisi sheet data + sheet petunjuk lengkap
func TemplateHandler(c *fiber.Ctx) error {
	f := excelize.NewFile()
	sheet := "Peserta Didik"
	f.SetSheetName("Sheet1", sheet)

	for i, h := range importHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	lastCol, _ := excelize.ColumnNumberToName(len(importHeaders))
	secondCol, _ := excelize.ColumnNumberToName(2)
	f.SetColWidth(sheet, "A", "A", 6)
	f.SetColWidth(sheet, secondCol, lastCol, 18)
	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})
	customFmt := "@"
	if textStyle, err := f.NewStyle(&excelize.Style{CustomNumFmt: &customFmt}); err == nil {
		f.SetColStyle(sheet, fmt.Sprintf("A:%s", lastCol), textStyle)
	}

	if style := pesertaDidikExcelHeaderStyle(f); style != -1 {
		firstCell, _ := excelize.CoordinatesToCellName(1, 1)
		lastHeaderCell, _ := excelize.CoordinatesToCellName(len(importHeaders), 1)
		f.SetCellStyle(sheet, firstCell, lastHeaderCell, style)
	}
	if style := pesertaDidikExcelHeaderStyle(f); style != -1 {
		firstCell, _ := excelize.CoordinatesToCellName(1, 1)
		lastHeaderCell, _ := excelize.CoordinatesToCellName(len(importHeaders), 1)
		f.SetCellStyle(sheet, firstCell, lastHeaderCell, style)
	}

	// Sheet petunjuk pengisian
	petunjuk := "Petunjuk"
	f.NewSheet(petunjuk)
	instructions := []string{
		"PETUNJUK PENGISIAN TEMPLATE UNGGAH PESERTA DIDIK",
		"",
		"1. Jangan mengubah urutan atau nama kolom pada sheet 'Peserta Didik'.",
		"2. Kolom 'No' boleh dikosongkan, tidak diproses saat import.",
		"3. Kolom 'ID Tahun Pelajaran', 'ID Kelas', 'ID Rombel' WAJIB diisi ANGKA ID (bukan nama) sesuai data master. Cek ID-nya di halaman Tahun Pelajaran / Kelas / Rombel pada dashboard.",
		"4. Kolom NIS wajib diisi, maksimal 12 digit angka, tidak boleh duplikat antar baris maupun dengan data yang sudah ada di sistem.",
		"5. Kolom NISN, jika diisi, harus PERSIS 10 digit angka, tidak boleh duplikat.",
		"6. Kolom No. Ujian SD, jika diisi, maksimal 22 digit angka, tidak boleh duplikat.",
		"7. Kolom NIK Siswa, NIK Ayah, NIK Ibu, jika diisi, harus PERSIS 16 digit angka, masing-masing tidak boleh duplikat.",
		"8. Kolom No. KK, jika diisi, harus PERSIS 16 digit angka (boleh sama antar kakak-adik, tidak dicek duplikat).",
		"9. Format semua kolom tanggal: YYYY-MM-DD (contoh: 2012-05-17).",
		"10. Kolom 'Foto' dikosongkan saja.",
		"11. Jika salah satu baris bermasalah, SELURUH proses import akan dibatalkan (tidak ada data yang tersimpan sebagian).",
		"12. Beberapa kolom hanya menerima nilai tertentu (harus ditulis PERSIS seperti daftar di bawah, huruf besar/kecil dan spasi berpengaruh):",
		"",
		"Kolom 'Jenis Kelamin' — pilih salah satu:",
		"L | P",
		"",
		"Kolom 'Agama', 'Agama Ayah', 'Agama Ibu', 'Agama Wali' — pilih salah satu:",
		"Islam | Kristen | Katholik | Hindu | Budha | Khonghucu | Kepercayaan kpd Tuhan YME | lainnya",
		"",
		"Kolom 'Kewarganegaraan', 'Kewarganegaraan Ayah', 'Kewarganegaraan Ibu', 'Kewarganegaraan Wali' — pilih salah satu:",
		"Indonesia | Asing",
		"",
		"Kolom 'Golongan Darah' — pilih salah satu:",
		"A | B | AB | O",
		"",
		"Kolom 'Tinggal Bersama' — pilih salah satu:",
		"Bersama orang tua | Wali | Kost | Asrama | Panti asuhan | Pesantren | Lainnya",
		"",
		"Kolom 'Yatim/Piatu' — pilih salah satu:",
		"Yatim | Piatu | Yatim Piatu",
		"",
		"Kolom 'Bahasa Sehari-hari' — pilih salah satu:",
		"Indonesia | Asing | Daerah",
		"",
		"Kolom 'Penerima Bantuan' — pilih salah satu:",
		"KJP | KIP | BSM",
		"",
		"Kolom 'Moda Transportasi' — pilih salah satu:",
		"Jalan kaki | Angkutan umum/bus/pete-pete | Mobil/bus antar jemput | Kereta api | Ojek",
		"Andong/bendi/sado/dokar/delman/becak | Perahu penyeberangan/rakit/getek | Kuda | Sepeda",
		"Sepeda motor | Mobil pribadi | Lainnya",
		"",
		"Kolom 'Status Jasmani' — pilih salah satu:",
		"Normal | Cacat",
		"",
		"Kolom 'Kebutuhan Khusus' — pilih salah satu:",
		"Tidak ada | netra (A) | rungu (B) | grahita ringan (C) | grahita sedang (C1) | daksa ringan (D)",
		"daksa sedang (D1) | laras (E) | wicara (F) | hyperaktif (H) | cerdas istimewa (I)",
		"bakat istimewa (J) | kesulitan belajar (K) | narkoba (N) | indigo (O) | down syndrome (P) | autis (Q)",
		"",
		"Kolom 'Jalur Masuk' — pilih salah satu:",
		"PPDB | Prestasi Akademik | Prestasi Non Akademik | Prestasi Akademik Dalam",
		"Prestasi Akademik Luar | Afirmasi | Anak Panti | Inklusi | Domisili | Mutasi | Tahap Kedua | Pindahan",
		"",
		"Kolom 'Semester' — pilih salah satu:",
		"1 | 2 | 3 | 4 | 5 | 6",
		"",
		"Kolom 'Pendidikan Ayah', 'Pendidikan Ibu', 'Pendidikan Wali' — pilih salah satu:",
		"Tidak sekolah | TK / sederajat | SD / sederajat | SMP / sederajat | SMA / sederajat",
		"D1 | D2 | D3 | D4 | S1 | S2 | S2 Terapan | S3 | S3 Terapan | Sp-1 | Sp-2",
		"PAUD | Profesi | Paket A | Paket B | Paket C | Informal | Non formal | Putus SD | Lainnya",
		"",
		"Kolom 'Pekerjaan Ayah', 'Pekerjaan Ibu', 'Pekerjaan Wali' — pilih salah satu:",
		"Tidak bekerja | Nelayan | Petani | Peternak | PNS/TNI/Polri | Karyawan Swasta",
		"Pedagang Kecil | Pedagang Besar | Wiraswasta | Wirausaha | Buruh | Pensiunan",
		"Tenaga Kerja Indonesia | Karyawan BUMN | Tidak dapat diterapkan | Sudah Meninggal | Lainnya",
		"",
		"Kolom 'Kategori Penghasilan Ayah', 'Kategori Penghasilan Ibu', 'Kategori Penghasilan Wali' — pilih salah satu:",
		"Tidak Pasti | Kurang dari Rp. 500,000 | Rp. 500,000 - Rp. 999,999 | Rp. 1,000,000 - Rp. 1,999,999",
		"Rp. 2,000,000 - Rp. 4,999,999 | Rp. 5,000,000 - Rp. 20,000,000 | Lebih dari Rp. 20,000,000 | Tidak Berpenghasilan",
		"",
		"Kolom 'Status Ayah', 'Status Ibu' — pilih salah satu:",
		"Hidup | Meninggal",
	}
	for i, line := range instructions {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetCellValue(petunjuk, cell, line)
	}
	f.SetColWidth(petunjuk, "A", "A", 100)
	f.SetActiveSheet(0)

	filename := fmt.Sprintf("template-peserta-didik-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return f.Write(c.Response().BodyWriter())
}

// ImportHandler memproses file Excel upload, validasi penuh (all-or-nothing)
func ImportHandler(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "File Excel wajib diunggah"})
	}
	src, err := fileHeader.Open()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Gagal membuka file"})
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Format file tidak valid, pastikan file .xlsx"})
	}
	sheet := f.GetSheetList()[0]
	rows, err := f.GetRows(sheet)
	if err != nil || len(rows) < 2 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "File kosong atau tidak ada data untuk diimport"})
	}

	var errorsList []ImportRowError
	var validData []PesertaDidik
	var allNIS, allNISN, allNoUjian, allNIKSiswa, allNIKAyah, allNIKIbu []string
	seenNIS := map[string]int{}
	seenNISN := map[string]int{}
	seenNoUjian := map[string]int{}
	seenNIKSiswa := map[string]int{}
	seenNIKAyah := map[string]int{}
	seenNIKIbu := map[string]int{}

	for i, row := range rows[1:] {
		rowNum := i + 2
		if isRowBlank(row) {
			continue
		}

		idTP, _ := strconv.Atoi(getCell(row, 1))
		idKelas, _ := strconv.Atoi(getCell(row, 2))
		idRombel, _ := strconv.Atoi(getCell(row, 3))

		p := PesertaDidik{
			IDTahunPelajaran: idTP, IDKelas: idKelas, IDRombel: idRombel,
			NamaSiswa: getCell(row, 4), NIS: getCell(row, 5), NISN: getCell(row, 6),
			NoUjianSD: getCell(row, 7), NoSimdik: getCell(row, 8), JenisKelamin: getCell(row, 9),
			TempatLahirSiswa: getCell(row, 10), TanggalLahirSiswa: getCell(row, 11),
			NoAkte: getCell(row, 12), AkteKota: getCell(row, 13), AgamaSiswa: getCell(row, 14), WargaSiswa: getCell(row, 15),
			AnakKe: getCell(row, 16), SaudaraKandung: getCell(row, 17), SaudaraTiri: getCell(row, 18), SaudaraAngkat: getCell(row, 19),
			YatimPiatu: getCell(row, 20), BahasaSiswa: getCell(row, 21), Bantuan: getCell(row, 22), NoKK: getCell(row, 23),
			NIKSiswa: getCell(row, 24), AlamatSiswa: getCell(row, 25), RTSiswa: getCell(row, 26), RWSiswa: getCell(row, 27),
			KoposSiswa: getCell(row, 28), KelSiswa: getCell(row, 29), KecSiswa: getCell(row, 30),
			TelpRmhSiswa: getCell(row, 31), HPSiswa: getCell(row, 32), TinggalSiswa: getCell(row, 33), JarakSekolah: getCell(row, 34),
			GolDarah: getCell(row, 35), PernahSakit: getCell(row, 36), StatusJasmani: getCell(row, 37), KebutuhanKhusus: getCell(row, 38),
			TinggiBadan: getCell(row, 39), BeratBadan: getCell(row, 40), LingkarKepala: getCell(row, 41), ModaTransportasi: getCell(row, 42),
			LulusSD: getCell(row, 43), TanggalIjazah: getCell(row, 44), KotaIjazah: getCell(row, 45), NomorIjazah: getCell(row, 46), TahunIjazah: getCell(row, 47),
			NomorSKHUN: getCell(row, 48), TahunSKHUN: getCell(row, 49), NilaiIjazah: getCell(row, 50), NilaiNEM: getCell(row, 51), LamaBelajar: getCell(row, 52),
			Jalur: getCell(row, 53), Semester: getCell(row, 54), AsalSMP: getCell(row, 55), TanggalMasuk: getCell(row, 56), Email: getCell(row, 57),
			NamaAyah: getCell(row, 58), TempatLahirAyah: getCell(row, 59), TanggalLahirAyah: getCell(row, 60), AgamaAyah: getCell(row, 61), WargaAyah: getCell(row, 62),
			PendAyah: getCell(row, 63), PekAyah: getCell(row, 64), GajiAyah: getCell(row, 65), GajiAyahKtg: getCell(row, 66),
			NIKAyah: getCell(row, 67), AlamatAyah: getCell(row, 68), RTAyah: getCell(row, 69), RWAyah: getCell(row, 70), KoposAyah: getCell(row, 71),
			KelAyah: getCell(row, 72), KecAyah: getCell(row, 73), TelpAyah: getCell(row, 74), StatusAyah: getCell(row, 75),
			NamaIbu: getCell(row, 76), TempatLahirIbu: getCell(row, 77), TanggalLahirIbu: getCell(row, 78), AgamaIbu: getCell(row, 79), WargaIbu: getCell(row, 80),
			PendIbu: getCell(row, 81), PekIbu: getCell(row, 82), GajiIbu: getCell(row, 83), GajiIbuKtg: getCell(row, 84),
			NIKIbu: getCell(row, 85), AlamatIbu: getCell(row, 86), RTIbu: getCell(row, 87), RWIbu: getCell(row, 88), KoposIbu: getCell(row, 89),
			KelIbu: getCell(row, 90), KecIbu: getCell(row, 91), TelpIbu: getCell(row, 92), StatusIbu: getCell(row, 93),
			NamaWali: getCell(row, 94), TempatLahirWali: getCell(row, 95), TanggalLahirWali: getCell(row, 96), AgamaWali: getCell(row, 97), WargaWali: getCell(row, 98),
			PendWali: getCell(row, 99), PekWali: getCell(row, 100), GajiWali: getCell(row, 101), GajiWaliKtg: getCell(row, 102),
			AlamatWali: getCell(row, 103), RTWali: getCell(row, 104), RWWali: getCell(row, 105), KoposWali: getCell(row, 106),
			KelWali: getCell(row, 107), KecWali: getCell(row, 108), TelpWali: getCell(row, 109),
			HobiSeni: getCell(row, 110), HobiOrkes: getCell(row, 111), HobiOrganisasi: getCell(row, 112), HobiLain: getCell(row, 113),
			StatusKeluar: getCell(row, 114), AlasanKeluar: getCell(row, 115), TanggalKeluar: getCell(row, 116), NomorSurat: getCell(row, 117), Tujuan: getCell(row, 118),
			Foto: getCell(row, 119),
		}

		if p.NIS == "" {
			errorsList = append(errorsList, ImportRowError{rowNum, "NIS wajib diisi"})
			continue
		}
		if msg := validateDigitCell(p.NIS, "NIS", 12, false, true); msg != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, msg})
		}
		if msg := validateDigitCell(p.NISN, "NISN", 10, true, false); msg != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, msg})
		}
		if msg := validateDigitCell(p.NoUjianSD, "No. Ujian SD", 22, false, false); msg != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, msg})
		}
		if msg := validateDigitCell(p.NIKSiswa, "NIK Siswa", 16, true, false); msg != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, msg})
		}
		if msg := validateDigitCell(p.NIKAyah, "NIK Ayah", 16, true, false); msg != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, msg})
		}
		if msg := validateDigitCell(p.NIKIbu, "NIK Ibu", 16, true, false); msg != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, msg})
		}
		if msg := validateDigitCell(p.NoKK, "No. KK", 16, true, false); msg != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, msg})
		}
		if p.IDTahunPelajaran == 0 {
			errorsList = append(errorsList, ImportRowError{rowNum, "ID Tahun Pelajaran wajib diisi angka"})
		}
		if p.IDKelas == 0 {
			errorsList = append(errorsList, ImportRowError{rowNum, "ID Kelas wajib diisi angka"})
		}
		if p.IDRombel == 0 {
			errorsList = append(errorsList, ImportRowError{rowNum, "ID Rombel wajib diisi angka"})
		}

		// cek duplikat di dalam file
		if prev, ok := seenNIS[p.NIS]; ok {
			errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("NIS '%s' duplikat dengan baris %d", p.NIS, prev)})
		} else {
			seenNIS[p.NIS] = rowNum
		}
		checkDup := func(value string, seen map[string]int, label string) {
			if value == "" {
				return
			}
			if prev, ok := seen[value]; ok {
				errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("%s '%s' duplikat dengan baris %d", label, value, prev)})
			} else {
				seen[value] = rowNum
			}
		}
		checkDup(p.NISN, seenNISN, "NISN")
		checkDup(p.NoUjianSD, seenNoUjian, "No. Ujian SD")
		checkDup(p.NIKSiswa, seenNIKSiswa, "NIK Siswa")
		checkDup(p.NIKAyah, seenNIKAyah, "NIK Ayah")
		checkDup(p.NIKIbu, seenNIKIbu, "NIK Ibu")

		validData = append(validData, p)
		allNIS = append(allNIS, p.NIS)
		if p.NISN != "" {
			allNISN = append(allNISN, p.NISN)
		}
		if p.NoUjianSD != "" {
			allNoUjian = append(allNoUjian, p.NoUjianSD)
		}
		if p.NIKSiswa != "" {
			allNIKSiswa = append(allNIKSiswa, p.NIKSiswa)
		}
		if p.NIKAyah != "" {
			allNIKAyah = append(allNIKAyah, p.NIKAyah)
		}
		if p.NIKIbu != "" {
			allNIKIbu = append(allNIKIbu, p.NIKIbu)
		}
	}

	// cek bentrok dengan data yang sudah ada di database
	conflicts, err := GetExistingConflicts(allNIS, allNISN, allNoUjian, allNIKSiswa, allNIKAyah, allNIKIbu)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	for _, cf := range conflicts {
		if rowNum, ok := seenNIS[cf.NIS]; ok && cf.NIS != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("NIS '%s' sudah terdaftar di database (a.n. %s)", cf.NIS, cf.NamaSiswa)})
		}
		if rowNum, ok := seenNISN[cf.NISN]; ok && cf.NISN != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("NISN '%s' sudah terdaftar di database (a.n. %s)", cf.NISN, cf.NamaSiswa)})
		}
		if rowNum, ok := seenNoUjian[cf.NoUjianSD]; ok && cf.NoUjianSD != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("No. Ujian SD '%s' sudah terdaftar di database (a.n. %s)", cf.NoUjianSD, cf.NamaSiswa)})
		}
		if rowNum, ok := seenNIKSiswa[cf.NIKSiswa]; ok && cf.NIKSiswa != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("NIK Siswa '%s' sudah terdaftar di database (a.n. %s)", cf.NIKSiswa, cf.NamaSiswa)})
		}
		if rowNum, ok := seenNIKAyah[cf.NIKAyah]; ok && cf.NIKAyah != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("NIK Ayah '%s' sudah terdaftar di database (a.n. %s)", cf.NIKAyah, cf.NamaSiswa)})
		}
		if rowNum, ok := seenNIKIbu[cf.NIKIbu]; ok && cf.NIKIbu != "" {
			errorsList = append(errorsList, ImportRowError{rowNum, fmt.Sprintf("NIK Ibu '%s' sudah terdaftar di database (a.n. %s)", cf.NIKIbu, cf.NamaSiswa)})
		}
	}

	if len(errorsList) > 0 {
		sort.Slice(errorsList, func(i, j int) bool { return errorsList[i].Row < errorsList[j].Row })
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Ditemukan data bermasalah, seluruh proses import dibatalkan",
			"errors":  errorsList,
		})
	}
	if len(validData) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Tidak ada data valid untuk diimport"})
	}

	if err := ImportPesertaDidik(validData); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menyimpan data: " + err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("%d data peserta didik berhasil diimport", len(validData)),
		"data":    fiber.Map{"imported": len(validData)},
	})
}
