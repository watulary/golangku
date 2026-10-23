<template>
  <aside
    :class="[
      sidebarWidth,
      'flex h-screen flex-col transition-all duration-300',
      'bg-white dark:bg-slate-900',
      'border-r border-gray-100 dark:border-white/5',
    ]"
  >
    <!-- ==========================================================
    Header
    =========================================================== -->

    <div
      class="flex h-16 items-center justify-between dark:border-slate-800 px-4"
    >
      <!-- Logo -->

      <RouterLink to="/dashboard" class="flex min-w-0 items-center gap-3">
        <div
          class="relative flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-xl bg-white dark:bg-[#11172A] ring-1 ring-slate-900/5 dark:ring-white/5 shadow-xs"
        >
          <img
            v-if="settingStore.setting?.logo_url"
            :src="settingStore.setting.logo_url"
            :alt="settingStore.setting.nama_sekolah"
            class="h-8 w-8 object-contain"
          />

          <ShieldCheckIcon
            v-else
            class="h-5 w-5 text-[#041983] dark:text-indigo-300"
          />
        </div>

        <div v-if="!props.collapsed" class="min-w-0">
          <h1
            class="truncate text-[15px] font-semibold text-gray-700 dark:text-slate-300"
          >
            {{ settingStore.setting?.nama_sekolah }}
          </h1>
        </div>
      </RouterLink>

      <!-- Collapse -->

      <button
        v-if="!props.collapsed"
        @click="$emit('toggle-collapse')"
        class="rounded-lg p-2 text-gray-500 transition-all duration-200 hover:bg-[#041983]/10 hover:text-[#041983] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:text-slate-400 dark:hover:bg-white/5 dark:hover:text-indigo-300"
      >
        <ChevronLeftIcon class="h-4 w-4" />
      </button>

      <button
        v-else
        @click="$emit('toggle-collapse')"
        class="rounded-lg p-2 text-gray-500 transition-all duration-200 hover:bg-[#041983]/10 hover:text-[#041983] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:text-slate-400 dark:hover:bg-white/5 dark:hover:text-indigo-300"
      >
        <ChevronRightIcon class="h-4 w-4" />
      </button>
    </div>

    <!-- ==========================================================
    Menu
    =========================================================== -->
    <div class="relative flex-1 overflow-hidden">
      <!-- Scroll Area -->
      <div
        class="sidebar-scroll h-full overflow-y-auto overflow-x-visible px-3 py-5"
      >
        <nav class="space-y-0.5">
          <template v-for="menu in menus" :key="menu.title">
            <!-- ======================================
          Menu tanpa submenu
          ======================================= -->

            <div v-if="!menu.children" class="relative group">
              <RouterLink
                :to="menu.to"
                @click="closeAllMenus"
                @mouseenter="props.collapsed && showTooltip($event, menu.title)"
                @mouseleave="hideTooltip"
                class="relative flex h-10 items-center rounded-xl transition-all duration-200"
                :class="[
                  props.collapsed ? 'justify-center px-0' : 'px-3',
                  isActive(menu)
                    ? 'bg-[#041983]/10 text-[#041983] dark:bg-white/5 dark:text-indigo-300'
                    : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900 dark:text-slate-300 dark:hover:bg-white/5 dark:hover:text-white',
                ]"
              >
                <span
                  v-if="isActive(menu)"
                  class="absolute inset-y-1 left-0 w-1 rounded-r-full bg-linear-to-b from-[#041983] to-[#f2c808] dark:from-indigo-400 dark:to-[#f2c808]"
                />

                <component :is="menu.icon" class="h-5 w-5 shrink-0" />

                <span
                  v-if="!props.collapsed"
                  class="ml-3 text-[14px] font-medium"
                >
                  {{ menu.title }}
                </span>
              </RouterLink>
            </div>

            <!-- ======================================
          Parent
          ======================================= -->

            <div v-else class="relative group rounded-xl">
              <button
                @click="handleParentClick(menu)"
                @mouseenter="props.collapsed && showTooltip($event, menu.title)"
                @mouseleave="hideTooltip"
                class="relative flex h-10 w-full items-center rounded-xl transition-all duration-200"
                :class="[
                  props.collapsed ? 'justify-center px-0' : 'px-3',
                  isChildActive(menu.children)
                    ? 'bg-[#041983]/10 text-[#041983] dark:bg-white/5 dark:text-indigo-300'
                    : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900 dark:text-slate-300 dark:hover:bg-white/5 dark:hover:text-white',
                ]"
              >
                <span
                  v-if="isChildActive(menu.children)"
                  class="absolute inset-y-1 left-0 w-1 rounded-r-full bg-linear-to-b from-[#041983] to-[#f2c808] dark:from-indigo-400 dark:to-[#f2c808]"
                />

                <component :is="menu.icon" class="h-5 w-5 shrink-0" />

                <span
                  v-if="!props.collapsed"
                  class="ml-3 flex-1 text-left text-[14px] font-medium"
                >
                  {{ menu.title }}
                </span>

                <ChevronDownIcon
                  v-if="!props.collapsed"
                  class="h-4 w-4 transition-transform duration-200"
                  :class="{
                    'rotate-180': isOpen(menu.title),
                  }"
                />
              </button>

              <!-- Submenu -->
              <Transition
                enter-active-class="transition-all duration-200 ease-out"
                leave-active-class="transition-all duration-150 ease-in"
                enter-from-class="opacity-0 max-h-0"
                enter-to-class="opacity-100 max-h-96"
                leave-from-class="opacity-100 max-h-96"
                leave-to-class="opacity-0 max-h-0"
              >
                <div
                  v-if="!props.collapsed && isOpen(menu.title)"
                  class="relative mt-1 overflow-hidden"
                >
                  <!-- Vertical Line -->

                  <div
                    class="absolute left-5.25 top-0 bottom-0 w-px bg-linear-to-b from-[#041983]/40 dark:from-indigo-400/40 to-transparent"
                  />

                  <div class="space-y-0 pl-8">
                    <RouterLink
                      v-for="child in menu.children"
                      :key="child.to"
                      :to="child.to"
                      class="group flex h-9 items-center rounded-lg px-3 text-[13px] transition-all duration-200 hover:bg-gray-100 dark:hover:bg-white/5"
                      :class="
                        isChildLinkActive(child, menu.children)
                          ? 'font-semibold text-gray-700 dark:text-slate-200'
                          : 'font-medium text-gray-600 hover:text-gray-900 dark:text-slate-400 dark:hover:text-white'
                      "
                    >
                      <!-- Dot -->

                      <span
                        class="mr-3 h-2 w-2 rounded-full transition-all duration-200"
                        :class="
                          isChildLinkActive(child, menu.children)
                            ? 'bg-[#041983] dark:bg-indigo-400'
                            : 'bg-gray-400 dark:bg-slate-600'
                        "
                      />

                      <span class="truncate">
                        {{ child.title }}
                      </span>
                    </RouterLink>
                  </div>
                </div>
              </Transition>
            </div>
          </template>
        </nav>
      </div>
    </div>
    <DashboardTooltip
      :visible="tooltip.visible"
      :text="tooltip.text"
      :x="tooltip.x"
      :y="tooltip.y"
    />
  </aside>
</template>

<style scoped>
/* ==========================================================
Scrollbar
========================================================== */
.sidebar-scroll {
  height: 100%;
}

.sidebar-scroll::-webkit-scrollbar {
  width: 6px;
}

.sidebar-scroll::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar-scroll::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 9999px;
}

.dark .sidebar-scroll::-webkit-scrollbar-thumb {
  background: rgba(99, 102, 241, 0.25);
}

.sidebar-scroll::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

.dark .sidebar-scroll::-webkit-scrollbar-thumb:hover {
  background: rgba(99, 102, 241, 0.4);
}
</style>

<script setup>
import { ref, computed, watch } from "vue";
import { RouterLink, useRoute } from "vue-router";
import { useSettingStore } from "@/stores/setting";
import DashboardTooltip from "@/components/dashboard/DashboardTooltip.vue";
import { useDashboardMenu } from "@/composables/useDashboardMenu";

import {
  ShieldCheckIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ChevronDownIcon,
} from "@heroicons/vue/24/outline";

/* --------------------------------------------------------------------------
 * Store
 * -------------------------------------------------------------------------- */

const settingStore = useSettingStore();

/* --------------------------------------------------------------------------
 * Router
 * -------------------------------------------------------------------------- */

const route = useRoute();

/* --------------------------------------------------------------------------
 * Props
 * -------------------------------------------------------------------------- */

const props = defineProps({
  collapsed: {
    type: Boolean,
    default: false,
  },
});

/* --------------------------------------------------------------------------
 * Emits
 * -------------------------------------------------------------------------- */

const emit = defineEmits(["toggle-collapse"]);

/* --------------------------------------------------------------------------
 * Sidebar Width
 * -------------------------------------------------------------------------- */

const sidebarWidth = computed(() => (props.collapsed ? "w-20" : "w-64"));

/* --------------------------------------------------------------------------
 * Open Menu
 * -------------------------------------------------------------------------- */

const openMenus = ref({});

/* --------------------------------------------------------------------------
 * Menu
 * -------------------------------------------------------------------------- */

const { menus } = useDashboardMenu();

/* --------------------------------------------------------------------------
 * Active Menu
 * -------------------------------------------------------------------------- */

function isActive(menu) {
  if (!menu.to) return false;

  // Beranda hanya aktif saat benar-benar di /dashboard
  if (menu.to === "/dashboard") {
    return route.path === "/dashboard";
  }

  return route.path === menu.to || route.path.startsWith(menu.to + "/");
}

function isChildActive(children) {
  return children.some(
    (child) => route.path === child.to || route.path.startsWith(child.to + "/"),
  );
}

/* --------------------------------------------------------------------------
 * Toggle
 * -------------------------------------------------------------------------- */

function toggleMenu(title) {
  const isCurrentlyOpen = !!openMenus.value[title];

  // Tutup semua
  openMenus.value = {};

  // Jika sebelumnya tertutup, buka yang dipilih
  if (!isCurrentlyOpen) {
    openMenus.value[title] = true;
  }
}

function isOpen(title) {
  return !!openMenus.value[title];
}

function closeAllMenus() {
  openMenus.value = {};
}

/* --------------------------------------------------------------------------
 * Open Active Parent
 * -------------------------------------------------------------------------- */

function openActiveParent() {
  // Tutup semua submenu terlebih dahulu
  openMenus.value = {};

  // Buka hanya parent yang sedang aktif
  menus.value.forEach((menu) => {
    if (menu.children && isChildActive(menu.children)) {
      openMenus.value[menu.title] = true;
    }
  });
}

watch(() => route.path, openActiveParent, {
  immediate: true,
});

function handleParentClick(menu) {
  if (props.collapsed) {
    pendingOpenMenu.value = menu.title;

    emit("toggle-collapse");

    return;
  }

  toggleMenu(menu.title);
}

const pendingOpenMenu = ref(null);

watch(
  () => props.collapsed,
  (collapsed) => {
    if (collapsed) {
      openMenus.value = {};
      return;
    }

    openActiveParent();

    if (pendingOpenMenu.value) {
      openMenus.value[pendingOpenMenu.value] = true;
      pendingOpenMenu.value = null;
    }
  },
);

const tooltip = ref({
  visible: false,
  text: "",
  x: 0,
  y: 0,
});

function showTooltip(event, text) {
  tooltip.value.visible = true;
  tooltip.value.text = text;

  const rect = event.currentTarget.getBoundingClientRect();

  tooltip.value.x = rect.right + 12;
  tooltip.value.y = rect.top + rect.height / 2;
}

function hideTooltip() {
  tooltip.value.visible = false;
}

function isChildLinkActive(child, siblings) {
  const childTo = child.to;
  if (route.path === childTo) return true;
  if (!route.path.startsWith(childTo + "/")) return false;

  // Kalau ada sibling lain (mis. "Mutasi"/"Arsip") yang path-nya lebih spesifik
  // dan juga cocok dengan route saat ini, sibling itu yang menang — child ini
  // (mis. "Aktif") tidak dianggap aktif meski secara prefix string cocok.
  const hasMoreSpecificSibling = siblings.some((sibling) => {
    if (sibling.to === childTo) return false;
    if (!sibling.to.startsWith(childTo)) return false;
    return route.path === sibling.to || route.path.startsWith(sibling.to + "/");
  });

  return !hasMoreSpecificSibling;
}
</script>
