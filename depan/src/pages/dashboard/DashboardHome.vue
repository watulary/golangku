<template>
  <DashboardPageHeader />
  <!-- Hero banner -->
  <div
    class="relative mb-8 overflow-hidden rounded-2xl bg-white p-6 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5 sm:p-8"
  >
    <div
      class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
    ></div>
    <div>
      <p class="text-sm font-medium text-gray-500 dark:text-slate-400">
        {{ greeting }}
      </p>
      <h1
        class="mt-1 text-2xl font-bold tracking-tight text-gray-700 dark:text-slate-200"
      >
        Statistik Tahun Pelajaran {{ namaTahunPelajaranAktif }}
      </h1>
      <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
        Statistik dan rekapitulasi data sekolah secara ringkas.
      </p>
    </div>
  </div>
  <div class="space-y-8">
    <!-- Statistik Akademik -->
    <div>
      <p
        class="mb-3 text-xs font-semibold tracking-wider text-gray-400 dark:text-slate-500"
      >
        Statistik Akademik
      </p>
      <div v-if="loading" class="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <div
          v-for="n in 3"
          :key="n"
          class="rounded-2xl bg-white p-5 ring-1 ring-slate-900/5 dark:bg-slate-900 dark:ring-white/5"
        >
          <div
            class="h-3 w-20 animate-pulse rounded bg-gray-100 dark:bg-white/5"
          ></div>
          <div class="mt-3 flex items-center gap-3">
            <div
              class="h-16 w-16 shrink-0 animate-pulse rounded-full bg-gray-100 dark:bg-white/5"
            ></div>
            <div class="min-w-0 flex-1 space-y-2">
              <div
                class="h-3 w-24 animate-pulse rounded bg-gray-100 dark:bg-white/5"
              ></div>
              <div
                class="h-3 w-20 animate-pulse rounded bg-gray-100 dark:bg-white/5"
              ></div>
            </div>
          </div>
        </div>
      </div>
      <TransitionGroup
        v-else
        tag="div"
        name="stat-card"
        class="grid grid-cols-1 gap-4 sm:grid-cols-3"
      >
        <!-- Rombel: total + mini bar per tingkat -->
        <RouterLink
          key="rombel"
          to="/dashboard/rombel"
          class="group rounded-2xl bg-white p-5 ring-1 ring-slate-900/5 shadow-sm transition-all duration-300 hover:-translate-y-0.5 hover:shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:bg-slate-900 dark:ring-white/5"
        >
          <p
            class="text-xs font-semibold tracking-wide text-gray-400 dark:text-slate-500"
          >
            Rombel
          </p>
          <h2
            class="mt-1 text-3xl font-bold leading-tight tracking-tight text-gray-700 tabular-nums dark:text-slate-200"
          >
            {{ jumlahRombelTotal }}
          </h2>
          <div class="mt-3 flex items-end gap-2.5" style="height: 44px">
            <div
              v-for="(t, i) in perTingkatStats"
              :key="t.nama_kelas"
              class="flex flex-col items-center gap-1"
            >
              <div
                class="w-5 rounded-t"
                :style="{
                  height: barHeight(t.jumlah_rombel) + 'px',
                  backgroundColor: rombelBarColor(i),
                }"
              ></div>
              <span class="text-[11px] text-gray-400 dark:text-slate-200">
                Kelas {{ t.nama_kelas }}
              </span>
            </div>
            <span
              class="ml-1 text-[11px] tabular-nums text-gray-500 dark:text-slate-200"
            >
              {{ perTingkatStats.map((t) => t.jumlah_rombel).join(" · ") }}
              rombel
            </span>
          </div>
        </RouterLink>
        <!-- Peserta Didik: bar per tingkat (satu baris) -->
        <RouterLink
          key="peserta-didik"
          to="/dashboard/peserta-didik"
          class="group rounded-2xl bg-white p-5 ring-1 ring-slate-900/5 shadow-sm transition-all duration-300 hover:-translate-y-0.5 hover:shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:bg-slate-900 dark:ring-white/5"
        >
          <p
            class="text-xs font-semibold tracking-wide text-gray-400 dark:text-slate-500"
          >
            Peserta Didik
          </p>
          <h2
            class="mt-1 text-3xl font-bold leading-tight tracking-tight text-gray-700 tabular-nums dark:text-slate-200"
          >
            {{ jumlahMuridTotal }}
          </h2>
          <div
            class="mt-3 flex h-3 w-full overflow-hidden rounded-full bg-gray-100 dark:bg-white/10"
          >
            <div
              v-for="(t, i) in perTingkatStats"
              :key="t.nama_kelas"
              class="h-full first:rounded-l-full last:rounded-r-full"
              :style="{
                width: muridPct(t.jumlah_murid) + '%',
                backgroundColor: tingkatColor(i),
              }"
            ></div>
          </div>
          <div
            class="mt-2.5 flex flex-wrap gap-x-3 gap-y-1 text-[11px] text-gray-500 dark:text-slate-200"
          >
            <span
              v-for="(t, i) in perTingkatStats"
              :key="t.nama_kelas"
              class="flex items-center gap-1.5"
            >
              <span
                class="h-2 w-2 shrink-0 rounded-sm"
                :style="{ backgroundColor: tingkatColor(i) }"
              ></span>
              <span class="tabular-nums">
                Kelas {{ t.nama_kelas }} · {{ t.jumlah_murid }}
              </span>
            </span>
          </div>
        </RouterLink>
        <!-- Pegawai: radial ring bersarang (Guru & Tendik), agak besar -->
        <RouterLink
          key="tenaga-pendidik"
          to="/dashboard/pegawai"
          class="group rounded-2xl bg-white p-5 ring-1 ring-slate-900/5 shadow-sm transition-all duration-300 hover:-translate-y-0.5 hover:shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:bg-slate-900 dark:ring-white/5"
        >
          <p
            class="text-xs font-semibold tracking-wide text-gray-400 dark:text-slate-500"
          >
            Pegawai
          </p>
          <div class="mt-2 flex items-center gap-3">
            <svg
              width="76"
              height="76"
              viewBox="0 0 76 76"
              class="shrink-0 -rotate-90"
            >
              <circle
                v-for="ring in tenagaPendidikRings"
                :key="`${ring.label}-track`"
                cx="38"
                cy="38"
                :r="ring.r"
                fill="none"
                stroke-width="7"
                class="stroke-gray-100 dark:stroke-white/10"
              />
              <circle
                v-for="ring in tenagaPendidikRings"
                :key="ring.label"
                cx="38"
                cy="38"
                :r="ring.r"
                fill="none"
                :stroke="ring.color"
                stroke-width="7"
                stroke-linecap="round"
                :stroke-dasharray="`${ring.dash} ${ring.circumference}`"
              />
            </svg>
            <div
              class="flex flex-col gap-1 text-xs text-gray-500 dark:text-slate-200"
            >
              <span
                v-for="ring in tenagaPendidikRings"
                :key="`${ring.label}-legend`"
                class="flex items-center gap-1.5"
              >
                <span
                  class="h-2 w-2 shrink-0 rounded-sm"
                  :style="{ backgroundColor: ring.color }"
                ></span>
                <span class="tabular-nums">
                  {{ ring.label }} · {{ ring.value }}
                </span>
              </span>
            </div>
          </div>
        </RouterLink>
      </TransitionGroup>
    </div>
    <!-- Statistik Konten Website -->
    <div>
      <p
        class="mb-3 text-xs font-semibold tracking-wider text-gray-400 dark:text-slate-500"
      >
        Statistik Konten Website
      </p>
      <div
        v-if="loading"
        class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-4"
      >
        <div
          v-for="n in 4"
          :key="n"
          class="flex items-center gap-4 rounded-2xl bg-white p-5 ring-1 ring-slate-900/5 dark:bg-slate-900 dark:ring-white/5"
        >
          <div
            class="h-12 w-12 shrink-0 animate-pulse rounded-xl bg-gray-100 dark:bg-white/5"
          ></div>
          <div class="min-w-0 flex-1 space-y-2">
            <div
              class="h-3 w-16 animate-pulse rounded bg-gray-100 dark:bg-white/5"
            ></div>
            <div
              class="h-6 w-10 animate-pulse rounded bg-gray-100 dark:bg-white/5"
            ></div>
          </div>
        </div>
      </div>
      <TransitionGroup
        v-else
        tag="div"
        name="stat-card"
        class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-4"
      >
        <RouterLink
          v-for="item in contentStatistics"
          :key="item.title"
          :to="item.to"
          class="group flex items-center gap-4 rounded-2xl bg-white p-5 ring-1 ring-slate-900/5 shadow-sm transition-all duration-300 hover:-translate-y-0.5 hover:shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:bg-slate-900 dark:ring-white/5"
        >
          <!-- Icon -->
          <div
            :style="{ backgroundColor: item.color + '14', color: item.color }"
            class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl"
          >
            <component :is="item.icon" class="h-6 w-6" />
          </div>
          <!-- Body -->
          <div class="min-w-0 flex-1">
            <p
              class="truncate text-xs font-semibold tracking-wide text-gray-400 dark:text-slate-500"
            >
              {{ item.title }}
            </p>
            <h2
              class="text-3xl font-bold leading-tight tracking-tight text-gray-700 tabular-nums dark:text-slate-200"
            >
              {{ item.value }}
            </h2>
          </div>
          <ChevronRightIcon
            class="hidden h-4 w-4 shrink-0 text-gray-300 opacity-0 transition-all duration-300 group-hover:translate-x-0.5 group-hover:opacity-100 dark:text-slate-600 sm:block"
          />
        </RouterLink>
      </TransitionGroup>
    </div>
    <!-- Rekap per Rombel -->
    <div
      class="overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
    >
      <div
        class="flex flex-col gap-3 border-b border-gray-100 px-6 py-5 dark:border-white/5 sm:flex-row sm:items-center sm:justify-between"
      >
        <div class="flex items-start gap-3">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/8 text-[#041983] dark:bg-white/5 dark:text-indigo-300"
          >
            <ChartBarSquareIcon class="h-5 w-5" />
          </div>
          <div>
            <h2 class="text-lg font-semibold text-gray-700 dark:text-slate-200">
              Rekapitulasi Rombel &amp; Peserta Didik
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
              Jumlah peserta didik aktif berdasarkan rombel, jenis kelamin dan
              agama.
            </p>
          </div>
        </div>
        <div
          class="flex items-center gap-1.5 self-start rounded-full bg-gray-50 px-3 py-1.5 text-xs font-medium text-gray-500 ring-1 ring-slate-900/5 dark:bg-white/5 dark:text-slate-400 dark:ring-white/5 sm:self-auto"
        >
          <span class="relative flex h-2 w-2">
            <span
              class="absolute inline-flex h-full w-full motion-safe:animate-ping rounded-full bg-emerald-400 opacity-75"
            ></span>
            <span
              class="relative inline-flex h-2 w-2 rounded-full bg-emerald-500"
            ></span>
          </span>
          Data Terkini
        </div>
      </div>
      <div class="max-h-[70vh] overflow-auto">
        <table class="min-w-full w-full">
          <thead class="sticky top-0 z-20 bg-gray-50/90 dark:bg-[#0B0F1A]">
            <tr>
              <th
                scope="col"
                class="sticky left-0 z-30 bg-gray-50 px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:bg-[#0B0F1A] dark:text-slate-400"
              >
                Rombel
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-center text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Peserta Didik
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-center text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Laki-laki
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-center text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Perempuan
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-center text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Islam
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-center text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Kristen
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-center text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Katholik
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-center text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Hindu
              </th>
              <th
                scope="col"
                class="px-4 py-3 text-center text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
              >
                Budha
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 dark:divide-white/5">
            <!-- Skeleton saat memuat -->
            <template v-if="loading">
              <tr v-for="n in 5" :key="`skeleton-${n}`">
                <td
                  class="sticky left-0 z-10 bg-white px-4 py-3 dark:bg-slate-900"
                >
                  <div
                    class="h-4 w-24 animate-pulse rounded bg-gray-100 dark:bg-white/5"
                  ></div>
                </td>
                <td v-for="c in 8" :key="c" class="px-4 py-3">
                  <div
                    class="mx-auto h-4 w-8 animate-pulse rounded bg-gray-100 dark:bg-white/5"
                  ></div>
                </td>
              </tr>
            </template>
            <template v-else>
              <tr v-if="statistikStore.rekapRombel.length === 0">
                <td colspan="9" class="px-4 py-16 text-center">
                  <div
                    class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
                  >
                    <InboxIcon class="h-10 w-10" />
                    <p class="text-sm font-medium">
                      Belum ada data rekapitulasi.
                    </p>
                  </div>
                </td>
              </tr>
              <template v-for="(item, idx) in rekapWithSubtotal" :key="idx">
                <!-- Baris rombel biasa -->
                <tr
                  v-if="item.type === 'row'"
                  class="group transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
                >
                  <td
                    class="sticky left-0 z-10 bg-white px-4 py-2.5 text-sm font-medium text-gray-700 group-hover:bg-gray-50 dark:bg-slate-900 dark:text-slate-200 dark:group-hover:bg-[#0B0F1A]/40"
                  >
                    {{ item.nama_rombel }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-gray-700 dark:text-slate-300"
                  >
                    {{ item.total }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-gray-700 dark:text-slate-300"
                  >
                    {{ item.laki }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-gray-700 dark:text-slate-300"
                  >
                    {{ item.perempuan }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-gray-700 dark:text-slate-300"
                  >
                    {{ item.islam }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-gray-700 dark:text-slate-300"
                  >
                    {{ item.kristen }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-gray-700 dark:text-slate-300"
                  >
                    {{ item.katholik }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-gray-700 dark:text-slate-300"
                  >
                    {{ item.hindu }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-gray-700 dark:text-slate-300"
                  >
                    {{ item.budha }}
                  </td>
                </tr>
                <!-- Baris subtotal per kelas -->
                <tr v-else class="bg-[#041983]/5 font-semibold dark:bg-white/5">
                  <td
                    class="sticky left-0 z-10 bg-[#041983]/8 px-4 py-2.5 text-sm text-[#041983] dark:bg-slate-900/60 dark:text-indigo-300"
                  >
                    Total: {{ item.jumlah_rombel }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-[#041983] dark:text-indigo-300"
                  >
                    {{ item.total }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-[#041983] dark:text-indigo-300"
                  >
                    {{ item.laki }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-[#041983] dark:text-indigo-300"
                  >
                    {{ item.perempuan }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-[#041983] dark:text-indigo-300"
                  >
                    {{ item.islam }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-[#041983] dark:text-indigo-300"
                  >
                    {{ item.kristen }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-[#041983] dark:text-indigo-300"
                  >
                    {{ item.katholik }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-[#041983] dark:text-indigo-300"
                  >
                    {{ item.hindu }}
                  </td>
                  <td
                    class="px-4 py-2.5 text-center text-sm tabular-nums text-[#041983] dark:text-indigo-300"
                  >
                    {{ item.budha }}
                  </td>
                </tr>
              </template>
              <tr
                v-if="statistikStore.rekapRombel.length > 0"
                class="bg-[#f2c808]/10 font-bold dark:bg-white/10"
              >
                <td
                  class="sticky left-0 z-10 bg-[#f2c808]/18 px-4 py-3 text-sm text-[#041983] dark:bg-slate-900/80 dark:text-amber-300"
                >
                  Total: {{ grandTotal.jumlah_rombel }}
                </td>
                <td
                  class="px-4 py-3 text-center text-sm tabular-nums text-[#041983] dark:text-amber-300"
                >
                  {{ grandTotal.total }}
                </td>
                <td
                  class="px-4 py-3 text-center text-sm tabular-nums text-[#041983] dark:text-amber-300"
                >
                  {{ grandTotal.laki }}
                </td>
                <td
                  class="px-4 py-3 text-center text-sm tabular-nums text-[#041983] dark:text-amber-300"
                >
                  {{ grandTotal.perempuan }}
                </td>
                <td
                  class="px-4 py-3 text-center text-sm tabular-nums text-[#041983] dark:text-amber-300"
                >
                  {{ grandTotal.islam }}
                </td>
                <td
                  class="px-4 py-3 text-center text-sm tabular-nums text-[#041983] dark:text-amber-300"
                >
                  {{ grandTotal.kristen }}
                </td>
                <td
                  class="px-4 py-3 text-center text-sm tabular-nums text-[#041983] dark:text-amber-300"
                >
                  {{ grandTotal.katholik }}
                </td>
                <td
                  class="px-4 py-3 text-center text-sm tabular-nums text-[#041983] dark:text-amber-300"
                >
                  {{ grandTotal.hindu }}
                </td>
                <td
                  class="px-4 py-3 text-center text-sm tabular-nums text-[#041983] dark:text-amber-300"
                >
                  {{ grandTotal.budha }}
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
<script setup>
import { computed, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import { useStatistikStore } from "../../stores/statistik";
import { useTahunPelajaranStore } from "../../stores/tahunpelajaran";
import {
  NewspaperIcon,
  MegaphoneIcon,
  PhotoIcon,
  VideoCameraIcon,
  ChartBarSquareIcon,
  ChevronRightIcon,
  InboxIcon,
} from "@heroicons/vue/24/outline";
const statistikStore = useStatistikStore();
const tahunPelajaranStore = useTahunPelajaranStore();
const loading = ref(true);
onMounted(async () => {
  try {
    await Promise.all([
      statistikStore.fetchStatistikDashboard(),
      statistikStore.fetchRekapRombel(),
      tahunPelajaranStore.fetchAll(),
    ]);
  } catch (err) {
    console.error("Gagal memuat data dashboard:", err);
  } finally {
    loading.value = false;
  }
});
// Sapaan dinamis mengikuti jam perangkat, bikin halaman terasa lebih hidup.
const greeting = computed(() => {
  const hour = new Date().getHours();
  if (hour < 11) return "Selamat pagi,";
  if (hour < 15) return "Selamat siang,";
  if (hour < 18) return "Selamat sore,";
  return "Selamat malam,";
});
// Nama tahun pelajaran yang sedang aktif (kolom `aktif` di tbl_tahunpelajaran),
// dipakai sebagai judul, bukan hardcode string tahun.
const namaTahunPelajaranAktif = computed(() => {
  const aktif = tahunPelajaranStore.list.find((tp) => tp.aktif);
  return aktif?.nama_tahunpelajaran ?? "-";
});
// Kartu statistik konten website (sisi kanan Rekap Rombel) — jumlah_rombel,
// jumlah_murid, jumlah_guru, jumlah_tendik dipindah ke kartu akademik custom
// di bawah (donut/mini-bar), tidak lagi di sini.
const contentStatistics = computed(() => [
  {
    title: "Berita",
    value: statistikStore.statistik?.jumlah_berita ?? 0,
    icon: NewspaperIcon,
    color: "#F97316",
    to: "/dashboard/berita",
  },
  {
    title: "Pengumuman",
    value: statistikStore.statistik?.jumlah_pengumuman ?? 0,
    icon: MegaphoneIcon,
    color: "#A855F7",
    to: "/dashboard/pengumuman",
  },
  {
    title: "Galeri Foto",
    value: statistikStore.statistik?.jumlah_galerifoto ?? 0,
    icon: PhotoIcon,
    color: "#14B8A6",
    to: "/dashboard/galeri",
  },
  {
    title: "Video",
    value: statistikStore.statistik?.jumlah_video ?? 0,
    icon: VideoCameraIcon,
    color: "#EC4899",
    to: "/dashboard/video",
  },
]);
// ===== Kartu Statistik Akademik (Rombel / Peserta Didik / Tenaga Pendidik) =====
const jumlahRombelTotal = computed(
  () => statistikStore.statistik?.jumlah_rombel ?? 0,
);
const jumlahMuridTotal = computed(
  () => statistikStore.statistik?.jumlah_murid ?? 0,
);
const jumlahGuru = computed(() => statistikStore.statistik?.jumlah_guru ?? 0);
const jumlahTendik = computed(
  () => statistikStore.statistik?.jumlah_tendik ?? 0,
);
// Breakdown per tingkat (kelas 7/8/9) diturunkan langsung dari data Rekap
// Rombel yang sudah dimuat (statistikStore.rekapRombel) — TIDAK perlu
// endpoint/kolom backend baru, cukup dikelompokkan ulang di sini.
const perTingkatStats = computed(() => {
  const rows = statistikStore.rekapRombel;
  const map = new Map();
  const order = [];
  rows.forEach((row) => {
    if (!map.has(row.nama_kelas)) {
      map.set(row.nama_kelas, {
        nama_kelas: row.nama_kelas,
        jumlah_rombel: 0,
        jumlah_murid: 0,
      });
      order.push(row.nama_kelas);
    }
    const entry = map.get(row.nama_kelas);
    entry.jumlah_rombel += 1;
    entry.jumlah_murid += row.total;
  });
  return order.map((k) => map.get(k));
});
// Gradasi biru muda->tua per tingkat (identitas warna sama dgn kartu Peserta
// Didik lama #38BDF8), dipakai jg utk bar Rombel & legend donut.
const TINGKAT_COLORS = ["#7DD3FC", "#38BDF8", "#0369A1", "#0C4A6E"];
function tingkatColor(index) {
  return TINGKAT_COLORS[index % TINGKAT_COLORS.length];
}
// Warna bar Rombel per tingkat: kelas 7 hijau, kelas 8 kuning, kelas 9 merah
const ROMBEL_BAR_COLORS = ["#22C55E", "#FACC15", "#EF4444"];
function rombelBarColor(index) {
  return ROMBEL_BAR_COLORS[index % ROMBEL_BAR_COLORS.length];
}
function barHeight(value) {
  const max = Math.max(1, ...perTingkatStats.value.map((t) => t.jumlah_rombel));
  return Math.max(4, Math.round((value / max) * 36));
}
// Persentase lebar tiap segmen stacked-bar Peserta Didik
function muridPct(value) {
  const total = jumlahMuridTotal.value || 1;
  return (value / total) * 100;
}
// Radial ring bersarang Pegawai (agak besar): Guru = cincin luar, Tendik =
// cincin dalam, panjang busur tiap cincin proporsional thd bagiannya dari
// total (Guru+Tendik).
const TENAGA_RING_RADII = [32, 22];
const tenagaPendidikRings = computed(() => {
  const total = jumlahGuru.value + jumlahTendik.value || 1;
  const items = [
    { label: "Guru", value: jumlahGuru.value, color: "#FB7185" },
    { label: "Tendik", value: jumlahTendik.value, color: "#FACC15" },
  ];
  return items.map((item, i) => {
    const r = TENAGA_RING_RADII[i];
    const circumference = 2 * Math.PI * r;
    const dash = (item.value / total) * circumference;
    return { ...item, r, circumference, dash };
  });
});
// Kelompokkan baris per kelas, lalu sisipkan baris subtotal setelah baris
// terakhir tiap kelas. Hasilnya array campuran: { type: 'row', ... } dan
// { type: 'subtotal', ... } yang di-render berbeda di template.
const rekapWithSubtotal = computed(() => {
  const rows = statistikStore.rekapRombel;
  const result = [];
  let currentKelas = null;
  let subtotal = null;
  const pushSubtotal = () => {
    if (subtotal) result.push({ type: "subtotal", ...subtotal });
  };
  rows.forEach((row) => {
    if (row.nama_kelas !== currentKelas) {
      pushSubtotal();
      currentKelas = row.nama_kelas;
      subtotal = {
        nama_kelas: row.nama_kelas,
        jumlah_rombel: 0,
        total: 0,
        laki: 0,
        perempuan: 0,
        islam: 0,
        kristen: 0,
        katholik: 0,
        hindu: 0,
        budha: 0,
      };
    }
    result.push({ type: "row", ...row });
    subtotal.jumlah_rombel += 1;
    subtotal.total += row.total;
    subtotal.laki += row.laki;
    subtotal.perempuan += row.perempuan;
    subtotal.islam += row.islam;
    subtotal.kristen += row.kristen;
    subtotal.katholik += row.katholik;
    subtotal.hindu += row.hindu;
    subtotal.budha += row.budha;
  });
  pushSubtotal();
  return result;
});
const grandTotal = computed(() => {
  const rows = statistikStore.rekapRombel;
  return rows.reduce(
    (acc, row) => {
      acc.jumlah_rombel += 1;
      acc.total += row.total;
      acc.laki += row.laki;
      acc.perempuan += row.perempuan;
      acc.islam += row.islam;
      acc.kristen += row.kristen;
      acc.katholik += row.katholik;
      acc.hindu += row.hindu;
      acc.budha += row.budha;
      return acc;
    },
    {
      jumlah_rombel: 0,
      total: 0,
      laki: 0,
      perempuan: 0,
      islam: 0,
      kristen: 0,
      katholik: 0,
      hindu: 0,
      budha: 0,
    },
  );
});
</script>
<style scoped>
.stat-card-enter-active {
  transition:
    opacity 0.35s ease,
    transform 0.35s ease;
}
.stat-card-enter-from {
  opacity: 0;
  transform: translateY(6px);
}
@media (prefers-reduced-motion: reduce) {
  .stat-card-enter-active {
    transition: none;
  }
}
</style>
