<template>
  <!-- SECTION: Carousel -->
  <section class="mb-12">
    <div
      class="relative overflow-hidden rounded-3xl shadow-xl h-62.5 md:h-87.5"
    >
      <!-- Slider -->
      <div
        :style="{ transform: `translateX(-${currentIndex * 100}%)` }"
        class="flex transition-transform duration-500 ease-in-out h-full"
      >
        <div
          v-for="(slide, index) in slides"
          :key="index"
          class="shrink-0 w-full h-full"
        >
          <img
            :src="slide.url"
            :alt="slide.title"
            class="w-full h-full object-cover"
          />
        </div>
      </div>

      <!-- Overlay -->
      <div
        class="absolute inset-0 bg-linear-to-r from-black/25 via-black/15 to-black/5"
      ></div>

      <!-- Indicator -->
      <div class="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2">
        <button
          v-for="(_, index) in slides"
          :key="index"
          @click="currentIndex = index"
          class="h-2 rounded-full transition-all duration-300"
          :class="
            currentIndex === index ? 'w-6 bg-[#f2c808]' : 'w-2 bg-white/70'
          "
        />
      </div>
    </div>
  </section>

  <!-- QUICK ACCESS -->
  <section class="mb-12 -mt-8 relative z-20">
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <router-link
        to="/"
        class="bg-white rounded-2xl shadow-lg p-5 text-center hover:-translate-y-1 hover:shadow-xl transition"
      >
        <div class="text-3xl mb-2">
          <BuildingOffice2Icon class="w-10 h-10 mx-auto text-green-500" />
        </div>
        <div class="font-semibold text-gray-700">
          {{ statistikStore.statistik?.jumlah_rombel || 0 }} Rombel
        </div>
      </router-link>

      <router-link
        to="/"
        class="bg-white rounded-2xl shadow-lg p-5 text-center hover:-translate-y-1 hover:shadow-xl transition"
      >
        <div class="text-3xl mb-2">
          <UserGroupIcon class="w-10 h-10 mx-auto text-blue-500" />
        </div>
        <div class="font-semibold text-gray-700">
          {{ statistikStore.statistik?.jumlah_murid || 0 }} Murid
        </div>
      </router-link>

      <router-link
        to="/"
        class="bg-white rounded-2xl shadow-lg p-5 text-center hover:-translate-y-1 hover:shadow-xl transition"
      >
        <div class="text-3xl mb-2">
          <AcademicCapIcon class="w-10 h-10 mx-auto text-red-500" />
        </div>
        <div class="font-semibold text-gray-700">
          {{ statistikStore.statistik?.jumlah_guru || 0 }} Guru
        </div>
      </router-link>

      <router-link
        to="/"
        class="bg-white rounded-2xl shadow-lg p-5 text-center hover:-translate-y-1 hover:shadow-xl transition"
      >
        <div class="text-3xl mb-2">
          <UsersIcon class="w-10 h-10 mx-auto text-[#f2c808]" />
        </div>
        <div class="font-semibold text-gray-700">
          {{ statistikStore.statistik?.jumlah_tendik || 0 }} Tendik
        </div>
      </router-link>
    </div>
  </section>

  <!-- SECTION: Profil Sekolah -->
  <section class="mb-12">
    <div class="inline-block">
      <h2 class="text-2xl font-extrabold text-gray-600 pb-2">Profil</h2>

      <div class="h-0.5 flex overflow-hidden rounded-full">
        <div class="w-1/2 bg-[#041983]"></div>
        <div class="w-1/2 bg-[#f2c808]"></div>
      </div>
    </div>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
      <div
        v-if="setting?.profil_sekolah"
        class="bg-white p-4 rounded-xl shadow-md"
      >
        <div
          @click="openVideo(setting.profil_sekolah)"
          class="block relative overflow-hidden rounded-lg group cursor-pointer"
        >
          <img
            loading="lazy"
            :src="`https://img.youtube.com/vi/${setting.profil_sekolah}/hqdefault.jpg`"
            alt="Profil Sekolah"
            class="w-full h-64 object-cover"
          />

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
      </div>

      <div class="bg-white p-6 rounded-xl shadow-md">
        <div class="flex gap-4 h-full">
          <div class="flex-1">
            <h3 class="text-xl font-bold text-[#041983] mb-3">
              {{ sambutanStore.sambutan?.judul_sambutan }}
            </h3>

            <p class="text-gray-700 leading-7 line-clamp-7 text-justify">
              {{ previewSambutan }}
            </p>

            <router-link
              to="/sambutan"
              class="text-blue-600 text-sm font-medium mt-3 inline-flex items-center gap-1 hover:gap-2 transition-all duration-300"
            >
              Baca Selengkapnya
              <span>→</span>
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- SECTION: Berita & Pengumuman -->
  <section class="mb-12">
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- BERITA -->
      <div class="lg:col-span-2">
        <div class="flex items-center justify-between mb-6">
          <div class="inline-block">
            <h2 class="text-2xl font-extrabold text-gray-600 pb-2">Berita</h2>

            <div class="h-0.5 flex overflow-hidden rounded-full">
              <div class="w-1/2 bg-[#041983]"></div>
              <div class="w-1/2 bg-[#f2c808]"></div>
            </div>
          </div>

          <router-link
            to="/berita"
            class="text-md font-medium text-gray-700 mt-3 inline-flex items-center gap-1 hover:gap-2 transition-all duration-300"
          >
            Lainnya
            <span>→</span>
          </router-link>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <router-link
            v-for="item in beritaStore.berita.slice(0, 6)"
            :key="item.id_berita"
            :to="`/berita/${item.slug_berita}`"
            class="group relative rounded-2xl overflow-hidden shadow-lg"
          >
            <img
              loading="lazy"
              :src="item.foto_url"
              :alt="item.judul_berita"
              class="h-64 w-full object-cover transition duration-500 group-hover:scale-110"
            />

            <!-- Overlay -->
            <div
              class="absolute inset-0 bg-linear-to-r from-black/25 via-black/15 to-black/5"
            ></div>

            <div class="absolute bottom-0 p-5 text-white">
              <div class="flex items-center gap-1.5 text-xs opacity-80 mb-2">
                <CalendarDaysIcon class="w-3.5 h-3.5" />

                <span>
                  {{
                    new Date(item.tanggal_publish).toLocaleDateString("id-ID", {
                      day: "numeric",
                      month: "long",
                      year: "numeric",
                    })
                  }}
                </span>
              </div>
              <h4 class="font-bold text-lg line-clamp-2">
                {{ item.judul_berita }}
              </h4>
            </div>
          </router-link>
        </div>
      </div>

      <!-- PENGUMUMAN -->
      <div>
        <div class="flex items-center justify-between mb-6">
          <div class="inline-block">
            <h2 class="text-2xl font-extrabold text-gray-600 pb-2">
              Pengumuman
            </h2>

            <div class="h-0.5 flex overflow-hidden rounded-full">
              <div class="w-1/2 bg-[#041983]"></div>
              <div class="w-1/2 bg-[#f2c808]"></div>
            </div>
          </div>

          <router-link
            to="/pengumuman"
            class="text-md font-medium text-gray-700 mt-3 inline-flex items-center gap-1 hover:gap-2 transition-all duration-300"
          >
            Lainnya
            <span>→</span>
          </router-link>
        </div>

        <div class="bg-white rounded-2xl shadow-lg overflow-hidden">
          <router-link
            v-for="item in pengumumanStore.pengumuman.slice(0, 5)"
            :key="item.id_pengumuman"
            :to="`/pengumuman/${item.slug_pengumuman}`"
            class="block p-4 border-b border-gray-100 hover:bg-gray-50 transition"
          >
            <div class="flex items-center gap-1.5 text-xs opacity-80 mb-2">
              <CalendarDaysIcon class="w-3.5 h-3.5" />

              <span>
                {{
                  new Date(item.tanggal_publish).toLocaleDateString("id-ID", {
                    day: "numeric",
                    month: "long",
                    year: "numeric",
                  })
                }}
              </span>
            </div>
            <h4
              class="font-semibold text-[#041983] line-clamp-2 hover:text-[#f2c808]"
            >
              {{ item.judul_pengumuman }}
            </h4>
          </router-link>
        </div>
      </div>
    </div>
  </section>

  <!-- SECTION: Foto -->
  <section class="mb-12">
    <div class="flex items-center justify-between mb-6">
      <div class="inline-block">
        <h2 class="text-2xl font-extrabold text-gray-600 pb-2">Galeri Foto</h2>

        <div class="h-0.5 flex overflow-hidden rounded-full">
          <div class="w-1/2 bg-[#041983]"></div>
          <div class="w-1/2 bg-[#f2c808]"></div>
        </div>
      </div>

      <router-link
        to="/foto"
        class="text-md font-medium text-gray-700 mt-3 inline-flex items-center gap-1 hover:gap-2 transition-all duration-300"
      >
        Lainnya
        <span>→</span>
      </router-link>
    </div>

    <div
      v-if="fotoEkskulStore.fotoekskul.length >= 3"
      class="grid grid-cols-1 md:grid-cols-3 gap-6"
    >
      <!-- FOTO BESAR -->
      <div
        @click="showImage(0)"
        class="md:col-span-2 group relative overflow-hidden rounded-2xl shadow-lg cursor-pointer"
      >
        <img
          loading="lazy"
          :src="fotoEkskulStore.fotoekskul[0].foto_url"
          :alt="fotoEkskulStore.fotoekskul[0].nama"
          class="w-full aspect-4/3 object-cover transition duration-500 group-hover:scale-105"
        />

        <!-- Overlay -->
        <div
          class="absolute inset-0 bg-linear-to-r from-black/25 via-black/15 to-black/5"
        ></div>

        <!-- Overlay Perbesar -->
        <div
          class="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition duration-300"
        >
          <span
            class="flex items-center gap-2 text-white text-sm font-medium px-4 py-2 rounded-full border border-white backdrop-blur-sm"
          >
            <MagnifyingGlassPlusIcon class="w-5 h-5" />
            Perbesar
          </span>
        </div>

        <div class="absolute bottom-0 left-0 p-5 text-white">
          <div class="flex items-center gap-1.5 text-xs opacity-80 mb-2">
            <CalendarDaysIcon class="w-3.5 h-3.5" />

            <span>
              {{
                new Date(
                  fotoEkskulStore.fotoekskul[0]?.waktu,
                ).toLocaleDateString("id-ID", {
                  day: "numeric",
                  month: "long",
                  year: "numeric",
                })
              }}
            </span>
          </div>
          <h3 class="font-bold text-lg">
            {{ fotoEkskulStore.fotoekskul[0].nama }}
          </h3>

          <!-- <p class="text-sm opacity-90">
            {{ fotoEkskulStore.fotoekskul[0].kegiatan }}
          </p> -->
        </div>
      </div>

      <!-- FOTO KANAN -->
      <div class="grid grid-rows-2 gap-6">
        <div
          v-for="(item, idx) in fotoEkskulStore.fotoekskul.slice(1, 3)"
          :key="item.id"
          @click="showImage(idx + 1)"
          class="group relative overflow-hidden rounded-2xl shadow-lg cursor-pointer"
        >
          <img
            loading="lazy"
            :src="item.foto_url"
            :alt="item.nama"
            class="w-full aspect-4/3 object-cover transition duration-500 group-hover:scale-105"
          />

          <!-- Overlay -->
          <div
            class="absolute inset-0 bg-linear-to-r from-black/25 via-black/15 to-black/5"
          ></div>

          <div
            class="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition duration-300"
          >
            <span
              class="flex items-center gap-2 text-white text-sm font-medium px-4 py-2 rounded-full border border-white backdrop-blur-sm"
            >
              <MagnifyingGlassPlusIcon class="w-5 h-5" />
              Perbesar
            </span>
          </div>

          <div class="absolute bottom-0 left-0 p-4 text-white">
            <div class="flex items-center gap-1.5 text-xs opacity-80 mb-1">
              <CalendarDaysIcon class="w-3.5 h-3.5" />

              <span>
                {{
                  new Date(item.waktu).toLocaleDateString("id-ID", {
                    day: "numeric",
                    month: "long",
                    year: "numeric",
                  })
                }}
              </span>
            </div>
            <h4 class="font-bold text-lg">
              {{ item.nama }}
            </h4>

            <!-- <p class="text-xs opacity-90">
              {{ item.kegiatan }}
            </p> -->
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- SECTION: Video -->
  <section class="mb-4">
    <div class="flex items-center justify-between mb-6">
      <div class="inline-block">
        <h2 class="text-2xl font-extrabold text-gray-600 pb-2">Video</h2>

        <div class="h-0.5 flex overflow-hidden rounded-full">
          <div class="w-1/2 bg-[#041983]"></div>
          <div class="w-1/2 bg-[#f2c808]"></div>
        </div>
      </div>

      <router-link
        to="/video"
        class="text-md font-medium text-gray-700 mt-3 inline-flex items-center gap-1 hover:gap-2 transition-all duration-300"
      >
        Lainnya
        <span>→</span>
      </router-link>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div
        v-for="item in videoStore.videos"
        :key="item.id"
        @click="openVideo(item.video)"
        class="group relative overflow-hidden rounded-2xl shadow-lg cursor-pointer"
      >
        <img
          loading="lazy"
          :src="item.thumbnail"
          :alt="item.nama"
          class="h-64 w-full object-cover transition duration-500 group-hover:scale-110"
        />

        <!-- Overlay -->
        <div
          class="absolute inset-0 bg-linear-to-r from-black/25 via-black/15 to-black/5"
        ></div>

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

        <div
          class="absolute bottom-0 left-0 right-0 p-5 bg-linear-to-t from-black/90 to-transparent"
        >
          <div class="absolute bottom-0 left-0 p-4 text-white">
            <div class="flex items-center gap-1.5 text-xs opacity-80 mb-1">
              <CalendarDaysIcon class="w-3.5 h-3.5" />

              <span>
                {{
                  new Date(item.waktu).toLocaleDateString("id-ID", {
                    day: "numeric",
                    month: "long",
                    year: "numeric",
                  })
                }}
              </span>
            </div>
            <h4 class="text-white font-bold text-lg line-clamp-1">
              {{ item.nama }}
            </h4>
          </div>
        </div>
      </div>
    </div>
  </section>
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

  <VueEasyLightbox
    :visible="visible"
    :imgs="images"
    :index="index"
    @hide="visible = false"
  />
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
import { useFotoBannerStore } from "../../stores/fotobanner";
import { useSambutanStore } from "../../stores/sambutan";
import { useBeritaStore } from "../../stores/berita";
import { usePengumumanStore } from "../../stores/pengumuman";
import { useFotoEkskulStore } from "../../stores/fotoekskul";
import { useVideoStore } from "../../stores/video";
import { useSettingStore } from "../../stores/setting";
import { useStatistikStore } from "../../stores/statistik";
import VueEasyLightbox from "vue-easy-lightbox";
import {
  AcademicCapIcon,
  UsersIcon,
  UserGroupIcon,
  BuildingOffice2Icon,
  CalendarDaysIcon,
  MagnifyingGlassPlusIcon,
} from "@heroicons/vue/24/outline";

const bannerStore = useFotoBannerStore();
const sambutanStore = useSambutanStore();
const beritaStore = useBeritaStore();
const pengumumanStore = usePengumumanStore();
const fotoEkskulStore = useFotoEkskulStore();
const videoStore = useVideoStore();
const settingStore = useSettingStore();
const statistikStore = useStatistikStore();

//lightbox
const visible = ref(false);
const index = ref(0);

const images = computed(() => {
  return fotoEkskulStore.fotoekskul.map((item) => item.foto_url);
});

const showImage = (idx) => {
  index.value = idx;
  visible.value = true;
};
//end lightbox

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

const setting = computed(() => {
  return settingStore.setting || {};
});

const previewSambutan = computed(() => {
  const temp = document.createElement("div");
  temp.innerHTML = sambutanStore.sambutan?.isi_sambutan || "";

  return (temp.textContent || temp.innerText || "").replace(/\s+/g, " ").trim();
});

const slides = computed(() => {
  return bannerStore.banners.map((item) => ({
    url: item.gambar_url,
    title: item.nama,
    kegiatan: item.kegiatan,
  }));
});

watch(
  slides,
  (value) => {
    console.log("SLIDES =", value);
  },
  { immediate: true },
);

const currentIndex = ref(0);

let intervalId = null;

const nextSlide = () => {
  if (!slides.value.length) return;

  currentIndex.value = (currentIndex.value + 1) % slides.value.length;
};

const prevSlide = () => {
  if (!slides.value.length) return;

  currentIndex.value =
    (currentIndex.value - 1 + slides.value.length) % slides.value.length;
};

const startAutoPlay = () => {
  if (!slides.value.length) return;

  intervalId = setInterval(nextSlide, 5000);
};

const stopAutoPlay = () => {
  if (intervalId) clearInterval(intervalId);
};

onMounted(async () => {
  await Promise.all([
    bannerStore.fetchBanners(),
    sambutanStore.fetchSambutan(),
    beritaStore.fetchBerita(),
    pengumumanStore.fetchPengumuman(),
    fotoEkskulStore.fetchFotoEkskul(),
    videoStore.fetchVideos(),
    statistikStore.fetchStatistik(),
  ]);

  startAutoPlay();
});

onUnmounted(() => {
  stopAutoPlay();
});
</script>
