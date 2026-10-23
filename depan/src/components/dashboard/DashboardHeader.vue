<template>
  <header
    class="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-gray-100 bg-white/80 px-4 backdrop-blur-md transition-colors duration-300 dark:border-white/5 dark:bg-slate-900 lg:px-6"
  >
    <!-- ============================= -->
    <!-- Left -->
    <!-- ============================= -->
    <div class="flex min-w-0 items-center gap-3">
      <!-- Mobile Sidebar Button -->
      <button
        @click="$emit('toggle-sidebar')"
        class="shrink-0 rounded-lg p-2 text-gray-600 transition hover:bg-[#041983]/10 hover:text-[#041983] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:text-slate-300 dark:hover:bg-white/5 dark:hover:text-indigo-300 lg:hidden"
      >
        <Bars3Icon class="h-6 w-6" />
      </button>

      <!-- Date Time -->
      <div
        class="flex items-center gap-2 rounded-xl bg-gray-50 px-3 py-1.5 text-sm font-medium text-gray-600 ring-1 ring-slate-900/5 dark:bg-white/5 dark:text-slate-300 dark:ring-white/5"
      >
        <CalendarDaysIcon
          class="h-5 w-5 shrink-0 text-gray-400 dark:text-slate-500"
        />
        <!-- Mobile: tanggal singkat saja -->
        <span class="sm:hidden">{{ currentDateShort }}</span>
        <!-- Tablet & Desktop -->
        <span class="hidden sm:inline">{{ currentDateTime }}</span>
      </div>
    </div>

    <!-- ============================= -->
    <!-- Right -->
    <!-- ============================= -->
    <div class="flex shrink-0 items-center gap-2">
      <!-- Theme -->
      <button
        @click="toggleTheme"
        title="Ganti tema"
        class="shrink-0 rounded-xl p-2 text-gray-600 transition hover:bg-[#041983]/10 hover:text-[#041983] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:text-slate-300 dark:hover:bg-white/5 dark:hover:text-indigo-300"
      >
        <MoonIcon v-if="theme === 'light'" class="h-5 w-5" />
        <SunIcon v-else class="h-5 w-5" />
      </button>

      <!-- Divider -->
      <div class="hidden h-8 w-px bg-gray-100 dark:bg-white/10 sm:block"></div>

      <!-- User -->
      <div
        class="flex shrink-0 items-center gap-2 rounded-xl px-1 py-2 sm:gap-3 sm:px-2"
      >
        <!-- Avatar + nama: klik untuk membuka Profil Akun -->
        <RouterLink
          to="/dashboard/profil"
          title="Profil Akun"
          class="flex shrink-0 items-center gap-2 rounded-xl p-1 transition hover:bg-[#041983]/10 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:hover:bg-white/5 sm:gap-3"
        >
          <!-- Avatar -->
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-full bg-[#041983] text-sm font-semibold text-white"
          >
            <img
              v-if="user.foto_url"
              :src="user.foto_url"
              class="h-full w-full object-cover"
            />
            <span v-else>
              {{ initials }}
            </span>
          </div>
          <!-- User Info -->
          <div class="hidden text-left md:block">
            <p class="text-sm font-semibold text-gray-700 dark:text-slate-300">
              {{ user.nama_pengguna }}
            </p>
            <p class="text-xs text-gray-500 dark:text-slate-400">
              {{ authStore.displayRole }}
            </p>
          </div>
        </RouterLink>
        <!-- Logout -->
        <button
          @click="logout"
          title="Keluar"
          class="shrink-0 rounded-lg p-2 text-gray-500 transition hover:bg-red-50 hover:text-red-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40 dark:text-slate-400 dark:hover:bg-red-500/10 dark:hover:text-red-400"
        >
          <ArrowRightOnRectangleIcon class="h-5 w-5" />
        </button>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed, ref, onMounted, onBeforeUnmount } from "vue";
import { RouterLink } from "vue-router";
import { useAuthStore } from "@/stores/auth";

import {
  Bars3Icon,
  MoonIcon,
  SunIcon,
  ArrowRightOnRectangleIcon,
  CalendarDaysIcon,
} from "@heroicons/vue/24/outline";

import { useTheme } from "@/composables/useTheme";

/*
|--------------------------------------------------------------------------
| Props & Emits
|--------------------------------------------------------------------------
*/

defineProps({
  sidebarOpen: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["toggle-sidebar", "logout"]);

/*
|--------------------------------------------------------------------------
| Theme
|--------------------------------------------------------------------------
*/

const { theme, toggleTheme } = useTheme();

/*
|--------------------------------------------------------------------------
| Logout
|--------------------------------------------------------------------------
*/

function logout() {
  emit("logout");
}

/*
|--------------------------------------------------------------------------
| Date Time
|--------------------------------------------------------------------------
*/

const now = ref(new Date());

let timer = null;

onMounted(() => {
  timer = setInterval(() => {
    now.value = new Date();
  }, 60000);
});

onBeforeUnmount(() => {
  clearInterval(timer);
});

const currentDateTime = computed(() => {
  const date = new Intl.DateTimeFormat("id-ID", {
    weekday: "long",
    day: "numeric",
    month: "long",
    year: "numeric",
  }).format(now.value);

  const time = new Intl.DateTimeFormat("id-ID", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  }).format(now.value);

  return `${date} • ${time} WIB`;
});

const currentDateShort = computed(() => {
  const day = String(now.value.getDate()).padStart(2, "0");
  const month = String(now.value.getMonth() + 1).padStart(2, "0");
  const year = now.value.getFullYear();
  return `${day}-${month}-${year}`;
});

/*
|--------------------------------------------------------------------------
| User
|--------------------------------------------------------------------------
*/

const authStore = useAuthStore();

const user = computed(() => authStore.user ?? {});

const initials = computed(() => {
  const nama = user.value?.nama_pengguna ?? "";

  return nama
    .split(" ")
    .map((item) => item[0])
    .join("")
    .substring(0, 2)
    .toUpperCase();
});
</script>
