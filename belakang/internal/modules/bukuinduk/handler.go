package bukuinduk

import (
	"fmt"
	"strconv"
	"time"

	"sekolah/internal/auth"

	"github.com/gofiber/fiber/v2"
)

type requestBulkCetak struct {
	TahunPelajaran string           `json:"tahun_pelajaran"`
	IDs            []string         `json:"ids"`
	Opsi           CetakOpsiRequest `json:"opsi"`
}

func parseSumber(raw string) (SumberData, error) {
	switch SumberData(raw) {
	case SumberAktif, SumberMutasi, SumberArsip:
		return SumberData(raw), nil
	default:
		return "", fmt.Errorf("sumber tidak valid: %s", raw)
	}
}

func fetchOne(sumber SumberData, rawID string) (*BukuIndukData, error) {
	switch sumber {
	case SumberAktif:
		return GetAktifByNIS(rawID)
	case SumberMutasi:
		return GetMutasiByNIS(rawID)
	case SumberArsip:
		kodeSiswa, err := strconv.Atoi(rawID)
		if err != nil {
			return nil, fmt.Errorf("kode_siswa tidak valid: %s", rawID)
		}
		return GetArsipByKodeSiswa(kodeSiswa)
	default:
		return nil, fmt.Errorf("sumber tidak dikenal")
	}
}

// cetakOpsiBase = Opsi Cetak default yg dipakai kalau tombol Cetak (di daftar Peserta Didik)
// tidak mengirim override apa pun - yaitu Opsi Cetak yg disimpan admin lewat menu Buku Induk >
// Pengaturan. Gagal baca (mis. tabel belum ada) tidak menggagalkan cetak, fallback ke
// DefaultCetakOpsi() (cetak semua lembar).
func cetakOpsiBase() CetakOpsi {
	opsi, err := GetCetakOpsiSetting()
	if err != nil {
		return DefaultCetakOpsi()
	}

	return opsi
}

func CetakSatuHandler(c *fiber.Ctx) error {
	sumber, err := parseSumber(c.Params("sumber"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	id := c.Params("id")
	tahunPelajaran := c.Query("tahun_pelajaran")
	if tahunPelajaran == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "tahun_pelajaran wajib diisi"})
	}

	opsi := parseCetakOpsiQuery(c).ToCetakOpsiWithBase(cetakOpsiBase())
	if !opsi.adaLembarDicetak() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Pilih minimal 1 lembar yang akan dicetak."})
	}

	data, err := fetchOne(sumber, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Data siswa tidak ditemukan"})
	}

	return renderAndSend(c, []BukuIndukData{*data}, tahunPelajaran, fmt.Sprintf("buku-induk-%s", data.NIS), opsi)
}

func CetakBulkHandler(c *fiber.Ctx) error {
	sumber, err := parseSumber(c.Params("sumber"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	var req requestBulkCetak
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Body tidak valid"})
	}
	if req.TahunPelajaran == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "tahun_pelajaran wajib diisi"})
	}
	if len(req.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Tidak ada siswa yang dipilih"})
	}

	opsi := req.Opsi.ToCetakOpsiWithBase(cetakOpsiBase())
	if !opsi.adaLembarDicetak() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Pilih minimal 1 lembar yang akan dicetak."})
	}

	var list []BukuIndukData
	for _, rawID := range req.IDs {
		data, err := fetchOne(sumber, rawID)
		if err != nil {
			continue
		}
		list = append(list, *data)
	}
	if len(list) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Tidak ada data siswa yang valid ditemukan"})
	}

	return renderAndSend(c, list, req.TahunPelajaran, "buku-induk-massal", opsi)
}

func renderAndSend(c *fiber.Ctx, siswaList []BukuIndukData, tahunPelajaran string, filenamePrefix string, opsi CetakOpsi) error {
	profil, err := GetSekolahProfil()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mengambil profil sekolah"})
	}
	kepsek, err := GetKepalaSekolahInfo()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data kepala sekolah"})
	}

	claims := c.Locals("claims").(*auth.Claims)
	namaPetugas, err := GetNamaPengguna(claims.IDPengguna)
	if err != nil {
		namaPetugas = ""
	}

	var nipPetugas string
	if claims.NUPTK != "" {
		nipPetugas, err = GetNIPPegawai(claims.NUPTK)
		if err != nil {
			nipPetugas = "" // gagal lookup tidak menggagalkan seluruh cetak
		}
	}

	minNIS, maxNIS := computeRentangNIS(siswaList)

	pdfBytes, err := GenerateBukuIndukPDF(BukuIndukPDFInput{
		SiswaList:       siswaList,
		Profil:          *profil,
		Kepsek:          *kepsek,
		TahunPelajaran:  tahunPelajaran,
		NamaPetugas:     namaPetugas,
		NIPPetugas:      nipPetugas,
		RentangNISAwal:  minNIS,
		RentangNISAkhir: maxNIS,
		Opsi:            opsi,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal membuat PDF: " + err.Error()})
	}

	filename := fmt.Sprintf("%s-%s.pdf", filenamePrefix, time.Now().Format("2006-01-02-15-04-05"))
	c.Set("Access-Control-Expose-Headers", "Content-Disposition")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Set("Content-Type", "application/pdf")
	return c.Send(pdfBytes)
}

func computeRentangNIS(siswaList []BukuIndukData) (string, string) {
	if len(siswaList) == 0 {
		return "", ""
	}
	minNIS, maxNIS := siswaList[0].NIS, siswaList[0].NIS
	for _, s := range siswaList[1:] {
		if s.NIS < minNIS {
			minNIS = s.NIS
		}
		if s.NIS > maxNIS {
			maxNIS = s.NIS
		}
	}
	return minNIS, maxNIS
}
