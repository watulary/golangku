package kelulusan

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"

	"sekolah/internal/auth"
)

const fotoDir = "./storage/fotokelulusan"

// defaultFotoFile - foto fallback kalau siswa belum punya fotolulus. Tidak
// pernah boleh terhapus lewat BulkDeleteFotoHandler (lihat guard di sana).
const defaultFotoFile = "peserta_didik.jpg"

var allowedFotoExt = map[string]bool{".jpg": true, ".jpeg": true, ".png": true}

// ---------- ADMIN: list kelulusan ----------

func ListAdminHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	sortBy := c.Query("sort_by", "nama_siswa")
	sortDir := c.Query("sort_dir", "asc")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	rows, err := GetAllAdmin(search, sortBy, sortDir, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	for i := range rows {
		if rows[i].FotoLulus != "" {
			rows[i].FotoURL = "/storage/fotokelulusan/" + rows[i].FotoLulus
		} else {
			rows[i].FotoURL = "/storage/fotokelulusan/" + defaultFotoFile // foto default, sama seperti index.php lama
		}
	}
	total, err := CountAdmin(search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	totalPage := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"success": true,
		"data":    rows,
		"pagination": fiber.Map{
			"page": page, "limit": limit, "total": total, "total_page": totalPage,
		},
	})
}

func GetOneAdminHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")
	row, err := GetByNISAdmin(nis)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data tidak ditemukan."})
	}
	if row.FotoLulus != "" {
		row.FotoURL = "/storage/fotokelulusan/" + row.FotoLulus
	} else {
		row.FotoURL = "/storage/fotokelulusan/" + defaultFotoFile
	}
	return c.JSON(fiber.Map{"success": true, "data": row})
}

// ---------- ADMIN: pengaturan pengumuman & isi SKL (kelulusan_waktu) ----------

func GetWaktuHandler(c *fiber.Ctx) error {
	w, err := GetWaktu()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": w})
}

var pengumumanFields = []string{
	"aturwaktu", "aturform", "heroarea1", "heroarea2", "judul", "nomor",
	"pengantar", "penutup", "akhir", "tanggal", "mengetahui", "kepsek",
	"nipkepsek", "tahun_pelajaran",
}

func UpdatePengumumanHandler(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*auth.Claims)

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid."})
	}

	old, err := GetWaktu()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	changed := diffFields(waktuToMap(old), body, pengumumanFields)
	if len(changed) == 0 {
		return c.JSON(fiber.Map{"success": true, "message": "Tidak ada perubahan data.", "changed": []string{}})
	}

	data := map[string]interface{}{"id_pengguna": claims.IDPengguna}
	for _, f := range pengumumanFields {
		data[f] = body[f]
	}
	if err := UpdatePengumuman(data); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Ubah data berhasil.", "changed": changed})
}

var sklFields = []string{
	"kop1_skl", "kop2_skl", "kop3_skl", "kop4_skl", "kop5_skl", "kop6_skl", "kop7_skl", "kop8_skl",
	"isi1_skl", "isi2_skl", "isi3_skl", "isi4_skl", "isi5_skl",
	"tgl_skl", "sekolah_skl", "kepsek_skl", "nipkepsek_skl", "ttd_skl",
}

func UpdateSKLHandler(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*auth.Claims)

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid."})
	}

	old, err := GetWaktu()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	changed := diffFields(waktuToMap(old), body, sklFields)
	if len(changed) == 0 {
		return c.JSON(fiber.Map{"success": true, "message": "Tidak ada perubahan data.", "changed": []string{}})
	}

	data := map[string]interface{}{"id_pengguna": claims.IDPengguna}
	for _, f := range sklFields {
		data[f] = body[f]
	}
	if err := UpdateSKL(data); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Ubah data berhasil.", "changed": changed})
}

func waktuToMap(w *KelulusanWaktu) map[string]interface{} {
	return map[string]interface{}{
		"aturwaktu": w.AturWaktu, "aturform": w.AturForm, "heroarea1": w.HeroArea1, "heroarea2": w.HeroArea2,
		"judul": w.Judul, "nomor": w.Nomor, "pengantar": w.Pengantar, "penutup": w.Penutup, "akhir": w.Akhir,
		"tanggal": w.Tanggal, "mengetahui": w.Mengetahui, "kepsek": w.Kepsek, "nipkepsek": w.NipKepsek,
		"tahun_pelajaran": w.TahunPelajaran,
		"kop1_skl":        w.Kop1SKL, "kop2_skl": w.Kop2SKL, "kop3_skl": w.Kop3SKL, "kop4_skl": w.Kop4SKL,
		"kop5_skl": w.Kop5SKL, "kop6_skl": w.Kop6SKL, "kop7_skl": w.Kop7SKL, "kop8_skl": w.Kop8SKL,
		"isi1_skl": w.Isi1SKL, "isi2_skl": w.Isi2SKL, "isi3_skl": w.Isi3SKL, "isi4_skl": w.Isi4SKL, "isi5_skl": w.Isi5SKL,
		"tgl_skl": w.TglSKL, "sekolah_skl": w.SekolahSKL, "kepsek_skl": w.KepsekSKL,
		"nipkepsek_skl": w.NipKepsekSKL, "ttd_skl": w.TtdSKL,
	}
}

// diffFields - bandingkan old vs body hanya utk fields yg diberikan; pola
// "cek changed dulu" dipakai di semua UpdateHandler modul lain.
func diffFields(old map[string]interface{}, body map[string]interface{}, fields []string) []string {
	var changed []string
	for _, f := range fields {
		oldVal := fmt.Sprintf("%v", old[f])
		newVal := fmt.Sprintf("%v", body[f])
		if oldVal != newVal {
			changed = append(changed, f)
		}
	}
	return changed
}

// ---------- ADMIN: CRUD manual per-siswa ----------

var kelulusanFields = []string{
	"id_tahunpelajaran", "id_kelas", "id_rombel", "nopes", "pa_s", "ppkn_s", "bind_s", "bing_s",
	"mtk_s", "ipa_s", "ips_s", "penjas_s", "senbud_s", "prakarya_s", "jumlah", "rata_s", "ket", "fotolulus",
}

func CreateHandler(c *fiber.Ctx) error {
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid."})
	}
	nis, _ := body["nis"].(string)
	if nis == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "NIS wajib diisi."})
	}
	ok, err := ExistsPesertaDidik(nis)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if !ok {
		return c.Status(422).JSON(fiber.Map{"success": false, "message": "NIS tidak terdaftar di Peserta Didik."})
	}
	dup, err := ExistsNIS(nis)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if dup {
		return c.Status(422).JSON(fiber.Map{"success": false, "message": "NIS sudah punya data kelulusan."})
	}

	data := map[string]interface{}{"nis": nis}
	for _, f := range kelulusanFields {
		data[f] = body[f]
	}
	if err := Create(data); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Tambah data berhasil."})
}

func UpdateHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid."})
	}

	old, err := GetByNIS(nis)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data tidak ditemukan."})
	}
	oldMap := map[string]interface{}{
		"id_tahunpelajaran": old.IDTahunPelajaran, "id_kelas": old.IDKelas, "id_rombel": old.IDRombel,
		"nopes": old.Nopes, "pa_s": old.PaS, "ppkn_s": old.PpknS, "bind_s": old.BindS, "bing_s": old.BingS,
		"mtk_s": old.MtkS, "ipa_s": old.IpaS, "ips_s": old.IpsS, "penjas_s": old.PenjasS,
		"senbud_s": old.SenbudS, "prakarya_s": old.PrakaryaS, "jumlah": old.Jumlah, "rata_s": old.RataS,
		"ket": old.Ket, "fotolulus": old.FotoLulus,
	}
	changed := diffFields(oldMap, body, kelulusanFields)
	if len(changed) == 0 {
		return c.JSON(fiber.Map{"success": true, "message": "Tidak ada perubahan data.", "changed": []string{}})
	}

	data := map[string]interface{}{}
	for _, f := range kelulusanFields {
		data[f] = body[f]
	}
	if err := Update(nis, data); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Ubah data berhasil.", "changed": changed})
}

func DeleteHandler(c *fiber.Ctx) error {
	nis := c.Params("nis")
	if err := Delete(nis); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Hapus data berhasil."})
}

func BulkDeleteHandler(c *fiber.Ctx) error {
	var body struct {
		NIS []string `json:"nis"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid."})
	}
	if err := BulkDelete(body.NIS); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": fmt.Sprintf("%d data berhasil dihapus.", len(body.NIS))})
}

// ---------- ADMIN: import excel ----------

var excelCols = []string{
	"id_tahunpelajaran", "id_kelas", "id_rombel", "nis", "nopes", "pa_s", "ppkn_s", "bind_s", "bing_s",
	"mtk_s", "ipa_s", "ips_s", "penjas_s", "senbud_s", "prakarya_s", "jumlah", "rata_s", "ket", "fotolulus",
}

func ImportExcelHandler(c *fiber.Ctx) error {
	type errItem struct {
		Row     int    `json:"row"`
		Message string `json:"message"`
	}

	fh, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "File excel wajib dipilih."})
	}
	f, err := fh.Open()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "File tidak bisa dibuka."})
	}
	defer f.Close()

	xf, err := excelize.OpenReader(f)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Format file excel tidak valid."})
	}
	sheet := xf.GetSheetName(0)
	rowsData, err := xf.GetRows(sheet)
	if err != nil || len(rowsData) <= 1 {
		return c.Status(422).JSON(fiber.Map{"success": false, "message": "File excel kosong."})
	}

	var errs []errItem
	nisCheck := map[string]bool{}
	nopesCheck := map[string]bool{}
	fotoCheck := map[string]bool{}
	type rowRec struct {
		nis  string
		data map[string]interface{}
	}
	var valid []rowRec

	label := map[string]string{
		"id_tahunpelajaran": "ID Tahun Pelajaran", "id_kelas": "ID Kelas", "id_rombel": "ID Rombel",
		"nis": "NIS", "nopes": "No Ujian Peserta", "pa_s": "Nilai PA", "ppkn_s": "Nilai PPKn",
		"bind_s": "Nilai Bind", "bing_s": "Nilai Bing", "mtk_s": "Nilai MTK", "ipa_s": "Nilai IPA",
		"ips_s": "Nilai IPS", "penjas_s": "Nilai Penjas", "senbud_s": "Nilai Senbud",
		"prakarya_s": "Nilai Prakarya", "jumlah": "Nilai Jumlah", "rata_s": "Nilai Rerata",
		"ket": "Keterangan", "fotolulus": "Foto",
	}

	for i, row := range rowsData {
		if i == 0 {
			continue
		}
		vals := map[string]string{}
		allEmpty := true
		for idx, col := range excelCols {
			v := ""
			if idx < len(row) {
				v = strings.TrimSpace(row[idx])
			}
			vals[col] = v
			if v != "" {
				allEmpty = false
			}
		}
		if allEmpty {
			continue
		}

		lineNo := i + 1
		for _, col := range excelCols {
			if col == "fotolulus" {
				continue // opsional -- kosong berarti pakai foto default (peserta_didik.jpg) saat ditampilkan
			}
			if vals[col] == "" {
				errs = append(errs, errItem{lineNo, label[col] + " kosong."})
			}
		}
		if vals["nis"] != "" {
			if nisCheck[vals["nis"]] {
				errs = append(errs, errItem{lineNo, "NIS (" + vals["nis"] + ") terdapat data ganda."})
			}
			nisCheck[vals["nis"]] = true
		}
		if vals["nopes"] != "" {
			if nopesCheck[vals["nopes"]] {
				errs = append(errs, errItem{lineNo, "No Peserta (" + vals["nopes"] + ") terdapat data ganda."})
			}
			nopesCheck[vals["nopes"]] = true
		}
		// fotolulus opsional; nilai kosong ATAU sama dengan foto default
		// (peserta_didik.jpg) TIDAK divalidasi unik -- keduanya berarti
		// "belum ditentukan fotonya", wajar dipakai berulang oleh banyak
		// siswa. Nama file foto ASLI (hasil upload ZIP) tetap wajib unik.
		if vals["fotolulus"] != "" && vals["fotolulus"] != defaultFotoFile {
			if fotoCheck[vals["fotolulus"]] {
				errs = append(errs, errItem{lineNo, "Foto (" + vals["fotolulus"] + ") terdapat data ganda."})
			}
			fotoCheck[vals["fotolulus"]] = true
		}
		if vals["nis"] != "" {
			if ok, _ := ExistsNIS(vals["nis"]); ok {
				errs = append(errs, errItem{lineNo, "NIS (" + vals["nis"] + ") sudah terdaftar di database."})
			}
			if ok, _ := ExistsPesertaDidik(vals["nis"]); !ok {
				errs = append(errs, errItem{lineNo, "NIS (" + vals["nis"] + ") tidak terdaftar di Peserta Didik."})
			}
		}
		if vals["nopes"] != "" {
			if ok, _ := ExistsNopes(vals["nopes"]); ok {
				errs = append(errs, errItem{lineNo, "No Peserta (" + vals["nopes"] + ") sudah terdaftar di database."})
			}
		}
		if vals["fotolulus"] != "" && vals["fotolulus"] != defaultFotoFile {
			if ok, _ := ExistsFotoLulus(vals["fotolulus"]); ok {
				errs = append(errs, errItem{lineNo, "Foto (" + vals["fotolulus"] + ") sudah terdaftar di database."})
			}
		}

		data := map[string]interface{}{}
		for _, col := range excelCols {
			data[col] = vals[col]
		}
		valid = append(valid, rowRec{nis: vals["nis"], data: data})
	}

	if len(errs) > 0 {
		return c.Status(422).JSON(fiber.Map{"success": false, "message": "Impor gagal, periksa kembali data.", "errors": errs})
	}
	if len(valid) == 0 {
		return c.Status(422).JSON(fiber.Map{"success": false, "message": "File excel kosong."})
	}

	for _, r := range valid {
		if err := Create(r.data); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
		}
	}

	return c.JSON(fiber.Map{"success": true, "message": fmt.Sprintf("%d Peserta didik berhasil diunggah.", len(valid))})
}

func TemplateHandler(c *fiber.Ctx) error {
	xf := excelize.NewFile()
	defer xf.Close()
	sheet := "Kelulusan"
	xf.SetSheetName("Sheet1", sheet)

	// Header dibuat human-readable (bukan nama kolom DB mentah) -- ImportExcelHandler
	// membaca berdasarkan POSISI kolom, bukan nama header, jadi ini murni kosmetik,
	// sama seperti pola label di template Kenaikan.
	header := []string{
		"ID Tahun Pelajaran", "ID Kelas", "ID Rombel", "NIS", "No Peserta Ujian",
		"Nilai PA", "Nilai PPKn", "Nilai B.Indo", "Nilai B.Inggris", "Nilai MTK",
		"Nilai IPA", "Nilai IPS", "Nilai Penjas", "Nilai Senbud", "Nilai Prakarya",
		"Jumlah Nilai", "Rerata Nilai", "Keterangan", "Nama File Foto",
	}

	// Semua kolom dipaksa format Text (NumFmt 49 = "@") supaya NIS/No Peserta
	// TIDAK otomatis dikonversi jadi angka saat diisi -- nol di depan (mis.
	// NIS "012749") akan hilang kalau kolomnya General/Number.
	textStyle, errStyle := xf.NewStyle(&excelize.Style{NumFmt: 49})
	lastCell, _ := excelize.CoordinatesToCellName(len(header), 1)
	lastCol := lastCell[:len(lastCell)-1]
	if errStyle == nil {
		xf.SetColStyle(sheet, "A:"+lastCol, textStyle)
	}

	headerStyle, _ := xf.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"2563EB"}, Pattern: 1},
	})
	for i, h := range header {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		xf.SetCellValue(sheet, cell, h)
	}
	xf.SetCellStyle(sheet, "A1", lastCell, headerStyle)
	xf.SetPanes(sheet, &excelize.Panes{Freeze: true, Split: false, XSplit: 0, YSplit: 1, TopLeftCell: "A2", ActivePane: "bottomLeft"})

	// Sheet Petunjuk -- sama pola dgn template Kenaikan.
	xf.NewSheet("Petunjuk")
	xf.SetCellValue("Petunjuk", "A1", "Aturan pengisian:")
	xf.SetCellValue("Petunjuk", "A2", "1. ID Tahun Pelajaran, ID Kelas, ID Rombel harus sesuai data referensi yang berlaku (lihat daftar di bawah).")
	xf.SetCellValue("Petunjuk", "A3", "2. NIS harus sudah terdaftar di data Peserta Didik.")
	xf.SetCellValue("Petunjuk", "A4", "3. NIS tidak boleh duplikat, dan tidak boleh sudah ada di data Kelulusan.")
	xf.SetCellValue("Petunjuk", "A5", "4. No Peserta Ujian tidak boleh duplikat, dan tidak boleh sudah ada di data Kelulusan.")
	xf.SetCellValue("Petunjuk", "A6", "5. Semua nilai wajib diisi (format angka, contoh: 83.98).")
	xf.SetCellValue("Petunjuk", "A7", "6. Keterangan wajib diisi (contoh: Lulus, Tidak Lulus).")
	xf.SetCellValue("Petunjuk", "A8", "7. Nama File Foto bersifat OPSIONAL -- boleh dikosongkan atau diisi \"peserta_didik.jpg\" (keduanya berarti \"belum ada foto\", boleh berulang di banyak baris). Kalau diisi nama file lain, harus sudah ada di Galeri Foto Kelulusan dan tidak boleh duplikat (unggah foto dulu lewat halaman Galeri Foto sebelum import).")

	// Daftar referensi ditulis DINAMIS (bukan cuma label statis) -- query
	// data aktif saat ini, satu baris per item, di bawah masing-masing label.
	row := 10
	writeRefList := func(label string, items []RefItem, fetchErr error) {
		xf.SetCellValue("Petunjuk", fmt.Sprintf("A%d", row), label)
		row++
		if fetchErr != nil || len(items) == 0 {
			xf.SetCellValue("Petunjuk", fmt.Sprintf("A%d", row), "(gagal memuat / belum ada data aktif)")
			row++
		} else {
			for _, it := range items {
				xf.SetCellValue("Petunjuk", fmt.Sprintf("A%d", row), fmt.Sprintf("%d - %s", it.ID, it.Nama))
				row++
			}
		}
		row++ // baris kosong pemisah antar grup referensi
	}

	tahunList, errTP := GetActiveTahunPelajaranRef()
	writeRefList("Referensi Tahun Pelajaran (id - nama):", tahunList, errTP)

	kelasList, errKelas := GetActiveKelasRef()
	writeRefList("Referensi Kelas (id - nama):", kelasList, errKelas)

	rombelList, errRombel := GetActiveRombelRef()
	writeRefList("Referensi Rombel (id - nama):", rombelList, errRombel)

	activeIdx, _ := xf.GetSheetIndex(sheet)
	xf.SetActiveSheet(activeIdx)

	buf, err := xf.WriteToBuffer()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	filename := fmt.Sprintf("template-kelulusan-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	return c.Send(buf.Bytes())
}

func ExportHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	rows, err := GetAllForExport(search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	xf := excelize.NewFile()
	sheet := "Sheet1"
	header := []string{"NIS", "Nama Siswa", "Kelas", "Rombel", "Tahun Pelajaran", "No Peserta",
		"PA", "PPKn", "B.Indo", "B.Ing", "MTK", "IPA", "IPS", "Penjas", "Senbud", "Prakarya",
		"Jumlah", "Rerata", "Keterangan"}
	for i, h := range header {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		xf.SetCellValue(sheet, cell, h)
	}
	for r, row := range rows {
		vals := []interface{}{row.NIS, row.NamaSiswa, row.NamaKelas, row.NamaRombel, row.NamaTahunPelajaran,
			row.Nopes, row.PaS, row.PpknS, row.BindS, row.BingS, row.MtkS, row.IpaS, row.IpsS,
			row.PenjasS, row.SenbudS, row.PrakaryaS, row.Jumlah, row.RataS, row.Ket}
		for i, v := range vals {
			cell, _ := excelize.CoordinatesToCellName(i+1, r+2)
			xf.SetCellValue(sheet, cell, v)
		}
	}

	// Semua kolom dipaksa format Text -- sama alasannya dgn TemplateHandler:
	// NIS/No Peserta yg berawalan nol harus tetap utuh saat file dibuka.
	textStyle, err := xf.NewStyle(&excelize.Style{NumFmt: 49})
	if err == nil {
		lastCol, _ := excelize.CoordinatesToCellName(len(header), 1)
		colRange := lastCol[:len(lastCol)-1]
		xf.SetColStyle(sheet, "A:"+colRange, textStyle)
	}

	buf, err := xf.WriteToBuffer()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	filename := fmt.Sprintf("kelulusan-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	return c.Send(buf.Bytes())
}

// ---------- ADMIN: galeri foto kelulusan (upload ZIP massal) ----------

func ListFotoHandler(c *fiber.Ctx) error {
	rows, err := ListFoto()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": rows})
}

func UploadFotoZipHandler(c *fiber.Ctx) error {
	fh, err := c.FormFile("gambar_zip")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "File zip wajib dipilih."})
	}
	if !strings.HasSuffix(strings.ToLower(fh.Filename), ".zip") {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "File harus berformat zip."})
	}
	const maxZip = 10 << 20 // 10 MB, sama seperti batas lama
	if fh.Size > maxZip {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Ukuran zip maksimal 10 MB."})
	}

	if err := os.MkdirAll(fotoDir, 0o755); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("kelulusan-%d.zip", time.Now().UnixNano()))
	if err := c.SaveFile(fh, tmpPath); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	defer os.Remove(tmpPath)

	zr, err := zip.OpenReader(tmpPath)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ZIP gagal dibuka."})
	}
	defer zr.Close()

	var saved []string
	for _, entry := range zr.File {
		if entry.FileInfo().IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name))
		if !allowedFotoExt[ext] {
			continue
		}
		safeName := filepath.Base(entry.Name)

		rc, err := entry.Open()
		if err != nil {
			continue
		}
		dst, err := os.Create(filepath.Join(fotoDir, safeName))
		if err != nil {
			rc.Close()
			continue
		}
		io.Copy(dst, rc)
		dst.Close()
		rc.Close()

		saved = append(saved, safeName)
	}

	if len(saved) == 0 {
		return c.Status(422).JSON(fiber.Map{"success": false, "message": "Tidak ada file jpg, jpeg, atau png di dalam zip."})
	}
	if err := InsertBatchFoto(saved); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": fmt.Sprintf("%d foto berhasil diunggah.", len(saved))})
}

func BulkDeleteFotoHandler(c *fiber.Ctx) error {
	var body struct {
		ID []int `json:"id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid."})
	}
	if len(body.ID) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Foto belum dipilih."})
	}

	rows, err := GetFotoByIDs(body.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	// Guard lapis kedua (frontend sudah menyembunyikan checkbox-nya juga):
	// foto default TIDAK PERNAH ikut dihapus, baik file fisik maupun baris DB,
	// meski ID-nya somehow terkirim di request.
	var deletableIDs []int
	skippedDefault := false
	for _, row := range rows {
		if row.Gambar == defaultFotoFile {
			skippedDefault = true
			continue
		}
		os.Remove(filepath.Join(fotoDir, row.Gambar))
		deletableIDs = append(deletableIDs, row.ID)
	}

	if err := DeleteFotoRows(deletableIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	message := fmt.Sprintf("%d foto berhasil dihapus.", len(deletableIDs))
	if skippedDefault {
		message += " Foto default (peserta_didik.jpg) dilewati dan tidak dihapus."
	}
	return c.JSON(fiber.Map{"success": true, "message": message})
}

// ---------- PESERTA DIDIK: data sendiri ----------
// Pola RESMI: reuse fungsi admin (GetByNISAdmin) dipanggil pakai claims.NIS.
// Kalau belum ada data, kembalikan null (bukan 404) biar frontend tampilkan "nihil".

func GetMeHandler(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*auth.Claims)

	waktu, err := GetWaktu()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	row, err := GetByNISAdmin(claims.NIS)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    fiber.Map{"waktu": waktu, "kelulusan": nil},
		})
	}
	if row.FotoLulus != "" {
		row.FotoURL = "/storage/fotokelulusan/" + row.FotoLulus
	} else {
		row.FotoURL = "/storage/fotokelulusan/" + defaultFotoFile
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"waktu": waktu, "kelulusan": row},
	})
}

// ---------- Cetak SKL (PDF) — Peserta Didik (punya sendiri) + Admin/Pegawai (siapa saja) ----------

func CetakPDFHandler(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*auth.Claims)
	nis := c.Params("nis")

	isAdmin := strings.EqualFold(claims.Level, "Admin")
	isPegawai := strings.EqualFold(claims.Level, "Pegawai")
	if !isAdmin && !isPegawai && claims.NIS != nis {
		return c.Status(403).JSON(fiber.Map{"success": false, "message": "Anda tidak berhak mencetak SKL siswa lain."})
	}

	cetak, err := GetCetakData(nis)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data kelulusan tidak ditemukan."})
	}
	skl, err := GetWaktu()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	pdfBytes, err := GenerateSKL(cetak, skl)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s.pdf"`, cetak.NamaSiswa))
	return c.Send(pdfBytes)
}
