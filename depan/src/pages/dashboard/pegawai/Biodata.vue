<template>
  <DashboardPageHeader
    title="Biodata Diri"
    subtitle="Data lengkap kepegawaian yang tercatat di sekolah."
  />

  <!-- Loading -->
  <div
    v-if="pegawaiStore.meLoading && !detail"
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
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>
      <div
        class="group relative h-32 w-24 flex-shrink-0 overflow-hidden rounded-lg border border-gray-200 bg-gray-50 dark:border-slate-700 dark:bg-slate-800"
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
          {{ detail.nama_pegawai }}
        </p>
        <p class="text-sm text-gray-500 dark:text-slate-400">
          NUPTK: {{ detail.nuptk }} · {{ detail.jabatan }}
        </p>
        <p class="text-sm text-gray-500 dark:text-slate-400">
          {{ detail.pns_non || "-" }}
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
          class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/[0.08] text-[#041983] dark:bg-white/5 dark:text-indigo-300"
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
import { usePegawaiStore } from "@/stores/pegawai";
import { formatDateNumeric } from "@/utils/date";
import {
  BriefcaseIcon,
  IdentificationIcon,
  MapPinIcon,
  DocumentTextIcon,
  UserGroupIcon,
  UserIcon,
  MagnifyingGlassPlusIcon,
} from "@heroicons/vue/24/outline";

const pegawaiStore = usePegawaiStore();
const detail = ref(null);
const visibleLightbox = ref(false);
const images = computed(() =>
  detail.value?.foto_url ? [detail.value.foto_url] : [],
);

const sections = [
  {
    title: "Data Kepegawaian",
    description: "Identitas jabatan dan status kepegawaian.",
    icon: BriefcaseIcon,
    fields: [
      { key: "jabatan", label: "Jabatan" },
      { key: "strata", label: "Strata Pendidikan" },
      { key: "pns_non", label: "Status Kepegawaian" },
      { key: "nikki", label: "NIKKI" },
      { key: "simpkbid", label: "SIMPKB ID" },
      { key: "dapodikid", label: "Dapodik ID" },
      { key: "nip_lama", label: "NIP Lama" },
      { key: "nip_baru", label: "NIP Baru" },
      { key: "nrk", label: "NRK" },
      { key: "nrg", label: "NRG" },
      { key: "karpeg", label: "Karpeg" },
      { key: "karis", label: "Karis/Karsu" },
    ],
  },
  {
    title: "Data Pribadi",
    description: "Identitas dan data pribadi.",
    icon: IdentificationIcon,
    fields: [
      { key: "jenis_kelamin", label: "Jenis Kelamin" },
      { key: "tempat_lahir", label: "Tempat Lahir" },
      { key: "tanggal_lahir", label: "Tanggal Lahir", type: "date" },
      { key: "agama", label: "Agama" },
      { key: "nik", label: "NIK" },
      { key: "no_kk", label: "No. KK" },
      { key: "npwp", label: "NPWP" },
      { key: "nama_ibukdg", label: "Nama Ibu Kandung" },
      { key: "nm_suami_istri", label: "Nama Suami/Istri" },
      { key: "pek_suami_istri", label: "Pekerjaan Suami/Istri" },
    ],
  },
  {
    title: "Alamat & Kontak",
    description: "Alamat tempat tinggal dan kontak.",
    icon: MapPinIcon,
    fields: [
      { key: "alamat", label: "Alamat" },
      { key: "rt", label: "RT" },
      { key: "rw", label: "RW" },
      { key: "kel", label: "Kelurahan/Desa" },
      { key: "kec", label: "Kecamatan" },
      { key: "kopos", label: "Kode Pos" },
      { key: "telp_rmh", label: "Telp Rumah" },
      { key: "telp_hp", label: "No. HP" },
      { key: "email", label: "Email" },
    ],
  },
  {
    title: "Riwayat SK",
    description: "Riwayat surat keputusan dan TMT kepegawaian.",
    icon: DocumentTextIcon,
    fields: [
      { key: "sk_sekolah", label: "SK Sekolah" },
      { key: "tmt_sekolah", label: "TMT Sekolah" },
      { key: "sk_nonkki", label: "SK Non KKI" },
      { key: "tmt_nonkki", label: "TMT Non KKI" },
      { key: "sk_kki", label: "SK KKI" },
      { key: "tmt_kki", label: "TMT KKI" },
      { key: "sk_cpns", label: "SK CPNS" },
      { key: "tmt_cpns", label: "TMT CPNS" },
      { key: "sk_pns", label: "SK PNS" },
      { key: "tmt_pns", label: "TMT PNS" },
      { key: "tmt_kgb", label: "TMT KGB" },
      { key: "pangkat_kgb", label: "Pangkat KGB" },
      { key: "gaji_kgb", label: "Gaji KGB" },
      { key: "sk_kgb", label: "SK KGB" },
      { key: "tglsk_kgb", label: "Tanggal SK KGB" },
    ],
  },
  {
    title: "Status Kepegawaian Saat Ini",
    description: "Status keluar/aktif dan keterangan terkait.",
    icon: UserGroupIcon,
    fields: [
      {
        key: "status_keluar",
        label: "Status",
        format: (v) => (v ? "Keluar" : "Aktif"),
      },
      { key: "keterangan_keluar", label: "Keterangan Keluar" },
      { key: "tanggal_keluar", label: "Tanggal Keluar", type: "date" },
      { key: "nomor_surat", label: "Nomor Surat" },
    ],
  },
];

function formatValue(field, value) {
  if (field.format) return field.format(value) || "-";
  if (!value) return "-";
  if (field.type === "date") return formatDateNumeric(value);
  return value;
}

onMounted(async () => {
  try {
    detail.value = await pegawaiStore.fetchMePegawai();
  } catch {
    // detail tetap null, ditangani blok "Not found" di template
  }
});
</script>
