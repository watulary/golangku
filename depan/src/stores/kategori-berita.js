import { defineStore } from "pinia";
import api from "../services/api";

export const useKategoriBeritaStore = defineStore("kategoriBerita", {
  state: () => ({
    // dropdown (form berita)
    list: [],
    loading: false,
    loaded: false,

    // admin
    adminKategoriBerita: [],
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

    sortBy: "urutan",
    sortDir: "asc",
  }),

  actions: {
    async fetchKategoriBerita(force = false) {
      if (this.loaded && !force) return;

      this.loading = true;

      try {
        const response = await api.get("/kategori-berita");

        this.list = response.data.data || [];
        this.loaded = true;
      } catch (error) {
        console.error("Gagal mengambil kategori berita:", error);
      } finally {
        this.loading = false;
      }
    },

    async fetchAdminKategoriBerita() {
      this.adminLoading = true;

      try {
        const response = await api.get("/kategori-berita/list", {
          params: {
            search: this.filters.q,
            page: this.pagination.page,
            limit: this.pagination.limit,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
          },
        });

        this.adminKategoriBerita = response.data.data || [];
        this.pagination.total = response.data.total || 0;
        this.pagination.limit = response.data.limit || 10;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );

        // kalau halaman yang sedang aktif sudah di luar jangkauan
        // (mis. setelah delete item terakhir di halaman ini), mundur & fetch ulang
        if (this.pagination.page > this.pagination.total_page) {
          this.pagination.page = this.pagination.total_page;
          return this.fetchAdminKategoriBerita();
        }

        this.pagination.page = response.data.page || this.pagination.page;
      } catch (error) {
        console.error("Gagal mengambil kategori berita admin:", error);
      } finally {
        this.adminLoading = false;
      }
    },

    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAdminKategoriBerita();
    },

    toggleSort(column) {
      if (this.sortBy === column) {
        this.sortDir = this.sortDir === "asc" ? "desc" : "asc";
      } else {
        this.sortBy = column;
        this.sortDir = "asc";
      }

      this.pagination.page = 1;
      this.fetchAdminKategoriBerita();
    },

    async fetchKategoriBeritaById(id) {
      this.adminLoading = true;

      try {
        const response = await api.get(`/kategori-berita/${id}`);

        this.adminDetail = response.data.data;

        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil kategori berita:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async createKategoriBerita(payload) {
      this.adminLoading = true;

      try {
        const response = await api.post("/kategori-berita", payload);

        return response.data;
      } catch (error) {
        console.error("Gagal membuat kategori berita:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async updateKategoriBerita(id, payload) {
      this.adminLoading = true;

      try {
        const response = await api.put(`/kategori-berita/${id}`, payload);

        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui kategori berita:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async deleteKategoriBerita(id) {
      this.adminLoading = true;

      try {
        const response = await api.delete(`/kategori-berita/${id}`);

        await this.fetchAdminKategoriBerita();

        return response.data;
      } catch (error) {
        console.error("Gagal menghapus kategori berita:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },
  },
});
