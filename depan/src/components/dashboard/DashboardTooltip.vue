<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="visible"
        class="pointer-events-none fixed z-[9999] -translate-y-1/2"
        :style="style"
      >
        <div
          class="rounded-xl bg-white px-3 py-1.5 text-sm font-medium text-gray-700 shadow-xl ring-1 ring-slate-900/5 dark:bg-[#11172A] dark:text-white dark:shadow-xs dark:ring-white/10"
        >
          {{ text }}
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
@media (prefers-reduced-motion: reduce) {
  .transition {
    transition: none !important;
  }
}
</style>

<script setup>
import { computed } from "vue";

const props = defineProps({
  visible: Boolean,

  text: {
    type: String,
    default: "",
  },

  x: {
    type: Number,
    default: 0,
  },

  y: {
    type: Number,
    default: 0,
  },
});

const style = computed(() => ({
  left: `${props.x}px`,
  top: `${props.y}px`,
}));
</script>
