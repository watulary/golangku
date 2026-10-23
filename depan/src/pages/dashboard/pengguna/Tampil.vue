<template>
  <DashboardPageHeader />
  <DashboardBulkActionBar
    :selected-count="penggunaStore.selectedIds.length"
    @activate="openActivateDialog"
    @deactivate="openDeactivateDialog"
    @reset-password="openResetPasswordDialog"
    @delete="openDeleteDialog"
    @clear="penggunaStore.clearSelection()"
  />
  <div class="space-y-6">
    <div
      class="relative overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
    >
      <div
        class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <!-- Header -->
      <div class="border-b border-gray-100 px-6 py-5 dark:border-white/5">
        <!-- ========================= -->
        <!-- Baris 1 -->
        <!-- ========================= -->

        <div
          class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/8 text-[#041983] dark:bg-white/5 dark:text-indigo-300"
            >
              <UsersIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Daftar Pengguna
            </h2>
          </div>

          <div class="relative w-64 max-w-55 sm:max-w-xs md:max-w-sm">
            <MagnifyingGlassIcon
              class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />

            <input
              v-model="penggunaStore.filters.q"
              @input="searchData"
              type="text"
              placeholder="Pencarian..."
              title="Bisa cari nama, email, atau tanggal log (cth: 2026 / 08-2026 / 19-08-2026)"
              class="w-full rounded-xl border border-gray-300 bg-white py-2.5 pl-10 pr-4 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-100 dark:placeholder:text-slate-500"
            />
          </div>
        </div>

        <!-- ========================= -->
        <!-- Baris 2 -->
        <!-- ========================= -->

        <div
          class="mt-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
        >
          <!-- kiri -->
          <div class="flex flex-wrap items-center gap-4">
            <div class="flex items-center gap-2">
              <span class="text-sm text-slate-500 dark:text-slate-400">
                Tampilkan
              </span>

              <select
                :value="penggunaStore.pagination.limit"
                @change="penggunaStore.changeLimit($event.target.value)"
                class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
              >
                <option :value="10">10</option>
                <option :value="25">25</option>
                <option :value="50">50</option>
                <option :value="100">100</option>
                <option :value="500">500</option>
                <option :value="1000">1000</option>
              </select>
            </div>

            <!-- Filter Level -->
            <div class="flex items-center gap-2">
              <span class="text-sm text-slate-500 dark:text-slate-400">
                Level
              </span>

              <select
                v-model="penggunaStore.filters.level"
                @change="filterData"
                class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
              >
                <option value="">Semua Level</option>
                <option value="Admin">Super Admin</option>
                <option value="Staf">Admin</option>
                <option value="Pegawai">Pegawai</option>
                <option value="Pesertadidik">Peserta Didik</option>
              </select>
            </div>
          </div>

          <!-- kanan -->
          <div class="flex items-center gap-2">
            <!-- Tambah Pengguna -->
            <RouterLink
              to="/dashboard/pengguna/tambah"
              title="Tambah Pengguna"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-blue-500 hover:bg-blue-50 hover:text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-blue-500 dark:hover:bg-blue-500/10 dark:hover:text-blue-400"
            >
              <PlusIcon class="h-5 w-5" />
            </RouterLink>

            <!-- Upload Pengguna -->
            <button
              type="button"
              @click="showUploadDialog = true"
              title="Unggah Pengguna"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-yellow-500 hover:bg-yellow-50 hover:text-yellow-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-yellow-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-yellow-500 dark:hover:bg-yellow-500/10 dark:hover:text-yellow-300"
            >
              <ArrowUpTrayIcon class="h-5 w-5" />
            </button>

            <!-- Export Excel -->
            <button
              type="button"
              @click="exportExcel"
              title="Unduh Pengguna"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-emerald-500 dark:hover:bg-emerald-500/10 dark:hover:text-emerald-400"
            >
              <ArrowDownTrayIcon class="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>

      <!-- Table -->
      <div v-if="!loading && pengguna.length === 0" class="py-16 text-center">
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <MagnifyingGlassIcon v-if="isSearching" class="h-10 w-10" />
          <InboxIcon v-else class="h-10 w-10" />
          <p class="text-sm font-medium">
            {{
              isSearching ? "Data tidak ditemukan." : "Belum ada data pengguna."
            }}
          </p>
        </div>
      </div>

      <div v-else class="max-h-[70vh] overflow-auto">
        <table class="min-w-275 w-full">
          <thead class="sticky top-0 z-20 bg-gray-50/90 dark:bg-[#0B0F1A]">
            <tr>
              <th scope="col" class="w-12 px-4 py-3 text-center">
                <input
                  ref="selectAllCheckbox"
                  type="checkbox"
                  :checked="penggunaStore.isAllSelected"
                  @change="penggunaStore.toggleSelectAll()"
                  class="h-4 w-4 cursor-pointer rounded border-gray-300 text-emerald-600 focus:ring-emerald-500 dark:border-slate-600 dark:bg-[#0B0F1A]"
                />
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Foto
              </th>

              <th scope="col" class="px-4 py-3 text-left">
                <button
                  type="button"
                  @click="toggleSort('nama_pengguna')"
                  class="flex items-center gap-1 text-[11px] font-semibold tracking-wider text-gray-500 transition hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200"
                >
                  Nama Pengguna
                  <component
                    :is="sortIcon('nama_pengguna')"
                    class="h-3.5 w-3.5"
                  />
                </button>
              </th>

              <th scope="col" class="px-4 py-3 text-left">
                <button
                  type="button"
                  @click="toggleSort('email')"
                  class="flex items-center gap-1 text-[11px] font-semibold tracking-wider text-gray-500 transition hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200"
                >
                  Email
                  <component :is="sortIcon('email')" class="h-3.5 w-3.5" />
                </button>
              </th>

              <th scope="col" class="px-4 py-3 text-left">
                <button
                  type="button"
                  @click="toggleSort('level')"
                  class="flex items-center gap-1 text-[11px] font-semibold tracking-wider text-gray-500 transition hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200"
                >
                  Level
                  <component :is="sortIcon('level')" class="h-3.5 w-3.5" />
                </button>
              </th>

              <th scope="col" class="px-4 py-3 text-left">
                <button
                  type="button"
                  @click="toggleSort('aktif')"
                  class="flex items-center gap-1 text-[11px] font-semibold tracking-wider text-gray-500 transition hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200"
                >
                  Status
                  <component :is="sortIcon('aktif')" class="h-3.5 w-3.5" />
                </button>
              </th>

              <th scope="col" class="px-4 py-3 text-left">
                <button
                  type="button"
                  @click="toggleSort('password_changed')"
                  class="flex items-center gap-1 text-[11px] font-semibold tracking-wider text-gray-500 transition hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200"
                >
                  Password
                  <component
                    :is="sortIcon('password_changed')"
                    class="h-3.5 w-3.5"
                  />
                </button>
              </th>

              <th scope="col" class="px-4 py-3 text-left">
                <button
                  type="button"
                  @click="toggleSort('last_login')"
                  class="flex items-center gap-1 text-[11px] font-semibold tracking-wider text-gray-500 transition hover:text-gray-700 dark:text-slate-400 dark:hover:text-slate-200"
                >
                  Tanggal Log
                  <component :is="sortIcon('last_login')" class="h-3.5 w-3.5" />
                </button>
              </th>
            </tr>
          </thead>

          <tbody class="divide-y divide-gray-100 dark:divide-white/5">
            <!-- Skeleton saat memuat -->
            <template v-if="loading">
              <tr v-for="n in 8" :key="`skeleton-${n}`">
                <td class="px-4 py-3 text-center">
                  <div
                    class="mx-auto h-4 w-4 animate-pulse rounded bg-gray-100 dark:bg-white/5"
                  ></div>
                </td>
                <td class="px-4 py-3">
                  <div
                    class="h-9 w-9 animate-pulse rounded-full bg-gray-100 dark:bg-white/5"
                  ></div>
                </td>
                <td v-for="c in 6" :key="c" class="px-4 py-3">
                  <div
                    class="h-4 w-20 animate-pulse rounded bg-gray-100 dark:bg-white/5"
                  ></div>
                </td>
              </tr>
            </template>

            <tr
              v-else
              v-for="item in pengguna"
              :key="item.id_pengguna"
              class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
            >
              <td class="w-12 px-4 py-2.5 text-center">
                <input
                  type="checkbox"
                  :checked="
                    penggunaStore.selectedIds.includes(item.id_pengguna)
                  "
                  @change="penggunaStore.toggleSelection(item.id_pengguna)"
                  class="h-4 w-4 cursor-pointer rounded border-gray-300 text-emerald-600 focus:ring-emerald-500 dark:border-slate-600 dark:bg-[#0B0F1A]"
                />
              </td>
              <!-- Foto -->
              <td class="px-4 py-2.5">
                <button
                  v-if="item.foto_url"
                  type="button"
                  title="Klik untuk perbesar"
                  class="group block cursor-pointer"
                  @click="openLightbox(item)"
                >
                  <img
                    :src="item.foto_url"
                    class="h-9 w-9 rounded-full object-cover transition duration-200 group-hover:opacity-80"
                  />
                </button>

                <div
                  v-else
                  class="flex h-9 w-9 items-center justify-center rounded-full bg-gray-200 font-semibold text-gray-700 dark:bg-slate-700 dark:text-white"
                >
                  {{ item.nama_pengguna.charAt(0) }}
                </div>
              </td>

              <!-- Nama -->
              <td
                class="px-4 py-2.5 text-xs font-medium text-gray-600 dark:text-slate-200"
              >
                <RouterLink
                  :to="`/dashboard/pengguna/ubah/${item.id_pengguna}`"
                  class="hover:text-[#041983] hover:underline dark:hover:text-indigo-300"
                >
                  {{ item.nama_pengguna }}
                </RouterLink>
              </td>

              <!-- Email -->
              <td
                class="px-4 py-2.5 text-sm leading-5 text-gray-500 dark:text-slate-400"
              >
                {{ item.email }}
              </td>

              <!-- Level -->
              <td class="px-4 py-2.5">
                <span
                  class="inline-flex items-center rounded-full px-2.5 py-1 text-xs leading-5"
                  :class="levelBadge(item.level).class"
                >
                  {{ levelBadge(item.level).text }}
                </span>
              </td>

              <!-- Status -->
              <td class="px-4 py-2.5">
                <span
                  v-if="item.aktif"
                  class="rounded-full bg-emerald-100 px-3 py-1 text-xs leading-5 text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300"
                >
                  Aktif
                </span>

                <span
                  v-else
                  class="rounded-full bg-red-100 px-3 py-1 text-xs leading-5 text-red-700 dark:bg-red-500/15 dark:text-red-300"
                >
                  Tidak Aktif
                </span>
              </td>

              <!-- Password -->
              <td class="px-4 py-2.5">
                <span
                  v-if="item.password_changed"
                  class="rounded-full bg-emerald-100 px-3 py-1 text-xs leading-5 text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300"
                >
                  Diubah
                </span>

                <span
                  v-else
                  class="rounded-full bg-amber-100 px-3 py-1 text-xs leading-5 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300"
                >
                  Standar
                </span>
              </td>

              <!-- Login -->
              <td
                class="px-4 py-2.5 text-xs leading-5 text-gray-500 dark:text-slate-400"
              >
                {{ formatDateTimeNumeric(item.last_login) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <DashboardPagination
        v-if="!loading && pengguna.length > 0"
        :pagination="penggunaStore.pagination"
        @change-page="changePage"
      />
    </div>
  </div>
  <DashboardConfirmDialog
    v-model="showActivateDialog"
    title="Aktifkan Pengguna"
    :message="`Anda akan mengaktifkan ${penggunaStore.selectedIds.length} pengguna.\n\nMereka dapat login kembali.`"
    confirm-text="Aktifkan"
    confirm-color="green"
    :loading="loadingAction"
    @confirm="bulkActivate"
  />
  <DashboardConfirmDialog
    v-model="showDeactivateDialog"
    title="Nonaktifkan Pengguna"
    :message="`Anda akan menonaktifkan ${penggunaStore.selectedIds.length} pengguna.\n\nMereka tidak akan bisa login sampai diaktifkan kembali.`"
    confirm-text="Nonaktifkan"
    confirm-color="amber"
    :loading="loadingAction"
    @confirm="bulkDeactivate"
  />
  <DashboardConfirmDialog
    v-model="showResetPasswordDialog"
    title="Reset Password"
    :message="`Anda akan mereset password ${penggunaStore.selectedIds.length} pengguna.\n\nPassword kembali ke password default (Abcd1234*).\n\nMereka diwajibkan mengganti password setelah login.`"
    confirm-text="Reset Password"
    confirm-color="blue"
    :loading="loadingAction"
    @confirm="bulkResetPassword"
  />
  <DashboardConfirmDialog
    v-model="showDeleteDialog"
    title="Hapus Pengguna"
    :message="`Anda akan menghapus ${penggunaStore.selectedIds.length} pengguna.\n\nData yang telah dihapus tidak dapat dikembalikan.`"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="loadingAction"
    @confirm="bulkDelete"
  />
  <DashboardUploadDialog
    v-model="showUploadDialog"
    title="Unggah Pengguna"
    subtitle="Unggah data pengguna menggunakan Template (.xlsx)"
    :upload-action="penggunaStore.uploadExcel"
    @success="toast.success($event)"
    @download-template="downloadTemplate"
  />

  <VueEasyLightbox
    :visible="visibleLightbox"
    :imgs="lightboxImages"
    :index="0"
    @hide="visibleLightbox = false"
  />
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { RouterLink } from "vue-router";
import { storeToRefs } from "pinia";
import VueEasyLightbox from "vue-easy-lightbox";
import api from "@/services/api";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import DashboardBulkActionBar from "@/components/dashboard/DashboardBulkActionBar.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";
import DashboardUploadDialog from "@/components/dashboard/DashboardUploadDialog.vue";
import { usePenggunaStore } from "@/stores/pengguna";
import { useToast } from "@/composables/useToast";
import { formatDateTimeNumeric } from "@/utils/date";
import {
  MagnifyingGlassIcon,
  ArrowDownTrayIcon,
  ArrowUpTrayIcon,
  PlusIcon,
  UsersIcon,
  InboxIcon,
  ChevronUpIcon,
  ChevronDownIcon,
  ChevronUpDownIcon,
} from "@heroicons/vue/24/outline";

const penggunaStore = usePenggunaStore();

const toast = useToast();

const { pengguna, loading } = storeToRefs(penggunaStore);

const { fetchData } = penggunaStore;

const isSearching = computed(() => {
  return (
    penggunaStore.filters.q.trim() !== "" || penggunaStore.filters.level !== ""
  );
});

// ===== Lightbox foto (pola sama spt GaleriFoto.vue / pesertadidik Tampil.vue) =====
const visibleLightbox = ref(false);
const lightboxImages = ref([]);
function openLightbox(item) {
  if (!item.foto_url) return;
  lightboxImages.value = [item.foto_url];
  visibleLightbox.value = true;
}

const selectAllCheckbox = ref(null);

// =========================================
// Confirmation Dialog
// =========================================

const showActivateDialog = ref(false);

const showDeactivateDialog = ref(false);

const showResetPasswordDialog = ref(false);

const showDeleteDialog = ref(false);

const showUploadDialog = ref(false);

const loadingAction = ref(false);

onMounted(async () => {
  await fetchData();
});

// =========================================
// Open Dialog
// =========================================

function openActivateDialog() {
  if (penggunaStore.selectedIds.length === 0) return;

  showActivateDialog.value = true;
}

function openDeactivateDialog() {
  if (penggunaStore.selectedIds.length === 0) return;

  showDeactivateDialog.value = true;
}

function openResetPasswordDialog() {
  if (penggunaStore.selectedIds.length === 0) return;

  showResetPasswordDialog.value = true;
}

function openDeleteDialog() {
  if (penggunaStore.selectedIds.length === 0) return;

  showDeleteDialog.value = true;
}

// =========================================
// Upload Excel
// =========================================

async function downloadTemplate() {
  try {
    const response = await api.get("/pengguna/template", {
      responseType: "blob",
    });

    const disposition = response.headers["content-disposition"];

    let filename = "template-pengguna.xlsx";

    if (disposition) {
      const match = disposition.match(/filename="?([^"]+)"?/);

      if (match) {
        filename = match[1];
      }
    }

    const url = window.URL.createObjectURL(new Blob([response.data]));

    const link = document.createElement("a");

    link.href = url;

    link.download = filename;

    document.body.appendChild(link);

    link.click();

    link.remove();

    window.URL.revokeObjectURL(url);
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal mengunduh template.");
  }
}

// =========================================
// Bulk Activate
// =========================================

async function bulkActivate() {
  loadingAction.value = true;

  try {
    const response = await penggunaStore.bulkActivate();

    showActivateDialog.value = false;

    toast.success(response.message);
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal mengaktifkan pengguna.",
    );
  } finally {
    loadingAction.value = false;
  }
}

// =========================================
// Bulk Deactivate
// =========================================

async function bulkDeactivate() {
  loadingAction.value = true;

  try {
    const response = await penggunaStore.bulkDeactivate();

    showDeactivateDialog.value = false;

    toast.success(response.message);
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menonaktifkan pengguna.",
    );
  } finally {
    loadingAction.value = false;
  }
}

async function bulkResetPassword() {
  loadingAction.value = true;

  try {
    const response = await penggunaStore.bulkResetPassword();

    showResetPasswordDialog.value = false;

    toast.success(response.message);
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal mereset password.");
  } finally {
    loadingAction.value = false;
  }
}

// =========================================
// Bulk Delete
// =========================================

async function bulkDelete() {
  loadingAction.value = true;

  try {
    const response = await penggunaStore.bulkDelete();

    toast.success(response.message);
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menghapus pengguna.");
  } finally {
    penggunaStore.clearSelection();
    showDeleteDialog.value = false;
    loadingAction.value = false;
  }
}

async function exportExcel() {
  try {
    await penggunaStore.exportExcel();

    toast.success("Data pengguna berhasil diekspor.");
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal mengekspor data pengguna.",
    );
  }
}

let searchTimeout = null;

function searchData() {
  clearTimeout(searchTimeout);

  searchTimeout = setTimeout(async () => {
    penggunaStore.pagination.page = 1;

    await penggunaStore.fetchData();
  }, 300);
}

// =========================================
// Filter Level
// =========================================

async function filterData() {
  penggunaStore.pagination.page = 1;

  await penggunaStore.fetchData();
}

// =========================================
// Sort Kolom
// =========================================

async function toggleSort(column) {
  if (penggunaStore.filters.sort === column) {
    penggunaStore.filters.order =
      penggunaStore.filters.order === "asc" ? "desc" : "asc";
  } else {
    penggunaStore.filters.sort = column;
    penggunaStore.filters.order = "asc";
  }

  penggunaStore.pagination.page = 1;

  await penggunaStore.fetchData();
}

function sortIcon(column) {
  if (penggunaStore.filters.sort !== column) {
    return ChevronUpDownIcon;
  }

  return penggunaStore.filters.order === "asc"
    ? ChevronUpIcon
    : ChevronDownIcon;
}

function levelBadge(level) {
  switch (level) {
    case "Admin":
      return {
        text: "Super Admin",
        class:
          "bg-red-100 text-red-700 ring-1 ring-red-200 dark:bg-red-500/10 dark:text-red-300 dark:ring-red-500/20",
      };

    case "Staf":
      return {
        text: "Admin",
        class:
          "bg-blue-100 text-blue-700 ring-1 ring-blue-200 dark:bg-blue-500/10 dark:text-blue-300 dark:ring-blue-500/20",
      };

    case "Pegawai":
      return {
        text: "Pegawai",
        class:
          "bg-emerald-100 text-emerald-700 ring-1 ring-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-300 dark:ring-emerald-500/20",
      };

    case "Pesertadidik":
      return {
        text: "Peserta Didik",
        class:
          "bg-violet-100 text-violet-700 ring-1 ring-violet-200 dark:bg-violet-500/10 dark:text-violet-300 dark:ring-violet-500/20",
      };

    default:
      return {
        text: level,
        class:
          "bg-slate-100 text-slate-700 ring-1 ring-slate-200 dark:bg-slate-700/30 dark:text-slate-300 dark:ring-slate-600",
      };
  }
}

async function changePage(page) {
  penggunaStore.pagination.page = page;
  await penggunaStore.fetchData();
}

watch(
  () => penggunaStore.isIndeterminate,
  (value) => {
    if (selectAllCheckbox.value) {
      selectAllCheckbox.value.indeterminate = value;
    }
  },
  { immediate: true },
);
</script>
