import { defineStore } from "pinia";
import api from "../services/api";

export const useFotoBannerStore = defineStore("fotobanner", {
  state: () => ({
    banners: [],
    loading: false,

    // cache
    loaded: false,

    // admin
    adminBanners: [],
    adminLoading: false,
    adminDetail: null,
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      total_page: 1,
    },
    filters: {
      q: "",
    },

    // sort
    sortBy: "waktu",
    sortDir: "desc",
  }),

  actions: {
    async fetchBanners(force = false) {
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/public/fotobanner");

        this.banners = response.data.data || [];
        this.loaded = true;
      } catch (error) {
        console.error("Gagal mengambil banner:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchAdminBanners() {
      this.adminLoading = true;

      try {
        const response = await api.get("/fotobanner", {
          params: {
            search: this.filters.q,
            page: this.pagination.page,
            limit: this.pagination.limit,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
          },
        });

        this.adminBanners = response.data.data || [];
        this.pagination.total = response.data.total || 0;
        this.pagination.page = response.data.page || 1;
        this.pagination.limit = response.data.limit || 10;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );
      } catch (error) {
        console.error("Gagal mengambil banner admin:", error);
      } finally {
        this.adminLoading = false;
      }
    },

    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAdminBanners();
    },

    async fetchBannerById(id) {
      this.adminLoading = true;

      try {
        const response = await api.get(`/fotobanner/${id}`);

        this.adminDetail = response.data.data;

        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil banner:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async createBanner(formData) {
      this.adminLoading = true;

      try {
        const response = await api.post("/fotobanner", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        return response.data;
      } catch (error) {
        console.error("Gagal membuat banner:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async updateBanner(id, formData) {
      this.adminLoading = true;

      try {
        const response = await api.put(`/fotobanner/${id}`, formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui banner:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async deleteBanner(id) {
      this.adminLoading = true;

      try {
        const response = await api.delete(`/fotobanner/${id}`);

        this.adminBanners = this.adminBanners.filter((b) => b.id !== id);

        return response.data;
      } catch (error) {
        console.error("Gagal menghapus banner:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    toggleSort(column) {
      if (this.sortBy === column) {
        this.sortDir = this.sortDir === "asc" ? "desc" : "asc";
      } else {
        this.sortBy = column;
        this.sortDir = "asc";
      }

      this.pagination.page = 1;
      this.fetchAdminBanners();
    },

    clearCache() {
      this.loaded = false;
    },
  },
});
