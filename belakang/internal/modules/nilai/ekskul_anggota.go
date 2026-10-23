package nilai

import "sekolah/internal/database"

// KeanggotaanSiswa = satu baris layar Keanggotaan Ekskul: siswa 1 rombel + status keanggotaannya
// pada ekskul & periode (tahun pelajaran + semester) yang dipilih.
type KeanggotaanSiswa struct {
	Nis       string `json:"nis"`
	NamaSiswa string `json:"nama_siswa"`
	Anggota   bool   `json:"anggota"`
}

// KeanggotaanSimpanRequest = simpan keanggotaan ekskul 1 rombel + 1 periode + 1 ekskul sekaligus.
// NisAnggota = daftar NIS (dari rombel ini) yang berstatus anggota; sisanya (anggota rombel ini yang
// tidak disebut) otomatis dianggap bukan anggota pada periode ini.
type KeanggotaanSimpanRequest struct {
	IDTahunPelajaran int      `json:"id_tahunpelajaran"`
	Semester         string   `json:"semester"`
	IDRombel         int      `json:"id_rombel"`
	IDEkskul         int      `json:"id_ekskul"`
	NisAnggota       []string `json:"nis_anggota"`
}

// GetDaftarKeanggotaan = siswa 1 rombel + status keanggotaan ekskul itu pada periode ini.
func GetDaftarKeanggotaan(idTP int, semester string, idRombel, idEkskul int) ([]KeanggotaanSiswa, error) {
	siswa, err := GetSiswaRombel(idRombel)
	if err != nil {
		return nil, err
	}

	anggota, err := GetKeanggotaanEkskul(idTP, semester, idEkskul)
	if err != nil {
		return nil, err
	}

	list := make([]KeanggotaanSiswa, 0, len(siswa))
	for _, s := range siswa {
		list = append(list, KeanggotaanSiswa{
			Nis:       s.Nis,
			NamaSiswa: s.NamaSiswa,
			Anggota:   anggota[s.Nis],
		})
	}

	return list, nil
}

// SaveKeanggotaanEkskul MENGGANTI seluruh keanggotaan ekskul itu pada tahun pelajaran + semester itu
// untuk siswa-siswa rombel ini dengan nisAnggota (dalam 1 transaksi). nisAnggota kosong = semua
// anggota periode itu (utk rombel ini) dihapus.
func SaveKeanggotaanEkskul(idTP int, semester string, idEkskul int, nisRombel, nisAnggota []string) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, nis := range nisRombel {
		if _, err := tx.Exec(
			`DELETE FROM ekskul_anggota WHERE nis = ? AND id_tahunpelajaran = ? AND semester = ? AND id_ekskul = ?`,
			nis, idTP, semester, idEkskul,
		); err != nil {
			return err
		}
	}

	for _, nis := range nisAnggota {
		if _, err := tx.Exec(
			`INSERT INTO ekskul_anggota (nis, id_tahunpelajaran, semester, id_ekskul) VALUES (?, ?, ?, ?)`,
			nis, idTP, semester, idEkskul,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetKeanggotaanEkskul = NIS anggota ekskul itu pada tahun pelajaran + semester itu (tbl_anggota_ekskul biasa,
// tanpa periode, sengaja TIDAK ikut di sini - itu dianggap "keanggotaan berjalan", bukan riwayat per periode).
func GetKeanggotaanEkskul(idTP int, semester string, idEkskul int) (map[string]bool, error) {
	var list []string

	err := database.DB.Select(&list,
		`SELECT nis FROM ekskul_anggota WHERE id_tahunpelajaran = ? AND semester = ? AND id_ekskul = ?`,
		idTP, semester, idEkskul)
	if err != nil {
		return nil, err
	}

	hasil := make(map[string]bool, len(list))
	for _, nis := range list {
		hasil[nis] = true
	}

	return hasil, nil
}

// ---------- untuk PDF Buku Induk ----------

// GetNamaEkskulPeriode = nama ekskul (distinct) yang PERNAH diikuti siswa ini menurut keanggotaan
// per periode (ekskul_anggota), lintas semua tahun pelajaran & semester. Dipakai supaya ekskul yang
// pernah diikuti tapi belum sempat dinilai tetap tampil sebagai baris kosong di tabel Ekstrakurikuler,
// walau siswa itu sekarang sudah bukan anggota ekskul tersebut lagi (mis. sudah pindah ekskul).
func GetNamaEkskulPeriode(nis string) ([]string, error) {
	var list []string

	err := database.DB.Select(&list, `
		SELECT DISTINCT e.nama_ekskul
		FROM ekskul_anggota a
		JOIN tbl_ekskul e ON e.id_ekskul = a.id_ekskul
		WHERE a.nis = ?
		ORDER BY e.nama_ekskul ASC
	`, nis)

	return list, err
}
