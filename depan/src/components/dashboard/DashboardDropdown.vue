<template>
  <div class="relative" ref="dropdown">
    <!-- Trigger -->
    <div @click="toggle">
      <slot name="trigger" />
    </div>

    <!-- Desktop -->
    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div
        v-if="open"
        class="absolute left-0 top-11 z-50 hidden w-60 overflow-hidden rounded-2xl bg-white shadow-xl ring-1 ring-slate-900/5 md:block dark:bg-slate-900 dark:ring-white/10"
      >
        <slot />
      </div>
    </Transition>

    <!-- Mobile -->
    <Transition
      enter-active-class="transition duration-200"
      enter-from-class="translate-y-full"
      enter-to-class="translate-y-0"
      leave-active-class="transition duration-150"
      leave-from-class="translate-y-0"
      leave-to-class="translate-y-full"
    >
      <div
        v-if="open"
        class="fixed inset-0 z-999 flex items-end bg-black/40 md:hidden"
      >
        <div
          class="w-full rounded-t-2xl bg-white p-4 ring-1 ring-slate-900/5 dark:bg-slate-900 dark:ring-white/10"
        >
          <div
            class="mx-auto mb-4 h-1.5 w-12 rounded-full bg-gray-300 dark:bg-slate-700"
          />

          <slot />

          <button
            type="button"
            class="mt-3 w-full rounded-xl border border-gray-200 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:border-slate-700 dark:text-slate-200 dark:hover:bg-slate-800"
            @click="close"
          >
            Tutup
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
@media (prefers-reduced-motion: reduce) {
  .transition {
    transition: none !important;
  }
}
</style>

<script setup>
import { ref, onMounted, onBeforeUnmount } from "vue";

const open = ref(false);
const dropdown = ref(null);

function toggle() {
  open.value = !open.value;
}

function close() {
  open.value = false;
}

function handleClickOutside(event) {
  if (dropdown.value && !dropdown.value.contains(event.target)) {
    close();
  }
}

onMounted(() => {
  document.addEventListener("click", handleClickOutside);
});

onBeforeUnmount(() => {
  document.removeEventListener("click", handleClickOutside);
});

defineExpose({
  close,
});
</script>
