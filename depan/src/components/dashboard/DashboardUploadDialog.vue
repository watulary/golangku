<template>
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0 scale-95"
    enter-to-class="opacity-100 scale-100"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100 scale-100"
    leave-to-class="opacity-0 scale-95"
  >
    <div
      v-if="modelValue"
      class="fixed inset-0 z-100 flex items-center justify-center p-4"
    >
      <!-- Backdrop -->
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm"></div>
      <!-- Dialog -->
      <div
        class="relative z-10 w-full max-w-lg overflow-hidden rounded-2xl bg-white shadow-2xl ring-1 ring-slate-900/5 dark:bg-slate-900 dark:ring-white/10"
      >
        <!-- Header -->
        <div
          class="flex items-center gap-3 border-b border-gray-100 px-6 py-5 dark:border-white/5"
        >
          <div
            class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-blue-100 text-blue-600 dark:bg-blue-500/10 dark:text-blue-300"
          >
            <ArrowUpTrayIcon class="h-5 w-5" />
          </div>
          <div>
            <h2 class="text-lg font-semibold text-gray-800 dark:text-slate-200">
              {{ title }}
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
              {{ subtitle }}
            </p>
          </div>
        </div>
        <!-- Body -->
        <div class="space-y-5 px-6 py-5">
          <!-- Template -->
          <div
            class="rounded-xl border border-gray-100 p-4 dark:border-white/5 dark:bg-[#0B0F1A]"
          >
            <p class="text-sm text-gray-700 dark:text-gray-300">
              Unduh Template (.xlsx) dibawah ini;
            </p>
            <button
              type="button"
              class="mt-2 inline-flex items-center gap-2 rounded-lg border border-emerald-300 bg-emerald-50 px-3 py-2 text-sm font-medium text-emerald-700 transition hover:bg-emerald-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 dark:border-emerald-500/20 dark:bg-emerald-500/10 dark:text-emerald-400 dark:hover:bg-emerald-500/20"
              @click="emit('download-template')"
            >
              <DocumentArrowDownIcon class="h-5 w-5" />
              <span>Template (.xlsx)</span>
            </button>
          </div>
          <!-- File -->
          <div>
            <input
              ref="fileInput"
              type="file"
              accept=".xlsx"
              class="block w-full rounded-xl border border-gray-300 bg-white px-3 py-2 text-sm text-gray-700 outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 file:mr-4 file:rounded-lg file:border-0 file:bg-blue-600 file:px-4 file:py-2 file:text-sm file:font-medium file:text-white hover:file:bg-blue-700 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200"
              @change="handleFileChange"
            />
          </div>
          <!-- Error -->
          <div
            v-if="errors.length"
            class="rounded-xl border border-red-200 bg-red-50 p-4 dark:border-red-500/20 dark:bg-red-500/10"
          >
            <h3
              class="mb-2 text-sm font-semibold text-red-700 dark:text-red-300"
            >
              Unggah dibatalkan karena terdapat kesalahan:
            </h3>
            <ul
              class="max-h-52 space-y-1 overflow-y-auto text-sm text-red-600 dark:text-red-300"
            >
              <li v-for="(item, index) in errors" :key="index">
                • Baris {{ item.row }} — {{ item.message }}
              </li>
            </ul>
          </div>
        </div>
        <!-- Footer -->
        <div
          class="flex items-center gap-3 border-t border-gray-100 bg-white px-6 py-4 dark:border-white/5 dark:bg-slate-900"
        >
          <button
            type="button"
            class="rounded-xl bg-blue-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="!selectedFile || loading"
            @click="upload"
          >
            <span v-if="loading">Proses...</span>
            <span v-else>Unggah</span>
          </button>
          <button
            type="button"
            class="rounded-xl border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:opacity-50 dark:border-slate-600 dark:text-slate-200 dark:hover:bg-slate-800"
            :disabled="loading"
            @click="closeDialog"
          >
            Batal
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
@media (prefers-reduced-motion: reduce) {
  .transition {
    transition: none !important;
  }
}
</style>

<script setup>
import { ref, watch } from "vue";
import {
  DocumentArrowDownIcon,
  ArrowUpTrayIcon,
} from "@heroicons/vue/24/outline";

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: "Unggah Data",
  },
  subtitle: {
    type: String,
    default: "Unggah data menggunakan Template (.xlsx)",
  },
  // Fungsi async(file) -> Promise<{ message, ... }>, dilempar dari komponen pemanggil.
  // Dialog ini tidak lagi tahu store mana yang dipakai (generik lintas modul).
  uploadAction: {
    type: Function,
    required: true,
  },
});
const emit = defineEmits(["update:modelValue", "success", "download-template"]);

const fileInput = ref(null);
const selectedFile = ref(null);
const loading = ref(false);
const errors = ref([]);

function handleFileChange(event) {
  const files = event.target.files;
  if (!files || files.length === 0) {
    selectedFile.value = null;
    return;
  }
  selectedFile.value = files[0];
  errors.value = [];
}

function closeDialog() {
  if (loading.value) return;
  emit("update:modelValue", false);
}

async function upload() {
  if (!selectedFile.value) return;
  loading.value = true;
  errors.value = [];
  try {
    const response = await props.uploadAction(selectedFile.value);
    emit("success", response.message);
    emit("update:modelValue", false);
    selectedFile.value = null;
    if (fileInput.value) {
      fileInput.value.value = "";
    }
  } catch (error) {
    if (error.response?.data?.errors) {
      errors.value = error.response.data.errors;
    } else {
      errors.value = [
        {
          row: "-",
          message:
            error.response?.data?.message || "Terjadi kesalahan saat upload.",
        },
      ];
    }
  } finally {
    loading.value = false;
  }
}

watch(
  () => props.modelValue,
  (value) => {
    if (value) {
      selectedFile.value = null;
      errors.value = [];
      if (fileInput.value) {
        fileInput.value.value = "";
      }
    }
  },
);
</script>
