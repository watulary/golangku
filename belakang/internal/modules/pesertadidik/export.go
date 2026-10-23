package pesertadidik

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

// pesertaDidikExcelHeaders — urutan HARUS sinkron dengan values di ExportHandler.
// 3 kolom pertama (Tahun Pelajaran/Kelas/Rombel) diisi NAMA hasil JOIN, bukan ID mentah.
var pesertaDidikExcelHeaders = []string{
	"No", "Tahun Pelajaran", "Kelas", "Rombel",
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
	"Foto", "Status",
}

func pesertaDidikExcelHeaderStyle(f *excelize.File) int {
	style, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"1E40AF"}},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "left"},
	})
	if err != nil {
		return -1
	}
	return style
}

// buildPesertaDidikExcelFile membangun workbook Excel dari data peserta didik.
// Dipakai bersama oleh export aktif & export mutasi (kolomnya identik).
func buildPesertaDidikExcelFile(sheet string, data []PesertaDidik) *excelize.File {
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", sheet)

	for i, h := range pesertaDidikExcelHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, row := range data {
		r := i + 2
		statusLabel := "Aktif"
		if row.StatusKeluar != "" {
			statusLabel = "Keluar"
		}
		values := []interface{}{
			i + 1,
			row.NamaTahunPelajaran, row.NamaKelas, row.NamaRombel,
			row.NamaSiswa, row.NIS, row.NISN, row.NoUjianSD, row.NoSimdik, row.JenisKelamin,
			row.TempatLahirSiswa, row.TanggalLahirSiswa, row.NoAkte, row.AkteKota, row.AgamaSiswa, row.WargaSiswa,
			row.AnakKe, row.SaudaraKandung, row.SaudaraTiri, row.SaudaraAngkat, row.YatimPiatu,
			row.BahasaSiswa, row.Bantuan, row.NoKK, row.NIKSiswa,
			row.AlamatSiswa, row.RTSiswa, row.RWSiswa, row.KoposSiswa, row.KelSiswa, row.KecSiswa,
			row.TelpRmhSiswa, row.HPSiswa, row.TinggalSiswa, row.JarakSekolah,
			row.GolDarah, row.PernahSakit, row.StatusJasmani, row.KebutuhanKhusus,
			row.TinggiBadan, row.BeratBadan, row.LingkarKepala, row.ModaTransportasi,
			row.LulusSD, row.TanggalIjazah, row.KotaIjazah, row.NomorIjazah, row.TahunIjazah,
			row.NomorSKHUN, row.TahunSKHUN, row.NilaiIjazah, row.NilaiNEM, row.LamaBelajar,
			row.Jalur, row.Semester, row.AsalSMP, row.TanggalMasuk, row.Email,
			row.NamaAyah, row.TempatLahirAyah, row.TanggalLahirAyah, row.AgamaAyah, row.WargaAyah,
			row.PendAyah, row.PekAyah, row.GajiAyah, row.GajiAyahKtg,
			row.NIKAyah, row.AlamatAyah, row.RTAyah, row.RWAyah, row.KoposAyah,
			row.KelAyah, row.KecAyah, row.TelpAyah, row.StatusAyah,
			row.NamaIbu, row.TempatLahirIbu, row.TanggalLahirIbu, row.AgamaIbu, row.WargaIbu,
			row.PendIbu, row.PekIbu, row.GajiIbu, row.GajiIbuKtg,
			row.NIKIbu, row.AlamatIbu, row.RTIbu, row.RWIbu, row.KoposIbu,
			row.KelIbu, row.KecIbu, row.TelpIbu, row.StatusIbu,
			row.NamaWali, row.TempatLahirWali, row.TanggalLahirWali, row.AgamaWali, row.WargaWali,
			row.PendWali, row.PekWali, row.GajiWali, row.GajiWaliKtg,
			row.AlamatWali, row.RTWali, row.RWWali, row.KoposWali, row.KelWali, row.KecWali, row.TelpWali,
			row.HobiSeni, row.HobiOrkes, row.HobiOrganisasi, row.HobiLain,
			row.StatusKeluar, row.AlasanKeluar, row.TanggalKeluar, row.NomorSurat, row.Tujuan,
			row.Foto, statusLabel,
		}
		for colIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, r)
			f.SetCellValue(sheet, cell, v)
		}
	}

	lastCol, _ := excelize.ColumnNumberToName(len(pesertaDidikExcelHeaders))
	f.SetColWidth(sheet, "A", "A", 6)
	secondCol, _ := excelize.ColumnNumberToName(2)
	f.SetColWidth(sheet, secondCol, lastCol, 18)

	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	// Format seluruh kolom jadi Text (@), supaya NIS/NISN/NIK/dll tidak diotak-atik Excel
	customFmt := "@"
	if textStyle, err := f.NewStyle(&excelize.Style{CustomNumFmt: &customFmt}); err == nil {
		f.SetColStyle(sheet, fmt.Sprintf("A:%s", lastCol), textStyle)
	}

	if style := pesertaDidikExcelHeaderStyle(f); style != -1 {
		firstCell, _ := excelize.CoordinatesToCellName(1, 1)
		lastHeaderCell, _ := excelize.CoordinatesToCellName(len(pesertaDidikExcelHeaders), 1)
		f.SetCellStyle(sheet, firstCell, lastHeaderCell, style)
	}

	return f
}

// ExportHandler mengekspor peserta didik AKTIF, mengikuti filter yang sedang aktif di halaman list.
func ExportHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	idTahunPelajaran, _ := strconv.Atoi(c.Query("id_tahunpelajaran", "0"))
	idKelas, _ := strconv.Atoi(c.Query("id_kelas", "0"))
	idRombel, _ := strconv.Atoi(c.Query("id_rombel", "0"))

	data, err := GetAllForExport(search, idTahunPelajaran, idKelas, idRombel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	f := buildPesertaDidikExcelFile("Peserta Didik", data)

	filename := fmt.Sprintf("data-peserta-didik-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	return f.Write(c.Response().BodyWriter())
}

// ExportMutasiHandler mengekspor SEMUA data mutasi sekaligus, TANPA filter (search/tahun/kelas/rombel).
func ExportMutasiHandler(c *fiber.Ctx) error {
	data, err := GetAllMutasiForExport()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	f := buildPesertaDidikExcelFile("Mutasi", data)

	filename := fmt.Sprintf("daftar-mutasi-peserta-didik-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	return f.Write(c.Response().BodyWriter())
}

// ExportArsipHandler mengekspor data arsip, mengikuti filter yang sedang aktif
// di halaman Arsip (search + tahun pelajaran + kelas + rombel).
func ExportArsipHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	idTahunPelajaran, _ := strconv.Atoi(c.Query("id_tahunpelajaran", "0"))
	idKelas, _ := strconv.Atoi(c.Query("id_kelas", "0"))
	idRombel, _ := strconv.Atoi(c.Query("id_rombel", "0"))

	data, err := GetAllArsipForExport(search, idTahunPelajaran, idKelas, idRombel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	f := buildPesertaDidikExcelFile("Arsip", data)

	filename := fmt.Sprintf("daftar-arsip-peserta-didik-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	return f.Write(c.Response().BodyWriter())
}
