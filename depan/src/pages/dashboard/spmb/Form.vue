<template>
  <DashboardPageHeader />
  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <!-- Section dinamis -->
    <section
      v-for="(section, index) in sections"
      :key="section.title"
      :class="cardClass"
      class="relative"
    >
      <div
        v-if="index === 0"
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
          <input
            v-else
            type="text"
            v-model="form[field.key]"
            :class="inputClass"
            :disabled="field.key === 'no_peserta' && isEdit"
            placeholder="Ketik disini..."
          />
          <p
            v-if="fieldErrors[field.key]"
            class="mt-1 text-xs text-red-600 dark:text-red-400"
          >
            {{ fieldErrors[field.key] }}
          </p>
          <p
            v-else-if="duplicateWarnings[field.key]"
            class="mt-1 text-xs text-yellow-600 dark:text-yellow-400"
          >
            Data ganda a.n. {{ duplicateWarnings[field.key].nama_peserta }}
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
          :to="{ name: 'dashboard-spmb' }"
          class="inline-flex items-center justify-center rounded-xl border border-gray-300 px-5 py-2.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
        >
          Batal
        </RouterLink>
      </div>
    </div>
  </form>
</template>

<script setup>
import { onMounted, reactive, ref, watch } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";
import FormField from "@/components/form/FormField.vue";
import { useToast } from "@/composables/useToast";
import { useSpmbStore } from "@/stores/spmb";
import {
  ArrowPathIcon,
  IdentificationIcon,
  AcademicCapIcon,
} from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const toast = useToast();
const spmbStore = useSpmbStore();

const isEdit = !!route.params.id;
const loading = ref(false);

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";
const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500 disabled:bg-gray-50 dark:disabled:bg-white/5 disabled:text-gray-400";

const OPT_KELAMIN = [
  { value: "L", label: "Laki-laki" },
  { value: "P", label: "Perempuan" },
];

const OPT_JALUR = [
  { value: 1, label: "Prestasi Akademik" },
  { value: 2, label: "Prestasi Non Akademik" },
  { value: 3, label: "Penyandang Disabilitas" },
  { value: 4, label: "KJP Plus" },
  { value: 5, label: "Domisili" },
  { value: 6, label: "Mutasi" },
  { value: 7, label: "Tahap Kedua" },
];

const OPT_LAPOR_DIRI = [
  { value: 0, label: "Belum Lapor" },
  { value: 1, label: "Sudah Lapor" },
  { value: 2, label: "Undur Diri" },
];

// ===== Konfigurasi field per section (data-driven) =====
const sections = [
  {
    title: "Data Peserta",
    description: "Identitas peserta didik yang sudah diterima di sekolah.",
    icon: IdentificationIcon,
    fields: [
      { key: "no_peserta", label: "No. Peserta", required: true },
      {
        key: "nama_peserta",
        label: "Nama Peserta",
        required: true,
        full: true,
      },
      {
        key: "kelamin",
        label: "Jenis Kelamin",
        type: "select",
        options: OPT_KELAMIN,
        required: true,
      },
      {
        key: "asal_sekolah",
        label: "Asal Sekolah",
        required: true,
        full: true,
      },
    ],
  },
  {
    title: "Data SPMB",
    description: "Jalur penerimaan dan status lapor diri peserta.",
    icon: AcademicCapIcon,
    fields: [
      {
        key: "jalur",
        label: "Jalur",
        type: "select",
        options: OPT_JALUR,
        required: true,
      },
      {
        key: "lapor_diri",
        label: "Status Lapor Diri",
        type: "select",
        options: OPT_LAPOR_DIRI,
        required: true,
      },
    ],
  },
];

// ===== Form state =====
const form = reactive({
  no_peserta: "",
  nama_peserta: "",
  kelamin: "",
  asal_sekolah: "",
  jalur: "",
  lapor_diri: "",
});

// ===== Duplikat & error (khusus no_peserta) =====
const duplicateWarnings = reactive({
  no_peserta: null,
});
const fieldErrors = reactive({
  no_peserta: "",
});

const debounceTimers = {};
async function checkDuplicateNoPeserta(value) {
  try {
    const exclude = isEdit ? route.params.id : "";
    const result = await spmbStore.checkDuplicate(value, exclude);
    if (result?.exists) duplicateWarnings.no_peserta = result;
  } catch {
    // non-blocking, diamkan kalau gagal
  }
}
watch(
  () => form.no_peserta,
  (value) => {
    if (isEdit) return;
    clearTimeout(debounceTimers.no_peserta);
    duplicateWarnings.no_peserta = null;
    fieldErrors.no_peserta = "";
    if (!value) return;
    debounceTimers.no_peserta = setTimeout(
      () => checkDuplicateNoPeserta(value),
      300,
    );
  },
);

// ===== Load data (mode edit) =====
async function loadData() {
  if (!isEdit) return;
  loading.value = true;
  try {
    const data = await spmbStore.fetchSpmbById(route.params.id);
    Object.assign(form, data);
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat data peserta SPMB.",
    );
  } finally {
    loading.value = false;
  }
}

// ===== Submit =====
function buildPayload() {
  return {
    no_peserta: form.no_peserta,
    nama_peserta: form.nama_peserta,
    kelamin: form.kelamin,
    asal_sekolah: form.asal_sekolah,
    jalur: Number(form.jalur),
    lapor_diri: Number(form.lapor_diri),
  };
}

async function submit() {
  if (loading.value) return;
  if (!isEdit && !form.no_peserta.trim()) {
    toast.error("No. peserta wajib diisi.");
    return;
  }
  if (!form.nama_peserta.trim()) {
    toast.error("Nama peserta wajib diisi.");
    return;
  }
  if (!form.kelamin) {
    toast.error("Jenis kelamin wajib dipilih.");
    return;
  }
  if (!form.asal_sekolah.trim()) {
    toast.error("Asal sekolah wajib diisi.");
    return;
  }
  if (!form.jalur) {
    toast.error("Jalur wajib dipilih.");
    return;
  }
  if (form.lapor_diri === "" || form.lapor_diri === null) {
    toast.error("Status lapor diri wajib dipilih.");
    return;
  }
  if (duplicateWarnings.no_peserta) {
    toast.error(
      "No. peserta sudah terdaftar pada data lain, silakan periksa kembali.",
    );
    return;
  }

  loading.value = true;
  try {
    const payload = buildPayload();
    if (isEdit) {
      const response = await spmbStore.updateSpmb(route.params.id, payload);
      const dataChanged = response.changed && response.changed.length > 0;
      if (dataChanged) {
        toast.success("Data peserta SPMB berhasil diperbarui.");
      } else {
        toast.warning("Tidak ada perubahan data.");
      }
    } else {
      await spmbStore.createSpmb(payload);
      toast.success("Peserta SPMB berhasil ditambahkan.");
    }
    router.push({ name: "dashboard-spmb" });
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menyimpan data peserta SPMB.",
    );
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  await loadData();
});
</script>
