<template>
  <DashboardBulkAction
    :selected-count="selectedNis.length"
    @clear="selectedNis = []"
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
        <div
          class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/8 text-[#041983] dark:bg-white/5 dark:text-indigo-300"
            >
              <ArrowsRightLeftIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Daftar Peserta Didik Pindah / Lainnya
            </h2>
          </div>
          <div class="relative w-64 max-w-55 sm:max-w-xs md:max-w-sm">
            <MagnifyingGlassIcon
              class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />
            <input
              v-model="pdStore.mutasiFilters.q"
              @input="searchData"
              type="text"
              placeholder="Pencarian..."
              class="w-full rounded-xl border border-gray-300 bg-white py-2.5 pl-10 pr-4 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500"
            />
          </div>
        </div>
        <div class="mt-5 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span class="text-sm text-slate-500 dark:text-slate-400">
              Tampilkan
            </span>
            <select
              :value="pdStore.mutasiPagination.limit"
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
          <div class="flex items-center gap-2">
            <RouterLink
              :to="{ name: 'dashboard-pesertadidik' }"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-amber-500 hover:bg-amber-50 hover:text-amber-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-amber-500 dark:hover:bg-amber-500/10 dark:hover:text-amber-400"
            >
              <ArrowLeftIcon class="h-5 w-5" />
            </RouterLink>
            <button
              type="button"
              title="Unduh Peserta Didik Mutasi"
              @click="handleExport"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-emerald-500 dark:hover:bg-emerald-500/10 dark:hover:text-emerald-400"
            >
              <ArrowDownTrayIcon class="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div
        v-if="!pdStore.mutasiLoading && pdStore.adminMutasi.length === 0"
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <MagnifyingGlassIcon v-if="isSearching" class="h-10 w-10" />
          <InboxIcon v-else class="h-10 w-10" />
          <p class="text-sm font-medium">
            {{
              isSearching ? "Data tidak ditemukan." : "Belum ada data mutasi."
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
                  Foto
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Rombel
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  NIS
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Nama Peserta Didik
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Status
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Tanggal Keluar
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  No Surat
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Tujuan
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Alasan
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-white/5">
              <!-- Skeleton saat memuat -->
              <template v-if="pdStore.mutasiLoading">
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
                v-for="item in pdStore.adminMutasi"
                :key="item.nis"
                class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
              >
                <td class="px-4 py-2.5">
                  <input
                    type="checkbox"
                    :checked="selectedNis.includes(item.nis)"
                    @change="toggleSelect(item.nis)"
                    class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </td>
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
                    <UserIcon class="h-5 w-5 text-gray-300" />
                  </div>
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.nama_rombel }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.nis }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
                >
                  <RouterLink
                    :to="{
                      name: 'dashboard-pesertadidik-mutasi-detil',
                      params: { nis: item.nis },
                    }"
                    class="hover:text-[#041983] hover:underline dark:hover:text-indigo-300"
                  >
                    {{ item.nama_siswa }}
                  </RouterLink>
                </td>
                <td class="px-4 py-2.5">
                  <span
                    class="inline-flex items-center rounded-full bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-700 dark:bg-amber-500/10 dark:text-amber-400"
                  >
                    {{ item.status_keluar || "-" }}
                  </span>
                </td>
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
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.tujuan || "-" }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.alasan_keluar || "-" }}
                </td>
                <td class="px-4 py-2.5">
                  <button
                    type="button"
                    title="Cetak Buku Induk"
                    @click="cetakSatu(item.nis)"
                    :disabled="bukuIndukStore.cetakLoading"
                    class="rounded-lg border border-slate-300 p-1.5 text-slate-500 transition hover:border-sky-500 hover:bg-sky-50 hover:text-sky-600 disabled:cursor-not-allowed disabled:opacity-50 dark:border-slate-700 dark:text-slate-400 dark:hover:border-sky-500 dark:hover:bg-sky-500/10"
                  >
                    <ArrowPathIcon
                      v-if="cetakingId === item.nis"
                      class="h-4 w-4 animate-spin"
                    />
                    <PrinterIcon v-else class="h-4 w-4" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <DashboardPagination
          :pagination="pdStore.mutasiPagination"
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
import { computed, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";
import VueEasyLightbox from "vue-easy-lightbox";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import DashboardBulkAction from "@/components/dashboard/DashboardBulkActionBar.vue";
import { usePesertaDidikStore } from "@/stores/pesertadidik";
import { useBukuIndukStore } from "@/stores/bukuinduk";
import { useTahunPelajaranStore } from "@/stores/tahunpelajaran";
import { useToast } from "@/composables/useToast";
import { formatDateNumeric } from "@/utils/date";
import {
  ArrowLeftIcon,
  MagnifyingGlassIcon,
  UserIcon,
  ArrowDownTrayIcon,
  ArrowsRightLeftIcon,
  InboxIcon,
  PrinterIcon,
  ArrowPathIcon,
} from "@heroicons/vue/24/outline";

const pdStore = usePesertaDidikStore();
const bukuIndukStore = useBukuIndukStore();
const tahunPelajaranStore = useTahunPelajaranStore();

const toast = useToast();

const selectedNis = ref([]);
const isAllSelected = computed(
  () =>
    pdStore.adminMutasi.length > 0 &&
    pdStore.adminMutasi.every((item) => selectedNis.value.includes(item.nis)),
);

function toggleSelect(nis) {
  const idx = selectedNis.value.indexOf(nis);
  if (idx === -1) {
    selectedNis.value.push(nis);
  } else {
    selectedNis.value.splice(idx, 1);
  }
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedNis.value = selectedNis.value.filter(
      (nis) => !pdStore.adminMutasi.some((item) => item.nis === nis),
    );
  } else {
    const pageNis = pdStore.adminMutasi.map((item) => item.nis);
    selectedNis.value = [...new Set([...selectedNis.value, ...pageNis])];
  }
}

// tahun pelajaran utk label Cover/Bagian V saat cetak — halaman Mutasi
// tidak punya filter tahun, jadi pakai tahun yg sedang aktif sbg fallback
const resolvedTahunPelajaran = computed(() => {
  const aktif = tahunPelajaranStore.list.find((t) => t.aktif === 1);
  return aktif ? aktif.nama_tahunpelajaran : "";
});

const cetakingId = ref(null);

async function cetakSatu(nis) {
  if (!resolvedTahunPelajaran.value) {
    toast.error("Tidak ada tahun pelajaran aktif.");
    return;
  }
  cetakingId.value = nis;
  try {
    await bukuIndukStore.cetakSatu("mutasi", nis, resolvedTahunPelajaran.value);
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal mencetak buku induk.");
  } finally {
    cetakingId.value = null;
  }
}

async function cetakBulk() {
  if (!resolvedTahunPelajaran.value) {
    toast.error("Tidak ada tahun pelajaran aktif.");
    return;
  }
  try {
    await bukuIndukStore.cetakBulk(
      "mutasi",
      selectedNis.value,
      resolvedTahunPelajaran.value,
    );
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal mencetak buku induk massal.",
    );
  }
}

const isSearching = computed(() => pdStore.mutasiFilters.q.trim() !== "");

// ===== Lightbox foto (pola sama spt GaleriFoto.vue) =====
const visibleLightbox = ref(false);
const lightboxImages = ref([]);
function openLightbox(item) {
  if (!item.foto_url) return;
  lightboxImages.value = [item.foto_url];
  visibleLightbox.value = true;
}

function changeLimit(limit) {
  pdStore.mutasiPagination.limit = Number(limit);
  pdStore.mutasiPagination.page = 1;
  pdStore.fetchAdminMutasi();
}

let searchTimeout = null;
function searchData() {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    pdStore.mutasiPagination.page = 1;
    pdStore.fetchAdminMutasi();
  }, 300);
}

async function changePage(page) {
  pdStore.mutasiPagination.page = page;
  await pdStore.fetchAdminMutasi();
}

async function handleExport() {
  try {
    await pdStore.exportMutasi();
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal export data mutasi.");
  }
}

onMounted(() => {
  tahunPelajaranStore.pagination.limit = 9999;
  tahunPelajaranStore.fetchAll();
  pdStore.fetchAdminMutasi();
});
</script>
