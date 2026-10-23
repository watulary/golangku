<template>
  <DashboardPageHeader
    title="Guru Wali Saya"
    subtitle="Informasi guru wali kelasmu."
  />

  <div
    class="relative overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
  >
    <div
      class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
    ></div>

    <!-- Loading -->
    <div
      v-if="guruwaliStore.siswaLoading"
      class="flex flex-col items-center justify-center py-20"
    >
      <div
        class="h-10 w-10 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"
      />

      <p class="mt-4 text-sm text-gray-500 dark:text-slate-400">
        Memuat data...
      </p>
    </div>

    <!-- Empty -->
    <div
      v-else-if="guruwaliStore.siswaList.length === 0"
      class="py-16 text-center"
    >
      <UserGroupIcon
        class="mx-auto h-10 w-10 text-gray-400 dark:text-slate-500"
      />
      <p class="mt-3 text-sm text-gray-500 dark:text-slate-400">
        Kamu belum terdaftar di guru wali manapun.
      </p>
    </div>

    <!-- List -->
    <ul v-else class="divide-y divide-gray-100 dark:divide-white/5">
      <li
        v-for="item in guruwaliStore.siswaList"
        :key="item.id_guruwali"
        class="flex items-center justify-between gap-4 px-6 py-4"
      >
        <div class="flex items-center gap-3">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-blue-50 text-blue-600 dark:bg-blue-500/10 dark:text-blue-400"
          >
            <UserGroupIcon class="h-5 w-5" />
          </div>

          <div>
            <p class="text-sm font-medium text-gray-800 dark:text-slate-200">
              {{ item.nama_guruwali }}
            </p>
            <p class="text-xs text-gray-500 dark:text-slate-400">
              Guru: {{ item.nama_pegawai || "—" }}
            </p>
          </div>
        </div>

        <span
          :class="[
            'inline-flex shrink-0 items-center rounded-full px-2.5 py-1 text-xs font-medium',
            item.aktif
              ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400'
              : 'bg-gray-100 text-gray-600 dark:bg-slate-800 dark:text-slate-400',
          ]"
        >
          {{ item.aktif ? "Aktif" : "Nonaktif" }}
        </span>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { onMounted } from "vue";

import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";

import { useGuruWaliStore } from "@/stores/guruwali";

import { UserGroupIcon } from "@heroicons/vue/24/outline";

const guruwaliStore = useGuruWaliStore();

onMounted(() => {
  guruwaliStore.fetchMe();
});
</script>
