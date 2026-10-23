package bukuinduk

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"

	"sekolah/internal/database"
)

// Wali kelas per TINGKAT (7/8/9) per siswa, tabel walikelas_tingkat. Dipakai halaman Data Nilai
// Kelas 7/8/9 di PDF supaya nama wali kelas sesuai kelasnya waktu itu, bukan wali kelas rombel
// siswa saat ini. Belum diisi = nama/NIP dikosongkan (ditandatangani manual), TANPA fallback ke
// wali kelas rombel sekarang.

// WaliKelasTingkatInfo = nama & NIP wali kelas untuk 1 tingkat.
type WaliKelasTingkatInfo struct {
	Nama string
	NIP  string
}

func tingkatValid(t int) bool { return t == 7 || t == 8 || t == 9 }

// GetWaliKelasPerTingkat mengembalikan wali kelas siswa untuk tiap tingkat yg sudah diisi.
func GetWaliKelasPerTingkat(nis string) (map[int]WaliKelasTingkatInfo, error) {
	var rows []struct {
		Tingkat int    `db:"tingkat"`
		Nama    string `db:"nama_walikelas"`
		NIP     string `db:"nip"`
	}

	err := database.DB.Select(&rows, `
		SELECT tingkat, COALESCE(nama_walikelas,'') AS nama_walikelas, COALESCE(nip,'') AS nip
		FROM walikelas_tingkat WHERE nis = ?`, nis)
	if err != nil {
		return nil, err
	}

	out := map[int]WaliKelasTingkatInfo{}
	for _, r := range rows {
		out[r.Tingkat] = WaliKelasTingkatInfo{Nama: r.Nama, NIP: r.NIP}
	}

	return out, nil
}

const upsertWaliTingkatSQL = `
	INSERT INTO walikelas_tingkat (nis, tingkat, nama_walikelas, nip) VALUES (?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE nama_walikelas = VALUES(nama_walikelas), nip = VALUES(nip)`

func jsonErr(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{"success": false, "message": msg})
}

// --- GET /bukuinduk/walikelas-tingkat/pegawai?q= : cari guru utk dropdown ---

func SearchPegawaiWaliTingkatHandler(c *fiber.Ctx) error {
	q := "%" + strings.TrimSpace(c.Query("q")) + "%"

	var rows []struct {
		NUPTK string `json:"nuptk" db:"nuptk"`
		Nama  string `json:"nama" db:"nama_pegawai"`
		NIP   string `json:"nip" db:"nip_baru"`
	}

	err := database.DB.Select(&rows, `
		SELECT COALESCE(nuptk,'') AS nuptk, COALESCE(nama_pegawai,'') AS nama_pegawai, COALESCE(nip_baru,'') AS nip_baru
		FROM pegawai
		WHERE jabatan = 'Guru' AND (nama_pegawai LIKE ? OR nuptk LIKE ? OR nip_baru LIKE ?)
		ORDER BY nama_pegawai LIMIT 30`, q, q, q)
	if err != nil {
		return jsonErr(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"success": true, "data": rows})
}

// --- GET /bukuinduk/walikelas-tingkat?id_rombel=&tingkat= : siswa rombel + wali kelas tercatat ---

func ListWaliTingkatHandler(c *fiber.Ctx) error {
	idRombel, _ := strconv.Atoi(c.Query("id_rombel"))
	tingkat, _ := strconv.Atoi(c.Query("tingkat"))

	if idRombel == 0 || !tingkatValid(tingkat) {
		return jsonErr(c, fiber.StatusBadRequest, "id_rombel dan tingkat (7/8/9) wajib diisi")
	}

	rows := []struct {
		NIS  string `json:"nis" db:"nis"`
		Nama string `json:"nama_siswa" db:"nama_siswa"`
		Wali string `json:"nama_walikelas" db:"nama_walikelas"`
		NIP  string `json:"nip" db:"nip"`
	}{}

	err := database.DB.Select(&rows, `
		SELECT pd.nis, COALESCE(pd.nama_siswa,'') AS nama_siswa,
			COALESCE(w.nama_walikelas,'') AS nama_walikelas, COALESCE(w.nip,'') AS nip
		FROM pesertadidik pd
		LEFT JOIN walikelas_tingkat w ON w.nis = pd.nis AND w.tingkat = ?
		WHERE pd.id_rombel = ?
		ORDER BY pd.nama_siswa`, tingkat, idRombel)
	if err != nil {
		return jsonErr(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"success": true, "data": rows})
}

// --- PUT /bukuinduk/walikelas-tingkat/rombel : isi wali kelas utk SEMUA siswa 1 rombel ---

type requestWaliRombel struct {
	IDRombel int    `json:"id_rombel"`
	Tingkat  int    `json:"tingkat"`
	NUPTK    string `json:"nuptk"` // dipilih dari pegawai (nama & NIP diambil server)
	Nama     string `json:"nama"`  // dipakai kalau NUPTK kosong (isian manual)
	NIP      string `json:"nip"`
}

func SaveWaliTingkatRombelHandler(c *fiber.Ctx) error {
	var req requestWaliRombel
	if err := c.BodyParser(&req); err != nil {
		return jsonErr(c, fiber.StatusBadRequest, "Body tidak valid")
	}

	if req.IDRombel == 0 || !tingkatValid(req.Tingkat) {
		return jsonErr(c, fiber.StatusBadRequest, "id_rombel dan tingkat (7/8/9) wajib diisi")
	}

	nama, nip := strings.TrimSpace(req.Nama), strings.TrimSpace(req.NIP)

	if req.NUPTK != "" {
		var p struct {
			Nama string `db:"nama_pegawai"`
			NIP  string `db:"nip_baru"`
		}

		if err := database.DB.Get(&p, `
			SELECT COALESCE(nama_pegawai,'') AS nama_pegawai, COALESCE(nip_baru,'') AS nip_baru
			FROM pegawai WHERE nuptk = ?`, req.NUPTK); err != nil {
			return jsonErr(c, fiber.StatusBadRequest, "Pegawai tidak ditemukan")
		}

		nama, nip = p.Nama, p.NIP
	}

	if nama == "" {
		return jsonErr(c, fiber.StatusBadRequest, "Pilih guru atau isi nama wali kelas")
	}

	var niss []string
	if err := database.DB.Select(&niss, `SELECT nis FROM pesertadidik WHERE id_rombel = ?`, req.IDRombel); err != nil {
		return jsonErr(c, fiber.StatusInternalServerError, err.Error())
	}

	if len(niss) == 0 {
		return jsonErr(c, fiber.StatusBadRequest, "Rombel ini tidak punya siswa")
	}

	tx, err := database.DB.Beginx()
	if err != nil {
		return jsonErr(c, fiber.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback() //nolint:errcheck

	for _, nis := range niss {
		if _, err := tx.Exec(upsertWaliTingkatSQL, nis, req.Tingkat, nama, nip); err != nil {
			return jsonErr(c, fiber.StatusInternalServerError, "Gagal menyimpan: "+err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		return jsonErr(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Wali kelas Kelas %d disimpan untuk %d siswa.", req.Tingkat, len(niss)),
	})
}

// --- DELETE /bukuinduk/walikelas-tingkat/rombel?id_rombel=&tingkat= : kosongkan ---

func ClearWaliTingkatRombelHandler(c *fiber.Ctx) error {
	idRombel, _ := strconv.Atoi(c.Query("id_rombel"))
	tingkat, _ := strconv.Atoi(c.Query("tingkat"))

	if idRombel == 0 || !tingkatValid(tingkat) {
		return jsonErr(c, fiber.StatusBadRequest, "id_rombel dan tingkat (7/8/9) wajib diisi")
	}

	_, err := database.DB.Exec(`
		DELETE FROM walikelas_tingkat
		WHERE tingkat = ? AND nis IN (SELECT nis FROM pesertadidik WHERE id_rombel = ?)`, tingkat, idRombel)
	if err != nil {
		return jsonErr(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"success": true, "message": "Wali kelas rombel ini untuk Kelas " + strconv.Itoa(tingkat) + " dikosongkan."})
}

// --- GET /bukuinduk/walikelas-tingkat/template : template Excel impor ---

func TemplateWaliTingkatHandler(c *fiber.Ctx) error {
	f := excelize.NewFile()
	defer f.Close()

	sheet := f.GetSheetName(0)
	_ = f.SetSheetRow(sheet, "A1", &[]string{"NIS", "Nama Wali Kelas", "NIP"})
	_ = f.SetSheetRow(sheet, "A2", &[]string{"13659", "Nama Guru, S.Pd", "198001012005011001"})
	_ = f.SetColWidth(sheet, "A", "A", 16)
	_ = f.SetColWidth(sheet, "B", "B", 36)
	_ = f.SetColWidth(sheet, "C", "C", 24)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return jsonErr(c, fiber.StatusInternalServerError, err.Error())
	}

	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	c.Set("Content-Disposition", `attachment; filename="template-walikelas-tingkat.xlsx"`)
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	return c.Send(buf.Bytes())
}

// --- POST /bukuinduk/walikelas-tingkat/import (multipart: file, tingkat) ---
// Kolom: NIS | Nama Wali Kelas | NIP. Dipakai utk Kelas 7/8 yg rombel lamanya sudah bubar
// (siswa sudah tercampur), jadi wali kelas diisi per NIS.

func ImportWaliTingkatHandler(c *fiber.Ctx) error {
	tingkat, _ := strconv.Atoi(c.FormValue("tingkat"))
	if !tingkatValid(tingkat) {
		return jsonErr(c, fiber.StatusBadRequest, "tingkat (7/8/9) wajib diisi")
	}

	fh, err := c.FormFile("file")
	if err != nil {
		return jsonErr(c, fiber.StatusBadRequest, "File Excel wajib diunggah")
	}

	src, err := fh.Open()
	if err != nil {
		return jsonErr(c, fiber.StatusBadRequest, "File tidak bisa dibaca")
	}
	defer src.Close()

	x, err := excelize.OpenReader(src)
	if err != nil {
		return jsonErr(c, fiber.StatusBadRequest, "Bukan file Excel (.xlsx) yang valid")
	}
	defer x.Close()

	sheet := x.GetSheetName(0)

	rows, err := x.GetRows(sheet)
	if err != nil || len(rows) < 2 {
		return jsonErr(c, fiber.StatusBadRequest, "File kosong (baris pertama harus header: NIS, Nama Wali Kelas, NIP)")
	}

	type item struct{ nis, nama, nip string }

	var items []item

	cell := func(r []string, i int) string {
		if i < len(r) {
			return strings.TrimSpace(r[i])
		}

		return ""
	}

	for _, r := range rows[1:] {
		nis := cell(r, 0)
		if nis == "" {
			continue
		}

		items = append(items, item{nis: nis, nama: cell(r, 1), nip: cell(r, 2)})
	}

	if len(items) == 0 {
		return jsonErr(c, fiber.StatusBadRequest, "Tidak ada baris data (kolom NIS kosong semua)")
	}

	// NIS yg dikenal (aktif / mutasi / arsip)
	allNIS := make([]string, 0, len(items))
	for _, it := range items {
		allNIS = append(allNIS, it.nis)
	}

	query, args, err := sqlx.In(`
		SELECT nis FROM pesertadidik WHERE nis IN (?)
		UNION SELECT nis FROM pesertadidik_mutasi WHERE nis IN (?)
		UNION SELECT nis FROM pesertadidik_arsip WHERE nis IN (?)`, allNIS, allNIS, allNIS)
	if err != nil {
		return jsonErr(c, fiber.StatusInternalServerError, err.Error())
	}

	var dikenal []string
	if err := database.DB.Select(&dikenal, query, args...); err != nil {
		return jsonErr(c, fiber.StatusInternalServerError, err.Error())
	}

	ada := map[string]bool{}
	for _, n := range dikenal {
		ada[n] = true
	}

	tx, err := database.DB.Beginx()
	if err != nil {
		return jsonErr(c, fiber.StatusInternalServerError, err.Error())
	}
	defer tx.Rollback() //nolint:errcheck

	tersimpan, dilewati := 0, []string{}

	for _, it := range items {
		if !ada[it.nis] {
			dilewati = append(dilewati, it.nis+" (NIS tidak ditemukan)")
			continue
		}

		if it.nama == "" {
			dilewati = append(dilewati, it.nis+" (nama wali kelas kosong)")
			continue
		}

		if _, err := tx.Exec(upsertWaliTingkatSQL, it.nis, tingkat, it.nama, it.nip); err != nil {
			return jsonErr(c, fiber.StatusInternalServerError, "Gagal menyimpan: "+err.Error())
		}

		tersimpan++
	}

	if err := tx.Commit(); err != nil {
		return jsonErr(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"message":   fmt.Sprintf("Impor selesai: %d siswa tersimpan, %d dilewati.", tersimpan, len(dilewati)),
		"tersimpan": tersimpan,
		"dilewati":  dilewati,
	})
}
