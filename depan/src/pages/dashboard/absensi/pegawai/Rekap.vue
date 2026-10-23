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
        <div class="flex items-center gap-3">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/8 text-[#041983] dark:bg-white/5 dark:text-indigo-300"
          >
            <TableCellsIcon class="h-5 w-5" />
          </div>
          <div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Rekap Bulanan Absensi
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
              Matriks kehadiran peserta didik per bulan.
            </p>
          </div>
        </div>
        <div class="mt-5 flex flex-wrap items-end gap-3">
          <div>
            <label
              class="mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400"
              >Bulan</label
            >
            <select
              v-model.number="filters.bulan"
              @change="tampilkan"
              class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
            >
              <option
                v-for="(nama, idx) in namaBulan"
                :key="idx"
                :value="idx + 1"
              >
                {{ nama }}
              </option>
            </select>
          </div>
          <div>
            <label
              class="mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400"
              >Tahun</label
            >
            <input
              type="number"
              v-model.number="filters.tahun"
              @change="tampilkan"
              class="w-24 rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
            />
          </div>
          <div>
            <label
              class="mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400"
              >Kelas</label
            >
            <select
              v-model="filters.id_kelas"
              @change="onKelasChange"
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
          </div>
          <div>
            <label
              class="mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400"
              >Rombel</label
            >
            <select
              v-model="filters.id_rombel"
              @change="tampilkan"
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
          <button
            type="button"
            @click="handleExport"
            :disabled="
              !absensiStore.rekapSiswa.length || absensiStore.rekapExportLoading
            "
            title="Unduh Rekap"
            class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 disabled:cursor-not-allowed disabled:opacity-50 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-emerald-500 dark:hover:bg-emerald-500/10 dark:hover:text-emerald-400"
          >
            <ArrowDownTrayIcon class="h-5 w-5" />
          </button>
        </div>
      </div>

      <div
        v-if="absensiStore.rekapLoading"
        class="flex flex-col items-center justify-center py-20"
      >
        <div
          class="h-10 w-10 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"
        />
        <p class="mt-4 text-sm text-gray-500 dark:text-slate-400">
          Memuat data...
        </p>
      </div>

      <div
        v-else-if="absensiStore.rekapSiswa.length === 0"
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <InboxIcon class="h-10 w-10" />
          <p class="text-sm font-medium">
            Tidak ada data siswa untuk filter ini.
          </p>
        </div>
      </div>

      <div v-else class="max-h-[70vh] overflow-auto">
        <table class="w-full border-collapse text-sm">
          <thead class="sticky top-0 z-20 bg-gray-50 dark:bg-[#0B0F1A]">
            <tr>
              <th
                scope="col"
                class="px-3 py-2 text-left text-xs font-semibold text-gray-500 dark:text-slate-400"
              >
                No
              </th>
              <th
                scope="col"
                class="px-3 py-2 text-left text-xs font-semibold text-gray-500 dark:text-slate-400"
              >
                Rombel
              </th>
              <th
                scope="col"
                class="sticky left-0 z-30 w-40 bg-gray-50 px-3 py-2 text-left text-xs font-semibold text-gray-500 dark:bg-[#0B0F1A] dark:text-slate-400"
              >
                Nama Peserta Didik
              </th>
              <th
                v-for="d in absensiStore.rekapJumlahHari"
                :key="d"
                scope="col"
                :class="[
                  'w-8 px-1 py-2 text-center text-xs font-semibold text-gray-500 dark:text-slate-400',
                  isAkhirPekan(d) ? 'bg-gray-100 dark:bg-white/5' : '',
                ]"
              >
                {{ d }}
              </th>
              <th
                scope="col"
                class="w-12 border-l border-gray-200 px-2 py-2 text-center text-xs font-semibold text-green-600 dark:border-white/10 dark:text-green-400"
              >
                H
              </th>
              <th
                scope="col"
                class="w-12 px-2 py-2 text-center text-xs font-semibold text-yellow-600 dark:text-yellow-400"
              >
                S
              </th>
              <th
                scope="col"
                class="w-12 px-2 py-2 text-center text-xs font-semibold text-blue-600 dark:text-blue-400"
              >
                I
              </th>
              <th
                scope="col"
                class="w-12 px-2 py-2 text-center text-xs font-semibold text-red-600 dark:text-red-400"
              >
                A
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 dark:divide-white/5">
            <tr
              v-for="(s, idx) in absensiStore.rekapSiswa"
              :key="s.nis"
              class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
            >
              <td class="px-3 py-2 text-xs text-gray-600 dark:text-slate-400">
                {{ idx + 1 }}
              </td>
              <td class="px-3 py-2 text-xs text-gray-600 dark:text-slate-400">
                {{ s.nama_rombel }}
              </td>
              <td
                class="sticky left-0 z-10 w-40 truncate bg-white px-3 py-2 text-xs font-medium text-gray-800 dark:bg-slate-900 dark:text-slate-200"
              >
                {{ s.nama_siswa }}
              </td>
              <td
                v-for="(kode, d) in s.harian"
                :key="d"
                :class="[
                  'w-8 px-1 py-2 text-center text-xs font-semibold',
                  kodeColorClass(kode),
                  isAkhirPekan(d + 1) ? 'bg-gray-50 dark:bg-white/3' : '',
                ]"
              >
                {{ kode || "-" }}
              </td>
              <td
                class="border-l border-gray-100 px-2 py-2 text-center text-xs font-semibold text-green-600 dark:border-white/5 dark:text-green-400"
              >
                {{ hitungKode(s.harian, "H") }}
              </td>
              <td
                class="px-2 py-2 text-center text-xs font-semibold text-yellow-600 dark:text-yellow-400"
              >
                {{ hitungKode(s.harian, "S") }}
              </td>
              <td
                class="px-2 py-2 text-center text-xs font-semibold text-blue-600 dark:text-blue-400"
              >
                {{ hitungKode(s.harian, "I") }}
              </td>
              <td
                class="px-2 py-2 text-center text-xs font-semibold text-red-600 dark:text-red-400"
              >
                {{ hitungKode(s.harian, "A") }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed, onMounted } from "vue";
import { useAbsensiStore } from "@/stores/absensi";
import { useKelasStore } from "@/stores/kelas";
import { useRombelStore } from "@/stores/rombel";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import { useToast } from "@/composables/useToast";
import {
  ArrowDownTrayIcon,
  TableCellsIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const absensiStore = useAbsensiStore();
const kelasStore = useKelasStore();
const rombelStore = useRombelStore();
const toast = useToast();

const namaBulan = [
  "Januari",
  "Februari",
  "Maret",
  "April",
  "Mei",
  "Juni",
  "Juli",
  "Agustus",
  "September",
  "Oktober",
  "November",
  "Desember",
];

const now = new Date();
const filters = reactive({
  bulan: now.getMonth() + 1,
  tahun: now.getFullYear(),
  id_kelas: 0,
  id_rombel: 0,
});

const kelasOptions = computed(() =>
  (kelasStore.list || []).filter((k) => k.aktif === 1),
);
const rombelOptions = computed(() =>
  (rombelStore.list || []).filter((r) => r.aktif === 1),
);

function kodeColorClass(kode) {
  const colors = {
    H: "text-green-600 dark:text-green-400",
    S: "text-yellow-600 dark:text-yellow-400",
    I: "text-blue-600 dark:text-blue-400",
    A: "text-red-600 dark:text-red-400",
  };
  return colors[kode] || "text-gray-300 dark:text-slate-700";
}

function hitungKode(harian, kode) {
  return harian.filter((k) => k === kode).length;
}

function isAkhirPekan(hari) {
  const tanggal = new Date(filters.tahun, filters.bulan - 1, hari);
  const hariKe = tanggal.getDay(); // 0 = Minggu, 6 = Sabtu
  return hariKe === 0 || hariKe === 6;
}

onMounted(async () => {
  kelasStore.pagination.limit = 100;
  rombelStore.pagination.limit = 500;
  await kelasStore.fetchAll();
  filters.id_kelas = 0;
  rombelStore.changeFilterKelas(null);
  await rombelStore.fetchAll();
  filters.id_rombel = 0;
  await tampilkan();
});

async function onKelasChange() {
  filters.id_rombel = 0;
  rombelStore.changeFilterKelas(filters.id_kelas || null);
  await rombelStore.fetchAll();
  await tampilkan();
}

async function tampilkan() {
  try {
    await absensiStore.fetchRekapBulanan(
      filters.bulan,
      filters.tahun,
      filters.id_kelas,
      filters.id_rombel,
    );
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal memuat rekap bulanan.");
  }
}

async function handleExport() {
  try {
    await absensiStore.exportRekapBulanan(
      filters.bulan,
      filters.tahun,
      filters.id_kelas,
      filters.id_rombel,
    );
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal mengunduh rekap bulanan.",
    );
  }
}
</script>
