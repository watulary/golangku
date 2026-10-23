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
              Wali Kelas Per Tingkat
            </h2>
            <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
              Catat wali kelas Kelas 7, 8, dan 9 untuk setiap siswa agar
              tercetak benar di halaman Data Nilai Buku Induk. Yang belum diisi
              dicetak kosong (untuk tanda tangan manual).
            </p>
          </div>
        </div>
      </div>

      <div class="space-y-8 px-6 py-6">
        <!-- Tingkat -->
        <div>
          <h3 class="text-sm font-semibold text-gray-700 dark:text-slate-200">
            Tingkat yang dicatat
          </h3>
          <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
            Berlaku untuk kedua cara input di bawah. Wali kelas yang disimpan
            akan tercetak di halaman Data Nilai Kelas {{ tingkat }}.
          </p>
          <div class="mt-3 flex flex-wrap gap-3">
            <label
              v-for="t in [7, 8, 9]"
              :key="t"
              class="inline-flex cursor-pointer items-center gap-2 rounded-xl border px-4 py-2 text-sm font-medium shadow-sm transition"
              :class="
                tingkat === t
                  ? 'border-blue-500 bg-blue-50 text-blue-700 dark:border-blue-400 dark:bg-blue-500/10 dark:text-blue-300'
                  : 'border-slate-300 text-slate-700 hover:border-blue-400 dark:border-slate-600 dark:text-slate-200'
              "
            >
              <input
                type="radio"
                :value="t"
                v-model.number="tingkat"
                class="h-4 w-4 border-gray-300 text-blue-600 focus:ring-blue-500"
                @change="onTingkatChange"
              />
              Kelas {{ t }}
            </label>
          </div>
        </div>

        <!-- Mode A -->
        <div class="border-t border-gray-100 pt-6 dark:border-white/5">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-slate-200">
            A. Isi per rombel
          </h3>
          <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
            Pilih rombel dan guru, lalu simpan. Nama dan NIP diterapkan ke semua
            siswa di rombel tersebut.
          </p>

          <div class="mt-3 grid gap-4 md:grid-cols-2">
            <div>
              <label :class="labelCls">Rombel</label>
              <select
                v-model.number="idRombel"
                :class="inputCls"
                @change="loadSiswa"
              >
                <option :value="0">Pilih</option>
                <option
                  v-for="r in rombelOptions"
                  :key="r.id_rombel"
                  :value="r.id_rombel"
                >
                  {{ r.nama_rombel }}
                </option>
              </select>
            </div>
            <div class="relative">
              <label :class="labelCls">Wali Kelas</label>
              <input
                v-model="guruQuery"
                type="text"
                placeholder="Ketik Nama / NIP..."
                :class="inputCls"
                @input="onGuruInput"
              />
              <ul
                v-if="showGuruList && store.pegawaiOptions.length"
                class="absolute z-20 mt-1 max-h-56 w-full overflow-auto rounded-xl border border-slate-200 bg-white shadow-lg dark:border-white/10 dark:bg-slate-900"
              >
                <li
                  v-for="p in store.pegawaiOptions"
                  :key="p.nuptk || p.nama"
                  class="cursor-pointer px-3 py-2 text-sm text-slate-700 hover:bg-blue-50 dark:text-slate-200 dark:hover:bg-slate-800"
                  @click="pilihGuru(p)"
                >
                  {{ p.nama }}
                  <span class="text-xs text-slate-500 dark:text-slate-400">{{
                    p.nip ? "NIP. " + p.nip : ""
                  }}</span>
                </li>
              </ul>
              <p
                v-if="guruTerpilih"
                class="mt-1 text-xs text-emerald-600 dark:text-emerald-400"
              >
                Dipilih: {{ guruTerpilih.nama }}
                <span v-if="guruTerpilih.nip"
                  >(NIP. {{ guruTerpilih.nip }})</span
                >
              </p>
            </div>
          </div>

          <div class="mt-4 flex flex-wrap gap-2">
            <button
              :class="btnPrimary"
              class="inline-flex items-center gap-2"
              :disabled="store.saving || !idRombel || !guruTerpilih"
              @click="simpanRombel"
            >
              <ArrowPathIcon v-if="store.saving" class="h-4 w-4 animate-spin" />
              Simpan untuk semua siswa rombel
            </button>
            <button
              :class="btnGhost"
              :disabled="!idRombel"
              @click="kosongkanRombel"
            >
              Kosongkan
            </button>
          </div>

          <div
            v-if="idRombel"
            class="mt-4 overflow-x-auto rounded-xl border border-slate-200 dark:border-white/10"
          >
            <table class="w-full text-sm text-slate-700 dark:text-slate-200">
              <thead>
                <tr
                  class="border-b border-slate-200 bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500 dark:border-white/10 dark:bg-white/5 dark:text-slate-400"
                >
                  <th class="px-3 py-2">NIS</th>
                  <th class="px-3 py-2">Nama Siswa</th>
                  <th class="px-3 py-2">Wali Kelas Kelas {{ tingkat }}</th>
                  <th class="px-3 py-2">NIP</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="s in store.siswa"
                  :key="s.nis"
                  class="border-b border-slate-100 last:border-0 dark:border-white/5"
                >
                  <td class="px-3 py-2">{{ s.nis }}</td>
                  <td class="px-3 py-2">{{ s.nama_siswa }}</td>
                  <td class="px-3 py-2">{{ s.nama_walikelas || "-" }}</td>
                  <td class="px-3 py-2">{{ s.nip || "-" }}</td>
                </tr>
                <tr v-if="!store.siswa.length && !store.loading">
                  <td
                    colspan="4"
                    class="px-3 py-4 text-center text-slate-500 dark:text-slate-400"
                  >
                    Tidak ada siswa.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Mode B -->
        <div class="border-t border-gray-100 pt-6 dark:border-white/5">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-slate-200">
            B. Impor Excel (per NIS)
          </h3>
          <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">
            Untuk Kelas 7 dan 8 yang rombel lamanya sudah bubar. Kolom: NIS,
            Nama Wali Kelas, NIP. Data disimpan untuk Kelas {{ tingkat }}.
          </p>
          <button :class="btnPrimary" class="mt-3" @click="showUpload = true">
            Unggah Excel Kelas {{ tingkat }}
          </button>
          <div
            v-if="store.lastImport?.dilewati?.length"
            class="mt-3 rounded-xl border border-amber-200 bg-amber-50 p-3 text-xs text-amber-800 dark:border-amber-500/20 dark:bg-amber-500/10 dark:text-amber-300"
          >
            <p class="mb-1 font-semibold">
              Dilewati ({{ store.lastImport.dilewati.length }}):
            </p>
            <ul class="max-h-40 list-disc overflow-auto pl-5">
              <li v-for="d in store.lastImport.dilewati" :key="d">{{ d }}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>

  <DashboardUploadDialog
    v-model="showUpload"
    :title="`Unggah Wali Kelas - Kelas ${tingkat}`"
    subtitle="Kolom: NIS, Nama Wali Kelas, NIP (.xlsx)"
    :upload-action="uploadAction"
    @download-template="unduhTemplate"
    @success="onImportSuccess"
  />
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardUploadDialog from "@/components/dashboard/DashboardUploadDialog.vue";
import { UserGroupIcon, ArrowPathIcon } from "@heroicons/vue/24/outline";
import { useToast } from "@/composables/useToast";
import { useRombelStore } from "@/stores/rombel";
import { useWaliKelasTingkatStore } from "@/stores/walikelastingkat";

const store = useWaliKelasTingkatStore();
const rombelStore = useRombelStore();
const toast = useToast();

const inputCls =
  "w-full rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm text-slate-700 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 placeholder:text-slate-400 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500 dark:focus:border-blue-400";
const labelCls =
  "mb-1 block text-xs font-medium text-slate-500 dark:text-slate-400";
const btnPrimary =
  "rounded-xl bg-blue-600 px-4 py-1.5 text-sm font-semibold text-white shadow-sm transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50";
const btnGhost =
  "rounded-xl border border-slate-300 px-4 py-1.5 text-sm font-medium text-slate-700 transition hover:border-blue-500 hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-50 dark:border-slate-600 dark:text-slate-200 dark:hover:border-blue-400 dark:hover:bg-slate-800";

const tingkat = ref(7);
const idRombel = ref(0);
const guruQuery = ref("");
const guruTerpilih = ref(null);
const showGuruList = ref(false);

const rombelOptions = computed(() =>
  (rombelStore.list || []).filter((r) => r.aktif === 1),
);

let timer = null;
function onGuruInput() {
  guruTerpilih.value = null;
  showGuruList.value = true;
  clearTimeout(timer);
  timer = setTimeout(() => store.searchPegawai(guruQuery.value), 250);
}

function pilihGuru(p) {
  guruTerpilih.value = p;
  guruQuery.value = p.nama;
  showGuruList.value = false;
}

async function loadSiswa() {
  if (!idRombel.value) return;
  try {
    await store.fetchSiswa(idRombel.value, tingkat.value);
  } catch (e) {
    toast.error(e.response?.data?.message || "Gagal memuat siswa.");
  }
}

function onTingkatChange() {
  loadSiswa();
}

async function simpanRombel() {
  try {
    const res = await store.saveRombel({
      id_rombel: idRombel.value,
      tingkat: tingkat.value,
      nuptk: guruTerpilih.value.nuptk,
      nama: guruTerpilih.value.nama,
      nip: guruTerpilih.value.nip,
    });
    toast.success(res.message || "Tersimpan.");
    await loadSiswa();
  } catch (e) {
    toast.error(e.response?.data?.message || "Gagal menyimpan.");
  }
}

async function kosongkanRombel() {
  try {
    const res = await store.clearRombel(idRombel.value, tingkat.value);
    toast.success(res.message || "Dikosongkan.");
    await loadSiswa();
  } catch (e) {
    toast.error(e.response?.data?.message || "Gagal mengosongkan.");
  }
}

const showUpload = ref(false);

// dialog memanggil uploadAction(file); tingkat diambil saat unggah
function uploadAction(file) {
  return store.importExcel(file, tingkat.value);
}

async function onImportSuccess(message) {
  toast.success(message || "Impor selesai.");
  if (idRombel.value) await loadSiswa();
}

async function unduhTemplate() {
  try {
    await store.downloadTemplate();
  } catch {
    toast.error("Gagal mengunduh template.");
  }
}

onMounted(() => {
  rombelStore.pagination.limit = 500;
  rombelStore.fetchAll();
});
</script>
