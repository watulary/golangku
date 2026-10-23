package nilai

import (
	"strings"

	"sekolah/internal/database"
)

// Jenis ujian: "Sekolah" = Nilai Ujian Sekolah (semua mapel, tabel 1 halaman 2);
// "Daerah" = Nilai Ujian Sekolah Standar Daerah (4 mapel, tabel 2 halaman 2).
const (
	UjianSekolah = "Sekolah"
	UjianDaerah  = "Daerah"
)

// Bobot Nilai Sekolah = bobotRaport * rata-rata raport + bobotUjian * nilai ujian sekolah.
// Ubah di sini kalau rumus resmi sekolah berbeda (jumlah bobot harus 1).
const (
	bobotRaport = 0.5
	bobotUjian  = 0.5
)

func isJenisUjianValid(j string) bool {
	return j == UjianSekolah || j == UjianDaerah
}

// UjianSiswa = satu baris layar input ujian.
type UjianSiswa struct {
	IDUjian    int      `db:"id_ujian" json:"id_ujian"`
	Nis        string   `db:"nis" json:"nis"`
	NamaSiswa  string   `db:"nama_siswa" json:"nama_siswa"`
	NamaKelas  string   `db:"nama_kelas" json:"nama_kelas"`
	NamaRombel string   `db:"nama_rombel" json:"nama_rombel"`
	NilaiAkhir *float64 `db:"nilai_akhir" json:"nilai_akhir"`
	Capaian    string   `db:"capaian_pembelajaran" json:"capaian_pembelajaran"`
}

type UjianBulkItem struct {
	Nis        string   `json:"nis"`
	NilaiAkhir *float64 `json:"nilai_akhir"`
	Capaian    string   `json:"capaian_pembelajaran"`
}

type UjianBulkRequest struct {
	IDTahunPelajaran int             `json:"id_tahunpelajaran"`
	Jenis            string          `json:"jenis"`
	IDRombel         int             `json:"id_rombel"`
	IDMapel          int             `json:"id_mapel"`
	IDKelasCatat     int             `json:"id_kelas_catat"`
	Items            []UjianBulkItem `json:"items"`
}

func GetDaftarUjian(idTP int, jenis string, idRombel, idMapel int) ([]UjianSiswa, error) {
	list := []UjianSiswa{}

	err := database.DB.Select(&list, `
		SELECT
			IFNULL(n.id_ujian, 0) AS id_ujian,
			p.nis AS nis,
			p.nama_siswa AS nama_siswa,
			COALESCE(k.nama_kelas, '') AS nama_kelas,
			COALESCE(r.nama_rombel, '') AS nama_rombel,
			n.nilai_akhir AS nilai_akhir,
			COALESCE(n.capaian_pembelajaran, '') AS capaian_pembelajaran
		FROM pesertadidik p
		LEFT JOIN nilai_ujian n
			ON n.nis = p.nis
			AND n.id_tahunpelajaran = ?
			AND n.jenis = ?
			AND n.id_mapel = ?
		LEFT JOIN tbl_kelas k ON k.id_kelas = p.id_kelas
		LEFT JOIN tbl_rombel r ON r.id_rombel = p.id_rombel
		WHERE p.id_rombel = ?
		ORDER BY p.nama_siswa ASC
	`, idTP, jenis, idMapel, idRombel)

	return list, err
}

func isEmptyUjianItem(it UjianBulkItem) bool {
	return it.NilaiAkhir == nil && strings.TrimSpace(it.Capaian) == ""
}

// BulkSaveUjian: baris berisi -> upsert (unik nis+tahun+jenis+mapel); baris kosong -> hapus nilai lama.
func BulkSaveUjian(idTP int, jenis string, idKelas, idMapel int, items []UjianBulkItem) (saved int, cleared int, err error) {
	tx, err := database.DB.Beginx()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	for _, it := range items {
		if isEmptyUjianItem(it) {
			if _, err := tx.Exec(
				`DELETE FROM nilai_ujian WHERE nis = ? AND id_tahunpelajaran = ? AND jenis = ? AND id_mapel = ?`,
				it.Nis, idTP, jenis, idMapel,
			); err != nil {
				return 0, 0, err
			}

			cleared++

			continue
		}

		if _, err := tx.Exec(
			`INSERT INTO nilai_ujian
				(nis, id_tahunpelajaran, jenis, id_kelas, id_mapel, nilai_akhir, capaian_pembelajaran)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				id_kelas = VALUES(id_kelas),
				nilai_akhir = VALUES(nilai_akhir),
				capaian_pembelajaran = VALUES(capaian_pembelajaran)`,
			it.Nis, idTP, jenis, idKelas, idMapel, it.NilaiAkhir, nullStr(it.Capaian),
		); err != nil {
			return 0, 0, err
		}

		saved++
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}

	return saved, cleared, nil
}

// ---------- untuk PDF Buku Induk halaman 2 ----------

// SelUjian = nilai + capaian 1 mapel.
type SelUjian struct {
	Nilai   *float64
	Capaian string
}

// DataHalaman2 = semua angka halaman 2 Data Nilai untuk 1 siswa. Semua map berkunci id_mapel (mapel daun).
type DataHalaman2 struct {
	Mapel []MapelPohon

	Raport  map[int]float64  // rata-rata nilai akhir semua semester yang pernah diinput
	Sekolah map[int]SelUjian // Nilai Ujian Sekolah
	Daerah  map[int]SelUjian // Nilai Ujian Sekolah Standar Daerah
	Gabung  map[int]float64  // Nilai Sekolah = bobotRaport*Raport + bobotUjian*Sekolah (hanya kalau keduanya ada)
	Akhir   map[int]float64  // Nilai Akhir = rata-rata(Daerah, Gabung) (hanya kalau keduanya ada)

	Keputusan *Keputusan // nil kalau belum diisi

	JumlahRaport, RataRaport   *float64
	JumlahSekolah, RataSekolah *float64
	JumlahGabung, RataGabung   *float64
}

// GetDataHalaman2 menghitung semuanya dari nilai_mapel (rata-rata raport, otomatis) dan nilai_ujian.
func GetDataHalaman2(nis string) (*DataHalaman2, error) {
	pohon, err := GetPohonMapel()
	if err != nil {
		return nil, err
	}

	d := &DataHalaman2{
		Mapel:   pohon,
		Raport:  map[int]float64{},
		Sekolah: map[int]SelUjian{},
		Daerah:  map[int]SelUjian{},
		Gabung:  map[int]float64{},
		Akhir:   map[int]float64{},
	}

	rr := []struct {
		IDMapel int     `db:"id_mapel"`
		Rata    float64 `db:"rata"`
	}{}

	if err := database.DB.Select(&rr, `
		SELECT id_mapel, AVG(nilai_akhir) AS rata
		FROM nilai_mapel
		WHERE nis = ? AND nilai_akhir IS NOT NULL
		GROUP BY id_mapel
	`, nis); err != nil {
		return nil, err
	}

	for _, r := range rr {
		d.Raport[r.IDMapel] = round2(r.Rata)
	}

	uj := []struct {
		Jenis   string   `db:"jenis"`
		IDMapel int      `db:"id_mapel"`
		Nilai   *float64 `db:"nilai_akhir"`
		Capaian string   `db:"capaian"`
	}{}

	// kalau ada beberapa tahun (mis. tinggal kelas), tahun terbaru menang (urut naik -> menimpa)
	if err := database.DB.Select(&uj, `
		SELECT u.jenis AS jenis, u.id_mapel AS id_mapel, u.nilai_akhir AS nilai_akhir,
			COALESCE(u.capaian_pembelajaran, '') AS capaian
		FROM nilai_ujian u
		JOIN tbl_tahunpelajaran tp ON tp.id_tahunpelajaran = u.id_tahunpelajaran
		WHERE u.nis = ?
		ORDER BY tp.nama_tahunpelajaran ASC
	`, nis); err != nil {
		return nil, err
	}

	for _, u := range uj {
		sel := SelUjian{Nilai: u.Nilai, Capaian: u.Capaian}
		if u.Jenis == UjianDaerah {
			d.Daerah[u.IDMapel] = sel
		} else {
			d.Sekolah[u.IDMapel] = sel
		}
	}

	for id, r := range d.Raport {
		if s, ok := d.Sekolah[id]; ok && s.Nilai != nil {
			d.Gabung[id] = round2(bobotRaport*r + bobotUjian**s.Nilai)
		}
	}

	for id, g := range d.Gabung {
		if s, ok := d.Daerah[id]; ok && s.Nilai != nil {
			d.Akhir[id] = round2((*s.Nilai + g) / 2)
		}
	}

	ujianNilai := map[int]float64{}
	for id, s := range d.Sekolah {
		if s.Nilai != nil {
			ujianNilai[id] = *s.Nilai
		}
	}

	kep, err := GetKeputusan(nis)
	if err != nil {
		return nil, err
	}

	d.Keputusan = kep

	d.JumlahRaport, d.RataRaport = jumlahRata(pohon, d.Raport)
	d.JumlahSekolah, d.RataSekolah = jumlahRata(pohon, ujianNilai)
	d.JumlahGabung, d.RataGabung = jumlahRata(pohon, d.Gabung)

	return d, nil
}

// jumlahRata = jumlah & rata-rata per mapel utama; mapel dengan sub-mapel dihitung SATU kali
// (rata-rata sub-mapelnya yang terisi), sama seperti Jumlah/Rata-rata di halaman 1. nil kalau belum ada nilai.
func jumlahRata(pohon []MapelPohon, nilai map[int]float64) (*float64, *float64) {
	jumlah, n := 0.0, 0

	for _, m := range pohon {
		if len(m.Subs) == 0 {
			if v, ok := nilai[m.ID]; ok {
				jumlah += v
				n++
			}

			continue
		}

		total, cnt := 0.0, 0

		for _, s := range m.Subs {
			if v, ok := nilai[s.ID]; ok {
				total += v
				cnt++
			}
		}

		if cnt > 0 {
			jumlah += total / float64(cnt)
			n++
		}
	}

	if n == 0 {
		return nil, nil
	}

	j := round2(jumlah)
	r := round2(jumlah / float64(n))

	return &j, &r
}
