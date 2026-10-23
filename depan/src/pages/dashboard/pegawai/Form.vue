<template>
  <DashboardPageHeader
    :title="isEdit ? 'Ubah Pegawai' : 'Tambah Pegawai'"
    subtitle="Kelola data kepegawaian sekolah."
  />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <!-- Foto -->
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Foto Pegawai"
        description="Foto profil pegawai yang tampil pada daftar dan detil."
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
              alt="Foto Pegawai"
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
            v-if="isEdit && form.foto"
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

    <!-- Section dinamis -->
    <section
      v-for="section in visibleSections"
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
            :type="field.type || 'text'"
            v-model="form[field.key]"
            :class="[
              inputClass,
              field.type === 'date'
                ? 'dark:[&::-webkit-calendar-picker-indicator]:invert'
                : '',
            ]"
            :disabled="field.key === 'nuptk' && isEdit"
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
            Data ganda a.n. {{ duplicateWarnings[field.key].nama_pegawai }}
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
          :to="{ name: 'dashboard-pegawai' }"
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
    message="Foto pegawai akan dihapus. Lanjutkan?"
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
import FormField from "@/components/form/FormField.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";
import { useToast } from "@/composables/useToast";
import { usePegawaiStore } from "@/stores/pegawai";
import VueEasyLightbox from "vue-easy-lightbox";

import {
  ArrowPathIcon,
  ArrowRightOnRectangleIcon,
  BanknotesIcon,
  DocumentTextIcon,
  IdentificationIcon,
  MapPinIcon,
  PencilSquareIcon,
  PhoneIcon,
  PhotoIcon,
  UserIcon,
  UserGroupIcon,
  XMarkIcon,
  MagnifyingGlassPlusIcon,
} from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const pegawaiStore = usePegawaiStore();
const toast = useToast();

const isEdit = !!route.params.nuptk;

const fotoInput = ref(null);
const selectedFoto = ref(null);
const loading = ref(false);

const showDeleteFotoDialog = ref(false);
const deletingFoto = ref(false);
const visibleLightbox = ref(false);
const images = computed(() => (form.foto_url ? [form.foto_url] : []));

const fieldErrors = reactive({
  nuptk: "",
  nip_baru: "",
  nik: "",
  no_kk: "",
});

const duplicateWarnings = reactive({
  nuptk: null,
  nip_baru: null,
  nik: null,
});

const digitRules = { nuptk: 16, nip_baru: 18, nik: 16, no_kk: 16 };
const debounceTimers = {};

const visibleSections = computed(() =>
  sections.filter((section) => !(section.hideOnCreate && !isEdit)),
);

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";

const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";

const form = reactive({
  nuptk: "",
  nama_pegawai: "",
  jabatan: "",
  strata: "",
  pns_non: "",
  nikki: "",
  simpkbid: "",
  dapodikid: "",
  nip_lama: "",
  nip_baru: "",
  nrk: "",
  nrg: "",
  tempat_lahir: "",
  tanggal_lahir: "",
  jenis_kelamin: "",
  agama: "",
  nik: "",
  no_kk: "",
  alamat: "",
  rt: "",
  rw: "",
  kel: "",
  kec: "",
  kopos: "",
  telp_rmh: "",
  telp_hp: "",
  email: "",
  sk_sekolah: "",
  tmt_sekolah: "",
  sk_nonkki: "",
  tmt_nonkki: "",
  sk_kki: "",
  tmt_kki: "",
  sk_cpns: "",
  tmt_cpns: "",
  sk_pns: "",
  tmt_pns: "",
  tmt_kgb: "",
  pangkat_kgb: "",
  gaji_kgb: "",
  sk_kgb: "",
  tglsk_kgb: "",
  karpeg: "",
  karis: "",
  nama_ibukdg: "",
  nm_suami_istri: "",
  pek_suami_istri: "",
  npwp: "",
  status_keluar: "",
  keterangan_keluar: "",
  tanggal_keluar: "",
  nomor_surat: "",
  foto: "",
  foto_url: "", // client-only, tidak dikirim
});

const sections = [
  {
    title: "Identitas Utama",
    description: "Data pokok identitas pegawai.",
    icon: IdentificationIcon,
    fields: [
      { key: "nuptk", label: "NUPTK", required: true },
      { key: "nip_baru", label: "NIP" },
      {
        key: "nama_pegawai",
        label: "Nama Pegawai",
        required: true,
        full: true,
      },
      {
        key: "jabatan",
        label: "Jabatan",
        required: true,
        type: "select",
        options: [
          { value: "Guru", label: "Guru" },
          { value: "Tatausaha", label: "Tatausaha" },
        ],
      },
      {
        key: "strata",
        label: "Pendidikan Terakhir",
        type: "select",
        options: [
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
        ],
      },
      {
        key: "pns_non",
        label: "Status Kepegawaian",
        type: "select",
        options: [
          { value: "PNS", label: "PNS" },
          { value: "CPNS", label: "CPNS" },
          { value: "PPPK", label: "PPPK" },
          { value: "PPPK-PW", label: "PPPK-PW" },
          { value: "KKI", label: "KKI" },
          { value: "Non KKI", label: "Non KKI" },
          { value: "Guru Honor Sekolah", label: "Guru Honor Sekolah" },
          {
            value: "Honor Daerah TK.I Provinsi",
            label: "Honor Daerah TK.I Provinsi",
          },
        ],
      },
    ],
  },
  {
    title: "Data Kepegawaian",
    description: "Nomor induk dan identitas kepegawaian.",
    icon: DocumentTextIcon,
    fields: [
      { key: "nikki", label: "NIKKI" },
      { key: "simpkbid", label: "SIMPKB ID" },
      { key: "dapodikid", label: "Dapodik ID" },
      { key: "nip_lama", label: "NIP Lama" },
      { key: "nrk", label: "NRK" },
      { key: "nrg", label: "NRG" },
    ],
  },
  {
    title: "Data Pribadi",
    description: "Data diri pegawai.",
    icon: UserIcon,
    fields: [
      { key: "tempat_lahir", label: "Tempat Lahir" },
      { key: "tanggal_lahir", label: "Tanggal Lahir", type: "date" },
      {
        key: "jenis_kelamin",
        label: "Jenis Kelamin",
        type: "select",
        options: [
          { value: "L", label: "Laki-laki" },
          { value: "P", label: "Perempuan" },
        ],
      },
      {
        key: "agama",
        label: "Agama",
        type: "select",
        options: [
          { value: "Islam", label: "Islam" },
          { value: "Kristen", label: "Kristen" },
          { value: "Katholik", label: "Katholik" },
          { value: "Hindu", label: "Hindu" },
          { value: "Budha", label: "Budha" },
          { value: "Khonghucu", label: "Khonghucu" },
          {
            value: "Kepercayaan kpd Tuhan YME",
            label: "Kepercayaan kpd Tuhan YME",
          },
          { value: "lainnya", label: "lainnya" },
        ],
      },
      { key: "nik", label: "NIK" },
      { key: "no_kk", label: "No. KK" },
    ],
  },
  {
    title: "Alamat",
    description: "Alamat tempat tinggal pegawai.",
    icon: MapPinIcon,
    fields: [
      { key: "alamat", label: "Alamat", type: "textarea", full: true },
      { key: "rt", label: "RT" },
      { key: "rw", label: "RW" },
      { key: "kel", label: "Kelurahan" },
      { key: "kec", label: "Kecamatan" },
      { key: "kopos", label: "Kode Pos" },
    ],
  },
  {
    title: "Kontak",
    description: "Nomor telepon dan email yang bisa dihubungi.",
    icon: PhoneIcon,
    fields: [
      { key: "telp_rmh", label: "Telp. Rumah" },
      { key: "telp_hp", label: "No. HP" },
      { key: "email", label: "Email", type: "email" },
    ],
  },
  {
    title: "Riwayat SK & TMT",
    description: "Surat keputusan dan terhitung mulai tanggal.",
    icon: DocumentTextIcon,
    fields: [
      { key: "sk_sekolah", label: "SK Sekolah" },
      { key: "tmt_sekolah", label: "TMT Sekolah", type: "date" },
      { key: "sk_nonkki", label: "SK Non KKI" },
      { key: "tmt_nonkki", label: "TMT Non KKI", type: "date" },
      { key: "sk_kki", label: "SK KKI" },
      { key: "tmt_kki", label: "TMT KKI", type: "date" },
      { key: "sk_cpns", label: "SK CPNS" },
      { key: "tmt_cpns", label: "TMT CPNS", type: "date" },
      { key: "sk_pns", label: "SK PNS" },
      { key: "tmt_pns", label: "TMT PNS", type: "date" },
    ],
  },
  {
    title: "Kenaikan Gaji Berkala (KGB)",
    description: "Riwayat kenaikan gaji berkala terakhir.",
    icon: BanknotesIcon,
    fields: [
      { key: "tmt_kgb", label: "TMT KGB", type: "date" },
      {
        key: "pangkat_kgb",
        label: "Pangkat/Gol",
        type: "select",
        options: [
          { value: "I/a", label: "I/a" },
          { value: "I/b", label: "I/b" },
          { value: "I/c", label: "I/c" },
          { value: "I/d", label: "I/d" },
          { value: "II/a", label: "II/a" },
          { value: "II/b", label: "II/b" },
          { value: "II/c", label: "II/c" },
          { value: "II/d", label: "II/d" },
          { value: "III/a", label: "III/a" },
          { value: "III/b", label: "III/b" },
          { value: "III/c", label: "III/c" },
          { value: "III/d", label: "III/d" },
          { value: "IV/a", label: "IV/a" },
          { value: "IV/b", label: "IV/b" },
          { value: "IV/c", label: "IV/c" },
          { value: "IV/d", label: "IV/d" },
          { value: "IV/e", label: "IV/e" },
        ],
      },
      { key: "gaji_kgb", label: "Gaji" },
      { key: "sk_kgb", label: "SK KGB" },
      { key: "tglsk_kgb", label: "Tanggal SK KGB", type: "date" },
    ],
  },
  {
    title: "Dokumen & Keluarga",
    description: "Dokumen kepegawaian dan data keluarga.",
    icon: UserGroupIcon,
    fields: [
      { key: "karpeg", label: "Karpeg" },
      { key: "karis", label: "Karis/Karsu" },
      { key: "nama_ibukdg", label: "Nama Ibu Kandung" },
      { key: "nm_suami_istri", label: "Nama Suami/Istri" },
      {
        key: "pek_suami_istri",
        label: "Pekerjaan Suami/Istri",
        type: "select",
        options: [
          { value: "Tidak bekerja", label: "Tidak bekerja" },
          { value: "Nelayan", label: "Nelayan" },
          { value: "Petani", label: "Petani" },
          { value: "Peternak", label: "Peternak" },
          { value: "PNS/TNI/Polri", label: "PNS/TNI/Polri" },
          { value: "Karyawan Swasta", label: "Karyawan Swasta" },
          { value: "Pedagang Kecil", label: "Pedagang Kecil" },
          { value: "Pedagang Besar", label: "Pedagang Besar" },
          { value: "Wiraswasta", label: "Wiraswasta" },
          { value: "Wirausaha", label: "Wirausaha" },
          { value: "Buruh", label: "Buruh" },
          { value: "Pensiunan", label: "Pensiunan" },
          { value: "Tenaga Kerja Indonesia", label: "Tenaga Kerja Indonesia" },
          { value: "Karyawan BUMN", label: "Karyawan BUMN" },
          { value: "Tidak dapat diterapkan", label: "Tidak dapat diterapkan" },
          { value: "Sudah Meninggal", label: "Sudah Meninggal" },
          { value: "Lainnya", label: "Lainnya" },
        ],
      },
      { key: "npwp", label: "NPWP" },
    ],
  },
  {
    title: "Status Kepegawaian",
    description: "Isi hanya jika pegawai sudah tidak aktif.",
    icon: ArrowRightOnRectangleIcon,
    hideOnCreate: true,
    fields: [
      {
        key: "status_keluar",
        label: "Status Keluar",
        type: "select",
        options: [{ value: "Keluar", label: "Keluar" }],
      },
      { key: "nomor_surat", label: "Nomor Surat" },
      { key: "tanggal_keluar", label: "Tanggal Keluar", type: "date" },
      {
        key: "keterangan_keluar",
        label: "Keterangan",
        type: "textarea",
        full: true,
      },
    ],
  },
];

// =====================================
// Load Data (mode edit)
// =====================================

async function loadData() {
  if (!isEdit) return;

  loading.value = true;

  try {
    const data = await pegawaiStore.fetchPegawaiByNuptk(route.params.nuptk);

    Object.assign(form, data, { foto_url: data.foto_url || "" });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal memuat data pegawai.");
  } finally {
    loading.value = false;
  }
}

// =====================================
// Foto
// =====================================

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

function showImage() {
  visibleLightbox.value = true;
}

async function handleDeleteFoto() {
  deletingFoto.value = true;
  try {
    const res = await pegawaiStore.deleteFotoPegawai(route.params.nuptk);
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

// =====================================
// Validasi Digit & Cek Duplikat
// =====================================

function validateDigitLocal(key, value) {
  const length = digitRules[key];

  if (!value) {
    fieldErrors[key] = "";
    return true;
  }

  if (!/^\d+$/.test(value) || value.length !== length) {
    fieldErrors[key] = `Harus ${length} digit angka.`;
    return false;
  }

  fieldErrors[key] = "";
  return true;
}

async function checkDuplicateField(key, value) {
  if (key === "nuptk" && isEdit) return; // NUPTK tidak bisa diubah saat edit

  try {
    const exclude = isEdit ? route.params.nuptk : "";
    const result = await pegawaiStore.checkDuplicate(key, value, exclude);

    if (result.exists) {
      duplicateWarnings[key] = result;
    }
  } catch {
    // Pengecekan duplikat bersifat bantuan (non-blocking), diamkan kalau gagal
  }
}

function watchWithDuplicateCheck(key) {
  watch(
    () => form[key],
    (value) => {
      clearTimeout(debounceTimers[key]);
      duplicateWarnings[key] = null;

      if (!validateDigitLocal(key, value)) return;

      debounceTimers[key] = setTimeout(() => {
        checkDuplicateField(key, value);
      }, 300);
    },
  );
}

watchWithDuplicateCheck("nuptk");
watchWithDuplicateCheck("nip_baru");
watchWithDuplicateCheck("nik");

// No. KK cuma perlu validasi panjang, tidak perlu cek duplikat
watch(
  () => form.no_kk,
  (value) => validateDigitLocal("no_kk", value),
);

// =====================================
// Submit
// =====================================

function buildPayload() {
  const { foto_url, foto, ...rest } = form;
  return rest;
}

async function submit() {
  if (loading.value) return;

  if (!isEdit && !form.nuptk.trim()) {
    toast.error("NUPTK wajib diisi.");
    return;
  }

  if (!form.nama_pegawai.trim()) {
    toast.error("Nama pegawai wajib diisi.");
    return;
  }

  if (!form.jabatan.trim()) {
    toast.error("Jabatan wajib diisi.");
    return;
  }

  if (
    !validateDigitLocal("nuptk", form.nuptk) ||
    !validateDigitLocal("nip_baru", form.nip_baru) ||
    !validateDigitLocal("nik", form.nik) ||
    !validateDigitLocal("no_kk", form.no_kk)
  ) {
    toast.error(
      "Periksa kembali NUPTK/NIP/NIK/No. KK — harus sesuai jumlah digit.",
    );
    return;
  }

  if (
    duplicateWarnings.nuptk ||
    duplicateWarnings.nip_baru ||
    duplicateWarnings.nik
  ) {
    toast.error(
      "Ada data yang ganda dengan pegawai lain, silakan periksa kembali.",
    );
    return;
  }

  loading.value = true;

  try {
    const payload = buildPayload();
    let nuptk = route.params.nuptk;

    let dataChanged = true; // default true untuk mode tambah

    if (isEdit) {
      const response = await pegawaiStore.updatePegawai(nuptk, payload);
      dataChanged = response.changed !== false;
    } else {
      await pegawaiStore.createPegawai(payload);
      nuptk = form.nuptk;
    }

    if (selectedFoto.value) {
      const fd = new FormData();
      fd.append("foto", selectedFoto.value);
      await pegawaiStore.updateFotoPegawai(nuptk, fd);
    }

    // Tentukan pesan akhir berdasarkan kombinasi: data teks berubah? foto diganti?
    if (!isEdit) {
      toast.success("Pegawai berhasil ditambahkan.");
    } else if (dataChanged && selectedFoto.value) {
      toast.success("Data dan foto pegawai berhasil diperbarui.");
    } else if (dataChanged) {
      toast.success("Data pegawai berhasil diperbarui.");
    } else if (selectedFoto.value) {
      toast.success("Foto pegawai berhasil diperbarui.");
    } else {
      toast.warning("Tidak ada perubahan data.");
    }

    router.push({ name: "dashboard-pegawai" });
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menyimpan data pegawai.",
    );
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  await loadData();
});
</script>
