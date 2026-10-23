<template>
  <DashboardPageHeader
    title="Biodata Diri"
    subtitle="Data lengkap peserta didik yang tercatat di sekolah."
  />

  <!-- Loading -->
  <div
    v-if="pdStore.siswaLoading && !detail"
    class="flex flex-col items-center justify-center py-20"
  >
    <div
      class="h-10 w-10 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"
    />
    <p class="mt-4 text-sm text-gray-500 dark:text-slate-400">Memuat data...</p>
  </div>

  <!-- Not found -->
  <div
    v-else-if="!detail"
    class="rounded-2xl bg-white p-10 text-center ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
  >
    <p class="text-sm text-gray-500 dark:text-slate-400">
      Data belum tersedia. Silakan hubungi bagian tata usaha.
    </p>
  </div>

  <div v-else class="space-y-6">
    <!-- Header identitas -->
    <div
      class="relative flex flex-col gap-4 overflow-hidden rounded-2xl bg-white p-6 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5 sm:flex-row sm:items-center"
    >
      <div
        class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
      ></div>
      <div
        class="group relative h-32 w-24 shrink-0 overflow-hidden rounded-lg border border-gray-200 bg-gray-50 dark:border-slate-700 dark:bg-slate-800"
        :class="{ 'cursor-pointer': detail.foto_url }"
        @click="detail.foto_url && (visibleLightbox = true)"
      >
        <img
          v-if="detail.foto_url"
          :src="detail.foto_url"
          class="h-full w-full object-cover transition duration-200 group-hover:opacity-80"
        />
        <div
          v-else
          class="flex h-full w-full items-center justify-center text-gray-300 dark:text-slate-600"
        >
          <UserIcon class="h-12 w-12" />
        </div>
        <div
          v-if="detail.foto_url"
          class="pointer-events-none absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition duration-200 group-hover:opacity-100"
        >
          <MagnifyingGlassPlusIcon class="h-6 w-6 text-white" />
        </div>
      </div>
      <div>
        <p class="text-lg font-semibold text-gray-700 dark:text-slate-200">
          {{ detail.nama_siswa }}
        </p>
        <p class="text-sm text-gray-500 dark:text-slate-400">
          NIS: {{ detail.nis }} · NISN: {{ detail.nisn || "-" }}
        </p>
        <p class="text-sm text-gray-500 dark:text-slate-400">
          {{ detail.nama_tahunpelajaran }} · {{ detail.nama_kelas }} /
          {{ detail.nama_rombel }}
        </p>
      </div>
    </div>

    <!-- Section detail (read-only) -->
    <section
      v-for="section in sections"
      :key="section.title"
      class="overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
    >
      <div
        class="flex items-center gap-3 border-b border-gray-100 px-6 py-4 dark:border-white/5"
      >
        <div
          class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/8 text-[#041983] dark:bg-white/5 dark:text-indigo-300"
        >
          <component :is="section.icon" class="h-5 w-5" />
        </div>
        <div>
          <h3 class="text-sm font-semibold text-gray-700 dark:text-slate-200">
            {{ section.title }}
          </h3>
          <p class="text-xs text-gray-400 dark:text-slate-500">
            {{ section.description }}
          </p>
        </div>
      </div>
      <dl
        class="grid grid-cols-1 gap-x-6 gap-y-4 p-6 sm:grid-cols-2 lg:grid-cols-3"
      >
        <div v-for="field in section.fields" :key="field.key">
          <dt class="text-xs font-medium text-gray-400 dark:text-slate-500">
            {{ field.label }}
          </dt>
          <dd class="mt-0.5 text-sm text-gray-700 dark:text-slate-200">
            {{ formatValue(field, detail[field.key]) }}
          </dd>
        </div>
      </dl>
    </section>
  </div>

  <VueEasyLightbox
    :visible="visibleLightbox"
    :imgs="images"
    :index="0"
    @hide="visibleLightbox = false"
  />
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import VueEasyLightbox from "vue-easy-lightbox";
import { usePesertaDidikStore } from "@/stores/pesertadidik";
import { formatDateNumeric } from "@/utils/date";
import {
  AcademicCapIcon,
  HeartIcon,
  IdentificationIcon,
  MapPinIcon,
  UserGroupIcon,
  UserIcon,
  MagnifyingGlassPlusIcon,
} from "@heroicons/vue/24/outline";

const pdStore = usePesertaDidikStore();
const detail = ref(null);
const visibleLightbox = ref(false);
const images = computed(() =>
  detail.value?.foto_url ? [detail.value.foto_url] : [],
);

// Sama persis struktur sections dengan ArsipDetil.vue/MutasiDetil.vue — kalau field
// di sana berubah, samakan juga di sini.
const sections = [
  {
    title: "Data Peserta Didik",
    description: "Identitas dan data pribadi.",
    icon: IdentificationIcon,
    fields: [
      { key: "nisn", label: "NISN" },
      { key: "no_kk", label: "No. KK" },
      { key: "nik_siswa", label: "NIK" },
      { key: "no_ujiansd", label: "No. Ujian SD" },
      { key: "no_simdik", label: "No. Simdik" },
      { key: "jenis_kelamin", label: "Jenis Kelamin" },
      { key: "tempat_lahirsiswa", label: "Tempat Lahir" },
      { key: "tanggal_lahirsiswa", label: "Tanggal Lahir", type: "date" },
      { key: "no_akte", label: "No. Akte" },
      { key: "akte_kota", label: "Kota Akte" },
      { key: "agama_siswa", label: "Agama" },
      { key: "warga_siswa", label: "Kewarganegaraan" },
      { key: "anak_ke", label: "Anak Ke" },
      { key: "saudara_kandung", label: "Saudara Kandung" },
      { key: "saudara_tiri", label: "Saudara Tiri" },
      { key: "saudara_angkat", label: "Saudara Angkat" },
      { key: "yatim_piatu", label: "Yatim/Piatu" },
      { key: "bahasa_siswa", label: "Bahasa Sehari-hari" },
      { key: "bantuan", label: "Penerima Bantuan" },
      { key: "email", label: "Email" },
      { key: "gol_darah", label: "Golongan Darah" },
      { key: "pernah_sakit", label: "Riwayat Sakit" },
      { key: "status_jasmani", label: "Status Jasmani" },
      { key: "kebutuhan_khusus", label: "Kebutuhan Khusus" },
      { key: "tinggi_badan", label: "Tinggi Badan (cm)" },
      { key: "berat_badan", label: "Berat Badan (kg)" },
      { key: "lingkar_kepala", label: "Lingkar Kepala (cm)" },
    ],
  },
  {
    title: "Alamat & Kontak",
    description: "Alamat tempat tinggal dan kontak.",
    icon: MapPinIcon,
    fields: [
      { key: "alamat_siswa", label: "Alamat" },
      { key: "rt_siswa", label: "RT" },
      { key: "rw_siswa", label: "RW" },
      { key: "kel_siswa", label: "Kelurahan/Desa" },
      { key: "kec_siswa", label: "Kecamatan" },
      { key: "kopos_siswa", label: "Kode Pos" },
      { key: "telprmh_siswa", label: "Telp Rumah" },
      { key: "hp_siswa", label: "No. HP" },
      { key: "tinggal_siswa", label: "Tinggal Bersama" },
      { key: "jarak_sekolah", label: "Jarak ke Sekolah (km)" },
      { key: "moda_transportasi", label: "Moda Transportasi" },
    ],
  },
  {
    title: "Akademik & Kelulusan",
    description: "Data asal sekolah, kelulusan, dan ijazah.",
    icon: AcademicCapIcon,
    fields: [
      { key: "lulus_sd", label: "Lulus Dari (SD)" },
      { key: "tanggal_ijazah", label: "Tanggal Ijazah", type: "date" },
      { key: "kota_ijazah", label: "Kota Ijazah" },
      { key: "nomor_ijazah", label: "Nomor Ijazah" },
      { key: "tahun_ijazah", label: "Tahun Ijazah" },
      { key: "nomor_skhun", label: "Nomor SKHUN" },
      { key: "tahun_skhun", label: "Tahun SKHUN" },
      { key: "nilai_ijazah", label: "Nilai Ijazah" },
      { key: "nilai_nem", label: "Nilai NEM" },
      { key: "lama_belajar", label: "Lama Belajar (tahun)" },
      { key: "jalur", label: "Jalur Masuk" },
      { key: "semester", label: "Semester" },
      { key: "asal_smp", label: "Asal Sekolah" },
      { key: "tanggal_masuk", label: "Tanggal Masuk", type: "date" },
    ],
  },
  {
    title: "Data Ayah",
    description: "Identitas dan data orang tua (Ayah).",
    icon: UserIcon,
    fields: [
      { key: "nama_ayah", label: "Nama" },
      { key: "nik_ayah", label: "NIK" },
      { key: "tempat_lahirayah", label: "Tempat Lahir" },
      { key: "tanggal_lahirayah", label: "Tanggal Lahir", type: "date" },
      { key: "agama_ayah", label: "Agama" },
      { key: "warga_ayah", label: "Kewarganegaraan" },
      { key: "pend_ayah", label: "Pendidikan Terakhir" },
      { key: "pek_ayah", label: "Pekerjaan" },
      { key: "gaji_ayah", label: "Penghasilan" },
      { key: "gaji_ayahktg", label: "Kategori Penghasilan" },
      { key: "status_ayah", label: "Status" },
      { key: "alamat_ayah", label: "Alamat" },
      { key: "rt_ayah", label: "RT" },
      { key: "rw_ayah", label: "RW" },
      { key: "kopos_ayah", label: "Kode Pos" },
      { key: "kel_ayah", label: "Kelurahan/Desa" },
      { key: "kec_ayah", label: "Kecamatan" },
      { key: "telp_ayah", label: "No. Telp/HP" },
    ],
  },
  {
    title: "Data Ibu",
    description: "Identitas dan data orang tua (Ibu).",
    icon: UserIcon,
    fields: [
      { key: "nama_ibu", label: "Nama" },
      { key: "nik_ibu", label: "NIK" },
      { key: "tempat_lahiribu", label: "Tempat Lahir" },
      { key: "tanggal_lahiribu", label: "Tanggal Lahir", type: "date" },
      { key: "agama_ibu", label: "Agama" },
      { key: "warga_ibu", label: "Kewarganegaraan" },
      { key: "pend_ibu", label: "Pendidikan Terakhir" },
      { key: "pek_ibu", label: "Pekerjaan" },
      { key: "gaji_ibu", label: "Penghasilan" },
      { key: "gaji_ibuktg", label: "Kategori Penghasilan" },
      { key: "status_ibu", label: "Status" },
      { key: "alamat_ibu", label: "Alamat" },
      { key: "rt_ibu", label: "RT" },
      { key: "rw_ibu", label: "RW" },
      { key: "kopos_ibu", label: "Kode Pos" },
      { key: "kel_ibu", label: "Kelurahan/Desa" },
      { key: "kec_ibu", label: "Kecamatan" },
      { key: "telp_ibu", label: "No. Telp/HP" },
    ],
  },
  {
    title: "Data Wali",
    description: "Identitas dan data wali siswa (jika ada).",
    icon: UserGroupIcon,
    fields: [
      { key: "nama_wali", label: "Nama" },
      { key: "tempat_lahirwali", label: "Tempat Lahir" },
      { key: "tanggal_lahirwali", label: "Tanggal Lahir", type: "date" },
      { key: "agama_wali", label: "Agama" },
      { key: "warga_wali", label: "Kewarganegaraan" },
      { key: "pend_wali", label: "Pendidikan Terakhir" },
      { key: "pek_wali", label: "Pekerjaan" },
      { key: "gaji_wali", label: "Penghasilan" },
      { key: "gaji_waliktg", label: "Kategori Penghasilan" },
      { key: "alamat_wali", label: "Alamat" },
      { key: "rt_wali", label: "RT" },
      { key: "rw_wali", label: "RW" },
      { key: "kopos_wali", label: "Kode Pos" },
      { key: "kel_wali", label: "Kelurahan/Desa" },
      { key: "kec_wali", label: "Kecamatan" },
      { key: "telp_wali", label: "No. Telp/HP" },
    ],
  },
  {
    title: "Hobi & Minat",
    description: "Minat dan bakat di luar akademik.",
    icon: HeartIcon,
    fields: [
      { key: "hobi_seni", label: "Seni" },
      { key: "hobi_orkes", label: "Olahraga" },
      { key: "hobi_organisasi", label: "Organisasi" },
      { key: "hobi_lain", label: "Lainnya" },
    ],
  },
];

function formatValue(field, value) {
  if (!value) return "-";
  if (field.type === "date") return formatDateNumeric(value);
  return value;
}

onMounted(async () => {
  try {
    detail.value = await pdStore.fetchPesertaDidikMe();
  } catch {
    // detail tetap null, ditangani blok "Not found" di template
  }
});
</script>
