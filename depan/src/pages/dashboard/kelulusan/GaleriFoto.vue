<template>
  <DashboardPageHeader />

  <DashboardBulkAction
    :selected-count="selectedIds.length"
    @clear="selectedIds = []"
  >
    <template #actions>
      <button
        type="button"
        @click="showDeleteDialog = true"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-red-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-red-700"
      >
        <TrashIcon class="h-5 w-5" />
        Hapus
      </button>
    </template>
  </DashboardBulkAction>

  <div class="space-y-6">
    <div
      class="relative overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
    >
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <div class="border-b border-gray-100 px-6 py-5 dark:border-white/5">
        <div
          class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/[0.08] text-[#041983] dark:bg-white/5 dark:text-indigo-300"
            >
              <PhotoIcon class="h-5 w-5" />
            </div>
            <div>
              <h2
                class="text-xl font-semibold text-gray-700 dark:text-slate-200"
              >
                Galeri Foto Kelulusan
              </h2>
              <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
                Unggah file ZIP berisi foto (jpg, jpeg, png) maksimal 10 MB.
              </p>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <label
              v-if="selectableFoto.length > 0"
              class="flex cursor-pointer items-center gap-2 text-sm text-gray-600 dark:text-slate-300"
            >
              <input
                type="checkbox"
                :checked="isAllSelected"
                @change="toggleSelectAll"
                class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              />
              Pilih Semua
            </label>
            <RouterLink
              :to="{ name: 'dashboard-kelulusan' }"
              title="Kembali ke Data Kelulusan"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-slate-400 hover:bg-slate-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-400/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-slate-500 dark:hover:bg-slate-800"
            >
              <ArrowLeftIcon class="h-5 w-5" />
            </RouterLink>
            <button
              type="button"
              title="Unggah ZIP Foto"
              @click="showUploadDialog = true"
              class="inline-flex items-center justify-center gap-2 rounded-xl bg-blue-600 px-4 py-2.5 text-sm font-medium text-white shadow-sm transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40"
            >
              <ArrowUpTrayIcon class="h-5 w-5" />
              Unggah
            </button>
          </div>
        </div>
      </div>

      <div v-if="kelulusanStore.fotoLoading" class="p-6">
        <div class="flex flex-wrap gap-3">
          <div
            v-for="n in 10"
            :key="`skeleton-${n}`"
            class="h-[3cm] w-[2cm] shrink-0 animate-pulse rounded-lg bg-gray-100 dark:bg-white/5"
          ></div>
        </div>
      </div>

      <div
        v-else-if="kelulusanStore.fotoList.length === 0"
        class="py-16 text-center"
      >
        <PhotoIcon
          class="mx-auto h-12 w-12 text-gray-300 dark:text-slate-600"
        />
        <p class="mt-3 text-sm text-gray-500 dark:text-slate-400">
          Belum ada foto.
        </p>
      </div>

      <div v-else class="max-h-[70vh] overflow-auto p-6">
        <div class="flex flex-wrap gap-3">
          <div
            v-for="foto in kelulusanStore.fotoList"
            :key="foto.id"
            class="group relative w-[2cm] shrink-0 overflow-hidden rounded-lg border border-gray-200 bg-gray-50 dark:border-slate-700 dark:bg-slate-800"
          >
            <label
              v-if="foto.gambar !== DEFAULT_FOTO"
              class="absolute left-1 top-1 z-10"
            >
              <input
                type="checkbox"
                :checked="selectedIds.includes(foto.id)"
                @change="toggleSelect(foto.id)"
                class="h-3.5 w-3.5 rounded border-gray-300 text-blue-600 shadow focus:ring-blue-500"
              />
            </label>
            <span
              v-else
              title="Foto default, tidak bisa dihapus"
              class="absolute left-1 top-1 z-10 rounded-full bg-amber-100 px-1.5 py-0.5 text-[8px] font-medium text-amber-700 shadow dark:bg-amber-500/20 dark:text-amber-300"
            >
              Default
            </span>
            <button
              type="button"
              title="Klik untuk perbesar"
              class="block h-[3cm] w-[2cm] cursor-pointer"
              @click="openLightbox(foto)"
            >
              <img
                :src="fotoUrl(foto.gambar)"
                :alt="foto.gambar"
                class="h-full w-full object-cover transition duration-200 group-hover:opacity-80"
              />
            </button>
            <p
              class="truncate border-t border-gray-200 bg-white px-1 py-1 text-[9px] leading-tight text-gray-600 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-400"
              :title="foto.gambar"
            >
              {{ foto.gambar }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>

  <DashboardConfirmDialog
    v-model="showDeleteDialog"
    title="Hapus Foto"
    :message="`Anda akan menghapus ${selectedIds.length} foto dari galeri.\n\nFoto yang masih dipakai di data kelulusan (fotolulus) sebaiknya diganti dulu sebelum dihapus. Data yang telah dihapus tidak dapat dikembalikan.`"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="loadingAction"
    @confirm="handleDelete"
  />

  <!-- Dialog Upload ZIP (custom, bukan DashboardUploadDialog — itu hardcode utk .xlsx) -->
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0 scale-95"
    enter-to-class="opacity-100 scale-100"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100 scale-100"
    leave-to-class="opacity-0 scale-95"
  >
    <div
      v-if="showUploadDialog"
      class="fixed inset-0 z-100 flex items-center justify-center"
    >
      <div
        class="absolute inset-0 bg-black/50 backdrop-blur-sm"
        @click="!uploadLoading && (showUploadDialog = false)"
      ></div>
      <div
        class="relative z-10 w-full max-w-lg overflow-hidden rounded-2xl bg-white shadow-2xl ring-1 ring-slate-900/5 dark:bg-slate-900 dark:ring-white/10"
      >
        <div class="border-b border-gray-100 px-6 py-4 dark:border-white/5">
          <h2 class="text-lg font-semibold text-gray-800 dark:text-slate-200">
            Unggah Foto Kelulusan
          </h2>
          <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
            Unggah file ZIP berisi foto (jpg, jpeg, png), maksimal 10 MB.
          </p>
        </div>
        <div class="space-y-3 px-6 py-5">
          <input
            ref="zipInput"
            type="file"
            accept=".zip"
            class="block w-full rounded-xl border border-gray-300 bg-white px-3 py-2 text-sm text-gray-700 outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 file:mr-4 file:rounded-lg file:border-0 file:bg-blue-600 file:px-4 file:py-2 file:text-sm file:font-medium file:text-white hover:file:bg-blue-700 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200"
            @change="handleZipChange"
          />
          <p v-if="uploadError" class="text-sm text-red-600 dark:text-red-400">
            {{ uploadError }}
          </p>
        </div>
        <div
          class="flex items-center gap-3 border-t border-gray-100 bg-gray-50 px-6 py-4 dark:border-white/5 dark:bg-white/5"
        >
          <button
            type="button"
            class="rounded-xl bg-blue-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="!selectedZip || uploadLoading"
            @click="handleUploadZip"
          >
            <span v-if="uploadLoading">Proses...</span>
            <span v-else>Unggah</span>
          </button>
          <button
            type="button"
            class="rounded-xl border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 transition hover:bg-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:opacity-50 dark:border-slate-600 dark:text-slate-200 dark:hover:bg-slate-800"
            :disabled="uploadLoading"
            @click="showUploadDialog = false"
          >
            Batal
          </button>
        </div>
      </div>
    </div>
  </Transition>

  <VueEasyLightbox
    :visible="visibleLightbox"
    :imgs="lightboxImages"
    :index="0"
    @hide="visibleLightbox = false"
  />
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";
import VueEasyLightbox from "vue-easy-lightbox";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";
import DashboardBulkAction from "@/components/dashboard/DashboardBulkActionBar.vue";
import { useToast } from "@/composables/useToast";
import { useKelulusanStore } from "@/stores/kelulusan";
import {
  ArrowLeftIcon,
  ArrowUpTrayIcon,
  PhotoIcon,
  TrashIcon,
} from "@heroicons/vue/24/outline";

const kelulusanStore = useKelulusanStore();
const toast = useToast();

const showDeleteDialog = ref(false);
const showUploadDialog = ref(false);
const loadingAction = ref(false);

// Foto default fallback (dipakai kalau siswa belum punya fotolulus) TIDAK
// BOLEH terhapus, meski somehow tercantum di galeri -- diproteksi di sini
// (seleksi) sekaligus di backend (BulkDeleteFotoHandler) sbg lapis kedua.
const DEFAULT_FOTO = "peserta_didik.jpg";

const selectedIds = ref([]);
const selectableFoto = computed(() =>
  kelulusanStore.fotoList.filter((f) => f.gambar !== DEFAULT_FOTO),
);
const isAllSelected = computed(
  () =>
    selectableFoto.value.length > 0 &&
    selectableFoto.value.every((f) => selectedIds.value.includes(f.id)),
);
function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedIds.value = [];
  } else {
    selectedIds.value = selectableFoto.value.map((f) => f.id);
  }
}
function toggleSelect(id) {
  const idx = selectedIds.value.indexOf(id);
  if (idx === -1) {
    selectedIds.value.push(id);
  } else {
    selectedIds.value.splice(idx, 1);
  }
}

function fotoUrl(gambar) {
  return `/storage/fotokelulusan/${gambar}`;
}

// ===== Lightbox preview =====
const visibleLightbox = ref(false);
const lightboxImages = ref([]);
function openLightbox(foto) {
  lightboxImages.value = [fotoUrl(foto.gambar)];
  visibleLightbox.value = true;
}

async function handleDelete() {
  loadingAction.value = true;
  try {
    const res = await kelulusanStore.bulkDeleteFoto(selectedIds.value);
    toast.success(res?.message || "Foto terpilih berhasil dihapus.");
    selectedIds.value = [];
    showDeleteDialog.value = false;
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menghapus foto.");
  } finally {
    loadingAction.value = false;
  }
}

// ===== Upload ZIP =====
const zipInput = ref(null);
const selectedZip = ref(null);
const uploadLoading = ref(false);
const uploadError = ref("");

function handleZipChange(event) {
  const file = event.target.files?.[0];
  uploadError.value = "";
  if (!file) {
    selectedZip.value = null;
    return;
  }
  if (!file.name.toLowerCase().endsWith(".zip")) {
    uploadError.value = "File harus berformat .zip.";
    selectedZip.value = null;
    event.target.value = "";
    return;
  }
  const maxSize = 10 * 1024 * 1024;
  if (file.size > maxSize) {
    uploadError.value = "Ukuran zip maksimal 10 MB.";
    selectedZip.value = null;
    event.target.value = "";
    return;
  }
  selectedZip.value = file;
}

async function handleUploadZip() {
  if (!selectedZip.value) return;
  uploadLoading.value = true;
  uploadError.value = "";
  try {
    const res = await kelulusanStore.uploadFotoZip(selectedZip.value);
    toast.success(res?.message || "Foto berhasil diunggah.");
    showUploadDialog.value = false;
    selectedZip.value = null;
    if (zipInput.value) zipInput.value.value = "";
  } catch (error) {
    uploadError.value =
      error.response?.data?.message || "Gagal mengunggah foto.";
  } finally {
    uploadLoading.value = false;
  }
}

onMounted(() => {
  kelulusanStore.fetchFoto();
});
</script>
