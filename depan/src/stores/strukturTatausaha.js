import { defineStore } from "pinia";
import api from "../services/api";

export const useStrukturTatausahaStore = defineStore("strukturtatausaha", {
  state: () => ({
    strukturtatausaha: null,
    loading: false,

    // cache
    loaded: false,

    // dashboard
    dashboardStrukturTatausaha: null,
  }),

  actions: {
    async fetchStrukturTatausaha(force = false) {
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/struktur-tatausaha");

        this.strukturtatausaha = response.data.data;
        this.loaded = true;
      } catch (error) {
        console.error("Gagal mengambil struktur tata usaha:", error);
      } finally {
        this.loading = false;
      }
    },

    async incrementHits() {
      try {
        await api.post("/public/struktur-tatausaha/hits");
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

    async getDashboardStrukturTatausaha() {
      const response = await api.get("/struktur-tatausaha");
      this.dashboardStrukturTatausaha = response.data.data;
      return response.data;
    },

    async updateDashboardStrukturTatausaha(formData) {
      const response = await api.put("/struktur-tatausaha", formData);
      this.dashboardStrukturTatausaha = response.data.data;

      this.clearCache();

      return response.data;
    },
  },
});
