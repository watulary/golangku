import { defineStore } from "pinia";
import api from "@/services/api";

export const usePenggunaStore = defineStore("pengguna", {
  state: () => ({
    loading: false,

    pengguna: [],

    // ID pengguna yang dipilih
    selectedIds: [],

    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      total_page: 0,
    },

    filters: {
      q: "",
      level: "",
      aktif: "",
      sort: "last_login",
      order: "desc",
    },
  }),

  getters: {
    // Semua data pada halaman ini terpilih
    isAllSelected(state) {
      return (
        state.pengguna.length > 0 &&
        state.selectedIds.length === state.pengguna.length
      );
    },

    // Sebagian dipilih (untuk checkbox indeterminate)
    isIndeterminate(state) {
      return (
        state.selectedIds.length > 0 &&
        state.selectedIds.length < state.pengguna.length
      );
    },
  },

  actions: {
    async fetchData() {
      this.loading = true;

      try {
        const response = await api.get("/pengguna", {
          params: {
            page: this.pagination.page,
            limit: this.pagination.limit,

            q: this.filters.q,
            level: this.filters.level,
            aktif: this.filters.aktif,

            sort: this.filters.sort,
            order: this.filters.order,
          },
        });

        this.pengguna = response.data.data ?? [];
        this.pagination = response.data.pagination;

        // Reset pilihan ketika data berubah
        this.selectedIds = [];
      } finally {
        this.loading = false;
      }
    },

    async refresh() {
      await this.fetchData();
    },

    async changePage(page) {
      this.pagination.page = page;
      await this.fetchData();
    },

    async changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      await this.fetchData();
    },

    async nextPage() {
      if (this.pagination.page < this.pagination.total_page) {
        this.pagination.page++;
        await this.fetchData();
      }
    },

    async prevPage() {
      if (this.pagination.page > 1) {
        this.pagination.page--;
        await this.fetchData();
      }
    },

    // ===============================
    // Bulk Selection
    // ===============================

    toggleSelection(id) {
      const index = this.selectedIds.indexOf(id);

      if (index === -1) {
        this.selectedIds.push(id);
      } else {
        this.selectedIds.splice(index, 1);
      }
    },

    toggleSelectAll() {
      if (this.isAllSelected) {
        this.selectedIds = [];
      } else {
        this.selectedIds = this.pengguna.map((item) => item.id_pengguna);
      }
    },

    clearSelection() {
      this.selectedIds = [];
    },

    async bulkActivate() {
      const { data } = await api.post("/pengguna/bulk/activate", {
        ids: this.selectedIds,
      });

      this.clearSelection();
      await this.fetchData();

      return data;
    },

    async bulkDeactivate() {
      const { data } = await api.post("/pengguna/bulk/deactivate", {
        ids: this.selectedIds,
      });

      this.clearSelection();
      await this.fetchData();

      return data;
    },

    async bulkResetPassword() {
      const { data } = await api.post("/pengguna/bulk/reset-password", {
        ids: this.selectedIds,
      });

      this.clearSelection();
      await this.fetchData();

      return data;
    },

    async bulkDelete() {
      const { data } = await api.post("/pengguna/bulk/delete", {
        ids: this.selectedIds,
      });

      this.clearSelection();
      await this.fetchData();

      return data;
    },

    // ===============================
    // Create Pengguna
    // ===============================

    async createPengguna(payload) {
      const { data } = await api.post("/pengguna", payload);

      return data;
    },

    // ===============================
    // Check Duplicate (realtime, Tambah/Ubah)
    // ===============================

    async checkDuplicate(field, value, exclude = "") {
      const { data } = await api.get("/pengguna/check-duplicate", {
        params: { field, value, exclude },
      });

      return data;
    },

    // ===============================
    // Get Pengguna
    // ===============================

    async getPengguna(id) {
      const { data } = await api.get(`/pengguna/${id}`);

      return data.data;
    },

    // ===============================
    // Update Pengguna
    // ===============================

    async updatePengguna(id, payload) {
      const { data } = await api.put(`/pengguna/${id}`, payload);

      return data;
    },

    // ===============================
    // Foto Pengguna
    // ===============================

    async updateFotoPengguna(id, formData) {
      const { data } = await api.put(`/pengguna/${id}/foto`, formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });

      return data;
    },

    async deleteFotoPengguna(id) {
      const { data } = await api.delete(`/pengguna/${id}/foto`);

      return data;
    },

    async exportExcel() {
      const response = await api.get("/pengguna/export-excel", {
        responseType: "blob",
      });

      // Ambil nama file dari header
      const disposition = response.headers["content-disposition"];

      let filename = "pengguna.xlsx";

      if (disposition) {
        const match = disposition.match(/filename="?([^"]+)"?/);

        if (match) {
          filename = match[1];
        }
      }

      // Download
      const url = window.URL.createObjectURL(new Blob([response.data]));

      const link = document.createElement("a");

      link.href = url;
      link.download = filename;

      document.body.appendChild(link);

      link.click();

      link.remove();

      window.URL.revokeObjectURL(url);
    },

    async uploadExcel(file) {
      const formData = new FormData();

      formData.append("file", file);

      const { data } = await api.post("/pengguna/upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });

      // Refresh tabel setelah upload berhasil
      await this.fetchData();

      return data;
    },
  },
});
