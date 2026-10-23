package ekskul

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"sekolah/internal/auth"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

// isForeignKeyError: MySQL error 1451 (delete parent row gagal karena masih
// dirujuk baris lain) / 1216 (varian lama). Dipakai supaya pesan hapus gagal
// karena ekskul sudah punya nilai (nilai_ekskul) atau keanggotaan per periode
// (ekskul_anggota) tidak menampilkan error SQL mentah ke user.
func isForeignKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError

	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1451 || mysqlErr.Number == 1216
	}

	return false
}

// ================= EKSKUL =================

func GetAdminListHandler(c *fiber.Ctx) error {
	search := c.Query("search", "")
	statusAktif, _ := strconv.Atoi(c.Query("aktif", "-1"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sort_by", "ekskul")
	sortDir := c.Query("sort_dir", "asc")

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	data, err := GetAllAdmin(search, statusAktif, sortBy, sortDir, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	total, err := CountAdmin(search, statusAktif)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func GetByIDHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid.",
		})
	}

	data, err := GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Ekskul tidak ditemukan.",
			})
		}

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

func SearchPegawaiHandler(c *fiber.Ctx) error {
	search := c.Query("q", "")

	data, err := SearchPegawai(search, 15)
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

// GetMeHandler mengembalikan daftar ekskul yang diikuti siswa yang sedang login.
// Kalau belum ikut ekskul manapun, data dikembalikan kosong (bukan 404/error).
func GetMeHandler(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*auth.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Tidak terautentikasi.",
		})
	}

	data, err := GetEkskulBySiswa(claims.NIS)
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

func validateEkskul(nama, nuptk string, excludeID int) (int, string, bool) {
	pembinaValid, err := PembinaValid(nuptk)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error(), false
	}
	if !pembinaValid {
		return fiber.StatusBadRequest, "Pembina tidak ditemukan atau bukan Guru.", false
	}

	namaTaken, err := ExistsByNama(nama, excludeID)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error(), false
	}
	if namaTaken {
		return fiber.StatusConflict, "Nama ekskul sudah digunakan.", false
	}

	return 0, "", true
}

func CreateHandler(c *fiber.Ctx) error {
	var input Ekskul

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	input.NamaEkskul = strings.TrimSpace(input.NamaEkskul)
	input.NUPTK = strings.TrimSpace(input.NUPTK)

	if input.NamaEkskul == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama ekskul wajib diisi.",
		})
	}

	if input.NUPTK == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Pembina wajib dipilih.",
		})
	}

	if status, msg, ok := validateEkskul(input.NamaEkskul, input.NUPTK, 0); !ok {
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

	newID, err := Create(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data, err := GetByID(int(newID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Ekskul berhasil ditambahkan.",
		"data":    data,
	})
}

func UpdateHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid.",
		})
	}

	var input Ekskul

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	input.NamaEkskul = strings.TrimSpace(input.NamaEkskul)
	input.NUPTK = strings.TrimSpace(input.NUPTK)

	if input.NamaEkskul == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Nama ekskul wajib diisi.",
		})
	}

	if input.NUPTK == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Pembina wajib dipilih.",
		})
	}

	oldData, err := GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Ekskul tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if oldData.NamaEkskul == input.NamaEkskul &&
		oldData.NUPTK == input.NUPTK &&
		oldData.Aktif == input.Aktif {
		return c.JSON(fiber.Map{
			"success": true,
			"changed": false,
			"message": "Tidak ada perubahan data.",
			"data":    oldData,
		})
	}

	if status, msg, ok := validateEkskul(input.NamaEkskul, input.NUPTK, id); !ok {
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"message": msg,
		})
	}

	if err := Update(id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	data, err := GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"changed": true,
		"message": "Ekskul berhasil diperbarui.",
		"data":    data,
	})
}

func DeleteHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid.",
		})
	}

	if _, err := GetByID(id); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Ekskul tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	total, err := CountAnggotaByEkskul(id, "", 0, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	if total > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "Ekskul masih memiliki anggota, tidak bisa dihapus.",
		})
	}

	if err := Delete(id); err != nil {
		if isForeignKeyError(err) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Ekskul ini sudah punya nilai ekstrakurikuler atau riwayat keanggotaan yang tercatat, sehingga tidak bisa dihapus. Nonaktifkan saja ekskul ini (lewat menu edit) jika sudah tidak dipakai.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Ekskul berhasil dihapus.",
	})
}

// ================= ANGGOTA EKSKUL =================

func GetAnggotaListHandler(c *fiber.Ctx) error {
	idEkskul, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID ekskul tidak valid.",
		})
	}

	if _, err := GetByID(idEkskul); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Ekskul tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sort_by", "rombel")
	sortDir := c.Query("sort_dir", "asc")
	idKelas, _ := strconv.Atoi(c.Query("id_kelas", "0"))
	idRombel, _ := strconv.Atoi(c.Query("id_rombel", "0"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	data, err := GetAnggotaByEkskul(idEkskul, search, sortBy, sortDir, limit, offset, idKelas, idRombel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	total, err := CountAnggotaByEkskul(idEkskul, search, idKelas, idRombel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func SearchPesertaDidikHandler(c *fiber.Ctx) error {
	search := c.Query("q", "")

	data, err := SearchPesertaDidik(search, 15)
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

func AddAnggotaHandler(c *fiber.Ctx) error {
	idEkskul, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID ekskul tidak valid.",
		})
	}

	if _, err := GetByID(idEkskul); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Ekskul tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	var input struct {
		NIS string `json:"nis"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	input.NIS = strings.TrimSpace(input.NIS)

	if input.NIS == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Siswa wajib dipilih.",
		})
	}

	siswaValid, err := SiswaAktifValid(input.NIS)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	if !siswaValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Siswa tidak ditemukan atau tidak aktif tahun berjalan.",
		})
	}

	sudahAnggota, err := AnggotaExists(idEkskul, input.NIS)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	if sudahAnggota {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "Siswa sudah terdaftar di ekskul ini.",
		})
	}

	if _, err := AddAnggota(idEkskul, input.NIS); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Anggota berhasil ditambahkan.",
	})
}

func RemoveAnggotaHandler(c *fiber.Ctx) error {
	idAnggota, err := strconv.Atoi(c.Params("anggota_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID anggota tidak valid.",
		})
	}

	if err := RemoveAnggota(idAnggota); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Anggota tidak ditemukan.",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Anggota berhasil dihapus.",
	})
}

// BulkRemoveAnggotaHandler menghapus banyak anggota sekaligus (pola sama spt bulk-delete pengguna).
func BulkRemoveAnggotaHandler(c *fiber.Ctx) error {
	var input struct {
		IDs []int `json:"ids"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Payload tidak valid.",
		})
	}

	if len(input.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Tidak ada anggota yang dipilih.",
		})
	}

	if err := BulkRemoveAnggota(input.IDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("%d anggota berhasil dihapus.", len(input.IDs)),
	})
}

// =========================================
// Export Excel: Daftar Ekskul
// =========================================

func ExportEkskulHandler(c *fiber.Ctx) error {
	data, err := GetAllForExport()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	f := excelize.NewFile()
	sheet := "Ekskul"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{
		"No",
		"Nama Ekskul",
		"Pembina",
		"NUPTK",
		"Status",
		"Jumlah Anggota",
	}

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
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), row.NamaEkskul)
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
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"1E40AF"},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "left",
		},
	})

	if err == nil {
		f.SetCellStyle(sheet, "A1", "F1", headerStyle)
	}

	filename := fmt.Sprintf(
		"daftar-ekskul-%s.xlsx",
		time.Now().Format("2006-01-02-15-04-05"),
	)

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return f.Write(c.Response().BodyWriter())
}

// =========================================
// Export Excel: Daftar Anggota (Semua Ekskul, sheet per ekskul)
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
		clean = "Ekskul"
	}

	return clean
}

func ExportAnggotaHandler(c *fiber.Ctx) error {
	ekskulList, err := GetAllForExport()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if len(ekskulList) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Belum ada data ekskul untuk diexport.",
		})
	}

	anggotaList, err := GetAllAnggotaForExport()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Kelompokkan anggota per id_ekskul
	anggotaByEkskul := make(map[int][]AnggotaEkskul)
	for _, a := range anggotaList {
		anggotaByEkskul[a.IDEkskul] = append(anggotaByEkskul[a.IDEkskul], a)
	}

	f := excelize.NewFile()

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"1E40AF"},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "left",
		},
	})

	usedSheetNames := make(map[string]int)

	for idx, ek := range ekskulList {
		sheetName := sanitizeSheetName(ek.NamaEkskul)

		// Hindari nama sheet duplikat (kalau ada 2 ekskul dengan nama mirip setelah sanitasi)
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

		// Info pembina di baris atas
		f.SetCellValue(sheetName, "A1", "Ekskul")
		f.SetCellValue(sheetName, "B1", ek.NamaEkskul)
		f.SetCellValue(sheetName, "A2", "Pembina")
		f.SetCellValue(sheetName, "B2", ek.NamaPegawai)
		f.SetCellValue(sheetName, "A3", "NUPTK")
		f.SetCellValue(sheetName, "B3", ek.NUPTK)

		headers := []string{"No", "NIS", "Nama Siswa", "Rombel"}
		headerRow := 5

		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, headerRow)
			f.SetCellValue(sheetName, cell, h)
		}

		anggota := anggotaByEkskul[ek.IDEkskul]

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

	filename := fmt.Sprintf(
		"daftar-anggota-ekskul-%s.xlsx",
		time.Now().Format("2006-01-02-15-04-05"),
	)

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return f.Write(c.Response().BodyWriter())
}

// =========================================
// Export Excel: Daftar Anggota (satu ekskul, dipakai di halaman Anggota.vue)
// =========================================

func ExportAnggotaByEkskulHandler(c *fiber.Ctx) error {
	idEkskul, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid.",
		})
	}

	ek, err := GetByID(idEkskul)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Ekskul tidak ditemukan.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	search := c.Query("search", "")
	idKelas, _ := strconv.Atoi(c.Query("id_kelas", "0"))
	idRombel, _ := strconv.Atoi(c.Query("id_rombel", "0"))

	anggota, err := GetAnggotaByEkskulForExport(idEkskul, search, idKelas, idRombel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	f := excelize.NewFile()
	sheet := sanitizeSheetName(ek.NamaEkskul)
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "Ekskul")
	f.SetCellValue(sheet, "B1", ek.NamaEkskul)
	f.SetCellValue(sheet, "A2", "Pembina")
	f.SetCellValue(sheet, "B2", ek.NamaPegawai)
	f.SetCellValue(sheet, "A3", "NUPTK")
	f.SetCellValue(sheet, "B3", ek.NUPTK)

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
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"1E40AF"},
		},
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

	filename := fmt.Sprintf(
		"anggota-%s-%s.xlsx",
		sanitizeSheetName(ek.NamaEkskul),
		time.Now().Format("2006-01-02-15-04-05"),
	)

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")

	return f.Write(c.Response().BodyWriter())
}

// =========================================
// Import Anggota via Template Excel
// =========================================

// DownloadTemplateAnggotaHandler menyediakan file Excel template import anggota
// (sheet data + sheet petunjuk), polanya disamakan dengan pesertadidik.TemplateHandler
// (package pesertadidik, import.go): sheet data hanya berisi header (tanpa contoh data),
// freeze header hanya di sheet data, sheet Petunjuk berupa teks 1 kolom (bukan tabel).
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
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"1E40AF"},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "left",
		},
	}); err == nil {
		f.SetCellStyle(sheet, "A1", "A1", headerStyle)
	}

	// Sheet petunjuk pengisian (teks 1 kolom, sama seperti pesertadidik.TemplateHandler)
	petunjuk := "Petunjuk"
	f.NewSheet(petunjuk)

	instructions := []string{
		"PETUNJUK PENGISIAN TEMPLATE UNGGAH ANGGOTA EKSKUL",
		"",
		"1. Jangan mengubah urutan atau nama kolom pada sheet 'Anggota'.",
		"2. Kolom 'NIS' wajib diisi, sesuai NIS peserta didik yang terdaftar.",
		"3. NIS tidak boleh duplikat.",
		"4. NIS yang sudah terdaftar sebagai anggota ekskul ini akan ditolak.",
		"5. Jika salah satu baris bermasalah, seluruh proses unggah akan dibatalkan.",
	}

	for i, line := range instructions {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetCellValue(petunjuk, cell, line)
	}

	f.SetColWidth(petunjuk, "A", "A", 100)
	f.SetActiveSheet(0)

	filename := fmt.Sprintf("template-anggota-ekskul-%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))

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
// pola disamakan dengan pesertadidik.ImportHandler (package pesertadidik, import.go).
func ImportAnggotaHandler(c *fiber.Ctx) error {
	idEkskul, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID ekskul tidak valid.",
		})
	}

	if _, err := GetByID(idEkskul); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Ekskul tidak ditemukan.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
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

	f, err := excelize.OpenReader(src)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Format file tidak valid, pastikan file .xlsx"})
	}

	sheet := f.GetSheetList()[0]

	rows, err := f.GetRows(sheet)
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

		siswaValid, err := SiswaAktifValid(nis)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		if !siswaValid {
			errorsList = append(errorsList, anggotaImportRowError{
				Row:     rowNum,
				Message: fmt.Sprintf("NIS '%s' tidak ditemukan atau tidak aktif tahun berjalan", nis),
			})
			continue
		}

		sudahAnggota, err := AnggotaExists(idEkskul, nis)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		if sudahAnggota {
			errorsList = append(errorsList, anggotaImportRowError{
				Row:     rowNum,
				Message: fmt.Sprintf("NIS '%s' sudah terdaftar di ekskul ini.", nis),
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
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Tidak ada data valid untuk diimport",
		})
	}

	for _, nis := range validNIS {
		if _, err := AddAnggota(idEkskul, nis); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal menyimpan data: " + err.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("%d anggota berhasil ditambahkan", len(validNIS)),
		"data":    fiber.Map{"imported": len(validNIS)},
	})
}
