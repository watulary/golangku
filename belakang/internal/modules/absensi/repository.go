package absensi

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"sekolah/internal/database"
)

// GetDaftarAbsensi mengambil daftar siswa dalam satu rombel untuk tanggal tertentu,
// hasil LEFT JOIN pesertadidik -> absensi, urut alfabetis by nama_siswa.
// Siswa yang belum punya record absensi di tanggal ini tampil dengan status default "Hadir".
func GetDaftarAbsensi(tanggal string, idKelas, idRombel int) ([]AbsensiSiswa, error) {
	var list []AbsensiSiswa
	query := `
		SELECT
			IFNULL(a.id, 0) AS id,
			p.nis AS nis,
			p.nama_siswa AS nama_siswa,
			k.nama_kelas AS nama_kelas,
			r.nama_rombel AS nama_rombel,
			COALESCE(a.status, 'Hadir') AS status,
			COALESCE(a.keterangan, '') AS keterangan
		FROM pesertadidik p
		LEFT JOIN absensi a ON a.nis = p.nis AND a.tanggal = ?
		LEFT JOIN tbl_kelas k ON k.id_kelas = p.id_kelas
		LEFT JOIN tbl_rombel r ON r.id_rombel = p.id_rombel
		WHERE (? = 0 OR p.id_kelas = ?)
		  AND (? = 0 OR p.id_rombel = ?)
		ORDER BY k.nama_kelas ASC, r.nama_rombel ASC, p.nama_siswa ASC
	`
	err := database.DB.Select(&list, query, tanggal, idKelas, idKelas, idRombel, idRombel)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// BulkUpsert menyimpan/mengupdate absensi utk banyak siswa sekaligus dalam 1 transaction.
// Kalau record (nis, tanggal) sudah ada -> UPDATE, kalau belum -> INSERT.
func BulkUpsert(tanggal string, items []AbsensiBulkItem) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range items {
		var existingID int
		err := tx.Get(&existingID, `SELECT id FROM absensi WHERE nis = ? AND tanggal = ?`, item.Nis, tanggal)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if existingID > 0 {
			_, err = tx.Exec(
				`UPDATE absensi SET status = ?, keterangan = ? WHERE id = ?`,
				item.Status, item.Keterangan, existingID,
			)
		} else {
			_, err = tx.Exec(
				`INSERT INTO absensi (nis, tanggal, status, keterangan) VALUES (?, ?, ?, ?)`,
				item.Nis, tanggal, item.Status, item.Keterangan,
			)
		}
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

var statusKeSingkatan = map[string]string{
	"Hadir": "H", "Sakit": "S", "Izin": "I", "Alpa": "A",
}

// GetRekapBulanan mengambil matriks siswa x tanggal (1..akhir bulan) utk bulan/tahun tertentu.
// idKelas/idRombel = 0 berarti semua.
func GetRekapBulanan(bulan, tahun, idKelas, idRombel int) ([]RekapBulananSiswa, int, error) {
	daysInMonth := time.Date(tahun, time.Month(bulan)+1, 0, 0, 0, 0, 0, time.UTC).Day()

	type siswaRow struct {
		Nis        string `db:"nis"`
		NamaSiswa  string `db:"nama_siswa"`
		NamaKelas  string `db:"nama_kelas"`
		NamaRombel string `db:"nama_rombel"`
	}
	var siswaList []siswaRow
	querySiswa := `
		SELECT p.nis AS nis, p.nama_siswa AS nama_siswa,
			k.nama_kelas AS nama_kelas, r.nama_rombel AS nama_rombel
		FROM pesertadidik p
		LEFT JOIN tbl_kelas k ON k.id_kelas = p.id_kelas
		LEFT JOIN tbl_rombel r ON r.id_rombel = p.id_rombel
		WHERE (? = 0 OR p.id_kelas = ?)
		  AND (? = 0 OR p.id_rombel = ?)
		ORDER BY k.nama_kelas ASC, r.nama_rombel ASC, p.nama_siswa ASC
	`
	if err := database.DB.Select(&siswaList, querySiswa, idKelas, idKelas, idRombel, idRombel); err != nil {
		return nil, 0, err
	}

	type absensiRow struct {
		Nis     string `db:"nis"`
		Tanggal string `db:"tanggal"` // format YYYY-MM-DD
		Status  string `db:"status"`
	}
	var absensiRows []absensiRow
	tanggalAwal := fmt.Sprintf("%04d-%02d-01", tahun, bulan)
	tanggalAkhir := fmt.Sprintf("%04d-%02d-%02d", tahun, bulan, daysInMonth)
	queryAbsensi := `SELECT nis, tanggal, status FROM absensi WHERE tanggal BETWEEN ? AND ?`
	if err := database.DB.Select(&absensiRows, queryAbsensi, tanggalAwal, tanggalAkhir); err != nil {
		return nil, 0, err
	}

	// nis -> hari (1..N) -> kode singkatan
	absensiMap := make(map[string]map[int]string)
	for _, a := range absensiRows {
		if len(a.Tanggal) < 10 {
			continue
		}
		day, err := strconv.Atoi(a.Tanggal[8:10])
		if err != nil {
			continue
		}
		if absensiMap[a.Nis] == nil {
			absensiMap[a.Nis] = make(map[int]string)
		}
		absensiMap[a.Nis][day] = statusKeSingkatan[a.Status]
	}

	result := make([]RekapBulananSiswa, 0, len(siswaList))
	for _, s := range siswaList {
		harian := make([]string, daysInMonth)
		for d := 1; d <= daysInMonth; d++ {
			harian[d-1] = absensiMap[s.Nis][d]
		}
		result = append(result, RekapBulananSiswa{
			Nis: s.Nis, NamaSiswa: s.NamaSiswa, NamaKelas: s.NamaKelas, NamaRombel: s.NamaRombel,
			Harian: harian,
		})
	}
	return result, daysInMonth, nil
}

// GetAbsensiSaya mengambil seluruh record absensi milik 1 siswa dalam 1 bulan, urut tanggal ASC.
// Belum ada record di bulan itu -> slice kosong (bukan error), konsisten dgn pola "return kosong" Kenaikan/Kelulusan.
func GetAbsensiSaya(nis string, bulan, tahun int) ([]AbsensiSayaHarian, error) {
	var list []AbsensiSayaHarian
	daysInMonth := time.Date(tahun, time.Month(bulan)+1, 0, 0, 0, 0, 0, time.UTC).Day()
	tanggalAwal := fmt.Sprintf("%04d-%02d-01", tahun, bulan)
	tanggalAkhir := fmt.Sprintf("%04d-%02d-%02d", tahun, bulan, daysInMonth)

	query := `
		SELECT tanggal, status, COALESCE(keterangan, '') AS keterangan
		FROM absensi
		WHERE nis = ? AND tanggal BETWEEN ? AND ?
		ORDER BY tanggal ASC
	`
	if err := database.DB.Select(&list, query, nis, tanggalAwal, tanggalAkhir); err != nil {
		return nil, err
	}
	return list, nil
}
