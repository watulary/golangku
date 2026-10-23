package bukuinduk

import "github.com/go-pdf/fpdf"

type petunjukItem struct {
	Nomor string
	Teks  string
}

var petunjukUmum = []petunjukItem{
	{"1.", "Buku induk register ialah buku yang berisi data pribadi dan data perkembangan pendidikan peserta didik selama bersekolah di suatu sekolah, dan berlaku sejak peserta didik tersebut diterima di sekolah yang bersangkutan baik sebagai peserta didik baru maupun sebagai peserta didik pindahan."},
	{"2.", "Buku induk register diisi pada saat peserta didik diterima dan selanjutnya diikuti pencatatan perkembangan peserta didik sesuai dengan butir-butir dalam lembaran buku induk register ini."},
	{"3.", "Pengisian dilakukan dengan menggunakan tinta hitam bukan menggunakan bolpoin."},
	{"4.", "Hanya diisikan data yang dapat dijaminkan kebenarannya."},
	{"5.", "Dalam pengisian tidak dibenarkan ada coretan."},
	{"6.", "Pengelolaan buku induk register di tandatangani oleh kepala sekolah."},
}

type petunjukKhususSection struct {
	Judul string
	Poin  []petunjukItem
}

var petunjukKhusus = []petunjukKhususSection{
	{Judul: "A. Keterangan Pribadi", Poin: []petunjukItem{
		{"1.", "Nama lengkap peserta didik. Nama lengkap dimulai dengan nama peserta didik, nama keluarga atau marga (kalau ada), kemudian gelar. Contoh: nama anak Alzarefa, nama ayah Reva, gelar Doktorandus (Drs). Nama anak Nesya, nama keluarga Santoso, ditulis Nesya Santoso. Nama anak Khansa, nama ayah Santoso, gelar Datuk Mangkuto, ditulis Khansa Santoso Datuk Mangkuto."},
		{"2.", "Nama panggilan ditanyakan kepada siswa yang bersangkutan."},
		{"3.", "Jenis kelamin, ditulis laki-laki atau perempuan."},
		{"4.", "Tempat lahir ditulis nama ibu kota/nama kabupaten, tanggal lahir ditulis sesuai akte kelahiran atau kenal lahir, misalnya 12 Januari 2004."},
		{"5.", "Agama, ditulis sesuai agama peserta didik: Islam, Protestan, Katolik, Hindu, Budha, atau Konghucu."},
		{"6.", "Status dalam keluarga: anak kandung, anak tiri, anak angkat."},
		{"7.", "Anak keberapa, ditulis kesatu, kedua, ketiga, dan sebagainya."},
		{"8.", "Jumlah saudara: kandung (seayah-seibu), tiri (lain ayah dan/atau ibu), angkat (diangkat keluarga lain)."},
		{"9.", "Anak yatim/piatu/yatim-piatu, cukup jelas."},
		{"10.", "Bahasa sehari-hari di rumah, cukup jelas."},
		{"11.", "Penempatan pas foto 3x4 diisi silang sesuai nomor urut, setiap siswa membubuhi tanda tangan di bawah fotonya & stempel sekolah saat masuk dan meninggalkan sekolah."},
	}},
	{Judul: "B. Keterangan Tempat Tinggal Peserta Didik", Poin: []petunjukItem{
		{"12.", "Alamat diisi lengkap: jalan/gang, nomor rumah, RT/RW, kelurahan/desa, kecamatan, nomor telepon/HP."},
		{"13.", "Bila tinggal pada wali/menumpang, jelaskan hubungannya; bila di asrama, sebutkan nama asrama & badan/yayasan penyelenggara."},
		{"14.", "Jarak dari tempat tinggal ke sekolah, dalam satuan meter/km."},
		{"15.", "Moda ke sekolah: sepeda, sepeda motor, bis, angkot, becak, jalan kaki, dll."},
	}},
	{Judul: "C. Keterangan Jasmani Dan Kesehatan Peserta Didik", Poin: []petunjukItem{
		{"16.", "Berat badan diisi saat diterima & meninggalkan sekolah (kg)."},
		{"17.", "Tinggi badan diisi saat diterima & meninggalkan sekolah (cm)."},
		{"18.", "Golongan darah diisi berdasarkan keterangan PMI."},
		{"19.", "Penyakit yang pernah diderita, misalnya TBC, cacar, dll."},
		{"20.", "Cacat/kelainan jasmaniah (gagap, rabun jauh, mata juling, buta warna, kurang lengkap panca indra, pincang, dll), dikonsultasikan dgn guru BP."},
	}},
	{Judul: "D. Keterangan Pendidikan Sebelumnya", Poin: []petunjukItem{
		{"21.", "Nama sekolah sebelum masuk ke sekolah ini."},
		{"22.", "Nama sekolah sebelum pindah ke sekolah ini."},
	}},
	{Judul: "E. Keterangan Tentang Orang Tua Kandung", Poin: []petunjukItem{{"23.", "Cukup jelas."}}},
	{Judul: "F. Keterangan Tentang Wali", Poin: []petunjukItem{{"24.", "Cukup jelas."}}},
	{Judul: "G. Keterangan Intelegensi, Kepribadian Dan Prestasi Siswa", Poin: []petunjukItem{
		{"25.", "Diisi berdasarkan hasil tes."},
		{"26.", "Cukup jelas."},
		{"27.", "Diisi hasil prestasi yang pernah dicapai siswa."},
		{"28.", "Cukup jelas."},
	}},
	{Judul: "H. Keterangan Kehadiran Siswa", Poin: []petunjukItem{
		{"29.", "Kolom diisi berurutan: kelas, jumlah hadir, persentase hadir per semester, jumlah sakit, izin, alpa, total tidak hadir, persentase tidak hadir, jumlah hari efektif. Contoh: Hadir 230 hari (95,5%), Sakit 4 hari (1,7%), Izin 3 hari (1,25%), Alpa 3 hari (1,25%)."},
	}},
	{Judul: "I. Meninggalkan Sekolah", Poin: []petunjukItem{
		{"30.", "Cukup jelas."}, {"31.", "Cukup jelas."}, {"32.", "Cukup jelas."},
	}},
	{Judul: "J. Lain-Lain", Poin: []petunjukItem{{"33.", "Diisi data penting selama peserta didik bersekolah di sekolah ini."}}},
	{Judul: "K. Salinan Laporan Hasil Pencapaian Kompetensi Peserta Didik", Poin: []petunjukItem{}},
}

// renderNumberedItem menulis satu poin bernomor dgn hanging indent:
// nomor di kolom kiri tetap (numWidth), teks di kolom kanan yg wrap
// otomatis, baris lanjutan sejajar dgn AWAL TEKS (bukan sejajar nomor).
func renderNumberedItem(pdf *fpdf.Fpdf, item petunjukItem, numWidth, fullWidth, lineHeight float64) {
	x := pdf.GetX()
	y := pdf.GetY()
	pdf.SetXY(x, y)
	pdf.CellFormat(numWidth, lineHeight, item.Nomor, "", 0, "L", false, 0, "")
	pdf.SetXY(x+numWidth, y)
	pdf.MultiCell(fullWidth-numWidth, lineHeight, item.Teks, "", "L", false)
	pdf.SetX(x)
}

func renderPetunjukPages(pdf *fpdf.Fpdf) {
	pdf.AddPage()
	centerText(pdf, "III. PETUNJUK PENGISIAN BUKU INDUK", 13, true)
	pdf.Ln(4)

	fullWidth := pageWidth - marginL - marginR
	numWidth := 8.0
	lineHeight := 5.0

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 6, "Petunjuk Umum", "", 1, "L", false, 0, "")
	pdf.Ln(2)
	pdf.SetFont("Arial", "", 10)
	for _, item := range petunjukUmum {
		renderNumberedItem(pdf, item, numWidth, fullWidth, lineHeight)
		pdf.Ln(1.5)
	}

	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 6, "Petunjuk Khusus", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	for _, section := range petunjukKhusus {
		pdf.SetFont("Arial", "B", 10)
		pdf.MultiCell(fullWidth, lineHeight, section.Judul, "", "L", false)
		pdf.SetFont("Arial", "", 10)
		for _, item := range section.Poin {
			renderNumberedItem(pdf, item, numWidth, fullWidth, lineHeight)
			pdf.Ln(1)
		}
		pdf.Ln(3)
	}
}
