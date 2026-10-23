<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="modelValue"
        class="fixed inset-0 z-9999 flex items-center justify-center p-4"
      >
        <!-- Overlay -->
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" />

        <!-- Dialog -->
        <Transition
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="opacity-0 scale-95 translate-y-2"
          enter-to-class="opacity-100 scale-100 translate-y-0"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="opacity-100 scale-100 translate-y-0"
          leave-to-class="opacity-0 scale-95 translate-y-2"
        >
          <div
            v-if="modelValue"
            class="relative w-full max-w-md overflow-hidden rounded-2xl bg-white shadow-2xl ring-1 ring-slate-900/5 dark:bg-slate-900 dark:ring-white/10"
          >
            <!-- Header -->
            <div
              class="flex items-center gap-3 border-b border-gray-100 px-6 py-5 dark:border-white/5"
            >
              <!-- Icon -->
              <div
                class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl"
                :class="iconBackgroundClass"
              >
                <component :is="icon" class="h-6 w-6" :class="iconColorClass" />
              </div>

              <div class="flex-1">
                <h2 class="text-lg font-semibold text-gray-800 dark:text-white">
                  {{ title }}
                </h2>

                <p
                  v-if="subtitle"
                  class="mt-1 text-sm text-gray-500 dark:text-slate-400"
                >
                  {{ subtitle }}
                </p>
              </div>
            </div>

            <!-- Body -->
            <div class="px-6 py-5">
              <p
                class="whitespace-pre-line text-sm leading-6 text-gray-600 dark:text-slate-300"
              >
                {{ message }}
              </p>
            </div>

            <!-- Footer -->
            <div
              class="flex flex-col-reverse gap-3 border-t border-gray-100 bg-gray-50 px-6 py-4 sm:flex-row dark:border-white/5 dark:bg-white/5"
            >
              <!-- Konfirmasi -->
              <button
                type="button"
                @click="confirm"
                :disabled="loading"
                class="inline-flex items-center justify-center rounded-xl px-4 py-2.5 text-sm font-medium text-white transition focus-visible:outline-none focus-visible:ring-2 disabled:cursor-not-allowed disabled:opacity-60"
                :class="[confirmButtonClass, confirmButtonRingClass]"
              >
                <ArrowPathIcon
                  v-if="loading"
                  class="mr-2 h-4 w-4 animate-spin"
                />

                {{ loading ? "Memproses..." : confirmText }}
              </button>
              <!-- Batal -->
              <button
                type="button"
                @click="cancel"
                :disabled="loading"
                class="inline-flex items-center justify-center rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200 dark:hover:bg-slate-800"
              >
                {{ cancelText }}
              </button>
            </div>
          </div>
        </Transition>
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
import { computed, onBeforeUnmount, onMounted } from "vue";

import {
  CheckCircleIcon,
  ExclamationTriangleIcon,
  NoSymbolIcon,
  TrashIcon,
  KeyIcon,
  ArrowPathIcon,
} from "@heroicons/vue/24/outline";

// =========================================
// Props
// =========================================

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },

  title: {
    type: String,
    default: "Konfirmasi",
  },

  subtitle: {
    type: String,
    default: "",
  },

  message: {
    type: String,
    default: "",
  },

  confirmText: {
    type: String,
    default: "Ya",
  },

  cancelText: {
    type: String,
    default: "Batal",
  },

  confirmColor: {
    type: String,
    default: "blue",
  },

  loading: {
    type: Boolean,
    default: false,
  },
});

// =========================================
// Emits
// =========================================

const emit = defineEmits(["update:modelValue", "confirm", "cancel"]);

// =========================================
// Close Dialog
// =========================================

function close() {
  emit("update:modelValue", false);
}

function cancel() {
  if (props.loading) return;

  close();
  emit("cancel");
}

function confirm() {
  if (props.loading) return;

  emit("confirm");
}

// =========================================
// ESC Key
// =========================================

function handleKeydown(event) {
  if (!props.modelValue) return;

  // Abaikan tombol ESC
  if (event.key === "Escape") {
    event.preventDefault();
  }
}

onMounted(() => {
  window.addEventListener("keydown", handleKeydown);
});

onBeforeUnmount(() => {
  window.removeEventListener("keydown", handleKeydown);
});

// =========================================
// Icon
// =========================================

const icon = computed(() => {
  switch (props.confirmColor) {
    case "green":
      return CheckCircleIcon;

    case "amber":
      return NoSymbolIcon;

    case "red":
      return TrashIcon;

    case "blue":
      return KeyIcon;

    default:
      return ExclamationTriangleIcon;
  }
});

// =========================================
// Icon Background
// =========================================

const iconBackgroundClass = computed(() => {
  switch (props.confirmColor) {
    case "green":
      return "bg-emerald-100 dark:bg-emerald-900/30";

    case "amber":
      return "bg-amber-100 dark:bg-amber-900/30";

    case "red":
      return "bg-red-100 dark:bg-red-900/30";

    case "blue":
      return "bg-blue-100 dark:bg-blue-900/30";

    default:
      return "bg-[#041983]/10 dark:bg-white/5";
  }
});

// =========================================
// Icon Color
// =========================================

const iconColorClass = computed(() => {
  switch (props.confirmColor) {
    case "green":
      return "text-emerald-600 dark:text-emerald-400";

    case "amber":
      return "text-amber-600 dark:text-amber-400";

    case "red":
      return "text-red-600 dark:text-red-400";

    case "blue":
      return "text-blue-600 dark:text-blue-400";

    default:
      return "text-[#041983] dark:text-indigo-300";
  }
});

// =========================================
// Confirm Button
// =========================================

const confirmButtonClass = computed(() => {
  switch (props.confirmColor) {
    case "green":
      return "bg-emerald-600 hover:bg-emerald-700";

    case "amber":
      return "bg-amber-500 hover:bg-amber-600";

    case "red":
      return "bg-red-600 hover:bg-red-700";

    case "blue":
      return "bg-blue-600 hover:bg-blue-700";

    default:
      return "bg-[#041983] hover:bg-[#041983]/90";
  }
});

// Ring focus mengikuti warna tombol konfirmasi, biar tetap terasa satu kesatuan
const confirmButtonRingClass = computed(() => {
  switch (props.confirmColor) {
    case "green":
      return "focus-visible:ring-emerald-500/40";

    case "amber":
      return "focus-visible:ring-amber-500/40";

    case "red":
      return "focus-visible:ring-red-500/40";

    case "blue":
      return "focus-visible:ring-blue-500/40";

    default:
      return "focus-visible:ring-[#041983]/40";
  }
});
</script>
