<template>
  <DashboardPageHeader
    title="Guru Wali"
    subtitle="Kelola data guru wali dan peserta didik yang menjadi anggotanya."
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
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-blue-100 text-blue-600 dark:bg-blue-500/10 dark:text-blue-300"
            >
              <AcademicCapIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Daftar Guru Wali
            </h2>
          </div>

          <div class="flex flex-col gap-3 sm:flex-row">
            <select
              :value="guruwaliStore.filters.aktif"
              @change="guruwaliStore.changeFilterAktif($event.target.value)"
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
                v-model="guruwaliStore.filters.q"
                @input="searchData"
                type="text"
                placeholder="Pencarian..."
                class="w-full rounded-xl border border-gray-300 bg-white py-2.5 pl-10 pr-4 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-100 dark:placeholder:text-slate-500"
              />
            </div>
          </div>
        </div>

        <!-- Baris 2 -->
        <div
          class="mt-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
        >
          <!-- kiri -->
          <div class="flex items-center gap-2">
            <span class="text-sm text-slate-500 dark:text-slate-400">
              Tampilkan
            </span>

            <select
              :value="guruwaliStore.pagination.limit"
              @change="guruwaliStore.changeLimit($event.target.value)"
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

          <!-- kanan: TOMBOL TAMBAH ada di sini -->
          <div class="flex items-center gap-2">
            <RouterLink
              :to="{ name: 'dashboard-guruwali-tambah' }"
              title="Tambah Guru Wali"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-blue-500 hover:bg-blue-50 hover:text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-blue-500 dark:hover:bg-blue-500/10 dark:hover:text-blue-400"
            >
              <PlusIcon class="h-5 w-5" />
            </RouterLink>

            <!-- Export Daftar Guru Wali -->
            <button
              type="button"
              @click="exportGuruWali"
              title="Unduh Daftar Guru Wali"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-emerald-500 dark:hover:bg-emerald-500/10 dark:hover:text-emerald-400"
            >
              <ArrowDownTrayIcon class="h-5 w-5" />
            </button>

            <!-- Export Daftar Anggota (semua guru wali) -->
            <button
              type="button"
              @click="exportAnggota"
              title="Unduh Daftar Anggota Semua Guru Wali"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-emerald-500 dark:hover:bg-emerald-500/10 dark:hover:text-emerald-400"
            >
              <DocumentArrowDownIcon class="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div
        v-if="!guruwaliStore.loading && guruwaliStore.list.length === 0"
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <MagnifyingGlassIcon v-if="isSearching" class="h-10 w-10" />
          <InboxIcon v-else class="h-10 w-10" />
          <p class="text-sm font-medium">
            {{ isSearching ? "Data tidak ditemukan." : "Belum ada guru wali." }}
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
                  @click="guruwaliStore.toggleSort('nama')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Nama Wali
                  <component :is="sortIcon('nama')" class="h-4 w-4" />
                </button>
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="guruwaliStore.toggleSort('pegawai')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Guru
                  <component :is="sortIcon('pegawai')" class="h-4 w-4" />
                </button>
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="guruwaliStore.toggleSort('nuptk')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  NUPTK
                  <component :is="sortIcon('nuptk')" class="h-4 w-4" />
                </button>
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="guruwaliStore.toggleSort('aktif')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Status
                  <component :is="sortIcon('aktif')" class="h-4 w-4" />
                </button>
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="guruwaliStore.toggleSort('anggota')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Anggota
                  <component :is="sortIcon('anggota')" class="h-4 w-4" />
                </button>
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-right text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Aksi
              </th>
            </tr>
          </thead>

          <tbody class="divide-y divide-gray-100 dark:divide-white/5">
            <!-- Skeleton saat memuat -->
            <template v-if="guruwaliStore.loading">
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
              v-for="item in guruwaliStore.list"
              :key="item.id_guruwali"
              class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
            >
              <td
                class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
              >
                <RouterLink
                  :to="{
                    name: 'dashboard-guruwali-edit',
                    params: { id: item.id_guruwali },
                  }"
                  class="hover:text-[#041983] hover:underline dark:hover:text-indigo-300"
                >
                  {{ item.nama_guruwali }}
                </RouterLink>
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

              <td class="px-4 py-2.5">
                <RouterLink
                  :to="{
                    name: 'dashboard-guruwali-anggota',
                    params: { id: item.id_guruwali },
                  }"
                  title="Kelola Anggota"
                  class="relative inline-flex rounded-lg p-2 text-gray-500 transition hover:bg-blue-50 hover:text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:text-slate-400 dark:hover:bg-blue-500/10 dark:hover:text-blue-400"
                >
                  <UserGroupIcon class="h-4 w-4" />
                  <span
                    v-if="item.jumlah_anggota > 0"
                    class="absolute -right-1.5 -top-1.5 inline-flex h-4 min-w-4 items-center justify-center rounded-full bg-blue-600 px-1 text-[10px] font-semibold leading-none text-white dark:bg-blue-500"
                  >
                    {{ item.jumlah_anggota }}
                  </span>
                </RouterLink>
              </td>

              <td class="px-4 py-2.5">
                <div class="flex items-center justify-end gap-1">
                  <button
                    type="button"
                    title="Hapus"
                    @click="confirmDelete(item)"
                    class="rounded-lg p-2 text-gray-500 transition hover:bg-red-50 hover:text-red-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40 dark:text-slate-400 dark:hover:bg-red-500/10"
                  >
                    <TrashIcon class="h-4 w-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <DashboardPagination
          :pagination="guruwaliStore.pagination"
          @change-page="changePage"
        />
      </div>
    </div>
  </div>

  <DashboardConfirmDialog
    v-model="showDeleteDialog"
    title="Hapus Guru Wali"
    :message="`Anda akan menghapus guru wali &quot;${itemToDelete?.nama_guruwali}&quot;.\n\nSeluruh data anggota di dalamnya juga akan ikut terhapus dan tidak dapat dikembalikan.`"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="loadingAction"
    @confirm="handleDelete"
  />
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";

import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";

import { useToast } from "@/composables/useToast";
import { useGuruWaliStore } from "@/stores/guruwali";

import {
  ChevronUpIcon,
  ChevronDownIcon,
  ChevronUpDownIcon,
  MagnifyingGlassIcon,
  PlusIcon,
  TrashIcon,
  UserGroupIcon,
  AcademicCapIcon,
  InboxIcon,
  ArrowDownTrayIcon,
  DocumentArrowDownIcon,
} from "@heroicons/vue/24/outline";

const guruwaliStore = useGuruWaliStore();
const toast = useToast();

const showDeleteDialog = ref(false);
const itemToDelete = ref(null);
const loadingAction = ref(false);

const isSearching = computed(() => guruwaliStore.filters.q.trim() !== "");

let searchTimeout = null;

function searchData() {
  clearTimeout(searchTimeout);

  searchTimeout = setTimeout(async () => {
    guruwaliStore.pagination.page = 1;
    await guruwaliStore.fetchAll();
  }, 300);
}

async function changePage(page) {
  guruwaliStore.pagination.page = page;
  await guruwaliStore.fetchAll();
}

async function exportGuruWali() {
  try {
    await guruwaliStore.exportGuruWali();
    toast.success("Daftar guru wali berhasil diekspor.");
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal mengekspor daftar guru wali.",
    );
  }
}

async function exportAnggota() {
  try {
    await guruwaliStore.exportAnggota();
    toast.success("Daftar anggota guru wali berhasil diekspor.");
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal mengekspor daftar anggota.",
    );
  }
}

function confirmDelete(item) {
  itemToDelete.value = item;
  showDeleteDialog.value = true;
}

async function handleDelete() {
  if (!itemToDelete.value) return;

  loadingAction.value = true;

  try {
    await guruwaliStore.remove(itemToDelete.value.id_guruwali);
    toast.success("Guru wali berhasil dihapus.");
    showDeleteDialog.value = false;
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menghapus guru wali.");
  } finally {
    loadingAction.value = false;
    itemToDelete.value = null;
  }
}

function sortIcon(column) {
  if (guruwaliStore.sortBy !== column) return ChevronUpDownIcon;
  return guruwaliStore.sortDir === "asc" ? ChevronUpIcon : ChevronDownIcon;
}

onMounted(() => {
  guruwaliStore.fetchAll();
});
</script>
