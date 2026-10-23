import { defineStore } from "pinia";
import api from "../services/api";

function triggerDownload(res) {
  const disposition = res.headers["content-disposition"] || "";
  const match = disposition.match(/filename="?([^"]+)"?/);
  const filename = match ? match[1] : "data.xlsx";

  const url = window.URL.createObjectURL(new Blob([res.data]));
  const link = document.createElement("a");
  link.href = url;
  link.setAttribute("download", filename);
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);
}

export const useSpmbStore = defineStore("spmb", {
  state: () => ({
    // ==== Publik ====
    data: [],
    loading: false,
    loaded: false,

    // ==== Admin ====
    adminSpmb: [],
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
      jalur: 0,
      lapor_diri: -1,
    },
    sortBy: "nama_peserta",
    sortDir: "asc",

    importLoading: false,
    importErrors: [],
  }),
  actions: {
    // ==== Publik ====
    async fetchData(force = false) {
      if (this.loaded && !force) return;
      this.loading = true;
      try {
        const res = await api.get("/public/spmb");
        this.data = res.data.data || [];
        this.loaded = true;
      } catch (error) {
        console.error("Gagal mengambil data SPMB", error);
      } finally {
        this.loading = false;
      }
    },
    clearCache() {
      this.data = [];
      this.loaded = false;
    },

    // ==== Admin ====
    async fetchAdminSpmb() {
      this.adminLoading = true;
      try {
        const res = await api.get("/spmb", {
          params: {
            search: this.filters.q,
            jalur: this.filters.jalur,
            lapor_diri: this.filters.lapor_diri,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
            limit: this.pagination.limit,
            page: this.pagination.page,
          },
        });
        this.adminSpmb = res.data.data || [];
        this.pagination.total = res.data.pagination?.total || 0;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );
        if (this.pagination.page > this.pagination.total_page) {
          this.pagination.page = this.pagination.total_page;
          await this.fetchAdminSpmb();
          return;
        }
      } catch (error) {
        console.error("Gagal mengambil data admin SPMB", error);
      } finally {
        this.adminLoading = false;
      }
    },

    async fetchSpmbById(noPeserta) {
      this.adminLoading = true;
      try {
        const res = await api.get(`/spmb/${noPeserta}`);
        this.adminDetail = res.data.data;
        return this.adminDetail;
      } catch (error) {
        console.error("Gagal mengambil detail SPMB", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async createSpmb(payload) {
      this.adminLoading = true;
      try {
        await api.post("/spmb", payload);
      } catch (error) {
        console.error("Gagal menambah data SPMB", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async updateSpmb(noPeserta, payload) {
      this.adminLoading = true;
      try {
        const res = await api.put(`/spmb/${noPeserta}`, payload);
        return res.data;
      } catch (error) {
        console.error("Gagal memperbarui data SPMB", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async deleteSpmb(noPeserta) {
      this.adminLoading = true;
      try {
        await api.delete(`/spmb/${noPeserta}`);
        await this.fetchAdminSpmb();
      } catch (error) {
        console.error("Gagal menghapus data SPMB", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    // ==== Bulk actions ====
    async bulkUpdateLaporDiri(noPesertaList, laporDiri) {
      this.adminLoading = true;
      try {
        const res = await api.post("/spmb/bulk-lapor-diri", {
          no_peserta: noPesertaList,
          lapor_diri: laporDiri,
        });
        await this.fetchAdminSpmb();
        return res.data;
      } catch (error) {
        console.error("Gagal memperbarui status lapor diri", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async bulkDelete(noPesertaList) {
      this.adminLoading = true;
      try {
        const res = await api.post("/spmb/bulk-delete", {
          no_peserta: noPesertaList,
        });
        await this.fetchAdminSpmb();
        return res.data;
      } catch (error) {
        console.error("Gagal menghapus data terpilih", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    // ==== Import / Export / Template ====
    async importSpmb(file) {
      this.importLoading = true;
      this.importErrors = [];
      try {
        const formData = new FormData();
        formData.append("file", file);
        const res = await api.post("/spmb/import", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        await this.fetchAdminSpmb();
        return res.data;
      } catch (error) {
        if (error.response?.data?.errors) {
          this.importErrors = error.response.data.errors;
        }
        console.error("Gagal mengimpor data SPMB", error);
        throw error;
      } finally {
        this.importLoading = false;
      }
    },

    async downloadTemplate() {
      const res = await api.get("/spmb/template", { responseType: "blob" });
      triggerDownload(res);
    },

    async exportSpmb() {
      const res = await api.get("/spmb/export", {
        params: {
          search: this.filters.q,
          jalur: this.filters.jalur,
          lapor_diri: this.filters.lapor_diri,
        },
        responseType: "blob",
      });
      triggerDownload(res);
    },

    // ==== Sort & filter helper ====
    toggleSort(column) {
      if (this.sortBy === column) {
        this.sortDir = this.sortDir === "asc" ? "desc" : "asc";
      } else {
        this.sortBy = column;
        this.sortDir = "asc";
      }
      this.pagination.page = 1;
      this.fetchAdminSpmb();
    },
    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAdminSpmb();
    },
  },
});
