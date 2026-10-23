import { defineStore } from "pinia";
import api from "../services/api";

export const usePegawaiStore = defineStore("pegawai", {
  state: () => ({
    // admin
    adminPegawai: [],
    adminLoading: false,
    adminDetail: null,
    meData: null,
    meLoading: false,
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      total_page: 1,
    },
    filters: {
      q: "",
      status: 1, // -1 = semua, 1 = aktif, 0 = keluar
      kategori: "Guru", // "" = semua, "guru", "tu"
      pnsNon: "", // "" = semua
    },
    // sort
    sortBy: "nama",
    sortDir: "asc",

    // autocomplete (search-as-you-type)
    options: [],
    optionsLoading: false,

    // export
    exportLoading: false,

    // import
    importLoading: false,
    importErrors: [], // [{row, message}, ...] diisi kalau backend menolak import
  }),

  actions: {
    async fetchAdminPegawai() {
      this.adminLoading = true;

      try {
        const response = await api.get("/pegawai", {
          params: {
            search: this.filters.q,
            status: this.filters.status,
            kategori: this.filters.kategori,
            pns_non: this.filters.pnsNon,
            page: this.pagination.page,
            limit: this.pagination.limit,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
          },
        });

        this.adminPegawai = response.data.data || [];
        this.pagination.total = response.data.total || 0;
        this.pagination.page = response.data.page || 1;
        this.pagination.limit = response.data.limit || 10;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );

        // pagination out-of-range: kalau halaman sekarang melewati total_page
        // (mis. setelah hapus item terakhir di halaman terakhir), mundur & fetch ulang
        if (this.pagination.page > this.pagination.total_page) {
          this.pagination.page = this.pagination.total_page;
          await this.fetchAdminPegawai();
        }
      } catch (error) {
        console.error("Gagal mengambil data pegawai:", error);
      } finally {
        this.adminLoading = false;
      }
    },

    async importPegawai(file) {
      this.importLoading = true;
      this.importErrors = [];
      try {
        const formData = new FormData();
        formData.append("file", file);
        const response = await api.post("/pegawai/import", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        // refetch supaya list & pagination langsung update dengan data baru
        await this.fetchAdminPegawai();
        return response.data; // { success: true, message, data: { imported } }
      } catch (error) {
        // backend mengembalikan { success: false, message, errors: [...] } saat 400
        if (error.response?.data?.errors) {
          this.importErrors = error.response.data.errors;
        }
        console.error("Gagal import pegawai:", error);
        throw error;
      } finally {
        this.importLoading = false;
      }
    },

    async fetchMePegawai() {
      this.meLoading = true;
      try {
        const response = await api.get("/pegawai/me");
        this.meData = response.data.data;
        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil biodata pegawai:", error);
        throw error;
      } finally {
        this.meLoading = false;
      }
    },

    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAdminPegawai();
    },

    toggleSort(column) {
      if (this.sortBy === column) {
        this.sortDir = this.sortDir === "asc" ? "desc" : "asc";
      } else {
        this.sortBy = column;
        this.sortDir = "asc";
      }

      this.pagination.page = 1;
      this.fetchAdminPegawai();
    },

    async fetchPegawaiByNuptk(nuptk) {
      this.adminLoading = true;

      try {
        const response = await api.get(`/pegawai/${nuptk}`);

        this.adminDetail = response.data.data;

        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil pegawai:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async checkDuplicate(field, value, exclude = "") {
      try {
        const response = await api.get("/pegawai/check-duplicate", {
          params: { field, value, exclude },
        });
        return response.data.data; // { exists, nuptk?, nama_pegawai? }
      } catch (error) {
        console.error("Gagal cek duplikat pegawai:", error);
        throw error;
      }
    },

    async createPegawai(payload) {
      this.adminLoading = true;

      try {
        const response = await api.post("/pegawai", payload);

        // reset ke halaman 1 supaya data baru langsung terlihat
        this.pagination.page = 1;

        return response.data;
      } catch (error) {
        console.error("Gagal membuat pegawai:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async updatePegawai(nuptk, payload) {
      this.adminLoading = true;

      try {
        const response = await api.put(`/pegawai/${nuptk}`, payload);

        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui pegawai:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async updateFotoPegawai(nuptk, formData) {
      this.adminLoading = true;

      try {
        const response = await api.put(`/pegawai/${nuptk}/foto`, formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui foto pegawai:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async deleteFotoPegawai(nuptk) {
      this.adminLoading = true;
      try {
        const response = await api.delete(`/pegawai/${nuptk}/foto`);
        return response.data.data; // { foto_url }
      } catch (error) {
        console.error("Gagal menghapus foto pegawai:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async deletePegawai(nuptk) {
      this.adminLoading = true;

      try {
        const response = await api.delete(`/pegawai/${nuptk}`);

        // refetch penuh supaya list & pagination.total selalu sinkron
        await this.fetchAdminPegawai();

        return response.data;
      } catch (error) {
        console.error("Gagal menghapus pegawai:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async searchPegawaiOption(q) {
      this.optionsLoading = true;

      try {
        const response = await api.get("/pegawai/search", {
          params: { q },
        });

        this.options = response.data.data || [];

        return this.options;
      } catch (error) {
        console.error("Gagal mencari pegawai:", error);
        throw error;
      } finally {
        this.optionsLoading = false;
      }
    },

    async exportPegawai() {
      this.exportLoading = true;
      try {
        const response = await api.get("/pegawai/export", {
          params: {
            search: this.filters.q,
            status: this.filters.status,
            kategori: this.filters.kategori,
            pns_non: this.filters.pnsNon,
          },
          responseType: "blob",
        });
        const disposition = response.headers["content-disposition"] || "";
        const match = disposition.match(/filename="?([^"]+)"?/);
        const filename = match ? match[1] : "daftar-pegawai.xlsx";
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", filename);
        document.body.appendChild(link);
        link.click();
        link.remove();
        window.URL.revokeObjectURL(url);
      } catch (error) {
        console.error("Gagal export pegawai:", error);
        throw error;
      } finally {
        this.exportLoading = false;
      }
    },

    async downloadTemplate() {
      try {
        const response = await api.get("/pegawai/template", {
          responseType: "blob",
        });
        const disposition = response.headers["content-disposition"] || "";
        const match = disposition.match(/filename="?([^"]+)"?/);
        const filename = match ? match[1] : "template-import-pegawai.xlsx";
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", filename);
        document.body.appendChild(link);
        link.click();
        link.remove();
        window.URL.revokeObjectURL(url);
      } catch (error) {
        console.error("Gagal mengunduh template pegawai:", error);
        throw error;
      }
    },
  },
});
