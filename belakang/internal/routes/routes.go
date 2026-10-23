package routes

import (
	"sekolah/internal/auth"
	"sekolah/internal/middleware"
	"sekolah/internal/modules/absensi"
	"sekolah/internal/modules/berita"
	"sekolah/internal/modules/bukuinduk"
	"sekolah/internal/modules/ekskul"
	"sekolah/internal/modules/fotobanner"
	"sekolah/internal/modules/fotoekskul"
	"sekolah/internal/modules/guruwali"
	"sekolah/internal/modules/kategoriberita"
	"sekolah/internal/modules/kelas"
	"sekolah/internal/modules/kelulusan"
	"sekolah/internal/modules/kenaikan"
	"sekolah/internal/modules/mapel"
	"sekolah/internal/modules/nilai"
	"sekolah/internal/modules/pegawai"
	"sekolah/internal/modules/pengguna"
	"sekolah/internal/modules/pengumuman"
	"sekolah/internal/modules/pesertadidik"
	"sekolah/internal/modules/rombel"
	"sekolah/internal/modules/sambutan"
	"sekolah/internal/modules/sejarah"
	"sekolah/internal/modules/setting"
	"sekolah/internal/modules/spmb"
	"sekolah/internal/modules/statistik"
	"sekolah/internal/modules/strukturakademik"
	"sekolah/internal/modules/strukturkomite"
	"sekolah/internal/modules/strukturtatausaha"
	"sekolah/internal/modules/suratmasuk"
	"sekolah/internal/modules/tahunpelajaran"
	"sekolah/internal/modules/tiptap"
	"sekolah/internal/modules/tka"
	"sekolah/internal/modules/unggahdokumen"
	"sekolah/internal/modules/video"
	"sekolah/internal/modules/visimisi"
	"sekolah/internal/modules/walikelas"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {

	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"app":     "Sekolah",
			"status":  "ok",
		})
	})

	// =====================================================
	// PUBLIC API
	// =====================================================

	app.Get("/api/public/home", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
		})
	})

	app.Get("/api/public/setting", setting.Get)
	app.Get("/api/public/favicon", setting.Favicon)
	app.Get("/api/public/fotobanner", fotobanner.Get)
	app.Get("/api/public/sambutan", sambutan.Get)
	app.Post("/api/public/sambutan/hits", sambutan.IncrementHits)

	app.Get("/api/public/berita", berita.GetLatestHandler)
	app.Get("/api/public/berita/list", berita.GetAllHandler)
	app.Get("/api/public/berita/detil/:slug", berita.GetDetailHandler)
	app.Post("/api/public/berita/hits/:id", berita.IncrementHitsHandler)

	app.Get("/api/public/pengumuman", pengumuman.GetLatestHandler)
	app.Get("/api/public/pengumuman/list", pengumuman.GetAllHandler)
	app.Get("/api/public/pengumuman/detil/:slug", pengumuman.GetDetailHandler)
	app.Post("/api/public/pengumuman/hits/:id", pengumuman.IncrementHitsHandler)

	app.Get("/api/public/fotoekskul", fotoekskul.Get)
	app.Get("/api/public/fotoekskul/list", fotoekskul.GetAllHandler)

	app.Get("/api/public/video", video.Get)
	app.Get("/api/public/video/list", video.GetAllHandler)

	app.Get("/api/public/statistik", statistik.GetPublicHandler)
	app.Get("/api/public/statistik/rekap-rombel", statistik.GetRekapRombelHandler)

	app.Get("/api/public/sejarah", sejarah.GetHandler)
	app.Post("/api/public/sejarah/hits", sejarah.IncrementHitsHandler)

	app.Get("/api/public/visimisi", visimisi.GetHandler)
	app.Post("/api/public/visimisi/hits", visimisi.IncrementHitsHandler)

	app.Get("/api/public/struktur-akademik", strukturakademik.GetHandler)
	app.Post("/api/public/struktur-akademik/hits", strukturakademik.IncrementHitsHandler)

	app.Get("/api/public/struktur-tatausaha", strukturtatausaha.GetHandler)
	app.Post("/api/public/struktur-tatausaha/hits", strukturtatausaha.IncrementHitsHandler)

	app.Get("/api/public/struktur-komite", strukturkomite.GetHandler)
	app.Post("/api/public/struktur-komite/hits", strukturkomite.IncrementHitsHandler)

	app.Get("/api/public/informasi-dokumen", berita.GetInformasiDokumenHandler)
	app.Get("/api/public/informasi-berkala", berita.GetInformasiBerkalaHandler)
	app.Get("/api/public/informasi-setiap-saat", berita.GetInformasiSetiapSaatHandler)
	app.Get("/api/public/informasi-serta-merta", berita.GetInformasiSertaMertaHandler)
	app.Get("/api/public/informasi-kontak", berita.GetInformasiKontakHandler)

	app.Get("/api/public/spmb", spmb.GetHandler)
	app.Get("/api/public/unduh", unggahdokumen.GetHandler)

	app.Post("/api/public/login", auth.Login)

	app.Get(
		"/api/public/me",
		middleware.Auth(),
		auth.Me,
	)

	// =====================================================
	// PRIVATE API (JWT)
	// =====================================================

	api := app.Group("/api", middleware.Auth())

	// Statistik dashboard (beda dari /api/public/statistik yg tanpa filter
	// utk halaman publik): versi ini butuh claims utk filter per-pengguna
	// (Super Admin lihat semua, level lain cuma lihat konten miliknya).
	api.Get("/statistik", statistik.GetHandler)

	api.Post(
		"/tiptap/upload",
		middleware.RequireRole("Admin", "Staf"),
		tiptap.UploadImage,
	)

	api.Get(
		"/pengguna",
		middleware.RequireRole("Admin"),
		pengguna.ListHandler,
	)

	api.Post(
		"/pengguna",
		middleware.RequireRole("Admin"),
		pengguna.CreateHandler,
	)

	api.Post(
		"/pengguna/bulk/activate",
		middleware.RequireRole("Admin"),
		pengguna.BulkActivateHandler,
	)

	api.Post(
		"/pengguna/bulk/deactivate",
		middleware.RequireRole("Admin"),
		pengguna.BulkDeactivateHandler,
	)

	api.Post(
		"/pengguna/bulk/reset-password",
		middleware.RequireRole("Admin"),
		pengguna.BulkResetPasswordHandler,
	)

	api.Post(
		"/pengguna/bulk/delete",
		middleware.RequireRole("Admin"),
		pengguna.BulkDeleteHandler,
	)

	api.Post(
		"/pengguna/change-password",
		pengguna.ChangePasswordHandler,
	)

	api.Put(
		"/pengguna/foto",
		pengguna.UpdateFotoPenggunaHandler,
	)

	api.Delete(
		"/pengguna/foto",
		pengguna.DeleteFotoPenggunaHandler,
	)

	api.Get(
		"/pengguna/export-excel",
		middleware.RequireRole("Admin"),
		pengguna.ExportExcelHandler,
	)

	api.Get(
		"/pengguna/template",
		middleware.RequireRole("Admin"),
		pengguna.DownloadTemplateHandler,
	)

	api.Post(
		"/pengguna/upload",
		middleware.RequireRole("Admin"),
		pengguna.UploadHandler,
	)

	api.Get(
		"/pengguna/check-duplicate",
		middleware.RequireRole("Admin"),
		pengguna.CheckDuplicateHandler,
	)

	api.Get(
		"/pengguna/:id",
		middleware.RequireRole("Admin"),
		pengguna.GetPenggunaHandler,
	)

	api.Put(
		"/pengguna/:id",
		middleware.RequireRole("Admin"),
		pengguna.UpdatePenggunaHandler,
	)

	api.Put(
		"/pengguna/:id/foto",
		middleware.RequireRole("Admin"),
		pengguna.UpdateFotoPenggunaByIDHandler,
	)

	api.Delete(
		"/pengguna/:id/foto",
		middleware.RequireRole("Admin"),
		pengguna.DeleteFotoPenggunaByIDHandler,
	)

	api.Get(
		"/setting",
		middleware.RequireRole("Admin", "Staf"),
		setting.GetDashboardHandler,
	)

	api.Put(
		"/setting",
		middleware.RequireRole("Admin", "Staf"),
		setting.UpdateDashboardHandler,
	)

	api.Get(
		"/sambutan",
		middleware.RequireRole("Admin", "Staf"),
		sambutan.GetDashboardHandler,
	)

	api.Put(
		"/sambutan",
		middleware.RequireRole("Admin", "Staf"),
		sambutan.UpdateDashboardHandler,
	)

	api.Get(
		"/sejarah",
		middleware.RequireRole("Admin", "Staf"),
		sejarah.GetDashboardHandler,
	)

	api.Put(
		"/sejarah",
		middleware.RequireRole("Admin", "Staf"),
		sejarah.UpdateDashboardHandler,
	)

	api.Get(
		"/visimisi",
		middleware.RequireRole("Admin", "Staf"),
		visimisi.GetDashboardHandler,
	)

	api.Put(
		"/visimisi",
		middleware.RequireRole("Admin", "Staf"),
		visimisi.UpdateDashboardHandler,
	)

	api.Get(
		"/struktur-akademik",
		middleware.RequireRole("Admin", "Staf"),
		strukturakademik.GetDashboardHandler,
	)

	api.Put(
		"/struktur-akademik",
		middleware.RequireRole("Admin", "Staf"),
		strukturakademik.UpdateDashboardHandler,
	)

	api.Get(
		"/struktur-tatausaha",
		middleware.RequireRole("Admin", "Staf"),
		strukturtatausaha.GetDashboardHandler,
	)

	api.Put(
		"/struktur-tatausaha",
		middleware.RequireRole("Admin", "Staf"),
		strukturtatausaha.UpdateDashboardHandler,
	)

	api.Get(
		"/struktur-komite",
		middleware.RequireRole("Admin", "Staf"),
		strukturkomite.GetDashboardHandler,
	)

	api.Put(
		"/struktur-komite",
		middleware.RequireRole("Admin", "Staf"),
		strukturkomite.UpdateDashboardHandler,
	)

	api.Get(
		"/berita",
		middleware.RequireRole("Admin", "Staf"),
		berita.GetAdminListHandler,
	)

	api.Get(
		"/berita/:id",
		middleware.RequireRole("Admin", "Staf"),
		berita.GetByIDHandler,
	)

	api.Post(
		"/berita",
		middleware.RequireRole("Admin", "Staf"),
		berita.CreateHandler,
	)

	api.Put(
		"/berita/:id",
		middleware.RequireRole("Admin", "Staf"),
		berita.UpdateHandler,
	)

	api.Delete(
		"/berita/:id",
		middleware.RequireRole("Admin", "Staf"),
		berita.DeleteHandler,
	)

	api.Get(
		"/kategori-berita",
		middleware.RequireRole("Admin", "Staf"),
		kategoriberita.GetAllHandler,
	)

	api.Get(
		"/pengumuman",
		middleware.RequireRole("Admin", "Staf"),
		pengumuman.GetAdminListHandler,
	)

	api.Get(
		"/pengumuman/:id",
		middleware.RequireRole("Admin", "Staf"),
		pengumuman.GetByIDHandler,
	)

	api.Post(
		"/pengumuman",
		middleware.RequireRole("Admin", "Staf"),
		pengumuman.CreateHandler,
	)

	api.Put(
		"/pengumuman/:id",
		middleware.RequireRole("Admin", "Staf"),
		pengumuman.UpdateHandler,
	)

	api.Delete(
		"/pengumuman/:id",
		middleware.RequireRole("Admin", "Staf"),
		pengumuman.DeleteHandler,
	)

	api.Get(
		"/fotoekskul",
		middleware.RequireRole("Admin", "Staf"),
		fotoekskul.GetAdminListHandler,
	)

	api.Get(
		"/fotoekskul/:id",
		middleware.RequireRole("Admin", "Staf"),
		fotoekskul.GetByIDHandler,
	)

	api.Post(
		"/fotoekskul",
		middleware.RequireRole("Admin", "Staf"),
		fotoekskul.CreateHandler,
	)

	api.Put(
		"/fotoekskul/:id",
		middleware.RequireRole("Admin", "Staf"),
		fotoekskul.UpdateHandler,
	)

	api.Delete(
		"/fotoekskul/:id",
		middleware.RequireRole("Admin", "Staf"),
		fotoekskul.DeleteHandler,
	)

	api.Get(
		"/fotobanner",
		middleware.RequireRole("Admin", "Staf"),
		fotobanner.GetAdminListHandler,
	)

	api.Get(
		"/fotobanner/:id",
		middleware.RequireRole("Admin", "Staf"),
		fotobanner.GetByIDHandler,
	)

	api.Post(
		"/fotobanner",
		middleware.RequireRole("Admin", "Staf"),
		fotobanner.CreateHandler,
	)

	api.Put(
		"/fotobanner/:id",
		middleware.RequireRole("Admin", "Staf"),
		fotobanner.UpdateHandler,
	)

	api.Delete(
		"/fotobanner/:id",
		middleware.RequireRole("Admin", "Staf"),
		fotobanner.DeleteHandler,
	)

	api.Get(
		"/video",
		middleware.RequireRole("Admin", "Staf"),
		video.GetAdminListHandler,
	)

	api.Get(
		"/video/:id",
		middleware.RequireRole("Admin", "Staf"),
		video.GetByIDHandler,
	)

	api.Post(
		"/video",
		middleware.RequireRole("Admin", "Staf"),
		video.CreateHandler,
	)

	api.Put(
		"/video/:id",
		middleware.RequireRole("Admin", "Staf"),
		video.UpdateHandler,
	)

	api.Delete(
		"/video/:id",
		middleware.RequireRole("Admin", "Staf"),
		video.DeleteHandler,
	)

	api.Get(
		"/suratmasuk",
		middleware.RequireRole("Admin", "Staf"),
		suratmasuk.GetAdminListHandler,
	)

	api.Get(
		"/suratmasuk/:id",
		middleware.RequireRole("Admin", "Staf"),
		suratmasuk.GetByIDHandler,
	)

	api.Post(
		"/suratmasuk",
		middleware.RequireRole("Admin", "Staf"),
		suratmasuk.CreateHandler,
	)

	api.Put(
		"/suratmasuk/:id",
		middleware.RequireRole("Admin", "Staf"),
		suratmasuk.UpdateHandler,
	)

	api.Delete(
		"/suratmasuk/:id",
		middleware.RequireRole("Admin", "Staf"),
		suratmasuk.DeleteHandler,
	)

	api.Get(
		"/unggahdokumen",
		middleware.RequireRole("Admin", "Staf"),
		unggahdokumen.GetAdminListHandler,
	)

	api.Get(
		"/unggahdokumen/:id",
		middleware.RequireRole("Admin", "Staf"),
		unggahdokumen.GetByIDHandler,
	)

	api.Post(
		"/unggahdokumen",
		middleware.RequireRole("Admin", "Staf"),
		unggahdokumen.CreateHandler,
	)

	api.Put(
		"/unggahdokumen/:id",
		middleware.RequireRole("Admin", "Staf"),
		unggahdokumen.UpdateHandler,
	)

	api.Delete(
		"/unggahdokumen/:id",
		middleware.RequireRole("Admin", "Staf"),
		unggahdokumen.DeleteHandler,
	)

	api.Get(
		"/tahunpelajaran",
		middleware.RequireRole("Admin", "Staf", "Pegawai", "Pesertadidik"),
		tahunpelajaran.GetAdminListHandler,
	)

	api.Get(
		"/tahunpelajaran/:id",
		middleware.RequireRole("Admin", "Staf"),
		tahunpelajaran.GetByIDHandler,
	)

	api.Post(
		"/tahunpelajaran",
		middleware.RequireRole("Admin", "Staf"),
		tahunpelajaran.CreateHandler,
	)

	api.Put(
		"/tahunpelajaran/:id",
		middleware.RequireRole("Admin", "Staf"),
		tahunpelajaran.UpdateHandler,
	)

	api.Delete(
		"/tahunpelajaran/:id",
		middleware.RequireRole("Admin", "Staf"),
		tahunpelajaran.DeleteHandler,
	)

	api.Put(
		"/tahunpelajaran/:id/aktifkan",
		middleware.RequireRole("Admin", "Staf"),
		tahunpelajaran.SetActiveHandler,
	)

	api.Get(
		"/mapel",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		mapel.GetAdminListHandler,
	)

	// route statis "/mapel/induk" HARUS di atas "/mapel/:id"
	api.Get(
		"/mapel/induk",
		middleware.RequireRole("Admin", "Staf"),
		mapel.GetIndukHandler,
	)

	api.Get(
		"/mapel/dinilai",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		mapel.GetDinilaiHandler,
	)

	api.Get("/nilai/daftar",
		middleware.RequireRole("Admin", "Staf"),
		nilai.GetDaftarHandler,
	)
	api.Get("/nilai/template",
		middleware.RequireRole("Admin", "Staf"),
		nilai.TemplateHandler,
	)
	api.Post("/nilai/bulk",
		middleware.RequireRole("Admin", "Staf"),
		nilai.BulkHandler,
	)
	api.Post("/nilai/import",
		middleware.RequireRole("Admin", "Staf"),
		nilai.ImportHandler,
	)
	api.Get("/nilai/rekap",
		middleware.RequireRole("Admin", "Staf"),
		nilai.RekapHandler,
	)

	api.Get("/nilai-ekskul/ekskul",
		middleware.RequireRole("Admin", "Staf"),
		nilai.GetEkskulOpsiHandler,
	)
	api.Get("/nilai-ekskul/daftar",
		middleware.RequireRole("Admin", "Staf"),
		nilai.GetDaftarEkskulHandler,
	)
	api.Get("/nilai-ekskul/template",
		middleware.RequireRole("Admin", "Staf"),
		nilai.TemplateEkskulHandler,
	)
	api.Post("/nilai-ekskul/bulk",
		middleware.RequireRole("Admin", "Staf"),
		nilai.BulkEkskulHandler,
	)
	api.Post("/nilai-ekskul/import",
		middleware.RequireRole("Admin", "Staf"),
		nilai.ImportEkskulHandler,
	)

	api.Get("/keanggotaan-ekskul/ekskul",
		middleware.RequireRole("Admin", "Staf"),
		nilai.GetEkskulOpsiHandler,
	)
	api.Get("/keanggotaan-ekskul/daftar",
		middleware.RequireRole("Admin", "Staf"),
		nilai.GetDaftarKeanggotaanHandler,
	)
	api.Get("/keanggotaan-ekskul/template",
		middleware.RequireRole("Admin", "Staf"),
		nilai.TemplateKeanggotaanHandler,
	)
	api.Post("/keanggotaan-ekskul/simpan",
		middleware.RequireRole("Admin", "Staf"),
		nilai.SimpanKeanggotaanHandler,
	)
	api.Post("/keanggotaan-ekskul/import",
		middleware.RequireRole("Admin", "Staf"),
		nilai.ImportKeanggotaanHandler,
	)

	api.Get("/prestasi/daftar",
		middleware.RequireRole("Admin", "Staf"),
		nilai.GetDaftarPrestasiHandler,
	)
	api.Get("/prestasi/template",
		middleware.RequireRole("Admin", "Staf"),
		nilai.TemplatePrestasiHandler,
	)
	api.Post("/prestasi/bulk",
		middleware.RequireRole("Admin", "Staf"),
		nilai.BulkPrestasiHandler,
	)
	api.Post("/prestasi/import",
		middleware.RequireRole("Admin", "Staf"),
		nilai.ImportPrestasiHandler,
	)

	api.Get("/nilai-ujian/daftar",
		middleware.RequireRole("Admin", "Staf"),
		nilai.GetDaftarUjianHandler,
	)
	api.Get("/nilai-ujian/template",
		middleware.RequireRole("Admin", "Staf"),
		nilai.TemplateUjianHandler,
	)
	api.Post("/nilai-ujian/bulk",
		middleware.RequireRole("Admin", "Staf"),
		nilai.BulkUjianHandler,
	)
	api.Post("/nilai-ujian/import",
		middleware.RequireRole("Admin", "Staf"),
		nilai.ImportUjianHandler,
	)

	api.Get("/keputusan-akhir/daftar",
		middleware.RequireRole("Admin", "Staf"),
		nilai.GetDaftarKeputusanHandler,
	)
	api.Get("/keputusan-akhir/template",
		middleware.RequireRole("Admin", "Staf"),
		nilai.TemplateKeputusanHandler,
	)
	api.Post("/keputusan-akhir/bulk",
		middleware.RequireRole("Admin", "Staf"),
		nilai.BulkKeputusanHandler,
	)
	api.Post("/keputusan-akhir/import",
		middleware.RequireRole("Admin", "Staf"),
		nilai.ImportKeputusanHandler,
	)

	api.Get("/keputusan-naikkelas/daftar",
		middleware.RequireRole("Admin", "Staf"),
		nilai.GetDaftarNaikKelasHandler,
	)
	api.Get("/keputusan-naikkelas/template",
		middleware.RequireRole("Admin", "Staf"),
		nilai.TemplateNaikKelasHandler,
	)
	api.Post("/keputusan-naikkelas/bulk",
		middleware.RequireRole("Admin", "Staf"),
		nilai.BulkNaikKelasHandler,
	)
	api.Post("/keputusan-naikkelas/import",
		middleware.RequireRole("Admin", "Staf"),
		nilai.ImportNaikKelasHandler,
	)

	api.Get("/p5/daftar",
		middleware.RequireRole("Admin", "Staf"),
		nilai.GetDaftarP5Handler,
	)
	api.Get("/p5/template",
		middleware.RequireRole("Admin", "Staf"),
		nilai.TemplateP5Handler,
	)
	api.Post("/p5/bulk",
		middleware.RequireRole("Admin", "Staf"),
		nilai.BulkP5Handler,
	)
	api.Post("/p5/import",
		middleware.RequireRole("Admin", "Staf"),
		nilai.ImportP5Handler,
	)

	api.Get(
		"/mapel/:id",
		middleware.RequireRole("Admin", "Staf"),
		mapel.GetByIDHandler,
	)

	api.Post(
		"/mapel",
		middleware.RequireRole("Admin", "Staf"),
		mapel.CreateHandler,
	)

	api.Put(
		"/mapel/:id",
		middleware.RequireRole("Admin", "Staf"),
		mapel.UpdateHandler,
	)

	api.Delete(
		"/mapel/:id",
		middleware.RequireRole("Admin", "Staf"),
		mapel.DeleteHandler,
	)

	api.Put(
		"/mapel/:id/status",
		middleware.RequireRole("Admin", "Staf"),
		mapel.ToggleStatusHandler,
	)

	api.Get(
		"/kelas",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		kelas.GetAdminListHandler,
	)

	api.Get(
		"/kelas/:id",
		middleware.RequireRole("Admin", "Staf"),
		kelas.GetByIDHandler,
	)

	api.Post(
		"/kelas",
		middleware.RequireRole("Admin", "Staf"),
		kelas.CreateHandler,
	)

	api.Put(
		"/kelas/:id",
		middleware.RequireRole("Admin", "Staf"),
		kelas.UpdateHandler,
	)

	api.Delete(
		"/kelas/:id",
		middleware.RequireRole("Admin", "Staf"),
		kelas.DeleteHandler,
	)

	api.Put(
		"/kelas/:id/status",
		middleware.RequireRole("Admin", "Staf"),
		kelas.ToggleStatusHandler,
	)

	api.Get(
		"/rombel",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		rombel.GetAdminListHandler,
	)

	api.Get(
		"/rombel/:id",
		middleware.RequireRole("Admin", "Staf"),
		rombel.GetByIDHandler,
	)

	api.Post(
		"/rombel",
		middleware.RequireRole("Admin", "Staf"),
		rombel.CreateHandler,
	)

	api.Put(
		"/rombel/:id",
		middleware.RequireRole("Admin", "Staf"),
		rombel.UpdateHandler,
	)

	api.Delete(
		"/rombel/:id",
		middleware.RequireRole("Admin", "Staf"),
		rombel.DeleteHandler,
	)

	api.Put(
		"/rombel/:id/status",
		middleware.RequireRole("Admin", "Staf"),
		rombel.ToggleStatusHandler,
	)

	api.Get(
		"/walikelas",
		middleware.RequireRole("Admin", "Staf"),
		walikelas.GetAdminListHandler,
	)

	api.Get(
		"/walikelas/pegawai/search",
		middleware.RequireRole("Admin", "Staf"),
		walikelas.SearchPegawaiHandler,
	)

	api.Get(
		"/walikelas/:id",
		middleware.RequireRole("Admin", "Staf"),
		walikelas.GetByIDHandler,
	)

	api.Post(
		"/walikelas",
		middleware.RequireRole("Admin", "Staf"),
		walikelas.CreateHandler,
	)

	api.Put(
		"/walikelas/:id",
		middleware.RequireRole("Admin", "Staf"),
		walikelas.UpdateHandler,
	)

	api.Delete(
		"/walikelas/:id",
		middleware.RequireRole("Admin", "Staf"),
		walikelas.DeleteHandler,
	)

	api.Get(
		"/guruwali",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		guruwali.GetAllAdminHandler,
	)

	api.Post(
		"/guruwali",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.CreateHandler,
	)

	api.Post(
		"/guruwali/bulk-delete",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.BulkDeleteHandler,
	)

	// route statis "/guruwali/pegawai/search" & "/guruwali/pesertadidik/search" HARUS di atas "/guruwali/:id"
	api.Get(
		"/guruwali/pegawai/search",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.SearchPegawaiHandler,
	)

	api.Get(
		"/guruwali/pesertadidik/search",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.SearchPesertaDidikHandler,
	)

	// route statis "/guruwali/export" & "/guruwali/anggota/export" HARUS di atas "/guruwali/:id"
	api.Get(
		"/guruwali/export",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.ExportGuruWaliHandler,
	)

	api.Get(
		"/guruwali/anggota/export",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.ExportAnggotaHandler,
	)

	api.Get(
		"/guruwali/anggota/template",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.DownloadTemplateAnggotaHandler,
	)

	api.Get(
		"/guruwali/me",
		middleware.RequireRole("Pesertadidik"),
		guruwali.GetMeHandler,
	)

	// route statis "/guruwali/anggota/..." HARUS di atas "/guruwali/:id/..."
	api.Post(
		"/guruwali/anggota/bulk-delete",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.BulkRemoveAnggotaHandler,
	)

	api.Delete(
		"/guruwali/anggota/:anggotaId",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.RemoveAnggotaHandler,
	)

	api.Get(
		"/guruwali/:id/anggota",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		guruwali.GetAnggotaHandler,
	)

	api.Post(
		"/guruwali/:id/anggota",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.AddAnggotaHandler,
	)

	// route dinamis :id HARUS di bawah semua route statis di atas
	api.Get(
		"/guruwali/:id",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.GetByIDHandler,
	)

	api.Put(
		"/guruwali/:id",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.UpdateHandler,
	)

	api.Delete(
		"/guruwali/:id",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.DeleteHandler,
	)

	// route 3 segmen (/guruwali/:id/anggota/export & /import) aman didaftarkan setelah /guruwali/:id
	api.Get(
		"/guruwali/:id/anggota/export",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.ExportAnggotaByGuruWaliHandler,
	)

	api.Post(
		"/guruwali/:id/anggota/import",
		middleware.RequireRole("Admin", "Staf"),
		guruwali.ImportAnggotaHandler,
	)

	api.Get(
		"/kategori-berita/list",
		middleware.RequireRole("Admin", "Staf"),
		kategoriberita.GetAllAdminHandler,
	)

	api.Get(
		"/kategori-berita/:id",
		middleware.RequireRole("Admin", "Staf"),
		kategoriberita.GetByIDHandler,
	)

	api.Post(
		"/kategori-berita",
		middleware.RequireRole("Admin", "Staf"),
		kategoriberita.CreateHandler,
	)

	api.Put(
		"/kategori-berita/:id",
		middleware.RequireRole("Admin", "Staf"),
		kategoriberita.UpdateHandler,
	)

	api.Delete(
		"/kategori-berita/:id",
		middleware.RequireRole("Admin", "Staf"),
		kategoriberita.DeleteHandler,
	)

	api.Get(
		"/ekskul",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		ekskul.GetAdminListHandler,
	)

	api.Get(
		"/ekskul/export",
		middleware.RequireRole("Admin", "Staf"),
		ekskul.ExportEkskulHandler,
	)

	api.Get(
		"/ekskul/anggota/export",
		middleware.RequireRole("Admin", "Staf"),
		ekskul.ExportAnggotaHandler,
	)

	api.Get(
		"/ekskul/anggota/template",
		middleware.RequireRole("Admin", "Staf"),
		ekskul.DownloadTemplateAnggotaHandler,
	)

	api.Get(
		"/ekskul/pegawai/search",
		middleware.RequireRole("Admin", "Staf"),
		ekskul.SearchPegawaiHandler,
	)

	api.Get(
		"/ekskul/pesertadidik/search",
		middleware.RequireRole("Admin", "Staf"),
		ekskul.SearchPesertaDidikHandler,
	)

	api.Get(
		"/ekskul/me",
		middleware.RequireRole("Pesertadidik"),
		ekskul.GetMeHandler,
	)

	// route dinamis :id HARUS di bawah semua route statis di atas
	api.Get(
		"/ekskul/:id",
		middleware.RequireRole("Admin", "Staf"),
		ekskul.GetByIDHandler,
	)

	api.Post(
		"/ekskul",
		middleware.RequireRole("Admin", "Staf"),
		ekskul.CreateHandler,
	)

	api.Put(
		"/ekskul/:id",
		middleware.RequireRole("Admin", "Staf"),
		ekskul.UpdateHandler,
	)

	api.Delete(
		"/ekskul/:id",
		middleware.RequireRole("Admin", "Staf"),
		ekskul.DeleteHandler,
	)

	api.Get(
		"/ekskul/:id/anggota/export",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		ekskul.ExportAnggotaByEkskulHandler,
	)

	api.Get(
		"/ekskul/:id/anggota",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		ekskul.GetAnggotaListHandler,
	)

	api.Post(
		"/ekskul/:id/anggota",
		middleware.RequireRole("Admin", "Staf"),
		ekskul.AddAnggotaHandler,
	)

	api.Post(
		"/ekskul/:id/anggota/import",
		middleware.RequireRole("Admin", "Staf"),
		ekskul.ImportAnggotaHandler,
	)

	api.Post(
		"/ekskul/anggota/bulk-delete",
		middleware.RequireRole("Admin", "Staf"),
		ekskul.BulkRemoveAnggotaHandler,
	)

	api.Delete(
		"/ekskul/anggota/:anggota_id",
		middleware.RequireRole("Admin", "Staf"),
		ekskul.RemoveAnggotaHandler,
	)

	api.Get(
		"/pegawai",
		middleware.RequireRole("Admin", "Staf"),
		pegawai.GetAdminListHandler,
	)
	api.Get(
		"/pegawai/export",
		middleware.RequireRole("Admin", "Staf"),
		pegawai.ExportHandler,
	)
	api.Get(
		"/pegawai/template",
		middleware.RequireRole("Admin", "Staf"),
		pegawai.TemplateHandler,
	)
	api.Get(
		"/pegawai/search",
		middleware.RequireRole("Admin", "Staf"),
		pegawai.SearchOptionHandler,
	)
	api.Get(
		"/pegawai/check-duplicate",
		middleware.RequireRole("Admin", "Staf"),
		pegawai.CheckDuplicateHandler,
	)
	api.Post(
		"/pegawai/import",
		middleware.RequireRole("Admin", "Staf"),
		pegawai.ImportHandler,
	)
	api.Post(
		"/pegawai",
		middleware.RequireRole("Admin", "Staf"),
		pegawai.CreateHandler,
	)
	api.Get(
		"/pegawai/me",
		middleware.RequireRole("Pegawai"),
		pegawai.GetMeHandler,
	)

	// route dinamis :nuptk HARUS di bawah semua route statis di atas
	api.Get(
		"/pegawai/:nuptk",
		middleware.RequireRole("Admin", "Staf"),
		pegawai.GetByIDHandler,
	)
	api.Put(
		"/pegawai/:nuptk",
		middleware.RequireRole("Admin", "Staf"),
		pegawai.UpdateHandler,
	)
	api.Delete(
		"/pegawai/:nuptk",
		middleware.RequireRole("Admin", "Staf"),
		pegawai.DeleteHandler,
	)
	api.Put(
		"/pegawai/:nuptk/foto",
		middleware.RequireRole("Admin", "Staf"),
		pegawai.UpdateFotoHandler,
	)
	api.Delete(
		"/pegawai/:nuptk/foto",
		middleware.RequireRole("Admin", "Staf"),
		pegawai.DeleteFotoHandler,
	)

	api.Get(
		"/pesertadidik",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.GetAdminListHandler,
	)
	api.Get(
		"/pesertadidik/mutasi",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.GetMutasiListHandler,
	)
	api.Get(
		"/pesertadidik/arsip",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.GetArsipListHandler,
	)
	api.Get(
		"/pesertadidik/arsip/export",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.ExportArsipHandler,
	)
	api.Get(
		"/pesertadidik/arsip/:id",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.GetArsipByIDHandler,
	)
	api.Get(
		"/pesertadidik/check-duplicate",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.CheckDuplicateHandler,
	)
	api.Post(
		"/pesertadidik",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.CreateHandler,
	)
	api.Get(
		"/pesertadidik/export",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.ExportHandler,
	)
	api.Get(
		"/pesertadidik/mutasi/export",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.ExportMutasiHandler,
	)
	api.Get(
		"/pesertadidik/mutasi/:nis",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.GetMutasiByIDHandler,
	)
	api.Get(
		"/pesertadidik/template",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.TemplateHandler,
	)
	api.Post(
		"/pesertadidik/import",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.ImportHandler,
	)
	api.Post(
		"/pesertadidik/bulk/delete",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.BulkDeleteHandler,
	)
	api.Post(
		"/pesertadidik/bulk/mutasi",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.BulkMutasiHandler,
	)
	api.Post(
		"/pesertadidik/bulk/arsip",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.BulkArsipHandler,
	)
	api.Get(
		"/pesertadidik/me",
		middleware.RequireRole("Pesertadidik"),
		pesertadidik.GetMeHandler,
	)

	// route dinamis :nis HARUS di bawah semua route statis di atas
	api.Get(
		"/pesertadidik/:nis",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.GetByIDHandler,
	)
	api.Put(
		"/pesertadidik/:nis",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.UpdateHandler,
	)
	api.Delete(
		"/pesertadidik/:nis",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.DeleteHandler,
	)
	api.Put(
		"/pesertadidik/:nis/foto",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.UploadFotoHandler,
	)
	api.Delete(
		"/pesertadidik/:nis/foto",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.DeleteFotoHandler,
	)
	api.Post(
		"/pesertadidik/:nis/mutasi",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.MutasiHandler,
	)
	api.Post(
		"/pesertadidik/:nis/arsip",
		middleware.RequireRole("Admin", "Staf"),
		pesertadidik.ArsipHandler,
	)
	api.Get(
		"/absensi/daftar",
		middleware.RequireRole("Admin", "Staf"),
		absensi.GetDaftarHandler,
	)
	api.Post("/absensi/bulk",
		middleware.RequireRole("Admin", "Staf"),
		absensi.BulkUpsertHandler,
	)
	api.Get("/absensi/export",
		middleware.RequireRole("Admin", "Staf"),
		absensi.ExportHandler,
	)
	api.Get("/absensi/rekap-bulanan",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		absensi.GetRekapBulananHandler,
	)
	api.Get("/absensi/rekap-bulanan/export",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		absensi.ExportRekapBulananHandler,
	)
	api.Get("/absensi/saya",
		middleware.RequireRole("Pesertadidik"),
		absensi.GetSayaHandler,
	)
	api.Get(
		"/spmb",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		spmb.GetAdminListHandler,
	)
	api.Get(
		"/spmb/template",
		middleware.RequireRole("Admin", "Staf"),
		spmb.TemplateHandler,
	)
	api.Get(
		"/spmb/export",
		middleware.RequireRole("Admin", "Staf"),
		spmb.ExportHandler,
	)
	api.Post(
		"/spmb/import",
		middleware.RequireRole("Admin", "Staf"),
		spmb.ImportHandler,
	)
	api.Post(
		"/spmb/bulk-lapor-diri",
		middleware.RequireRole("Admin", "Staf"),
		spmb.BulkUpdateLaporDiriHandler,
	)
	api.Post(
		"/spmb/bulk-delete",
		middleware.RequireRole("Admin", "Staf"),
		spmb.BulkDeleteHandler,
	)
	// route dinamis :id HARUS di bawah semua route statis di atas
	api.Get(
		"/spmb/:id",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		spmb.GetByIDHandler,
	)
	api.Post(
		"/spmb",
		middleware.RequireRole("Admin", "Staf"),
		spmb.CreateHandler,
	)
	api.Put(
		"/spmb/:id",
		middleware.RequireRole("Admin", "Staf"),
		spmb.UpdateHandler,
	)
	api.Delete(
		"/spmb/:id",
		middleware.RequireRole("Admin", "Staf"),
		spmb.DeleteHandler,
	)
	api.Get(
		"/kenaikan",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		kenaikan.GetAdminListHandler,
	)
	api.Get(
		"/kenaikan/template",
		middleware.RequireRole("Admin", "Staf"),
		kenaikan.TemplateHandler,
	)
	api.Get(
		"/kenaikan/export",
		middleware.RequireRole("Admin", "Staf"),
		kenaikan.ExportHandler,
	)
	api.Post(
		"/kenaikan/import",
		middleware.RequireRole("Admin", "Staf"),
		kenaikan.ImportHandler,
	)
	api.Post(
		"/kenaikan/bulk-delete",
		middleware.RequireRole("Admin", "Staf"),
		kenaikan.BulkDeleteHandler,
	)
	api.Post(
		"/kenaikan",
		middleware.RequireRole("Admin", "Staf"),
		kenaikan.CreateHandler,
	)
	api.Get(
		"/kenaikan/me",
		middleware.RequireRole("Pesertadidik"),
		kenaikan.GetSiswaHandler,
	)
	// route dinamis :nis HARUS di bawah semua route statis di atas
	api.Get(
		"/kenaikan/:nis",
		middleware.RequireRole("Admin", "Staf"),
		kenaikan.GetByIDHandler,
	)
	api.Put(
		"/kenaikan/:nis",
		middleware.RequireRole("Admin", "Staf"),
		kenaikan.UpdateHandler,
	)
	api.Delete(
		"/kenaikan/:nis",
		middleware.RequireRole("Admin", "Staf"),
		kenaikan.DeleteHandler,
	)
	api.Get(
		"/kenaikan-waktu",
		middleware.RequireRole("Admin", "Staf"),
		kenaikan.GetWaktuHandler,
	)
	api.Put(
		"/kenaikan-waktu",
		middleware.RequireRole("Admin", "Staf"),
		kenaikan.UpdateWaktuHandler,
	)
	api.Get(
		"/kelulusan",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		kelulusan.ListAdminHandler,
	)
	api.Get(
		"/kelulusan/waktu",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.GetWaktuHandler,
	)
	api.Put(
		"/kelulusan/waktu",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.UpdatePengumumanHandler,
	)
	api.Put(
		"/kelulusan/skl",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.UpdateSKLHandler,
	)
	api.Get(
		"/kelulusan/template",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.TemplateHandler,
	)
	api.Get(
		"/kelulusan/export",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.ExportHandler,
	)
	api.Post(
		"/kelulusan/import",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.ImportExcelHandler,
	)
	api.Post(
		"/kelulusan/bulk-delete",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.BulkDeleteHandler,
	)
	api.Get(
		"/kelulusan/foto",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.ListFotoHandler,
	)
	api.Post(
		"/kelulusan/foto/upload",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.UploadFotoZipHandler,
	)
	api.Post(
		"/kelulusan/foto/bulk-delete",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.BulkDeleteFotoHandler,
	)
	api.Post(
		"/kelulusan",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.CreateHandler,
	)
	api.Get(
		"/kelulusan/me",
		middleware.RequireRole("Pesertadidik"),
		kelulusan.GetMeHandler,
	)
	// cetak SKL: Peserta Didik (punya sendiri) ATAU Admin/Staf (siapa saja) —
	// makanya cuma middleware.Auth() (sudah ada di level "api"), guard
	// ownership per-role dicek DI DALAM handler, bukan lewat RequireRole di sini.
	api.Get(
		"/kelulusan/cetak/:nis",
		kelulusan.CetakPDFHandler,
	)
	// route dinamis :nis HARUS di bawah semua route statis di atas
	api.Get(
		"/kelulusan/:nis",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.GetOneAdminHandler,
	)
	api.Put(
		"/kelulusan/:nis",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.UpdateHandler,
	)
	api.Delete(
		"/kelulusan/:nis",
		middleware.RequireRole("Admin", "Staf"),
		kelulusan.DeleteHandler,
	)
	api.Get(
		"/tka",
		middleware.RequireRole("Admin", "Staf", "Pegawai"),
		tka.GetAllAdminHandler,
	)
	api.Get(
		"/tka/waktu",
		middleware.RequireRole("Admin", "Staf"),
		tka.GetWaktuHandler,
	)
	api.Put(
		"/tka/waktu",
		middleware.RequireRole("Admin", "Staf"),
		tka.UpdateWaktuHandler,
	)
	api.Get(
		"/tka/template",
		middleware.RequireRole("Admin", "Staf"),
		tka.TemplateHandler,
	)
	api.Get(
		"/tka/export",
		middleware.RequireRole("Admin", "Staf"),
		tka.ExportHandler,
	)
	api.Post(
		"/tka/import",
		middleware.RequireRole("Admin", "Staf"),
		tka.ImportHandler,
	)
	api.Post(
		"/tka/bulk-delete",
		middleware.RequireRole("Admin", "Staf"),
		tka.BulkDeleteHandler,
	)
	api.Get(
		"/tka/me",
		middleware.RequireRole("Pesertadidik"),
		tka.GetMeHandler,
	)
	api.Get(
		"/bukuinduk/cetak/:sumber/:id",
		middleware.RequireRole("Admin", "Staf"),
		bukuinduk.CetakSatuHandler,
	)
	api.Post(
		"/bukuinduk/cetak/:sumber/bulk",
		middleware.RequireRole("Admin", "Staf"),
		bukuinduk.CetakBulkHandler,
	)
	api.Get(
		"/bukuinduk/opsi-cetak",
		middleware.RequireRole("Admin", "Staf"),
		bukuinduk.GetOpsiCetakHandler,
	)
	api.Put(
		"/bukuinduk/opsi-cetak",
		middleware.RequireRole("Admin", "Staf"),
		bukuinduk.SaveOpsiCetakHandler,
	)

	// Wali kelas per tingkat (Buku Induk > Wali Kelas); route statis semuanya, tanpa :id
	api.Get(
		"/bukuinduk/walikelas-tingkat/pegawai",
		middleware.RequireRole("Admin", "Staf"),
		bukuinduk.SearchPegawaiWaliTingkatHandler,
	)
	api.Get(
		"/bukuinduk/walikelas-tingkat/template",
		middleware.RequireRole("Admin", "Staf"),
		bukuinduk.TemplateWaliTingkatHandler,
	)
	api.Post(
		"/bukuinduk/walikelas-tingkat/import",
		middleware.RequireRole("Admin", "Staf"),
		bukuinduk.ImportWaliTingkatHandler,
	)
	api.Put(
		"/bukuinduk/walikelas-tingkat/rombel",
		middleware.RequireRole("Admin", "Staf"),
		bukuinduk.SaveWaliTingkatRombelHandler,
	)
	api.Delete(
		"/bukuinduk/walikelas-tingkat/rombel",
		middleware.RequireRole("Admin", "Staf"),
		bukuinduk.ClearWaliTingkatRombelHandler,
	)
	api.Get(
		"/bukuinduk/walikelas-tingkat",
		middleware.RequireRole("Admin", "Staf"),
		bukuinduk.ListWaliTingkatHandler,
	)
}
