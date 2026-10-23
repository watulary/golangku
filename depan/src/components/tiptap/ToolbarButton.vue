<template>
  <button
    type="button"
    :disabled="disabled"
    :title="title"
    @click="$emit('click')"
    class="inline-flex h-9 w-9 items-center justify-center rounded-lg border text-sm transition-all duration-150"
    :class="[
      active
        ? 'border-blue-600 bg-blue-600 text-white'
        : 'border-gray-200 bg-white text-gray-600 hover:bg-gray-100 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700',
      disabled ? 'cursor-not-allowed opacity-40' : 'cursor-pointer',
    ]"
  >
    <slot />
  </button>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: "",
  },
  active: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
});

defineEmits(["click"]);
</script>
