import { defineStore } from "pinia";
import api from "../services/api";

export const useWaliKelasStore = defineStore("walikelas", {
  state: () => ({
    list: [],
    loading: false,

    detail: null,
    detailLoading: false,

    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      total_page: 1,
    },
    filters: {
      q: "",
      id_kelas: "",
    },

    sortBy: "kelas",
    sortDir: "asc",

    pegawaiOptions: [],
    pegawaiSearching: false,
  }),

  actions: {
    async fetchAll() {
      this.loading = true;

      try {
        const response = await api.get("/walikelas", {
          params: {
            search: this.filters.q,
            id_kelas: this.filters.id_kelas || undefined,
            page: this.pagination.page,
            limit: this.pagination.limit,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
          },
        });

        this.list = response.data.data || [];
        this.pagination.total = response.data.total || 0;
        this.pagination.page = response.data.page || 1;
        this.pagination.limit = response.data.limit || 10;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );
      } catch (error) {
        console.error("Gagal mengambil wali kelas:", error);
      } finally {
        this.loading = false;
      }
    },

    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAll();
    },

    changeFilterKelas(idKelas) {
      this.filters.id_kelas = idKelas;
      this.pagination.page = 1;
      this.fetchAll();
    },

    toggleSort(column) {
      if (this.sortBy === column) {
        this.sortDir = this.sortDir === "asc" ? "desc" : "asc";
      } else {
        this.sortBy = column;
        this.sortDir = "asc";
      }

      this.pagination.page = 1;
      this.fetchAll();
    },

    async fetchById(id) {
      this.detailLoading = true;

      try {
        const response = await api.get(`/walikelas/${id}`);

        this.detail = response.data.data;

        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil wali kelas:", error);
        throw error;
      } finally {
        this.detailLoading = false;
      }
    },

    async searchPegawai(query) {
      this.pegawaiSearching = true;

      try {
        const response = await api.get("/walikelas/pegawai/search", {
          params: { q: query },
        });

        this.pegawaiOptions = response.data.data || [];
      } catch (error) {
        console.error("Gagal mencari pegawai:", error);
      } finally {
        this.pegawaiSearching = false;
      }
    },

    async create(payload) {
      this.loading = true;

      try {
        const response = await api.post("/walikelas", payload);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal membuat wali kelas:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    async update(id, payload) {
      this.loading = true;

      try {
        const response = await api.put(`/walikelas/${id}`, payload);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui wali kelas:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    async remove(id) {
      this.loading = true;

      try {
        const response = await api.delete(`/walikelas/${id}`);

        await this.fetchAll();

        return response.data;
      } catch (error) {
        console.error("Gagal menghapus wali kelas:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },
  },
});
