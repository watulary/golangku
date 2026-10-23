import { defineStore } from "pinia";
import api from "../services/api";

export const useStrukturAkademikStore = defineStore("strukturakademik", {
  state: () => ({
    strukturakademik: null,
    loading: false,

    // cache
    loaded: false,

    // dashboard
    dashboardStrukturAkademik: null,
  }),

  actions: {
    async fetchStrukturAkademik(force = false) {
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/struktur-akademik");

        this.strukturakademik = response.data.data;
        this.loaded = true;
      } catch (error) {
        console.error("Gagal mengambil struktur akademik:", error);
      } finally {
        this.loading = false;
      }
    },

    async incrementHits() {
      try {
        await api.post("/public/struktur-akademik/hits");
      } catch (error) {
        console.error("Gagal update hits:", error);
      }
    },

    clearCache() {
      this.loaded = false;
    },

    // =====================================
    // Dashboard
    // =====================================

    async getDashboardStrukturAkademik() {
      const response = await api.get("/struktur-akademik");
      this.dashboardStrukturAkademik = response.data.data;
      return response.data;
    },

    async updateDashboardStrukturAkademik(formData) {
      const response = await api.put("/struktur-akademik", formData);
      this.dashboardStrukturAkademik = response.data.data;

      this.clearCache();

      return response.data;
    },
  },
});
