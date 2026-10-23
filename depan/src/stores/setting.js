import { defineStore } from "pinia";
import api from "../services/api";

export const useSettingStore = defineStore("setting", {
  state: () => ({
    // Publik
    setting: null,
    loading: false,
    loaded: false,

    // Dashboard Setting
    dashboardSetting: null,
    dashboardLoading: false,
  }),

  actions: {
    async fetchSetting() {
      // sudah pernah diambil
      if (this.loaded) return;

      this.loading = true;

      try {
        const response = await api.get("/public/setting");

        this.setting = response.data.data;
        this.loaded = true;
      } catch (error) {
        console.error(error);
      } finally {
        this.loading = false;
      }
    },

    // Dashboard setting
    async getDashboardSetting() {
      this.dashboardLoading = true;

      try {
        const response = await api.get("/setting");

        this.dashboardSetting = response.data.data;

        return response.data;
      } finally {
        this.dashboardLoading = false;
      }
    },

    async updateDashboardSetting(payload) {
      const response = await api.put("/setting", payload);

      if (response.data.data) {
        this.dashboardSetting = response.data.data;
      }

      return response.data;
    },
  },
});
