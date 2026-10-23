<template>
  <DashboardPageHeader />

  <div class="mx-auto w-full space-y-8">
    <!-- ===================== FORM 1: PENGUMUMAN ===================== -->
    <form class="space-y-6" @submit.prevent="submitPengumuman">
      <section :class="cardClass" class="relative">
        <div
          class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
        ></div>

        <SectionHeading
          title="Pengumuman Kelulusan"
          description="Pengumuman kelulusan yang tampil di halaman peserta didik."
        >
          <MegaphoneIcon class="h-5 w-5" />
        </SectionHeading>
        <div class="grid grid-cols-1 gap-6 p-6 sm:grid-cols-2">
          <FormField label="Waktu Mundur" class="sm:col-span-2">
            <input
              v-model="pengumumanForm.aturwaktu"
              type="text"
              :class="inputClass"
              placeholder="0000-00-00 00:00:00"
            />
          </FormField>
          <FormField label="Judul">
            <input
              v-model="pengumumanForm.judul"
              type="text"
              :class="inputClass"
              placeholder="PENGUMUMAN KELULUSAN"
            />
          </FormField>
          <FormField label="Nomor">
            <input
              v-model="pengumumanForm.nomor"
              type="text"
              :class="inputClass"
              placeholder="Nomor: 68/-1.851.522/2021"
            />
          </FormField>
          <FormField label="Pengantar" class="sm:col-span-2">
            <textarea
              v-model="pengumumanForm.pengantar"
              :class="inputClass"
              rows="3"
              placeholder="Berdasarkan hasil Rapat Pleno Kelulusan..."
            />
          </FormField>
          <FormField label="Pernyataan" class="sm:col-span-2">
            <textarea
              v-model="pengumumanForm.penutup"
              :class="inputClass"
              rows="3"
              placeholder="Dalam menempuh pendidikan di ..."
            />
          </FormField>
          <FormField label="Akhir Kata" class="sm:col-span-2">
            <textarea
              v-model="pengumumanForm.akhir"
              :class="inputClass"
              rows="3"
              placeholder="Demikianlah pengumuman kelulusan ini dibuat..."
            />
          </FormField>
          <FormField label="Tanggal">
            <input
              v-model="pengumumanForm.tanggal"
              type="text"
              :class="inputClass"
              placeholder="Jakarta, 7 Juni 2021"
            />
          </FormField>
          <FormField label="Mengetahui">
            <input
              v-model="pengumumanForm.mengetahui"
              type="text"
              :class="inputClass"
              placeholder="Kepala Sekolah"
            />
          </FormField>
          <FormField label="Nama Kepala Sekolah">
            <input
              v-model="pengumumanForm.kepsek"
              type="text"
              :class="inputClass"
              placeholder="Nama Kepala Sekolah"
            />
          </FormField>
          <FormField label="NIP Kepala Sekolah">
            <input
              v-model="pengumumanForm.nipkepsek"
              type="text"
              :class="inputClass"
              placeholder="NIP."
            />
          </FormField>
          <FormField label="Tahun Pelajaran">
            <input
              v-model="pengumumanForm.tahun_pelajaran"
              type="text"
              :class="inputClass"
              placeholder="2020/2021"
            />
          </FormField>
        </div>
      </section>

      <!-- Toggle tampilan (khusus Kelulusan, tidak ada padanan di Kenaikan) -->
      <section :class="cardClass">
        <SectionHeading
          title="Pengaturan Tampilan"
          description="Kontrol elemen yang muncul di halaman Peserta Didik dan cetak SKL."
        >
          <AdjustmentsHorizontalIcon class="h-5 w-5" />
        </SectionHeading>
        <div class="grid grid-cols-1 gap-6 p-6 sm:grid-cols-3">
          <FormField label="Ttd & Stempel pada SKL">
            <select v-model="pengumumanForm.heroarea1" :class="inputClass">
              <option value="1">Tampilkan</option>
              <option value="0">Sembunyikan</option>
            </select>
          </FormField>
          <FormField label="Foto Siswa pada SKL">
            <select v-model="pengumumanForm.aturform" :class="inputClass">
              <option value="1">Tampilkan</option>
              <option value="0">Sembunyikan</option>
            </select>
          </FormField>
          <FormField label="Tombol Unduh SKL (Peserta Didik)">
            <select v-model="pengumumanForm.heroarea2" :class="inputClass">
              <option value="1">Tampilkan</option>
              <option value="0">Sembunyikan</option>
            </select>
          </FormField>
        </div>
      </section>

      <!-- Aksi Form 1 -->
      <div
        class="flex flex-col gap-4 rounded-2xl bg-white px-5 py-4 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5 md:flex-row md:items-center md:justify-between"
      >
        <button
          type="submit"
          :disabled="loadingPengumuman"
          class="inline-flex items-center justify-center gap-2 rounded-xl bg-blue-600 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:bg-gray-400 dark:disabled:bg-slate-700"
        >
          <ArrowPathIcon
            v-if="loadingPengumuman"
            class="h-5 w-5 animate-spin"
          />
          {{ loadingPengumuman ? "Menyimpan..." : "Simpan Pengumuman" }}
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

    <!-- ===================== FORM 2: ISI SURAT SKL ===================== -->
    <form class="space-y-6" @submit.prevent="submitSKL">
      <section :class="cardClass" class="relative">
        <div
          class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
        ></div>

        <SectionHeading
          title="Isi Surat Keterangan Lulus (SKL)"
          description="Kop surat dan isi naskah yang dicetak pada dokumen SKL."
        >
          <DocumentTextIcon class="h-5 w-5" />
        </SectionHeading>
        <div class="grid grid-cols-1 gap-6 p-6">
          <FormField label="Kop Baris 1">
            <input
              v-model="sklForm.kop1_skl"
              type="text"
              :class="inputClass"
              placeholder="PEMERINTAH PROVINSI DKI JAKARTA"
            />
          </FormField>
          <FormField label="Kop Baris 2">
            <input
              v-model="sklForm.kop2_skl"
              type="text"
              :class="inputClass"
              placeholder="DINAS PENDIDIKAN"
            />
          </FormField>
          <FormField label="Kop Baris 3">
            <input
              v-model="sklForm.kop3_skl"
              type="text"
              :class="inputClass"
              placeholder="SMP NEGERI 254 JAKARTA"
            />
          </FormField>
          <FormField label="Kop Baris 4">
            <input
              v-model="sklForm.kop4_skl"
              type="text"
              :class="inputClass"
              placeholder="Alamat sekolah"
            />
          </FormField>
          <FormField label="Kop Baris 5">
            <input
              v-model="sklForm.kop5_skl"
              type="text"
              :class="inputClass"
              placeholder="Kelurahan / Kecamatan"
            />
          </FormField>
          <FormField label="Kop Baris 6">
            <input
              v-model="sklForm.kop6_skl"
              type="text"
              :class="inputClass"
              placeholder="Kota / Provinsi"
            />
          </FormField>
          <FormField label="Kop Baris 7 (Kontak)">
            <input
              v-model="sklForm.kop7_skl"
              type="text"
              :class="inputClass"
              placeholder="Telp / Email / Website"
            />
          </FormField>
          <FormField label="Kop Baris 8 (cadangan, tidak dicetak)">
            <input
              v-model="sklForm.kop8_skl"
              type="text"
              :class="inputClass"
              placeholder="Tidak digunakan saat ini"
            />
          </FormField>
        </div>
      </section>

      <section :class="cardClass">
        <SectionHeading
          title="Naskah SKL"
          description="Urutan tampil pada dokumen: Isi 1 → Isi 2 → Isi 3 → data siswa & nilai → Isi 4 → tabel nilai → Isi 5."
        >
          <ClipboardDocumentCheckIcon class="h-5 w-5" />
        </SectionHeading>
        <div class="grid grid-cols-1 gap-6 p-6">
          <FormField label="Isi 1 (Judul Surat)">
            <input
              v-model="sklForm.isi1_skl"
              type="text"
              :class="inputClass"
              placeholder="SURAT KETERANGAN LULUS"
            />
          </FormField>
          <FormField label="Isi 2 (Nomor Surat)">
            <input
              v-model="sklForm.isi2_skl"
              type="text"
              :class="inputClass"
              placeholder="Nomor: 68/-1.851.522/2021"
            />
          </FormField>
          <FormField label="Isi 3 (Pembuka)">
            <textarea
              v-model="sklForm.isi3_skl"
              :class="inputClass"
              rows="3"
              placeholder="Yang bertanda tangan di bawah ini Kepala Sekolah..."
            />
          </FormField>
          <FormField label="Isi 4 (Setelah Data Nilai)">
            <textarea
              v-model="sklForm.isi4_skl"
              :class="inputClass"
              rows="3"
              placeholder="dari Sekolah ... Tahun Pelajaran ..., setelah memenuhi seluruh kriteria dengan nilai sebagai berikut;"
            />
          </FormField>
          <FormField label="Isi 5 (Penutup)">
            <textarea
              v-model="sklForm.isi5_skl"
              :class="inputClass"
              rows="3"
              placeholder="Demikian surat keterangan lulus ini dibuat agar dapat dipergunakan sebagaimana mestinya."
            />
          </FormField>
        </div>
      </section>

      <section :class="cardClass">
        <SectionHeading
          title="Blok Tanda Tangan SKL"
          description="Tanggal, nama sekolah, dan identitas kepala sekolah pada blok tanda tangan."
        >
          <PencilSquareIcon class="h-5 w-5" />
        </SectionHeading>
        <div class="grid grid-cols-1 gap-6 p-6 sm:grid-cols-2">
          <FormField label="Tanggal">
            <input
              v-model="sklForm.tgl_skl"
              type="text"
              :class="inputClass"
              placeholder="Jakarta, 7 Juni 2021"
            />
          </FormField>
          <FormField label="Nama Sekolah">
            <input
              v-model="sklForm.sekolah_skl"
              type="text"
              :class="inputClass"
              placeholder="Kepala SMPN 254 Jakarta"
            />
          </FormField>
          <FormField label="Nama Kepala Sekolah">
            <input
              v-model="sklForm.kepsek_skl"
              type="text"
              :class="inputClass"
              placeholder="Suryoto, M.Pd"
            />
          </FormField>
          <FormField label="NIP Kepala Sekolah">
            <input
              v-model="sklForm.nipkepsek_skl"
              type="text"
              :class="inputClass"
              placeholder="NIP. 196210091987031006"
            />
          </FormField>
          <FormField label="Tanda Tangan">
            <input
              v-model="sklForm.ttd_skl"
              type="text"
              :class="inputClass"
              placeholder="ttd"
            />
          </FormField>
        </div>
      </section>

      <!-- Aksi Form 2 -->
      <div
        class="flex flex-col gap-4 rounded-2xl bg-white px-5 py-4 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5 md:flex-row md:items-center md:justify-between"
      >
        <button
          type="submit"
          :disabled="loadingSKL"
          class="inline-flex items-center justify-center gap-2 rounded-xl bg-blue-600 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:bg-gray-400 dark:disabled:bg-slate-700"
        >
          <ArrowPathIcon v-if="loadingSKL" class="h-5 w-5 animate-spin" />
          {{ loadingSKL ? "Menyimpan..." : "Simpan Isi SKL" }}
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
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import FormField from "@/components/form/FormField.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";
import { useToast } from "@/composables/useToast";
import { useKelulusanStore } from "@/stores/kelulusan";
import { formatDateTime } from "@/utils/date";
import {
  AdjustmentsHorizontalIcon,
  ArrowPathIcon,
  ClipboardDocumentCheckIcon,
  ClockIcon,
  DocumentTextIcon,
  MegaphoneIcon,
  PencilSquareIcon,
  UserCircleIcon,
} from "@heroicons/vue/24/outline";

const kelulusanStore = useKelulusanStore();
const toast = useToast();

const loadingPengumuman = ref(false);
const loadingSKL = ref(false);

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";
const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";

// ===== Form 1: Pengumuman =====
const pengumumanForm = reactive({
  aturwaktu: "",
  aturform: "1",
  heroarea1: "1",
  heroarea2: "1",
  judul: "",
  nomor: "",
  pengantar: "",
  penutup: "",
  akhir: "",
  tanggal: "",
  mengetahui: "",
  kepsek: "",
  nipkepsek: "",
  tahun_pelajaran: "",
});

// ===== Form 2: Isi Surat SKL =====
const sklForm = reactive({
  kop1_skl: "",
  kop2_skl: "",
  kop3_skl: "",
  kop4_skl: "",
  kop5_skl: "",
  kop6_skl: "",
  kop7_skl: "",
  kop8_skl: "",
  isi1_skl: "",
  isi2_skl: "",
  isi3_skl: "",
  isi4_skl: "",
  isi5_skl: "",
  tgl_skl: "",
  sekolah_skl: "",
  kepsek_skl: "",
  nipkepsek_skl: "",
  ttd_skl: "",
});

async function loadData() {
  try {
    const data = await kelulusanStore.fetchWaktu();
    Object.assign(pengumumanForm, {
      aturwaktu: data.aturwaktu || "",
      aturform: data.aturform || "1",
      heroarea1: data.heroarea1 || "1",
      heroarea2: data.heroarea2 || "1",
      judul: data.judul || "",
      nomor: data.nomor || "",
      pengantar: data.pengantar || "",
      penutup: data.penutup || "",
      akhir: data.akhir || "",
      tanggal: data.tanggal || "",
      mengetahui: data.mengetahui || "",
      kepsek: data.kepsek || "",
      nipkepsek: data.nipkepsek || "",
      tahun_pelajaran: data.tahun_pelajaran || "",
    });
    Object.assign(sklForm, {
      kop1_skl: data.kop1_skl || "",
      kop2_skl: data.kop2_skl || "",
      kop3_skl: data.kop3_skl || "",
      kop4_skl: data.kop4_skl || "",
      kop5_skl: data.kop5_skl || "",
      kop6_skl: data.kop6_skl || "",
      kop7_skl: data.kop7_skl || "",
      kop8_skl: data.kop8_skl || "",
      isi1_skl: data.isi1_skl || "",
      isi2_skl: data.isi2_skl || "",
      isi3_skl: data.isi3_skl || "",
      isi4_skl: data.isi4_skl || "",
      isi5_skl: data.isi5_skl || "",
      tgl_skl: data.tgl_skl || "",
      sekolah_skl: data.sekolah_skl || "",
      kepsek_skl: data.kepsek_skl || "",
      nipkepsek_skl: data.nipkepsek_skl || "",
      ttd_skl: data.ttd_skl || "",
    });
  } catch (error) {
    toast.error(
      error.response?.data?.message ||
        "Gagal memuat pengaturan pengumuman kelulusan.",
    );
  }
}

async function submitPengumuman() {
  if (loadingPengumuman.value) return;
  loadingPengumuman.value = true;
  try {
    const response = await kelulusanStore.updateWaktu({ ...pengumumanForm });
    const dataChanged = response.changed && response.changed.length > 0;
    if (dataChanged) {
      toast.success(
        response.message || "Pengumuman kelulusan berhasil diperbarui.",
      );
    } else {
      toast.warning(response.message || "Tidak ada perubahan data.");
    }
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menyimpan pengumuman kelulusan.",
    );
  } finally {
    loadingPengumuman.value = false;
  }
}

async function submitSKL() {
  if (loadingSKL.value) return;
  loadingSKL.value = true;
  try {
    const response = await kelulusanStore.updateSKL({ ...sklForm });
    const dataChanged = response.changed && response.changed.length > 0;
    if (dataChanged) {
      toast.success(response.message || "Isi surat SKL berhasil diperbarui.");
    } else {
      toast.warning(response.message || "Tidak ada perubahan data.");
    }
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menyimpan isi surat SKL.",
    );
  } finally {
    loadingSKL.value = false;
  }
}

const lastUpdatedBy = computed(() => {
  return kelulusanStore.waktu?.nama_pengguna || "Belum diketahui";
});
const lastUpdatedAt = computed(() => {
  if (!kelulusanStore.waktu?.waktu) {
    return "Belum pernah diperbarui";
  }
  return formatDateTime(kelulusanStore.waktu.waktu);
});

onMounted(loadData);
</script>
