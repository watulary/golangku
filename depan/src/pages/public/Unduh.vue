<template>
  <div class="bg-white rounded-2xl shadow-lg p-6 md:p-8">
    <!-- Header + Search -->
    <div
      class="flex flex-col md:flex-row md:items-end md:justify-between gap-4 mb-8"
    >
      <div class="self-start">
        <h2 class="text-2xl font-extrabold text-gray-600">Unduh Dokumen</h2>

        <p class="text-sm text-gray-500 mt-2">
          Total: {{ filteredData.length }}
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

    <!-- Loading -->
    <div v-if="unggahdokumenStore.loading" class="text-center py-16">
      <p class="text-gray-500">Memuat data...</p>
    </div>

    <!-- Empty -->
    <div v-else-if="filteredData.length === 0" class="text-center py-16">
      <p class="text-gray-500">Data tidak ditemukan.</p>
    </div>

    <!-- Table -->
    <div v-else class="overflow-x-auto">
      <table class="w-full text-sm border-collapse">
        <thead>
          <tr class="bg-slate-50 text-gray-700 border-b border-slate-200">
            <th class="p-3 text-center">No</th>
            <th class="p-3 text-left">Dokumen</th>
            <th class="p-3 text-left">Oleh</th>
            <th class="p-3 text-left">Tanggal</th>
            <th class="p-3 text-left">Unduh</th>
          </tr>
        </thead>

        <tbody>
          <tr
            v-for="(item, index) in paginatedData"
            :key="item.id"
            class="border-b border-gray-100 hover:bg-gray-50"
          >
            <td class="p-3 text-center text-gray-700">
              {{ (currentPage - 1) * perPage + index + 1 }}
            </td>

            <td class="p-3 text-gray-700">
              {{ item.nama }}
            </td>

            <td class="p-3 text-gray-700">
              {{ item.nama_pengguna }}
            </td>

            <td class="p-3 text-gray-700">
              {{
                new Date(item.waktu).toLocaleDateString("id-ID", {
                  day: "numeric",
                  month: "long",
                  year: "numeric",
                })
              }}
            </td>

            <td class="p-3 text-center text-gray-700">
              <a
                :href="item.dokumen_url"
                title="Unduh Dokumen"
                class="inline-flex items-center justify-center w-9 h-9 rounded-lg border border-gray-300 text-gray-600 hover:border-[#041983] hover:text-[#041983] transition"
              >
                <ArrowDownTrayIcon class="w-5 h-5" />
              </a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div
      v-if="totalPages > 1"
      class="flex justify-center items-center gap-2 mt-10"
    >
      <button
        @click="currentPage--"
        :disabled="currentPage === 1"
        class="p-2 rounded-lg border border-gray-200 text-[#041983] hover:bg-[#041983] hover:text-[#f2c808] disabled:opacity-50"
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
        class="p-2 rounded-lg border border-gray-200 text-[#041983] hover:bg-[#041983] hover:text-[#f2c808] disabled:opacity-50"
      >
        <ChevronRightIcon class="w-5 h-5" />
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useUnggahdokumenStore } from "../../stores/unggahdokumen";

import {
  MagnifyingGlassIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ArrowDownTrayIcon,
} from "@heroicons/vue/24/outline";

const unggahdokumenStore = useUnggahdokumenStore();

const keyword = ref("");

// pagination
const currentPage = ref(1);
const perPage = 20;

onMounted(async () => {
  await unggahdokumenStore.fetchData();
});

const filteredData = computed(() => {
  if (!keyword.value.trim()) {
    return unggahdokumenStore.data;
  }

  const q = keyword.value.toLowerCase();

  return unggahdokumenStore.data.filter((item) => {
    const tanggal = new Date(item.waktu)
      .toLocaleDateString("id-ID", {
        day: "numeric",
        month: "long",
        year: "numeric",
      })
      .toLowerCase();

    return (
      item.nama?.toLowerCase().includes(q) ||
      item.nama_pengguna?.toLowerCase().includes(q) ||
      tanggal.includes(q)
    );
  });
});

const totalPages = computed(() => {
  return Math.ceil(filteredData.value.length / perPage);
});

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * perPage;

  return filteredData.value.slice(start, start + perPage);
});

watch(keyword, () => {
  currentPage.value = 1;
});
</script>
