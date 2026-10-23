<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
    >
      <div
        class="relative w-full max-w-xl rounded-xl border border-gray-100 bg-white shadow-xl dark:border-white/5 dark:bg-slate-900"
      >
        <button
          type="button"
          title="Tutup"
          class="absolute right-4 top-4 flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition hover:bg-gray-100 hover:text-gray-900 dark:text-slate-400 dark:hover:bg-slate-800 dark:hover:text-white"
          @click="closePopup"
        >
          <XMarkIcon class="h-6 w-6" />
        </button>
        <div
          class="flex items-start gap-4 border-b border-gray-100 px-7 py-6 dark:border-white/5"
        >
          <span
            class="mt-0.5 rounded-lg bg-yellow-50 p-2 text-[#f2c808] dark:bg-blue-500/10 dark:text-blue-400"
          >
            <MegaphoneIcon class="h-6 w-6" />
          </span>
          <h3 class="pr-8 text-lg font-semibold text-[#041983] dark:text-white">
            {{ popupJudul || "Pengumuman" }}
          </h3>
        </div>
        <div class="space-y-5 p-7">
          <p
            v-if="popupIsi"
            class="whitespace-pre-line text-base leading-relaxed text-gray-600 dark:text-slate-300"
          >
            {{ popupIsi }}
          </p>
          <a
            v-if="popupLink"
            :href="popupLink"
            rel="noopener noreferrer"
            class="inline-flex items-center justify-center gap-1 rounded-xl bg-blue-600 px-5 py-3 text-sm font-medium text-white transition-all duration-300 hover:gap-2 hover:bg-blue-700"
            @click="closePopup"
          >
            Baca Selengkapnya
            <span>→</span>
          </a>
        </div>
      </div>
    </div>
  </Teleport>
</template>
<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useRoute } from "vue-router";
import { XMarkIcon, MegaphoneIcon } from "@heroicons/vue/24/outline";
import { useSettingStore } from "@/stores/setting";
// Pasang komponen ini SEKALI di root layout publik (mis. PublicLayout.vue),
// bukan di tiap halaman -- fetchSetting() di store sudah ada guard `loaded`
// jadi aman dipanggil berkali-kali tanpa fetch ulang ke server.
const settingStore = useSettingStore();
const route = useRoute();
const visible = ref(false);
// Sumber data: heroarea1-4 pada tabel setting (sama seperti yg dipakai
// section "Popup Pengumuman" di Website.vue dashboard).
const popupAktif = computed(() =>
  ["1", "true", "yes"].includes(
    String(settingStore.setting?.heroarea1 ?? "").toLowerCase(),
  ),
);
const popupJudul = computed(() => settingStore.setting?.heroarea2 || "");
const popupIsi = computed(() => settingStore.setting?.heroarea3 || "");
const popupLink = computed(() => settingStore.setting?.heroarea4 || "");

// Popup ini berfungsi sbg PENGINGAT khusus di halaman utama ("/") --
// BUKAN "tampil sekali lalu disimpan permanen". Jadi TIDAK ADA lagi
// status "sudah ditutup" yg disimpan (dulu pakai localStorage) --
// visibility sepenuhnya mengikuti path saat ini: tampil tiap kali
// path persis "/", dan otomatis sembunyi di path lain manapun
// (mis. /profil/sejarah).
function isEligible() {
  return !!(popupAktif.value && (popupJudul.value || popupIsi.value));
}

function closePopup() {
  visible.value = false;
}

// Route berubah lewat navigasi SPA (klik menu/router-link) tanpa reload
// halaman -- watcher ini yg menampilkan/menyembunyikan popup saat itu.
watch(
  () => route.path,
  (path) => {
    visible.value = path === "/" && isEligible();
  },
);

// Saat pertama kali app di-mount / halaman di-reload -- tunggu data
// setting selesai diambil dulu baru evaluasi apakah popup boleh tampil.
onMounted(async () => {
  await settingStore.fetchSetting();
  visible.value = route.path === "/" && isEligible();
});
</script>
