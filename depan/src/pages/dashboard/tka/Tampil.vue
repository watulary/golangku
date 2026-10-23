<template>
  <DashboardPageHeader />

  <DashboardBulkAction
    :selected-count="selectedNis.length"
    @clear="selectedNis = []"
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
        class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <div class="border-b border-gray-100 px-6 py-5 dark:border-white/5">
        <div
          class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/8 text-[#041983] dark:bg-white/5 dark:text-indigo-300"
            >
              <ClipboardDocumentCheckIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Data Nilai TKA
            </h2>
          </div>
          <div class="relative w-64 max-w-55 sm:max-w-xs md:max-w-sm">
            <MagnifyingGlassIcon
              class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />
            <input
              v-model="tkaStore.filters.q"
              @input="searchData"
              type="text"
              placeholder="Pencarian..."
              class="w-full rounded-xl border border-gray-300 bg-white py-2.5 pl-10 pr-4 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500"
            />
          </div>
        </div>
        <div
          class="mt-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
        >
          <div class="flex flex-wrap items-center gap-3">
            <div class="flex items-center gap-2">
              <span class="text-sm text-slate-500 dark:text-slate-400">
                Tampilkan
              </span>
              <select
                :value="tkaStore.pagination.limit"
                @change="tkaStore.changeLimit($event.target.value)"
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
          </div>
          <div class="flex items-center gap-2">
            <RouterLink
              :to="{ name: 'dashboard-tka-pengumuman' }"
              title="Pengaturan Pengumuman"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-fuchsia-500 hover:bg-fuchsia-50 hover:text-fuchsia-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-fuchsia-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-fuchsia-500 dark:hover:bg-fuchsia-500/10 dark:hover:text-fuchsia-400"
            >
              <MegaphoneIcon class="h-5 w-5" />
            </RouterLink>
            <button
              type="button"
              title="Unggah Data TKA"
              @click="showUploadDialog = true"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-yellow-500 hover:bg-yellow-50 hover:text-yellow-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-yellow-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-yellow-500 dark:hover:bg-yellow-500/10 dark:hover:text-yellow-300"
            >
              <ArrowUpTrayIcon class="h-5 w-5" />
            </button>
            <button
              type="button"
              title="Unduh Data TKA"
              @click="handleExport"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-emerald-500 dark:hover:bg-emerald-500/10 dark:hover:text-emerald-400"
            >
              <ArrowDownTrayIcon class="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>

      <div
        v-if="!tkaStore.adminLoading && tkaStore.adminTka.length === 0"
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <MagnifyingGlassIcon v-if="isSearching" class="h-10 w-10" />
          <InboxIcon v-else class="h-10 w-10" />
          <p class="text-sm font-medium">
            {{ isSearching ? "Data tidak ditemukan." : "Belum ada data TKA." }}
          </p>
        </div>
      </div>

      <div v-else>
        <div class="max-h-[70vh] overflow-auto">
          <table class="min-w-225 w-full">
            <thead class="sticky top-0 z-20 bg-gray-50/90 dark:bg-[#0B0F1A]">
              <tr>
                <th class="w-10 px-4 py-3">
                  <input
                    type="checkbox"
                    :checked="isAllSelected"
                    @change="toggleSelectAll"
                    class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="tkaStore.toggleSort('nama_rombel')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Rombel
                    <component :is="sortIcon('nama_rombel')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="tkaStore.toggleSort('nis')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    NIS
                    <component :is="sortIcon('nis')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="tkaStore.toggleSort('nopes')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    No Peserta
                    <component :is="sortIcon('nopes')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="tkaStore.toggleSort('nama_siswa')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Nama Peserta Didik
                    <component :is="sortIcon('nama_siswa')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="tkaStore.toggleSort('mtk')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Matematika
                    <component :is="sortIcon('mtk')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="tkaStore.toggleSort('bind')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Bhs Indonesia
                    <component :is="sortIcon('bind')" class="h-4 w-4" />
                  </button>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-white/5">
              <!-- Skeleton saat memuat -->
              <template v-if="tkaStore.adminLoading">
                <tr v-for="n in 6" :key="`skeleton-${n}`">
                  <td v-for="c in 7" :key="c" class="px-4 py-3">
                    <div
                      class="h-4 w-20 animate-pulse rounded bg-gray-100 dark:bg-white/5"
                    ></div>
                  </td>
                </tr>
              </template>

              <tr
                v-else
                v-for="item in tkaStore.adminTka"
                :key="item.nis"
                class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
              >
                <td class="px-4 py-2.5">
                  <input
                    type="checkbox"
                    :checked="selectedNis.includes(item.nis)"
                    @change="toggleSelect(item.nis)"
                    class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.nama_rombel || "-" }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.nis }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.nopes }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
                >
                  {{ item.nama_siswa || "-" }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.mtk }}
                  <span class="text-gray-400">({{ item.ket_mtk }})</span>
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.bind }}
                  <span class="text-gray-400">({{ item.ket_bind }})</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <DashboardPagination
          :pagination="tkaStore.pagination"
          @change-page="changePage"
        />
      </div>
    </div>
  </div>

  <DashboardConfirmDialog
    v-model="showDeleteDialog"
    title="Hapus Data TKA"
    :message="`Anda akan menghapus ${selectedNis.length} data TKA.\n\nData yang telah dihapus tidak dapat dikembalikan.`"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="loadingAction"
    @confirm="handleDelete"
  />

  <DashboardUploadDialog
    v-model="showUploadDialog"
    title="Unggah Data TKA"
    subtitle="Unggah data nilai TKA menggunakan Template (.xlsx)"
    :upload-action="tkaStore.importTka"
    @download-template="handleDownloadTemplate"
    @success="handleImportSuccess"
  />
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";
import DashboardUploadDialog from "@/components/dashboard/DashboardUploadDialog.vue";
import DashboardBulkAction from "@/components/dashboard/DashboardBulkActionBar.vue";
import { useToast } from "@/composables/useToast";
import { useTkaStore } from "@/stores/tka";
import {
  ChevronDownIcon,
  ChevronUpIcon,
  ChevronUpDownIcon,
  MagnifyingGlassIcon,
  TrashIcon,
  ArrowDownTrayIcon,
  ArrowUpTrayIcon,
  MegaphoneIcon,
  ClipboardDocumentCheckIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const tkaStore = useTkaStore();
const toast = useToast();

const showDeleteDialog = ref(false);
const showUploadDialog = ref(false);
const loadingAction = ref(false);

const isSearching = computed(() => tkaStore.filters.q.trim() !== "");

const selectedNis = ref([]);
const isAllSelected = computed(
  () =>
    tkaStore.adminTka.length > 0 &&
    tkaStore.adminTka.every((item) => selectedNis.value.includes(item.nis)),
);
function toggleSelect(nis) {
  const idx = selectedNis.value.indexOf(nis);
  if (idx === -1) {
    selectedNis.value.push(nis);
  } else {
    selectedNis.value.splice(idx, 1);
  }
}
function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedNis.value = selectedNis.value.filter(
      (nis) => !tkaStore.adminTka.some((item) => item.nis === nis),
    );
  } else {
    const pageNis = tkaStore.adminTka.map((item) => item.nis);
    selectedNis.value = [...new Set([...selectedNis.value, ...pageNis])];
  }
}

let searchTimeout = null;
function searchData() {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(async () => {
    tkaStore.pagination.page = 1;
    await tkaStore.fetchAdminTka();
  }, 300);
}
async function changePage(page) {
  tkaStore.pagination.page = page;
  await tkaStore.fetchAdminTka();
}
function sortIcon(column) {
  if (tkaStore.sortBy !== column) return ChevronUpDownIcon;
  return tkaStore.sortDir === "asc" ? ChevronUpIcon : ChevronDownIcon;
}

async function handleExport() {
  try {
    await tkaStore.exportTka();
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal export data TKA.");
  }
}
async function handleDownloadTemplate() {
  try {
    await tkaStore.downloadTemplate();
  } catch (error) {
    toast.error("Gagal mengunduh template import.");
  }
}
function handleImportSuccess(message) {
  toast.success(message || "Data TKA berhasil diimport.");
}

async function handleDelete() {
  loadingAction.value = true;
  try {
    const res = await tkaStore.bulkDelete(selectedNis.value);
    toast.success(res?.message || "Data TKA berhasil dihapus.");
    selectedNis.value = [];
    showDeleteDialog.value = false;
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menghapus data TKA.");
  } finally {
    loadingAction.value = false;
  }
}

onMounted(() => {
  tkaStore.fetchAdminTka();
});
</script>
