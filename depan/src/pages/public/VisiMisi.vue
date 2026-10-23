<template>
  <div class="bg-white rounded-2xl shadow-lg p-6 md:p-8">
    <div class="grid grid-cols-1 xl:grid-cols-4 gap-8">
      <!-- KONTEN UTAMA -->
      <div class="xl:col-span-3">
        <!-- Judul -->
        <h1 class="text-xl font-bold text-[#041983] mb-3">
          {{ visiMisiStore.visimisi?.judul }}
        </h1>

        <!-- Meta -->
        <div
          class="flex flex-wrap items-center gap-4 mt-4 text-xs text-gray-500"
        >
          <span class="flex items-center gap-1.5">
            <CalendarDaysIcon class="w-3.5 h-3.5 shrink-0" />

            {{
              new Date(visiMisiStore.visimisi?.waktu).toLocaleDateString(
                "id-ID",
                {
                  day: "numeric",
                  month: "long",
                  year: "numeric",
                },
              )
            }}
          </span>

          <span class="flex items-center gap-1.5">
            <UserIcon class="w-3.5 h-3.5 shrink-0" />
            {{ visiMisiStore.visimisi?.nama_pengguna }}
          </span>

          <span class="flex items-center gap-1.5">
            <EyeIcon class="w-3.5 h-3.5 shrink-0" />
            {{ formatHits(visiMisiStore.visimisi?.hits || 0) }}
          </span>
        </div>

        <!-- Isi visimisi -->
        <div
          class="visimisi-content mt-8 text-gray-700"
          v-html="visiMisiStore.visimisi?.keterangan"
        ></div>
      </div>

      <!-- SIDEBAR -->
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
            v-for="item in beritaStore.berita.slice(0, 5)"
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
                  <span>{{
                    new Date(item.tanggal_publish).toLocaleDateString("id-ID", {
                      day: "numeric",
                      month: "long",
                      year: "numeric",
                    })
                  }}</span>
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
            v-for="item in pengumumanStore.pengumuman.slice(0, 5)"
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
                  <span>{{
                    new Date(item.tanggal_publish).toLocaleDateString("id-ID", {
                      day: "numeric",
                      month: "long",
                      year: "numeric",
                    })
                  }}</span>
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
import { onMounted } from "vue";
import { useVisiMisiStore } from "../../stores/visimisi";
import { useBeritaStore } from "../../stores/berita";
import { usePengumumanStore } from "../../stores/pengumuman";
import { formatHits } from "../../utils/format";
import { UserIcon, CalendarDaysIcon, EyeIcon } from "@heroicons/vue/24/outline";

const visiMisiStore = useVisiMisiStore();
const beritaStore = useBeritaStore();
const pengumumanStore = usePengumumanStore();

onMounted(async () => {
  const viewed = sessionStorage.getItem("visimisi_viewed");

  if (!viewed) {
    await visiMisiStore.incrementHits();
    sessionStorage.setItem("visimisi_viewed", "1");
  }

  await Promise.all([
    visiMisiStore.fetchVisiMisi(),
    beritaStore.fetchBerita(),
    pengumumanStore.fetchPengumuman(),
  ]);
});
</script>

<style scoped>
:deep(.visimisi-content p) {
  margin-bottom: 1rem;
  line-height: 1.9;
  /* text-align: justify; */
}

:deep(.visimisi-content img) {
  max-width: 100%;
  height: auto;
  border-radius: 0.75rem;
}

:deep(.visimisi-content img[data-align="left"]) {
  float: left;
  margin: 0.5rem 1rem 1rem 0;
}

:deep(.visimisi-content img[data-align="center"]) {
  display: block;
  margin-left: auto;
  margin-right: auto;
}

:deep(.visimisi-content img[data-align="right"]) {
  float: right;
  margin: 0.5rem 0 1rem 1rem;
}

:deep(.visimisi-content h1),
:deep(.visimisi-content h2),
:deep(.visimisi-content h3) {
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
  font-weight: 700;
  color: #041983;
}

:deep(.visimisi-content ul) {
  list-style-type: disc;
  padding-left: 1.5rem;
  margin-bottom: 1rem;
}

:deep(.visimisi-content ol) {
  list-style-type: decimal;
  padding-left: 1.5rem;
  margin-bottom: 1rem;
}

:deep(.visimisi-content li) {
  margin-bottom: 0.5rem;
}

:deep(.visimisi-content a) {
  color: rgb(37 99 235); /* blue-600 */
  text-decoration: none;
}

:deep(.visimisi-content a:hover) {
  color: rgb(29 78 216); /* blue-700 */
}

:deep(.visimisi-content table) {
  width: 100%;
  margin: 1rem 0;
  border-collapse: collapse;
}

:deep(.visimisi-content th),
:deep(.visimisi-content td) {
  padding: 0.75rem;
  border: 1px solid #d1d5db;
}

:deep(.visimisi-content th) {
  background: #f9fafb;
  font-weight: 600;
}

:deep(.dark .visimisi-content th),
:deep(.dark .visimisi-content td) {
  border-color: #334155;
}

:deep(.dark .visimisi-content th) {
  background: #1e293b;
}

:deep(.visimisi-content table) {
  width: max-content;
  min-width: 100%;
  border-collapse: collapse;
}

:deep(.visimisi-content th),
:deep(.visimisi-content td) {
  padding: 0.35rem 0.6rem;
  border: 1px solid #d1d5db;
  vertical-align: top;
}

:deep(.visimisi-content th p),
:deep(.visimisi-content td p) {
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
