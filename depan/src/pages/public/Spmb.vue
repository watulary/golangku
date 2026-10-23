<template>
  <div class="bg-white rounded-2xl shadow-lg p-6 md:p-8">
    <!-- Header + Search -->
    <div
      class="flex flex-col md:flex-row md:items-end md:justify-between gap-4 mb-8"
    >
      <div class="self-start">
        <h2 class="text-2xl font-extrabold text-gray-600">
          Hasil Akhir SPMB 2026/2027
        </h2>

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
    <div v-if="spmbStore.loading" class="text-center py-16">
      <p class="text-gray-500">Memuat data...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredData.length === 0" class="text-center py-16">
      <p class="text-gray-500">Data tidak ditemukan.</p>
    </div>

    <!-- Table -->
    <div v-else class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="bg-slate-50 text-gray-700 border-b border-slate-200">
            <th class="p-3 text-center">No</th>
            <th class="p-3 text-left">No Peserta</th>
            <th class="p-3 text-left">Nama Peserta</th>
            <th class="p-3 text-center">L/P</th>
            <th class="p-3 text-left">Asal Sekolah</th>
            <th class="p-3 text-left">Jalur</th>
            <th class="p-3 text-center">Lapor Diri</th>
          </tr>
        </thead>

        <tbody>
          <tr
            v-for="(item, index) in paginatedData"
            :key="item.no_peserta"
            class="border-b border-gray-100 hover:bg-gray-50"
          >
            <td class="p-3 text-center text-gray-700">
              {{ (currentPage - 1) * perPage + index + 1 }}
            </td>

            <td class="p-3 text-center text-gray-700">
              {{ item.no_peserta }}
            </td>

            <td class="p-3 text-gray-700">
              {{ item.nama_peserta }}
            </td>

            <td class="p-3 text-center text-gray-700">
              {{ item.kelamin }}
            </td>

            <td class="p-3 text-gray-700">
              {{ item.asal_sekolah }}
            </td>

            <td class="p-3 text-gray-700">
              {{ getJalur(item.jalur) }}
            </td>

            <td class="p-3 text-center text-gray-700">
              <span
                class="px-3 py-1 rounded-full text-xs font-semibold"
                :class="getLaporDiri(item.lapor_diri).class"
              >
                {{ getLaporDiri(item.lapor_diri).text }}
              </span>
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
import { useSpmbStore } from "../../stores/spmb";
import {
  MagnifyingGlassIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from "@heroicons/vue/24/outline";

const spmbStore = useSpmbStore();

const keyword = ref("");

// Pagination
const currentPage = ref(1);
const perPage = 100;

onMounted(async () => {
  await spmbStore.fetchData();
});

const filteredData = computed(() => {
  if (!keyword.value.trim()) {
    return spmbStore.data;
  }

  const q = keyword.value.toLowerCase();

  return spmbStore.data.filter((item) => {
    const jalur = getJalur(item.jalur).toLowerCase();
    const laporDiri = getLaporDiri(item.lapor_diri).text.toLowerCase();

    return (
      item.no_peserta.toLowerCase().includes(q) ||
      item.nama_peserta.toLowerCase().includes(q) ||
      item.asal_sekolah.toLowerCase().includes(q) ||
      item.kelamin.toLowerCase().includes(q) ||
      jalur.includes(q) ||
      laporDiri.includes(q)
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

const getJalur = (jalur) => {
  switch (jalur) {
    case 1:
      return "Prestasi Akademik";
    case 2:
      return "Prestasi Non Akademik";
    case 3:
      return "Penyandang Disabilitas";
    case 4:
      return "KJP Plus";
    case 5:
      return "Domisili";
    case 6:
      return "Mutasi";
    case 7:
      return "Tahap Kedua";
    default:
      return "Tidak Diketahui";
  }
};

const getLaporDiri = (status) => {
  switch (status) {
    case 0:
      return {
        text: "Belum",
        class: "bg-yellow-100 text-yellow-800",
      };

    case 1:
      return {
        text: "Sudah",
        class: "bg-green-100 text-green-800",
      };

    case 2:
      return {
        text: "Undur",
        class: "bg-red-100 text-red-800",
      };

    default:
      return {
        text: "-",
        class: "bg-gray-100 text-gray-800",
      };
  }
};
</script>
