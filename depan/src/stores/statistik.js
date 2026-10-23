import { defineStore } from "pinia";
import api from "../services/api";

export const useStatistikStore = defineStore("statistik", {
  state: () => ({
    statistik: null,
    loading: false,

    rekapRombel: [],
    rekapRombelLoading: false,

    // cache
    loaded: false,
    dashboardLoaded: false,
    rekapRombelLoaded: false,
  }),

  actions: {
    async fetchStatistik(force = false) {
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/statistik");

        this.statistik = response.data.data;
        this.loaded = true;
      } catch (error) {
        console.error("Gagal mengambil statistik:", error);
      } finally {
        this.loading = false;
      }
    },

    // Versi dashboard (beda dari fetchStatistik di atas yg dipakai halaman
    // publik) — endpoint /statistik butuh login & difilter per-pengguna
    // (Super Admin lihat semua). Cache terpisah (dashboardLoaded) supaya
    // tidak bentrok dgn cache versi publik.
    async fetchStatistikDashboard(force = false) {
      if (this.dashboardLoaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/statistik");

        this.statistik = response.data.data;
        this.dashboardLoaded = true;
      } catch (error) {
        console.error("Gagal mengambil statistik dashboard:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchRekapRombel(force = false) {
      if (this.rekapRombelLoaded && !force) return;

      this.rekapRombelLoading = true;

      try {
        const response = await api.get("/public/statistik/rekap-rombel");

        this.rekapRombel = response.data.data || [];
        this.rekapRombelLoaded = true;
      } catch (error) {
        console.error("Gagal mengambil rekap rombel:", error);
      } finally {
        this.rekapRombelLoading = false;
      }
    },

    clearCache() {
      this.loaded = false;
      this.dashboardLoaded = false;
      this.rekapRombelLoaded = false;
    },
  },
});
