<template>
  <DashboardPageHeader />
  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <!-- Foto -->
    <section v-if="isEdit && studentInfo" :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Foto Siswa"
        description="Foto profil diambil dari data peserta didik."
      >
        <PhotoIcon class="h-5 w-5" />
      </SectionHeading>
      <div class="flex items-center gap-4 p-6">
        <div class="group relative">
          <div
            class="h-32 w-24 overflow-hidden rounded-lg border border-gray-200 bg-gray-50 dark:border-slate-700 dark:bg-slate-800"
            :class="{ 'cursor-pointer': studentInfo.foto_url }"
            @click="studentInfo.foto_url && showImage()"
          >
            <img
              v-if="studentInfo.foto_url"
              :src="studentInfo.foto_url"
              :alt="studentInfo.nama_siswa"
              class="h-full w-full object-cover transition duration-200 group-hover:opacity-80"
            />
            <div
              v-else
              class="flex h-full w-full items-center justify-center text-gray-300 dark:text-slate-600"
            >
              <UserIcon class="h-12 w-12" />
            </div>
            <div
              v-if="studentInfo.foto_url"
              class="pointer-events-none absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition duration-200 group-hover:opacity-100"
            >
              <MagnifyingGlassPlusIcon class="h-5 w-5 text-white" />
            </div>
          </div>
          <button
            type="button"
            title="Ganti Foto"
            @click="fotoInput.click()"
            class="absolute -right-3 -bottom-3 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700 dark:hover:text-white"
          >
            <PencilSquareIcon class="h-4 w-4" />
          </button>
          <button
            v-if="!isDefaultFoto"
            type="button"
            title="Hapus Foto"
            @click="showDeleteFotoDialog = true"
            class="absolute -left-3 -bottom-3 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 hover:bg-red-50 hover:text-red-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-red-500/10 dark:hover:text-red-400"
          >
            <XMarkIcon class="h-4 w-4" />
          </button>
          <input
            ref="fotoInput"
            type="file"
            accept=".jpg,.jpeg,.png,.webp"
            class="hidden"
            @change="handleFotoChange"
          />
        </div>
        <div>
          <p class="text-sm font-medium text-gray-800 dark:text-slate-200">
            {{ studentInfo.nama_siswa || "-" }}
          </p>
          <p class="text-xs text-gray-500 dark:text-slate-400">
            NISN: {{ studentInfo.nisn || "-" }} &bull;
            {{ studentInfo.jenis_kelamin || "-" }}
          </p>
        </div>
      </div>
    </section>

    <!-- Section dinamis -->
    <section
      v-for="section in sections"
      :key="section.title"
      :class="cardClass"
      class="relative"
    >
      <div
        v-if="!(isEdit && studentInfo)"
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
          :to="{ name: 'dashboard-kenaikan' }"
          class="inline-flex items-center justify-center rounded-xl border border-gray-300 px-5 py-2.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
        >
          Batal
        </RouterLink>
      </div>
    </div>
  </form>
  <DashboardConfirmDialog
    v-model="showDeleteFotoDialog"
    title="Hapus Foto"
    message="Foto akan dikembalikan ke foto default. Lanjutkan?"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="deletingFoto"
    @confirm="handleDeleteFoto"
  />
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
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";
import FormField from "@/components/form/FormField.vue";
import { useToast } from "@/composables/useToast";
import { useKenaikanStore } from "@/stores/kenaikan";
import { usePesertaDidikStore } from "@/stores/pesertadidik";
import { useTahunPelajaranStore } from "@/stores/tahunpelajaran";
import { useKelasStore } from "@/stores/kelas";
import { useRombelStore } from "@/stores/rombel";
import VueEasyLightbox from "vue-easy-lightbox";
import {
  ArrowPathIcon,
  AcademicCapIcon,
  PhotoIcon,
  PencilSquareIcon,
  UserIcon,
  XMarkIcon,
  MagnifyingGlassPlusIcon,
} from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const toast = useToast();
const kenaikanStore = useKenaikanStore();
const pdStore = usePesertaDidikStore();
const tahunPelajaranStore = useTahunPelajaranStore();
const kelasStore = useKelasStore();
const rombelStore = useRombelStore();

const isEdit = !!route.params.nis;
const loading = ref(false);
const studentInfo = ref(null);

const FOTO_DEFAULT = "1636812404_b39deec2f158cb789567.jpg";
const isDefaultFoto = computed(
  () =>
    !studentInfo.value?.foto_url ||
    studentInfo.value.foto_url.endsWith(FOTO_DEFAULT),
);
const fotoInput = ref(null);
const showDeleteFotoDialog = ref(false);
const deletingFoto = ref(false);
const visibleLightbox = ref(false);
const lightboxImages = computed(() =>
  studentInfo.value?.foto_url ? [studentInfo.value.foto_url] : [],
);

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";
const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500 disabled:bg-gray-50 dark:disabled:bg-white/5 disabled:text-gray-400";

// ===== Form state =====
const form = reactive({
  nis: route.params.nis || "",
  id_tahunpelajaran: "",
  id_kelas: "",
  id_rombel: "",
  ket: "",
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
    title: "Data Kenaikan Kelas",
    description: "NIS siswa serta kelas/rombel tujuan hasil kenaikan.",
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
      {
        key: "ket",
        label: "Keterangan",
        type: "textarea",
        required: true,
        full: true,
      },
    ],
  },
]);

// ===== Load data referensi =====
async function loadReferensi() {
  await Promise.all([tahunPelajaranStore.fetchAll(), kelasStore.fetchAll()]);
}

// ===== Load data (mode edit) =====
async function loadData() {
  if (!isEdit) return;
  loading.value = true;
  try {
    const data = await kenaikanStore.fetchKenaikanById(route.params.nis);
    studentInfo.value = data;
    form.id_tahunpelajaran = data.id_tahunpelajaran;
    form.id_kelas = data.id_kelas;
    form.ket = data.ket;
    // muat rombel sesuai kelas siswa dulu, baru set id_rombel-nya
    if (data.id_kelas) {
      rombelStore.changeFilterKelas(data.id_kelas);
      await rombelStore.fetchAll();
    }
    form.id_rombel = data.id_rombel;
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal memuat data kenaikan.");
  } finally {
    loading.value = false;
  }
}

// ===== Foto (langsung tersimpan lewat pdStore, terpisah dari submit kenaikan) =====
function showImage() {
  visibleLightbox.value = true;
}
async function handleFotoChange(event) {
  const file = event.target.files?.[0];
  if (!file) return;
  const maxSize = 2 * 1024 * 1024;
  if (file.size > maxSize) {
    toast.error("Ukuran foto maksimal 2 MB.");
    event.target.value = "";
    return;
  }
  const allowed = ["image/jpeg", "image/png", "image/webp"];
  if (!allowed.includes(file.type)) {
    toast.error("Foto harus berformat JPG, JPEG, PNG atau WEBP.");
    event.target.value = "";
    return;
  }
  try {
    const fd = new FormData();
    fd.append("foto", file);
    const res = await pdStore.updateFotoPesertaDidik(route.params.nis, fd);
    studentInfo.value.foto_url = res.foto_url;
    toast.success("Foto siswa berhasil diperbarui.");
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal mengunggah foto.");
  } finally {
    event.target.value = "";
  }
}
async function handleDeleteFoto() {
  deletingFoto.value = true;
  try {
    const res = await pdStore.deleteFotoPesertaDidik(route.params.nis);
    studentInfo.value.foto_url = res.foto_url;
    toast.success("Foto siswa berhasil dihapus.");
    showDeleteFotoDialog.value = false;
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menghapus foto.");
  } finally {
    deletingFoto.value = false;
  }
}

// ===== Submit =====
function buildPayload() {
  return {
    id_tahunpelajaran: Number(form.id_tahunpelajaran),
    id_kelas: Number(form.id_kelas),
    id_rombel: Number(form.id_rombel),
    nis: form.nis,
    ket: form.ket,
  };
}
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
  if (!form.ket.trim()) {
    toast.error("Keterangan wajib diisi.");
    return;
  }

  loading.value = true;
  try {
    const payload = buildPayload();
    if (isEdit) {
      const response = await kenaikanStore.updateKenaikan(
        route.params.nis,
        payload,
      );
      const dataChanged = response.changed && response.changed.length > 0;
      if (dataChanged) {
        toast.success("Data kenaikan berhasil diperbarui.");
      } else {
        toast.warning("Tidak ada perubahan data.");
      }
    } else {
      await kenaikanStore.createKenaikan(payload);
      toast.success("Data kenaikan berhasil ditambahkan.");
    }
    router.push({ name: "dashboard-kenaikan" });
  } catch (error) {
    const message =
      error.response?.data?.message || "Gagal menyimpan data kenaikan.";
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
