<template>
  <div class="flex h-screen overflow-hidden bg-gray-100 dark:bg-[#0B0F1A]">
    <!-- Overlay Mobile -->
    <Transition
      enter-active-class="transition-opacity duration-300"
      leave-active-class="transition-opacity duration-200"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="sidebarOpen"
        class="fixed inset-0 z-40 bg-black/40 lg:hidden"
        @click="closeSidebar"
      />
    </Transition>

    <!-- Sidebar -->
    <div
      :class="[
        'fixed inset-y-0 left-0 z-50 transition-transform duration-300 lg:relative lg:translate-x-0',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full',
      ]"
    >
      <DashboardSidebar
        :collapsed="sidebarCollapsed"
        @toggle-collapse="toggleCollapse"
      />
    </div>

    <!-- Content -->
    <div class="flex min-w-0 flex-1 flex-col">
      <!-- Header -->
      <DashboardHeader @toggle-sidebar="toggleSidebar" @logout="logout" />

      <!-- Main -->
      <main
        class="flex-1 overflow-y-auto bg-gray-50 p-6 transition-colors duration-300 dark:bg-[#0B0F1A]"
      >
        <RouterView />
      </main>

      <!-- Footer -->
      <DashboardFooter />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from "vue";
import { RouterView } from "vue-router";
import { useTheme } from "@/composables/useTheme";
import { useSettingStore } from "@/stores/setting";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";

import DashboardSidebar from "@/components/dashboard/DashboardSidebar.vue";
import DashboardHeader from "@/components/dashboard/DashboardHeader.vue";
import DashboardFooter from "@/components/dashboard/DashboardFooter.vue";

const { initTheme } = useTheme();

onMounted(() => {
  initTheme();
});

const settingStore = useSettingStore();
const router = useRouter();
const authStore = useAuthStore();

onMounted(async () => {
  await Promise.all([
    settingStore.fetchSetting(),
    authStore.fetchMe(),
  ]);
});

/*
|--------------------------------------------------------------------------
| Mobile Breakpoint
|--------------------------------------------------------------------------
*/

const isDesktop = () => window.innerWidth >= 1024;

/*
|--------------------------------------------------------------------------
| Sidebar
|--------------------------------------------------------------------------
*/

const sidebarOpen = ref(isDesktop());
const sidebarCollapsed = ref(
  localStorage.getItem("sidebar-collapsed") === "true",
);

/*
|--------------------------------------------------------------------------
| Sidebar
|--------------------------------------------------------------------------
*/

function openSidebar() {
  sidebarOpen.value = true;
}

function closeSidebar() {
  if (!isDesktop()) {
    sidebarOpen.value = false;
  }
}

function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value;
}

function toggleCollapse() {
  sidebarCollapsed.value = !sidebarCollapsed.value;

  localStorage.setItem("sidebar-collapsed", sidebarCollapsed.value);
}

/*
|--------------------------------------------------------------------------
| Resize
|--------------------------------------------------------------------------
*/

function handleResize() {
  if (isDesktop()) {
    sidebarOpen.value = true;
  } else {
    sidebarOpen.value = false;
  }
}

/*
|--------------------------------------------------------------------------
| Logout
|--------------------------------------------------------------------------
*/

async function logout() {
  authStore.logout();

  await router.replace("/login");
}

/*
|--------------------------------------------------------------------------
| Lifecycle
|--------------------------------------------------------------------------
*/

onMounted(() => {
  window.addEventListener("resize", handleResize);
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", handleResize);
});
</script>