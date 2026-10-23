<template>
  <div>
    <!-- Label -->
    <label
      v-if="label"
      class="mb-2 block text-sm font-medium text-gray-700 dark:text-slate-300"
    >
      {{ label }}

      <span v-if="required" class="ml-1 text-red-500"> * </span>
    </label>

    <!-- Input -->
    <slot />

    <!-- Hint -->
    <p
      v-if="hint && !error"
      class="mt-2 text-xs text-gray-500 dark:text-slate-400"
    >
      {{ hint }}
    </p>

    <!-- Error -->
    <p
      v-if="error"
      class="mt-2 text-xs font-medium text-red-600 dark:text-red-400"
    >
      {{ error }}
    </p>
  </div>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    default: "",
  },

  required: {
    type: Boolean,
    default: false,
  },

  hint: {
    type: String,
    default: "",
  },

  error: {
    type: String,
    default: "",
  },
});
</script>
