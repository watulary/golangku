<template>
  <div class="bg-white rounded-2xl shadow-lg p-6 md:p-8">
    <div class="xl:col-span-3">
      <!-- Judul -->
      <h1 class="text-xl font-bold text-[#041983] mb-3">
        {{ strukturAkademikStore.strukturakademik?.judul }}
      </h1>

      <!-- Meta -->
      <div class="flex flex-wrap items-center gap-4 mt-4 text-xs text-gray-500">
        <span class="flex items-center gap-1.5">
          <CalendarDaysIcon class="w-3.5 h-3.5 shrink-0" />

          {{
            new Date(
              strukturAkademikStore.strukturakademik?.waktu,
            ).toLocaleDateString("id-ID", {
              day: "numeric",
              month: "long",
              year: "numeric",
            })
          }}
        </span>

        <span class="flex items-center gap-1.5">
          <UserIcon class="w-3.5 h-3.5 shrink-0" />
          {{ strukturAkademikStore.strukturakademik?.nama_pengguna }}
        </span>

        <span class="flex items-center gap-1.5">
          <EyeIcon class="w-3.5 h-3.5 shrink-0" />
          {{ formatHits(strukturAkademikStore.strukturakademik?.hits || 0) }}
        </span>
      </div>

      <!-- Gambar Struktur Akademik -->
      <div class="mt-8">
        <div
          class="relative group overflow-hidden rounded-xl shadow-xs border border-gray-100"
        >
          <img
            loading="lazy"
            v-if="strukturAkademikStore.strukturakademik?.foto_url"
            :src="strukturAkademikStore.strukturakademik.foto_url"
            :alt="strukturAkademikStore.strukturakademik.judul"
            class="w-full"
          />

          <div
            class="absolute inset-0 bg-black/30 opacity-0 group-hover:opacity-100 transition duration-300 flex items-center justify-center cursor-pointer"
            @click="showImage"
          >
            <span
              class="flex items-center gap-2 text-white text-sm font-medium px-4 py-2 rounded-full border border-white backdrop-blur-sm"
            >
              <MagnifyingGlassPlusIcon class="w-5 h-5" />
              Perbesar
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
  <VueEasyLightbox
    :visible="visible"
    :imgs="[strukturAkademikStore.strukturakademik?.foto_url]"
    :index="0"
    @hide="visible = false"
  />
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useStrukturAkademikStore } from "../../stores/strukturakademik";
import { formatHits } from "../../utils/format";
import VueEasyLightbox from "vue-easy-lightbox";
import {
  UserIcon,
  CalendarDaysIcon,
  EyeIcon,
  MagnifyingGlassPlusIcon,
} from "@heroicons/vue/24/outline";

const strukturAkademikStore = useStrukturAkademikStore();

const visible = ref(false);

const showImage = () => {
  visible.value = true;
};

onMounted(async () => {
  const viewed = sessionStorage.getItem("strukturakademik_viewed");

  if (!viewed) {
    await strukturAkademikStore.incrementHits();
    sessionStorage.setItem("strukturakademik_viewed", "1");
  }

  await Promise.all([strukturAkademikStore.fetchStrukturAkademik()]);
});
</script>
