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
          class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
        ></div>
        <SectionHeading
          title="Pengumuman TKA"
          description="Hasil Tes Kompetensi Akademik (TKA) akan tampil otomatis setelah waktu berikut."
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
          class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
        ></div>
        <SectionHeading
          title="Pengumuman TKA"
          description="Hasil Tes Kompetensi Akademik (TKA)"
        >
          <DocumentTextIcon class="h-5 w-5" />
        </SectionHeading>
        <div
          class="space-y-2 p-6 text-center text-sm text-gray-700 dark:text-slate-300"
        >
          <p class="font-bold underline">
            {{ waktu?.judul || "PENGUMUMAN HASIL TKA" }}
          </p>
          <p>{{ waktu?.nomor }}</p>
          <p class="whitespace-pre-line py-2">{{ waktu?.pengantar }}</p>

          <p class="pt-2">
            Nama Peserta Didik:
            <strong class="text-gray-900 dark:text-slate-100">{{
              tka?.nama_siswa || "nihil"
            }}</strong>
          </p>
          <p>
            NIS:
            <strong class="text-gray-900 dark:text-slate-100">{{
              tka?.nis || "nihil"
            }}</strong>
          </p>
          <p>
            Rombel:
            <strong class="text-gray-900 dark:text-slate-100">
              {{ tka?.nama_rombel || "nihil" }}
            </strong>
          </p>
          <p>
            No Peserta Ujian:
            <strong class="text-gray-900 dark:text-slate-100">{{
              tka?.nopes || "nihil"
            }}</strong>
          </p>

          <!-- Tabel Nilai -->
          <div class="pt-4 text-left">
            <table
              class="mx-auto w-full max-w-md border border-gray-200 text-sm dark:border-slate-700"
            >
              <thead class="bg-gray-50 dark:bg-[#0B0F1A]">
                <tr>
                  <th
                    class="border border-gray-200 px-3 py-2 dark:border-slate-700"
                  >
                    No
                  </th>
                  <th
                    class="border border-gray-200 px-3 py-2 dark:border-slate-700"
                  >
                    Mata Pelajaran
                  </th>
                  <th
                    class="border border-gray-200 px-3 py-2 dark:border-slate-700"
                  >
                    Nilai
                  </th>
                  <th
                    class="border border-gray-200 px-3 py-2 dark:border-slate-700"
                  >
                    Kategori
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td
                    class="border border-gray-200 px-3 py-2 text-center dark:border-slate-700"
                  >
                    1
                  </td>
                  <td
                    class="border border-gray-200 px-3 py-2 dark:border-slate-700"
                  >
                    Matematika
                  </td>
                  <td
                    class="border border-gray-200 px-3 py-2 text-center dark:border-slate-700"
                  >
                    {{ tka?.mtk || "nihil" }}
                  </td>
                  <td
                    class="border border-gray-200 px-3 py-2 text-center dark:border-slate-700"
                  >
                    {{ tka?.ket_mtk || "nihil" }}
                  </td>
                </tr>
                <tr>
                  <td
                    class="border border-gray-200 px-3 py-2 text-center dark:border-slate-700"
                  >
                    2
                  </td>
                  <td
                    class="border border-gray-200 px-3 py-2 dark:border-slate-700"
                  >
                    Bahasa Indonesia
                  </td>
                  <td
                    class="border border-gray-200 px-3 py-2 text-center dark:border-slate-700"
                  >
                    {{ tka?.bind || "nihil" }}
                  </td>
                  <td
                    class="border border-gray-200 px-3 py-2 text-center dark:border-slate-700"
                  >
                    {{ tka?.ket_bind || "nihil" }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="pt-4 whitespace-pre-line">
            {{ waktu?.penutup }}<br /><br />
            {{ waktu?.akhir }}<br /><br />
            {{ waktu?.tanggal }}<br />
            {{ waktu?.mengetahui }}<br /><br />
            {{ waktu?.ttd_skl }}<br /><br />
            {{ waktu?.kepsek }}<br />
            {{ waktu?.nipkepsek }}
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
import { useTkaStore } from "@/stores/tka";
import { ClockIcon, DocumentTextIcon } from "@heroicons/vue/24/outline";

const tkaStore = useTkaStore();
const loading = ref(true);
const now = ref(new Date());
let timer = null;

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";

onMounted(async () => {
  loading.value = true;
  try {
    await tkaStore.fetchTkaPesertadidik();
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

const waktu = computed(() => tkaStore.siswaData?.waktu || null);
const tka = computed(() => tkaStore.siswaData?.tka || null);

const targetDate = computed(() =>
  waktu.value?.aturwaktu
    ? new Date(waktu.value.aturwaktu.replace(" ", "T"))
    : null,
);
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
</script>
