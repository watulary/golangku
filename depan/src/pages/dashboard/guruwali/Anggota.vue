<template>
  <DashboardPageHeader
    :title="`Anggota Guru Wali${guruwaliStore.detail ? ' — ' + guruwaliStore.detail.nama_guruwali : ''}`"
    subtitle="Kelola peserta didik yang terdaftar pada guru wali ini."
  />

  <DashboardBulkAction
    :selected-count="selectedAnggota.length"
    @clear="selectedAnggota = []"
  >
    <template #actions>
      <button
        type="button"
        @click="showBulkDeleteDialog = true"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-red-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-red-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40"
      >
        <TrashIcon class="h-5 w-5" />
        Hapus
      </button>
    </template>
  </DashboardBulkAction>

  <div class="space-y-6">
    <!-- Tambah Anggota -->
    <section
      class="relative overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
    >
      <div
        class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Tambah Anggota"
        description="Cari peserta didik untuk didaftarkan ke guru wali ini."
      >
        <UserPlusIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="p-6">
        <FormField label="Peserta Didik">
          <div class="relative">
            <input
              v-model="siswaQuery"
              @input="onSiswaSearch"
              @focus="showSiswaDropdown = true"
              type="text"
              :class="inputClass"
              placeholder="Ketik Nama atau NIS..."
              autocomplete="off"
            />

            <div
              v-if="showSiswaDropdown && siswaQuery.trim() !== ''"
              class="absolute z-10 mt-1 max-h-56 w-full overflow-y-auto rounded-xl bg-white shadow-lg ring-1 ring-slate-900/5 dark:bg-[#0B0F1A] dark:ring-white/10"
            >
              <div
                v-if="guruwaliStore.pesertaDidikSearching"
                class="px-4 py-3 text-sm text-gray-500 dark:text-slate-400"
              >
                Mencari...
              </div>

              <div
                v-else-if="guruwaliStore.pesertaDidikOptions.length === 0"
                class="px-4 py-3 text-sm text-gray-500 dark:text-slate-400"
              >
                Peserta didik tidak ditemukan.
              </div>

              <button
                v-for="s in guruwaliStore.pesertaDidikOptions"
                :key="s.nis"
                type="button"
                @click="selectSiswa(s)"
                class="flex w-full flex-col items-start px-4 py-2.5 text-left transition hover:bg-gray-50 dark:hover:bg-slate-800"
              >
                <span
                  class="text-sm font-medium text-gray-800 dark:text-slate-200"
                >
                  {{ s.nama_siswa }}
                </span>
                <span class="text-xs text-gray-500 dark:text-slate-400">
                  {{ s.nama_rombel ? `${s.nama_rombel}/${s.nis}` : s.nis }}
                </span>
              </button>
            </div>
          </div>

          <p
            v-if="selectedSiswa.nis"
            class="mt-1.5 text-xs text-emerald-600 dark:text-emerald-400"
          >
            Terpilih: {{ selectedSiswa.nama_siswa }} ({{
              selectedSiswa.nama_rombel
                ? `${selectedSiswa.nama_rombel}/${selectedSiswa.nis}`
                : selectedSiswa.nis
            }})
          </p>
        </FormField>

        <div class="mt-4 flex flex-wrap items-center gap-2">
          <button
            type="button"
            :disabled="!selectedSiswa.nis || addingAnggota"
            @click="handleAddAnggota"
            class="inline-flex items-center justify-center gap-2 rounded-xl bg-blue-600 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:bg-gray-400 dark:disabled:bg-slate-700"
          >
            <ArrowPathIcon v-if="addingAnggota" class="h-5 w-5 animate-spin" />
            <PlusIcon v-else class="h-5 w-5" />
            {{ addingAnggota ? "Menambahkan..." : "Tambah Anggota" }}
          </button>

          <button
            type="button"
            title="Unggah Anggota"
            @click="showUploadDialog = true"
            class="inline-flex shrink-0 items-center justify-center rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-yellow-500 hover:bg-yellow-50 hover:text-yellow-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-yellow-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-yellow-500 dark:hover:bg-yellow-500/10 dark:hover:text-yellow-300"
          >
            <ArrowUpTrayIcon class="h-5 w-5" />
          </button>
        </div>
      </div>
    </section>

    <!-- Daftar Anggota -->
    <div
      class="overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
    >
      <div class="border-b border-gray-100 px-6 py-5 dark:border-white/5">
        <!-- Baris 1 -->
        <div
          class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-blue-100 text-blue-600 dark:bg-blue-500/10 dark:text-blue-300"
            >
              <UserGroupIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Daftar Anggota
            </h2>
          </div>

          <div class="relative w-64 max-w-55 sm:max-w-xs md:max-w-sm">
            <MagnifyingGlassIcon
              class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />

            <input
              v-model="guruwaliStore.anggotaFilters.q"
              @input="searchAnggota"
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
          <!-- kiri: Tampilkan + filter Kelas & Rombel -->
          <div class="flex flex-wrap items-center gap-3">
            <div class="flex items-center gap-2">
              <span class="text-sm text-slate-500 dark:text-slate-400">
                Tampilkan
              </span>

              <select
                :value="guruwaliStore.anggotaPagination.limit"
                @change="
                  guruwaliStore.changeLimitAnggota(
                    $event.target.value,
                    idGuruWali,
                  )
                "
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
                {{ k.nama_kelas }}
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
                {{ r.nama_rombel }}
              </option>
            </select>
          </div>

          <!-- kanan -->
          <div class="flex items-center gap-2">
            <button
              type="button"
              title="Unduh Excel"
              :disabled="exportingAnggota"
              @click="handleExportAnggota"
              class="inline-flex shrink-0 items-center justify-center rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-emerald-500 dark:hover:bg-emerald-500/10 dark:hover:text-emerald-400"
            >
              <ArrowPathIcon
                v-if="exportingAnggota"
                class="h-5 w-5 animate-spin"
              />
              <ArrowDownTrayIcon v-else class="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div
        v-if="
          !guruwaliStore.anggotaLoading &&
          guruwaliStore.anggotaList.length === 0
        "
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <MagnifyingGlassIcon v-if="isSearchingAnggota" class="h-10 w-10" />
          <InboxIcon v-else class="h-10 w-10" />
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
      <div v-else>
        <div class="max-h-[70vh] overflow-auto">
          <table class="w-full min-w-160">
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
              </tr>
            </thead>

            <tbody class="divide-y divide-gray-100 dark:divide-white/5">
              <!-- Skeleton saat memuat -->
              <template v-if="guruwaliStore.anggotaLoading">
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
                v-for="item in guruwaliStore.anggotaList"
                :key="item.id"
                class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
              >
                <td class="w-10 px-4 py-2.5 text-center">
                  <input
                    type="checkbox"
                    :checked="selectedAnggota.includes(item.id)"
                    @change="toggleSelect(item.id)"
                    class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
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

        <DashboardPagination
          :pagination="guruwaliStore.anggotaPagination"
          @change-page="changePage"
        />
      </div>
    </div>

    <RouterLink
      :to="{ name: 'dashboard-guruwali' }"
      class="inline-flex items-center justify-center rounded-xl border border-gray-300 px-5 py-2.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
    >
      Kembali
    </RouterLink>
  </div>

  <DashboardConfirmDialog
    v-model="showBulkDeleteDialog"
    title="Hapus Anggota"
    :message="`Anda akan mengeluarkan ${selectedAnggota.length} anggota dari guru wali ini.\n\nData yang telah dihapus tidak dapat dikembalikan.`"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="bulkDeleteLoading"
    @confirm="handleBulkDelete"
  />

  <DashboardUploadDialog
    v-model="showUploadDialog"
    title="Unggah Anggota"
    subtitle="Unggah data anggota menggunakan Template (.xlsx)"
    :upload-action="(file) => guruwaliStore.importAnggota(idGuruWali, file)"
    @success="toast.success($event)"
    @download-template="downloadTemplate"
  />
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { RouterLink, useRoute } from "vue-router";

import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";
import DashboardUploadDialog from "@/components/dashboard/DashboardUploadDialog.vue";
import DashboardBulkAction from "@/components/dashboard/DashboardBulkActionBar.vue";
import FormField from "@/components/form/FormField.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";

import { useToast } from "@/composables/useToast";
import { useGuruWaliStore } from "@/stores/guruwali";
import { useKelasStore } from "@/stores/kelas";
import { useRombelStore } from "@/stores/rombel";

import {
  ArrowDownTrayIcon,
  ArrowPathIcon,
  ArrowUpTrayIcon,
  MagnifyingGlassIcon,
  PlusIcon,
  TrashIcon,
  UserPlusIcon,
  UserGroupIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const route = useRoute();
const guruwaliStore = useGuruWaliStore();
const kelasStore = useKelasStore();
const rombelStore = useRombelStore();
const toast = useToast();

const idGuruWali = route.params.id;

const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";

const exportingAnggota = ref(false);

const isSearchingAnggota = computed(
  () => guruwaliStore.anggotaFilters.q.trim() !== "",
);

// =====================================
// Filter Kelas & Rombel
// =====================================

const kelasOptions = computed(() =>
  (kelasStore.list || []).filter((k) => k.aktif === 1),
);
const rombelOptions = computed(() =>
  (rombelStore.list || []).filter((r) => r.aktif === 1),
);

const kelasFilter = ref(guruwaliStore.anggotaFilters.idKelas);
const rombelFilter = ref(guruwaliStore.anggotaFilters.idRombel);

function changeKelasFilter() {
  guruwaliStore.changeFilterKelasAnggota(kelasFilter.value, idGuruWali);
}
function changeRombelFilter() {
  guruwaliStore.changeFilterRombelAnggota(rombelFilter.value, idGuruWali);
}

// =====================================
// Checklist & Bulk Delete
// =====================================

const selectedAnggota = ref([]);
const showBulkDeleteDialog = ref(false);
const bulkDeleteLoading = ref(false);

const isAllSelected = computed(
  () =>
    guruwaliStore.anggotaList.length > 0 &&
    guruwaliStore.anggotaList.every((item) =>
      selectedAnggota.value.includes(item.id),
    ),
);

function toggleSelect(id) {
  const idx = selectedAnggota.value.indexOf(id);
  if (idx === -1) {
    selectedAnggota.value.push(id);
  } else {
    selectedAnggota.value.splice(idx, 1);
  }
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedAnggota.value = selectedAnggota.value.filter(
      (id) => !guruwaliStore.anggotaList.some((item) => item.id === id),
    );
  } else {
    const pageIds = guruwaliStore.anggotaList.map((item) => item.id);
    selectedAnggota.value = [
      ...new Set([...selectedAnggota.value, ...pageIds]),
    ];
  }
}

async function handleBulkDelete() {
  if (selectedAnggota.value.length === 0) return;

  bulkDeleteLoading.value = true;

  try {
    const res = await guruwaliStore.bulkRemoveAnggota(
      selectedAnggota.value,
      idGuruWali,
    );
    toast.success(res.message || "Anggota berhasil dihapus.");
    selectedAnggota.value = [];
    showBulkDeleteDialog.value = false;
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menghapus anggota.");
  } finally {
    bulkDeleteLoading.value = false;
  }
}

// =====================================
// Pencarian Siswa (autocomplete tambah anggota)
// =====================================

const siswaQuery = ref("");
const showSiswaDropdown = ref(false);
const addingAnggota = ref(false);

const selectedSiswa = reactive({
  nis: "",
  nama_siswa: "",
  nama_rombel: "",
});

let searchTimeout = null;

function onSiswaSearch() {
  showSiswaDropdown.value = true;
  selectedSiswa.nis = "";
  selectedSiswa.nama_siswa = "";
  selectedSiswa.nama_rombel = "";

  clearTimeout(searchTimeout);

  searchTimeout = setTimeout(() => {
    if (siswaQuery.value.trim() !== "") {
      guruwaliStore.searchPesertaDidik(siswaQuery.value.trim());
    }
  }, 300);
}

function selectSiswa(siswa) {
  selectedSiswa.nis = siswa.nis;
  selectedSiswa.nama_siswa = siswa.nama_siswa;
  selectedSiswa.nama_rombel = siswa.nama_rombel;
  siswaQuery.value = siswa.nama_siswa;
  showSiswaDropdown.value = false;
}

async function handleAddAnggota() {
  if (!selectedSiswa.nis || addingAnggota.value) return;

  addingAnggota.value = true;

  try {
    await guruwaliStore.addAnggota(idGuruWali, selectedSiswa.nis);
    toast.success("Anggota berhasil ditambahkan.");

    siswaQuery.value = "";
    selectedSiswa.nis = "";
    selectedSiswa.nama_siswa = "";
    selectedSiswa.nama_rombel = "";
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menambah anggota guru wali.",
    );
  } finally {
    addingAnggota.value = false;
  }
}

// =====================================
// Unggah Anggota via Template Excel
// =====================================

const showUploadDialog = ref(false);

async function downloadTemplate() {
  try {
    await guruwaliStore.downloadTemplateAnggota();
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal mengunduh template Excel.",
    );
  }
}

// =====================================
// Daftar Anggota
// =====================================

let anggotaSearchTimeout = null;

function searchAnggota() {
  clearTimeout(anggotaSearchTimeout);

  anggotaSearchTimeout = setTimeout(async () => {
    guruwaliStore.anggotaPagination.page = 1;
    await guruwaliStore.fetchAnggota(idGuruWali);
  }, 300);
}

async function changePage(page) {
  guruwaliStore.anggotaPagination.page = page;
  await guruwaliStore.fetchAnggota(idGuruWali);
}

async function handleExportAnggota() {
  if (exportingAnggota.value) return;

  exportingAnggota.value = true;

  try {
    await guruwaliStore.exportAnggotaByGuruWali(idGuruWali);
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal mengunduh data anggota.",
    );
  } finally {
    exportingAnggota.value = false;
  }
}

onMounted(async () => {
  await guruwaliStore.fetchById(idGuruWali);

  guruwaliStore.anggotaFilters.q = "";
  guruwaliStore.anggotaPagination.page = 1;
  await guruwaliStore.fetchAnggota(idGuruWali);

  kelasStore.pagination.limit = 100;
  rombelStore.pagination.limit = 500;
  kelasStore.fetchAll();
  rombelStore.fetchAll();
});
</script>
