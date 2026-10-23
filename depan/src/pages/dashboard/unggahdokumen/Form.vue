<template>
  <DashboardPageHeader
    :title="isEdit ? 'Edit Dokumen' : 'Tambah Dokumen'"
    subtitle="Kelola dokumen yang dapat diunduh oleh pengunjung."
  />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Dokumen"
        description="Data dan file dokumen yang dapat diunduh."
      >
        <DocumentTextIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="space-y-6 p-6">
        <!-- Dokumen -->
        <div
          class="flex min-h-40 items-center justify-center rounded-2xl border border-dashed border-gray-300 bg-gray-50 p-5 dark:border-slate-700 dark:bg-[#0B0F1A]"
        >
          <div class="group relative flex flex-col items-center gap-2">
            <a
              v-if="form.dokumen_url"
              :href="form.dokumen_url"
              rel="noopener"
              class="flex h-20 w-20 items-center justify-center rounded-lg bg-gray-100 text-gray-600 transition hover:bg-gray-200 hover:text-gray-800 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700 dark:hover:text-white"
            >
              <DocumentArrowDownIcon class="h-10 w-10" />
            </a>

            <DocumentTextIcon
              v-else
              class="h-15 w-15 text-gray-300 transition duration-200 group-hover:opacity-80 dark:text-slate-600"
            />

            <p
              v-if="selectedDokumen"
              class="max-w-52 truncate text-center text-xs text-gray-500 dark:text-slate-400"
            >
              {{ selectedDokumen.name }}
            </p>

            <button
              type="button"
              title="Ganti Dokumen"
              @click="dokumenInput.click()"
              class="absolute -right-4 -bottom-2 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700 dark:hover:text-white"
            >
              <PencilSquareIcon class="h-4 w-4" />
            </button>

            <input
              ref="dokumenInput"
              type="file"
              accept=".pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx"
              class="hidden"
              @change="handleDokumenChange"
            />
          </div>
        </div>

        <!-- Nama / Judul Dokumen -->
        <FormField label="Nama Dokumen">
          <input
            v-model="form.nama"
            type="text"
            :class="inputClass"
            placeholder="Ketik disini..."
          />
        </FormField>
      </div>

      <!-- Aksi -->
      <div
        class="flex items-center gap-3 border-t border-gray-100 px-6 py-4 dark:border-white/5"
      >
        <button
          type="submit"
          :disabled="loading"
          class="inline-flex items-center justify-center gap-2 rounded-xl bg-blue-600 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:bg-gray-400 dark:disabled:bg-slate-700"
        >
          <ArrowPathIcon v-if="loading" class="h-5 w-5 animate-spin" />
          {{ loading ? "Menyimpan..." : "Simpan" }}
        </button>

        <RouterLink
          :to="{ name: 'dashboard-unggahdokumen' }"
          class="inline-flex items-center justify-center rounded-xl border border-gray-300 px-5 py-2.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
        >
          Batal
        </RouterLink>
      </div>
    </section>
  </form>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";

import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import FormField from "@/components/form/FormField.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";

import { useToast } from "@/composables/useToast";
import { useUnggahdokumenStore } from "@/stores/unggahdokumen";

import {
  ArrowPathIcon,
  DocumentTextIcon,
  DocumentArrowDownIcon,
  PencilSquareIcon,
} from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const unggahdokumenStore = useUnggahdokumenStore();
const toast = useToast();

const isEdit = !!route.params.id;

const dokumenInput = ref(null);
const selectedDokumen = ref(null);
const loading = ref(false);

const form = reactive({
  id: null,
  nama: "",
  nama_pengguna: "",
  dokumen_url: "",
  waktu: "",
});

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";

const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";

const allowedExt = [".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"];

// =====================================
// Load Data (mode edit)
// =====================================

async function loadData() {
  if (!isEdit) return;

  loading.value = true;

  try {
    const data = await unggahdokumenStore.fetchUnggahdokumenById(
      route.params.id,
    );

    Object.assign(form, {
      id: data.id,
      nama: data.nama || "",
      nama_pengguna: data.nama_pengguna || "",
      dokumen_url: data.dokumen_url || "",
      waktu: data.waktu || "",
    });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal memuat data dokumen.");
  } finally {
    loading.value = false;
  }
}

// =====================================
// Submit
// =====================================

async function submit() {
  if (loading.value) return;

  if (!form.nama.trim()) {
    toast.error("Nama dokumen wajib diisi.");
    return;
  }

  if (!isEdit && !selectedDokumen.value) {
    toast.error("File dokumen wajib diunggah.");
    return;
  }

  loading.value = true;

  try {
    const data = new FormData();

    data.append("nama", form.nama);

    if (selectedDokumen.value) {
      data.append("dokumen", selectedDokumen.value);
    }

    if (isEdit) {
      const response = await unggahdokumenStore.updateUnggahdokumen(
        form.id,
        data,
      );

      if (response.changed) {
        toast.success(response.message || "Dokumen berhasil diperbarui.");
      } else {
        toast.warning(response.message || "Tidak ada perubahan data.");
      }
    } else {
      await unggahdokumenStore.createUnggahdokumen(data);
      toast.success("Dokumen berhasil ditambahkan.");
    }

    router.push({ name: "dashboard-unggahdokumen" });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menyimpan dokumen.");
  } finally {
    loading.value = false;
  }
}

function handleDokumenChange(event) {
  const file = event.target.files?.[0];

  if (!file) return;

  const maxSize = 2 * 1024 * 1024;
  const ext = "." + file.name.split(".").pop().toLowerCase();

  if (file.size > maxSize) {
    toast.error("Ukuran dokumen maksimal 2 MB.");
    event.target.value = "";
    return;
  }

  if (!allowedExt.includes(ext)) {
    toast.error("Format file tidak didukung.");
    event.target.value = "";
    return;
  }

  selectedDokumen.value = file;
  form.dokumen_url = "";
}

onMounted(async () => {
  await loadData();
});
</script>
