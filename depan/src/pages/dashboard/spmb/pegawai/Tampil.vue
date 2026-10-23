<template>
  <DashboardPageHeader
    title="Data SPMB"
    subtitle="Daftar peserta didik yang sudah diterima di sekolah (lihat saja)."
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

        <div class="mt-5 flex flex-wrap items-center gap-3">
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
          <table class="min-w-200 w-full">
            <thead class="sticky top-0 z-20 bg-gray-50/90 dark:bg-[#0B0F1A]">
              <tr>
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
                  <td v-for="c in 6" :key="c" class="px-4 py-3">
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
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.no_peserta }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
                >
                  {{ item.nama_peserta }}
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
</template>

<script setup>
import { computed, onMounted } from "vue";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import { useSpmbStore } from "@/stores/spmb";
import {
  ChevronDownIcon,
  ChevronUpIcon,
  ChevronUpDownIcon,
  MagnifyingGlassIcon,
  UserGroupIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const spmbStore = useSpmbStore();

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
      "bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400",
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

const jalurFilter = computed({
  get: () => spmbStore.filters.jalur,
  set: (val) => (spmbStore.filters.jalur = val),
});
const laporDiriFilter = computed({
  get: () => spmbStore.filters.lapor_diri,
  set: (val) => (spmbStore.filters.lapor_diri = val),
});

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

onMounted(() => {
  spmbStore.fetchAdminSpmb();
});
</script>
