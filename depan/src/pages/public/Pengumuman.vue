<template>
  <div class="bg-white rounded-2xl shadow-lg p-6 md:p-8">
    <!-- Header + Search -->
    <div
      class="flex flex-col md:flex-row md:items-end md:justify-between gap-4 mb-8"
    >
      <div class="self-start">
        <h2 class="text-2xl font-extrabold text-gray-600">Daftar Pengumuman</h2>

        <p class="text-sm text-gray-500 mt-2">
          Total: {{ filteredPengumuman.length }}
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
    <div v-if="filteredPengumuman.length === 0" class="text-center py-10">
      <p class="text-sm text-gray-500">Data tidak ditemukan.</p>
    </div>

    <!-- Grid Pengumuman -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <RouterLink
        v-for="item in paginatedPengumuman"
        :key="item.id_pengumuman"
        :to="`/pengumuman/${item.slug_pengumuman}`"
        class="bg-white rounded-2xl overflow-hidden border border-gray-100 shadow-sm hover:shadow-lg transition-all duration-300"
      >
        <!-- Foto -->
        <img
          loading="lazy"
          :src="item.foto_url"
          :alt="item.judul_pengumuman"
          class="w-full h-72 object-cover"
        />

        <!-- Content -->
        <div class="p-5">
          <p class="flex items-center gap-1.5 text-xs text-gray-500 mb-3">
            <CalendarDaysIcon class="w-3.5 h-3.5 shrink-0" />

            {{
              new Date(item.tanggal_publish).toLocaleDateString("id-ID", {
                day: "numeric",
                month: "long",
                year: "numeric",
              })
            }}
          </p>

          <h2 class="font-semibold text-[#041983] line-clamp-2">
            {{ item.judul_pengumuman }}
          </h2>
        </div>
      </RouterLink>
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
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { usePengumumanStore } from "../../stores/pengumuman";

import {
  CalendarDaysIcon,
  MagnifyingGlassIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from "@heroicons/vue/24/outline";

const pengumumanStore = usePengumumanStore();

const keyword = ref("");

const currentPage = ref(1);
const perPage = 8;

onMounted(() => {
  pengumumanStore.fetchAllPengumuman();
});

const filteredPengumuman = computed(() => {
  if (!keyword.value.trim()) {
    return pengumumanStore.semuaPengumuman;
  }

  return pengumumanStore.semuaPengumuman.filter((item) =>
    item.judul_pengumuman.toLowerCase().includes(keyword.value.toLowerCase()),
  );
});

const totalPages = computed(() =>
  Math.ceil(filteredPengumuman.value.length / perPage),
);

const paginatedPengumuman = computed(() => {
  const start = (currentPage.value - 1) * perPage;

  return filteredPengumuman.value.slice(start, start + perPage);
});

watch(keyword, () => {
  currentPage.value = 1;
});
</script>
