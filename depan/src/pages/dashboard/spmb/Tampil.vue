<template>
  <DashboardPageHeader />
  <DashboardBulkAction
    :selected-count="selectedNoPeserta.length"
    @clear="selectedNoPeserta = []"
  >
    <template #actions>
      <button
        type="button"
        @click="openBulkLaporDiri"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-amber-500 px-4 py-2 text-sm font-medium text-white transition hover:bg-amber-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500/40"
      >
        <PencilSquareIcon class="h-5 w-5" />
        Ubah Lapor Diri
      </button>
      <button
        type="button"
        @click="showDeleteDialog = true"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-red-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-red-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40"
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

      <!-- Header -->
      <div class="border-b border-gray-100 px-6 py-5 dark:border-white/5">
        <!-- Baris 1 -->
        <div
          class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/8 text-[#041983] dark:bg-white/5 dark:text-indigo-300"
            >
              <UserGroupIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Daftar Peserta SPMB
            </h2>
          </div>
          <div class="relative w-64 max-w-55 sm:max-w-xs md:max-w-sm">
            <MagnifyingGlassIcon
              class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />
            <input
              v-model="spmbStore.filters.q"
              @input="searchData"
              type="text"
              placeholder="Pencarian..."
              class="w-full rounded-xl border border-gray-300 bg-white py-2.5 pl-10 pr-4 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500"
            />
          </div>
        </div>
        <!-- Baris 2 -->
        <div
          class="mt-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
        >
          <!-- kiri -->
          <div class="flex flex-wrap items-center gap-3">
            <div class="flex items-center gap-2">
              <span class="text-sm text-slate-500 dark:text-slate-400">
                Tampilkan
              </span>
              <select
                :value="spmbStore.pagination.limit"
                @change="spmbStore.changeLimit($event.target.value)"
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
            <select
              v-model="jalurFilter"
              @change="changeJalurFilter"
              class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
            >
              <option :value="0">Semua Jalur</option>
              <option
                v-for="opt in JALUR_OPTIONS"
                :key="opt.value"
                :value="opt.value"
              >
                {{ opt.label }}
              </option>
            </select>
            <select
              v-model="laporDiriFilter"
              @change="changeLaporDiriFilter"
              class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
            >
              <option :value="-1">Semua Status</option>
              <option
                v-for="opt in LAPOR_DIRI_OPTIONS"
                :key="opt.value"
                :value="opt.value"
              >
                {{ opt.label }}
              </option>
            </select>
          </div>
          <!-- kanan -->
          <div class="flex items-center gap-2">
            <RouterLink
              :to="{ name: 'dashboard-spmb-tambah' }"
              title="Tambah Peserta SPMB"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-blue-500 hover:bg-blue-50 hover:text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-blue-500 dark:hover:bg-blue-500/10 dark:hover:text-blue-400"
            >
              <PlusIcon class="h-5 w-5" />
            </RouterLink>
            <button
              type="button"
              title="Unggah Peserta SPMB"
              @click="showUploadDialog = true"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-yellow-500 hover:bg-yellow-50 hover:text-yellow-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-yellow-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-yellow-500 dark:hover:bg-yellow-500/10 dark:hover:text-yellow-300"
            >
              <ArrowUpTrayIcon class="h-5 w-5" />
            </button>
            <button
              type="button"
              title="Unduh Peserta SPMB"
              @click="handleExport"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-emerald-500 dark:hover:bg-emerald-500/10 dark:hover:text-emerald-400"
            >
              <ArrowDownTrayIcon class="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>

      <!-- Loading -->
      <div
        v-if="spmbStore.adminLoading"
        class="flex flex-col items-center justify-center py-20"
      >
        <div
          class="h-10 w-10 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"
        />
        <p class="mt-4 text-sm text-gray-500 dark:text-slate-400">
          Memuat data...
        </p>
      </div>

      <!-- Empty -->
      <div
        v-if="!spmbStore.adminLoading && spmbStore.adminSpmb.length === 0"
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <MagnifyingGlassIcon v-if="isSearching" class="h-10 w-10" />
          <InboxIcon v-else class="h-10 w-10" />
          <p class="text-sm font-medium">
            {{
              isSearching
                ? "Data tidak ditemukan."
                : "Belum ada data peserta SPMB."
            }}
          </p>
        </div>
      </div>

      <!-- Table -->
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
                    @click="spmbStore.toggleSort('no_peserta')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    No. Peserta
                    <component :is="sortIcon('no_peserta')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="spmbStore.toggleSort('nama_peserta')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Nama Peserta
                    <component :is="sortIcon('nama_peserta')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="spmbStore.toggleSort('kelamin')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    L/P
                    <component :is="sortIcon('kelamin')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="spmbStore.toggleSort('asal_sekolah')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Asal Sekolah
                    <component :is="sortIcon('asal_sekolah')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="spmbStore.toggleSort('jalur')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Jalur
                    <component :is="sortIcon('jalur')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="spmbStore.toggleSort('lapor_diri')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Lapor Diri
                    <component :is="sortIcon('lapor_diri')" class="h-4 w-4" />
                  </button>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-white/5">
              <!-- Skeleton saat memuat -->
              <template v-if="spmbStore.adminLoading">
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
                v-for="item in spmbStore.adminSpmb"
                :key="item.no_peserta"
                class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
              >
                <td class="px-4 py-2.5">
                  <input
                    type="checkbox"
                    :checked="selectedNoPeserta.includes(item.no_peserta)"
                    @change="toggleSelect(item.no_peserta)"
                    class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.no_peserta }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
                >
                  <RouterLink
                    :to="{
                      name: 'dashboard-spmb-edit',
                      params: { id: item.no_peserta },
                    }"
                    class="hover:text-[#041983] hover:underline dark:hover:text-indigo-300"
                  >
                    {{ item.nama_peserta }}
                  </RouterLink>
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.kelamin }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.asal_sekolah }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ jalurLabel(item.jalur) }}
                </td>
                <td class="px-4 py-2.5 text-sm">
                  <span
                    v-if="laporDiriOption(item.lapor_diri)"
                    class="inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium"
                    :class="laporDiriOption(item.lapor_diri).badge"
                  >
                    {{ laporDiriOption(item.lapor_diri).label }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <DashboardPagination
          :pagination="spmbStore.pagination"
          @change-page="changePage"
        />
      </div>
    </div>
  </div>

  <!-- Dialog Delete -->
  <DashboardConfirmDialog
    v-model="showDeleteDialog"
    title="Hapus Peserta SPMB"
    :message="`Anda akan menghapus ${selectedNoPeserta.length} data peserta SPMB.\n\nData yang telah dihapus tidak dapat dikembalikan.`"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="loadingAction"
    @confirm="handleDelete"
  />

  <!-- Modal Ubah Lapor Diri -->
  <div
    v-if="showLaporDiriDialog"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
  >
    <div
      class="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl ring-1 ring-slate-900/5 dark:bg-slate-900 dark:ring-white/10"
    >
      <h3 class="text-lg font-semibold text-gray-800 dark:text-slate-200">
        Ubah Status Lapor Diri
      </h3>
      <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
        {{ selectedNoPeserta.length }} peserta dipilih
      </p>
      <form @submit.prevent="handleBulkLaporDiri" class="mt-4 space-y-4">
        <div>
          <label class="text-sm font-medium text-gray-700 dark:text-slate-300">
            Status Lapor Diri
          </label>
          <select
            v-model.number="laporDiriForm.lapor_diri"
            required
            class="mt-1 w-full rounded-xl border border-gray-300 px-3 py-2 text-sm outline-none transition focus:border-amber-500 focus:ring-2 focus:ring-amber-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200"
          >
            <option
              v-for="opt in LAPOR_DIRI_OPTIONS"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.label }}
            </option>
          </select>
        </div>
        <div class="flex justify-end gap-2 pt-2">
          <button
            type="button"
            @click="showLaporDiriDialog = false"
            class="rounded-xl border border-gray-300 px-4 py-2 text-sm font-medium text-gray-600 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
          >
            Batal
          </button>
          <button
            type="submit"
            :disabled="loadingAction"
            class="rounded-xl bg-amber-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-amber-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500/40 disabled:opacity-50"
          >
            Simpan
          </button>
        </div>
      </form>
    </div>
  </div>

  <DashboardUploadDialog
    v-model="showUploadDialog"
    title="Unggah Peserta SPMB"
    subtitle="Unggah data peserta SPMB menggunakan Template (.xlsx)"
    :upload-action="spmbStore.importSpmb"
    @download-template="handleDownloadTemplate"
    @success="handleImportSuccess"
  />
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { RouterLink } from "vue-router";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";
import DashboardUploadDialog from "@/components/dashboard/DashboardUploadDialog.vue";
import DashboardBulkAction from "@/components/dashboard/DashboardBulkActionBar.vue";
import { useToast } from "@/composables/useToast";
import { useSpmbStore } from "@/stores/spmb";
import {
  ChevronDownIcon,
  ChevronUpIcon,
  ChevronUpDownIcon,
  MagnifyingGlassIcon,
  PencilSquareIcon,
  PlusIcon,
  TrashIcon,
  ArrowDownTrayIcon,
  ArrowUpTrayIcon,
  UserGroupIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const spmbStore = useSpmbStore();
const toast = useToast();

const showDeleteDialog = ref(false);
const showUploadDialog = ref(false);
const showLaporDiriDialog = ref(false);
const loadingAction = ref(false);

const isSearching = computed(() => spmbStore.filters.q.trim() !== "");

const JALUR_OPTIONS = [
  { value: 1, label: "Prestasi Akademik" },
  { value: 2, label: "Prestasi Non Akademik" },
  { value: 3, label: "Penyandang Disabilitas" },
  { value: 4, label: "KJP Plus" },
  { value: 5, label: "Domisili" },
  { value: 6, label: "Mutasi" },
  { value: 7, label: "Tahap Kedua" },
];

const LAPOR_DIRI_OPTIONS = [
  {
    value: 0,
    label: "Belum Lapor",
    badge:
      "bg-yellow-50 text-yellow-700 dark:bg-yellow-500/10 dark:text-yellow-400",
  },
  {
    value: 1,
    label: "Sudah Lapor",
    badge:
      "bg-green-50 text-green-700 dark:bg-green-500/10 dark:text-green-400",
  },
  {
    value: 2,
    label: "Undur Diri",
    badge: "bg-red-50 text-red-700 dark:bg-red-500/10 dark:text-red-400",
  },
];

function jalurLabel(val) {
  return JALUR_OPTIONS.find((o) => o.value === val)?.label || "-";
}
function laporDiriOption(val) {
  return LAPOR_DIRI_OPTIONS.find((o) => o.value === val);
}

const selectedNoPeserta = ref([]);
const isAllSelected = computed(
  () =>
    spmbStore.adminSpmb.length > 0 &&
    spmbStore.adminSpmb.every((item) =>
      selectedNoPeserta.value.includes(item.no_peserta),
    ),
);
function toggleSelect(noPeserta) {
  const idx = selectedNoPeserta.value.indexOf(noPeserta);
  if (idx === -1) {
    selectedNoPeserta.value.push(noPeserta);
  } else {
    selectedNoPeserta.value.splice(idx, 1);
  }
}
function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedNoPeserta.value = selectedNoPeserta.value.filter(
      (noPeserta) =>
        !spmbStore.adminSpmb.some((item) => item.no_peserta === noPeserta),
    );
  } else {
    const pageNoPeserta = spmbStore.adminSpmb.map((item) => item.no_peserta);
    selectedNoPeserta.value = [
      ...new Set([...selectedNoPeserta.value, ...pageNoPeserta]),
    ];
  }
}

function openBulkLaporDiri() {
  laporDiriForm.lapor_diri = 1;
  showLaporDiriDialog.value = true;
}

const jalurFilter = ref(spmbStore.filters.jalur);
const laporDiriFilter = ref(spmbStore.filters.lapor_diri);
function changeJalurFilter() {
  spmbStore.filters.jalur = Number(jalurFilter.value);
  spmbStore.pagination.page = 1;
  spmbStore.fetchAdminSpmb();
}
function changeLaporDiriFilter() {
  spmbStore.filters.lapor_diri = Number(laporDiriFilter.value);
  spmbStore.pagination.page = 1;
  spmbStore.fetchAdminSpmb();
}

let searchTimeout = null;
function searchData() {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(async () => {
    spmbStore.pagination.page = 1;
    await spmbStore.fetchAdminSpmb();
  }, 300);
}

async function changePage(page) {
  spmbStore.pagination.page = page;
  await spmbStore.fetchAdminSpmb();
}

function sortIcon(column) {
  if (spmbStore.sortBy !== column) return ChevronUpDownIcon;
  return spmbStore.sortDir === "asc" ? ChevronUpIcon : ChevronDownIcon;
}

// ===== Ubah Lapor Diri =====
const laporDiriForm = reactive({
  lapor_diri: 1,
});

async function handleExport() {
  try {
    await spmbStore.exportSpmb();
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal export data peserta SPMB.",
    );
  }
}
async function handleDownloadTemplate() {
  try {
    await spmbStore.downloadTemplate();
  } catch (error) {
    toast.error("Gagal mengunduh template import.");
  }
}
function handleImportSuccess(message) {
  toast.success(message || "Data peserta SPMB berhasil diimport.");
}

async function handleDelete() {
  loadingAction.value = true;
  try {
    const res = await spmbStore.bulkDelete(selectedNoPeserta.value);
    toast.success(res?.message || "Data peserta SPMB berhasil dihapus.");
    selectedNoPeserta.value = [];
    showDeleteDialog.value = false;
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menghapus data peserta SPMB.",
    );
  } finally {
    loadingAction.value = false;
  }
}

async function handleBulkLaporDiri() {
  loadingAction.value = true;
  try {
    const res = await spmbStore.bulkUpdateLaporDiri(
      selectedNoPeserta.value,
      laporDiriForm.lapor_diri,
    );
    toast.success(res?.message || "Status lapor diri berhasil diperbarui.");
    selectedNoPeserta.value = [];
    showLaporDiriDialog.value = false;
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memperbarui status lapor diri.",
    );
  } finally {
    loadingAction.value = false;
  }
}

onMounted(() => {
  spmbStore.fetchAdminSpmb();
});
</script>
