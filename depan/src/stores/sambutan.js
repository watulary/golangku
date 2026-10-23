import { defineStore } from "pinia";
import api from "../services/api";

export const useSambutanStore = defineStore("sambutan", {
  state: () => ({
    // =====================================
    // Public
    // =====================================

    sambutan: null,
    loading: false,
    loaded: false,

    // =====================================
    // Dashboard
    // =====================================

    dashboardSambutan: null,
    dashboardLoading: false,
  }),

  actions: {
    // =====================================
    // Public
    // =====================================

    async fetchSambutan(force = false) {
      // Sudah pernah diambil
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/sambutan");

        this.sambutan = response.data.data;
        this.loaded = true;
      } catch (error) {
        console.error(error);
      } finally {
        this.loading = false;
      }
    },

    async incrementHits() {
      try {
        await api.post("/public/sambutan/hits");
      } catch (error) {
        console.error(error);
      }
    },

    clearCache() {
      this.loaded = false;
    },

    // =====================================
    // Dashboard
    // =====================================

    async getDashboardSambutan() {
      this.dashboardLoading = true;

      try {
        const response = await api.get("/sambutan");

        this.dashboardSambutan = response.data.data;

        return response.data;
      } finally {
        this.dashboardLoading = false;
      }
    },

    async updateDashboardSambutan(payload) {
      this.dashboardLoading = true;

      try {
        const response = await api.put("/sambutan", payload);

        if (response.data.data) {
          this.dashboardSambutan = response.data.data;
        }

        return response.data;
      } finally {
        this.dashboardLoading = false;
      }
    },
  },
});
