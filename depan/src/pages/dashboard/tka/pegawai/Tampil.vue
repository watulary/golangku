<template>
  <DashboardPageHeader />

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
        <div class="mt-5 flex items-center gap-2">
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
                    Nama Siswa
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
                  <td v-for="c in 6" :key="c" class="px-4 py-3">
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
</template>

<script setup>
import { computed, onMounted } from "vue";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import { useTkaStore } from "@/stores/tka";
import {
  ChevronDownIcon,
  ChevronUpIcon,
  ChevronUpDownIcon,
  MagnifyingGlassIcon,
  ClipboardDocumentCheckIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const tkaStore = useTkaStore();

const isSearching = computed(() => tkaStore.filters.q.trim() !== "");

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

onMounted(() => {
  tkaStore.fetchAdminTka();
});
</script>
