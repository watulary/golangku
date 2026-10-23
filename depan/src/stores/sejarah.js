import { defineStore } from "pinia";
import api from "../services/api";

export const useSejarahStore = defineStore("sejarah", {
  state: () => ({
    sejarah: null,
    loading: false,

    // cache
    loaded: false,

    // dashboard
    dashboardSejarah: null,
  }),

  actions: {
    async fetchSejarah(force = false) {
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/sejarah");

        this.sejarah = response.data.data;
        this.loaded = true;
      } catch (error) {
        console.error("Gagal mengambil sejarah:", error);
      } finally {
        this.loading = false;
      }
    },
    async incrementHits() {
      try {
        await api.post("/public/sejarah/hits");
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

    async getDashboardSejarah() {
      const response = await api.get("/sejarah");
      this.dashboardSejarah = response.data.data;
      return response.data;
    },

    async updateDashboardSejarah(payload) {
      const response = await api.put("/sejarah", payload);
      this.dashboardSejarah = response.data.data;

      // sinkronkan cache publik juga, biar tidak stale
      this.clearCache();

      return response.data;
    },
  },
});
