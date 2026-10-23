<template>
  <div v-if="!isHome" class="mb-3">
    <nav class="flex flex-wrap items-center gap-1.5 text-sm font-medium">
      <template v-for="(item, index) in breadcrumbItems" :key="index">
        <!-- Separator -->
        <span v-if="index > 0" class="text-gray-300 dark:text-slate-600">
          /
        </span>

        <!-- Segment -->
        <RouterLink
          v-if="item.to && index !== breadcrumbItems.length - 1"
          :to="item.to"
          class="flex items-center gap-1 rounded-sm text-gray-500 transition hover:text-[#041983] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:text-slate-400 dark:hover:text-indigo-300"
        >
          <HomeIcon v-if="index === 0" class="h-3.5 w-3.5 shrink-0" />
          {{ item.label }}
        </RouterLink>

        <span
          v-else
          class="flex items-center gap-1 font-semibold text-gray-800 dark:text-white"
        >
          <HomeIcon v-if="index === 0" class="h-3.5 w-3.5 shrink-0" />
          {{ item.label }}
        </span>
      </template>
    </nav>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { useRoute } from "vue-router";
import { HomeIcon } from "@heroicons/vue/24/outline";

const route = useRoute();

const isHome = computed(() => route.path === "/dashboard");

const breadcrumbItems = computed(() => {
  // Jika route memiliki breadcrumb khusus
  if (route.meta.breadcrumb) {
    return route.meta.breadcrumb;
  }

  return route.matched
    .filter((record) => record.meta?.title)
    .map((record) => ({
      label: record.meta.title,
      to: record.path,
    }))
    .filter(
      (item, index, arr) => index === 0 || item.label !== arr[index - 1].label,
    );
});
</script>
