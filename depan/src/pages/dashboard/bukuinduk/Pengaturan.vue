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
            <Cog6ToothIcon class="h-5 w-5" />
          </div>
          <div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Pengaturan Cetak Buku Induk
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
              Atur lembar mana saja yang ikut tercetak, dan boleh isi nama
              Kepala Sekolah / Wali Kelas kalau perlu tampil beda dari setting
              biasa. Berlaku otomatis setiap tombol Cetak Buku Induk dipakai
              (satuan maupun massal), sampai diubah lagi di sini.
            </p>
          </div>
        </div>
      </div>

      <div v-if="store.loading" class="space-y-3 px-6 py-8">
        <div
          v-for="n in 5"
          :key="n"
          class="h-4 w-full max-w-md animate-pulse rounded bg-gray-100 dark:bg-white/5"
        ></div>
      </div>

      <form v-else @submit.prevent="submit" class="space-y-8 px-6 py-6">
        <!-- Lembar yang dicetak -->
        <div>
          <h3 class="text-sm font-semibold text-gray-700 dark:text-slate-200">
            Lembar yang Dicetak
          </h3>
          <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
            Matikan lembar depan buku (biasanya cukup dicetak 1x di awal tahun
            pelajaran) kalau hanya ingin cetak arsip data nilai semester
            berjalan.
          </p>
          <div class="mt-3 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <label
              v-for="opt in lembarOptions"
              :key="opt.key"
              :class="checkboxRowClass"
            >
              <input
                type="checkbox"
                v-model="store.opsi[opt.key]"
                class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              />
              <span>
                <span
                  class="block text-sm font-medium text-gray-700 dark:text-slate-200"
                >
                  {{ opt.label }}
                </span>
                <span class="block text-xs text-slate-500 dark:text-slate-400">
                  {{ opt.desc }}
                </span>
              </span>
            </label>
          </div>
        </div>

        <!-- Tingkat Data Nilai -->
        <div>
          <h3 class="text-sm font-semibold text-gray-700 dark:text-slate-200">
            Tingkat Data Nilai
          </h3>
          <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
            Pilih tingkat kelas mana yang ikut dicetak di halaman Data
            Nilai/Ketidakhadiran. Kalau baru mengisi nilai Kelas 7, cukup
            centang Kelas 7 saja - halaman Kelas 8 & 9 yang masih kosong tidak
            ikut tercetak.
          </p>
          <div class="mt-3 flex flex-wrap gap-3">
            <label
              v-for="t in [7, 8, 9]"
              :key="t"
              class="inline-flex items-center gap-2 rounded-xl border border-slate-300 px-4 py-2 text-sm font-medium text-slate-700 shadow-sm dark:border-slate-600 dark:text-slate-200"
            >
              <input
                type="checkbox"
                :value="t"
                v-model="store.opsi.tingkat"
                class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              />
              Kelas {{ t }}
            </label>
          </div>
        </div>

        <!-- Override Kepala Sekolah & Wali Kelas -->
        <div>
          <h3 class="text-sm font-semibold text-gray-700 dark:text-slate-200">
            Nama Kepala Sekolah &amp; Wali Kelas (opsional)
          </h3>
          <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
            Kosongkan supaya tetap memakai data dari pengaturan standar. Isi
            hanya kalau untuk cetak ini butuh nama yang berbeda (mis. cetak
            arsip histori dengan kepala sekolah/wali kelas yang menjabat saat
            itu). Wali Kelas di sini hanya memengaruhi blok tanda tangan di
            halaman depan (Daftar Isi &amp; Nomor Urut), bukan data per siswa di
            halaman Data Nilai masing-masing.
          </p>
          <div class="mt-3 grid gap-4 sm:grid-cols-2">
            <div>
              <label :class="labelClass">Nama Kepala Sekolah</label>
              <input
                v-model="store.opsi.kepsek_nama"
                type="text"
                placeholder=""
                :class="inputClass"
              />
            </div>
            <div>
              <label :class="labelClass">NIP Kepala Sekolah</label>
              <input
                v-model="store.opsi.kepsek_nip"
                type="text"
                placeholder=""
                :class="inputClass"
              />
            </div>
            <div>
              <label :class="labelClass">Nama Wali Kelas</label>
              <input
                v-model="store.opsi.walikelas_nama"
                type="text"
                placeholder=""
                :class="inputClass"
              />
            </div>
            <div>
              <label :class="labelClass">NIP Wali Kelas</label>
              <input
                v-model="store.opsi.walikelas_nip"
                type="text"
                placeholder=""
                :class="inputClass"
              />
            </div>
          </div>
        </div>

        <div
          class="flex justify-end border-t border-gray-100 pt-5 dark:border-white/5"
        >
          <button
            type="submit"
            :disabled="store.saving"
            class="inline-flex items-center gap-2 rounded-xl bg-blue-600 px-5 py-2.5 text-sm font-medium text-white shadow-sm transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:opacity-50"
          >
            <ArrowPathIcon v-if="store.saving" class="h-4 w-4 animate-spin" />
            Simpan Pengaturan
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from "vue";
import { useOpsiCetakStore } from "@/stores/opsicetak";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import { useToast } from "@/composables/useToast";
import { Cog6ToothIcon, ArrowPathIcon } from "@heroicons/vue/24/outline";

const store = useOpsiCetakStore();
const toast = useToast();

const labelClass =
  "mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400";
const inputClass =
  "w-full rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm text-slate-700 shadow-sm transition focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200";
const checkboxRowClass =
  "flex items-start gap-3 rounded-xl border border-slate-200 p-3 dark:border-white/10";

const lembarOptions = [
  { key: "sampul", label: "Sampul", desc: "Halaman judul/cover buku." },
  { key: "daftar_isi", label: "Daftar Isi", desc: "Ringkasan isi buku." },
  {
    key: "pemeriksaan",
    label: "I. Pemeriksaan",
    desc: "Tabel pemeriksaan buku induk (kosong, diisi tangan).",
  },
  {
    key: "nomor_urut",
    label: "II. Nomor Urut Murid",
    desc: "Daftar ringkas seluruh siswa yang dicetak.",
  },
  {
    key: "petunjuk",
    label: "III. Petunjuk Pengisian",
    desc: "Petunjuk umum & khusus.",
  },
  {
    key: "lembar_iv",
    label: "IV. Biodata Siswa",
    desc: "Keterangan Siswa, Orang Tua, Perkembangan, dst (per siswa).",
  },
  {
    key: "halaman_akhir",
    label: "Halaman Akhir",
    desc: "Rata-rata/Ujian, Ekstrakurikuler, Keputusan Lulus/Tidak Lulus (per siswa).",
  },
];

onMounted(async () => {
  try {
    await store.fetchOpsi();
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat pengaturan cetak.",
    );
  }
});

async function submit() {
  try {
    const res = await store.saveOpsi();
    toast.success(res.message || "Pengaturan berhasil disimpan.");
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menyimpan pengaturan.");
  }
}
</script>
