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
        <div class="flex items-center gap-3">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/8 text-[#041983] dark:bg-white/5 dark:text-indigo-300"
          >
            <UserGroupIcon class="h-5 w-5" />
          </div>
          <div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Keanggotaan Ekstrakurikuler
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
              Catat siapa yang menjadi anggota tiap ekstrakurikuler pada tahun
              pelajaran &amp; semester tertentu (riwayat per periode), terpisah
              dari daftar keanggotaan berjalan.
            </p>
          </div>
        </div>

        <div
          class="mt-5 flex flex-col gap-3 xl:flex-row xl:items-end xl:justify-between"
        >
          <div class="flex flex-wrap items-end gap-3">
            <div>
              <label :class="labelClass">Tahun Pelajaran</label>
              <select
                v-model="filters.id_tahunpelajaran"
                @change="onFilterChange"
                :class="selectClass"
              >
                <option :value="0" disabled>Pilih</option>
                <option
                  v-for="t in tahunOptions"
                  :key="t.id_tahunpelajaran"
                  :value="t.id_tahunpelajaran"
                >
                  {{ t.nama_tahunpelajaran }}
                </option>
              </select>
            </div>
            <div>
              <label :class="labelClass">Semester</label>
              <select
                v-model="filters.semester"
                @change="onFilterChange"
                :class="selectClass"
              >
                <option value="Ganjil">Ganjil</option>
                <option value="Genap">Genap</option>
              </select>
            </div>
            <div>
              <label :class="labelClass">Kelas</label>
              <select
                v-model="filters.id_kelas"
                @change="onKelasChange"
                :class="selectClass"
              >
                <option :value="0" disabled>Pilih</option>
                <option
                  v-for="k in kelasOptions"
                  :key="k.id_kelas"
                  :value="k.id_kelas"
                >
                  {{ k.nama_kelas }}
                </option>
              </select>
            </div>
            <div>
              <label :class="labelClass">Rombel</label>
              <select
                v-model="filters.id_rombel"
                @change="onFilterChange"
                :class="selectClass"
              >
                <option :value="0" disabled>Pilih</option>
                <option
                  v-for="r in rombelOptions"
                  :key="r.id_rombel"
                  :value="r.id_rombel"
                >
                  {{ r.nama_rombel }}
                </option>
              </select>
            </div>
            <div>
              <label :class="labelClass">Ekstrakurikuler</label>
              <select
                v-model="filters.id_ekskul"
                @change="onFilterChange"
                :class="[selectClass, 'max-w-64']"
              >
                <option :value="0" disabled>Pilih</option>
                <option
                  v-for="e in keanggotaanStore.ekskulOptions"
                  :key="e.id_ekskul"
                  :value="e.id_ekskul"
                >
                  {{ e.nama_ekskul }}
                </option>
              </select>
            </div>
          </div>

          <div v-if="filterLengkap" class="flex items-center gap-2">
            <button
              type="button"
              @click="toggleSemua"
              class="rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm font-medium text-slate-600 shadow-sm transition hover:border-blue-500 hover:text-blue-600 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300"
            >
              {{ semuaTercentang ? "Batalkan Semua" : "Pilih Semua" }}
            </button>
            <button
              type="button"
              title="Unggah Keanggotaan Ekskul"
              @click="showUploadDialog = true"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-yellow-500 hover:bg-yellow-50 hover:text-yellow-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-yellow-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-yellow-500 dark:hover:bg-yellow-500/10 dark:hover:text-yellow-300"
            >
              <ArrowUpTrayIcon class="h-5 w-5" />
            </button>
            <button
              v-if="keanggotaanStore.daftar.length"
              type="button"
              @click="submit"
              :disabled="keanggotaanStore.submitting"
              class="inline-flex items-center gap-2 rounded-xl bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:opacity-50"
            >
              <ArrowPathIcon
                v-if="keanggotaanStore.submitting"
                class="h-4 w-4 animate-spin"
              />
              Simpan Keanggotaan
            </button>
          </div>
        </div>
      </div>

      <!-- Pencarian -->
      <div
        v-if="keanggotaanStore.daftar.length"
        class="flex flex-wrap items-center justify-between gap-2 border-b border-gray-100 px-6 py-3 dark:border-white/5"
      >
        <div class="relative w-full max-w-62.5">
          <MagnifyingGlassIcon
            class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400"
          />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Pencarian..."
            class="w-full rounded-xl border border-slate-300 bg-white py-1.5 pl-9 pr-3 text-sm text-slate-700 shadow-sm transition focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500"
          />
        </div>
        <p class="text-xs text-slate-500 dark:text-slate-400">
          Anggota
          <span class="font-semibold">{{ jumlahAnggota }}</span> dari
          {{ keanggotaanStore.daftar.length }} siswa rombel ini.
        </p>
      </div>

      <!-- State kosong -->
      <div v-if="!filterLengkap" class="py-16 text-center">
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <UserGroupIcon class="h-10 w-10" />
          <p class="text-sm font-medium">
            Pilih tahun pelajaran, semester, rombel, dan ekstrakurikuler.
          </p>
        </div>
      </div>

      <div
        v-else-if="
          !keanggotaanStore.loading && keanggotaanStore.daftar.length === 0
        "
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <InboxIcon class="h-10 w-10" />
          <p class="text-sm font-medium">Tidak ada siswa di rombel tersebut.</p>
        </div>
      </div>

      <div
        v-else-if="!keanggotaanStore.loading && filteredDaftar.length === 0"
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <MagnifyingGlassIcon class="h-10 w-10" />
          <p class="text-sm font-medium">Tidak ada hasil pencarian.</p>
        </div>
      </div>

      <div
        v-if="
          filterLengkap && (keanggotaanStore.loading || filteredDaftar.length)
        "
        class="max-h-[70vh] overflow-auto"
      >
        <table class="w-full min-w-125">
          <thead class="sticky top-0 z-20 bg-gray-50/90 dark:bg-[#0B0F1A]">
            <tr>
              <th :class="thClass">No</th>
              <th :class="thClass">NIS</th>
              <th :class="thClass">Nama Peserta Didik</th>
              <th :class="[thClass, 'w-32 text-center']">Anggota</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 dark:divide-white/5">
            <template v-if="keanggotaanStore.loading">
              <tr v-for="n in 6" :key="`skeleton-${n}`">
                <td v-for="c in 4" :key="c" class="px-4 py-3">
                  <div
                    class="h-4 w-20 animate-pulse rounded bg-gray-100 dark:bg-white/5"
                  ></div>
                </td>
              </tr>
            </template>

            <template v-else>
              <tr
                v-for="(s, idx) in filteredDaftar"
                :key="s.nis"
                class="align-top transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
              >
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ idx + 1 }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ s.nis }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
                >
                  {{ s.nama_siswa }}
                </td>
                <td class="px-4 py-2.5 text-center">
                  <input
                    v-model="s.anggota"
                    type="checkbox"
                    class="h-4 w-4 rounded border-slate-300 text-blue-600 focus:ring-blue-500"
                  />
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <DashboardUploadDialog
    v-model="showUploadDialog"
    title="Unggah Keanggotaan Ekskul"
    subtitle="Unggah status keanggotaan sesuai filter yang dipilih menggunakan Template (.xlsx)"
    :upload-action="(file) => keanggotaanStore.importExcel(file, params())"
    @download-template="handleTemplate"
    @success="handleImportSuccess"
  />
</template>

<script setup>
import { reactive, ref, computed, onMounted } from "vue";
import { onBeforeRouteLeave } from "vue-router";
import { useKeanggotaanEkskulStore } from "@/stores/keanggotaanekskul";
import { useKelasStore } from "@/stores/kelas";
import { useRombelStore } from "@/stores/rombel";
import { useTahunPelajaranStore } from "@/stores/tahunpelajaran";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardUploadDialog from "@/components/dashboard/DashboardUploadDialog.vue";
import { useToast } from "@/composables/useToast";
import {
  UserGroupIcon,
  ArrowUpTrayIcon,
  ArrowPathIcon,
  MagnifyingGlassIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const keanggotaanStore = useKeanggotaanEkskulStore();
const kelasStore = useKelasStore();
const rombelStore = useRombelStore();
const tahunStore = useTahunPelajaranStore();
const toast = useToast();

const labelClass =
  "mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400";
const selectClass =
  "rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200";
const thClass =
  "px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400";

const filters = reactive({
  id_tahunpelajaran: 0,
  semester: "Ganjil",
  id_kelas: 0,
  id_rombel: 0,
  id_ekskul: 0,
});

const showUploadDialog = ref(false);
const searchQuery = ref("");

const tahunOptions = computed(() => tahunStore.list || []);
const kelasOptions = computed(() =>
  (kelasStore.list || []).filter((k) => k.aktif === 1),
);
const rombelOptions = computed(() =>
  (rombelStore.list || []).filter((r) => r.aktif === 1),
);

const filterLengkap = computed(
  () =>
    filters.id_tahunpelajaran > 0 &&
    filters.id_rombel > 0 &&
    filters.id_ekskul > 0,
);

const filteredDaftar = computed(() => {
  const q = searchQuery.value.trim().toLowerCase();
  if (!q) return keanggotaanStore.daftar;
  return keanggotaanStore.daftar.filter(
    (s) =>
      s.nis.toLowerCase().includes(q) || s.nama_siswa.toLowerCase().includes(q),
  );
});

const jumlahAnggota = computed(
  () => keanggotaanStore.daftar.filter((s) => s.anggota).length,
);

const semuaTercentang = computed(
  () =>
    keanggotaanStore.daftar.length > 0 &&
    keanggotaanStore.daftar.every((s) => s.anggota),
);

function toggleSemua() {
  const nilai = !semuaTercentang.value;
  keanggotaanStore.daftar.forEach((s) => {
    s.anggota = nilai;
  });
}

function params() {
  return {
    id_tahunpelajaran: filters.id_tahunpelajaran,
    semester: filters.semester,
    id_rombel: filters.id_rombel,
    id_ekskul: filters.id_ekskul,
  };
}

onMounted(async () => {
  tahunStore.pagination.limit = 100;
  tahunStore.sortBy = "nama";
  tahunStore.sortDir = "desc";
  kelasStore.pagination.limit = 100;
  rombelStore.pagination.limit = 500;

  await Promise.all([
    tahunStore.fetchAll(),
    kelasStore.fetchAll(),
    keanggotaanStore.fetchEkskul(),
  ]);

  const aktif = (tahunStore.list || []).find((t) => t.aktif === 1);
  if (aktif) {
    filters.id_tahunpelajaran = aktif.id_tahunpelajaran;
    filters.semester = aktif.semester || "Ganjil";
  }

  keanggotaanStore.clear();
});

async function onKelasChange() {
  filters.id_rombel = 0;
  keanggotaanStore.clear();
  rombelStore.changeFilterKelas(filters.id_kelas || null);
  await rombelStore.fetchAll();
}

async function onFilterChange() {
  if (!confirmBuang()) return;
  await tampilkan();
}

// Filter berubah saat ada isian belum disimpan -> konfirmasi dulu.
function confirmBuang() {
  if (keanggotaanStore.daftar.length && keanggotaanStore.hasChanges()) {
    return window.confirm(
      "Ada perubahan yang belum disimpan. Lanjutkan dan buang perubahan?",
    );
  }
  return true;
}

onBeforeRouteLeave(() => {
  if (keanggotaanStore.daftar.length && keanggotaanStore.hasChanges()) {
    return window.confirm(
      "Ada perubahan yang belum disimpan. Tinggalkan halaman ini?",
    );
  }
});

async function tampilkan() {
  searchQuery.value = "";
  if (!filterLengkap.value) {
    keanggotaanStore.clear();
    return;
  }
  try {
    await keanggotaanStore.fetchDaftar(params());
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat data keanggotaan.",
    );
  }
}

async function submit() {
  if (!keanggotaanStore.hasChanges()) {
    toast.warning("Tidak ada perubahan data.");
    return;
  }
  try {
    const res = await keanggotaanStore.submit(params());
    toast.success(res.message || "Keanggotaan berhasil disimpan.");
    await keanggotaanStore.fetchDaftar(params());
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menyimpan keanggotaan.",
    );
  }
}

async function handleTemplate() {
  try {
    await keanggotaanStore.downloadTemplate(params());
  } catch (error) {
    toast.error("Gagal mengunduh template.");
  }
}

async function handleImportSuccess(message) {
  toast.success(message || "Impor keanggotaan berhasil.");
  await tampilkan();
}
</script>
