<template>
  <DashboardPageHeader />
  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <!-- Teks Pengumuman -->
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Pengumuman TKA"
        description="Pengumuman hasil TKA yang tampil di halaman peserta didik."
      >
        <MegaphoneIcon class="h-5 w-5" />
      </SectionHeading>
      <div class="grid grid-cols-1 gap-6 p-6 sm:grid-cols-2">
        <FormField label="Waktu Mundur" class="sm:col-span-2">
          <input
            v-model="form.aturwaktu"
            type="text"
            :class="inputClass"
            placeholder="0000-00-00 00:00:00"
          />
        </FormField>
        <FormField label="Judul">
          <input
            v-model="form.judul"
            type="text"
            :class="inputClass"
            placeholder="PENGUMUMAN HASIL TKA"
          />
        </FormField>
        <FormField label="Nomor">
          <input
            v-model="form.nomor"
            type="text"
            :class="inputClass"
            placeholder="Nomor: 68/-1.851.522/2022"
          />
        </FormField>
        <FormField label="Pengantar" class="sm:col-span-2">
          <textarea
            v-model="form.pengantar"
            :class="inputClass"
            rows="3"
            placeholder="Berdasarkan hasil Tes Kompetensi Akademik..."
          />
        </FormField>
        <FormField label="Pernyataan" class="sm:col-span-2">
          <textarea
            v-model="form.penutup"
            :class="inputClass"
            rows="3"
            placeholder="Dalam menempuh pendidikan di ..."
          />
        </FormField>
        <FormField label="Akhir Kata" class="sm:col-span-2">
          <textarea
            v-model="form.akhir"
            :class="inputClass"
            rows="3"
            placeholder="Demikianlah pengumuman hasil TKA ini dibuat..."
          />
        </FormField>
        <FormField label="Tanggal">
          <input
            v-model="form.tanggal"
            type="text"
            :class="inputClass"
            placeholder="Jakarta, 4 Juni 2022"
          />
        </FormField>
        <FormField label="Mengetahui">
          <input
            v-model="form.mengetahui"
            type="text"
            :class="inputClass"
            placeholder="Kepala Sekolah"
          />
        </FormField>
        <FormField label="Nama Kepala Sekolah">
          <input
            v-model="form.kepsek"
            type="text"
            :class="inputClass"
            placeholder="Nama Kepala Sekolah"
          />
        </FormField>
        <FormField label="NIP Kepala Sekolah">
          <input
            v-model="form.nipkepsek"
            type="text"
            :class="inputClass"
            placeholder="NIP."
          />
        </FormField>
        <FormField label="Tanda Tangan">
          <input
            v-model="form.ttd_skl"
            type="text"
            :class="inputClass"
            placeholder="ttd"
          />
        </FormField>
      </div>
    </section>
    <!-- Aksi -->
    <div
      class="flex flex-col gap-4 rounded-2xl bg-white px-5 py-4 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5 md:flex-row md:items-center md:justify-between"
    >
      <button
        type="submit"
        :disabled="loading"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-blue-600 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:bg-gray-400 dark:disabled:bg-slate-700"
      >
        <ArrowPathIcon v-if="loading" class="h-5 w-5 animate-spin" />
        {{ loading ? "Menyimpan..." : "Simpan" }}
      </button>
      <div class="space-y-2 md:text-right">
        <div class="flex items-center gap-2 md:justify-end">
          <UserCircleIcon class="h-5 w-5 text-gray-500 dark:text-slate-400" />
          <span class="text-xs text-gray-600 dark:text-slate-400">
            Diperbarui oleh:
            {{ lastUpdatedBy }}
          </span>
        </div>
        <div class="flex items-center gap-2 md:justify-end">
          <ClockIcon class="h-5 w-5 text-gray-500 dark:text-slate-400" />
          <span class="text-xs text-gray-600 dark:text-slate-400">
            {{ lastUpdatedAt }}
          </span>
        </div>
      </div>
    </div>
  </form>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import FormField from "@/components/form/FormField.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";
import { useToast } from "@/composables/useToast";
import { useTkaStore } from "@/stores/tka";
import { formatDateTime } from "@/utils/date";
import {
  ArrowPathIcon,
  ClockIcon,
  MegaphoneIcon,
  UserCircleIcon,
} from "@heroicons/vue/24/outline";

const tkaStore = useTkaStore();
const toast = useToast();
const loading = ref(false);

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";
const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";

const form = reactive({
  aturwaktu: "",
  aturform: "",
  judul: "",
  nomor: "",
  pengantar: "",
  penutup: "",
  akhir: "",
  tanggal: "",
  mengetahui: "",
  kepsek: "",
  nipkepsek: "",
  ttd_skl: "",
  // tahun_pelajaran tidak ditampilkan di form (mengikuti versi lama yang
  // juga menonaktifkan field ini), tapi tetap disimpan agar tidak tertimpa kosong.
  tahun_pelajaran: "",
});

async function loadData() {
  loading.value = true;
  try {
    const data = await tkaStore.fetchWaktu();
    Object.assign(form, {
      aturwaktu: data.aturwaktu || "",
      aturform: data.aturform || "",
      judul: data.judul || "",
      nomor: data.nomor || "",
      pengantar: data.pengantar || "",
      penutup: data.penutup || "",
      akhir: data.akhir || "",
      tanggal: data.tanggal || "",
      mengetahui: data.mengetahui || "",
      kepsek: data.kepsek || "",
      nipkepsek: data.nipkepsek || "",
      ttd_skl: data.ttd_skl || "",
      tahun_pelajaran: data.tahun_pelajaran || "",
    });
  } catch (error) {
    toast.error(
      error.response?.data?.message ||
        "Gagal memuat pengaturan pengumuman TKA.",
    );
  } finally {
    loading.value = false;
  }
}

async function submit() {
  if (loading.value) return;
  loading.value = true;
  try {
    const response = await tkaStore.updateWaktu({ ...form });
    const dataChanged = response.changed && response.changed.length > 0;
    if (dataChanged) {
      toast.success(response.message || "Pengumuman TKA berhasil diperbarui.");
    } else {
      toast.warning(response.message || "Tidak ada perubahan data.");
    }
  } catch (error) {
    toast.error(
      error.response?.data?.message ||
        "Gagal menyimpan pengaturan pengumuman TKA.",
    );
  } finally {
    loading.value = false;
  }
}

const lastUpdatedBy = computed(() => {
  return tkaStore.waktu?.nama_pengguna || "Belum diketahui";
});
const lastUpdatedAt = computed(() => {
  if (!tkaStore.waktu?.waktu) {
    return "Belum pernah diperbarui";
  }
  return formatDateTime(tkaStore.waktu.waktu);
});

onMounted(loadData);
</script>
