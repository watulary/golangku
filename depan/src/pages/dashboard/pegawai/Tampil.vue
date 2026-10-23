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
              <AcademicCapIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Daftar Pegawai
            </h2>
          </div>

          <div class="relative w-64 max-w-55 sm:max-w-xs md:max-w-sm">
            <MagnifyingGlassIcon
              class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />

            <input
              v-model="pegawaiStore.filters.q"
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
                :value="pegawaiStore.pagination.limit"
                @change="pegawaiStore.changeLimit($event.target.value)"
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
              v-model="statusFilter"
              @change="changeStatusFilter"
              class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
            >
              <option :value="-1">Semua Status</option>
              <option :value="1">Aktif</option>
              <option :value="0">Keluar</option>
            </select>

            <select
              v-model="kategoriFilter"
              @change="changeKategoriFilter"
              class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
            >
              <option value="">Semua Jabatan</option>
              <option value="Guru">Guru</option>
              <option value="Tatausaha">Tendik</option>
            </select>

            <select
              v-model="pnsNonFilter"
              @change="changePnsNonFilter"
              class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
            >
              <option value="">Semua Pegawai</option>
              <option value="PNS">PNS</option>
              <option value="CPNS">CPNS</option>
              <option value="PPPK">PPPK</option>
              <option value="PPPK-PW">PPPK-PW</option>
              <option value="KKI">KKI</option>
              <option value="Non KKI">Non KKI</option>
              <option value="Guru Honor Sekolah">Guru Honor Sekolah</option>
              <option value="Honor Daerah TK.I Provinsi">
                Honor Daerah TK.I Provinsi
              </option>
            </select>
          </div>

          <!-- kanan -->
          <div class="flex items-center gap-2">
            <RouterLink
              :to="{ name: 'dashboard-pegawai-tambah' }"
              title="Tambah Pegawai"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-blue-500 hover:bg-blue-50 hover:text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-blue-500 dark:hover:bg-blue-500/10 dark:hover:text-blue-400"
            >
              <PlusIcon class="h-5 w-5" />
            </RouterLink>
            <!-- ganti tombol Import lama dengan ini -->
            <button
              type="button"
              title="Unggah Pegawai"
              @click="showUploadDialog = true"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-yellow-500 hover:bg-yellow-50 hover:text-yellow-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-yellow-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-yellow-500 dark:hover:bg-yellow-500/10 dark:hover:text-yellow-300"
            >
              <ArrowUpTrayIcon class="h-5 w-5" />
            </button>
            <button
              type="button"
              title="Unduh Pegawai"
              @click="handleExport"
              :disabled="pegawaiStore.exportLoading"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-emerald-500 dark:hover:bg-emerald-500/10 dark:hover:text-emerald-400"
            >
              <ArrowDownTrayIcon class="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div
        v-if="
          !pegawaiStore.adminLoading && pegawaiStore.adminPegawai.length === 0
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
              isSearching ? "Data tidak ditemukan." : "Belum ada data pegawai."
            }}
          </p>
        </div>
      </div>

      <!-- Table -->
      <div v-else>
        <div class="max-h-[70vh] overflow-auto">
          <table class="min-w-250 w-full">
            <thead class="sticky top-0 z-20 bg-gray-50/90 dark:bg-[#0B0F1A]">
              <tr>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Foto
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="pegawaiStore.toggleSort('nip_baru')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    NIP
                    <component :is="sortIcon('nip_baru')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="pegawaiStore.toggleSort('nama')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Nama
                    <component :is="sortIcon('nama')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Pegawai
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Status
                </th>

                <!-- Kolom kondisional: Keluar vs Aktif/Semua -->
                <template v-if="isKeluarView">
                  <th
                    scope="col"
                    class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                  >
                    <button
                      type="button"
                      @click="pegawaiStore.toggleSort('tanggal_keluar')"
                      class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                    >
                      Tanggal Keluar
                      <component
                        :is="sortIcon('tanggal_keluar')"
                        class="h-4 w-4"
                      />
                    </button>
                  </th>
                  <th
                    scope="col"
                    class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                  >
                    Nomor Surat
                  </th>
                  <th
                    scope="col"
                    class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                  >
                    Keterangan
                  </th>
                </template>
                <template v-else>
                  <th
                    scope="col"
                    class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                  >
                    <button
                      type="button"
                      @click="pegawaiStore.toggleSort('jenis_kelamin')"
                      class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                    >
                      L/P
                      <component
                        :is="sortIcon('jenis_kelamin')"
                        class="h-4 w-4"
                      />
                    </button>
                  </th>
                  <th
                    scope="col"
                    class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                  >
                    <button
                      type="button"
                      @click="pegawaiStore.toggleSort('agama')"
                      class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                    >
                      Agama
                      <component :is="sortIcon('agama')" class="h-4 w-4" />
                    </button>
                  </th>
                  <th
                    scope="col"
                    class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                  >
                    <button
                      type="button"
                      @click="pegawaiStore.toggleSort('tempat_lahir')"
                      class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                    >
                      Tempat Lahir
                      <component
                        :is="sortIcon('tempat_lahir')"
                        class="h-4 w-4"
                      />
                    </button>
                  </th>
                  <th
                    scope="col"
                    class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                  >
                    <button
                      type="button"
                      @click="pegawaiStore.toggleSort('tanggal_lahir')"
                      class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                    >
                      Tanggal Lahir
                      <component
                        :is="sortIcon('tanggal_lahir')"
                        class="h-4 w-4"
                      />
                    </button>
                  </th>
                </template>

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
              <template v-if="pegawaiStore.adminLoading">
                <tr v-for="n in 6" :key="`skeleton-${n}`">
                  <td v-for="c in 10" :key="c" class="px-4 py-3">
                    <div
                      class="h-4 w-20 animate-pulse rounded bg-gray-100 dark:bg-white/5"
                    ></div>
                  </td>
                </tr>
              </template>

              <tr
                v-else
                v-for="item in pegawaiStore.adminPegawai"
                :key="item.nuptk"
                class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
              >
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
                      class="h-10 w-10 rounded-full object-cover transition duration-200 group-hover:opacity-80"
                    />
                  </button>
                  <div
                    v-else
                    class="flex h-10 w-10 items-center justify-center rounded-full bg-gray-100 dark:bg-slate-800"
                  >
                    <AcademicCapIcon class="h-5 w-5 text-gray-300" />
                  </div>
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.nip_baru || "-" }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
                >
                  <RouterLink
                    :to="{
                      name: 'dashboard-pegawai-edit',
                      params: { nuptk: item.nuptk },
                    }"
                    class="hover:text-[#041983] hover:underline dark:hover:text-indigo-300"
                  >
                    {{ item.nama_pegawai }}
                  </RouterLink>
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.pns_non || "-" }}
                </td>
                <td class="px-4 py-2.5">
                  <span
                    v-if="!item.status_keluar"
                    class="inline-flex items-center rounded-full bg-emerald-50 px-2.5 py-0.5 text-xs font-medium text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400"
                  >
                    Aktif
                  </span>
                  <span
                    v-else
                    class="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-xs font-medium text-gray-600 dark:bg-slate-800 dark:text-slate-400"
                  >
                    Keluar
                  </span>
                </td>

                <template v-if="isKeluarView">
                  <td
                    class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                  >
                    {{ formatDateNumeric(item.tanggal_keluar) }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                  >
                    {{ item.nomor_surat || "-" }}
                  </td>
                  <td
                    class="max-w-60 truncate px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                    :title="item.keterangan_keluar"
                  >
                    {{ item.keterangan_keluar || "-" }}
                  </td>
                </template>
                <template v-else>
                  <td
                    class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                  >
                    {{ item.jenis_kelamin }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                  >
                    {{ item.agama }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                  >
                    {{ item.tempat_lahir }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                  >
                    {{ formatDateNumeric(item.tanggal_lahir) }}
                  </td>
                </template>
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
        </div>
        <DashboardPagination
          :pagination="pegawaiStore.pagination"
          @change-page="changePage"
        />
      </div>
    </div>
  </div>

  <DashboardUploadDialog
    v-model="showUploadDialog"
    title="Unggah Pegawai"
    subtitle="Unggah data pegawai menggunakan Template (.xlsx)"
    :upload-action="pegawaiStore.importPegawai"
    @download-template="handleDownloadTemplate"
    @success="handleImportSuccess"
  />

  <DashboardConfirmDialog
    v-model="showDeleteDialog"
    title="Hapus Pegawai"
    :message="`Anda akan menghapus data pegawai &quot;${pegawaiToDelete?.nama_pegawai}&quot;.\n\nData yang telah dihapus tidak dapat dikembalikan.`"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="loadingAction"
    @confirm="handleDelete"
  />

  <VueEasyLightbox
    :visible="visibleLightbox"
    :imgs="lightboxImages"
    :index="0"
    @hide="visibleLightbox = false"
  />
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";
import VueEasyLightbox from "vue-easy-lightbox";

import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";
import DashboardUploadDialog from "@/components/dashboard/DashboardUploadDialog.vue";

import { useToast } from "@/composables/useToast";
import { usePegawaiStore } from "@/stores/pegawai";
import { formatDateNumeric } from "@/utils/date";

import {
  ArrowDownTrayIcon,
  ChevronUpIcon,
  ChevronDownIcon,
  ChevronUpDownIcon,
  MagnifyingGlassIcon,
  PlusIcon,
  TrashIcon,
  AcademicCapIcon,
  ArrowUpTrayIcon,
  DocumentArrowDownIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const pegawaiStore = usePegawaiStore();
const toast = useToast();

const showDeleteDialog = ref(false);
const loadingAction = ref(false);
const pegawaiToDelete = ref(null);

const statusFilter = ref(pegawaiStore.filters.status);

const isSearching = computed(() => pegawaiStore.filters.q.trim() !== "");

const isKeluarView = computed(() => Number(pegawaiStore.filters.status) === 0);

// ===== Lightbox foto (pola sama spt GaleriFoto.vue / pesertadidik Tampil.vue) =====
const visibleLightbox = ref(false);
const lightboxImages = ref([]);
function openLightbox(item) {
  if (!item.foto_url) return;
  lightboxImages.value = [item.foto_url];
  visibleLightbox.value = true;
}

const showUploadDialog = ref(false);

async function handleDownloadTemplate() {
  try {
    await pegawaiStore.downloadTemplate();
  } catch (error) {
    toast.error("Gagal mengunduh template import.");
  }
}

function handleImportSuccess(message) {
  toast.success(message || "Data pegawai berhasil diimport.");
}

let searchTimeout = null;

function searchData() {
  clearTimeout(searchTimeout);

  searchTimeout = setTimeout(async () => {
    pegawaiStore.pagination.page = 1;
    await pegawaiStore.fetchAdminPegawai();
  }, 300);
}

function changeStatusFilter() {
  pegawaiStore.filters.status = Number(statusFilter.value);
  pegawaiStore.pagination.page = 1;

  if (pegawaiStore.filters.status === 0) {
    // masuk ke tampilan "Keluar" -> auto-sort tanggal keluar terbaru
    pegawaiStore.sortBy = "tanggal_keluar";
    pegawaiStore.sortDir = "desc";
  } else {
    // keluar dari tampilan "Keluar" -> kembali ke sort default nama A-Z
    pegawaiStore.sortBy = "nama";
    pegawaiStore.sortDir = "asc";
  }

  pegawaiStore.fetchAdminPegawai();
}

async function changePage(page) {
  pegawaiStore.pagination.page = page;
  await pegawaiStore.fetchAdminPegawai();
}

function confirmDelete(item) {
  pegawaiToDelete.value = item;
  showDeleteDialog.value = true;
}

async function handleDelete() {
  if (!pegawaiToDelete.value) return;

  loadingAction.value = true;

  try {
    await pegawaiStore.deletePegawai(pegawaiToDelete.value.nuptk);
    toast.success("Data pegawai berhasil dihapus.");
    showDeleteDialog.value = false;
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menghapus data pegawai.",
    );
  } finally {
    loadingAction.value = false;
    pegawaiToDelete.value = null;
  }
}

async function handleExport() {
  try {
    await pegawaiStore.exportPegawai();
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal export data pegawai.");
  }
}

function sortIcon(column) {
  if (pegawaiStore.sortBy !== column) return ChevronUpDownIcon;
  return pegawaiStore.sortDir === "asc" ? ChevronUpIcon : ChevronDownIcon;
}

const kategoriFilter = ref(pegawaiStore.filters.kategori);

function changeKategoriFilter() {
  pegawaiStore.filters.kategori = kategoriFilter.value;
  pegawaiStore.pagination.page = 1;
  pegawaiStore.fetchAdminPegawai();
}

const pnsNonFilter = ref(pegawaiStore.filters.pnsNon);

function changePnsNonFilter() {
  pegawaiStore.filters.pnsNon = pnsNonFilter.value;
  pegawaiStore.pagination.page = 1;
  pegawaiStore.fetchAdminPegawai();
}
onMounted(() => pegawaiStore.fetchAdminPegawai());
</script>
