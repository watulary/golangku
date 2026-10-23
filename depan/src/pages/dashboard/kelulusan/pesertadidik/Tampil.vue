<template>
  <DashboardPageHeader />
  <div class="mx-auto w-full space-y-6">
    <div v-if="loading" class="flex flex-col items-center justify-center py-20">
      <div
        class="h-10 w-10 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"
      />
      <p class="mt-4 text-sm text-gray-500 dark:text-slate-400">
        Memuat data...
      </p>
    </div>
    <template v-else>
      <!-- Countdown -->
      <section v-if="!sudahWaktunya" :class="cardClass" class="relative">
        <div
          class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
        ></div>
        <SectionHeading
          title="Pengumuman Kelulusan"
          description="Hasil kelulusan akan tampil otomatis setelah waktu berikut."
        >
          <ClockIcon class="h-5 w-5" />
        </SectionHeading>
        <div class="flex justify-center gap-6 p-8">
          <div
            v-for="(val, label) in {
              Hari: sisaWaktu.hari,
              Jam: sisaWaktu.jam,
              Menit: sisaWaktu.menit,
              Detik: sisaWaktu.detik,
            }"
            :key="label"
            class="text-center"
          >
            <div class="text-3xl font-bold text-blue-600 dark:text-blue-400">
              {{ val }}
            </div>
            <div class="text-xs text-gray-500 dark:text-slate-400">
              {{ label }}
            </div>
          </div>
        </div>
      </section>

      <!-- Naskah Pengumuman (satu alur utuh) -->
      <section v-else :class="cardClass" class="relative">
        <div
          class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
        ></div>
        <SectionHeading
          title="Pengumuman Kelulusan"
          description="Pengumuman resmi hasil kelulusan dari sekolah."
        >
          <DocumentTextIcon class="h-5 w-5" />
        </SectionHeading>
        <div
          class="space-y-2 p-6 text-center text-sm text-gray-700 dark:text-slate-300"
        >
          <p class="font-bold underline">
            {{ waktu?.judul || "PENGUMUMAN KELULUSAN" }}
          </p>
          <p>{{ waktu?.nomor }}</p>
          <p class="whitespace-pre-line py-2">{{ waktu?.pengantar }}</p>

          <img
            :src="fotoURL || `/storage/fotokelulusan/${defaultFoto}`"
            alt="Foto peserta didik"
            class="mx-auto h-38.75 w-30 rounded-lg border border-gray-200 object-cover dark:border-slate-700"
          />

          <p class="pt-2">
            Nama Peserta Didik:
            <strong class="text-gray-900 dark:text-slate-100">{{
              kelulusan?.nama_siswa || "nihil"
            }}</strong>
          </p>
          <p>
            Tempat, Tanggal Lahir:
            <strong class="text-gray-900 dark:text-slate-100">
              {{ kelulusan?.tempat_lahirsiswa || "nihil" }},
              {{ formatTanggalLahir(kelulusan?.tanggal_lahirsiswa) }}
            </strong>
          </p>
          <p>
            Kelas:
            <strong class="text-gray-900 dark:text-slate-100">{{
              kelulusan?.nama_kelas || "nihil"
            }}</strong>
          </p>
          <p>
            Rombel:
            <strong class="text-gray-900 dark:text-slate-100">{{
              kelulusan?.nama_rombel || "nihil"
            }}</strong>
          </p>
          <p>
            NIS:
            <strong class="text-gray-900 dark:text-slate-100">{{
              kelulusan?.nis || "nihil"
            }}</strong>
          </p>
          <p>
            NISN:
            <strong class="text-gray-900 dark:text-slate-100">{{
              kelulusan?.nisn || "nihil"
            }}</strong>
          </p>
          <p>
            No Peserta:
            <strong class="text-gray-900 dark:text-slate-100">{{
              kelulusan?.nopes || "nihil"
            }}</strong>
          </p>

          <p class="pt-2">
            Dinyatakan:
            <strong class="text-gray-900 dark:text-slate-100">{{
              kelulusan?.ket || "nihil"
            }}</strong>
          </p>

          <div class="pt-4 whitespace-pre-line">
            {{ waktu?.penutup }}<br /><br />
            {{ waktu?.akhir }}<br /><br />
            {{ waktu?.tanggal }}<br />
            {{ waktu?.mengetahui }}<br /><br />
            {{ waktu?.ttd_skl }}<br /><br />
            {{ waktu?.kepsek }}<br />
            {{ waktu?.nipkepsek }}
          </div>

          <!-- Tombol Unduh SKL: hanya kalau heroarea2 aktif dan data kelulusan ada -->
          <div v-if="tampilTombolUnduh" class="pt-6">
            <button
              type="button"
              :disabled="cetakLoading"
              @click="handleUnduh"
              class="inline-flex items-center justify-center gap-2 rounded-xl bg-emerald-600 px-5 py-2.5 text-sm font-medium text-white shadow-sm transition hover:bg-emerald-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 disabled:cursor-not-allowed disabled:bg-gray-400 dark:disabled:bg-slate-700"
            >
              <ArrowPathIcon v-if="cetakLoading" class="h-5 w-5 animate-spin" />
              <DocumentArrowDownIcon v-else class="h-5 w-5" />
              {{ cetakLoading ? "Menyiapkan..." : "Unduh SKL" }}
            </button>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from "vue";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";
import { useToast } from "@/composables/useToast";
import { useKelulusanStore } from "@/stores/kelulusan";
import {
  ArrowPathIcon,
  ClockIcon,
  DocumentArrowDownIcon,
  DocumentTextIcon,
} from "@heroicons/vue/24/outline";

const kelulusanStore = useKelulusanStore();
const toast = useToast();

const loading = ref(true);
const cetakLoading = ref(false);
const now = ref(new Date());
let timer = null;

const defaultFoto = "peserta_didik.jpg";
const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";

onMounted(async () => {
  loading.value = true;
  try {
    await kelulusanStore.fetchKelulusanPesertadidik();
  } catch (error) {
    toast.error(
      error.response?.data?.message ||
        "Gagal mengambil data kelulusan. Coba muat ulang halaman.",
    );
  } finally {
    loading.value = false;
  }
  timer = setInterval(() => {
    now.value = new Date();
  }, 1000);
});
onUnmounted(() => {
  if (timer) clearInterval(timer);
});

// Backend return { success, data: { waktu, kelulusan } }, sama pola dgn
// /kenaikan/me yg return { success, message, data: { waktu, kenaikan } }.
const waktu = computed(() => kelulusanStore.siswaData?.data?.waktu || null);
const kelulusan = computed(
  () => kelulusanStore.siswaData?.data?.kelulusan || null,
);

const targetDate = computed(() => {
  const raw = waktu.value?.aturwaktu;
  if (!raw) return null;
  // Backend selalu format "YYYY-MM-DD HH:MM:SS" (waktu WIB/Asia-Jakarta),
  // TANPA info zona. Kalau langsung new Date(raw.replace(" ","T")), JS
  // menafsirkannya sebagai waktu LOKAL perangkat yang membuka halaman ini
  // (bisa salah kalau OS/browser bukan zona WIB) -- jadi offset +07:00
  // dipaksa eksplisit di sini supaya perbandingan waktu selalu konsisten.
  const iso = raw.replace(" ", "T") + (raw.includes("+") ? "" : "+07:00");
  const parsed = new Date(iso);
  return isNaN(parsed.getTime()) ? null : parsed;
});
const sudahWaktunya = computed(() => {
  if (!targetDate.value) return false;
  return now.value.getTime() >= targetDate.value.getTime();
});
const sisaWaktu = computed(() => {
  if (!targetDate.value) return { hari: 0, jam: 0, menit: 0, detik: 0 };
  const distance = Math.max(
    0,
    targetDate.value.getTime() - now.value.getTime(),
  );
  return {
    hari: Math.floor(distance / 86400000),
    jam: Math.floor((distance % 86400000) / 3600000),
    menit: Math.floor((distance % 3600000) / 60000),
    detik: Math.floor((distance % 60000) / 1000),
  };
});

const bulanIndo = [
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
function formatTanggalLahir(tgl) {
  if (!tgl) return "nihil";
  const [thn, bln, hr] = tgl.split("-");
  if (!thn || !bln || !hr) return "nihil";
  return `${parseInt(hr, 10)} ${bulanIndo[parseInt(bln, 10) - 1]} ${thn}`;
}

const fotoURL = computed(() => kelulusan.value?.foto_url || null);

// heroarea2 mengontrol tombol unduh, sama seperti index.php lama:
// !empty($waktululus['heroarea2']) && !empty($kelulusan['nis'])
const tampilTombolUnduh = computed(() => {
  const heroarea2 = waktu.value?.heroarea2;
  return (
    heroarea2 !== "" &&
    heroarea2 !== "0" &&
    heroarea2 != null &&
    !!kelulusan.value?.nis
  );
});

async function handleUnduh() {
  if (!kelulusan.value?.nis) return;
  cetakLoading.value = true;
  try {
    await kelulusanStore.cetakSKL(kelulusan.value.nis);
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal mengunduh SKL.");
  } finally {
    cetakLoading.value = false;
  }
}
</script>
