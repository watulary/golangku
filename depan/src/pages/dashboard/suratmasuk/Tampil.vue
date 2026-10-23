<template>
  <DashboardPageHeader />

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
              <InboxArrowDownIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Daftar Surat Masuk
            </h2>
          </div>

          <div class="relative w-64 max-w-55 sm:max-w-xs md:max-w-sm">
            <MagnifyingGlassIcon
              class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />

            <input
              v-model="suratMasukStore.filters.q"
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
          <div class="flex items-center gap-2">
            <span class="text-sm text-slate-500 dark:text-slate-400">
              Tampilkan
            </span>

            <select
              :value="suratMasukStore.pagination.limit"
              @change="suratMasukStore.changeLimit($event.target.value)"
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

          <!-- kanan -->
          <div class="flex items-center gap-2">
            <RouterLink
              :to="{ name: 'dashboard-suratmasuk-tambah' }"
              title="Tambah Surat Masuk"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-blue-500 hover:bg-blue-50 hover:text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-blue-500 dark:hover:bg-blue-500/10 dark:hover:text-blue-400"
            >
              <PlusIcon class="h-5 w-5" />
            </RouterLink>
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div
        v-if="
          !suratMasukStore.adminLoading &&
          suratMasukStore.adminSuratMasuk.length === 0
        "
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <MagnifyingGlassIcon v-if="isSearching" class="h-10 w-10" />
          <InboxIcon v-else class="h-10 w-10" />
          <p class="text-sm font-medium">
            {{
              isSearching ? "Data tidak ditemukan." : "Belum ada surat masuk."
            }}
          </p>
        </div>
      </div>

      <!-- Table -->
      <div v-else class="max-h-[70vh] overflow-auto">
        <table class="min-w-275 w-full">
          <thead class="sticky top-0 z-20 bg-gray-50/90 dark:bg-[#0B0F1A]">
            <tr>
              <th
                scope="col"
                class="px-4 py-3 text-[11px] font-semibold text-gray-500 dark:text-slate-400"
              >
                <div class="flex items-center justify-center">
                  <ArrowDownTrayIcon class="h-4 w-4" />
                </div>
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="suratMasukStore.toggleSort('nama')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Nama Surat
                  <component :is="sortIcon('nama')" class="h-4 w-4" />
                </button>
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="suratMasukStore.toggleSort('nomor')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Nomor Surat
                  <component :is="sortIcon('nomor')" class="h-4 w-4" />
                </button>
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="suratMasukStore.toggleSort('asal')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Asal Surat
                  <component :is="sortIcon('asal')" class="h-4 w-4" />
                </button>
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="suratMasukStore.toggleSort('tujuan')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Tujuan
                  <component :is="sortIcon('tujuan')" class="h-4 w-4" />
                </button>
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="suratMasukStore.toggleSort('tanggal')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Tanggal Surat
                  <component :is="sortIcon('tanggal')" class="h-4 w-4" />
                </button>
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="suratMasukStore.toggleSort('tanggal_terima')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Tanggal Diterima
                  <component :is="sortIcon('tanggal_terima')" class="h-4 w-4" />
                </button>
              </th>

              <th
                v-if="isAdmin"
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Pengguna
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
            <template v-if="suratMasukStore.adminLoading">
              <tr v-for="n in 6" :key="`skeleton-${n}`">
                <td v-for="c in 9" :key="c" class="px-4 py-3">
                  <div
                    class="h-4 w-20 animate-pulse rounded bg-gray-100 dark:bg-white/5"
                  ></div>
                </td>
              </tr>
            </template>

            <tr
              v-else
              v-for="item in suratMasukStore.adminSuratMasuk"
              :key="item.id"
              class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
            >
              <td class="px-4 py-2.5">
                <a
                  v-if="item.dokumen_url"
                  :href="item.dokumen_url"
                  rel="noopener"
                  title="Unduh"
                  class="flex h-10 w-10 items-center justify-center rounded-lg bg-gray-100 text-gray-600 transition hover:bg-gray-200 hover:text-gray-800 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700 dark:hover:text-white"
                >
                  <DocumentArrowDownIcon class="h-5 w-5" />
                </a>
                <div
                  v-else
                  class="flex h-10 w-10 items-center justify-center rounded-lg bg-gray-100 dark:bg-slate-800"
                >
                  <DocumentTextIcon class="h-5 w-5 text-gray-300" />
                </div>
              </td>

              <td
                class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
              >
                <RouterLink
                  :to="{
                    name: 'dashboard-suratmasuk-edit',
                    params: { id: item.id },
                  }"
                  class="hover:text-[#041983] hover:underline dark:hover:text-indigo-300"
                >
                  {{ item.nama }}
                </RouterLink>
              </td>

              <td class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400">
                {{ item.nomor || "—" }}
              </td>

              <td class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400">
                {{ item.asal || "—" }}
              </td>

              <td class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400">
                {{ item.tujuan || "—" }}
              </td>

              <td class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400">
                {{ formatDateNumeric(item.tanggal) }}
              </td>

              <td class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400">
                {{ formatDateNumeric(item.tanggal_terima) }}
              </td>

              <td
                v-if="isAdmin"
                class="px-4 py-2.5 text-xs text-gray-600 dark:text-slate-400"
              >
                {{ item.nama_pengguna || "—" }}
              </td>

              <td class="px-4 py-2.5">
                <div class="flex items-center justify-end">
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
          :pagination="suratMasukStore.pagination"
          @change-page="changePage"
        />
      </div>
    </div>
  </div>
  <DashboardConfirmDialog
    v-model="showDeleteDialog"
    title="Hapus Surat Masuk"
    :message="`Anda akan menghapus surat masuk &quot;${suratMasukToDelete?.nama}&quot;.\n\nData yang telah dihapus tidak dapat dikembalikan.`"
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
import { useSuratMasukStore } from "@/stores/suratmasuk";
import { useAuthStore } from "@/stores/auth";
import { formatDateNumeric } from "@/utils/date";

import {
  ChevronUpIcon,
  ChevronDownIcon,
  ChevronUpDownIcon,
  DocumentArrowDownIcon,
  DocumentTextIcon,
  MagnifyingGlassIcon,
  PlusIcon,
  TrashIcon,
  ArrowDownTrayIcon,
  InboxArrowDownIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const suratMasukStore = useSuratMasukStore();
const toast = useToast();
const authStore = useAuthStore();

const showDeleteDialog = ref(false);
const loadingAction = ref(false);
const suratMasukToDelete = ref(null);

const isAdmin = computed(() => authStore.role === "Admin");

const isSearching = computed(() => suratMasukStore.filters.q.trim() !== "");

let searchTimeout = null;

function searchData() {
  clearTimeout(searchTimeout);

  searchTimeout = setTimeout(async () => {
    suratMasukStore.pagination.page = 1;
    await suratMasukStore.fetchAdminSuratMasuk();
  }, 300);
}

async function changePage(page) {
  suratMasukStore.pagination.page = page;
  await suratMasukStore.fetchAdminSuratMasuk();
}

function confirmDelete(item) {
  suratMasukToDelete.value = item;
  showDeleteDialog.value = true;
}

async function handleDelete() {
  if (!suratMasukToDelete.value) return;

  loadingAction.value = true;

  try {
    await suratMasukStore.deleteSuratMasuk(suratMasukToDelete.value.id);
    toast.success("Surat masuk berhasil dihapus.");
    showDeleteDialog.value = false;
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menghapus surat masuk.",
    );
  } finally {
    loadingAction.value = false;
    suratMasukToDelete.value = null;
  }
}

function sortIcon(column) {
  if (suratMasukStore.sortBy !== column) return ChevronUpDownIcon;
  return suratMasukStore.sortDir === "asc" ? ChevronUpIcon : ChevronDownIcon;
}

onMounted(() => suratMasukStore.fetchAdminSuratMasuk());
</script>
