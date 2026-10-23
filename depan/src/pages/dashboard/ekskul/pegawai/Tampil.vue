<template>
  <DashboardPageHeader
    title="Ekstrakurikuler"
    subtitle="Daftar ekstrakurikuler beserta pembina dan jumlah anggotanya."
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
        <!-- Baris 1 -->
        <div
          class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/8 text-[#041983] dark:bg-white/5 dark:text-indigo-300"
            >
              <TrophyIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Daftar Ekstrakurikuler
            </h2>
          </div>

          <div class="flex flex-col gap-3 sm:flex-row">
            <select
              :value="ekskulStore.filters.aktif"
              @change="ekskulStore.changeFilterAktif($event.target.value)"
              class="rounded-xl border border-slate-300 bg-white px-3 py-2.5 text-sm text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
            >
              <option value="">Semua Status</option>
              <option value="1">Aktif</option>
              <option value="0">Nonaktif</option>
            </select>

            <div class="relative w-64 max-w-55 sm:max-w-xs md:max-w-sm">
              <MagnifyingGlassIcon
                class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
              />

              <input
                v-model="ekskulStore.filters.q"
                @input="searchData"
                type="text"
                placeholder="Pencarian..."
                class="w-full rounded-xl border border-gray-300 bg-white py-2.5 pl-10 pr-4 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500"
              />
            </div>
          </div>
        </div>

        <!-- Baris 2 -->
        <div class="mt-5 flex items-center gap-2">
          <span class="text-sm text-slate-500 dark:text-slate-400">
            Tampilkan
          </span>

          <select
            :value="ekskulStore.pagination.limit"
            @change="ekskulStore.changeLimit($event.target.value)"
            class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
          >
            <option :value="10">10</option>
            <option :value="25">25</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
          </select>
        </div>
      </div>

      <!-- Empty -->
      <div
        v-if="!ekskulStore.loading && ekskulStore.list.length === 0"
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <MagnifyingGlassIcon v-if="isSearching" class="h-10 w-10" />
          <InboxIcon v-else class="h-10 w-10" />
          <p class="text-sm font-medium">
            {{ isSearching ? "Data tidak ditemukan." : "Belum ada ekskul." }}
          </p>
        </div>
      </div>

      <!-- Table -->
      <div v-else class="max-h-[70vh] overflow-auto">
        <table class="w-full min-w-200">
          <thead class="sticky top-0 z-20 bg-gray-50/90 dark:bg-[#0B0F1A]">
            <tr>
              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="ekskulStore.toggleSort('ekskul')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Nama Ekstrakurikuler
                  <component :is="sortIcon('ekskul')" class="h-4 w-4" />
                </button>
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="ekskulStore.toggleSort('pegawai')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Pembina
                  <component :is="sortIcon('pegawai')" class="h-4 w-4" />
                </button>
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                NUPTK
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Status
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-center text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Anggota
              </th>
            </tr>
          </thead>

          <tbody class="divide-y divide-gray-100 dark:divide-white/5">
            <!-- Skeleton saat memuat -->
            <template v-if="ekskulStore.loading">
              <tr v-for="n in 6" :key="`skeleton-${n}`">
                <td v-for="c in 5" :key="c" class="px-4 py-3">
                  <div
                    class="h-4 w-20 animate-pulse rounded bg-gray-100 dark:bg-white/5"
                  ></div>
                </td>
              </tr>
            </template>

            <tr
              v-else
              v-for="item in ekskulStore.list"
              :key="item.id_ekskul"
              class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
            >
              <td
                class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
              >
                {{ item.nama_ekskul }}
              </td>

              <td class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400">
                {{ item.nama_pegawai || "—" }}
              </td>

              <td class="px-4 py-2.5 text-xs text-gray-500 dark:text-slate-400">
                {{ item.nuptk }}
              </td>

              <td class="px-4 py-2.5">
                <span
                  :class="[
                    'inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium',
                    item.aktif
                      ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400'
                      : 'bg-gray-100 text-gray-600 dark:bg-slate-800 dark:text-slate-400',
                  ]"
                >
                  {{ item.aktif ? "Aktif" : "Nonaktif" }}
                </span>
              </td>

              <td class="px-4 py-2.5 text-center">
                <button
                  type="button"
                  @click="openAnggota(item)"
                  title="Lihat Anggota"
                  class="relative inline-flex rounded-lg p-2 text-gray-500 transition hover:bg-blue-50 hover:text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:text-slate-400 dark:hover:bg-blue-500/10 dark:hover:text-blue-400"
                >
                  <UserGroupIcon class="h-4 w-4" />
                  <span
                    v-if="item.jumlah_anggota > 0"
                    class="absolute -right-1.5 -top-1.5 inline-flex h-4 min-w-4 items-center justify-center rounded-full bg-blue-600 px-1 text-[10px] font-semibold leading-none text-white dark:bg-blue-500"
                  >
                    {{ item.jumlah_anggota }}
                  </span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <DashboardPagination
          :pagination="ekskulStore.pagination"
          @change-page="changePage"
        />
      </div>
    </div>
  </div>

  <!-- Modal Daftar Anggota (read-only) -->
  <Teleport to="body">
    <div
      v-if="showAnggotaModal"
      @click.self="closeAnggotaModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4"
    >
      <div
        class="flex max-h-[85vh] w-full max-w-2xl flex-col overflow-hidden rounded-2xl bg-white shadow-xl ring-1 ring-slate-900/5 dark:bg-slate-900 dark:ring-white/10"
      >
        <!-- Header Modal -->
        <div
          class="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-white/5"
        >
          <div>
            <h3 class="text-lg font-semibold text-gray-800 dark:text-slate-200">
              Anggota Ekstrakurikuler
            </h3>
            <p class="text-sm text-gray-500 dark:text-slate-400">
              {{ selectedEkskul?.nama_ekskul }}
            </p>
          </div>

          <button
            type="button"
            @click="closeAnggotaModal"
            class="rounded-lg p-2 text-gray-500 transition hover:bg-gray-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:text-slate-400 dark:hover:bg-slate-800"
          >
            <XMarkIcon class="h-5 w-5" />
          </button>
        </div>

        <!-- Search -->
        <div class="border-b border-gray-100 px-6 py-4 dark:border-white/5">
          <div class="relative">
            <MagnifyingGlassIcon
              class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />

            <input
              v-model="ekskulStore.anggotaFilters.q"
              @input="searchAnggota"
              type="text"
              placeholder="Pencarian..."
              class="w-full rounded-xl border border-gray-300 bg-white py-2.5 pl-10 pr-4 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500"
            />
          </div>
        </div>

        <!-- Isi -->
        <div class="flex-1 overflow-y-auto">
          <!-- Loading -->
          <div
            v-if="ekskulStore.anggotaLoading"
            class="flex flex-col items-center justify-center py-16"
          >
            <div
              class="h-8 w-8 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"
            />
            <p class="mt-3 text-sm text-gray-500 dark:text-slate-400">
              Memuat data...
            </p>
          </div>

          <!-- Empty -->
          <div
            v-else-if="ekskulStore.anggotaList.length === 0"
            class="py-16 text-center"
          >
            <div
              class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
            >
              <InboxIcon class="h-10 w-10" />
              <p class="text-sm font-medium">
                {{
                  isSearchingAnggota
                    ? "Data tidak ditemukan."
                    : "Belum ada anggota terdaftar."
                }}
              </p>
            </div>
          </div>

          <!-- Table -->
          <table v-else class="w-full">
            <thead class="sticky top-0 bg-gray-50 dark:bg-[#0B0F1A]">
              <tr>
                <th
                  scope="col"
                  class="px-6 py-2.5 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Rombel
                </th>
                <th
                  scope="col"
                  class="px-4 py-2.5 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  NIS
                </th>
                <th
                  scope="col"
                  class="px-4 py-2.5 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Nama Peserta Didik
                </th>
              </tr>
            </thead>

            <tbody class="divide-y divide-gray-100 dark:divide-white/5">
              <tr v-for="item in ekskulStore.anggotaList" :key="item.id">
                <td
                  class="px-6 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.nama_rombel || "—" }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.nis }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
                >
                  {{ item.nama_siswa }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div class="border-t border-gray-100 px-6 py-3 dark:border-white/5">
          <DashboardPagination
            :pagination="ekskulStore.anggotaPagination"
            @change-page="changeAnggotaPage"
          />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";

import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";

import { useEkskulStore } from "@/stores/ekskul";

import {
  ChevronUpIcon,
  ChevronDownIcon,
  ChevronUpDownIcon,
  MagnifyingGlassIcon,
  UserGroupIcon,
  XMarkIcon,
  TrophyIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const ekskulStore = useEkskulStore();

const isSearching = computed(() => ekskulStore.filters.q.trim() !== "");

let searchTimeout = null;

function searchData() {
  clearTimeout(searchTimeout);

  searchTimeout = setTimeout(async () => {
    ekskulStore.pagination.page = 1;
    await ekskulStore.fetchAll();
  }, 300);
}

async function changePage(page) {
  ekskulStore.pagination.page = page;
  await ekskulStore.fetchAll();
}

function sortIcon(column) {
  if (ekskulStore.sortBy !== column) return ChevronUpDownIcon;
  return ekskulStore.sortDir === "asc" ? ChevronUpIcon : ChevronDownIcon;
}

// =====================================
// Modal Daftar Anggota (read-only)
// =====================================

const showAnggotaModal = ref(false);
const selectedEkskul = ref(null);

const isSearchingAnggota = computed(
  () => ekskulStore.anggotaFilters.q.trim() !== "",
);

let anggotaSearchTimeout = null;

async function openAnggota(item) {
  selectedEkskul.value = item;
  ekskulStore.anggotaFilters.q = "";
  ekskulStore.anggotaPagination.page = 1;
  showAnggotaModal.value = true;

  await ekskulStore.fetchAnggota(item.id_ekskul);
}

function closeAnggotaModal() {
  showAnggotaModal.value = false;
  selectedEkskul.value = null;
}

function searchAnggota() {
  clearTimeout(anggotaSearchTimeout);

  anggotaSearchTimeout = setTimeout(async () => {
    ekskulStore.anggotaPagination.page = 1;
    await ekskulStore.fetchAnggota(selectedEkskul.value.id_ekskul);
  }, 300);
}

async function changeAnggotaPage(page) {
  ekskulStore.anggotaPagination.page = page;
  await ekskulStore.fetchAnggota(selectedEkskul.value.id_ekskul);
}

onMounted(() => {
  ekskulStore.fetchAll();
});
</script>
