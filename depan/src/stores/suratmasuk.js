import { defineStore } from "pinia";
import api from "../services/api";

export const useSuratMasukStore = defineStore("suratmasuk", {
  state: () => ({
    // admin
    adminSuratMasuk: [],
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
    async fetchAdminSuratMasuk() {
      this.adminLoading = true;

      try {
        const response = await api.get("/suratmasuk", {
          params: {
            search: this.filters.q,
            page: this.pagination.page,
            limit: this.pagination.limit,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
          },
        });

        this.adminSuratMasuk = response.data.data || [];
        this.pagination.total = response.data.total || 0;
        this.pagination.page = response.data.page || 1;
        this.pagination.limit = response.data.limit || 10;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );
      } catch (error) {
        console.error("Gagal mengambil surat masuk admin:", error);
      } finally {
        this.adminLoading = false;
      }
    },

    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAdminSuratMasuk();
    },

    async fetchSuratMasukById(id) {
      this.adminLoading = true;

      try {
        const response = await api.get(`/suratmasuk/${id}`);

        this.adminDetail = response.data.data;

        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil surat masuk:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async createSuratMasuk(formData) {
      this.adminLoading = true;

      try {
        const response = await api.post("/suratmasuk", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        return response.data;
      } catch (error) {
        console.error("Gagal membuat surat masuk:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async updateSuratMasuk(id, formData) {
      this.adminLoading = true;

      try {
        const response = await api.put(`/suratmasuk/${id}`, formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui surat masuk:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async deleteSuratMasuk(id) {
      this.adminLoading = true;

      try {
        const response = await api.delete(`/suratmasuk/${id}`);

        this.adminSuratMasuk = this.adminSuratMasuk.filter((s) => s.id !== id);

        return response.data;
      } catch (error) {
        console.error("Gagal menghapus surat masuk:", error);
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
      this.fetchAdminSuratMasuk();
    },

    clearCache() {
      // tidak ada cache publik untuk surat masuk
    },
  },
});
