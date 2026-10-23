<script setup>
import { FaceFrownIcon, HomeIcon } from "@heroicons/vue/24/outline";
</script>

<template>
  <div class="flex flex-col items-center justify-center text-center px-4 py-16">
    <div
      class="mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-red-100"
    >
      <FaceFrownIcon class="h-10 w-10 text-red-500" />
    </div>

    <h1 class="text-6xl font-bold text-gray-800">404</h1>
    <h2 class="mt-2 text-xl font-semibold text-gray-700">
      Halaman tidak ditemukan!
    </h2>
    <p class="mt-2 text-gray-500 max-w-md">
      Maaf, halaman yang Anda cari tidak tersedia atau belum dibuat.
    </p>

    <router-link
      to="/dashboard"
      class="mt-6 inline-flex items-center gap-2 rounded-lg bg-blue-600 px-5 py-2.5 text-white font-medium hover:bg-blue-700 transition"
    >
      <HomeIcon class="h-5 w-5" />
      Kembali
    </router-link>
  </div>
</template>
