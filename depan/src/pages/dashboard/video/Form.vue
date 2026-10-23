<template>
  <DashboardPageHeader />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Video"
        description="Nama, kegiatan, dan ID video YouTube yang tampil pada halaman video."
      >
        <DocumentTextIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="space-y-6 p-6">
        <!-- Preview Thumbnail -->
        <div
          class="flex min-h-56 items-center justify-center rounded-2xl border border-dashed border-gray-300 bg-gray-50 p-5 dark:border-slate-700 dark:bg-[#0B0F1A]"
        >
          <button
            type="button"
            :disabled="!thumbnailPreview"
            title="Preview Video"
            @click="openPreview"
            class="group relative block rounded-2xl focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40"
            :class="thumbnailPreview ? 'cursor-pointer' : 'cursor-default'"
          >
            <img
              v-if="thumbnailPreview"
              :src="thumbnailPreview"
              alt="Thumbnail Video"
              class="h-32 w-56 rounded-2xl object-cover"
              @error="thumbnailError = true"
            />

            <VideoCameraIcon
              v-else
              class="h-20 w-20 text-gray-300 dark:text-slate-600"
            />

            <div
              v-if="thumbnailPreview"
              class="absolute inset-0 flex items-center justify-center rounded-2xl bg-black/20 opacity-80 transition group-hover:bg-black/40 group-hover:opacity-100"
            >
              <PlayCircleIcon class="h-10 w-10 text-white drop-shadow" />
            </div>
          </button>
        </div>

        <!-- Nama -->
        <FormField label="Judul Video">
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

        <!-- ID Video YouTube -->
        <FormField label="ID Video YouTube">
          <input
            v-model="form.video"
            type="text"
            :class="inputClass"
            placeholder="JABDMzV5Yxk"
            maxlength="11"
            @input="thumbnailError = false"
          />

          <p class="mt-1.5 text-xs text-gray-500 dark:text-slate-400">
            Contoh:
            <code class="rounded bg-gray-100 px-1 py-0.5 dark:bg-slate-800">
              youtube.com/watch?v=<strong>JABDMzV5Yxk</strong>
            </code>
          </p>
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
            :to="{ name: 'dashboard-video' }"
            class="inline-flex items-center justify-center rounded-xl border border-gray-300 px-5 py-2.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
          >
            Batal
          </RouterLink>
        </div>
      </div>
    </section>
  </form>

  <!-- Preview Modal -->
  <Teleport to="body">
    <div
      v-if="showPreview"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4"
      @click.self="closePreview"
    >
      <div class="w-full max-w-2xl">
        <div class="mb-2 flex items-center justify-between">
          <h3 class="text-sm font-medium text-white">
            {{ form.nama || "Preview Video" }}
          </h3>

          <button
            type="button"
            title="Tutup"
            @click="closePreview"
            class="rounded-lg p-1.5 text-white/80 transition hover:bg-white/10 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40"
          >
            <XMarkIcon class="h-5 w-5" />
          </button>
        </div>

        <div class="aspect-video w-full overflow-hidden rounded-xl bg-black">
          <iframe
            v-if="previewURL"
            :src="previewURL"
            class="h-full w-full"
            allow="
              accelerometer;
              autoplay;
              clipboard-write;
              encrypted-media;
              gyroscope;
              picture-in-picture;
            "
            allowfullscreen
          />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";

import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import FormField from "@/components/form/FormField.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";

import { useToast } from "@/composables/useToast";
import { useVideoStore } from "@/stores/video";

import {
  ArrowPathIcon,
  DocumentTextIcon,
  PlayCircleIcon,
  VideoCameraIcon,
  XMarkIcon,
} from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const videoStore = useVideoStore();
const toast = useToast();

const isEdit = !!route.params.id;

const loading = ref(false);
const thumbnailError = ref(false);

const showPreview = ref(false);

const youtubeIdRegex = /^[a-zA-Z0-9_-]{11}$/;

const form = reactive({
  id: null,
  nama: "",
  kegiatan: "",
  video: "",
  nama_pengguna: "",
  waktu: "",
});

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";

const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";

const thumbnailPreview = computed(() => {
  const id = form.video.trim();

  if (!youtubeIdRegex.test(id) || thumbnailError.value) return "";

  return `https://img.youtube.com/vi/${id}/mqdefault.jpg`;
});

// ID YouTube yang valid dari form, dipakai untuk embed player
const validVideoId = computed(() => {
  const id = form.video.trim();
  return youtubeIdRegex.test(id) ? id : "";
});

// URL embed player untuk modal preview
const previewURL = computed(() => {
  if (!validVideoId.value) return "";
  return `https://www.youtube.com/embed/${validVideoId.value}?autoplay=1`;
});

function openPreview() {
  if (!validVideoId.value) return;
  showPreview.value = true;
}

function closePreview() {
  showPreview.value = false;
}

// =====================================
// Load Data (mode edit)
// =====================================

async function loadData() {
  if (!isEdit) return;

  loading.value = true;

  try {
    const data = await videoStore.fetchVideoById(route.params.id);

    Object.assign(form, {
      id: data.id,
      nama: data.nama || "",
      kegiatan: data.kegiatan || "",
      video: data.video || "",
      nama_pengguna: data.nama_pengguna || "",
      waktu: data.waktu || "",
    });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal memuat data video.");
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

  const videoId = form.video.trim();

  if (!videoId) {
    toast.error("ID video YouTube wajib diisi.");
    return;
  }

  if (!youtubeIdRegex.test(videoId)) {
    toast.error("ID video YouTube tidak valid (harus 11 karakter).");
    return;
  }

  loading.value = true;

  try {
    const payload = {
      nama: form.nama,
      kegiatan: form.kegiatan,
      video: videoId,
    };

    if (isEdit) {
      const response = await videoStore.updateVideo(form.id, payload);

      if (response.changed) {
        toast.success(response.message || "Video berhasil diperbarui.");
      } else {
        toast.warning(response.message || "Tidak ada perubahan data.");
      }
    } else {
      await videoStore.createVideo(payload);
      toast.success("Video berhasil ditambahkan.");
    }

    router.push({ name: "dashboard-video" });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menyimpan video.");
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  await loadData();
});
</script>
