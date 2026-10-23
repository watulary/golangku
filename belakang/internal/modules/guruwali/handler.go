package guruwali

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"sekolah/internal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

// ================= GURU WALI =================

func GetAllAdminHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	aktif := c.Query("aktif", "")
	sortBy := c.Query("sort_by", "id")
	sortDir := c.Query("sort_dir", "asc")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	items, err := GetAllAdmin(search, aktif, sortBy, sortDir, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data guru wali"})
	}
	total, err := CountAdmin(search, aktif)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghitung data guru wali"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    items,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func GetByIDHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}
	g, err := GetByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data guru wali"})
	}
	if g == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data guru wali tidak ditemukan"})
	}
	namaPegawai, _ := GetNamaPegawai(g.NUPTK) // gagal ambil nama tidak menggagalkan response utama
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id_guruwali":   g.IDGuruWali,
			"nama_guruwali": g.NamaGuruWali,
			"nuptk":         g.NUPTK,
			"aktif":         g.Aktif,
			"nama_pegawai":  namaPegawai,
		},
	})
}

// GetMeHandler mengembalikan guru wali yang membawahi siswa yang sedang login.
// Kalau belum terdaftar di guru wali manapun, data dikembalikan kosong (bukan 404/error).
func GetMeHandler(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Tidak terautentikasi.",
		})
	}

	data, err := GetGuruWaliBySiswa(claims.NIS)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func CreateHandler(c *fiber.Ctx) error {
	var data GuruWali
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid"})
	}

	data.NamaGuruWali = strings.TrimSpace(data.NamaGuruWali)
	data.NUPTK = strings.TrimSpace(data.NUPTK)

	if data.NamaGuruWali == "" || data.NUPTK == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Nama dan NUPTK wajib diisi"})
	}

	namaTaken, err := ExistsByNama(data.NamaGuruWali, 0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal memeriksa nama guru wali"})
	}
	if namaTaken {
		return c.Status(409).JSON(fiber.Map{"success": false, "message": "Nama guru wali sudah digunakan"})
	}

	nuptkTaken, err := ExistsByNUPTK(data.NUPTK, 0)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal memeriksa data guru"})
	}
	if nuptkTaken {
		return c.Status(409).JSON(fiber.Map{"success": false, "message": "Guru ini sudah menjabat sebagai guru wali lain"})
	}

	id, err := Create(data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menyimpan data guru wali"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Data guru wali berhasil ditambahkan", "data": fiber.Map{"id_guruwali": id}})
}

func UpdateHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}
	old, err := GetByID(id)
	if err != nil || old == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data guru wali tidak ditemukan"})
	}
	var data GuruWali
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid"})
	}

	data.NamaGuruWali = strings.TrimSpace(data.NamaGuruWali)
	data.NUPTK = strings.TrimSpace(data.NUPTK)

	if data.NamaGuruWali == "" || data.NUPTK == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Nama dan NUPTK wajib diisi"})
	}

	var changed []string
	if old.NamaGuruWali != data.NamaGuruWali {
		changed = append(changed, "nama_guruwali")
	}
	if old.NUPTK != data.NUPTK {
		changed = append(changed, "nuptk")
	}
	if old.Aktif != data.Aktif {
		changed = append(changed, "aktif")
	}

	if len(changed) == 0 {
		return c.JSON(fiber.Map{"success": true, "message": "Tidak ada perubahan data.", "changed": []string{}})
	}

	// Nama cuma dicek ulang kalau memang berubah — hindari nolak update yang namanya tetap sama
	// (excludeID=id supaya baris yang sedang diedit tidak dianggap bentrok dgn dirinya sendiri).
	if old.NamaGuruWali != data.NamaGuruWali {
		namaTaken, err := ExistsByNama(data.NamaGuruWali, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal memeriksa nama guru wali"})
		}
		if namaTaken {
			return c.Status(409).JSON(fiber.Map{"success": false, "message": "Nama guru wali sudah digunakan"})
		}
	}

	// NUPTK juga cuma dicek ulang kalau memang berubah, pola sama seperti cek nama di atas.
	if old.NUPTK != data.NUPTK {
		nuptkTaken, err := ExistsByNUPTK(data.NUPTK, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal memeriksa data guru"})
		}
		if nuptkTaken {
			return c.Status(409).JSON(fiber.Map{"success": false, "message": "Guru ini sudah menjabat sebagai guru wali lain"})
		}
	}

	if err := Update(id, data); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal memperbarui data guru wali"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Data guru wali berhasil diperbarui", "changed": changed})
}

func DeleteHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}
	if err := Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghapus data guru wali"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Data guru wali berhasil dihapus"})
}

type bulkDeleteRequest struct {
	IDs []int `json:"ids"`
}

func BulkDeleteHandler(c *fiber.Ctx) error {
	var body bulkDeleteRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid"})
	}
	if err := BulkDelete(body.IDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghapus data guru wali"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Data guru wali terpilih berhasil dihapus"})
}

// SearchPegawaiHandler dipakai autocomplete NUPTK di Form.vue
func SearchPegawaiHandler(c *fiber.Ctx) error {
	q := c.Query("q", "")
	items, err := SearchPegawai(q)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mencari pegawai"})
	}
	return c.JSON(fiber.Map{"success": true, "data": items})
}

// ================= ANGGOTA =================

func GetAnggotaHandler(c *fiber.Ctx) error {
	idGuruWali, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}
	search := c.Query("search", "")
	sortBy := c.Query("sort_by", "rombel")
	sortDir := c.Query("sort_dir", "asc")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	idKelas, _ := strconv.Atoi(c.Query("id_kelas", "0"))
	idRombel, _ := strconv.Atoi(c.Query("id_rombel", "0"))
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	items, err := GetAnggotaByGuruWali(idGuruWali, search, sortBy, sortDir, idKelas, idRombel, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data anggota"})
	}
	total, err := CountAnggotaByGuruWali(idGuruWali, search, idKelas, idRombel)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghitung data anggota"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    items,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

type addAnggotaRequest struct {
	NIS string `json:"nis"`
}

func AddAnggotaHandler(c *fiber.Ctx) error {
	idGuruWali, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}
	var body addAnggotaRequest
	if err := c.BodyParser(&body); err != nil || body.NIS == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "NIS wajib diisi"})
	}

	exists, err := ExistsPesertaDidik(body.NIS)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal memeriksa data peserta didik"})
	}
	if !exists {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "NIS tidak terdaftar sebagai peserta didik"})
	}

	dup, err := IsAnggotaExists(idGuruWali, body.NIS)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal memeriksa data anggota"})
	}
	if dup {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Peserta didik ini sudah terdaftar sebagai anggota wali"})
	}

	if _, err := AddAnggota(idGuruWali, body.NIS); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menambahkan anggota"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Peserta didik berhasil ditambahkan sebagai anggota wali"})
}

func RemoveAnggotaHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("anggotaId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}
	if err := RemoveAnggota(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghapus anggota"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Anggota berhasil dihapus"})
}

func BulkRemoveAnggotaHandler(c *fiber.Ctx) error {
	var body bulkDeleteRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Data tidak valid"})
	}
	if err := BulkRemoveAnggota(body.IDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghapus anggota terpilih"})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Anggota terpilih berhasil dihapus"})
}

// SearchPesertaDidikHandler dipakai autocomplete "Tambah Anggota" di Anggota.vue
func SearchPesertaDidikHandler(c *fiber.Ctx) error {
	q := c.Query("q", "")
	items, err := SearchPesertaDidik(q)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mencari peserta didik"})
	}
	return c.JSON(fiber.Map{"success": true, "data": items})
}

// =========================================
// Export Excel: Daftar Guru Wali
// =========================================

func ExportGuruWaliHandler(c *fiber.Ctx) error {
	data, err := GetAllForExport("")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	f := excelize.NewFile()
	sheet := "Guru Wali"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"No", "Nama Guru Wali", "Guru", "NUPTK", "Status", "Jumlah Anggota"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, row := range data {
		r := i + 2

		status := "Nonaktif"
		if row.Aktif {
			status = "Aktif"
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), row.NamaGuruWali)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", r), row.NamaPegawai)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", r), row.NUPTK)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", r), status)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", r), row.JumlahAnggota)
	}

	f.SetColWidth(sheet, "A", "A", 6)
	f.SetColWidth(sheet, "B", "B", 30)
	f.SetColWidth(sheet, "C", "C", 30)
	f.SetColWidth(sheet, "D", "D", 22)
	f.SetColWidth(sheet, "E", "E", 12)
	f.SetColWidth(sheet, "F", "F", 16)

	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"1E40AF"}},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "left",
		},
	})
	if err == nil {
		f.SetCellStyle(sheet, "A1", "F1", headerStyle)
	}

	filename := fmt.Sprintf("daftar-guru-wali-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return f.Write(c.Response().BodyWriter())
}

// =========================================
// Export Excel: Daftar Anggota (Semua Guru Wali, sheet per guru wali)
// =========================================

// sanitizeSheetName membuang karakter yang tidak diperbolehkan Excel pada nama sheet
// dan membatasi panjangnya maksimal 31 karakter.
func sanitizeSheetName(name string) string {
	replacer := strings.NewReplacer(
		"\\", "-",
		"/", "-",
		"?", "",
		"*", "",
		"[", "(",
		"]", ")",
		":", "-",
	)

	clean := replacer.Replace(name)
	clean = strings.TrimSpace(clean)

	if len(clean) > 31 {
		clean = clean[:31]
	}

	if clean == "" {
		clean = "Guru Wali"
	}

	return clean
}

func ExportAnggotaHandler(c *fiber.Ctx) error {
	guruWaliList, err := GetAllForExport("")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	if len(guruWaliList) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Belum ada data guru wali untuk diexport."})
	}

	anggotaList, err := GetAllAnggotaForExport()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	// Kelompokkan anggota per id_guruwali
	anggotaByGuruWali := make(map[int][]AnggotaGuruWaliListItem)
	for _, a := range anggotaList {
		anggotaByGuruWali[a.IDGuruWali] = append(anggotaByGuruWali[a.IDGuruWali], a)
	}

	f := excelize.NewFile()

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"1E40AF"}},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "left",
		},
	})

	usedSheetNames := make(map[string]int)

	for idx, gw := range guruWaliList {
		sheetName := sanitizeSheetName(gw.NamaGuruWali)

		// Hindari nama sheet duplikat (kalau ada 2 guru wali dengan nama mirip setelah sanitasi)
		if count, exists := usedSheetNames[sheetName]; exists {
			usedSheetNames[sheetName] = count + 1
			suffix := fmt.Sprintf(" (%d)", count+1)
			if len(sheetName)+len(suffix) > 31 {
				sheetName = sheetName[:31-len(suffix)]
			}
			sheetName += suffix
		} else {
			usedSheetNames[sheetName] = 1
		}

		if idx == 0 {
			f.SetSheetName("Sheet1", sheetName)
		} else {
			f.NewSheet(sheetName)
		}

		// Info guru wali di baris atas
		f.SetCellValue(sheetName, "A1", "Guru Wali")
		f.SetCellValue(sheetName, "B1", gw.NamaGuruWali)
		f.SetCellValue(sheetName, "A2", "Guru")
		f.SetCellValue(sheetName, "B2", gw.NamaPegawai)
		f.SetCellValue(sheetName, "A3", "NUPTK")
		f.SetCellValue(sheetName, "B3", gw.NUPTK)

		headers := []string{"No", "NIS", "Nama Siswa", "Rombel"}
		headerRow := 5

		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, headerRow)
			f.SetCellValue(sheetName, cell, h)
		}

		anggota := anggotaByGuruWali[gw.IDGuruWali]

		if len(anggota) == 0 {
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", headerRow+1), "Belum ada anggota terdaftar.")
		}

		for i, a := range anggota {
			r := headerRow + 1 + i

			f.SetCellValue(sheetName, fmt.Sprintf("A%d", r), i+1)
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", r), a.NIS)
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", r), a.NamaSiswa)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", r), a.NamaRombel)
		}

		f.SetColWidth(sheetName, "A", "A", 20)
		f.SetColWidth(sheetName, "B", "B", 15)
		f.SetColWidth(sheetName, "C", "C", 30)
		f.SetColWidth(sheetName, "D", "D", 18)

		if headerStyle != 0 {
			lastCol := fmt.Sprintf("D%d", headerRow)
			f.SetCellStyle(sheetName, fmt.Sprintf("A%d", headerRow), lastCol, headerStyle)
		}

		_ = f.SetPanes(sheetName, &excelize.Panes{
			Freeze:      true,
			YSplit:      headerRow,
			TopLeftCell: fmt.Sprintf("A%d", headerRow+1),
			ActivePane:  "bottomLeft",
		})
	}

	f.SetActiveSheet(0)

	filename := fmt.Sprintf("daftar-anggota-guru-wali-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return f.Write(c.Response().BodyWriter())
}

// =========================================
// Export Excel: Daftar Anggota (satu guru wali, dipakai di halaman Anggota.vue)
// =========================================

func ExportAnggotaByGuruWaliHandler(c *fiber.Ctx) error {
	idGuruWali, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID tidak valid."})
	}

	gw, err := GetByID(idGuruWali)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if gw == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Guru wali tidak ditemukan."})
	}
	namaPegawai, _ := GetNamaPegawai(gw.NUPTK)

	search := c.Query("search", "")
	idKelas, _ := strconv.Atoi(c.Query("id_kelas", "0"))
	idRombel, _ := strconv.Atoi(c.Query("id_rombel", "0"))

	anggota, err := GetAnggotaByGuruWaliForExport(idGuruWali, search, idKelas, idRombel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	f := excelize.NewFile()
	sheet := sanitizeSheetName(gw.NamaGuruWali)
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "Guru Wali")
	f.SetCellValue(sheet, "B1", gw.NamaGuruWali)
	f.SetCellValue(sheet, "A2", "Guru")
	f.SetCellValue(sheet, "B2", namaPegawai)
	f.SetCellValue(sheet, "A3", "NUPTK")
	f.SetCellValue(sheet, "B3", gw.NUPTK)

	headers := []string{"No", "Rombel", "NIS", "Nama Peserta Didik"}
	headerRow := 5

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, headerRow)
		f.SetCellValue(sheet, cell, h)
	}

	if len(anggota) == 0 {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", headerRow+1), "Belum ada anggota terdaftar.")
	}

	for i, a := range anggota {
		r := headerRow + 1 + i

		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), a.NamaRombel)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", r), a.NIS)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", r), a.NamaSiswa)
	}

	// Kolom NIS diberi style Text biar nol di depan tidak hilang / tidak jadi notasi ilmiah
	if len(anggota) > 0 {
		if textStyle, styleErr := f.NewStyle(&excelize.Style{
			CustomNumFmt: func() *string { s := "@"; return &s }(),
		}); styleErr == nil {
			firstCell := fmt.Sprintf("C%d", headerRow+1)
			lastCell := fmt.Sprintf("C%d", headerRow+len(anggota))
			f.SetCellStyle(sheet, firstCell, lastCell, textStyle)
		}
	}

	f.SetColWidth(sheet, "A", "A", 6)
	f.SetColWidth(sheet, "B", "B", 18)
	f.SetColWidth(sheet, "C", "C", 15)
	f.SetColWidth(sheet, "D", "D", 30)

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"1E40AF"}},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "left",
		},
	})
	if err == nil {
		lastHeaderCell := fmt.Sprintf("D%d", headerRow)
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", headerRow), lastHeaderCell, headerStyle)
	}

	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      headerRow,
		TopLeftCell: fmt.Sprintf("A%d", headerRow+1),
		ActivePane:  "bottomLeft",
	})

	filename2 := fmt.Sprintf(
		"anggota-%s-%s.xlsx",
		sanitizeSheetName(gw.NamaGuruWali),
		time.Now().Format("2006-01-02-15-04-05"),
	)

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename2))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return f.Write(c.Response().BodyWriter())
}

// =========================================
// Import Anggota via Template Excel
// =========================================

// DownloadTemplateAnggotaHandler menyediakan file Excel template import anggota
// (sheet data + sheet petunjuk), polanya disamakan dengan modul Ekskul.
func DownloadTemplateAnggotaHandler(c *fiber.Ctx) error {
	f := excelize.NewFile()

	sheet := "Anggota"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "NIS")

	f.SetColWidth(sheet, "A", "A", 20)

	_ = f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	// Kolom NIS diberi style Text biar nol di depan tidak hilang / tidak jadi notasi ilmiah
	customFmt := "@"
	if textStyle, err := f.NewStyle(&excelize.Style{CustomNumFmt: &customFmt}); err == nil {
		f.SetColStyle(sheet, "A", textStyle)
	}

	if headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"1E40AF"}},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "left",
		},
	}); err == nil {
		f.SetCellStyle(sheet, "A1", "A1", headerStyle)
	}

	// Sheet petunjuk pengisian (teks 1 kolom, sama seperti pola pesertadidik/ekskul)
	petunjuk := "Petunjuk"
	f.NewSheet(petunjuk)

	instructions := []string{
		"PETUNJUK PENGISIAN TEMPLATE UNGGAH ANGGOTA GURU WALI",
		"",
		"1. Jangan mengubah urutan atau nama kolom pada sheet 'Anggota'.",
		"2. Kolom 'NIS' wajib diisi, sesuai NIS peserta didik yang terdaftar.",
		"3. NIS tidak boleh duplikat.",
		"4. NIS yang sudah terdaftar sebagai anggota guru wali ini akan ditolak.",
		"5. Jika salah satu baris bermasalah, seluruh proses unggah akan dibatalkan.",
	}

	for i, line := range instructions {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetCellValue(petunjuk, cell, line)
	}

	f.SetColWidth(petunjuk, "A", "A", 100)
	f.SetActiveSheet(0)

	filename := fmt.Sprintf("template-anggota-guru-wali-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	return f.Write(c.Response().BodyWriter())
}

// anggotaImportRowError merepresentasikan satu baris bermasalah saat import,
// bentuknya {row, message} sesuai kontrak DashboardUploadDialog.vue.
type anggotaImportRowError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

// ImportAnggotaHandler memproses file Excel upload, validasi penuh (all-or-nothing) —
// pola disamakan dengan ImportAnggotaHandler modul Ekskul.
func ImportAnggotaHandler(c *fiber.Ctx) error {
	idGuruWali, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID guru wali tidak valid."})
	}

	gw, err := GetByID(idGuruWali)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if gw == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Guru wali tidak ditemukan."})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "File Excel wajib diunggah"})
	}

	src, err := fileHeader.Open()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Gagal membuka file"})
	}
	defer src.Close()

	f2, err := excelize.OpenReader(src)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Format file tidak valid, pastikan file .xlsx"})
	}

	sheet2 := f2.GetSheetList()[0]

	rows, err := f2.GetRows(sheet2)
	if err != nil || len(rows) < 2 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "File kosong atau tidak ada data untuk diimport"})
	}

	var errorsList []anggotaImportRowError
	var validNIS []string
	seenNIS := map[string]int{}

	for i, row := range rows[1:] {
		rowNum := i + 2

		if len(row) == 0 {
			continue
		}

		nis := strings.TrimSpace(row[0])
		if nis == "" {
			continue
		}

		if prev, ok := seenNIS[nis]; ok {
			errorsList = append(errorsList, anggotaImportRowError{
				Row:     rowNum,
				Message: fmt.Sprintf("NIS '%s' duplikat dengan baris %d", nis, prev),
			})
			continue
		}
		seenNIS[nis] = rowNum

		exists, err := ExistsPesertaDidik(nis)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
		}
		if !exists {
			errorsList = append(errorsList, anggotaImportRowError{
				Row:     rowNum,
				Message: fmt.Sprintf("NIS '%s' tidak terdaftar sebagai peserta didik", nis),
			})
			continue
		}

		sudahAnggota, err := IsAnggotaExists(idGuruWali, nis)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
		}
		if sudahAnggota {
			errorsList = append(errorsList, anggotaImportRowError{
				Row:     rowNum,
				Message: fmt.Sprintf("NIS '%s' sudah terdaftar di guru wali ini.", nis),
			})
			continue
		}

		validNIS = append(validNIS, nis)
	}

	if len(errorsList) > 0 {
		sort.Slice(errorsList, func(i, j int) bool { return errorsList[i].Row < errorsList[j].Row })

		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Ditemukan data bermasalah, seluruh proses import dibatalkan",
			"errors":  errorsList,
		})
	}

	if len(validNIS) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Tidak ada data valid untuk diimport"})
	}

	for _, nis := range validNIS {
		if _, err := AddAnggota(idGuruWali, nis); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal menyimpan data: " + err.Error()})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("%d anggota berhasil ditambahkan", len(validNIS)),
		"data":    fiber.Map{"imported": len(validNIS)},
	})
}
