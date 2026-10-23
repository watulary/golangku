<template>
  <DashboardPageHeader />
  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <!-- Foto -->
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Foto Peserta Didik"
        description="Foto profil peserta didik yang tampil pada daftar dan detil."
      >
        <PhotoIcon class="h-5 w-5" />
      </SectionHeading>
      <div class="flex justify-start p-6">
        <div class="group relative">
          <div
            class="h-32 w-24 overflow-hidden rounded-lg border border-gray-200 bg-gray-50 dark:border-slate-700 dark:bg-slate-800"
            :class="{ 'cursor-pointer': form.foto_url }"
            @click="form.foto_url && showImage()"
          >
            <img
              v-if="form.foto_url"
              :src="form.foto_url"
              alt="Foto Peserta Didik"
              class="h-full w-full object-cover transition duration-200 group-hover:opacity-80"
            />
            <div
              v-else
              class="flex h-full w-full items-center justify-center text-gray-300 dark:text-slate-600"
            >
              <UserIcon class="h-12 w-12" />
            </div>
            <div
              v-if="form.foto_url"
              class="pointer-events-none absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition duration-200 group-hover:opacity-100"
            >
              <MagnifyingGlassPlusIcon class="h-5 w-5 text-white" />
            </div>
          </div>
          <button
            type="button"
            title="Ganti Foto"
            @click="fotoInput.click()"
            class="absolute -right-3 -bottom-3 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700 dark:hover:text-white"
          >
            <PencilSquareIcon class="h-4 w-4" />
          </button>
          <button
            v-if="isEdit && form.foto && form.foto !== FOTO_DEFAULT"
            type="button"
            title="Hapus Foto"
            @click="showDeleteFotoDialog = true"
            class="absolute -left-3 -bottom-3 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 hover:bg-red-50 hover:text-red-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-red-500/10 dark:hover:text-red-400"
          >
            <XMarkIcon class="h-4 w-4" />
          </button>
          <input
            ref="fotoInput"
            type="file"
            accept=".jpg,.jpeg,.png,.webp"
            class="hidden"
            @change="handleFotoChange"
          />
        </div>
      </div>
    </section>

    <!-- Data Akademik (relasi kelas/rombel/tahun pelajaran) -->
    <section :class="cardClass">
      <SectionHeading
        title="Data Akademik"
        description="Penempatan tahun pelajaran, kelas, dan rombel peserta didik."
      >
        <AcademicCapIcon class="h-5 w-5" />
      </SectionHeading>
      <div class="grid grid-cols-1 gap-6 p-6 sm:grid-cols-3">
        <FormField label="Tahun Pelajaran *">
          <select v-model.number="form.id_tahunpelajaran" :class="inputClass">
            <option :value="0">Pilih</option>
            <option
              v-for="tp in tahunPelajaranOptions"
              :key="tp.id_tahunpelajaran"
              :value="tp.id_tahunpelajaran"
            >
              {{ tp.nama_tahunpelajaran }}
            </option>
          </select>
        </FormField>
        <FormField label="Kelas *">
          <select v-model.number="form.id_kelas" :class="inputClass">
            <option :value="0">Pilih</option>
            <option
              v-for="k in kelasOptions"
              :key="k.id_kelas"
              :value="k.id_kelas"
            >
              {{ k.nama_kelas }}
            </option>
          </select>
        </FormField>
        <FormField label="Rombel *">
          <select v-model.number="form.id_rombel" :class="inputClass">
            <option :value="0">Pilih</option>
            <option
              v-for="r in rombelOptionsForKelas"
              :key="r.id_rombel"
              :value="r.id_rombel"
            >
              {{ r.nama_rombel }}
            </option>
          </select>
        </FormField>
      </div>
    </section>

    <!-- Section dinamis -->
    <section
      v-for="section in sections"
      :key="section.title"
      :class="cardClass"
    >
      <SectionHeading :title="section.title" :description="section.description">
        <component :is="section.icon" class="h-5 w-5" />
      </SectionHeading>
      <div class="grid grid-cols-1 gap-6 p-6 sm:grid-cols-2">
        <FormField
          v-for="field in section.fields"
          :key="field.key"
          :label="field.label + (field.required ? ' *' : '')"
          :class="field.full ? 'sm:col-span-2' : ''"
        >
          <select
            v-if="field.type === 'select'"
            v-model="form[field.key]"
            :class="inputClass"
          >
            <option value="">Pilih</option>
            <option
              v-for="opt in field.options"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.label }}
            </option>
          </select>
          <textarea
            v-else-if="field.type === 'textarea'"
            v-model="form[field.key]"
            rows="3"
            :class="inputClass"
            placeholder="Ketik disini..."
          />
          <input
            v-else
            :type="
              field.type === 'date'
                ? 'date'
                : field.type === 'number'
                  ? 'number'
                  : 'text'
            "
            v-model="form[field.key]"
            :class="[
              inputClass,
              field.type === 'date'
                ? 'dark:[&::-webkit-calendar-picker-indicator]:invert'
                : '',
            ]"
            :disabled="field.key === 'nis' && isEdit"
            :placeholder="field.type === 'date' ? '' : 'Ketik disini...'"
          />
          <p
            v-if="fieldErrors[field.key]"
            class="mt-1 text-xs text-red-600 dark:text-red-400"
          >
            {{ fieldErrors[field.key] }}
          </p>
          <p
            v-else-if="duplicateWarnings[field.key]"
            class="mt-1 text-xs text-yellow-600 dark:text-yellow-400"
          >
            Data ganda a.n. {{ duplicateWarnings[field.key].nama_siswa }}
          </p>
        </FormField>
      </div>
    </section>

    <!-- Aksi -->
    <div
      class="flex flex-col gap-4 rounded-2xl bg-white px-5 py-4 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5 md:flex-row md:items-center md:justify-between"
    >
      <div class="flex gap-3">
        <button
          type="submit"
          :disabled="loading"
          class="inline-flex items-center justify-center gap-2 rounded-xl bg-blue-600 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:bg-gray-400 dark:disabled:bg-slate-700"
        >
          <ArrowPathIcon v-if="loading" class="h-5 w-5 animate-spin" />
          {{ loading ? "Menyimpan..." : "Simpan" }}
        </button>
        <RouterLink
          :to="{ name: 'dashboard-pesertadidik' }"
          class="inline-flex items-center justify-center rounded-xl border border-gray-300 px-5 py-2.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
        >
          Batal
        </RouterLink>
      </div>
    </div>
  </form>
  <DashboardConfirmDialog
    v-model="showDeleteFotoDialog"
    title="Hapus Foto"
    message="Foto akan dikembalikan ke foto default. Lanjutkan?"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="deletingFoto"
    @confirm="handleDeleteFoto"
  />
  <VueEasyLightbox
    :visible="visibleLightbox"
    :imgs="images"
    :index="0"
    @hide="visibleLightbox = false"
  />
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";
import FormField from "@/components/form/FormField.vue";
import { useToast } from "@/composables/useToast";
import { usePesertaDidikStore } from "@/stores/pesertadidik";
import { useTahunPelajaranStore } from "@/stores/tahunpelajaran";
import { useKelasStore } from "@/stores/kelas";
import { useRombelStore } from "@/stores/rombel";
import VueEasyLightbox from "vue-easy-lightbox";
import {
  AcademicCapIcon,
  ArrowPathIcon,
  BanknotesIcon,
  DocumentTextIcon,
  HeartIcon,
  IdentificationIcon,
  MapPinIcon,
  PencilSquareIcon,
  PhotoIcon,
  UserGroupIcon,
  UserIcon,
  XMarkIcon,
  MagnifyingGlassPlusIcon,
} from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const toast = useToast();

const pdStore = usePesertaDidikStore();
const tahunPelajaranStore = useTahunPelajaranStore();
const kelasStore = useKelasStore();
const rombelStore = useRombelStore();

const showDeleteFotoDialog = ref(false);
const deletingFoto = ref(false);
const FOTO_DEFAULT = "1636812404_b39deec2f158cb789567.jpg";

const isEdit = !!route.params.nis;
const fotoInput = ref(null);
const selectedFoto = ref(null);

const visibleLightbox = ref(false);
const images = computed(() => (form.foto_url ? [form.foto_url] : []));

const loading = ref(false);

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";
const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";

const tahunPelajaranOptions = computed(() =>
  (tahunPelajaranStore.list || []).filter((tp) => tp.aktif === 1),
);

const kelasOptions = computed(() =>
  (kelasStore.list || []).filter((k) => k.aktif === 1),
);
const rombelOptionsForKelas = computed(() => {
  const all = rombelStore.list || [];
  const filtered = all.filter(
    (r) => !form.id_kelas || r.id_kelas === form.id_kelas,
  );
  const aktifList = filtered.filter((r) => r.aktif === 1);

  const current = all.find((r) => r.id_rombel === form.id_rombel);
  if (current && current.aktif !== 1) {
    return [current, ...aktifList];
  }
  return aktifList;
});

// ===== Opsi enum (value/label, konsisten sama pola pegawai) =====
const OPT_JK = [
  { value: "L", label: "Laki-laki" },
  { value: "P", label: "Perempuan" },
];
const OPT_AGAMA = [
  "Islam",
  "Kristen",
  "Katholik",
  "Hindu",
  "Budha",
  "Khonghucu",
  "Kepercayaan kpd Tuhan YME",
  "lainnya",
].map((v) => ({ value: v, label: v }));
const OPT_WARGA = ["Indonesia", "Asing"].map((v) => ({
  value: v,
  label: v,
}));
const OPT_GOLDARAH = ["A", "B", "AB", "O"].map((v) => ({ value: v, label: v }));
const OPT_TINGGAL = [
  "Bersama orang tua",
  "Wali",
  "Kost",
  "Asrama",
  "Panti asuhan",
  "Pesantren",
  "Lainnya",
].map((v) => ({ value: v, label: v }));
const OPT_PENDIDIKAN = [
  { value: "Tidak sekolah", label: "Tidak sekolah" },
  { value: "TK / sederajat", label: "TK / sederajat" },
  { value: "SD / sederajat", label: "SD / sederajat" },
  { value: "SMP / sederajat", label: "SMP / sederajat" },
  { value: "SMA / sederajat", label: "SMA / sederajat" },
  { value: "D1", label: "D1" },
  { value: "D2", label: "D2" },
  { value: "D3", label: "D3" },
  { value: "D4", label: "D4" },
  { value: "S1", label: "S1" },
  { value: "S2", label: "S2" },
  { value: "S2 Terapan", label: "S2 Terapan" },
  { value: "S3", label: "S3" },
  { value: "S3 Terapan", label: "S3 Terapan" },
  { value: "Sp-1", label: "Sp-1" },
  { value: "Sp-2", label: "Sp-2" },
  { value: "PAUD", label: "PAUD" },
  { value: "Profesi", label: "Profesi" },
  { value: "Paket A", label: "Paket A" },
  { value: "Paket B", label: "Paket B" },
  { value: "Paket C", label: "Paket C" },
  { value: "Informal", label: "Informal" },
  { value: "Non formal", label: "Non formal" },
  { value: "Putus SD", label: "Putus SD" },
  { value: "Lainnya", label: "Lainnya" },
];
const OPT_BANTUAN = ["KJP", "KIP", "BSM"].map((v) => ({ value: v, label: v }));
const OPT_YATIM = ["Yatim", "Piatu", "Yatim Piatu"].map((v) => ({
  value: v,
  label: v,
}));
const OPT_BAHASA = ["Indonesia", "Daerah", "Asing"].map((v) => ({
  value: v,
  label: v,
}));
const OPT_TRANSPORTASI = [
  "Jalan kaki",
  "Angkutan umum/bus/pete-pete",
  "Mobil/bus antar jemput",
  "Kereta api",
  "Ojek",
  "Andong/bendi/sado/dokar/delman/becak",
  "Perahu penyeberangan/rakit/getek",
  "Kuda",
  "Sepeda",
  "Sepeda motor",
  "Mobil pribadi",
  "Lainnya",
].map((v) => ({ value: v, label: v }));
const OPT_JASMANI = ["Normal", "Cacat"].map((v) => ({ value: v, label: v }));
const OPT_KEBUTUHAN_KHUSUS = [
  "Tidak ada",
  "netra (A)",
  "rungu (B)",
  "grahita ringan (C)",
  "grahita sedang (C1)",
  "daksa ringan (D)",
  "daksa sedang (D1)",
  "laras (E)",
  "wicara (F)",
  "hyperaktif (H)",
  "cerdas istimewa (I)",
  "bakat istimewa (J)",
  "kesulitan belajar (K)",
  "narkoba (N)",
  "indigo (O)",
  "down syndrome (P)",
  "autis (Q)",
].map((v) => ({ value: v, label: v }));
const OPT_JALUR = [
  "PPDB",
  "Prestasi Akademik",
  "Prestasi Non Akademik",
  "Prestasi Akademik Dalam",
  "Prestasi Akademik Luar",
  "Afirmasi",
  "Anak Panti",
  "Inklusi",
  "Domisili",
  "Mutasi",
  "Tahap Kedua",
  "Pindahan",
].map((v) => ({ value: v, label: v }));
const OPT_SEMESTER = ["1", "2", "3", "4", "5", "6"].map((v) => ({
  value: v,
  label: v,
}));
const OPT_PEKERJAAN = [
  "Tidak bekerja",
  "Nelayan",
  "Petani",
  "Peternak",
  "PNS/TNI/Polri",
  "Karyawan Swasta",
  "Pedagang Kecil",
  "Pedagang Besar",
  "Wiraswasta",
  "Wirausaha",
  "Buruh",
  "Pensiunan",
  "Tenaga Kerja Indonesia",
  "Karyawan BUMN",
  "Tidak dapat diterapkan",
  "Sudah Meninggal",
  "Lainnya",
].map((v) => ({ value: v, label: v }));
const OPT_KATEGORI_PENGHASILAN = [
  "Tidak Pasti",
  "Kurang dari Rp. 500,000",
  "Rp. 500,000 - Rp. 999,999",
  "Rp. 1,000,000 - Rp. 1,999,999",
  "Rp. 2,000,000 - Rp. 4,999,999",
  "Rp. 5,000,000 - Rp. 20,000,000",
  "Lebih dari Rp. 20,000,000",
  "Tidak Berpenghasilan",
].map((v) => ({ value: v, label: v }));
const OPT_STATUS_ORTU = ["Hidup", "Meninggal"].map((v) => ({
  value: v,
  label: v,
}));

// ===== Konfigurasi field per section (data-driven) =====
const sections = [
  {
    title: "Data Siswa",
    description: "Identitas dan data pribadi siswa.",
    icon: IdentificationIcon,
    fields: [
      { key: "nis", label: "NIS", required: true },
      {
        key: "nama_siswa",
        label: "Nama Peserta Didik",
        required: true,
        full: true,
      },
      { key: "nisn", label: "NISN" },
      { key: "no_ujiansd", label: "No. Sidanira" },
      { key: "no_simdik", label: "No. Simdik" },
      {
        key: "jenis_kelamin",
        label: "Jenis Kelamin",
        type: "select",
        options: OPT_JK,
        required: true,
      },
      { key: "tempat_lahirsiswa", label: "Tempat Lahir", required: true },
      {
        key: "tanggal_lahirsiswa",
        label: "Tanggal Lahir",
        type: "date",
        required: true,
      },
      { key: "no_akte", label: "No. Akte" },
      { key: "akte_kota", label: "Kota Akte" },
      {
        key: "agama_siswa",
        label: "Agama",
        type: "select",
        options: OPT_AGAMA,
        required: true,
      },
      {
        key: "warga_siswa",
        label: "Kewarganegaraan",
        type: "select",
        options: OPT_WARGA,
      },
      { key: "anak_ke", label: "Anak Ke", type: "number" },
      { key: "saudara_kandung", label: "Saudara Kandung", type: "number" },
      { key: "saudara_tiri", label: "Saudara Tiri", type: "number" },
      { key: "saudara_angkat", label: "Saudara Angkat", type: "number" },
      {
        key: "yatim_piatu",
        label: "Yatim/Piatu",
        type: "select",
        options: OPT_YATIM,
      },
      {
        key: "bahasa_siswa",
        label: "Bahasa Sehari-hari",
        type: "select",
        options: OPT_BAHASA,
      },
      {
        key: "bantuan",
        label: "Penerima Bantuan",
        type: "select",
        options: OPT_BANTUAN,
      },
      { key: "no_kk", label: "No. KK" },
      { key: "nik_siswa", label: "NIK Siswa" },
      { key: "email", label: "Email", type: "email" },
      {
        key: "gol_darah",
        label: "Golongan Darah",
        type: "select",
        options: OPT_GOLDARAH,
      },
      { key: "pernah_sakit", label: "Riwayat Sakit" },
      {
        key: "status_jasmani",
        label: "Status Jasmani",
        type: "select",
        options: OPT_JASMANI,
      },
      {
        key: "kebutuhan_khusus",
        label: "Kebutuhan Khusus",
        type: "select",
        options: OPT_KEBUTUHAN_KHUSUS,
      },
      { key: "tinggi_badan", label: "Tinggi Badan (cm)", type: "number" },
      { key: "berat_badan", label: "Berat Badan (kg)", type: "number" },
      { key: "lingkar_kepala", label: "Lingkar Kepala (cm)", type: "number" },
    ],
  },
  {
    title: "Alamat & Kontak Siswa",
    description: "Alamat tempat tinggal dan kontak siswa.",
    icon: MapPinIcon,
    fields: [
      { key: "alamat_siswa", label: "Alamat", type: "textarea", full: true },
      { key: "rt_siswa", label: "RT" },
      { key: "rw_siswa", label: "RW" },
      { key: "kel_siswa", label: "Kelurahan/Desa" },
      { key: "kec_siswa", label: "Kecamatan" },
      { key: "kopos_siswa", label: "Kode Pos" },
      { key: "telprmh_siswa", label: "Telp Rumah" },
      { key: "hp_siswa", label: "No. HP" },
      {
        key: "tinggal_siswa",
        label: "Tinggal Bersama",
        type: "select",
        options: OPT_TINGGAL,
      },
      { key: "jarak_sekolah", label: "Jarak ke Sekolah (km)" },
      {
        key: "moda_transportasi",
        label: "Moda Transportasi",
        type: "select",
        options: OPT_TRANSPORTASI,
      },
    ],
  },
  {
    title: "Akademik & Kelulusan",
    description: "Data asal sekolah, kelulusan, dan ijazah.",
    icon: AcademicCapIcon,
    fields: [
      { key: "lulus_sd", label: "Lulus Dari (SD)" },
      { key: "tanggal_ijazah", label: "Tanggal Ijazah", type: "date" },
      { key: "kota_ijazah", label: "Kota Ijazah" },
      { key: "nomor_ijazah", label: "Nomor Ijazah" },
      { key: "tahun_ijazah", label: "Tahun Ijazah" },
      { key: "nomor_skhun", label: "Nomor SKHUN" },
      { key: "tahun_skhun", label: "Tahun SKHUN" },
      { key: "nilai_ijazah", label: "Nilai Ijazah" },
      { key: "nilai_nem", label: "Nilai NEM" },
      { key: "lama_belajar", label: "Lama Belajar (tahun)", type: "number" },
      {
        key: "jalur",
        label: "Jalur Masuk",
        type: "select",
        options: OPT_JALUR,
      },
      {
        key: "semester",
        label: "Semester",
        type: "select",
        options: OPT_SEMESTER,
      },
      { key: "asal_smp", label: "Asal Sekolah" },
      { key: "tanggal_masuk", label: "Tanggal Masuk", type: "date" },
    ],
  },
  {
    title: "Data Ayah",
    description: "Identitas dan data orang tua (ayah).",
    icon: UserIcon,
    fields: [
      { key: "nama_ayah", label: "Nama Ayah", full: true },
      { key: "nik_ayah", label: "NIK Ayah" },
      { key: "tempat_lahirayah", label: "Tempat Lahir" },
      { key: "tanggal_lahirayah", label: "Tanggal Lahir", type: "date" },
      { key: "agama_ayah", label: "Agama", type: "select", options: OPT_AGAMA },
      {
        key: "warga_ayah",
        label: "Kewarganegaraan",
        type: "select",
        options: OPT_WARGA,
      },
      {
        key: "pend_ayah",
        label: "Pendidikan Terakhir",
        type: "select",
        options: OPT_PENDIDIKAN,
      },
      {
        key: "pek_ayah",
        label: "Pekerjaan",
        type: "select",
        options: OPT_PEKERJAAN,
      },
      { key: "gaji_ayah", label: "Penghasilan" },
      {
        key: "gaji_ayahktg",
        label: "Kategori Penghasilan",
        type: "select",
        options: OPT_KATEGORI_PENGHASILAN,
      },
      {
        key: "status_ayah",
        label: "Status",
        type: "select",
        options: OPT_STATUS_ORTU,
      },
      { key: "alamat_ayah", label: "Alamat", type: "textarea", full: true },
      { key: "rt_ayah", label: "RT" },
      { key: "rw_ayah", label: "RW" },
      { key: "kopos_ayah", label: "Kode Pos" },
      { key: "kel_ayah", label: "Kelurahan/Desa" },
      { key: "kec_ayah", label: "Kecamatan" },
      { key: "telp_ayah", label: "No. Telp/HP" },
    ],
  },
  {
    title: "Data Ibu",
    description: "Identitas dan data orang tua (ibu).",
    icon: UserIcon,
    fields: [
      { key: "nama_ibu", label: "Nama Ibu", full: true },
      { key: "nik_ibu", label: "NIK Ibu" },
      { key: "tempat_lahiribu", label: "Tempat Lahir" },
      { key: "tanggal_lahiribu", label: "Tanggal Lahir", type: "date" },
      { key: "agama_ibu", label: "Agama", type: "select", options: OPT_AGAMA },
      {
        key: "warga_ibu",
        label: "Kewarganegaraan",
        type: "select",
        options: OPT_WARGA,
      },
      {
        key: "pend_ibu",
        label: "Pendidikan Terakhir",
        type: "select",
        options: OPT_PENDIDIKAN,
      },
      {
        key: "pek_ibu",
        label: "Pekerjaan",
        type: "select",
        options: OPT_PEKERJAAN,
      },
      { key: "gaji_ibu", label: "Penghasilan" },
      {
        key: "gaji_ibuktg",
        label: "Kategori Penghasilan",
        type: "select",
        options: OPT_KATEGORI_PENGHASILAN,
      },
      {
        key: "status_ibu",
        label: "Status",
        type: "select",
        options: OPT_STATUS_ORTU,
      },
      { key: "alamat_ibu", label: "Alamat", type: "textarea", full: true },
      { key: "rt_ibu", label: "RT" },
      { key: "rw_ibu", label: "RW" },
      { key: "kopos_ibu", label: "Kode Pos" },
      { key: "kel_ibu", label: "Kelurahan/Desa" },
      { key: "kec_ibu", label: "Kecamatan" },
      { key: "telp_ibu", label: "No. Telp/HP" },
    ],
  },
  {
    title: "Data Wali",
    description: "Identitas dan data wali siswa (jika ada).",
    icon: UserGroupIcon,
    fields: [
      { key: "nama_wali", label: "Nama Wali", full: true },
      { key: "tempat_lahirwali", label: "Tempat Lahir" },
      { key: "tanggal_lahirwali", label: "Tanggal Lahir", type: "date" },
      { key: "agama_wali", label: "Agama", type: "select", options: OPT_AGAMA },
      {
        key: "warga_wali",
        label: "Kewarganegaraan",
        type: "select",
        options: OPT_WARGA,
      },
      {
        key: "pend_wali",
        label: "Pendidikan Terakhir",
        type: "select",
        options: OPT_PENDIDIKAN,
      },
      {
        key: "pek_wali",
        label: "Pekerjaan",
        type: "select",
        options: OPT_PEKERJAAN,
      },
      { key: "gaji_wali", label: "Penghasilan" },
      {
        key: "gaji_waliktg",
        label: "Kategori Penghasilan",
        type: "select",
        options: OPT_KATEGORI_PENGHASILAN,
      },
      { key: "alamat_wali", label: "Alamat", type: "textarea", full: true },
      { key: "rt_wali", label: "RT" },
      { key: "rw_wali", label: "RW" },
      { key: "kopos_wali", label: "Kode Pos" },
      { key: "kel_wali", label: "Kelurahan/Desa" },
      { key: "kec_wali", label: "Kecamatan" },
      { key: "telp_wali", label: "No. Telp/HP" },
    ],
  },
  {
    title: "Hobi & Minat",
    description: "Minat dan bakat siswa di luar akademik.",
    icon: HeartIcon,
    fields: [
      { key: "hobi_seni", label: "Seni" },
      { key: "hobi_orkes", label: "Olahraga" },
      { key: "hobi_organisasi", label: "Organisasi" },
      { key: "hobi_lain", label: "Lainnya" },
    ],
  },
];

// ===== Form state =====
const form = reactive({
  id_tahunpelajaran: 0,
  id_kelas: 0,
  id_rombel: 0,
  foto: "",
  foto_url: "", // client-only, tidak dikirim
});
sections.forEach((section) =>
  section.fields.forEach((field) => (form[field.key] = "")),
);

//duplikat warning
const duplicateWarnings = reactive({
  nis: null,
  nik_siswa: null,
  nik_ayah: null,
  nik_ibu: null,
  nisn: null,
  no_ujiansd: null,
});

const fieldErrors = reactive({
  nis: "",
  nisn: "",
  no_ujiansd: "",
  nik_siswa: "",
  nik_ayah: "",
  nik_ibu: "",
  no_kk: "",
});

// exact: true = harus PERSIS sejumlah digit. exact: false = MAKSIMAL sejumlah digit.
const digitRules = {
  nis: { length: 12, exact: false },
  nisn: { length: 10, exact: true },
  no_ujiansd: { length: 22, exact: false },
  nik_siswa: { length: 16, exact: true },
  nik_ayah: { length: 16, exact: true },
  nik_ibu: { length: 16, exact: true },
  no_kk: { length: 16, exact: true },
};

function validateDigitLocal(key, value) {
  const rule = digitRules[key];
  if (!rule) return true;
  if (!value) {
    fieldErrors[key] = "";
    return true;
  }
  if (!/^\d+$/.test(value)) {
    fieldErrors[key] = "Harus berupa angka.";
    return false;
  }
  if (rule.exact && value.length !== rule.length) {
    fieldErrors[key] = `Harus ${rule.length} digit angka.`;
    return false;
  }
  if (!rule.exact && value.length > rule.length) {
    fieldErrors[key] = `Maksimal ${rule.length} digit angka.`;
    return false;
  }
  fieldErrors[key] = "";
  return true;
}

const debounceTimers = {};

async function checkDuplicateField(key, value) {
  try {
    const exclude = isEdit ? route.params.nis : "";
    const result = await pdStore.checkDuplicate(key, value, exclude);
    if (result.exists) duplicateWarnings[key] = result;
  } catch {
    // non-blocking, diamkan kalau gagal
  }
}

Object.keys(duplicateWarnings).forEach((key) => {
  watch(
    () => form[key],
    (value) => {
      if (key === "nis" && isEdit) return;
      clearTimeout(debounceTimers[key]);
      duplicateWarnings[key] = null;

      if (!validateDigitLocal(key, value)) return; // format salah, jangan lanjut cek duplikat

      if (!value) return;
      debounceTimers[key] = setTimeout(
        () => checkDuplicateField(key, value),
        300,
      );
    },
  );
});

watch(
  () => form.no_kk,
  (value) => validateDigitLocal("no_kk", value),
);

// ===== Load data (mode edit) =====
async function loadData() {
  if (!isEdit) return;
  loading.value = true;
  try {
    const data = await pdStore.fetchPesertaDidikByNis(route.params.nis);
    Object.assign(form, data, { foto_url: data.foto_url || "" });
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat data peserta didik.",
    );
  } finally {
    loading.value = false;
  }
}

async function handleDeleteFoto() {
  deletingFoto.value = true;
  try {
    const res = await pdStore.deleteFotoPesertaDidik(route.params.nis);
    form.foto = "";
    form.foto_url = res.foto_url;
    selectedFoto.value = null;
    toast.success("Foto berhasil dihapus.");
    showDeleteFotoDialog.value = false;
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menghapus foto.");
  } finally {
    deletingFoto.value = false;
  }
}

// ===== Foto =====
function handleFotoChange(event) {
  const file = event.target.files?.[0];
  if (!file) return;
  const maxSize = 2 * 1024 * 1024;
  if (file.size > maxSize) {
    toast.error("Ukuran foto maksimal 2 MB.");
    event.target.value = "";
    return;
  }
  const allowed = ["image/jpeg", "image/png", "image/webp"];
  if (!allowed.includes(file.type)) {
    toast.error("Foto harus berformat JPG, JPEG, PNG atau WEBP.");
    event.target.value = "";
    return;
  }
  selectedFoto.value = file;
  form.foto_url = URL.createObjectURL(file);
}

// ===== Submit =====
function buildPayload() {
  const { foto_url, foto, ...rest } = form;
  const payload = { ...rest };

  // Field number (anak_ke, tinggi_badan, dll) otomatis jadi tipe Number oleh v-model
  // pada <input type="number">, padahal backend expect string (kolom varchar).
  // Stringify semua kecuali relasi id_* yang memang harus int.
  const intFields = ["id_tahunpelajaran", "id_kelas", "id_rombel"];
  Object.keys(payload).forEach((key) => {
    if (intFields.includes(key)) return;
    if (payload[key] !== null && payload[key] !== undefined) {
      payload[key] = String(payload[key]);
    }
  });

  return payload;
}
async function submit() {
  if (loading.value) return;
  if (!isEdit && !form.nis.trim()) {
    toast.error("NIS wajib diisi.");
    return;
  }
  if (!form.nama_siswa.trim()) {
    toast.error("Nama siswa wajib diisi.");
    return;
  }
  if (!form.id_tahunpelajaran || !form.id_kelas || !form.id_rombel) {
    toast.error("Tahun pelajaran, kelas, dan rombel wajib dipilih.");
    return;
  }
  const digitKeys = Object.keys(digitRules);
  for (const key of digitKeys) {
    if (!validateDigitLocal(key, form[key])) {
      toast.error(
        "Periksa kembali NIS/NISN/No. Ujian SD/NIK — format digit tidak sesuai.",
      );
      return;
    }
  }
  if (Object.values(duplicateWarnings).some((w) => w)) {
    toast.error(
      "Ada data yang ganda dengan siswa lain, silakan periksa kembali.",
    );
    return;
  }
  loading.value = true;
  try {
    const payload = buildPayload();
    let nis = route.params.nis;
    let dataChanged = true; // default true untuk mode tambah

    if (isEdit) {
      const response = await pdStore.updatePesertaDidik(nis, payload);
      dataChanged = !response.changed || response.changed.length > 0;
    } else {
      await pdStore.createPesertaDidik(payload);
      nis = form.nis;
    }

    if (selectedFoto.value) {
      const fd = new FormData();
      fd.append("foto", selectedFoto.value);
      await pdStore.updateFotoPesertaDidik(nis, fd);
    }

    // Tentukan pesan akhir berdasarkan kombinasi: data teks berubah? foto diganti?
    if (!isEdit) {
      toast.success("Peserta didik berhasil ditambahkan.");
    } else if (dataChanged && selectedFoto.value) {
      toast.success("Data dan foto peserta didik berhasil diperbarui.");
    } else if (dataChanged) {
      toast.success("Data peserta didik berhasil diperbarui.");
    } else if (selectedFoto.value) {
      toast.success("Foto peserta didik berhasil diperbarui.");
    } else {
      toast.warning("Tidak ada perubahan data.");
    }

    router.push({ name: "dashboard-pesertadidik" });
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menyimpan data peserta didik.",
    );
  } finally {
    loading.value = false;
  }
}

function showImage() {
  visibleLightbox.value = true;
}

onMounted(async () => {
  tahunPelajaranStore.pagination.limit = 50;
  kelasStore.pagination.limit = 100;
  rombelStore.pagination.limit = 500;

  tahunPelajaranStore.fetchAll();
  kelasStore.fetchAll();
  rombelStore.fetchAll();

  if (!isEdit && tahunPelajaranStore.aktif) {
    form.id_tahunpelajaran = tahunPelajaranStore.aktif.id_tahunpelajaran;
  }

  await loadData();
});
</script>
