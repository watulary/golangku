import { createRouter, createWebHistory } from "vue-router";
import PublicLayout from "../layouts/PublicLayout.vue";
import Home from "../pages/public/Home.vue";
import Sambutan from "../pages/public/Sambutan.vue";
import Berita from "../pages/public/Berita.vue";
import BeritaDetil from "../pages/public/BeritaDetil.vue";
import Pengumuman from "../pages/public/Pengumuman.vue";
import PengumumanDetil from "../pages/public/PengumumanDetil.vue";
import Foto from "../pages/public/Foto.vue";
import Video from "../pages/public/Video.vue";
import Sejarah from "../pages/public/Sejarah.vue";
import VisiMisi from "../pages/public/VisiMisi.vue";
import StrukturAkademik from "../pages/public/StrukturAkademik.vue";
import StrukturTatausaha from "../pages/public/StrukturTatausaha.vue";
import StrukturKomite from "../pages/public/StrukturKomite.vue";
import InformasiDokumen from "../pages/public/InformasiDokumen.vue";
import InformasiDokumenDetil from "../pages/public/InformasiDokumenDetil.vue";
import InformasiSetiapSaat from "../pages/public/InformasiSetiapSaat.vue";
import InformasiSetiapSaatDetil from "../pages/public/InformasiSetiapSaatDetil.vue";
import InformasiBerkala from "../pages/public/InformasiBerkala.vue";
import InformasiBerkalaDetil from "../pages/public/InformasiBerkalaDetil.vue";
import InformasiSertaMerta from "../pages/public/InformasiSertaMerta.vue";
import InformasiSertaMertaDetil from "../pages/public/InformasiSertaMertaDetil.vue";
import InformasiKontak from "../pages/public/InformasiKontak.vue";
import InformasiKontakDetil from "../pages/public/InformasiKontakDetil.vue";
import Spmb from "../pages/public/Spmb.vue";
import Unduh from "../pages/public/Unduh.vue";
import Lokasi from "../pages/public/Lokasi.vue";
import Login from "../pages/public/Login.vue";
import DashboardLayout from "@/layouts/DashboardLayout.vue";
import DashboardHome from "@/pages/dashboard/DashboardHome.vue";
import PenggunaTampil from "@/pages/dashboard/pengguna/Tampil.vue";
import PenggunaTambah from "@/pages/dashboard/pengguna/Tambah.vue";
import PenggunaUbah from "@/pages/dashboard/pengguna/Ubah.vue";
import ProfilTampil from "@/pages/dashboard/profil/Tampil.vue";
import SettingWebsite from "@/pages/dashboard/setting/Website.vue";
import SettingSambutan from "@/pages/dashboard/setting/Sambutan.vue";
import SettingSejarah from "@/pages/dashboard/setting/Sejarah.vue";
import SettingVisiMisi from "@/pages/dashboard/setting/VisiMisi.vue";
import SettingStrukturAkademik from "@/pages/dashboard/setting/StrukturAkademik.vue";
import SettingStrukturTatausaha from "@/pages/dashboard/setting/StrukturTatausaha.vue";
import SettingStrukturKomite from "@/pages/dashboard/setting/StrukturKomite.vue";
import BeritaList from "@/pages/dashboard/berita/Tampil.vue";
import BeritaForm from "@/pages/dashboard/berita/Form.vue";
import PengumumanList from "@/pages/dashboard/pengumuman/Tampil.vue";
import PengumumanForm from "@/pages/dashboard/pengumuman/Form.vue";
import FotoEkskulList from "@/pages/dashboard/fotoekskul/Tampil.vue";
import FotoEkskulForm from "@/pages/dashboard/fotoekskul/Form.vue";
import FotoBannerList from "@/pages/dashboard/fotobanner/Tampil.vue";
import FotoBannerForm from "@/pages/dashboard/fotobanner/Form.vue";
import VideoList from "@/pages/dashboard/video/Tampil.vue";
import VideoForm from "@/pages/dashboard/video/Form.vue";
import SuratMasukList from "@/pages/dashboard/suratmasuk/Tampil.vue";
import SuratMasukForm from "@/pages/dashboard/suratmasuk/Form.vue";
import UnggahdokumenList from "@/pages/dashboard/unggahdokumen/Tampil.vue";
import UnggahdokumenForm from "@/pages/dashboard/unggahdokumen/Form.vue";
import TahunPelajaranList from "@/pages/dashboard/tahunpelajaran/Tampil.vue";
import TahunPelajaranForm from "@/pages/dashboard/tahunpelajaran/Form.vue";
import KelasList from "@/pages/dashboard/kelas/Tampil.vue";
import KelasForm from "@/pages/dashboard/kelas/Form.vue";
import RombelList from "@/pages/dashboard/rombel/Tampil.vue";
import RombelForm from "@/pages/dashboard/rombel/Form.vue";
import WaliKelasList from "@/pages/dashboard/walikelas/Tampil.vue";
import WaliKelasForm from "@/pages/dashboard/walikelas/Form.vue";
import MapelList from "@/pages/dashboard/mapel/Tampil.vue";
import MapelForm from "@/pages/dashboard/mapel/Form.vue";
import GuruWaliList from "@/pages/dashboard/guruwali/Tampil.vue";
import GuruWaliForm from "@/pages/dashboard/guruwali/Form.vue";
import GuruWaliAnggota from "@/pages/dashboard/guruwali/Anggota.vue";
import GuruWaliPegawai from "@/pages/dashboard/guruwali/pegawai/Tampil.vue";
import GuruWaliPesertadidik from "@/pages/dashboard/guruwali/pesertadidik/Tampil.vue";
import KategoriBeritaList from "@/pages/dashboard/kategoriberita/Tampil.vue";
import KategoriBeritaForm from "@/pages/dashboard/kategoriberita/Form.vue";
import EkskulList from "@/pages/dashboard/ekskul/Tampil.vue";
import EkskulForm from "@/pages/dashboard/ekskul/Form.vue";
import EkskulAnggota from "@/pages/dashboard/ekskul/Anggota.vue";
import EkskulPegawai from "@/pages/dashboard/ekskul/pegawai/Tampil.vue";
import EkskulPesertadidik from "@/pages/dashboard/ekskul/pesertadidik/Tampil.vue";
import PegawaiList from "@/pages/dashboard/pegawai/Tampil.vue";
import PegawaiForm from "@/pages/dashboard/pegawai/Form.vue";
import PegawaiBiodata from "@/pages/dashboard/pegawai/Biodata.vue";
import PesertaDidikList from "@/pages/dashboard/pesertadidik/Tampil.vue";
import PesertaDidikForm from "@/pages/dashboard/pesertadidik/Form.vue";
import PesertaDidikMutasi from "@/pages/dashboard/pesertadidik/Mutasi.vue";
import PesertaDidikMutasiDetil from "@/pages/dashboard/pesertadidik/MutasiDetil.vue";
import PesertaDidikArsip from "@/pages/dashboard/pesertadidik/Arsip.vue";
import PesertaDidikArsipDetil from "@/pages/dashboard/pesertadidik/ArsipDetil.vue";
import AbsensiInput from "@/pages/dashboard/absensi/Input.vue";
import NilaiInput from "@/pages/dashboard/nilai/Input.vue";
import NilaiEkskulInput from "@/pages/dashboard/nilaiekskul/Input.vue";
import KeanggotaanEkskulInput from "@/pages/dashboard/keanggotaanekskul/Input.vue";
import PrestasiInput from "@/pages/dashboard/prestasi/Input.vue";
import NilaiUjianInput from "@/pages/dashboard/nilaiujian/Input.vue";
import KeputusanAkhirInput from "@/pages/dashboard/keputusanakhir/Input.vue";
import KeputusanNaikKelasInput from "@/pages/dashboard/keputusannaikkelas/Input.vue";
import BukuIndukPengaturan from "@/pages/dashboard/bukuinduk/Pengaturan.vue";
import BukuIndukWaliKelas from "@/pages/dashboard/bukuinduk/WaliKelas.vue";
import P5Input from "@/pages/dashboard/p5/Input.vue";
import AbsensiRekap from "@/pages/dashboard/absensi/Rekap.vue";
import AbsensiRekapPegawai from "@/pages/dashboard/absensi/pegawai/Rekap.vue";
import AbsensiPesertaDidik from "@/pages/dashboard/absensi/pesertadidik/Tampil.vue";
import PesertaDidikBiodata from "@/pages/dashboard/pesertadidik/Biodata.vue";
import SpmbList from "@/pages/dashboard/spmb/Tampil.vue";
import SpmbForm from "@/pages/dashboard/spmb/Form.vue";
import SpmbPegawaiList from "@/pages/dashboard/spmb/pegawai/Tampil.vue";
import KenaikanList from "@/pages/dashboard/kenaikan/Tampil.vue";
import KenaikanForm from "@/pages/dashboard/kenaikan/Form.vue";
import KenaikanPengumuman from "@/pages/dashboard/kenaikan/Pengumuman.vue";
import KenaikanPegawai from "@/pages/dashboard/kenaikan/pegawai/Tampil.vue";
import KenaikanPesertadidik from "@/pages/dashboard/kenaikan/pesertadidik/Tampil.vue";
import KelulusanList from "@/pages/dashboard/kelulusan/Tampil.vue";
import KelulusanForm from "@/pages/dashboard/kelulusan/Form.vue";
import KelulusanPengumuman from "@/pages/dashboard/kelulusan/Pengumuman.vue";
import KelulusanGaleriFoto from "@/pages/dashboard/kelulusan/GaleriFoto.vue";
import KelulusanPegawai from "@/pages/dashboard/kelulusan/pegawai/Tampil.vue";
import KelulusanPesertadidik from "@/pages/dashboard/kelulusan/pesertadidik/Tampil.vue";
import TkaList from "@/pages/dashboard/tka/Tampil.vue";
import TkaPengumuman from "@/pages/dashboard/tka/Pengumuman.vue";
import TkaPegawai from "@/pages/dashboard/tka/pegawai/Tampil.vue";
import TkaPesertadidik from "@/pages/dashboard/tka/pesertadidik/Tampil.vue";
import NotFound from "../pages/public/NotFound.vue";
import DashboardNotFound from "@/pages/dashboard/NotFound.vue";
import { useAuthStore } from "@/stores/auth";

const routes = [
  {
    path: "/",
    component: PublicLayout,

    children: [
      {
        path: "",
        name: "home",
        component: Home,
      },

      {
        path: "sambutan",
        name: "sambutan",
        component: Sambutan,
      },

      {
        path: "berita",
        name: "berita",
        component: Berita,
      },

      {
        path: "berita/:slug",
        name: "berita-detil",
        component: BeritaDetil,
      },

      {
        path: "pengumuman",
        name: "pengumuman",
        component: Pengumuman,
      },

      {
        path: "pengumuman/:slug",
        name: "pengumuman-detil",
        component: PengumumanDetil,
      },

      {
        path: "foto",
        name: "foto",
        component: Foto,
      },

      {
        path: "video",
        name: "video",
        component: Video,
      },

      {
        path: "profil/sejarah",
        name: "sejarah",
        component: Sejarah,
      },

      {
        path: "profil/visi-misi",
        name: "visi-misi",
        component: VisiMisi,
      },

      {
        path: "profil/struktur-akademik",
        name: "struktur-akademik",
        component: StrukturAkademik,
      },

      {
        path: "profil/struktur-tatausaha",
        name: "struktur-tatausaha",
        component: StrukturTatausaha,
      },

      {
        path: "profil/struktur-komite",
        name: "struktur-komite",
        component: StrukturKomite,
      },

      {
        path: "ppid/informasi-dokumen",
        name: "informasi-dokumen",
        component: InformasiDokumen,
      },

      {
        path: "ppid/informasi-dokumen/:slug",
        name: "informasi-dokumen-detil",
        component: InformasiDokumenDetil,
      },

      {
        path: "ppid/informasi-setiap-saat",
        name: "informasi-setiap-saat",
        component: InformasiSetiapSaat,
      },

      {
        path: "ppid/informasi-setiap-saat/:slug",
        name: "informasi-setiap-saat-detil",
        component: InformasiSetiapSaatDetil,
      },

      {
        path: "ppid/informasi-berkala",
        name: "informasi-berkala",
        component: InformasiBerkala,
      },

      {
        path: "ppid/informasi-berkala/:slug",
        name: "informasi-berkala-detil",
        component: InformasiBerkalaDetil,
      },

      {
        path: "ppid/informasi-serta-merta",
        name: "informasi-serta-merta",
        component: InformasiSertaMerta,
      },

      {
        path: "ppid/informasi-serta-merta/:slug",
        name: "informasi-serta-merta-detil",
        component: InformasiSertaMertaDetil,
      },

      {
        path: "ppid/informasi-kontak",
        name: "informasi-kontak",
        component: InformasiKontak,
      },

      {
        path: "ppid/informasi-kontak/:slug",
        name: "informasi-kontak-detil",
        component: InformasiKontakDetil,
      },

      {
        path: "spmb",
        name: "spmb",
        component: Spmb,
      },

      {
        path: "unduh",
        name: "unduh",
        component: Unduh,
      },

      {
        path: "lokasi",
        name: "lokasi",
        component: Lokasi,
      },

      {
        path: ":pathMatch(.*)*",
        name: "not-found",
        component: NotFound,
      },
    ],
  },

  // Login TANPA PublicLayout
  {
    path: "/login",
    name: "login",
    component: Login,
    meta: {
      guest: true,
    },
  },

  // Dashboard
  {
    path: "/dashboard",
    component: DashboardLayout,
    meta: {
      requiresAuth: true,
      title: "Beranda",
    },
    children: [
      {
        path: "",
        name: "dashboard",
        component: DashboardHome,
        meta: {
          title: "Beranda",
        },
      },

      {
        path: "profil",
        name: "profil",
        component: ProfilTampil,
        meta: {
          title: "Profil Akun",
        },
      },

      {
        path: "pengguna",
        name: "pengguna",
        component: PenggunaTampil,
        meta: {
          title: "Pengguna",
          roles: ["Admin"],
        },
      },

      {
        path: "pengguna/tambah",
        name: "pengguna-tambah",
        component: PenggunaTambah,
        meta: {
          title: "Tambah",
          breadcrumb: [
            {
              label: "Beranda",
              to: "/dashboard",
            },
            {
              label: "Pengguna",
              to: "/dashboard/pengguna",
            },
            {
              label: "Tambah",
            },
          ],
          roles: ["Admin"],
        },
      },

      {
        path: "pengguna/ubah/:id",
        name: "pengguna-ubah",
        component: PenggunaUbah,
        meta: {
          title: "Ubah",
          breadcrumb: [
            {
              label: "Beranda",
              to: "/dashboard",
            },
            {
              label: "Pengguna",
              to: "/dashboard/pengguna",
            },
            {
              label: "Ubah",
            },
          ],
          roles: ["Admin"],
        },
      },

      {
        path: "website",
        name: "website",
        component: SettingWebsite,
        meta: {
          title: "Website",
          breadcrumb: [
            {
              label: "Beranda",
              to: "/dashboard",
            },
            {
              label: "Pengaturan",
            },
            {
              label: "Website",
            },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "sambutan",
        name: "dashboard-sambutan",
        component: SettingSambutan,
        meta: {
          title: "Sambutan",
          breadcrumb: [
            {
              label: "Beranda",
              to: "/dashboard",
            },
            {
              label: "Pengaturan",
            },
            {
              label: "Sambutan",
            },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "sejarah",
        name: "dashboard-sejarah",
        component: SettingSejarah,
        meta: {
          title: "Sejarah",
          breadcrumb: [
            {
              label: "Beranda",
              to: "/dashboard",
            },
            {
              label: "Pengaturan",
            },
            {
              label: "Sejarah",
            },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "visimisi",
        name: "dashboard-visimisi",
        component: SettingVisiMisi,
        meta: {
          title: "Visi & Misi",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Pengaturan" },
            { label: "Visi & Misi" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "struktur-akademik",
        name: "dashboard-struktur-akademik",
        component: SettingStrukturAkademik,
        meta: {
          title: "Struktur Akademik",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Pengaturan" },
            { label: "Struktur Akademik" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "struktur-tatausaha",
        name: "dashboard-struktur-tatausaha",
        component: SettingStrukturTatausaha,
        meta: {
          title: "Struktur Tata Usaha",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Pengaturan" },
            { label: "Struktur Tata Usaha" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "struktur-komite",
        name: "dashboard-struktur-komite",
        component: SettingStrukturKomite,
        meta: {
          title: "Struktur Komite",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Pengaturan" },
            { label: "Struktur Komite" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "berita",
        name: "dashboard-berita",
        component: BeritaList,
        meta: {
          title: "Berita",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Berita" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "berita/tambah",
        name: "dashboard-berita-tambah",
        component: BeritaForm,
        meta: {
          title: "Tambah Berita",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Berita", to: "/dashboard/berita" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "berita/ubah/:id",
        name: "dashboard-berita-edit",
        component: BeritaForm,
        meta: {
          title: "Ubah Berita",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Berita", to: "/dashboard/berita" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "pengumuman",
        name: "dashboard-pengumuman",
        component: PengumumanList,
        meta: {
          title: "Pengumuman",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Pengumuman" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "pengumuman/tambah",
        name: "dashboard-pengumuman-tambah",
        component: PengumumanForm,
        meta: {
          title: "Tambah Pengumuman",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Pengumuman", to: "/dashboard/pengumuman" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "pengumuman/ubah/:id",
        name: "dashboard-pengumuman-edit",
        component: PengumumanForm,
        meta: {
          title: "Ubah Pengumuman",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Pengumuman", to: "/dashboard/pengumuman" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "galeri",
        name: "dashboard-fotoekskul",
        component: FotoEkskulList,
        meta: {
          title: "Galeri",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Galeri Foto" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "galeri/tambah",
        name: "dashboard-fotoekskul-tambah",
        component: FotoEkskulForm,
        meta: {
          title: "Tambah Foto",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Galeri", to: "/dashboard/galeri" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "galeri/ubah/:id",
        name: "dashboard-fotoekskul-edit",
        component: FotoEkskulForm,
        meta: {
          title: "Ubah Foto",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Galeri", to: "/dashboard/galeri" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "banner",
        name: "dashboard-fotobanner",
        component: FotoBannerList,
        meta: {
          title: "Banner",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Banner" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "banner/tambah",
        name: "dashboard-fotobanner-tambah",
        component: FotoBannerForm,
        meta: {
          title: "Tambah Banner",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Banner", to: "/dashboard/banner" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "banner/ubah/:id",
        name: "dashboard-fotobanner-edit",
        component: FotoBannerForm,
        meta: {
          title: "Ubah Banner",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Banner", to: "/dashboard/banner" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "video",
        name: "dashboard-video",
        component: VideoList,
        meta: {
          title: "Video",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Video" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "video/tambah",
        name: "dashboard-video-tambah",
        component: VideoForm,
        meta: {
          title: "Tambah Video",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Video", to: "/dashboard/video" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "video/ubah/:id",
        name: "dashboard-video-edit",
        component: VideoForm,
        meta: {
          title: "Ubah Video",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Video", to: "/dashboard/video" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "surat-masuk",
        name: "dashboard-suratmasuk",
        component: SuratMasukList,
        meta: {
          title: "Surat Masuk",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Surat Masuk" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "surat-masuk/tambah",
        name: "dashboard-suratmasuk-tambah",
        component: SuratMasukForm,
        meta: {
          title: "Tambah Surat Masuk",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Surat Masuk", to: "/dashboard/surat-masuk" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "surat-masuk/ubah/:id",
        name: "dashboard-suratmasuk-edit",
        component: SuratMasukForm,
        meta: {
          title: "Ubah Surat Masuk",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Surat Masuk", to: "/dashboard/surat-masuk" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "dokumen",
        name: "dashboard-unggahdokumen",
        component: UnggahdokumenList,
        meta: {
          title: "Dokumen",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Unggah Dokumen" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "dokumen/tambah",
        name: "dashboard-unggahdokumen-tambah",
        component: UnggahdokumenForm,
        meta: {
          title: "Tambah Dokumen",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Dokumen", to: "/dashboard/dokumen" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "dokumen/ubah/:id",
        name: "dashboard-unggahdokumen-edit",
        component: UnggahdokumenForm,
        meta: {
          title: "Ubah Dokumen",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Konten" },
            { label: "Dokumen", to: "/dashboard/dokumen" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "tahun-pelajaran",
        name: "dashboard-tahunpelajaran",
        component: TahunPelajaranList,
        meta: {
          title: "Tahun Pelajaran",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Tahun Pelajaran" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "tahun-pelajaran/tambah",
        name: "dashboard-tahunpelajaran-tambah",
        component: TahunPelajaranForm,
        meta: {
          title: "Tambah Tahun Pelajaran",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Tahun Pelajaran", to: "/dashboard/tahun-pelajaran" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "tahun-pelajaran/ubah/:id",
        name: "dashboard-tahunpelajaran-edit",
        component: TahunPelajaranForm,
        meta: {
          title: "Ubah Tahun Pelajaran",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Tahun Pelajaran", to: "/dashboard/tahun-pelajaran" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "kelas",
        name: "dashboard-kelas",
        component: KelasList,
        meta: {
          title: "Kelas",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Kelas" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "kelas/tambah",
        name: "dashboard-kelas-tambah",
        component: KelasForm,
        meta: {
          title: "Tambah Kelas",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Kelas", to: "/dashboard/kelas" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "kelas/ubah/:id",
        name: "dashboard-kelas-edit",
        component: KelasForm,
        meta: {
          title: "Ubah Kelas",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Kelas", to: "/dashboard/kelas" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "rombel",
        name: "dashboard-rombel",
        component: RombelList,
        meta: {
          title: "Rombel",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Rombel" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "rombel/tambah",
        name: "dashboard-rombel-tambah",
        component: RombelForm,
        meta: {
          title: "Tambah Rombel",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Rombel", to: "/dashboard/rombel" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "rombel/ubah/:id",
        name: "dashboard-rombel-edit",
        component: RombelForm,
        meta: {
          title: "Ubah Rombel",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Rombel", to: "/dashboard/rombel" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "wali-kelas",
        name: "dashboard-walikelas",
        component: WaliKelasList,
        meta: {
          title: "Wali Kelas",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Wali Kelas" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "wali-kelas/tambah",
        name: "dashboard-walikelas-tambah",
        component: WaliKelasForm,
        meta: {
          title: "Tambah Wali Kelas",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Wali Kelas", to: "/dashboard/wali-kelas" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "wali-kelas/ubah/:id",
        name: "dashboard-walikelas-edit",
        component: WaliKelasForm,
        meta: {
          title: "Ubah Wali Kelas",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Wali Kelas", to: "/dashboard/wali-kelas" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "mapel",
        name: "dashboard-mapel",
        component: MapelList,
        meta: {
          title: "Mata Pelajaran",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Mata Pelajaran" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "mapel/tambah",
        name: "dashboard-mapel-tambah",
        component: MapelForm,
        meta: {
          title: "Tambah Mata Pelajaran",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Mata Pelajaran", to: "/dashboard/mapel" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "mapel/ubah/:id",
        name: "dashboard-mapel-edit",
        component: MapelForm,
        meta: {
          title: "Ubah Mata Pelajaran",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Mata Pelajaran", to: "/dashboard/mapel" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "guru-wali",
        name: "dashboard-guruwali",
        component: GuruWaliList,
        meta: {
          title: "Guru Wali",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Guru Wali" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "guru-wali/tambah",
        name: "dashboard-guruwali-tambah",
        component: GuruWaliForm,
        meta: {
          title: "Tambah Guru Wali",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Guru Wali", to: "/dashboard/guru-wali" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "guru-wali/ubah/:id",
        name: "dashboard-guruwali-edit",
        component: GuruWaliForm,
        meta: {
          title: "Ubah Guru Wali",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Guru Wali", to: "/dashboard/guru-wali" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "guru-wali/:id/anggota",
        name: "dashboard-guruwali-anggota",
        component: GuruWaliAnggota,
        meta: {
          title: "Anggota Guru Wali",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Guru Wali", to: "/dashboard/guru-wali" },
            { label: "Anggota" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "guru-wali/pegawai",
        name: "dashboard-guruwali-pegawai",
        component: GuruWaliPegawai,
        meta: {
          title: "Guru Wali",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Guru Wali" },
          ],
          roles: ["Pegawai"],
        },
      },

      {
        path: "guru-wali/pesertadidik",
        name: "dashboard-guruwali-pesertadidik",
        component: GuruWaliPesertadidik,
        meta: {
          title: "Guru Wali",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Guru Wali" },
          ],
          roles: ["Pesertadidik"],
        },
      },

      {
        path: "kategori-berita",
        name: "dashboard-kategoriberita",
        component: KategoriBeritaList,
        meta: {
          title: "Kategori Berita",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Kategori Berita" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "kategori-berita/tambah",
        name: "dashboard-kategoriberita-tambah",
        component: KategoriBeritaForm,
        meta: {
          title: "Tambah Kategori Berita",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Kategori Berita", to: "/dashboard/kategori-berita" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "kategori-berita/ubah/:id",
        name: "dashboard-kategoriberita-edit",
        component: KategoriBeritaForm,
        meta: {
          title: "Ubah Kategori Berita",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Kategori Berita", to: "/dashboard/kategori-berita" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "ekstrakurikuler",
        name: "dashboard-ekskul",
        component: EkskulList,
        meta: {
          title: "Ekskul",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Ekstrakurikuler" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "ekstrakurikuler/tambah",
        name: "dashboard-ekskul-tambah",
        component: EkskulForm,
        meta: {
          title: "Tambah Ekskul",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Ekstrakurikuler", to: "/dashboard/ekstrakurikuler" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "ekstrakurikuler/ubah/:id",
        name: "dashboard-ekskul-edit",
        component: EkskulForm,
        meta: {
          title: "Ubah Ekskul",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Ekstrakurikuler", to: "/dashboard/ekstrakurikuler" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "ekstrakurikuler/:id/anggota",
        name: "dashboard-ekskul-anggota",
        component: EkskulAnggota,
        meta: {
          title: "Anggota Ekskul",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Data Referensi" },
            { label: "Ekstrakurikuler", to: "/dashboard/ekstrakurikuler" },
            { label: "Anggota" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "ekstrakurikuler/pegawai",
        name: "dashboard-ekskul-pegawai",
        component: EkskulPegawai,
        meta: {
          title: "Ekstrakurikuler",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Ekstrakurikuler" },
          ],
          roles: ["Pegawai"],
        },
      },
      {
        path: "ekstrakurikuler/pesertadidik",
        name: "dashboard-ekskul-pesertadidik",
        component: EkskulPesertadidik,
        meta: {
          title: "Ekstrakurikuler",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Ekstrakurikuler" },
          ],
          roles: ["Pesertadidik"],
        },
      },

      {
        path: "pegawai",
        name: "dashboard-pegawai",
        component: PegawaiList,
        meta: {
          title: "Pegawai",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Pegawai" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "pegawai/biodata",
        name: "dashboard-pegawai-biodata",
        component: PegawaiBiodata,
        meta: {
          title: "Biodata Diri",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Biodata Diri" },
          ],
          roles: ["Pegawai"],
        },
      },

      {
        path: "pegawai/tambah",
        name: "dashboard-pegawai-tambah",
        component: PegawaiForm,
        meta: {
          title: "Tambah Pegawai",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Pegawai", to: "/dashboard/pegawai" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "pegawai/ubah/:nuptk",
        name: "dashboard-pegawai-edit",
        component: PegawaiForm,
        meta: {
          title: "Ubah Pegawai",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Pegawai", to: "/dashboard/pegawai" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "pesertadidik",
        name: "dashboard-pesertadidik",
        component: PesertaDidikList,
        meta: {
          title: "Peserta Didik",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Peserta Didik" },
            { label: "Aktif" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "pesertadidik/tambah",
        name: "dashboard-pesertadidik-tambah",
        component: PesertaDidikForm,
        meta: {
          title: "Tambah Peserta Didik",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Peserta Didik", to: "/dashboard/pesertadidik" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "pesertadidik/ubah/:nis",
        name: "dashboard-pesertadidik-edit",
        component: PesertaDidikForm,
        meta: {
          title: "Ubah Peserta Didik",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Peserta Didik", to: "/dashboard/pesertadidik" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "pesertadidik/pindah/:nis",
        name: "dashboard-pesertadidik-mutasi-detil",
        component: PesertaDidikMutasiDetil,
        meta: {
          title: "Detil Mutasi Peserta Didik",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Peserta Didik", to: "/dashboard/pesertadidik" },
            { label: "Pindah", to: "/dashboard/pesertadidik/pindah" },
            { label: "Detil" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "pesertadidik/pindah",
        name: "dashboard-pesertadidik-mutasi",
        component: PesertaDidikMutasi,
        meta: {
          title: "Mutasi Peserta Didik",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Peserta Didik", to: "/dashboard/pesertadidik" },
            { label: "Pindah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "pesertadidik/arsip",
        name: "dashboard-pesertadidik-arsip",
        component: PesertaDidikArsip,
        meta: {
          title: "Arsip Peserta Didik",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Peserta Didik", to: "/dashboard/pesertadidik" },
            { label: "Arsip" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "pesertadidik/arsip/:id",
        name: "dashboard-pesertadidik-arsip-detil",
        component: PesertaDidikArsipDetil,
        meta: {
          title: "Detil Arsip Peserta Didik",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Peserta Didik", to: "/dashboard/pesertadidik" },
            { label: "Arsip", to: "/dashboard/pesertadidik/arsip" },
            { label: "Detil" },
          ],
          roles: ["Admin", "Staf"],
        },
      },

      {
        path: "absensi",
        name: "dashboard-absensi",
        component: AbsensiInput,
        meta: {
          title: "Absensi",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Absensi Peserta Didik" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "nilai",
        name: "dashboard-nilai",
        component: NilaiInput,
        meta: {
          title: "Nilai Peserta Didik",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Buku Induk" },
            { label: "Nilai Mata Pelajaran" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "nilai-ekskul",
        name: "dashboard-nilai-ekskul",
        component: NilaiEkskulInput,
        meta: {
          title: "Nilai Ekstrakurikuler",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Buku Induk" },
            { label: "Nilai Ekstrakurikuler" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "keanggotaan-ekskul",
        name: "dashboard-keanggotaan-ekskul",
        component: KeanggotaanEkskulInput,
        meta: {
          title: "Keanggotaan Ekstrakurikuler",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Buku Induk" },
            { label: "Keanggotaan Ekskul" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "prestasi",
        name: "dashboard-prestasi",
        component: PrestasiInput,
        meta: {
          title: "Catatan Prestasi",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Buku Induk" },
            { label: "Catatan Prestasi" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "p5",
        name: "dashboard-p5",
        component: P5Input,
        meta: {
          title: "Projek P5",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Buku Induk" },
            { label: "Projek P5" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "nilai-ujian",
        name: "dashboard-nilai-ujian",
        component: NilaiUjianInput,
        meta: {
          title: "Nilai Ujian",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Buku Induk" },
            { label: "Nilai Ujian" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "keputusan-akhir",
        name: "dashboard-keputusan-akhir",
        component: KeputusanAkhirInput,
        meta: {
          title: "Keputusan Akhir",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Buku Induk" },
            { label: "Keputusan Akhir" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "keputusan-naikkelas",
        name: "dashboard-keputusan-naikkelas",
        component: KeputusanNaikKelasInput,
        meta: {
          title: "Keputusan Naik Kelas",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Buku Induk" },
            { label: "Keputusan Naik Kelas" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "bukuinduk-pengaturan",
        name: "dashboard-bukuinduk-pengaturan",
        component: BukuIndukPengaturan,
        meta: {
          title: "Pengaturan Cetak Buku Induk",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Buku Induk" },
            { label: "Pengaturan" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "bukuinduk-walikelas",
        name: "dashboard-bukuinduk-walikelas",
        component: BukuIndukWaliKelas,
        meta: {
          title: "Wali Kelas Per Tingkat",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Buku Induk" },
            { label: "Wali Kelas" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "absensi/rekap",
        name: "dashboard-absensi-rekap",
        component: AbsensiRekap,
        meta: {
          title: "Rekap Bulanan Absensi",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Absensi Peserta Didik", to: "/dashboard/absensi" },
            { label: "Rekap Bulanan" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "absensi/pegawai",
        name: "dashboard-absensi-rekap-pegawai",
        component: AbsensiRekapPegawai,
        meta: {
          title: "Rekap Bulanan Absensi",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Absensi" },
          ],
          roles: ["Pegawai"],
        },
      },
      {
        path: "absensi-HAPUS/pesertadidik",
        name: "dashboard-absensi-pesertadidik",
        component: AbsensiPesertaDidik,
        meta: {
          title: "Absensi",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Absensi" },
          ],
          roles: ["Pesertadidik"],
        },
      },
      {
        path: "peserta-didik/biodata",
        name: "dashboard-pesertadidik-biodata",
        component: PesertaDidikBiodata,
        meta: {
          title: "Biodata Diri",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Biodata Diri" },
          ],
          roles: ["Pesertadidik"],
        },
      },
      {
        path: "spmb",
        name: "dashboard-spmb",
        component: SpmbList,
        meta: {
          title: "SPMB",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "SPMB" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "spmb/tambah",
        name: "dashboard-spmb-tambah",
        component: SpmbForm,
        meta: {
          title: "Tambah Peserta SPMB",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "SPMB", to: "/dashboard/spmb" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "spmb/ubah/:id",
        name: "dashboard-spmb-edit",
        component: SpmbForm,
        meta: {
          title: "Ubah Peserta SPMB",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "SPMB", to: "/dashboard/spmb" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "spmb/pegawai",
        name: "dashboard-spmb-pegawai",
        component: SpmbPegawaiList,
        meta: {
          title: "SPMB",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "SPMB" },
          ],
          roles: ["Pegawai"],
        },
      },

      {
        path: "kenaikan",
        name: "dashboard-kenaikan",
        component: KenaikanList,
        meta: {
          title: "Kenaikan Kelas",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Kenaikan Kelas" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "kenaikan/tambah",
        name: "dashboard-kenaikan-tambah",
        component: KenaikanForm,
        meta: {
          title: "Tambah Data Kenaikan",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Kenaikan Kelas", to: "/dashboard/kenaikan" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "kenaikan/ubah/:nis",
        name: "dashboard-kenaikan-edit",
        component: KenaikanForm,
        meta: {
          title: "Ubah Data Kenaikan",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Kenaikan Kelas", to: "/dashboard/kenaikan" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "kenaikan/pengumuman",
        name: "dashboard-kenaikan-pengumuman",
        component: KenaikanPengumuman,
        meta: {
          title: "Pengaturan Pengumuman Kenaikan",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Kenaikan Kelas", to: "/dashboard/kenaikan" },
            { label: "Pengumuman" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "kenaikan/pegawai",
        name: "dashboard-kenaikan-pegawai",
        component: KenaikanPegawai,
        meta: {
          title: "Kenaikan Kelas",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Kenaikan Kelas" },
          ],
          roles: ["Pegawai"],
        },
      },
      {
        path: "kenaikan/pesertadidik",
        name: "dashboard-kenaikan-pesertadidik",
        component: KenaikanPesertadidik,
        meta: {
          title: "Kenaikan Kelas",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Kenaikan Kelas" },
          ],
          roles: ["Pesertadidik"],
        },
      },

      {
        path: "kelulusan",
        name: "dashboard-kelulusan",
        component: KelulusanList,
        meta: {
          title: "Kelulusan",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Kelulusan" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "kelulusan/tambah",
        name: "dashboard-kelulusan-tambah",
        component: KelulusanForm,
        meta: {
          title: "Tambah Data Kelulusan",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Kelulusan", to: "/dashboard/kelulusan" },
            { label: "Tambah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "kelulusan/ubah/:nis",
        name: "dashboard-kelulusan-edit",
        component: KelulusanForm,
        meta: {
          title: "Ubah Data Kelulusan",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Kelulusan", to: "/dashboard/kelulusan" },
            { label: "Ubah" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "kelulusan/pengumuman",
        name: "dashboard-kelulusan-pengumuman",
        component: KelulusanPengumuman,
        meta: {
          title: "Pengaturan Pengumuman & SKL",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Kelulusan", to: "/dashboard/kelulusan" },
            { label: "Pengumuman" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "kelulusan/galeri-foto",
        name: "dashboard-kelulusan-galeri-foto",
        component: KelulusanGaleriFoto,
        meta: {
          title: "Galeri Foto Kelulusan",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Kelulusan", to: "/dashboard/kelulusan" },
            { label: "Galeri Foto" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "kelulusan/pegawai",
        name: "dashboard-kelulusan-pegawai",
        component: KelulusanPegawai,
        meta: {
          title: "Kelulusan",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Kelulusan" },
          ],
          roles: ["Pegawai"],
        },
      },
      {
        path: "kelulusan/pesertadidik",
        name: "dashboard-kelulusan-pesertadidik",
        component: KelulusanPesertadidik,
        meta: {
          title: "Kelulusan",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "Kelulusan" },
          ],
          roles: ["Pesertadidik"],
        },
      },

      {
        path: "tka",
        name: "dashboard-tka",
        component: TkaList,
        meta: {
          title: "TKA",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "TKA" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "tka/pengumuman",
        name: "dashboard-tka-pengumuman",
        component: TkaPengumuman,
        meta: {
          title: "Pengaturan Pengumuman TKA",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "TKA", to: "/dashboard/tka" },
            { label: "Pengumuman" },
          ],
          roles: ["Admin", "Staf"],
        },
      },
      {
        path: "tka/pegawai",
        name: "dashboard-tka-pegawai",
        component: TkaPegawai,
        meta: {
          title: "TKA",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "TKA" },
          ],
          roles: ["Pegawai"],
        },
      },
      {
        path: "tka/pesertadidik",
        name: "dashboard-tka-pesertadidik",
        component: TkaPesertadidik,
        meta: {
          title: "TKA",
          breadcrumb: [
            { label: "Beranda", to: "/dashboard" },
            { label: "Fitur" },
            { label: "TKA" },
          ],
          roles: ["Pesertadidik"],
        },
      },

      {
        path: ":pathMatch(.*)*",
        name: "dashboard-not-found",
        component: DashboardNotFound,
        meta: {
          title: "Halaman Tidak Ditemukan",
        },
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,

  scrollBehavior() {
    return {
      top: 0,
    };
  },
});

router.beforeEach(async (to) => {
  const authStore = useAuthStore();

  // Route membutuhkan login
  if (to.meta.requiresAuth) {
    if (!authStore.token) {
      return "/login";
    }

    try {
      await authStore.fetchMe();

      // ======================================================
      // Force Change Password
      // ======================================================

      if (authStore.mustChangePassword && to.path !== "/dashboard/profil") {
        return "/dashboard/profil";
      }

      if (to.meta.roles) {
        if (!to.meta.roles.includes(authStore.role)) {
          return "/dashboard";
        }
      }

      // ======================================================
      // Peserta Didik tidak pernah melihat Beranda (DashboardHome)
      // ======================================================

      if (to.name === "dashboard" && authStore.role === "Pesertadidik") {
        return "/dashboard/peserta-didik/biodata";
      }

      return true;
    } catch {
      authStore.logout();

      return "/login";
    }
  }

  // Route hanya untuk tamu
  if (to.meta.guest) {
    if (authStore.token) {
      try {
        await authStore.fetchMe();

        return "/dashboard";
      } catch {
        authStore.logout();

        return true;
      }
    }
  }

  return true;
});

export default router;
