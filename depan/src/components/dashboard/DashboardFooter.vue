<template>
  <footer
    class="flex h-16 items-center justify-center border-t border-gray-100 bg-white px-4 transition-colors duration-300 dark:border-white/5 dark:bg-slate-900 lg:px-6"
  >
    <div
      class="flex flex-row flex-wrap items-center justify-center gap-x-1.5 gap-y-1 text-center text-xs text-gray-400 dark:text-slate-400"
    >
      <div class="flex items-center gap-1">
        <span>© {{ currentYear }} With</span>
        <HeartIcon
          class="h-4 w-4 shrink-0 animate-pulse text-[#FB7185] drop-shadow-xs"
        />
        <span>Tim Tendik</span>
      </div>
      <span class="hidden sm:inline">•</span>
      <span>
        {{ settingStore.setting?.nama_sekolah }}
      </span>
    </div>
  </footer>
</template>

<script setup>
import { computed } from "vue";
import { useSettingStore } from "@/stores/setting";
import { HeartIcon } from "@heroicons/vue/24/solid";

const settingStore = useSettingStore();

const currentYear = computed(() => new Date().getFullYear());
</script>
