import {
  HomeIcon,
  UsersIcon,
  Cog6ToothIcon,
  CircleStackIcon,
  AcademicCapIcon,
  UserGroupIcon,
  BriefcaseIcon,
  Squares2X2Icon,
  BookOpenIcon,
} from "@heroicons/vue/24/outline";

/*
|--------------------------------------------------------------------------
| Roles
|--------------------------------------------------------------------------
*/

export const ROLES = {
  ADMIN: "Admin",
  STAF: "Staf",
  PEGAWAI: "Pegawai",
  PESERTA_DIDIK: "Pesertadidik",
};

const ALL_ROLES = Object.values(ROLES);

/*
|--------------------------------------------------------------------------
| Dashboard Menu
|--------------------------------------------------------------------------
*/

export const dashboardMenus = [
  /*
  |--------------------------------------------------------------------------
  | Dashboard
  |--------------------------------------------------------------------------
  */
  {
    title: "Beranda",
    icon: HomeIcon,
    to: "/dashboard",
    roles: [ROLES.ADMIN, ROLES.STAF, ROLES.PEGAWAI],
  },

  /*
  |--------------------------------------------------------------------------
  | Profil
  | (Profil Akun diakses dari header: klik avatar/nama pengguna)
  |--------------------------------------------------------------------------
  */
  {
    title: "Biodata Diri",
    icon: AcademicCapIcon,
    to: "/dashboard/pegawai/biodata",
    roles: [ROLES.PEGAWAI],
  },

  {
    title: "Biodata Diri",
    icon: UserGroupIcon,
    to: "/dashboard/peserta-didik/biodata",
    roles: [ROLES.PESERTA_DIDIK],
  },

  /*
  |--------------------------------------------------------------------------
  | Pengguna
  |--------------------------------------------------------------------------
  */
  {
    title: "Pengguna",
    icon: UsersIcon,
    to: "/dashboard/pengguna",
    roles: [ROLES.ADMIN],
  },

  /*
  |--------------------------------------------------------------------------
  | Pengaturan
  |--------------------------------------------------------------------------
  */
  {
    title: "Pengaturan",
    icon: Cog6ToothIcon,
    roles: [ROLES.ADMIN, ROLES.STAF],

    children: [
      {
        title: "Website",
        to: "/dashboard/website",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Sambutan",
        to: "/dashboard/sambutan",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Sejarah",
        to: "/dashboard/sejarah",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Visi Misi",
        to: "/dashboard/visimisi",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Struktur Akademik",
        to: "/dashboard/struktur-akademik",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Struktur Tata Usaha",
        to: "/dashboard/struktur-tatausaha",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Struktur Komite",
        to: "/dashboard/struktur-komite",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
    ],
  },

  /*
  |--------------------------------------------------------------------------
  | Data Referensi
  |--------------------------------------------------------------------------
  */
  {
    title: "Data Referensi",
    icon: CircleStackIcon,
    roles: [ROLES.ADMIN, ROLES.STAF],

    children: [
      {
        title: "Tahun Pelajaran",
        to: "/dashboard/tahun-pelajaran",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Kelas",
        to: "/dashboard/kelas",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Rombel",
        to: "/dashboard/rombel",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Wali Kelas",
        to: "/dashboard/wali-kelas",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Mata Pelajaran",
        to: "/dashboard/mapel",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Kategori Berita",
        to: "/dashboard/kategori-berita",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
    ],
  },

  /*
  |--------------------------------------------------------------------------
  | Pegawai
  |--------------------------------------------------------------------------
  */
  {
    title: "Pegawai",
    icon: AcademicCapIcon,
    roles: [ROLES.ADMIN, ROLES.STAF],

    children: [
      {
        title: "Guru & Tendik",
        to: "/dashboard/pegawai",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
    ],
  },

  /*
  |--------------------------------------------------------------------------
  | Peserta Didik
  |--------------------------------------------------------------------------
  */
  {
    title: "Peserta Didik",
    icon: UserGroupIcon,
    roles: [ROLES.ADMIN, ROLES.STAF],

    children: [
      {
        title: "Aktif",
        to: "/dashboard/pesertadidik",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Pindah / Lainnya",
        to: "/dashboard/pesertadidik/pindah",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Arsip",
        to: "/dashboard/pesertadidik/arsip",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
    ],
  },

  /*
  |--------------------------------------------------------------------------
  | Buku Induk
  |--------------------------------------------------------------------------
  */
  {
    title: "Buku Induk",
    icon: BookOpenIcon,
    roles: [ROLES.ADMIN, ROLES.STAF],

    children: [
      {
        title: "Pengaturan",
        to: "/dashboard/bukuinduk-pengaturan",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Wali Kelas",
        to: "/dashboard/bukuinduk-walikelas",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Nilai Mata Pelajaran",
        to: "/dashboard/nilai",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Nilai Ekstrakurikuler",
        to: "/dashboard/nilai-ekskul",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Keanggotaan Ekskul",
        to: "/dashboard/keanggotaan-ekskul",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Catatan Prestasi",
        to: "/dashboard/prestasi",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Projek P5",
        to: "/dashboard/p5",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Nilai Ujian",
        to: "/dashboard/nilai-ujian",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Keputusan Naik Kelas",
        to: "/dashboard/keputusan-naikkelas",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Keputusan Akhir",
        to: "/dashboard/keputusan-akhir",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
    ],
  },

  /*
  |--------------------------------------------------------------------------
  | Konten
  |--------------------------------------------------------------------------
  */
  {
    title: "Konten",
    icon: BriefcaseIcon,
    roles: [ROLES.ADMIN, ROLES.STAF],

    children: [
      {
        title: "Berita",
        to: "/dashboard/berita",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Pengumuman",
        to: "/dashboard/pengumuman",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Galeri Foto",
        to: "/dashboard/galeri",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Banner",
        to: "/dashboard/banner",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Video",
        to: "/dashboard/video",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Surat Masuk",
        to: "/dashboard/surat-masuk",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },

      {
        title: "Unggah Dokumen",
        to: "/dashboard/dokumen",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
    ],
  },

  /*
  |--------------------------------------------------------------------------
  | Fitur
  |--------------------------------------------------------------------------
  */
  {
    title: "Fitur",
    icon: Squares2X2Icon,
    roles: ALL_ROLES,

    children: [
      {
        title: "Absensi Peserta Didik",
        to: "/dashboard/absensi",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Ekstrakurikuler",
        to: "/dashboard/ekstrakurikuler",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Guru Wali",
        to: "/dashboard/guru-wali",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Absensi",
        to: "/dashboard/absensi/pegawai",
        roles: [ROLES.PEGAWAI],
      },
      // {
      //   title: "Absensi",
      //   to: "/dashboard/absensi/pesertadidik",
      //   roles: [ROLES.PESERTA_DIDIK],
      // },
      {
        title: "Ekstrakurikuler",
        to: "/dashboard/ekstrakurikuler/pegawai",
        roles: [ROLES.PEGAWAI],
      },
      {
        title: "Guru Wali",
        to: "/dashboard/guru-wali/pegawai",
        roles: [ROLES.PEGAWAI],
      },
      {
        title: "Ekstrakurikuler",
        to: "/dashboard/ekstrakurikuler/pesertadidik",
        roles: [ROLES.PESERTA_DIDIK],
      },
      {
        title: "Guru Wali",
        to: "/dashboard/guru-wali/pesertadidik",
        roles: [ROLES.PESERTA_DIDIK],
      },
      {
        title: "Kenaikan Kelas",
        to: "/dashboard/kenaikan",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Kenaikan Kelas",
        to: "/dashboard/kenaikan/pegawai",
        roles: [ROLES.PEGAWAI],
      },
      {
        title: "Kenaikan Kelas",
        to: "/dashboard/kenaikan/pesertadidik",
        roles: [ROLES.PESERTA_DIDIK],
      },

      {
        title: "Kelulusan",
        to: "/dashboard/kelulusan",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "Kelulusan",
        to: "/dashboard/kelulusan/pegawai",
        roles: [ROLES.PEGAWAI],
      },
      {
        title: "Kelulusan",
        to: "/dashboard/kelulusan/pesertadidik",
        roles: [ROLES.PESERTA_DIDIK],
      },

      {
        title: "TKA",
        to: "/dashboard/tka",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "TKA",
        to: "/dashboard/tka/pegawai",
        roles: [ROLES.PEGAWAI],
      },
      {
        title: "TKA",
        to: "/dashboard/tka/pesertadidik",
        roles: [ROLES.PESERTA_DIDIK],
      },

      {
        title: "SPMB",
        to: "/dashboard/spmb",
        roles: [ROLES.ADMIN, ROLES.STAF],
      },
      {
        title: "SPMB",
        to: "/dashboard/spmb/pegawai",
        roles: [ROLES.PEGAWAI],
      },
    ],
  },
];
