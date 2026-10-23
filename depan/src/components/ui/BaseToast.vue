<template>
  <div
    class="relative overflow-hidden rounded-2xl border-l-4 bg-white shadow-lg ring-1 ring-slate-900/5 transition-colors dark:bg-slate-900 dark:ring-white/10"
    :class="containerClass"
  >
    <div class="flex items-start gap-3 px-4 py-3">
      <!-- Icon -->

      <component
        :is="icon"
        class="mt-0.5 h-6 w-6 shrink-0"
        :class="iconColor"
      />

      <!-- Content -->

      <div class="min-w-0 flex-1">
        <p class="text-sm font-medium text-gray-800 dark:text-slate-100">
          {{ toast.message }}
        </p>
      </div>

      <!-- Close -->

      <button
        type="button"
        @click="$emit('close', toast.id)"
        class="rounded-sm text-gray-400 transition hover:text-gray-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gray-400/40 dark:text-slate-500 dark:hover:text-slate-300 dark:focus-visible:ring-slate-500/40"
      >
        <XMarkIcon class="h-5 w-5" />
      </button>
    </div>

    <!-- Progress -->

    <div
      class="toast-progress-bar absolute bottom-0 left-0 h-1"
      :class="progressColor"
      :style="{
        animationDuration: `${toast.duration}ms`,
      }"
    />
  </div>
</template>

<script setup>
import { computed } from "vue";
import {
  CheckCircleIcon,
  XCircleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  XMarkIcon,
} from "@heroicons/vue/24/solid";

const props = defineProps({
  toast: {
    type: Object,
    required: true,
  },
});

defineEmits(["close"]);

const config = computed(() => {
  switch (props.toast.type) {
    case "success":
      return {
        icon: CheckCircleIcon,
        iconColor: "text-emerald-500",
        progressColor: "bg-emerald-500",
        containerClass: "border-l-emerald-500",
      };

    case "error":
      return {
        icon: XCircleIcon,
        iconColor: "text-red-500",
        progressColor: "bg-red-500",
        containerClass: "border-l-red-500",
      };

    case "warning":
      return {
        icon: ExclamationTriangleIcon,
        iconColor: "text-amber-500",
        progressColor: "bg-amber-500",
        containerClass: "border-l-amber-500",
      };

    default:
      return {
        icon: InformationCircleIcon,
        iconColor: "text-blue-500",
        progressColor: "bg-blue-500",
        containerClass: "border-l-blue-500",
      };
  }
});

const icon = computed(() => config.value.icon);

const iconColor = computed(() => config.value.iconColor);

const progressColor = computed(() => config.value.progressColor);

const containerClass = computed(() => config.value.containerClass);
</script>

<style scoped>
@keyframes toast-progress {
  from {
    width: 100%;
  }

  to {
    width: 0%;
  }
}

.toast-progress-bar {
  animation-name: toast-progress;
  animation-timing-function: linear;
  animation-fill-mode: forwards;
}

@media (prefers-reduced-motion: reduce) {
  .toast-progress-bar {
    animation: none;
  }
}
</style>
