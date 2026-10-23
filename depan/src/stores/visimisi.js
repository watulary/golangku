import { defineStore } from "pinia";
import api from "../services/api";

export const useVisiMisiStore = defineStore("visimisi", {
  state: () => ({
    visimisi: null,
    loading: false,

    // cache
    loaded: false,

    // dashboard
    dashboardVisiMisi: null,
  }),

  actions: {
    async fetchVisiMisi(force = false) {
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/visimisi");

        this.visimisi = response.data.data;
        this.loaded = true;
      } catch (error) {
        console.error("Gagal mengambil visi misi:", error);
      } finally {
        this.loading = false;
      }
    },

    async incrementHits() {
      try {
        await api.post("/public/visimisi/hits");
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

    async getDashboardVisiMisi() {
      const response = await api.get("/visimisi");
      this.dashboardVisiMisi = response.data.data;
      return response.data;
    },

    async updateDashboardVisiMisi(payload) {
      const response = await api.put("/visimisi", payload);
      this.dashboardVisiMisi = response.data.data;

      this.clearCache();

      return response.data;
    },
  },
});
