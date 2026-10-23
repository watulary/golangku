package nilai

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"

	"sekolah/internal/database"
)

// pilihKelasCatat menentukan kelas tempat nilai / prestasi dicatat.
// idKelasCatat <= 0 = ikut kelas rombel yang dipilih (perilaku biasa). Kalau diisi, kelas itulah yang
// dipakai; gunanya untuk mengisi data tahun lalu siswa yang sekarang sudah di kelas lebih tinggi
// (rombel sekarang dipilih supaya siswanya tampil, kelas dicatat sebagai kelas saat itu).
func pilihKelasCatat(idKelasCatat, idKelasRombel int, namaKelasRombel string) (int, string, *fiber.Error) {
	if idKelasCatat <= 0 {
		return idKelasRombel, namaKelasRombel, nil
	}

	var nama string

	err := database.DB.Get(&nama, `SELECT nama_kelas FROM tbl_kelas WHERE id_kelas = ? LIMIT 1`, idKelasCatat)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", fiber.NewError(fiber.StatusBadRequest, "Kelas untuk pencatatan tidak ditemukan.")
		}

		return 0, "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return idKelasCatat, nama, nil
}
