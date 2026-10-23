<template>
  <div
    class="pointer-events-none fixed right-10 top-17 z-9999 flex flex-col gap-3"
  >
    <TransitionGroup tag="div" name="toast">
      <BaseToast
        v-for="item in toastStore.toasts.value"
        :key="item.id"
        :toast="item"
        class="pointer-events-auto"
        @close="toastStore.remove"
      />
    </TransitionGroup>
  </div>
</template>

<script setup>
import { useToast } from "@/composables/useToast";
import BaseToast from "./BaseToast.vue";

const toastStore = useToast();
</script>
