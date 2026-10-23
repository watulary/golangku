<template>
  <router-view />
  <ToastContainer />
</template>
<script setup>
import { watch, onMounted } from "vue";
import ToastContainer from "@/components/ui/ToastContainer.vue";
import { useSettingStore } from "@/stores/setting";
import api from "@/services/api";

const settingStore = useSettingStore();

onMounted(() => {
  settingStore.fetchSetting();
});

watch(
  () => settingStore.setting?.nama_sekolah,
  (nama) => {
    document.title = nama || "Sekolah";
  },
  { immediate: true },
);

// Favicon dinamis dari setting.logo_sekolah. Gambar persegi (proporsi logo
// terjaga) dibuat oleh backend di GET /public/favicon.
function setFavicon(logoSekolah) {
  let link = document.querySelector("link[rel~='icon']");
  if (!link) {
    link = document.createElement("link");
    link.rel = "icon";
    document.head.appendChild(link);
  }

  // Logo kosong/dihapus: kosongkan favicon (tanpa request ke server).
  // Ganti dengan "/favicon.png" kalau kamu punya favicon bawaan.
  if (!logoSekolah) {
    link.removeAttribute("type");
    link.href = "data:,";
    return;
  }

  const base = (api.defaults.baseURL || "/api").replace(/\/$/, "");

  link.type = "image/png";
  // ?v= = nama file logo (berubah tiap upload) supaya cache ikut ter-refresh.
  link.href = `${base}/public/favicon?v=${encodeURIComponent(logoSekolah)}`;
}

// Tunggu setting selesai dimuat (null = belum dimuat), supaya favicon tidak
// sempat dikosongkan sebelum datanya ada.
watch(
  () => settingStore.setting,
  (setting) => {
    if (setting) setFavicon(setting.logo_sekolah);
  },
  { immediate: true },
);
</script>
