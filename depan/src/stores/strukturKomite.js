import { defineStore } from "pinia";
import api from "../services/api";

export const useStrukturKomiteStore = defineStore("strukturkomite", {
  state: () => ({
    strukturkomite: null,
    loading: false,

    // cache
    loaded: false,

    // dashboard
    dashboardStrukturKomite: null,
  }),

  actions: {
    async fetchStrukturKomite(force = false) {
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/struktur-komite");

        this.strukturkomite = response.data.data;
        this.loaded = true;
      } catch (error) {
        console.error("Gagal mengambil struktur komite:", error);
      } finally {
        this.loading = false;
      }
    },

    async incrementHits() {
      try {
        await api.post("/public/struktur-komite/hits");
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

    async getDashboardStrukturKomite() {
      const response = await api.get("/struktur-komite");
      this.dashboardStrukturKomite = response.data.data;
      return response.data;
    },

    async updateDashboardStrukturKomite(formData) {
      const response = await api.put("/struktur-komite", formData);
      this.dashboardStrukturKomite = response.data.data;

      this.clearCache();

      return response.data;
    },
  },
});
