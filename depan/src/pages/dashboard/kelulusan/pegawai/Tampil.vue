<template>
  <DashboardPageHeader />

  <div class="space-y-6">
    <div
      class="relative overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
    >
      <div
        class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <div class="border-b border-gray-100 px-6 py-5 dark:border-white/5">
        <div
          class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/8 text-[#041983] dark:bg-white/5 dark:text-indigo-300"
            >
              <UserIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Data Kelulusan
            </h2>
          </div>
          <div class="relative w-64 max-w-55 sm:max-w-xs md:max-w-sm">
            <MagnifyingGlassIcon
              class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />
            <input
              v-model="kelulusanStore.filters.q"
              @input="searchData"
              type="text"
              placeholder="Pencarian..."
              class="w-full rounded-xl border border-gray-300 bg-white py-2.5 pl-10 pr-4 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500"
            />
          </div>
        </div>
        <div
          class="mt-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
        >
          <div class="flex flex-wrap items-center gap-3">
            <div class="flex items-center gap-2">
              <span class="text-sm text-slate-500 dark:text-slate-400">
                Tampilkan
              </span>
              <select
                :value="kelulusanStore.pagination.limit"
                @change="kelulusanStore.changeLimit($event.target.value)"
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
        </div>
      </div>

      <div
        v-if="
          !kelulusanStore.adminLoading &&
          kelulusanStore.adminKelulusan.length === 0
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
              isSearching
                ? "Data tidak ditemukan."
                : "Belum ada data kelulusan."
            }}
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
                  Foto
                </th>

                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="kelulusanStore.toggleSort('nama_rombel')"
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
                    @click="kelulusanStore.toggleSort('nis')"
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
                    @click="kelulusanStore.toggleSort('nopes')"
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
                    @click="kelulusanStore.toggleSort('nama_siswa')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Nama Peserta Didik
                    <component :is="sortIcon('nama_siswa')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="kelulusanStore.toggleSort('rata_s')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Rerata
                    <component :is="sortIcon('rata_s')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="kelulusanStore.toggleSort('ket')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Keterangan
                    <component :is="sortIcon('ket')" class="h-4 w-4" />
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
              <template v-if="kelulusanStore.adminLoading">
                <tr v-for="n in 6" :key="`skeleton-${n}`">
                  <td v-for="c in 7" :key="c" class="px-4 py-3">
                    <div
                      class="h-4 w-20 animate-pulse rounded bg-gray-100 dark:bg-white/5"
                    ></div>
                  </td>
                </tr>
              </template>

              <tr
                v-else
                v-for="item in kelulusanStore.adminKelulusan"
                :key="item.nis"
                class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
              >
                <td class="px-4 py-2.5">
                  <button
                    type="button"
                    title="Klik untuk perbesar"
                    class="group block cursor-pointer"
                    @click="openLightbox(item)"
                  >
                    <img
                      :src="item.foto_url"
                      :alt="item.nama_siswa"
                      class="h-9 w-9 rounded-full object-cover ring-1 ring-gray-200 transition duration-200 group-hover:opacity-80 dark:ring-white/10"
                    />
                  </button>
                </td>

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
                  {{ item.nopes || "-" }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
                >
                  {{ item.nama_siswa || "-" }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.rata_s || "-" }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.ket || "-" }}
                </td>
                <td class="px-4 py-2.5 text-right">
                  <button
                    type="button"
                    title="Cetak SKL"
                    @click="handleCetak(item.nis)"
                    :disabled="cetakLoadingNis === item.nis"
                    class="inline-flex items-center justify-center rounded-lg border border-slate-300 bg-white p-2 text-slate-600 shadow-sm transition hover:border-blue-500 hover:bg-blue-50 hover:text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:opacity-50 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-blue-500 dark:hover:bg-blue-500/10 dark:hover:text-blue-400"
                  >
                    <PrinterIcon class="h-4 w-4" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <DashboardPagination
          :pagination="kelulusanStore.pagination"
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
import VueEasyLightbox from "vue-easy-lightbox";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import { useToast } from "@/composables/useToast";
import { useKelulusanStore } from "@/stores/kelulusan";
import {
  ChevronDownIcon,
  ChevronUpIcon,
  ChevronUpDownIcon,
  MagnifyingGlassIcon,
  PrinterIcon,
  UserIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const kelulusanStore = useKelulusanStore();
const toast = useToast();

const cetakLoadingNis = ref(null);

const isSearching = computed(() => kelulusanStore.filters.q.trim() !== "");

// ===== Lightbox foto (pola sama spt GaleriFoto.vue) =====
const visibleLightbox = ref(false);
const lightboxImages = ref([]);
function openLightbox(item) {
  if (!item.foto_url) return;
  lightboxImages.value = [item.foto_url];
  visibleLightbox.value = true;
}

let searchTimeout = null;
function searchData() {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(async () => {
    kelulusanStore.pagination.page = 1;
    await kelulusanStore.fetchAdminKelulusan();
  }, 300);
}
async function changePage(page) {
  kelulusanStore.pagination.page = page;
  await kelulusanStore.fetchAdminKelulusan();
}
function sortIcon(column) {
  if (kelulusanStore.sortBy !== column) return ChevronUpDownIcon;
  return kelulusanStore.sortDir === "asc" ? ChevronUpIcon : ChevronDownIcon;
}

async function handleCetak(nis) {
  cetakLoadingNis.value = nis;
  try {
    await kelulusanStore.cetakSKL(nis);
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal mencetak SKL.");
  } finally {
    cetakLoadingNis.value = null;
  }
}

onMounted(() => {
  kelulusanStore.fetchAdminKelulusan();
});
</script>
