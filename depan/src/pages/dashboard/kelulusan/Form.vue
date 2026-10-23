<template>
  <DashboardPageHeader />
  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <!-- Info Siswa (read-only, identifikasi saat edit) -->
    <section v-if="isEdit && studentInfo" :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Informasi Siswa"
        description="Data diambil dari Peserta Didik, hanya untuk identifikasi."
      >
        <UserIcon class="h-5 w-5" />
      </SectionHeading>
      <div class="grid grid-cols-1 gap-4 p-6 sm:grid-cols-3">
        <div>
          <p class="text-xs text-gray-500 dark:text-slate-400">Nama Siswa</p>
          <p class="text-sm font-medium text-gray-800 dark:text-slate-200">
            {{ studentInfo.nama_siswa || "-" }}
          </p>
        </div>
        <div>
          <p class="text-xs text-gray-500 dark:text-slate-400">
            Kelas / Rombel
          </p>
          <p class="text-sm font-medium text-gray-800 dark:text-slate-200">
            {{ studentInfo.nama_kelas || "-" }} /
            {{ studentInfo.nama_rombel || "-" }}
          </p>
        </div>
        <div>
          <p class="text-xs text-gray-500 dark:text-slate-400">
            Tahun Pelajaran
          </p>
          <p class="text-sm font-medium text-gray-800 dark:text-slate-200">
            {{ studentInfo.nama_tahunpelajaran || "-" }}
          </p>
        </div>
      </div>
    </section>

    <!-- Section dinamis -->
    <section
      v-for="(section, index) in sections"
      :key="section.title"
      :class="cardClass"
      class="relative"
    >
      <div
        v-if="index === 0 && !(isEdit && studentInfo)"
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading :title="section.title" :description="section.description">
        <component :is="section.icon" class="h-5 w-5" />
      </SectionHeading>
      <div class="grid grid-cols-1 gap-6 p-6 sm:grid-cols-2">
        <FormField
          v-for="field in section.fields"
          :key="field.key"
          :label="field.label + (field.required ? ' *' : '')"
          :class="field.full ? 'sm:col-span-2' : ''"
        >
          <select
            v-if="field.type === 'select'"
            v-model="form[field.key]"
            :class="inputClass"
            :disabled="field.disabled"
            @change="field.onChange && field.onChange()"
          >
            <option value="">--Pilih--</option>
            <option
              v-for="opt in field.options"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.label }}
            </option>
          </select>
          <textarea
            v-else-if="field.type === 'textarea'"
            v-model="form[field.key]"
            :class="inputClass"
            rows="3"
            placeholder="Ketik disini..."
          />
          <input
            v-else
            type="text"
            v-model="form[field.key]"
            :class="inputClass"
            :disabled="field.key === 'nis' && isEdit"
            placeholder="Ketik disini..."
          />
          <p
            v-if="fieldErrors[field.key]"
            class="mt-1 text-xs text-red-600 dark:text-red-400"
          >
            {{ fieldErrors[field.key] }}
          </p>
        </FormField>
      </div>
    </section>

    <!-- Foto Kelulusan (dipilih dari galeri, bukan upload di sini — kelola foto di halaman Galeri Foto) -->
    <section :class="cardClass">
      <SectionHeading
        title="Foto Kelulusan (untuk SKL)"
        description="Pilih foto dari Galeri Foto Kelulusan. Kelola/unggah foto baru lewat halaman Galeri Foto."
      >
        <PhotoIcon class="h-5 w-5" />
      </SectionHeading>
      <div class="flex flex-col gap-4 p-6 sm:flex-row sm:items-start">
        <button
          type="button"
          :title="fotoPreviewUrl ? 'Klik untuk perbesar' : ''"
          :class="[
            'group h-32 w-24 shrink-0 overflow-hidden rounded-lg border border-gray-200 bg-gray-50 dark:border-slate-700 dark:bg-slate-800',
            fotoPreviewUrl ? 'cursor-pointer' : 'cursor-default',
          ]"
          @click="fotoPreviewUrl && openLightbox()"
        >
          <img
            v-if="fotoPreviewUrl"
            :src="fotoPreviewUrl"
            alt="Foto kelulusan terpilih"
            class="h-full w-full object-cover transition duration-200 group-hover:opacity-80"
          />
          <div
            v-else
            class="flex h-full w-full items-center justify-center text-gray-300 dark:text-slate-600"
          >
            <UserIcon class="h-10 w-10" />
          </div>
        </button>
        <div class="w-full max-w-sm">
          <FormField label="Pilih Foto">
            <select v-model="form.fotolulus" :class="inputClass">
              <option value="">--Belum dipilih--</option>
              <option
                v-for="foto in kelulusanStore.fotoList"
                :key="foto.id"
                :value="foto.gambar"
              >
                {{ foto.gambar }}
              </option>
            </select>
          </FormField>
          <p
            v-if="kelulusanStore.fotoList.length === 0"
            class="mt-2 text-xs text-amber-600 dark:text-amber-400"
          >
            Galeri kosong.
            <RouterLink
              :to="{ name: 'dashboard-kelulusan-galeri-foto' }"
              class="underline hover:text-amber-700 dark:hover:text-amber-300"
            >
              Unggah foto
            </RouterLink>
          </p>
        </div>
      </div>
    </section>

    <!-- Aksi -->
    <div
      class="flex flex-col gap-4 rounded-2xl bg-white px-5 py-4 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5 md:flex-row md:items-center md:justify-between"
    >
      <div class="flex gap-3">
        <button
          type="submit"
          :disabled="loading"
          class="inline-flex items-center justify-center gap-2 rounded-xl bg-blue-600 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:bg-gray-400 dark:disabled:bg-slate-700"
        >
          <ArrowPathIcon v-if="loading" class="h-5 w-5 animate-spin" />
          {{ loading ? "Menyimpan..." : "Simpan" }}
        </button>
        <RouterLink
          :to="{ name: 'dashboard-kelulusan' }"
          class="inline-flex items-center justify-center rounded-xl border border-gray-300 px-5 py-2.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
        >
          Batal
        </RouterLink>
      </div>
    </div>
  </form>

  <VueEasyLightbox
    :visible="visibleLightbox"
    :imgs="lightboxImages"
    :index="0"
    @hide="visibleLightbox = false"
  />
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import VueEasyLightbox from "vue-easy-lightbox";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";
import FormField from "@/components/form/FormField.vue";
import { useToast } from "@/composables/useToast";
import { useKelulusanStore } from "@/stores/kelulusan";
import { useTahunPelajaranStore } from "@/stores/tahunpelajaran";
import { useKelasStore } from "@/stores/kelas";
import { useRombelStore } from "@/stores/rombel";
import {
  ArrowPathIcon,
  AcademicCapIcon,
  ClipboardDocumentCheckIcon,
  PhotoIcon,
  UserIcon,
} from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const toast = useToast();
const kelulusanStore = useKelulusanStore();
const tahunPelajaranStore = useTahunPelajaranStore();
const kelasStore = useKelasStore();
const rombelStore = useRombelStore();

const isEdit = !!route.params.nis;
const loading = ref(false);
const studentInfo = ref(null);

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";
const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500 disabled:bg-gray-50 dark:disabled:bg-white/5 disabled:text-gray-400";

// ===== Foto kelulusan (dipilih dari galeri) =====
const FOTO_BASE_URL = "/storage/fotokelulusan/";
const fotoPreviewUrl = computed(() =>
  form.fotolulus ? FOTO_BASE_URL + form.fotolulus : "",
);

// ===== Lightbox foto (pola sama spt GaleriFoto.vue) =====
const visibleLightbox = ref(false);
const lightboxImages = ref([]);
function openLightbox() {
  if (!fotoPreviewUrl.value) return;
  lightboxImages.value = [fotoPreviewUrl.value];
  visibleLightbox.value = true;
}

// ===== Form state =====
const form = reactive({
  nis: route.params.nis || "",
  id_tahunpelajaran: "",
  id_kelas: "",
  id_rombel: "",
  nopes: "",
  pa_s: "",
  ppkn_s: "",
  bind_s: "",
  bing_s: "",
  mtk_s: "",
  ipa_s: "",
  ips_s: "",
  penjas_s: "",
  senbud_s: "",
  prakarya_s: "",
  jumlah: "",
  rata_s: "",
  ket: "",
  fotolulus: "",
});
const fieldErrors = reactive({
  nis: "",
});

// ===== Opsi dropdown (FK referensi, hanya yang aktif) =====
const tahunPelajaranOptions = computed(() =>
  tahunPelajaranStore.list
    .filter((t) => t.aktif === 1)
    .map((t) => ({ value: t.id_tahunpelajaran, label: t.nama_tahunpelajaran })),
);
const kelasOptions = computed(() =>
  kelasStore.list
    .filter((k) => k.aktif === 1)
    .map((k) => ({ value: k.id_kelas, label: k.nama_kelas })),
);
const rombelOptions = computed(() =>
  rombelStore.list
    .filter((r) => r.aktif === 1)
    .map((r) => ({ value: r.id_rombel, label: r.nama_rombel })),
);
const ketOptions = [
  { value: "Lulus", label: "Lulus" },
  { value: "Tidak Lulus", label: "Tidak Lulus" },
];

// ===== Cascading: Rombel difilter berdasarkan Kelas terpilih =====
async function onKelasChange() {
  form.id_rombel = "";
  if (form.id_kelas) {
    rombelStore.changeFilterKelas(form.id_kelas);
    await rombelStore.fetchAll();
  } else {
    rombelStore.changeFilterKelas(null);
    await rombelStore.fetchAll();
  }
}

// ===== Konfigurasi field per section (data-driven) =====
const sections = computed(() => [
  {
    title: "Data Kelulusan",
    description: "NIS siswa, kelas/rombel, no peserta, dan status kelulusan.",
    icon: AcademicCapIcon,
    fields: [
      { key: "nis", label: "NIS", required: true },
      {
        key: "id_tahunpelajaran",
        label: "Tahun Pelajaran",
        type: "select",
        options: tahunPelajaranOptions.value,
        required: true,
      },
      {
        key: "id_kelas",
        label: "Kelas",
        type: "select",
        options: kelasOptions.value,
        required: true,
        onChange: onKelasChange,
      },
      {
        key: "id_rombel",
        label: "Rombel",
        type: "select",
        options: rombelOptions.value,
        required: true,
        disabled: !form.id_kelas,
      },
      { key: "nopes", label: "No Peserta Ujian", required: true },
      {
        key: "ket",
        label: "Keterangan",
        type: "select",
        options: ketOptions,
        required: true,
      },
    ],
  },
  {
    title: "Nilai Ujian Sekolah",
    description: "Nilai per mata pelajaran serta jumlah dan rerata.",
    icon: ClipboardDocumentCheckIcon,
    fields: [
      { key: "pa_s", label: "Pendidikan Agama & Budi Pekerti", required: true },
      { key: "ppkn_s", label: "PPKn", required: true },
      { key: "bind_s", label: "Bahasa Indonesia", required: true },
      { key: "bing_s", label: "Bahasa Inggris", required: true },
      { key: "mtk_s", label: "Matematika", required: true },
      { key: "ipa_s", label: "IPA", required: true },
      { key: "ips_s", label: "IPS", required: true },
      { key: "penjas_s", label: "PJOK", required: true },
      { key: "senbud_s", label: "Seni Budaya", required: true },
      { key: "prakarya_s", label: "Informatika", required: true },
      { key: "jumlah", label: "Jumlah Nilai", required: true },
      { key: "rata_s", label: "Rerata Nilai", required: true },
    ],
  },
]);

// ===== Load data referensi =====
async function loadReferensi() {
  await Promise.all([
    tahunPelajaranStore.fetchAll(),
    kelasStore.fetchAll(),
    kelulusanStore.fetchFoto(),
  ]);
}

// ===== Load data (mode edit) =====
async function loadData() {
  if (!isEdit) return;
  loading.value = true;
  try {
    const data = await kelulusanStore.fetchKelulusanById(route.params.nis);
    studentInfo.value = data;
    form.id_tahunpelajaran = data.id_tahunpelajaran;
    form.id_kelas = data.id_kelas;
    form.nopes = data.nopes;
    form.pa_s = data.pa_s;
    form.ppkn_s = data.ppkn_s;
    form.bind_s = data.bind_s;
    form.bing_s = data.bing_s;
    form.mtk_s = data.mtk_s;
    form.ipa_s = data.ipa_s;
    form.ips_s = data.ips_s;
    form.penjas_s = data.penjas_s;
    form.senbud_s = data.senbud_s;
    form.prakarya_s = data.prakarya_s;
    form.jumlah = data.jumlah;
    form.rata_s = data.rata_s;
    form.ket = data.ket;
    form.fotolulus = data.fotolulus;
    // muat rombel sesuai kelas siswa dulu, baru set id_rombel-nya
    if (data.id_kelas) {
      rombelStore.changeFilterKelas(data.id_kelas);
      await rombelStore.fetchAll();
    }
    form.id_rombel = data.id_rombel;
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat data kelulusan.",
    );
  } finally {
    loading.value = false;
  }
}

// ===== Submit =====
function buildPayload() {
  return {
    id_tahunpelajaran: Number(form.id_tahunpelajaran),
    id_kelas: Number(form.id_kelas),
    id_rombel: Number(form.id_rombel),
    nis: form.nis,
    nopes: form.nopes,
    pa_s: form.pa_s,
    ppkn_s: form.ppkn_s,
    bind_s: form.bind_s,
    bing_s: form.bing_s,
    mtk_s: form.mtk_s,
    ipa_s: form.ipa_s,
    ips_s: form.ips_s,
    penjas_s: form.penjas_s,
    senbud_s: form.senbud_s,
    prakarya_s: form.prakarya_s,
    jumlah: form.jumlah,
    rata_s: form.rata_s,
    ket: form.ket,
    fotolulus: form.fotolulus,
  };
}

const requiredTextFields = [
  ["nopes", "No peserta ujian"],
  ["pa_s", "Nilai Pendidikan Agama"],
  ["ppkn_s", "Nilai PPKn"],
  ["bind_s", "Nilai Bahasa Indonesia"],
  ["bing_s", "Nilai Bahasa Inggris"],
  ["mtk_s", "Nilai Matematika"],
  ["ipa_s", "Nilai IPA"],
  ["ips_s", "Nilai IPS"],
  ["penjas_s", "Nilai PJOK"],
  ["senbud_s", "Nilai Seni Budaya"],
  ["prakarya_s", "Nilai Informatika"],
  ["jumlah", "Jumlah nilai"],
  ["rata_s", "Rerata nilai"],
];

async function submit() {
  if (loading.value) return;
  fieldErrors.nis = "";

  if (!isEdit && !form.nis.trim()) {
    toast.error("NIS wajib diisi.");
    return;
  }
  if (!form.id_tahunpelajaran) {
    toast.error("Tahun pelajaran wajib dipilih.");
    return;
  }
  if (!form.id_kelas) {
    toast.error("Kelas wajib dipilih.");
    return;
  }
  if (!form.id_rombel) {
    toast.error("Rombel wajib dipilih.");
    return;
  }
  if (!form.ket) {
    toast.error("Keterangan wajib dipilih.");
    return;
  }
  for (const [key, label] of requiredTextFields) {
    if (!String(form[key]).trim()) {
      toast.error(`${label} wajib diisi.`);
      return;
    }
  }

  loading.value = true;
  try {
    const payload = buildPayload();
    if (isEdit) {
      const response = await kelulusanStore.updateKelulusan(
        route.params.nis,
        payload,
      );
      const dataChanged = response.changed && response.changed.length > 0;
      if (dataChanged) {
        toast.success("Data kelulusan berhasil diperbarui.");
      } else {
        toast.warning("Tidak ada perubahan data.");
      }
    } else {
      await kelulusanStore.createKelulusan(payload);
      toast.success("Data kelulusan berhasil ditambahkan.");
    }
    router.push({ name: "dashboard-kelulusan" });
  } catch (error) {
    const message =
      error.response?.data?.message || "Gagal menyimpan data kelulusan.";
    if (message.toLowerCase().includes("nis")) {
      fieldErrors.nis = message;
    }
    toast.error(message);
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  await loadReferensi();
  await loadData();
});
</script>
