<template>
  <DashboardPageHeader
    title="Struktur Komite"
    subtitle="Kelola gambar struktur komite yang tampil pada halaman website."
  />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <!-- Struktur Komite -->
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Struktur Komite"
        description="Judul dan gambar struktur komite yang tampil pada halaman website."
      >
        <PhotoIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="space-y-6 p-6">
        <!-- Gambar -->
        <div
          class="flex min-h-56 items-center justify-center rounded-2xl border border-dashed border-gray-300 bg-gray-50 p-5 dark:border-slate-700 dark:bg-[#0B0F1A]"
        >
          <div class="group relative">
            <img
              v-if="form.foto_url"
              :src="form.foto_url"
              alt="Struktur Komite"
              class="max-h-72 w-auto cursor-pointer rounded-2xl object-contain transition duration-200 group-hover:opacity-80"
              @click="openLightbox"
            />
            <PhotoIcon
              v-else
              class="h-20 w-20 text-gray-300 transition duration-200 group-hover:opacity-80 dark:text-slate-600"
            />
            <!-- Overlay Zoom -->
            <div
              v-if="form.foto_url"
              class="pointer-events-none absolute inset-0 flex items-center justify-center rounded-2xl bg-black/0 opacity-0 transition duration-200 group-hover:bg-black/30 group-hover:opacity-100"
            >
              <MagnifyingGlassPlusIcon class="h-5 w-5 text-white" />
            </div>
            <!-- Overlay Kompresi -->
            <div
              v-if="processingFoto"
              class="pointer-events-none absolute inset-0 flex items-center justify-center rounded-2xl bg-black/50"
            >
              <ArrowPathIcon class="h-6 w-6 animate-spin text-white" />
            </div>
            <button
              type="button"
              title="Ganti Gambar"
              :disabled="processingFoto"
              @click="gambarInput.click()"
              class="absolute -right-4 -bottom-4 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700 dark:hover:text-white"
            >
              <PencilSquareIcon class="h-4 w-4" />
            </button>
            <!-- Hapus (kembali ke gambar default; disembunyikan bila sudah default) -->
            <button
              v-if="form.foto_url && !isDefaultFoto"
              type="button"
              title="Hapus Gambar"
              :disabled="processingFoto"
              @click="removeFoto"
              class="absolute -left-4 -bottom-4 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 hover:bg-red-50 hover:text-red-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-red-500/10 dark:hover:text-red-400"
            >
              <XMarkIcon class="h-4 w-4" />
            </button>
            <input
              ref="gambarInput"
              type="file"
              accept=".jpg,.jpeg,.png,.webp"
              class="hidden"
              @change="handleGambarChange"
            />
          </div>
        </div>

        <!-- Judul -->
        <FormField label="Judul">
          <input
            v-model="form.judul"
            type="text"
            :class="inputClass"
            placeholder="Ketik disini..."
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
            Diperbarui oleh: {{ lastUpdatedBy }}
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
  <VueEasyLightbox
    :visible="lightboxVisible"
    :imgs="lightboxImg"
    @hide="lightboxVisible = false"
  />
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import FormField from "@/components/form/FormField.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";
import { useToast } from "@/composables/useToast";
import { useStrukturKomiteStore } from "@/stores/strukturkomite";
import VueEasyLightbox from "vue-easy-lightbox";
import { formatDateTime } from "@/utils/date";
import { compressImage } from "@/utils/image";
import {
  ArrowPathIcon,
  ClockIcon,
  MagnifyingGlassPlusIcon,
  PhotoIcon,
  UserCircleIcon,
  PencilSquareIcon,
  XMarkIcon,
} from "@heroicons/vue/24/outline";

const strukturKomiteStore = useStrukturKomiteStore();
const toast = useToast();

const gambarInput = ref(null);
const selectedGambar = ref(null);
const hapusGambar = ref(false);
const fotoBase = ref(""); // alamat folder gambar di server, diambil dari foto_url
const processingFoto = ref(false);
const loading = ref(false);
const lightboxVisible = ref(false);
const lightboxImg = ref("");

// Harus sama dengan defaultGambar di handler.go
const DEFAULT_GAMBAR = "1636812404_b39deec2f158cb789567.jpg";

const form = reactive({
  id: null,
  judul: "",
  foto_url: "",
});

const isDefaultFoto = computed(() => form.foto_url.endsWith(DEFAULT_GAMBAR));

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";

const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";

// =====================================
// Load Data
// =====================================

async function loadData() {
  loading.value = true;

  try {
    await strukturKomiteStore.getDashboardStrukturKomite();

    Object.assign(form, {
      id: strukturKomiteStore.dashboardStrukturKomite.id,
      judul: strukturKomiteStore.dashboardStrukturKomite.judul || "",
      foto_url: strukturKomiteStore.dashboardStrukturKomite.foto_url || "",
    });

    fotoBase.value = form.foto_url.replace(/[^/]+$/, "");
    selectedGambar.value = null;
    hapusGambar.value = false;
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat data struktur komite.",
    );
  } finally {
    loading.value = false;
  }
}

// =====================================
// Submit
// =====================================

async function submit() {
  if (loading.value) return;

  loading.value = true;

  try {
    const data = new FormData();

    data.append("judul", form.judul);

    if (selectedGambar.value) {
      data.append("gambar", selectedGambar.value);
    }

    if (hapusGambar.value) {
      data.append("hapus_foto", "1");
    }

    const response =
      await strukturKomiteStore.updateDashboardStrukturKomite(data);

    // Gambar sudah tersimpan/terhapus di server; jangan kirim ulang di simpan berikutnya.
    selectedGambar.value = null;
    hapusGambar.value = false;

    // Samakan tampilan dengan data server (mis. gambar default setelah dihapus).
    if (response.data?.foto_url) {
      form.foto_url = response.data.foto_url;
      fotoBase.value = form.foto_url.replace(/[^/]+$/, "");
    }

    if (response.changed) {
      toast.success(response.message || "Struktur komite berhasil diperbarui.");
    } else {
      toast.warning(response.message || "Tidak ada perubahan data.");
    }
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menyimpan struktur komite.",
    );
  } finally {
    loading.value = false;
  }
}

async function handleGambarChange(event) {
  const file = event.target.files?.[0];

  event.target.value = "";

  if (!file) return;

  const allowed = ["image/jpeg", "image/png", "image/webp"];

  if (!allowed.includes(file.type)) {
    toast.error("Gambar harus berformat JPG, JPEG, PNG atau WEBP.");
    return;
  }

  processingFoto.value = true;

  try {
    const compressed = await compressImage(file, 100);

    selectedGambar.value = compressed;
    hapusGambar.value = false;
    form.foto_url = URL.createObjectURL(compressed);
  } catch {
    toast.error("Gagal memproses gambar.");
  } finally {
    processingFoto.value = false;
  }
}

// =====================================
// Hapus Gambar (kembali ke default; diterapkan saat tombol Simpan ditekan)
// =====================================

function removeFoto() {
  selectedGambar.value = null;
  hapusGambar.value = true;
  form.foto_url = fotoBase.value ? fotoBase.value + DEFAULT_GAMBAR : "";
}

// =====================================
// Last Updated
// =====================================

const lastUpdatedBy = computed(() => {
  return (
    strukturKomiteStore.dashboardStrukturKomite?.nama_pengguna ||
    "Belum diketahui"
  );
});

const lastUpdatedAt = computed(() => {
  if (!strukturKomiteStore.dashboardStrukturKomite?.waktu) {
    return "Belum pernah diperbarui";
  }

  return formatDateTime(strukturKomiteStore.dashboardStrukturKomite.waktu);
});

function openLightbox() {
  if (!form.foto_url) return;
  lightboxImg.value = form.foto_url;
  lightboxVisible.value = true;
}

onMounted(loadData);
</script>
