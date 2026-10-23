<template>
  <DashboardPageHeader />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Pengumuman"
        description="Judul dan isi pengumuman yang tampil pada halaman website."
      >
        <DocumentTextIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="space-y-6 p-6">
        <!-- Foto -->
        <div
          class="flex min-h-56 items-center justify-center rounded-2xl border border-dashed border-gray-300 bg-gray-50 p-5 dark:border-slate-700 dark:bg-[#0B0F1A]"
        >
          <div class="group relative">
            <img
              v-if="form.foto_url"
              :src="form.foto_url"
              alt="Foto Pengumuman"
              class="h-32 w-48 cursor-pointer rounded-2xl object-cover transition duration-200 group-hover:opacity-80"
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
              title="Ganti Foto"
              :disabled="processingFoto"
              @click="fotoInput.click()"
              class="absolute -right-4 -bottom-4 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700 dark:hover:text-white"
            >
              <PencilSquareIcon class="h-4 w-4" />
            </button>
            <button
              v-if="form.foto_url && !isFotoDefault"
              type="button"
              title="Hapus Foto"
              :disabled="processingFoto"
              @click="showDeleteFotoDialog = true"
              class="absolute -left-4 -bottom-4 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 hover:bg-red-50 hover:text-red-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-red-500/10 dark:hover:text-red-400"
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

        <!-- Judul -->
        <FormField label="Judul Pengumuman">
          <input
            v-model="form.judul_pengumuman"
            type="text"
            :class="inputClass"
            placeholder="Ketik disini..."
          />
        </FormField>

        <!-- Status -->
        <FormField label="Status">
          <select v-model="form.status_pengumuman" :class="inputClass">
            <option value="Publish">Publish</option>
            <option value="Draft">Draft</option>
          </select>
        </FormField>

        <!-- Ringkasan -->
        <FormField label="Ringkasan">
          <textarea
            v-model="form.ringkasan"
            rows="3"
            maxlength="100"
            :class="inputClass"
            placeholder="Ketik disini..."
          />

          <p
            class="mt-1.5 text-right text-xs text-gray-400 dark:text-slate-500"
          >
            {{ form.ringkasan.length }}/100
          </p>
        </FormField>

        <!-- Isi -->
        <FormField label="Isi Pengumuman">
          <TiptapEditor v-model="form.isi_pengumuman" />
        </FormField>
      </div>
      <!-- Aksi -->
      <div
        class="flex items-center gap-3 border-t border-gray-100 px-6 py-4 dark:border-white/5"
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
            :to="{ name: 'dashboard-pengumuman' }"
            class="inline-flex items-center justify-center rounded-xl border border-gray-300 px-5 py-2.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
          >
            Batal
          </RouterLink>
        </div>
      </div>
    </section>
  </form>
  <DashboardConfirmDialog
    v-model="showDeleteFotoDialog"
    title="Hapus Foto"
    message="Foto pengumuman akan dihapus. Lanjutkan?"
    confirm-text="Hapus"
    confirm-color="red"
    @confirm="handleHapusFoto"
  />
  <!-- Lightbox Foto -->
  <VueEasyLightbox
    :visible="lightboxVisible"
    :imgs="lightboxImg"
    @hide="lightboxVisible = false"
  />
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import FormField from "@/components/form/FormField.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";
import TiptapEditor from "@/components/tiptap/TiptapEditor.vue";
import { useToast } from "@/composables/useToast";
import { compressImage } from "@/utils/image";
import { usePengumumanStore } from "@/stores/pengumuman";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";
import VueEasyLightbox from "vue-easy-lightbox";

import {
  ArrowPathIcon,
  DocumentTextIcon,
  MagnifyingGlassPlusIcon,
  PencilSquareIcon,
  PhotoIcon,
  XMarkIcon,
} from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const pengumumanStore = usePengumumanStore();
const toast = useToast();

const isEdit = !!route.params.id;

const fotoInput = ref(null);
const selectedFoto = ref(null);
const processingFoto = ref(false);
const loading = ref(false);
const hapusFoto = ref(false);
const showDeleteFotoDialog = ref(false);
const lightboxVisible = ref(false);
const lightboxImg = ref("");

const form = reactive({
  id_pengumuman: null,
  judul_pengumuman: "",
  ringkasan: "",
  isi_pengumuman: "",
  status_pengumuman: "Publish",
  foto_url: "",
  nama_pengguna: "",
  tanggal_publish: "",
});

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";
const FOTO_DEFAULT = "1636812404_b39deec2f158cb789567.jpg";
const isFotoDefault = computed(() => {
  return form.foto_url.endsWith(FOTO_DEFAULT);
});
const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";

// =====================================
// Load Data (mode edit)
// =====================================

async function loadData() {
  if (!isEdit) return;

  loading.value = true;

  try {
    const data = await pengumumanStore.fetchPengumumanById(route.params.id);

    Object.assign(form, {
      id_pengumuman: data.id_pengumuman,
      judul_pengumuman: data.judul_pengumuman || "",
      ringkasan: data.ringkasan || "",
      isi_pengumuman: data.isi_pengumuman || "",
      status_pengumuman: data.status_pengumuman || "Publish",
      foto_url: data.foto_url || "",
      nama_pengguna: data.nama_pengguna || "",
      tanggal_publish: data.tanggal_publish || "",
    });
    hapusFoto.value = false;
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat data pengumuman.",
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

  if (!form.judul_pengumuman.trim()) {
    toast.error("Judul pengumuman wajib diisi.");
    return;
  }

  if (!form.isi_pengumuman.trim()) {
    toast.error("Isi pengumuman wajib diisi.");
    return;
  }

  loading.value = true;

  try {
    const data = new FormData();

    data.append("judul_pengumuman", form.judul_pengumuman);
    data.append("ringkasan", form.ringkasan);
    data.append("isi_pengumuman", form.isi_pengumuman);
    data.append("status_pengumuman", form.status_pengumuman);

    if (selectedFoto.value) {
      data.append("foto", selectedFoto.value);
    } else if (hapusFoto.value) {
      data.append("hapus_foto", "1");
    }

    if (isEdit) {
      const response = await pengumumanStore.updatePengumuman(
        form.id_pengumuman,
        data,
      );

      if (response.changed) {
        toast.success(response.message || "Pengumuman berhasil diperbarui.");
      } else {
        toast.warning(response.message || "Tidak ada perubahan data.");
      }
    } else {
      await pengumumanStore.createPengumuman(data);
      toast.success("Pengumuman berhasil dibuat.");
    }

    router.push({ name: "dashboard-pengumuman" });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menyimpan pengumuman.");
  } finally {
    loading.value = false;
  }
}

async function handleFotoChange(event) {
  const file = event.target.files?.[0];

  event.target.value = "";

  if (!file) return;

  const allowed = ["image/jpeg", "image/png", "image/webp"];

  if (!allowed.includes(file.type)) {
    toast.error("Foto harus berformat JPG, JPEG, PNG atau WEBP.");
    return;
  }

  processingFoto.value = true;

  try {
    const compressed = await compressImage(file, 100);

    selectedFoto.value = compressed;
    hapusFoto.value = false;
    form.foto_url = URL.createObjectURL(compressed);
  } catch {
    toast.error("Gagal memproses foto.");
  } finally {
    processingFoto.value = false;
  }
}

function handleHapusFoto() {
  selectedFoto.value = null;
  hapusFoto.value = true;
  form.foto_url = "";
  if (fotoInput.value) {
    fotoInput.value.value = "";
  }
  showDeleteFotoDialog.value = false;
}

function openLightbox() {
  if (!form.foto_url) return;
  lightboxImg.value = form.foto_url;
  lightboxVisible.value = true;
}

onMounted(async () => {
  await loadData();
});
</script>
