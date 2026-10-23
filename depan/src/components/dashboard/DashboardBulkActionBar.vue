<template>
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0 -translate-y-2"
    enter-to-class="opacity-100 translate-y-0"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100 translate-y-0"
    leave-to-class="opacity-0 -translate-y-2"
  >
    <div
      v-if="selectedCount > 0"
      class="mb-5 overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
    >
      <!-- Header -->
      <div
        class="flex flex-col gap-3 border-b border-gray-100 px-5 py-4 sm:flex-row sm:items-center sm:justify-between dark:border-white/5"
      >
        <div class="flex items-center gap-2">
          <CheckCircleIcon
            class="h-5 w-5 text-[#041983] dark:text-indigo-300"
          />
          <span class="text-sm font-semibold text-slate-800 dark:text-white">
            {{ selectedCount }} Dipilih
          </span>
        </div>
        <button
          type="button"
          @click="$emit('clear')"
          class="inline-flex items-center gap-2 rounded-xl border border-gray-200 bg-white px-3 py-2 text-sm font-medium text-slate-600 transition-all duration-200 hover:bg-gray-50 hover:text-slate-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:bg-slate-700 dark:hover:text-slate-100"
        >
          <XMarkIcon class="h-4 w-4" />
          Batalkan
        </button>
      </div>
      <!-- Actions -->
      <div
        class="grid grid-cols-2 gap-3 p-4 lg:flex lg:flex-wrap lg:items-center"
      >
        <slot name="actions">
          <!-- Default: 4 aksi pengguna (tidak berubah, dipakai kalau slot tidak di-override) -->
          <button
            type="button"
            @click="$emit('activate')"
            class="inline-flex items-center justify-center gap-2 rounded-xl bg-emerald-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-emerald-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40"
          >
            <CheckCircleIcon class="h-5 w-5" />
            Aktifkan
          </button>
          <button
            type="button"
            @click="$emit('deactivate')"
            class="inline-flex items-center justify-center gap-2 rounded-xl bg-amber-500 px-4 py-2 text-sm font-medium text-white transition hover:bg-amber-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500/40"
          >
            <NoSymbolIcon class="h-5 w-5" />
            Nonaktifkan
          </button>
          <button
            type="button"
            @click="$emit('reset-password')"
            class="inline-flex items-center justify-center gap-2 rounded-xl bg-blue-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40"
          >
            <KeyIcon class="h-5 w-5" />
            Reset Password
          </button>
          <button
            type="button"
            @click="$emit('delete')"
            class="inline-flex items-center justify-center gap-2 rounded-xl bg-red-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-red-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40"
          >
            <TrashIcon class="h-5 w-5" />
            Hapus
          </button>
        </slot>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
@media (prefers-reduced-motion: reduce) {
  .transition,
  .transition-all {
    transition: none !important;
  }
}
</style>

<script setup>
import {
  CheckCircleIcon,
  XMarkIcon,
  NoSymbolIcon,
  KeyIcon,
  TrashIcon,
} from "@heroicons/vue/24/outline";

defineProps({
  selectedCount: {
    type: Number,
    default: 0,
  },
});
defineEmits(["activate", "deactivate", "reset-password", "delete", "clear"]);
</script>
