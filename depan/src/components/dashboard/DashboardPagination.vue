<template>
  <div
    v-if="pagination.total > 0"
    class="flex flex-col gap-4 border-t border-gray-100 px-6 py-4 dark:border-white/5 lg:flex-row lg:items-center lg:justify-between"
  >
    <!-- Info -->
    <div class="text-sm text-gray-500 dark:text-slate-400">
      Menampilkan
      <span
        class="font-semibold tabular-nums text-gray-700 dark:text-slate-200"
      >
        {{ startItem }}
      </span>
      –
      <span
        class="font-semibold tabular-nums text-gray-700 dark:text-slate-200"
      >
        {{ endItem }}
      </span>
      dari
      <span
        class="font-semibold tabular-nums text-gray-700 dark:text-slate-200"
      >
        {{ pagination.total }}
      </span>
      data
    </div>

    <!-- Pagination -->
    <div class="flex items-center justify-center gap-1">
      <!-- Previous -->
      <button
        class="inline-flex h-10 items-center gap-2 rounded-xl border px-3 py-2 text-sm font-medium transition disabled:cursor-not-allowed disabled:opacity-50 border-gray-200 bg-white text-gray-700 hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:hover:bg-slate-700"
        :disabled="pagination.page === 1"
        @click="$emit('change-page', pagination.page - 1)"
      >
        <ChevronLeftIcon class="h-4 w-4" />
      </button>

      <!-- Pages -->
      <template v-for="page in pages" :key="page">
        <button
          v-if="page !== '...'"
          @click="$emit('change-page', page)"
          class="h-10 w-10 rounded-xl text-sm font-medium tabular-nums transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40"
          :class="
            page === pagination.page
              ? 'bg-blue-600 text-white shadow-sm'
              : 'border border-gray-200 bg-white text-gray-700 hover:bg-gray-50 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:hover:bg-slate-700'
          "
        >
          {{ page }}
        </button>

        <span v-else class="px-2 text-gray-400"> ... </span>
      </template>

      <!-- Next -->
      <button
        class="inline-flex h-10 items-center gap-2 rounded-xl border px-3 py-2 text-sm font-medium transition disabled:cursor-not-allowed disabled:opacity-50 border-gray-200 bg-white text-gray-700 hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:hover:bg-slate-700"
        :disabled="pagination.page === pagination.total_page"
        @click="$emit('change-page', pagination.page + 1)"
      >
        <ChevronRightIcon class="h-4 w-4" />
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { ChevronLeftIcon, ChevronRightIcon } from "@heroicons/vue/24/outline";

const props = defineProps({
  pagination: {
    type: Object,
    required: true,
  },
});

defineEmits(["change-page"]);

const startItem = computed(() => {
  if (!props.pagination.total) return 0;

  return (props.pagination.page - 1) * props.pagination.limit + 1;
});

const endItem = computed(() => {
  return Math.min(
    props.pagination.page * props.pagination.limit,
    props.pagination.total,
  );
});

const pages = computed(() => {
  const total = props.pagination.total_page;
  const current = props.pagination.page;

  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1);
  }

  const result = [];

  result.push(1);

  if (current > 3) {
    result.push("...");
  }

  const start = Math.max(2, current - 1);
  const end = Math.min(total - 1, current + 1);

  for (let i = start; i <= end; i++) {
    result.push(i);
  }

  if (current < total - 2) {
    result.push("...");
  }

  result.push(total);

  return result;
});
</script>
