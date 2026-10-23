<template>
  <div
    v-if="pengumumanStore.detailPengumuman"
    class="bg-white rounded-2xl shadow-lg p-6 md:p-8"
  >
    <div class="grid grid-cols-1 xl:grid-cols-4 gap-8">
      <!-- KONTEN UTAMA -->
      <div class="xl:col-span-3">
        <!-- Judul -->
        <h1 class="text-xl font-bold text-[#041983] mb-3">
          {{ pengumumanStore.detailPengumuman.judul_pengumuman }}
        </h1>

        <!-- Meta -->
        <div
          class="flex flex-wrap items-center gap-4 mt-4 text-xs text-gray-500"
        >
          <span class="flex items-center gap-1.5">
            <CalendarDaysIcon class="w-3.5 h-3.5 shrink-0" />

            {{
              new Date(
                pengumumanStore.detailPengumuman.tanggal_publish,
              ).toLocaleDateString("id-ID", {
                day: "numeric",
                month: "long",
                year: "numeric",
              })
            }}
          </span>

          <span class="flex items-center gap-1.5">
            <UserIcon class="w-3.5 h-3.5 shrink-0" />
            {{ pengumumanStore.detailPengumuman.nama_pengguna }}
          </span>

          <span class="flex items-center gap-1.5">
            <EyeIcon class="w-3.5 h-3.5 shrink-0" />
            {{ formatHits(pengumumanStore.detailPengumuman.hits) }}
          </span>
        </div>

        <!-- Foto -->
        <div class="mt-8 flex justify-center">
          <img
            loading="lazy"
            :src="pengumumanStore.detailPengumuman.foto_url"
            :alt="pengumumanStore.detailPengumuman.judul_pengumuman"
            class="max-w-lg w-full rounded-xl shadow-md"
          />
        </div>

        <!-- Isi Pengumuman -->
        <div
          class="pengumuman-content mt-8 text-gray-700"
          v-html="isiPengumumanHtml"
        ></div>
      </div>

      <aside class="xl:border-l border-gray-100 xl:pl-6">
        <!-- BERITA -->
        <div class="mb-8">
          <div class="flex items-center justify-between mb-4">
            <div class="inline-block">
              <h2 class="text-xl font-bold text-gray-600 pb-2">Berita</h2>

              <div class="h-0.5 flex overflow-hidden rounded-full">
                <div class="w-1/2 bg-[#041983]"></div>
                <div class="w-1/2 bg-[#f2c808]"></div>
              </div>
            </div>

            <router-link
              to="/berita"
              class="text-sm font-medium text-gray-700 mt-3 inline-flex items-center gap-1 italic hover:gap-2 transition-all duration-300"
            >
              Lainnya
              <span>→</span>
            </router-link>
          </div>

          <div
            v-for="item in beritaLainnya"
            :key="item.id_berita"
            class="mb-4 pb-4 border-b border-gray-200 last:border-0"
          >
            <router-link
              :to="`/berita/${item.slug_berita}`"
              class="flex gap-3 group"
            >
              <img
                loading="lazy"
                :src="item.foto_url"
                :alt="item.judul_berita"
                class="w-20 h-16 rounded-lg object-cover shrink-0"
              />

              <div class="flex-1">
                <div
                  class="flex items-center gap-1.5 text-xs text-gray-500 mb-1"
                >
                  <CalendarDaysIcon class="w-3.5 h-3.5 shrink-0" />

                  <span>
                    {{
                      new Date(item.tanggal_publish).toLocaleDateString(
                        "id-ID",
                        {
                          day: "numeric",
                          month: "long",
                          year: "numeric",
                        },
                      )
                    }}
                  </span>
                </div>

                <div
                  class="font-medium text-sm text-[#041983] group-hover:text-[#f2c808] line-clamp-2"
                >
                  {{ item.judul_berita }}
                </div>
              </div>
            </router-link>
          </div>
        </div>

        <!-- PENGUMUMAN -->
        <div>
          <div class="flex items-center justify-between mb-4">
            <div class="inline-block">
              <h2 class="text-xl font-bold text-gray-600 pb-2">Pengumuman</h2>

              <div class="h-0.5 flex overflow-hidden rounded-full">
                <div class="w-1/2 bg-[#041983]"></div>
                <div class="w-1/2 bg-[#f2c808]"></div>
              </div>
            </div>

            <router-link
              to="/pengumuman"
              class="text-sm font-medium text-gray-700 mt-3 inline-flex items-center gap-1 italic hover:gap-2 transition-all duration-300"
            >
              Lainnya
              <span>→</span>
            </router-link>
          </div>

          <div
            v-for="item in pengumumanLainnya"
            :key="item.id_pengumuman"
            class="mb-4 pb-4 border-b border-gray-200 last:border-0"
          >
            <router-link
              :to="`/pengumuman/${item.slug_pengumuman}`"
              class="flex gap-3 group"
            >
              <img
                loading="lazy"
                :src="item.foto_url"
                :alt="item.judul_pengumuman"
                class="w-20 h-16 rounded-lg object-cover shrink-0"
              />

              <div class="flex-1">
                <div
                  class="flex items-center gap-1.5 text-xs text-gray-500 mb-1"
                >
                  <CalendarDaysIcon class="w-3.5 h-3.5 shrink-0" />

                  <span>
                    {{
                      new Date(item.tanggal_publish).toLocaleDateString(
                        "id-ID",
                        {
                          day: "numeric",
                          month: "long",
                          year: "numeric",
                        },
                      )
                    }}
                  </span>
                </div>

                <div
                  class="font-medium text-sm text-[#041983] group-hover:text-[#f2c808] line-clamp-2"
                >
                  {{ item.judul_pengumuman }}
                </div>
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
import { usePengumumanStore } from "../../stores/pengumuman";
import { formatHits } from "../../utils/format";
import { CalendarDaysIcon, EyeIcon, UserIcon } from "@heroicons/vue/24/outline";

const route = useRoute();

const beritaStore = useBeritaStore();
const pengumumanStore = usePengumumanStore();

const beritaLainnya = computed(() => {
  return beritaStore.berita
    .filter((item) => item.slug_berita !== route.params.slug)
    .slice(0, 5);
});

const pengumumanLainnya = computed(() => {
  return pengumumanStore.pengumuman
    .filter((item) => item.slug_pengumuman !== route.params.slug)
    .slice(0, 5);
});

// Load detail pengumuman + hit counter
const loadPengumuman = async (slug) => {
  await pengumumanStore.fetchDetailPengumuman(slug);

  const pengumuman = pengumumanStore.detailPengumuman;

  if (!pengumuman) return;

  const key = `pengumuman_${pengumuman.id_pengumuman}`;

  if (!sessionStorage.getItem(key)) {
    await pengumumanStore.incrementHits(pengumuman.id_pengumuman);

    sessionStorage.setItem(key, "1");
  }
};

const isiPengumumanHtml = computed(() => {
  const html = pengumumanStore.detailPengumuman?.isi_pengumuman;

  if (!html) return "";

  const wrapper = document.createElement("div");
  wrapper.innerHTML = html;

  wrapper.querySelectorAll("table").forEach((table) => {
    const scrollDiv = document.createElement("div");
    scrollDiv.className = "table-scroll";

    table.parentNode.insertBefore(scrollDiv, table);
    scrollDiv.appendChild(table);
  });

  return wrapper.innerHTML;
});

// Saat slug berubah
watch(
  () => route.params.slug,
  async (newSlug) => {
    if (!newSlug) return;

    await loadPengumuman(newSlug);

    window.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  },
);

// Pertama kali halaman dibuka
onMounted(async () => {
  await Promise.all([
    loadPengumuman(route.params.slug),
    pengumumanStore.fetchPengumuman(),
    beritaStore.fetchBerita(),
  ]);
});
</script>

<style scoped>
:deep(.pengumuman-content p) {
  margin-bottom: 1rem;
  line-height: 1.9;
  /* text-align: justify; */
}

:deep(.pengumuman-content img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.75rem;
}

:deep(.pengumuman-content img[data-align="left"]) {
  float: left;
  margin: 0.5rem 1rem 1rem 0;
}

:deep(.pengumuman-content img[data-align="center"]) {
  display: block;
  margin-left: auto;
  margin-right: auto;
}

:deep(.pengumuman-content img[data-align="right"]) {
  float: right;
  margin: 0.5rem 0 1rem 1rem;
}

:deep(.pengumuman-content h1),
:deep(.pengumuman-content h2),
:deep(.pengumuman-content h3) {
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
  font-weight: 700;
  color: #041983;
}

:deep(.pengumuman-content ul) {
  list-style-type: disc;
  padding-left: 1.5rem;
  margin-bottom: 1rem;
}

:deep(.pengumuman-content ol) {
  list-style-type: decimal;
  padding-left: 1.5rem;
  margin-bottom: 1rem;
}

:deep(.pengumuman-content li) {
  margin-bottom: 0.5rem;
}

:deep(.pengumuman-content pre) {
  white-space: normal;
  overflow: visible;
  background: transparent;
  padding: 0;
  margin: 0;
  font-family: inherit;
}

:deep(.pengumuman-content table) {
  width: max-content;
  min-width: 100%;
  border-collapse: collapse;
  margin: 0;
}

:deep(.pengumuman-content th),
:deep(.pengumuman-content td) {
  padding: 0.35rem 0.6rem;
  border: 1px solid #d1d5db;
  vertical-align: top;
}

:deep(.pengumuman-content th) {
  background: #f9fafb;
  font-weight: 600;
}

:deep(.pengumuman-content th p),
:deep(.pengumuman-content td p) {
  margin: 0;
}

:deep(.dark .pengumuman-content th),
:deep(.dark .pengumuman-content td) {
  border-color: #334155;
}

:deep(.dark .pengumuman-content th) {
  background: #1e293b;
}

:deep(.pengumuman-content a) {
  color: rgb(37 99 235); /* blue-600 */
  text-decoration: none;
}

:deep(.pengumuman-content a:hover) {
  color: rgb(29 78 216); /* blue-700 */
  text-decoration: underline;
}

:deep(.pengumuman-content) {
  overflow-wrap: break-word;
}

:deep(.table-scroll) {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  margin: 0 0 1rem 0;
}

:deep(.table-scroll table) {
  margin: 0;
}
</style>
