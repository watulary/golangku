package nilai

import (
	"fmt"
	"strconv"
	"strings"

	"sekolah/internal/database"
)

// Ketidakhadiran = rekap Sakit / Izin / Alpa (Tanpa Keterangan) 1 siswa dalam 1 semester.
// Dihitung langsung dari tabel absensi, tidak disimpan.
type Ketidakhadiran struct {
	Sakit int `json:"sakit"`
	Izin  int `json:"izin"`
	Alpa  int `json:"alpa"`

	// Tercatat = jumlah seluruh record absensi semester itu (termasuk Hadir).
	// 0 berarti absensi semester itu belum pernah diinput, jadi angka 0 di atas
	// belum bisa dianggap "tidak pernah absen".
	Tercatat int `json:"tercatat"`
}

// RentangSemester mengubah nama tahun pelajaran ("2024/2025") + semester menjadi rentang tanggal:
//   - Ganjil : 1 Juli - 31 Desember tahun pertama  -> 2024-07-01 s/d 2024-12-31
//   - Genap  : 1 Januari - 30 Juni tahun kedua      -> 2025-01-01 s/d 2025-06-30
func RentangSemester(namaTahun, semester string) (awal, akhir string, err error) {
	parts := strings.Split(strings.TrimSpace(namaTahun), "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("format tahun pelajaran %q tidak dikenali (harus 2024/2025)", namaTahun)
	}

	y1, e1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	y2, e2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if e1 != nil || e2 != nil {
		return "", "", fmt.Errorf("format tahun pelajaran %q tidak dikenali (harus 2024/2025)", namaTahun)
	}

	switch semester {
	case "Ganjil":
		return fmt.Sprintf("%04d-07-01", y1), fmt.Sprintf("%04d-12-31", y1), nil
	case "Genap":
		return fmt.Sprintf("%04d-01-01", y2), fmt.Sprintf("%04d-06-30", y2), nil
	}

	return "", "", fmt.Errorf("semester %q tidak valid", semester)
}

func (k *Ketidakhadiran) tambah(status string, jumlah int) {
	k.Tercatat += jumlah

	switch status {
	case "Sakit":
		k.Sakit = jumlah
	case "Izin":
		k.Izin = jumlah
	case "Alpa":
		k.Alpa = jumlah
	}
}

// GetKetidakhadiran = rekap 1 siswa utk 1 semester. Dipakai generator PDF Buku Induk.
func GetKetidakhadiran(nis, namaTahun, semester string) (Ketidakhadiran, error) {
	var hasil Ketidakhadiran

	awal, akhir, err := RentangSemester(namaTahun, semester)
	if err != nil {
		return hasil, err
	}

	rows := []struct {
		Status string `db:"status"`
		Jumlah int    `db:"jumlah"`
	}{}

	err = database.DB.Select(&rows, `
		SELECT status, COUNT(*) AS jumlah
		FROM absensi
		WHERE nis = ? AND tanggal BETWEEN ? AND ?
		GROUP BY status
	`, nis, awal, akhir)
	if err != nil {
		return hasil, err
	}

	for _, r := range rows {
		hasil.tambah(r.Status, r.Jumlah)
	}

	return hasil, nil
}

// GetKetidakhadiranRombel = rekap semua siswa 1 rombel (key = NIS), 1 query.
func GetKetidakhadiranRombel(idRombel int, namaTahun, semester string) (map[string]Ketidakhadiran, error) {
	awal, akhir, err := RentangSemester(namaTahun, semester)
	if err != nil {
		return nil, err
	}

	rows := []struct {
		Nis    string `db:"nis"`
		Status string `db:"status"`
		Jumlah int    `db:"jumlah"`
	}{}

	err = database.DB.Select(&rows, `
		SELECT a.nis AS nis, a.status AS status, COUNT(*) AS jumlah
		FROM absensi a
		JOIN pesertadidik p ON p.nis = a.nis
		WHERE p.id_rombel = ? AND a.tanggal BETWEEN ? AND ?
		GROUP BY a.nis, a.status
	`, idRombel, awal, akhir)
	if err != nil {
		return nil, err
	}

	hasil := map[string]Ketidakhadiran{}
	for _, r := range rows {
		k := hasil[r.Nis]
		k.tambah(r.Status, r.Jumlah)
		hasil[r.Nis] = k
	}

	return hasil, nil
}
