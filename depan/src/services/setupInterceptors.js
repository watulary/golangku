import api from "./api";

export function setupInterceptors(authStore, router) {
  api.interceptors.response.use(
    (response) => response,

    async (error) => {
      const status = error.response?.status;
      const code = error.response?.data?.code;
      const message = error.response?.data?.message;

      // ===============================
      // Akun dinonaktifkan
      // ===============================
      if (status === 401 && code === "ACCOUNT_DISABLED") {
        authStore.logout();

        alert(message || "Akun Anda telah dinonaktifkan.");

        if (router.currentRoute.value.path !== "/login") {
          router.replace("/login");
        }

        return Promise.reject(error);
      }

      // ===============================
      // Token tidak valid / kedaluwarsa
      // ===============================
      if (status === 401) {
        authStore.logout();

        if (router.currentRoute.value.path !== "/login") {
          router.replace("/login");
        }
      }

      return Promise.reject(error);
    },
  );
}
