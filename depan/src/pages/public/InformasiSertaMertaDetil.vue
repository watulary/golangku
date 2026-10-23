<template>
  <div
    v-if="beritaStore.detailBerita"
    class="bg-white rounded-2xl shadow-lg p-6 md:p-8"
  >
    <div class="grid grid-cols-1 xl:grid-cols-4 gap-8">
      <!-- KONTEN -->
      <div class="xl:col-span-3">
        <!-- Judul -->
        <h1 class="text-xl font-bold text-[#041983] mb-3">
          {{ beritaStore.detailBerita.judul_berita }}
        </h1>

        <!-- Meta -->
        <div
          class="flex flex-wrap items-center gap-4 mt-4 text-xs text-gray-500"
        >
          <span class="flex items-center gap-1.5">
            <CalendarDaysIcon class="w-3.5 h-3.5" />

            {{
              new Date(
                beritaStore.detailBerita.tanggal_publish,
              ).toLocaleDateString("id-ID", {
                day: "numeric",
                month: "long",
                year: "numeric",
              })
            }}
          </span>

          <span class="flex items-center gap-1.5">
            <UserIcon class="w-3.5 h-3.5" />
            {{ beritaStore.detailBerita.nama_pengguna }}
          </span>

          <span class="flex items-center gap-1.5">
            <EyeIcon class="w-3.5 h-3.5" />
            {{ formatHits(beritaStore.detailBerita.hits) }}
          </span>
        </div>

        <!-- Foto -->
        <div
          v-if="beritaStore.detailBerita.foto_url"
          class="mt-8 flex justify-center"
        >
          <img
            loading="lazy"
            :src="beritaStore.detailBerita.foto_url"
            :alt="beritaStore.detailBerita.judul_berita"
            class="max-w-lg w-full rounded-xl shadow-md"
          />
        </div>

        <!-- Isi -->
        <div
          class="berita-content mt-8 text-gray-700"
          v-html="beritaStore.detailBerita.isi_berita"
        ></div>
      </div>

      <!-- SIDEBAR -->
      <aside class="xl:border-l border-gray-100 xl:pl-6">
        <div>
          <div class="inline-block mb-4">
            <h2 class="text-xl font-bold text-gray-600 pb-2">
              Informasi Serta Merta
            </h2>

            <div class="h-0.5 flex overflow-hidden rounded-full">
              <div class="w-1/2 bg-[#041983]"></div>
              <div class="w-1/2 bg-[#f2c808]"></div>
            </div>
          </div>

          <div
            v-for="item in informasiLainnya"
            :key="item.id_berita"
            class="border-b border-gray-200 py-4 last:border-0"
          >
            <router-link
              :to="`/ppid/informasi-serta-merta/${item.slug_berita}`"
              class="group block"
            >
              <div
                class="text-sm font-medium text-[#041983] group-hover:text-[#f2c808] line-clamp-2"
              >
                {{ item.judul_berita }}
              </div>

              <div class="flex items-center gap-1.5 text-xs text-gray-500 mt-2">
                <CalendarDaysIcon class="w-3.5 h-3.5" />

                {{
                  new Date(item.tanggal_publish).toLocaleDateString("id-ID", {
                    day: "numeric",
                    month: "long",
                    year: "numeric",
                  })
                }}
              </div>
            </router-link>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { computed, watch, onMounted } from "vue";
import { useRoute } from "vue-router";
import { useBeritaStore } from "../../stores/berita";

import { formatHits } from "../../utils/format";

import { CalendarDaysIcon, EyeIcon, UserIcon } from "@heroicons/vue/24/outline";

const route = useRoute();
const beritaStore = useBeritaStore();

const informasiLainnya = computed(() => {
  return beritaStore.informasiSertaMerta
    .filter((item) => item.slug_berita !== route.params.slug)
    .slice(0, 5);
});

const loadInformasi = async (slug) => {
  await beritaStore.fetchDetailBerita(slug);

  const data = beritaStore.detailBerita;

  if (!data) return;

  const key = `informasi_serta_merta_${data.id_berita}`;

  if (!sessionStorage.getItem(key)) {
    await beritaStore.incrementHits(data.id_berita);

    sessionStorage.setItem(key, "1");
  }
};

watch(
  () => route.params.slug,
  async (newSlug) => {
    await loadInformasi(newSlug);

    window.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  },
);

onMounted(async () => {
  await Promise.all([
    loadInformasi(route.params.slug),
    beritaStore.fetchInformasiSertaMerta(),
  ]);
});
</script>

<style scoped>
:deep(.berita-content p) {
  margin-bottom: 1rem;
  line-height: 1.9;
}

:deep(.berita-content img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.75rem;
}

:deep(.berita-content img[data-align="left"]) {
  float: left;
  margin: 0.5rem 1rem 1rem 0;
}

:deep(.berita-content img[data-align="center"]) {
  display: block;
  margin-left: auto;
  margin-right: auto;
}

:deep(.berita-content img[data-align="right"]) {
  float: right;
  margin: 0.5rem 0 1rem 1rem;
}

:deep(.berita-content h1),
:deep(.berita-content h2),
:deep(.berita-content h3) {
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
  font-weight: 700;
  color: #041983;
}

:deep(.berita-content ul) {
  list-style-type: disc;
  padding-left: 1.5rem;
  margin-bottom: 1rem;
}

:deep(.berita-content ol) {
  list-style-type: decimal;
  padding-left: 1.5rem;
  margin-bottom: 1rem;
}

:deep(.berita-content li) {
  margin-bottom: 0.5rem;
}

:deep(.berita-content a) {
  color: rgb(37 99 235);
  text-decoration: none;
}

:deep(.berita-content a:hover) {
  color: rgb(29 78 216);
}

:deep(.berita-content table) {
  width: 100%;
  margin: 1rem 0;
  border-collapse: collapse;
}

:deep(.berita-content th),
:deep(.berita-content td) {
  padding: 0.75rem;
  border: 1px solid #d1d5db;
}

:deep(.berita-content th) {
  background: #f9fafb;
  font-weight: 600;
}

:deep(.dark .berita-content th),
:deep(.dark .berita-content td) {
  border-color: #334155;
}

:deep(.dark .berita-content th) {
  background: #1e293b;
}

:deep(.berita-content table) {
  width: max-content;
  min-width: 100%;
  border-collapse: collapse;
}

:deep(.berita-content th),
:deep(.berita-content td) {
  padding: 0.35rem 0.6rem;
  border: 1px solid #d1d5db;
  vertical-align: top;
}

:deep(.berita-content th p),
:deep(.berita-content td p) {
  margin: 0;
}

:deep(.table-scroll) {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  margin: 1rem 0;
}

:deep(.table-scroll table) {
  margin: 0;
}
</style>
