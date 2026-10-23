<template>
  <DashboardPageHeader />
  <div class="space-y-6">
    <div
      class="relative overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
    >
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>
      <div class="border-b border-gray-100 px-6 py-5 dark:border-white/5">
        <div class="flex items-center gap-3">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/[0.08] text-[#041983] dark:bg-white/5 dark:text-indigo-300"
          >
            <ClipboardDocumentCheckIcon class="h-5 w-5" />
          </div>
          <div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Absensi
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
              Riwayat kehadiran per bulan.
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
        </div>
      </div>

      <div v-if="absensiStore.siswaLoading">
        <div class="grid grid-cols-2 gap-3 px-6 py-5 sm:grid-cols-4">
          <div
            v-for="n in 4"
            :key="n"
            class="rounded-xl bg-gray-100 p-4 dark:bg-white/5"
          >
            <div
              class="mx-auto h-7 w-10 animate-pulse rounded bg-gray-200 dark:bg-white/10"
            ></div>
            <div
              class="mx-auto mt-2 h-3 w-12 animate-pulse rounded bg-gray-200 dark:bg-white/10"
            ></div>
          </div>
        </div>
        <div class="space-y-2 px-6 pb-6">
          <div
            v-for="n in 4"
            :key="n"
            class="h-10 animate-pulse rounded-lg bg-gray-100 dark:bg-white/5"
          ></div>
        </div>
      </div>

      <template v-else>
        <!-- Ringkasan -->
        <div class="grid grid-cols-2 gap-3 px-6 py-5 sm:grid-cols-4">
          <div
            class="rounded-xl border border-green-200 bg-green-50 p-4 text-center dark:border-green-500/20 dark:bg-green-500/10"
          >
            <p
              class="text-2xl font-bold tabular-nums text-green-700 dark:text-green-400"
            >
              {{ ringkasan.hadir }}
            </p>
            <p
              class="mt-1 text-xs font-medium text-green-700/80 dark:text-green-400/80"
            >
              Hadir
            </p>
          </div>
          <div
            class="rounded-xl border border-yellow-200 bg-yellow-50 p-4 text-center dark:border-yellow-500/20 dark:bg-yellow-500/10"
          >
            <p
              class="text-2xl font-bold tabular-nums text-yellow-700 dark:text-yellow-400"
            >
              {{ ringkasan.sakit }}
            </p>
            <p
              class="mt-1 text-xs font-medium text-yellow-700/80 dark:text-yellow-400/80"
            >
              Sakit
            </p>
          </div>
          <div
            class="rounded-xl border border-blue-200 bg-blue-50 p-4 text-center dark:border-blue-500/20 dark:bg-blue-500/10"
          >
            <p
              class="text-2xl font-bold tabular-nums text-blue-700 dark:text-blue-400"
            >
              {{ ringkasan.izin }}
            </p>
            <p
              class="mt-1 text-xs font-medium text-blue-700/80 dark:text-blue-400/80"
            >
              Izin
            </p>
          </div>
          <div
            class="rounded-xl border border-red-200 bg-red-50 p-4 text-center dark:border-red-500/20 dark:bg-red-500/10"
          >
            <p
              class="text-2xl font-bold tabular-nums text-red-700 dark:text-red-400"
            >
              {{ ringkasan.alpa }}
            </p>
            <p
              class="mt-1 text-xs font-medium text-red-700/80 dark:text-red-400/80"
            >
              Alpa
            </p>
          </div>
        </div>

        <!-- Daftar per tanggal -->
        <div
          v-if="absensiStore.siswaList.length === 0"
          class="py-16 text-center"
        >
          <div
            class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
          >
            <InboxIcon class="h-10 w-10" />
            <p class="text-sm font-medium">Belum ada data absensi bulan ini.</p>
          </div>
        </div>
        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead class="bg-gray-50 dark:bg-[#0B0F1A]">
              <tr>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Tanggal
                </th>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Status
                </th>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Keterangan
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-white/5">
              <tr
                v-for="item in absensiStore.siswaList"
                :key="item.tanggal"
                class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
              >
                <td
                  class="px-6 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ formatTanggal(item.tanggal) }}
                </td>
                <td class="px-6 py-2.5">
                  <span
                    :class="[
                      'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
                      statusBadgeClass(item.status),
                    ]"
                  >
                    {{ item.status }}
                  </span>
                </td>
                <td
                  class="px-6 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.keterangan || "-" }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed, onMounted } from "vue";
import { useAbsensiStore } from "@/stores/absensi";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import { useToast } from "@/composables/useToast";
import {
  ClipboardDocumentCheckIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";

const absensiStore = useAbsensiStore();
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
});

const ringkasan = computed(() => {
  const list = absensiStore.siswaList;
  return {
    hadir: list.filter((i) => i.status === "Hadir").length,
    sakit: list.filter((i) => i.status === "Sakit").length,
    izin: list.filter((i) => i.status === "Izin").length,
    alpa: list.filter((i) => i.status === "Alpa").length,
  };
});

function statusBadgeClass(status) {
  const colors = {
    Hadir:
      "bg-green-50 text-green-700 dark:bg-green-500/10 dark:text-green-400",
    Sakit:
      "bg-yellow-50 text-yellow-700 dark:bg-yellow-500/10 dark:text-yellow-400",
    Izin: "bg-blue-50 text-blue-700 dark:bg-blue-500/10 dark:text-blue-400",
    Alpa: "bg-red-50 text-red-700 dark:bg-red-500/10 dark:text-red-400",
  };
  return (
    colors[status] ||
    "bg-gray-100 text-gray-600 dark:bg-slate-800 dark:text-slate-400"
  );
}

function formatTanggal(tanggal) {
  const d = new Date(tanggal);
  return d.toLocaleDateString("id-ID", {
    weekday: "long",
    day: "numeric",
    month: "long",
    year: "numeric",
  });
}

async function tampilkan() {
  try {
    await absensiStore.fetchAbsensiSaya(filters.bulan, filters.tahun);
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal memuat absensi.");
  }
}

onMounted(tampilkan);
</script>
