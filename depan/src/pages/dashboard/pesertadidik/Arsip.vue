<template>
  <DashboardPageHeader />
  <DashboardBulkAction
    :selected-count="selectedIds.length"
    @clear="selectedIds = []"
  >
    <template #actions>
      <button
        type="button"
        @click="cetakBulk"
        :disabled="bukuIndukStore.cetakLoading"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-sky-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-sky-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-500/40 disabled:cursor-not-allowed disabled:opacity-60"
      >
        <ArrowPathIcon
          v-if="bukuIndukStore.cetakLoading"
          class="h-5 w-5 animate-spin"
        />
        <PrinterIcon v-else class="h-5 w-5" />
        {{ bukuIndukStore.cetakLoading ? "Mencetak..." : "Buku Induk" }}
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
        <div
          class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/8 text-[#041983] dark:bg-white/5 dark:text-indigo-300"
            >
              <ArchiveBoxIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Daftar Peserta Didik Arsip
            </h2>
          </div>
          <div class="relative w-64 max-w-55 sm:max-w-xs md:max-w-sm">
            <MagnifyingGlassIcon
              class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />
            <input
              v-model="pdStore.arsipFilters.q"
              @input="searchData"
              type="text"
              placeholder="Pencarian..."
              :disabled="!tahunFilter"
              class="w-full rounded-xl border border-gray-300 bg-white py-2.5 pl-10 pr-4 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 disabled:cursor-not-allowed disabled:bg-gray-50 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500 dark:disabled:bg-slate-800/50"
            />
          </div>
        </div>

        <div
          class="mt-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
        >
          <div class="flex flex-wrap items-center gap-3">
            <select
              v-model="tahunFilter"
              @change="changeTahunFilter"
              class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
            >
              <option :value="0">Pilih Tahun Pelajaran</option>
              <option
                v-for="tp in tahunPelajaranOptions"
                :key="tp.id_tahunpelajaran"
                :value="tp.id_tahunpelajaran"
              >
                {{ tp.nama_tahunpelajaran }}
              </option>
            </select>
            <template v-if="tahunFilter">
              <div class="flex items-center gap-2">
                <span class="text-sm text-slate-500 dark:text-slate-400">
                  Tampilkan
                </span>
                <select
                  :value="pdStore.arsipPagination.limit"
                  @change="changeLimit($event.target.value)"
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
                v-model="kelasFilter"
                @change="changeKelasFilter"
                class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
              >
                <option :value="0">Semua Kelas</option>
                <option
                  v-for="k in kelasOptions"
                  :key="k.id_kelas"
                  :value="k.id_kelas"
                >
                  {{ k.nama_kelas }}{{ k.aktif === 1 ? "" : " (Nonaktif)" }}
                </option>
              </select>
              <select
                v-model="rombelFilter"
                @change="changeRombelFilter"
                class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
              >
                <option :value="0">Semua Rombel</option>
                <option
                  v-for="r in rombelOptions"
                  :key="r.id_rombel"
                  :value="r.id_rombel"
                >
                  {{ r.nama_rombel }}{{ r.aktif === 1 ? "" : " (Nonaktif)" }}
                </option>
              </select>
              <label
                class="flex items-center gap-1.5 text-sm text-slate-600 dark:text-slate-300"
              >
                <input
                  type="checkbox"
                  :checked="isAllSelected"
                  @change="toggleSelectAll"
                  class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                />
                Pilih Semua
              </label>
            </template>
          </div>
          <div class="flex items-center gap-2">
            <RouterLink
              :to="{ name: 'dashboard-pesertadidik' }"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-amber-500 hover:bg-amber-50 hover:text-amber-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-amber-500 dark:hover:bg-amber-500/10 dark:hover:text-amber-400"
            >
              <ArrowLeftIcon class="h-5 w-5" />
            </RouterLink>
            <button
              type="button"
              title="Unduh Peserta Didik Arsip"
              @click="handleExport"
              :disabled="!tahunFilter"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-emerald-500 dark:hover:bg-emerald-500/10 dark:hover:text-emerald-400"
            >
              <ArrowDownTrayIcon class="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>

      <!-- Belum pilih tahun pelajaran -->
      <div v-if="!tahunFilter" class="py-16 text-center">
        <ArchiveBoxIcon
          class="mx-auto h-10 w-10 text-gray-300 dark:text-slate-700"
        />
        <p class="mt-3 text-sm text-gray-500 dark:text-slate-400">
          Pilih tahun pelajaran terlebih dahulu untuk melihat data arsip.
        </p>
      </div>

      <!-- Loading -->
      <div v-else-if="pdStore.arsipLoading" class="p-6">
        <div class="space-y-3">
          <div
            v-for="n in 6"
            :key="`skeleton-${n}`"
            class="flex items-center gap-4 py-1"
          >
            <div
              class="h-10 w-10 shrink-0 animate-pulse rounded-full bg-gray-100 dark:bg-white/5"
            ></div>
            <div class="flex-1 space-y-2">
              <div
                class="h-4 w-40 animate-pulse rounded bg-gray-100 dark:bg-white/5"
              ></div>
              <div
                class="h-3 w-56 animate-pulse rounded bg-gray-100 dark:bg-white/5"
              ></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div
        v-else-if="pdStore.adminArsip.length === 0"
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
                : "Belum ada data arsip untuk tahun pelajaran ini."
            }}
          </p>
        </div>
      </div>

      <!-- Grouped list -->
      <div v-else>
        <div
          v-for="group in groupedArsip"
          :key="group.label"
          class="border-b border-gray-100 last:border-b-0 dark:border-white/5"
        >
          <div
            class="sticky top-0 z-20 bg-gray-50 px-6 py-2.5 text-sm font-semibold text-gray-700 dark:bg-[#0B0F1A] dark:text-slate-200"
          >
            {{ group.label }}
            <span
              class="ml-2 text-xs font-normal text-gray-400 dark:text-slate-500"
            >
              ({{ group.items.length }} siswa)
            </span>
          </div>
          <div class="divide-y divide-gray-100 dark:divide-white/5">
            <RouterLink
              v-for="item in group.items"
              :key="item.id"
              :to="{
                name: 'dashboard-pesertadidik-arsip-detil',
                params: { id: item.id },
              }"
              class="flex items-center gap-4 px-6 py-3 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:-outline-offset-2 focus-visible:outline-[#041983] dark:hover:bg-[#0B0F1A]/40"
            >
              <input
                type="checkbox"
                :checked="selectedIds.includes(item.id)"
                @click.stop
                @change="toggleSelect(item.id)"
                class="h-4 w-4 shrink-0 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              />
              <button
                v-if="item.foto_url"
                type="button"
                title="Klik untuk perbesar"
                class="group block shrink-0 cursor-pointer"
                @click.stop.prevent="openLightbox(item)"
              >
                <img
                  :src="item.foto_url"
                  class="h-10 w-10 rounded-full object-cover transition duration-200 group-hover:opacity-80"
                />
              </button>
              <div
                v-else
                class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-gray-100 dark:bg-slate-800"
              >
                <UserIcon class="h-5 w-5 text-gray-300" />
              </div>
              <div class="min-w-0 flex-1">
                <p
                  class="truncate text-sm font-medium text-gray-800 dark:text-slate-200"
                >
                  {{ item.nama_siswa }}
                </p>
                <p class="text-xs text-gray-500 dark:text-slate-400">
                  NIS: {{ item.nis }} · NISN: {{ item.nisn || "-" }} ·
                  {{ item.jenis_kelamin }}
                </p>
              </div>
              <span
                class="shrink-0 rounded-full bg-purple-50 px-2.5 py-0.5 text-xs font-medium text-purple-700 dark:bg-purple-500/10 dark:text-purple-400"
              >
                {{ item.status_keluar || "Arsip" }}
              </span>
              <button
                type="button"
                title="Cetak Buku Induk"
                @click.stop.prevent="cetakSatu(item.id)"
                :disabled="bukuIndukStore.cetakLoading"
                class="shrink-0 rounded-lg border border-slate-300 p-1.5 text-slate-500 transition hover:border-sky-500 hover:bg-sky-50 hover:text-sky-600 disabled:cursor-not-allowed disabled:opacity-50 dark:border-slate-700 dark:text-slate-400 dark:hover:border-sky-500 dark:hover:bg-sky-500/10"
              >
                <ArrowPathIcon
                  v-if="cetakingId === item.id"
                  class="h-4 w-4 animate-spin"
                />
                <PrinterIcon v-else class="h-4 w-4" />
              </button>
              <ChevronRightIcon
                class="h-4 w-4 shrink-0 text-gray-300 dark:text-slate-600"
              />
            </RouterLink>
          </div>
        </div>

        <DashboardPagination
          :pagination="pdStore.arsipPagination"
          @change-page="changePage"
        />
      </div>
    </div>
  </div>

  <VueEasyLightbox
    :visible="visibleLightbox"
    :imgs="lightboxImages"
    :index="0"
    @hide="visibleLightbox = false"
  />
</template>

<script setup>
import { computed, ref } from "vue";
import { RouterLink } from "vue-router";
import VueEasyLightbox from "vue-easy-lightbox";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import DashboardBulkAction from "@/components/dashboard/DashboardBulkActionBar.vue";
import { usePesertaDidikStore } from "@/stores/pesertadidik";
import { useBukuIndukStore } from "@/stores/bukuinduk";
import { useTahunPelajaranStore } from "@/stores/tahunpelajaran";
import { useKelasStore } from "@/stores/kelas";
import { useRombelStore } from "@/stores/rombel";
import { useToast } from "@/composables/useToast";
import {
  ArchiveBoxIcon,
  ArrowLeftIcon,
  ChevronRightIcon,
  MagnifyingGlassIcon,
  UserIcon,
  ArrowDownTrayIcon,
  InboxIcon,
  PrinterIcon,
  ArrowPathIcon,
} from "@heroicons/vue/24/outline";

const pdStore = usePesertaDidikStore();
const bukuIndukStore = useBukuIndukStore();
const tahunPelajaranStore = useTahunPelajaranStore();
const kelasStore = useKelasStore();
const rombelStore = useRombelStore();

const toast = useToast();

const tahunPelajaranOptions = computed(() =>
  (tahunPelajaranStore.list || []).filter((tp) => tp.aktif !== 1),
);
// Arsip = data historis, id_kelas/id_rombel yang tersimpan bisa saja sudah
// tidak aktif lagi sekarang (struktur kelas/rombel berubah tiap tahun) — jadi
// TIDAK difilter aktif spt di Tampil.vue, tampilkan semua biar tetap bisa dipilih
const kelasOptions = computed(() => kelasStore.list || []);
const rombelOptions = computed(() => rombelStore.list || []);

const tahunFilter = ref(pdStore.arsipFilters.idTahunPelajaran || 0);
const kelasFilter = ref(pdStore.arsipFilters.idKelas || 0);
const rombelFilter = ref(pdStore.arsipFilters.idRombel || 0);
const isSearching = computed(() => pdStore.arsipFilters.q.trim() !== "");

// ===== Lightbox foto (pola sama spt GaleriFoto.vue) =====
const visibleLightbox = ref(false);
const lightboxImages = ref([]);
function openLightbox(item) {
  if (!item.foto_url) return;
  lightboxImages.value = [item.foto_url];
  visibleLightbox.value = true;
}

// Data dari backend sudah diurutkan per kelas -> rombel -> nama (lihat GetAllArsip),
// jadi tinggal dikelompokkan visual di sini. Catatan: pengelompokan hanya berlaku
// dalam 1 halaman pagination — kalau 1 kelas/rombel terpotong ke halaman berikutnya,
// grup yang sama bisa muncul lagi di halaman itu.
const groupedArsip = computed(() => {
  const groups = [];
  let current = null;
  for (const item of pdStore.adminArsip) {
    const label = `${item.nama_kelas || "-"} / ${item.nama_rombel || "-"}`;
    if (!current || current.label !== label) {
      current = { label, items: [] };
      groups.push(current);
    }
    current.items.push(item);
  }
  return groups;
});

// ===== Checklist & cetak Buku Induk =====
const selectedIds = ref([]);
const isAllSelected = computed(
  () =>
    pdStore.adminArsip.length > 0 &&
    pdStore.adminArsip.every((item) => selectedIds.value.includes(item.id)),
);

function toggleSelect(id) {
  const idx = selectedIds.value.indexOf(id);
  if (idx === -1) {
    selectedIds.value.push(id);
  } else {
    selectedIds.value.splice(idx, 1);
  }
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedIds.value = selectedIds.value.filter(
      (id) => !pdStore.adminArsip.some((item) => item.id === id),
    );
  } else {
    const pageIds = pdStore.adminArsip.map((item) => item.id);
    selectedIds.value = [...new Set([...selectedIds.value, ...pageIds])];
  }
}

// tahun pelajaran utk label Cover/Bagian V saat cetak — ikut tahun yg
// SEDANG DIPILIH di filter (lebih akurat drpd fallback tahun aktif,
// krn data arsip memang terikat ke tahun ajaran tertentu di masa lalu)
const resolvedTahunPelajaran = computed(() => {
  const tp = tahunPelajaranStore.list.find(
    (t) => t.id_tahunpelajaran === Number(tahunFilter.value),
  );
  return tp ? tp.nama_tahunpelajaran : "";
});

const cetakingId = ref(null);

async function cetakSatu(id) {
  if (!resolvedTahunPelajaran.value) {
    toast.error("Tahun pelajaran tidak valid.");
    return;
  }
  cetakingId.value = id;
  try {
    await bukuIndukStore.cetakSatu("arsip", id, resolvedTahunPelajaran.value);
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal mencetak buku induk.");
  } finally {
    cetakingId.value = null;
  }
}

async function cetakBulk() {
  if (!resolvedTahunPelajaran.value) {
    toast.error("Tahun pelajaran tidak valid.");
    return;
  }
  try {
    await bukuIndukStore.cetakBulk(
      "arsip",
      selectedIds.value,
      resolvedTahunPelajaran.value,
    );
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal mencetak buku induk massal.",
    );
  }
}

function changeTahunFilter() {
  pdStore.arsipFilters.idTahunPelajaran = Number(tahunFilter.value);
  pdStore.arsipPagination.page = 1;
  selectedIds.value = [];
  if (tahunFilter.value) {
    pdStore.fetchAdminArsip();
  } else {
    pdStore.adminArsip = [];
  }
}

function changeLimit(limit) {
  pdStore.arsipPagination.limit = Number(limit);
  pdStore.arsipPagination.page = 1;
  pdStore.fetchAdminArsip();
}

function changeKelasFilter() {
  pdStore.arsipFilters.idKelas = Number(kelasFilter.value);
  pdStore.arsipPagination.page = 1;
  pdStore.fetchAdminArsip();
}
function changeRombelFilter() {
  pdStore.arsipFilters.idRombel = Number(rombelFilter.value);
  pdStore.arsipPagination.page = 1;
  pdStore.fetchAdminArsip();
}

let searchTimeout = null;
function searchData() {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    pdStore.arsipPagination.page = 1;
    pdStore.fetchAdminArsip();
  }, 300);
}

async function changePage(page) {
  pdStore.arsipPagination.page = page;
  await pdStore.fetchAdminArsip();
}

async function handleExport() {
  try {
    await pdStore.exportArsip();
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal export data arsip.");
  }
}

tahunPelajaranStore.pagination.limit = 50;
tahunPelajaranStore.fetchAll();
kelasStore.pagination.limit = 100;
rombelStore.pagination.limit = 500;
kelasStore.fetchAll();
rombelStore.fetchAll();
</script>
