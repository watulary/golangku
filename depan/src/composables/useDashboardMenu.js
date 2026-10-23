import { computed } from "vue";

import { useAuthStore } from "../stores/auth";
import { dashboardMenus } from "../config/dashboard-menu";

export function useDashboardMenu() {
  const authStore = useAuthStore();

  const menus = computed(() => {
    const currentRole = authStore.role;

    if (!currentRole) return [];

    // Force Change Password
    if (authStore.mustChangePassword) {
      return dashboardMenus.filter((menu) => menu.forceChangePassword);
    }

    return (
      dashboardMenus
        .filter((menu) => menu.roles.includes(currentRole))
        .map((menu) => {
          // Menu tanpa submenu
          if (!menu.children) {
            return menu;
          }

          // Filter submenu berdasarkan role
          const children = menu.children.filter((child) =>
            child.roles.includes(currentRole),
          );

          return {
            ...menu,
            children,
          };
        })
        // Hilangkan parent jika semua submenu hilang
        .filter((menu) => {
          if (!menu.children) return true;

          return menu.children.length > 0;
        })
    );
  });

  return {
    menus,
  };
}
