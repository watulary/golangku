<template>
  <div
    class="flex items-start gap-3 border-b border-gray-100 px-6 py-5 dark:border-white/5"
  >
    <!-- Icon -->
    <div
      v-if="$slots.default"
      class="mt-0.5 flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/[0.08] text-[#041983] dark:bg-white/5 dark:text-indigo-300"
    >
      <slot />
    </div>

    <!-- Title -->
    <div class="min-w-0 flex-1">
      <h2 class="text-lg font-semibold text-gray-700 dark:text-slate-200">
        {{ title }}
      </h2>

      <p
        v-if="description"
        class="mt-1 text-sm text-gray-500 dark:text-slate-400"
      >
        {{ description }}
      </p>
    </div>

    <!-- Action -->
    <div v-if="$slots.action">
      <slot name="action" />
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true,
  },

  description: {
    type: String,
    default: "",
  },
});
</script>
