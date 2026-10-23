<template>
  <!-- Tombol -->
  <ToolbarButton title="Upload Gambar" :disabled="loading" @click="openPicker">
    <PhotoIcon v-if="!loading" class="h-5 w-5" />

    <ArrowPathIcon v-else class="h-5 w-5 animate-spin" />
  </ToolbarButton>

  <!-- File Picker -->
  <input
    ref="fileInput"
    type="file"
    accept=".jpg,.jpeg,.png,.webp"
    class="hidden"
    @change="handleFile"
  />
</template>

<script setup>
import { ref } from "vue";

import ToolbarButton from "./ToolbarButton.vue";

import { useToast } from "@/composables/useToast";
import { useTiptapStore } from "@/stores/tiptap";
import { compressImage } from "@/utils/image";

import { PhotoIcon, ArrowPathIcon } from "@heroicons/vue/24/outline";

const props = defineProps({
  editor: {
    type: Object,
    required: true,
  },
});

const toast = useToast();
const tiptapStore = useTiptapStore();

const fileInput = ref(null);
const loading = ref(false);

function openPicker() {
  fileInput.value?.click();
}

async function handleFile(event) {
  const file = event.target.files?.[0];

  if (!file) return;

  // ============================
  // Validasi MIME
  // ============================

  const allowed = ["image/jpeg", "image/png", "image/webp"];

  if (!allowed.includes(file.type)) {
    toast.error("Format gambar harus JPG, JPEG, PNG atau WEBP.");
    event.target.value = "";
    return;
  }

  loading.value = true;

  try {
    // ============================
    // Kompres Otomatis (maks 100 KB)
    // ============================

    const compressed = await compressImage(file, 100);

    const response = await tiptapStore.uploadImage(compressed);

    if (!props.editor) {
      toast.error("Editor belum siap.");
      return;
    }

    props.editor
      .chain()
      .focus()
      .setImage({
        src: response.url,
        alt: file.name,
      })
      .run();

    toast.success(response.message);
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal mengunggah gambar.");
  } finally {
    loading.value = false;

    if (fileInput.value) {
      fileInput.value.value = "";
    }
  }
}
</script>
