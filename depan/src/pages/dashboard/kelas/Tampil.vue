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
              <RectangleStackIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Daftar Kelas
            </h2>
          </div>

          <div class="relative w-64 max-w-55 sm:max-w-xs md:max-w-sm">
            <MagnifyingGlassIcon
              class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />

            <input
              v-model="kelasStore.filters.q"
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
              :value="kelasStore.pagination.limit"
              @change="kelasStore.changeLimit($event.target.value)"
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
              :to="{ name: 'dashboard-kelas-tambah' }"
              title="Tambah Kelas"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-blue-500 hover:bg-blue-50 hover:text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-blue-500 dark:hover:bg-blue-500/10 dark:hover:text-blue-400"
            >
              <PlusIcon class="h-5 w-5" />
            </RouterLink>
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div
        v-if="!kelasStore.loading && kelasStore.list.length === 0"
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <MagnifyingGlassIcon v-if="isSearching" class="h-10 w-10" />
          <InboxIcon v-else class="h-10 w-10" />
          <p class="text-sm font-medium">
            {{ isSearching ? "Data tidak ditemukan." : "Belum ada kelas." }}
          </p>
        </div>
      </div>

      <!-- Table -->
      <div v-else class="max-h-[70vh] overflow-auto">
        <table class="w-full min-w-150">
          <thead class="sticky top-0 z-20 bg-gray-50/90 dark:bg-[#0B0F1A]">
            <tr>
              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                ID
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="kelasStore.toggleSort('nama')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Kelas
                  <component :is="sortIcon('nama')" class="h-4 w-4" />
                </button>
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="kelasStore.toggleSort('aktif')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Status
                  <component :is="sortIcon('aktif')" class="h-4 w-4" />
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
            <template v-if="kelasStore.loading">
              <tr v-for="n in 6" :key="`skeleton-${n}`">
                <td v-for="c in 4" :key="c" class="px-4 py-3">
                  <div
                    class="h-4 w-20 animate-pulse rounded bg-gray-100 dark:bg-white/5"
                  ></div>
                </td>
              </tr>
            </template>

            <tr
              v-else
              v-for="item in kelasStore.list"
              :key="item.id_kelas"
              class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
            >
              <td class="px-4 py-2.5 text-sm text-gray-500 dark:text-slate-400">
                {{ item.id_kelas }}
              </td>
              <td
                class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
              >
                <RouterLink
                  :to="{
                    name: 'dashboard-kelas-edit',
                    params: { id: item.id_kelas },
                  }"
                  class="hover:text-[#041983] hover:underline dark:hover:text-indigo-300"
                >
                  {{ item.nama_kelas }}
                </RouterLink>
              </td>

              <td class="px-4 py-2.5">
                <button
                  type="button"
                  title="Ubah Status"
                  @click="confirmToggle(item)"
                  :class="[
                    'inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-1',
                    item.aktif === 1
                      ? 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200 focus-visible:ring-emerald-500/40 dark:bg-emerald-500/10 dark:text-emerald-400 dark:hover:bg-emerald-500/20'
                      : 'bg-gray-100 text-gray-600 hover:bg-gray-200 focus-visible:ring-gray-400/40 dark:bg-slate-800 dark:text-slate-400 dark:hover:bg-slate-700',
                  ]"
                >
                  <CheckCircleIcon
                    v-if="item.aktif === 1"
                    class="h-3.5 w-3.5"
                  />
                  {{ item.aktif === 1 ? "Aktif" : "Nonaktif" }}
                </button>
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
          :pagination="kelasStore.pagination"
          @change-page="changePage"
        />
      </div>
    </div>
  </div>

  <DashboardConfirmDialog
    v-model="showDeleteDialog"
    title="Hapus Kelas"
    :message="`Anda akan menghapus kelas &quot;${itemToDelete?.nama_kelas}&quot;.\n\nData yang telah dihapus tidak dapat dikembalikan.`"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="loadingAction"
    @confirm="handleDelete"
  />

  <DashboardConfirmDialog
    v-model="showToggleDialog"
    title="Ubah Status Kelas"
    :message="`Anda akan meng${itemToToggle?.aktif === 1 ? 'nonaktif' : 'aktif'}kan kelas &quot;${itemToToggle?.nama_kelas}&quot;.`"
    confirm-text="Ya, Lanjutkan"
    confirm-color="blue"
    :loading="loadingAction"
    @confirm="handleToggle"
  />
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";

import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";

import { useToast } from "@/composables/useToast";
import { useKelasStore } from "@/stores/kelas";

import {
  ChevronUpIcon,
  ChevronDownIcon,
  ChevronUpDownIcon,
  CheckCircleIcon,
  MagnifyingGlassIcon,
  PlusIcon,
  TrashIcon,
  RectangleStackIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const kelasStore = useKelasStore();
const toast = useToast();

const showDeleteDialog = ref(false);
const itemToDelete = ref(null);

const showToggleDialog = ref(false);
const itemToToggle = ref(null);

const loadingAction = ref(false);

const isSearching = computed(() => kelasStore.filters.q.trim() !== "");

let searchTimeout = null;

function searchData() {
  clearTimeout(searchTimeout);

  searchTimeout = setTimeout(async () => {
    kelasStore.pagination.page = 1;
    await kelasStore.fetchAll();
  }, 300);
}

async function changePage(page) {
  kelasStore.pagination.page = page;
  await kelasStore.fetchAll();
}

function confirmDelete(item) {
  itemToDelete.value = item;
  showDeleteDialog.value = true;
}

async function handleDelete() {
  if (!itemToDelete.value) return;

  loadingAction.value = true;

  try {
    await kelasStore.remove(itemToDelete.value.id_kelas);
    toast.success("Kelas berhasil dihapus.");
    showDeleteDialog.value = false;
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menghapus kelas.");
  } finally {
    loadingAction.value = false;
    itemToDelete.value = null;
  }
}

function confirmToggle(item) {
  itemToToggle.value = item;
  showToggleDialog.value = true;
}

async function handleToggle() {
  if (!itemToToggle.value) return;

  loadingAction.value = true;

  try {
    const response = await kelasStore.toggleStatus(itemToToggle.value.id_kelas);
    toast.success(response.message || "Status kelas berhasil diubah.");
    showToggleDialog.value = false;
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal mengubah status kelas.",
    );
  } finally {
    loadingAction.value = false;
    itemToToggle.value = null;
  }
}

function sortIcon(column) {
  if (kelasStore.sortBy !== column) return ChevronUpDownIcon;
  return kelasStore.sortDir === "asc" ? ChevronUpIcon : ChevronDownIcon;
}

onMounted(() => kelasStore.fetchAll());
</script>
