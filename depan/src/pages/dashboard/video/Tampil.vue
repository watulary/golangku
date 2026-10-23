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
              <VideoCameraIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Daftar Video
            </h2>
          </div>

          <div class="relative w-64 max-w-55 sm:max-w-xs md:max-w-sm">
            <MagnifyingGlassIcon
              class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />

            <input
              v-model="videoStore.filters.q"
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
              :value="videoStore.pagination.limit"
              @change="videoStore.changeLimit($event.target.value)"
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
              :to="{ name: 'dashboard-video-tambah' }"
              title="Tambah Video"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-blue-500 hover:bg-blue-50 hover:text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-blue-500 dark:hover:bg-blue-500/10 dark:hover:text-blue-400"
            >
              <PlusIcon class="h-5 w-5" />
            </RouterLink>
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div
        v-if="!videoStore.adminLoading && videoStore.adminVideos.length === 0"
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <MagnifyingGlassIcon v-if="isSearching" class="h-10 w-10" />
          <InboxIcon v-else class="h-10 w-10" />
          <p class="text-sm font-medium">
            {{ isSearching ? "Data tidak ditemukan." : "Belum ada video." }}
          </p>
        </div>
      </div>

      <!-- Table -->
      <div v-else class="max-h-[70vh] overflow-auto">
        <table class="min-w-225 w-full">
          <thead class="sticky top-0 z-20 bg-gray-50/90 dark:bg-[#0B0F1A]">
            <tr>
              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Video
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="videoStore.toggleSort('nama')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Judul Video
                  <component :is="sortIcon('nama')" class="h-4 w-4" />
                </button>
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="videoStore.toggleSort('kegiatan')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Kegiatan
                  <component :is="sortIcon('kegiatan')" class="h-4 w-4" />
                </button>
              </th>

              <th
                v-if="isAdmin"
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="videoStore.toggleSort('pengguna')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Pengguna
                  <component :is="sortIcon('pengguna')" class="h-4 w-4" />
                </button>
              </th>

              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                <button
                  type="button"
                  @click="videoStore.toggleSort('waktu')"
                  class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                >
                  Tanggal Dibuat
                  <component :is="sortIcon('waktu')" class="h-4 w-4" />
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
            <template v-if="videoStore.adminLoading">
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
              v-for="item in videoStore.adminVideos"
              :key="item.id"
              class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
            >
              <td class="px-4 py-2.5">
                <button
                  type="button"
                  title="Preview Video"
                  @click="openPreview(item)"
                  class="group relative block h-10 w-14 overflow-hidden rounded-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40"
                >
                  <img
                    v-if="item.thumbnail"
                    :src="item.thumbnail"
                    class="h-10 w-14 rounded-lg object-cover"
                  />
                  <div
                    v-else
                    class="flex h-10 w-14 items-center justify-center rounded-lg bg-gray-100 dark:bg-slate-800"
                  >
                    <VideoCameraIcon class="h-5 w-5 text-gray-300" />
                  </div>

                  <span
                    class="absolute inset-0 flex items-center justify-center bg-black/30 opacity-0 transition group-hover:opacity-100"
                  >
                    <PlayCircleIcon class="h-6 w-6 text-white" />
                  </span>
                </button>
              </td>
              <td
                class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
              >
                <RouterLink
                  :to="{
                    name: 'dashboard-video-edit',
                    params: { id: item.id },
                  }"
                  class="hover:text-[#041983] hover:underline dark:hover:text-indigo-300"
                >
                  {{ item.nama }}
                </RouterLink>
              </td>

              <td class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400">
                {{ item.kegiatan }}
              </td>

              <td
                v-if="isAdmin"
                class="px-4 py-2.5 text-xs text-gray-600 dark:text-slate-400"
              >
                {{ item.nama_pengguna || "—" }}
              </td>

              <td class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400">
                {{ formatDateTimeNumeric(item.waktu) }}
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
          :pagination="videoStore.pagination"
          @change-page="changePage"
        />
      </div>
    </div>
  </div>

  <DashboardConfirmDialog
    v-model="showDeleteDialog"
    title="Hapus Video"
    :message="`Anda akan menghapus video &quot;${videoToDelete?.nama}&quot;.\n\nData yang telah dihapus tidak dapat dikembalikan.`"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="loadingAction"
    @confirm="handleDelete"
  />

  <!-- Preview Modal -->
  <Teleport to="body">
    <div
      v-if="showPreview"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4"
      @click.self="closePreview"
    >
      <div class="w-full max-w-2xl">
        <div class="mb-2 flex items-center justify-between">
          <h3 class="text-sm font-medium text-white">
            {{ previewItem?.nama }}
          </h3>

          <button
            type="button"
            title="Tutup"
            @click="closePreview"
            class="rounded-lg p-1.5 text-white/80 transition hover:bg-white/10 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40"
          >
            <XMarkIcon class="h-5 w-5" />
          </button>
        </div>

        <div class="aspect-video w-full overflow-hidden rounded-xl bg-black">
          <iframe
            v-if="previewURL"
            :src="previewURL"
            class="h-full w-full"
            allow="
              accelerometer;
              autoplay;
              clipboard-write;
              encrypted-media;
              gyroscope;
              picture-in-picture;
            "
            allowfullscreen
          />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";

import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";

import { useToast } from "@/composables/useToast";
import { useVideoStore } from "@/stores/video";
import { useAuthStore } from "@/stores/auth";
import { formatDateTimeNumeric } from "@/utils/date";

import {
  ChevronUpIcon,
  ChevronDownIcon,
  ChevronUpDownIcon,
  MagnifyingGlassIcon,
  PlayCircleIcon,
  PlusIcon,
  TrashIcon,
  VideoCameraIcon,
  XMarkIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const videoStore = useVideoStore();
const toast = useToast();
const authStore = useAuthStore();

const showDeleteDialog = ref(false);
const loadingAction = ref(false);
const videoToDelete = ref(null);

const showPreview = ref(false);
const previewItem = ref(null);
const previewURL = ref("");

const isAdmin = computed(() => authStore.role === "Admin");

const isSearching = computed(() => videoStore.filters.q.trim() !== "");

let searchTimeout = null;

function searchData() {
  clearTimeout(searchTimeout);

  searchTimeout = setTimeout(async () => {
    videoStore.pagination.page = 1;
    await videoStore.fetchAdminVideos();
  }, 300);
}

async function changePage(page) {
  videoStore.pagination.page = page;
  await videoStore.fetchAdminVideos();
}

function confirmDelete(item) {
  videoToDelete.value = item;
  showDeleteDialog.value = true;
}

async function handleDelete() {
  if (!videoToDelete.value) return;

  loadingAction.value = true;

  try {
    await videoStore.deleteVideo(videoToDelete.value.id);
    toast.success("Video berhasil dihapus.");
    showDeleteDialog.value = false;
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menghapus video.");
  } finally {
    loadingAction.value = false;
    videoToDelete.value = null;
  }
}

function openPreview(item) {
  previewItem.value = item;
  previewURL.value = item.youtube_url || "";
  showPreview.value = true;
}

function closePreview() {
  showPreview.value = false;
  previewURL.value = "";
  previewItem.value = null;
}

function sortIcon(column) {
  if (videoStore.sortBy !== column) return ChevronUpDownIcon;
  return videoStore.sortDir === "asc" ? ChevronUpIcon : ChevronDownIcon;
}

onMounted(() => videoStore.fetchAdminVideos());
</script>
