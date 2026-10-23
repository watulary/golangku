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
        <div
          class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/8 text-[#041983] dark:bg-white/5 dark:text-indigo-300"
            >
              <ClipboardDocumentCheckIcon class="h-5 w-5" />
            </div>
            <div>
              <h2
                class="text-xl font-semibold text-gray-700 dark:text-slate-200"
              >
                Absensi Harian
              </h2>
              <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
                Catat absensi harian peserta didik per kelas/rombel.
              </p>
            </div>
          </div>
          <RouterLink
            :to="{ name: 'dashboard-absensi-rekap' }"
            title="Rekap Bulanan"
            class="inline-flex w-fit items-center gap-2 rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-[#041983] hover:bg-[#041983]/5 hover:text-[#041983] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-indigo-400 dark:hover:bg-white/5 dark:hover:text-indigo-300"
          >
            <TableCellsIcon class="h-5 w-5" />
          </RouterLink>
        </div>
        <div
          class="mt-5 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"
        >
          <div class="flex flex-wrap items-end gap-3">
            <div>
              <label
                class="mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400"
              >
                Tanggal
              </label>
              <input
                type="date"
                v-model="filters.tanggal"
                @change="tampilkan"
                class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition scheme-light hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200 dark:scheme-dark"
              />
            </div>
            <div>
              <label
                class="mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400"
              >
                Kelas
              </label>
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
              >
                Rombel
              </label>
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
          </div>
          <div class="flex items-center gap-2">
            <button
              v-if="hasRingkasan"
              type="button"
              @click="bukaRingkasan"
              title="Lihat Ringkasan Terakhir"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-blue-500 hover:bg-blue-50 hover:text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-blue-500 dark:hover:bg-blue-500/10 dark:hover:text-blue-400"
            >
              <ChartBarIcon class="h-5 w-5" />
            </button>
            <button
              v-if="absensiStore.daftarSiswa.length"
              type="button"
              @click="handleExport"
              :disabled="absensiStore.exportLoading"
              title="Unduh Laporan"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 disabled:opacity-50 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-emerald-500 dark:hover:bg-emerald-500/10 dark:hover:text-emerald-400"
            >
              <ArrowDownTrayIcon class="h-5 w-5" />
            </button>
            <button
              v-if="absensiStore.daftarSiswa.length"
              type="button"
              @click="setSemua('Hadir')"
              title="Tandai Semua Hadir"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-yellow-500 hover:bg-yellow-50 hover:text-yellow-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-yellow-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-yellow-500 dark:hover:bg-yellow-500/10 dark:hover:text-yellow-400"
            >
              <CheckCircleIcon class="h-5 w-5" />
            </button>
            <button
              v-if="absensiStore.daftarSiswa.length"
              type="button"
              @click="submit"
              :disabled="absensiStore.submitting"
              class="inline-flex items-center gap-2 rounded-xl bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:opacity-50"
            >
              <ArrowPathIcon
                v-if="absensiStore.submitting"
                class="h-4 w-4 animate-spin"
              />
              Simpan Absensi
            </button>
          </div>
        </div>
      </div>
      <!-- Filter Cepat Status -->
      <div
        v-if="absensiStore.daftarSiswa.length"
        class="flex flex-wrap items-center gap-2 border-b border-gray-100 px-6 py-3 dark:border-white/5"
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

        <div class="flex flex-wrap items-center gap-2">
          <button
            v-for="opt in statusFilterOptions"
            :key="opt.value"
            type="button"
            @click="statusFilter = opt.value"
            :class="[
              'inline-flex items-center gap-1.5 rounded-full border px-3 py-1 text-xs font-medium transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-1 focus-visible:ring-slate-400/40',
              statusFilter === opt.value
                ? opt.activeClass
                : 'border-slate-300 bg-white text-slate-600 hover:border-slate-400 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300',
            ]"
          >
            {{ opt.label }}
            <span
              :class="[
                'rounded-full bg-black/10 px-1.5 py-0.5 text-[10px] font-semibold',
                statusFilter === opt.value ? 'text-white' : opt.numberClass,
              ]"
            >
              {{ ringkasanSekarang[opt.key] }}
            </span>
          </button>
        </div>
      </div>
      <div
        v-if="!absensiStore.loading && !sudahDicoba"
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <CalendarDaysIcon class="h-10 w-10" />
          <p class="text-sm font-medium">
            Pilih tanggal, kelas/rombel lalu klik tampilkan.
          </p>
        </div>
      </div>

      <div
        v-else-if="
          !absensiStore.loading && absensiStore.daftarSiswa.length === 0
        "
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <InboxIcon class="h-10 w-10" />
          <p class="text-sm font-medium">
            Tidak ada data peserta didik untuk filter ini.
          </p>
        </div>
      </div>

      <div
        v-else-if="!absensiStore.loading && filteredDaftarSiswa.length === 0"
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <MagnifyingGlassIcon class="h-10 w-10" />
          <p class="text-sm font-medium">
            Tidak ada peserta didik dengan status "{{ statusFilter }}".
          </p>
        </div>
      </div>

      <div
        v-if="absensiStore.loading || filteredDaftarSiswa.length > 0"
        class="max-h-[70vh] overflow-auto"
      >
        <table class="w-full min-w-250">
          <thead class="sticky top-0 z-20 bg-gray-50/90 dark:bg-[#0B0F1A]">
            <tr>
              <th
                scope="col"
                class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                No
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
                Keterangan
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 dark:divide-white/5">
            <!-- Skeleton saat memuat -->
            <template v-if="absensiStore.loading">
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
              v-for="(s, idx) in filteredDaftarSiswa"
              :key="s.nis"
              class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
            >
              <td class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400">
                {{ idx + 1 }}
              </td>
              <td class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400">
                {{ s.nama_rombel }}
              </td>
              <td class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400">
                {{ s.nis }}
              </td>
              <td
                class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
              >
                {{ s.nama_siswa }}
              </td>
              <td class="px-4 py-2.5">
                <select
                  :value="s.status"
                  @change="
                    (e) => absensiStore.updateStatusLokal(s.nis, e.target.value)
                  "
                  :class="[
                    'rounded-xl border bg-white px-3 py-1.5 text-sm font-medium shadow-sm transition focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:bg-[#0B0F1A]',
                    statusSelectClass(s.status),
                  ]"
                >
                  <option value="Hadir">Hadir</option>
                  <option value="Sakit">Sakit</option>
                  <option value="Izin">Izin</option>
                  <option value="Alpa">Alpa</option>
                </select>
              </td>
              <td class="px-4 py-2.5">
                <input
                  type="text"
                  :value="s.keterangan"
                  @input="
                    (e) =>
                      absensiStore.updateKeteranganLokal(s.nis, e.target.value)
                  "
                  placeholder="Opsional"
                  class="w-full rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm text-slate-700 shadow-sm transition focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <!-- Modal Ringkasan Simpan -->
  <div
    v-if="showRingkasan"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
  >
    <div
      class="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl ring-1 ring-slate-900/5 dark:bg-slate-900 dark:ring-white/10"
    >
      <div class="flex items-center gap-3">
        <div
          class="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-green-100 dark:bg-green-500/10"
        >
          <CheckCircleIcon class="h-7 w-7 text-green-600 dark:text-green-400" />
        </div>
        <div>
          <h3 class="text-lg font-semibold text-gray-800 dark:text-slate-200">
            Absensi Tersimpan
          </h3>
          <p class="text-sm text-gray-500 dark:text-slate-400">
            Diproses:
            <span class="font-semibold text-gray-700 dark:text-slate-200">{{
              displayRingkasan.total
            }}</span>
            siswa
          </p>
        </div>
      </div>

      <div class="mt-5 grid grid-cols-2 gap-3">
        <button
          type="button"
          @click="filterDariRingkasan('Hadir')"
          class="rounded-xl border border-green-200 bg-green-50 p-4 text-center transition hover:border-green-400 hover:bg-green-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-green-500/40 dark:border-green-500/20 dark:bg-green-500/10 dark:hover:bg-green-500/20"
        >
          <p
            class="text-2xl font-bold tabular-nums text-green-700 dark:text-green-400"
          >
            {{ displayRingkasan.hadir }}
          </p>
          <p
            class="mt-1 text-xs font-medium text-green-700/80 dark:text-green-400/80"
          >
            Hadir
          </p>
        </button>
        <button
          type="button"
          @click="filterDariRingkasan('Sakit')"
          class="rounded-xl border border-yellow-200 bg-yellow-50 p-4 text-center transition hover:border-yellow-400 hover:bg-yellow-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-yellow-500/40 dark:border-yellow-500/20 dark:bg-yellow-500/10 dark:hover:bg-yellow-500/20"
        >
          <p
            class="text-2xl font-bold tabular-nums text-yellow-700 dark:text-yellow-400"
          >
            {{ displayRingkasan.sakit }}
          </p>
          <p
            class="mt-1 text-xs font-medium text-yellow-700/80 dark:text-yellow-400/80"
          >
            Sakit
          </p>
        </button>
        <button
          type="button"
          @click="filterDariRingkasan('Izin')"
          class="rounded-xl border border-blue-200 bg-blue-50 p-4 text-center transition hover:border-blue-400 hover:bg-blue-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-blue-500/20 dark:bg-blue-500/10 dark:hover:bg-blue-500/20"
        >
          <p
            class="text-2xl font-bold tabular-nums text-blue-700 dark:text-blue-400"
          >
            {{ displayRingkasan.izin }}
          </p>
          <p
            class="mt-1 text-xs font-medium text-blue-700/80 dark:text-blue-400/80"
          >
            Izin
          </p>
        </button>
        <button
          type="button"
          @click="filterDariRingkasan('Alpa')"
          class="rounded-xl border border-red-200 bg-red-50 p-4 text-center transition hover:border-red-400 hover:bg-red-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40 dark:border-red-500/20 dark:bg-red-500/10 dark:hover:bg-red-500/20"
        >
          <p
            class="text-2xl font-bold tabular-nums text-red-700 dark:text-red-400"
          >
            {{ displayRingkasan.alpa }}
          </p>
          <p
            class="mt-1 text-xs font-medium text-red-700/80 dark:text-red-400/80"
          >
            Alpa
          </p>
        </button>
      </div>

      <button
        type="button"
        @click="showRingkasan = false"
        class="mt-5 w-full rounded-xl bg-blue-600 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40"
      >
        Tutup
      </button>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from "vue";
import { useAbsensiStore } from "@/stores/absensi";
import { useKelasStore } from "@/stores/kelas";
import { useRombelStore } from "@/stores/rombel";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import { useToast } from "@/composables/useToast";
import { RouterLink } from "vue-router";
import {
  ArrowPathIcon,
  CheckCircleIcon,
  ArrowDownTrayIcon,
  TableCellsIcon,
  ChartBarIcon,
  MagnifyingGlassIcon,
  ClipboardDocumentCheckIcon,
  CalendarDaysIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const absensiStore = useAbsensiStore();
const kelasStore = useKelasStore();
const rombelStore = useRombelStore();
const toast = useToast();

// toISOString() selalu UTC — kalau dipakai langsung, jam 00:00-06:59 WIB
// akan salah tanggal (mundur 1 hari). Paksa eksplisit ke Asia/Jakarta di sini.
function tanggalHariIniWIB() {
  return new Intl.DateTimeFormat("en-CA", {
    timeZone: "Asia/Jakarta",
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).format(new Date());
}

const filters = reactive({
  tanggal: tanggalHariIniWIB(),
  id_kelas: 0,
  id_rombel: 0,
});

const sudahDicoba = ref(false);

const kelasOptions = computed(() =>
  (kelasStore.list || []).filter((k) => k.aktif === 1),
);
const rombelOptions = computed(() =>
  (rombelStore.list || []).filter((r) => r.aktif === 1),
);

//search + pencarian
const searchQuery = ref("");

const filteredDaftarSiswa = computed(() => {
  let list = absensiStore.daftarSiswa;
  if (statusFilter.value !== "Semua") {
    list = list.filter((s) => s.status === statusFilter.value);
  }
  const q = searchQuery.value.trim().toLowerCase();
  if (q) {
    list = list.filter(
      (s) =>
        s.nis.toLowerCase().includes(q) ||
        s.nama_siswa.toLowerCase().includes(q),
    );
  }
  return list;
});

function statusSelectClass(status) {
  const colors = {
    Hadir:
      "border-green-300 text-green-700 dark:border-green-600 dark:text-green-400",
    Sakit:
      "border-yellow-300 text-yellow-700 dark:border-yellow-600 dark:text-yellow-400",
    Izin: "border-blue-300 text-blue-700 dark:border-blue-600 dark:text-blue-400",
    Alpa: "border-red-300 text-red-700 dark:border-red-600 dark:text-red-400",
  };
  return (
    colors[status] ||
    "border-slate-300 text-slate-700 dark:border-slate-600 dark:text-slate-200"
  );
}

// ===== Ringkasan (live dari daftarSiswa yang sedang tampil, bukan snapshot submit) =====
const showRingkasan = ref(false);
const displayRingkasan = reactive({
  total: 0,
  hadir: 0,
  sakit: 0,
  izin: 0,
  alpa: 0,
});
const hasRingkasan = computed(() => absensiStore.daftarSiswa.length > 0);
const ringkasanSekarang = computed(() => {
  const saved = absensiStore.originalSiswa.filter((s) => s.id > 0);
  return {
    total: absensiStore.daftarSiswa.length,
    hadir: saved.filter((s) => s.status === "Hadir").length,
    sakit: saved.filter((s) => s.status === "Sakit").length,
    izin: saved.filter((s) => s.status === "Izin").length,
    alpa: saved.filter((s) => s.status === "Alpa").length,
  };
});

// ===== Filter cepat status =====
const statusFilter = ref("Semua");
const statusFilterOptions = [
  {
    value: "Semua",
    key: "total",
    label: "Semua",
    activeClass: "border-slate-600 bg-slate-600 text-white",
    numberClass: "text-slate-600 dark:text-slate-300",
  },
  {
    value: "Hadir",
    key: "hadir",
    label: "Hadir",
    activeClass: "border-green-600 bg-green-600 text-white",
    numberClass: "text-green-600 dark:text-green-400",
  },
  {
    value: "Sakit",
    key: "sakit",
    label: "Sakit",
    activeClass: "border-yellow-500 bg-yellow-500 text-white",
    numberClass: "text-yellow-600 dark:text-yellow-400",
  },
  {
    value: "Izin",
    key: "izin",
    label: "Izin",
    activeClass: "border-blue-600 bg-blue-600 text-white",
    numberClass: "text-blue-600 dark:text-blue-400",
  },
  {
    value: "Alpa",
    key: "alpa",
    label: "Alpa",
    activeClass: "border-red-600 bg-red-600 text-white",
    numberClass: "text-red-600 dark:text-red-400",
  },
];

function filterDariRingkasan(status) {
  statusFilter.value = status;
  showRingkasan.value = false;
}

function animateRingkasan() {
  const from = { ...displayRingkasan };
  const target = ringkasanSekarang.value;
  const duration = 700;
  const start = performance.now();
  function tick(now) {
    const progress = Math.min((now - start) / duration, 1);
    const eased = 1 - Math.pow(1 - progress, 3); // ease-out cubic
    Object.keys(target).forEach((key) => {
      displayRingkasan[key] = Math.round(
        from[key] + (target[key] - from[key]) * eased,
      );
    });
    if (progress < 1) requestAnimationFrame(tick);
  }
  requestAnimationFrame(tick);
}

function bukaRingkasan() {
  Object.assign(displayRingkasan, {
    total: 0,
    hadir: 0,
    sakit: 0,
    izin: 0,
    alpa: 0,
  });
  showRingkasan.value = true;
  animateRingkasan();
}

onMounted(async () => {
  kelasStore.pagination.limit = 100;
  rombelStore.pagination.limit = 500;
  await kelasStore.fetchAll();

  filters.id_kelas = 0; // Semua Kelas
  rombelStore.changeFilterKelas(null);
  await rombelStore.fetchAll();
  filters.id_rombel = 0; // Semua Rombel

  await tampilkan();
});

async function onKelasChange() {
  filters.id_rombel = 0; // balik ke "Semua Rombel" tiap ganti kelas manual
  rombelStore.changeFilterKelas(filters.id_kelas || null);
  await rombelStore.fetchAll();
  await tampilkan();
}

async function tampilkan() {
  sudahDicoba.value = true;
  statusFilter.value = "Semua";
  searchQuery.value = "";
  try {
    await absensiStore.fetchDaftar(
      filters.tanggal,
      filters.id_kelas,
      filters.id_rombel,
    );
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat daftar absensi.",
    );
  }
}

function setSemua(status) {
  absensiStore.daftarSiswa.forEach((s) =>
    absensiStore.updateStatusLokal(s.nis, status),
  );
  toast.success(`Semua ditandai ${status}.`);
}

async function submit() {
  if (!absensiStore.hasChanges()) {
    toast.warning("Tidak ada perubahan data.");
    return;
  }
  const items = absensiStore.daftarSiswa.map((s) => ({
    nis: s.nis,
    status: s.status,
    keterangan: s.keterangan,
  }));
  try {
    await absensiStore.submitBulk(
      filters.tanggal,
      filters.id_kelas,
      filters.id_rombel,
      items,
    );
    // refetch supaya id & status ter-update dari DB — dipakai buat hitung chip "sudah tersimpan"
    await absensiStore.fetchDaftar(
      filters.tanggal,
      filters.id_kelas,
      filters.id_rombel,
    );
    bukaRingkasan();
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menyimpan absensi.");
  }
}

async function handleExport() {
  try {
    const status = statusFilter.value !== "Semua" ? statusFilter.value : "";
    await absensiStore.exportAbsensi(
      filters.tanggal,
      filters.id_kelas,
      filters.id_rombel,
      status,
    );
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal mengunduh laporan absensi.",
    );
  }
}
</script>
