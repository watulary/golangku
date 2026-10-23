<template>
  <DashboardPageHeader />

  <div class="space-y-6">
    <div
      class="relative overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
    >
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <!-- Header -->
      <div class="border-b border-gray-100 px-6 py-5 dark:border-white/5">
        <div class="flex items-center gap-3">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-blue-100 text-blue-600 dark:bg-blue-500/10 dark:text-blue-300"
          >
            <PencilSquareIcon class="h-5 w-5" />
          </div>
          <div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Ubah Pengguna
            </h2>

            <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
              Perbarui data pengguna di bawah ini;
            </p>
          </div>
        </div>
      </div>

      <!-- Body -->
      <div class="space-y-6 p-6">
        <!-- Foto -->
        <div>
          <label
            class="mb-2 block text-sm font-medium text-gray-700 dark:text-slate-300"
          >
            Foto
          </label>

          <div class="group relative inline-block">
            <div
              class="h-32 w-24 overflow-hidden rounded-lg border border-gray-200 bg-gray-50 dark:border-slate-700 dark:bg-slate-800"
              :class="{ 'cursor-pointer': form.foto_url }"
              @click="form.foto_url && showImage()"
            >
              <img
                v-if="form.foto_url"
                :src="form.foto_url"
                alt="Foto Pengguna"
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

            <div
              v-if="processingFoto"
              class="pointer-events-none absolute inset-0 flex items-center justify-center rounded-lg bg-black/50"
            >
              <ArrowPathIcon class="h-6 w-6 animate-spin text-white" />
            </div>

            <button
              type="button"
              title="Ganti Foto"
              :disabled="processingFoto"
              @click="fotoInput.click()"
              class="absolute -right-3 -bottom-3 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:opacity-60 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700 dark:hover:text-white"
            >
              <PencilSquareIcon class="h-4 w-4" />
            </button>

            <button
              v-if="form.foto && form.foto !== FOTO_DEFAULT"
              type="button"
              title="Hapus Foto"
              :disabled="processingFoto"
              @click="showDeleteFotoDialog = true"
              class="absolute -left-3 -bottom-3 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 hover:bg-red-50 hover:text-red-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40 disabled:cursor-not-allowed disabled:opacity-60 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-red-500/10 dark:hover:text-red-400"
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

        <!-- Nama Pengguna -->
        <div>
          <label
            class="mb-2 block text-sm font-medium text-gray-700 dark:text-slate-300"
          >
            Nama
            <span class="text-gray-700 dark:text-slate-300">*</span>
          </label>

          <input
            v-model="form.nama_pengguna"
            type="text"
            autocomplete="off"
            class="w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-100"
            placeholder="Ketik disini..."
          />
        </div>

        <!-- Email -->
        <div>
          <label
            class="mb-2 block text-sm font-medium text-gray-700 dark:text-slate-300"
          >
            Email
            <span class="text-gray-700 dark:text-slate-300">*</span>
          </label>

          <input
            v-model="form.email"
            type="email"
            autocomplete="off"
            class="w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-100"
            placeholder="Ketik disini..."
          />
        </div>

        <!-- Level -->
        <div>
          <label
            class="mb-2 block text-sm font-medium text-gray-700 dark:text-slate-300"
          >
            Level
            <span class="text-gray-700 dark:text-slate-300">*</span>
          </label>

          <select
            v-model="form.level"
            class="w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-100"
          >
            <option value="">--Pilih--</option>
            <option value="Admin">Super Admin</option>
            <option value="Staf">Admin</option>
            <option value="Pegawai">Pegawai</option>
            <option value="Pesertadidik">Peserta Didik</option>
          </select>
        </div>

        <!-- NUPTK -->
        <div v-if="form.level !== 'Pesertadidik' && form.level !== ''">
          <label
            class="mb-2 block text-sm font-medium text-gray-700 dark:text-slate-300"
          >
            NUPTK
            <span class="text-gray-700 dark:text-slate-300">*</span>
          </label>

          <input
            v-model="form.nuptk"
            type="text"
            inputmode="numeric"
            autocomplete="off"
            class="w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-100"
            placeholder="Masukkan NUPTK"
          />

          <p class="mt-2 text-xs text-gray-500 dark:text-slate-400">
            NUPTK harus terdiri dari 16 digit angka.
          </p>

          <p
            v-if="nuptkDuplicate"
            class="mt-1 text-xs text-yellow-600 dark:text-yellow-400"
          >
            Data ganda a.n. {{ nuptkDuplicate.nama_pengguna }}
          </p>
        </div>

        <!-- NIS -->
        <div v-if="form.level === 'Pesertadidik'">
          <label
            class="mb-2 block text-sm font-medium text-gray-700 dark:text-slate-300"
          >
            NIS
            <span class="text-gray-700 dark:text-slate-300">*</span>
          </label>

          <input
            v-model="form.nis"
            type="text"
            inputmode="numeric"
            autocomplete="off"
            class="w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-100"
            placeholder="Masukkan NIS"
          />

          <p class="mt-2 text-xs text-gray-500 dark:text-slate-400">
            NIS hanya boleh berisi angka.
          </p>
        </div>
      </div>

      <!-- Footer -->
      <div
        class="flex flex-col-reverse gap-3 border-t border-gray-100 bg-white px-6 py-5 sm:flex-row dark:border-white/5 dark:bg-slate-900"
      >
        <button
          type="button"
          @click="submit"
          :disabled="loading"
          class="rounded-xl bg-blue-600 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:opacity-60"
        >
          <span v-if="loading">Menyimpan...</span>
          <span v-else>Simpan</span>
        </button>
        <button
          type="button"
          @click="kembali"
          class="rounded-xl border border-gray-300 bg-white px-5 py-2.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200 dark:hover:bg-slate-800"
        >
          Batal
        </button>
      </div>
    </div>
  </div>

  <VueEasyLightbox
    :visible="visibleLightbox"
    :imgs="lightboxImages"
    :index="0"
    @hide="visibleLightbox = false"
  />

  <DashboardConfirmDialog
    v-model="showDeleteFotoDialog"
    title="Hapus Foto"
    message="Foto pengguna akan dihapus. Lanjutkan?"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="deletingFoto"
    @confirm="handleDeleteFoto"
  />
</template>

<script setup>
import { computed, reactive, ref, watch, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";

import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";

import { usePenggunaStore } from "@/stores/pengguna";
import { useToast } from "@/composables/useToast";
import { compressImage } from "@/utils/image";
import VueEasyLightbox from "vue-easy-lightbox";

import {
  PencilSquareIcon,
  XMarkIcon,
  UserIcon,
  MagnifyingGlassPlusIcon,
  ArrowPathIcon,
} from "@heroicons/vue/24/outline";

const router = useRouter();
const route = useRoute();

const toast = useToast();

const penggunaStore = usePenggunaStore();

// =========================================
// State
// =========================================

const loading = ref(false);

const form = reactive({
  nama_pengguna: "",
  email: "",
  level: "",

  nuptk: "",
  nis: "",

  foto: "",
  foto_url: "",
});

// =========================================
// Foto
// =========================================

const FOTO_DEFAULT = "1636812404_b39deec2f158cb789567.jpg";

const fotoInput = ref(null);
const selectedFoto = ref(null);
const processingFoto = ref(false);

const showDeleteFotoDialog = ref(false);
const deletingFoto = ref(false);

const visibleLightbox = ref(false);
const lightboxImages = computed(() => (form.foto_url ? [form.foto_url] : []));

function showImage() {
  visibleLightbox.value = true;
}

async function handleFotoChange(event) {
  const file = event.target.files?.[0];

  event.target.value = "";

  if (!file) {
    return;
  }

  const allowed = ["image/jpeg", "image/png", "image/webp"];

  if (!allowed.includes(file.type)) {
    toast.error("Foto harus berformat JPG, JPEG, PNG atau WEBP.");
    return;
  }

  processingFoto.value = true;

  try {
    const compressed = await compressImage(file, 100);

    selectedFoto.value = compressed;
    form.foto_url = URL.createObjectURL(compressed);
  } catch {
    toast.error("Gagal memproses foto.");
  } finally {
    processingFoto.value = false;
  }
}

async function handleDeleteFoto() {
  deletingFoto.value = true;

  try {
    const res = await penggunaStore.deleteFotoPengguna(route.params.id);

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

// =========================================
// Ambil Data
// =========================================

onMounted(async () => {
  loading.value = true;

  try {
    const data = await penggunaStore.getPengguna(route.params.id);

    form.nama_pengguna = data.nama_pengguna;
    form.email = data.email;
    form.level = data.level;

    form.nuptk = data.nuptk || "";
    form.nis = data.nis || "";

    form.foto = data.foto || "";
    form.foto_url = data.foto_url || "";
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Data pengguna tidak ditemukan.",
    );

    router.push("/dashboard/pengguna");
  } finally {
    loading.value = false;
  }
});

// =========================================
// Cek Duplikat NUPTK (realtime)
// =========================================

const nuptkDuplicate = ref(null);

let nuptkTimeout = null;

// =========================================
// Watch Level
// =========================================

watch(
  () => form.level,
  (level) => {
    switch (level) {
      case "Pegawai":
        form.nis = "";
        break;

      case "Pesertadidik":
        form.nuptk = "";
        break;

      default:
        form.nis = "";
        break;
    }

    nuptkDuplicate.value = null;
  },
);

watch(
  () => form.nuptk,
  (value) => {
    clearTimeout(nuptkTimeout);
    nuptkDuplicate.value = null;

    if (!value) return;

    nuptkTimeout = setTimeout(async () => {
      try {
        const result = await penggunaStore.checkDuplicate(
          "nuptk",
          value,
          route.params.id,
        );

        if (result.exists) nuptkDuplicate.value = result;
      } catch {
        // non-blocking, diamkan kalau gagal
      }
    }, 300);
  },
);

// =========================================
// Kembali
// =========================================

function kembali() {
  router.push("/dashboard/pengguna");
}

// =========================================
// Submit
// =========================================

async function submit() {
  if (!form.nama_pengguna.trim()) {
    toast.error("Nama Pengguna wajib diisi.");
    return;
  }

  if (!form.email.trim()) {
    toast.error("Email wajib diisi.");
    return;
  }

  if (!form.level) {
    toast.error("Level wajib dipilih.");
    return;
  }

  if (form.level === "Pegawai" && !form.nuptk.trim()) {
    toast.error("NUPTK wajib diisi.");
    return;
  }

  if (form.level === "Pesertadidik" && !form.nis.trim()) {
    toast.error("NIS wajib diisi.");
    return;
  }

  if (nuptkDuplicate.value) {
    toast.error(
      "NUPTK sudah digunakan oleh pengguna lain, silakan periksa kembali.",
    );
    return;
  }

  loading.value = true;

  try {
    const response = await penggunaStore.updatePengguna(route.params.id, {
      nama_pengguna: form.nama_pengguna.trim(),
      email: form.email.trim(),
      level: form.level,

      // Admin, Staf, Pegawai boleh memiliki NUPTK
      nuptk: form.level !== "Pesertadidik" ? form.nuptk.trim() : "",

      // Hanya Peserta Didik yang memakai NIS
      nis: form.level === "Pesertadidik" ? form.nis.trim() : "",
    });

    const dataChanged = response.changed !== false;

    if (selectedFoto.value) {
      const fd = new FormData();
      fd.append("foto", selectedFoto.value);

      await penggunaStore.updateFotoPengguna(route.params.id, fd);
    }

    if (!dataChanged && !selectedFoto.value) {
      toast.warning(response.message);
      return;
    }

    if (dataChanged && selectedFoto.value) {
      toast.success("Data dan foto pengguna berhasil diperbarui.");
    } else if (dataChanged) {
      toast.success(response.message);
    } else {
      toast.success("Foto pengguna berhasil diperbarui.");
    }

    router.push("/dashboard/pengguna");
  } catch (error) {
    const message =
      error.response?.data?.message || "Gagal memperbarui pengguna.";

    if (message === "Tidak ada perubahan data.") {
      toast.warning(message);
    } else {
      toast.error(message);
    }
  } finally {
    loading.value = false;
  }
}
</script>
