<template>
  <div class="bg-white rounded-2xl shadow-lg p-6 md:p-8">
    <!-- Header + Search -->
    <div
      class="flex flex-col md:flex-row md:items-end md:justify-between gap-4 mb-8"
    >
      <div class="self-start">
        <h2 class="text-2xl font-extrabold text-gray-600">Daftar Video</h2>

        <p class="text-sm text-gray-500 mt-2">
          Total: {{ filteredVideos.length }}
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

    <!-- Empty -->
    <div v-if="filteredVideos.length === 0" class="text-center py-10">
      <p class="text-sm text-gray-500">Data tidak ditemukan.</p>
    </div>

    <!-- Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div
        v-for="item in paginatedVideos"
        :key="item.id"
        @click="openVideo(item.video)"
        class="group relative overflow-hidden rounded-2xl shadow-lg cursor-pointer"
      >
        <!-- Thumbnail -->
        <div class="relative">
          <img
            loading="lazy"
            :src="item.thumbnail"
            :alt="item.nama"
            class="w-full h-72 object-cover"
          />

          <!-- Play Button -->
          <div
            class="absolute inset-0 flex items-center justify-center bg-black/10"
          >
            <div
              class="w-12 h-12 rounded-full bg-red-600 flex items-center justify-center shadow-lg transition duration-300 group-hover:scale-110"
            >
              <svg
                class="w-6 h-6 text-white ml-0.5"
                fill="currentColor"
                viewBox="0 0 24 24"
              >
                <path d="M8 5v14l11-7z" />
              </svg>
            </div>
          </div>
        </div>

        <!-- Content -->
        <div class="p-5">
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

          <!-- <p class="text-sm text-gray-600 mt-2 line-clamp-2">
            {{ item.kegiatan }}
          </p> -->
        </div>
      </div>
    </div>

    <Teleport to="body">
      <div
        v-if="showVideoModal"
        class="fixed inset-0 z-50 bg-black/80 flex items-center justify-center p-4"
        @click="closeVideo"
      >
        <div class="relative w-full max-w-4xl" @click.stop>
          <button
            @click="closeVideo"
            class="absolute -top-10 right-0 text-white text-3xl"
          >
            ×
          </button>

          <div class="aspect-video">
            <iframe
              class="w-full h-full rounded-xl"
              :src="`https://www.youtube.com/embed/${selectedVideo}?autoplay=1`"
              title="YouTube video player"
              frameborder="0"
              allow="
                accelerometer;
                autoplay;
                clipboard-write;
                encrypted-media;
                gyroscope;
                picture-in-picture;
              "
              allowfullscreen
            ></iframe>
          </div>
        </div>
      </div>
    </Teleport>

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
        @click="currentPage = Math.min(currentPage + 1, totalPages)"
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
import { useVideoStore } from "../../stores/video";
import {
  CalendarDaysIcon,
  MagnifyingGlassIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from "@heroicons/vue/24/outline";

const videoStore = useVideoStore();

const keyword = ref("");

const currentPage = ref(1);
const perPage = 8;

//play youtube
const selectedVideo = ref(null);
const showVideoModal = ref(false);

const openVideo = (videoId) => {
  selectedVideo.value = videoId;
  showVideoModal.value = true;
};

const closeVideo = () => {
  showVideoModal.value = false;
  selectedVideo.value = null;
};
//end play youtube

onMounted(async () => {
  await videoStore.fetchAllVideos();
});

const filteredVideos = computed(() => {
  if (!keyword.value.trim()) {
    return videoStore.allVideos;
  }

  return videoStore.allVideos.filter(
    (item) =>
      item.nama?.toLowerCase().includes(keyword.value.toLowerCase()) ||
      item.kegiatan?.toLowerCase().includes(keyword.value.toLowerCase()),
  );
});

const totalPages = computed(() => {
  return Math.ceil(filteredVideos.value.length / perPage);
});

const paginatedVideos = computed(() => {
  const start = (currentPage.value - 1) * perPage;

  return filteredVideos.value.slice(start, start + perPage);
});

watch(keyword, () => {
  currentPage.value = 1;
});
</script>
