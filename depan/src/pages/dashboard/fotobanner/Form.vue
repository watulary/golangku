<template>
  <DashboardPageHeader />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Banner"
        description="Nama, kegiatan, dan gambar yang tampil pada halaman beranda."
      >
        <DocumentTextIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="space-y-6 p-6">
        <!-- Gambar -->
        <div
          class="flex min-h-56 items-center justify-center rounded-2xl border border-dashed border-gray-300 bg-gray-50 p-5 dark:border-slate-700 dark:bg-[#0B0F1A]"
        >
          <div class="group relative">
            <img
              v-if="form.gambar_url"
              :src="form.gambar_url"
              alt="Banner"
              class="h-32 w-48 cursor-pointer rounded-2xl object-cover transition duration-200 group-hover:opacity-80"
              @click="openLightbox"
            />
            <PhotoIcon
              v-else
              class="h-20 w-20 text-gray-300 transition duration-200 group-hover:opacity-80 dark:text-slate-600"
            />
            <!-- Overlay Zoom -->
            <div
              v-if="form.gambar_url"
              class="pointer-events-none absolute inset-0 flex items-center justify-center rounded-2xl bg-black/0 opacity-0 transition duration-200 group-hover:bg-black/30 group-hover:opacity-100"
            >
              <MagnifyingGlassPlusIcon class="h-5 w-5 text-white" />
            </div>
            <!-- Overlay Kompresi -->
            <div
              v-if="processingGambar"
              class="pointer-events-none absolute inset-0 flex items-center justify-center rounded-2xl bg-black/50"
            >
              <ArrowPathIcon class="h-6 w-6 animate-spin text-white" />
            </div>
            <button
              type="button"
              title="Ganti Gambar"
              :disabled="processingGambar"
              @click="gambarInput.click()"
              class="absolute -right-4 -bottom-4 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700 dark:hover:text-white"
            >
              <PencilSquareIcon class="h-4 w-4" />
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

        <!-- Nama -->
        <FormField label="Judul Banner">
          <input
            v-model="form.nama"
            type="text"
            :class="inputClass"
            placeholder="Ketik disini..."
          />
        </FormField>

        <!-- Kegiatan -->
        <FormField label="Kegiatan">
          <input
            v-model="form.kegiatan"
            type="text"
            :class="inputClass"
            placeholder="Ketik disini..."
          />
        </FormField>
      </div>
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
            :to="{ name: 'dashboard-fotobanner' }"
            class="inline-flex items-center justify-center rounded-xl border border-gray-300 px-5 py-2.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
          >
            Batal
          </RouterLink>
        </div>
      </div>
    </section>
  </form>
  <!-- Lightbox Gambar -->
  <VueEasyLightbox
    :visible="lightboxVisible"
    :imgs="lightboxImg"
    @hide="lightboxVisible = false"
  />
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import FormField from "@/components/form/FormField.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";
import { useToast } from "@/composables/useToast";
import { compressImage } from "@/utils/image";
import { useFotoBannerStore } from "@/stores/fotobanner";
import VueEasyLightbox from "vue-easy-lightbox";

import {
  ArrowPathIcon,
  DocumentTextIcon,
  MagnifyingGlassPlusIcon,
  PencilSquareIcon,
  PhotoIcon,
} from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const fotoBannerStore = useFotoBannerStore();
const toast = useToast();

const isEdit = !!route.params.id;

const gambarInput = ref(null);
const selectedGambar = ref(null);
const processingGambar = ref(false);
const loading = ref(false);
const lightboxVisible = ref(false);
const lightboxImg = ref("");

const form = reactive({
  id: null,
  nama: "",
  kegiatan: "",
  gambar_url: "",
  nama_pengguna: "",
  waktu: "",
});

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";

const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";

// =====================================
// Load Data (mode edit)
// =====================================

async function loadData() {
  if (!isEdit) return;

  loading.value = true;

  try {
    const data = await fotoBannerStore.fetchBannerById(route.params.id);

    Object.assign(form, {
      id: data.id,
      nama: data.nama || "",
      kegiatan: data.kegiatan || "",
      gambar_url: data.gambar_url || "",
      nama_pengguna: data.nama_pengguna || "",
      waktu: data.waktu || "",
    });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal memuat data banner.");
  } finally {
    loading.value = false;
  }
}

// =====================================
// Submit
// =====================================

async function submit() {
  if (loading.value) return;

  if (!form.nama.trim()) {
    toast.error("Nama wajib diisi.");
    return;
  }

  if (!form.kegiatan.trim()) {
    toast.error("Kegiatan wajib diisi.");
    return;
  }

  if (!isEdit && !selectedGambar.value) {
    toast.error("Gambar wajib diunggah.");
    return;
  }

  loading.value = true;

  try {
    const data = new FormData();

    data.append("nama", form.nama);
    data.append("kegiatan", form.kegiatan);

    if (selectedGambar.value) {
      data.append("gambar", selectedGambar.value);
    }

    if (isEdit) {
      const response = await fotoBannerStore.updateBanner(form.id, data);

      if (response.changed) {
        toast.success(response.message || "Banner berhasil diperbarui.");
      } else {
        toast.warning(response.message || "Tidak ada perubahan data.");
      }
    } else {
      await fotoBannerStore.createBanner(data);
      toast.success("Banner berhasil ditambahkan.");
    }

    router.push({ name: "dashboard-fotobanner" });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menyimpan banner.");
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

  processingGambar.value = true;

  try {
    const compressed = await compressImage(file, 100);

    selectedGambar.value = compressed;
    form.gambar_url = URL.createObjectURL(compressed);
  } catch {
    toast.error("Gagal memproses gambar.");
  } finally {
    processingGambar.value = false;
  }
}

function openLightbox() {
  if (!form.gambar_url) return;
  lightboxImg.value = form.gambar_url;
  lightboxVisible.value = true;
}

onMounted(async () => {
  await loadData();
});
</script>
