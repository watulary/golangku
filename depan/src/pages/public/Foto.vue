<template>
  <div class="bg-white rounded-2xl shadow-lg p-6 md:p-8">
    <!-- Header + Search -->
    <div
      class="flex flex-col md:flex-row md:items-end md:justify-between gap-4 mb-8"
    >
      <div class="self-start">
        <h2 class="text-2xl font-extrabold text-gray-600">Galeri Foto</h2>

        <p class="text-sm text-gray-500 mt-2">
          Total: {{ filteredFoto.length }}
        </p>

        <div class="mt-1 h-0.5 w-full flex overflow-hidden rounded-full">
          <div class="w-1/2 bg-[#041983]"></div>
          <div class="w-1/2 bg-[#f2c808]"></div>
        </div>
      </div>

      <div class="relative w-full md:w-64">
        <MagnifyingGlassIcon
          class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400"
        />

        <input
          v-model="keyword"
          type="text"
          placeholder="Pencarian..."
          class="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:border-[#f2c808]"
        />
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="filteredFoto.length === 0" class="text-center py-10">
      <p class="text-sm text-gray-500">Data tidak ditemukan.</p>
    </div>

    <!-- Grid Foto -->
    <div v-else class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-6">
      <div
        v-for="item in paginatedFoto"
        :key="item.id"
        @click="showImage(item)"
        class="group bg-white rounded-2xl overflow-hidden border border-gray-100 shadow-sm hover:shadow-lg transition-all duration-300 cursor-pointer"
      >
        <!-- Foto -->
        <div class="relative overflow-hidden">
          <img
            loading="lazy"
            :src="item.foto_url"
            :alt="item.nama"
            class="w-full aspect-square object-cover group-hover:scale-105 transition duration-300"
          />

          <div
            class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition duration-300 flex items-center justify-center"
          >
            <span
              class="flex items-center gap-2 text-white text-sm font-medium px-4 py-2 rounded-full border border-white backdrop-blur-sm"
            >
              <MagnifyingGlassPlusIcon class="w-5 h-5" />
              Perbesar
            </span>
          </div>
        </div>

        <!-- Content -->
        <div class="p-4">
          <p class="flex items-center gap-1.5 text-xs text-gray-500 mb-3">
            <CalendarDaysIcon class="w-3.5 h-3.5 shrink-0" />

            {{
              new Date(item.waktu).toLocaleDateString("id-ID", {
                day: "numeric",
                month: "long",
                year: "numeric",
              })
            }}
          </p>

          <h2 class="font-semibold text-[#041983] line-clamp-2">
            {{ item.nama }}
          </h2>

          <!-- <p class="text-sm text-gray-500 mt-1 line-clamp-1">
            {{ item.kegiatan }}
          </p> -->
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div
      v-if="totalPages > 1"
      class="flex justify-center items-center gap-2 mt-10"
    >
      <button
        @click="currentPage--"
        :disabled="currentPage === 1"
        class="p-2 rounded-lg border border-gray-200 text-[#041983] hover:bg-[#041983] hover:text-[#f2c808] disabled:opacity-50 disabled:hover:bg-white disabled:hover:text-[#041983]"
      >
        <ChevronLeftIcon class="w-5 h-5" />
      </button>

      <button
        v-for="page in totalPages"
        :key="page"
        @click="currentPage = page"
        :class="[
          'px-4 py-2 rounded-lg border',
          currentPage === page
            ? 'bg-[#041983] text-[#f2c808] border-[#041983]'
            : 'bg-white',
        ]"
      >
        {{ page }}
      </button>

      <button
        @click="currentPage++"
        :disabled="currentPage === totalPages"
        class="p-2 rounded-lg border border-gray-200 text-[#041983] hover:bg-[#041983] hover:text-[#f2c808] disabled:opacity-50 disabled:hover:bg-white disabled:hover:text-[#041983]"
      >
        <ChevronRightIcon class="w-5 h-5" />
      </button>
    </div>
  </div>
  <VueEasyLightbox
    :visible="visible"
    :imgs="images"
    :index="index"
    @hide="visible = false"
  />
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useFotoEkskulStore } from "../../stores/fotoekskul";
import VueEasyLightbox from "vue-easy-lightbox";

import {
  CalendarDaysIcon,
  MagnifyingGlassPlusIcon,
  MagnifyingGlassIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from "@heroicons/vue/24/outline";

const fotoStore = useFotoEkskulStore();

const keyword = ref("");

const visible = ref(false);
const index = ref(0);

const currentPage = ref(1);
const perPage = 8;

onMounted(() => {
  fotoStore.fetchAllFotoEkskul();
});

const filteredFoto = computed(() => {
  if (!keyword.value.trim()) {
    return fotoStore.fotoekskul;
  }

  const search = keyword.value.toLowerCase();

  return fotoStore.fotoekskul.filter(
    (item) =>
      item.nama.toLowerCase().includes(search) ||
      item.kegiatan.toLowerCase().includes(search),
  );
});

const totalPages = computed(() => {
  return Math.ceil(filteredFoto.value.length / perPage);
});

const paginatedFoto = computed(() => {
  const start = (currentPage.value - 1) * perPage;

  return filteredFoto.value.slice(start, start + perPage);
});

const images = computed(() => {
  return filteredFoto.value.map((item) => item.foto_url);
});

const showImage = (item) => {
  const idx = filteredFoto.value.findIndex((foto) => foto.id === item.id);

  index.value = idx;
  visible.value = true;
};

watch(keyword, () => {
  currentPage.value = 1;
});
</script>
