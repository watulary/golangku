import { defineStore } from "pinia";
import api from "../services/api";

export const usePesertaDidikStore = defineStore("pesertadidik", {
  state: () => ({
    // ===== Aktif (admin) =====
    adminPesertaDidik: [],
    adminLoading: false,
    adminDetail: null,
    siswaData: null,
    siswaLoading: false,
    pagination: {
      page: 1,
      limit: 10,
      total: 0,
      total_page: 1,
    },
    filters: {
      q: "",
      idTahunPelajaran: 0, // 0 = semua
      idKelas: 0,
      idRombel: 0,
    },
    sortBy: "nama_siswa",
    sortDir: "asc",

    importLoading: false,
    importErrors: [], //diisi kalau backend menolak import

    // ===== Mutasi =====
    adminMutasi: [],
    mutasiLoading: false,
    mutasiPagination: {
      page: 1,
      limit: 10,
      total: 0,
      total_page: 1,
    },
    mutasiFilters: { q: "" },

    mutasiDetail: null,

    // ===== Arsip =====
    adminArsip: [],
    arsipLoading: false,
    arsipDetail: null,
    arsipPagination: {
      page: 1,
      limit: 10,
      total: 0,
      total_page: 1,
    },
    arsipFilters: { q: "", idTahunPelajaran: 0, idKelas: 0, idRombel: 0 },
  }),

  actions: {
    // ===== Aktif =====
    async fetchAdminPesertaDidik() {
      this.adminLoading = true;
      try {
        const response = await api.get("/pesertadidik", {
          params: {
            search: this.filters.q,
            id_tahunpelajaran: this.filters.idTahunPelajaran || undefined,
            id_kelas: this.filters.idKelas || undefined,
            id_rombel: this.filters.idRombel || undefined,
            page: this.pagination.page,
            limit: this.pagination.limit,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
          },
        });
        this.adminPesertaDidik = response.data.data || [];
        this.pagination.total = response.data.pagination?.total || 0;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );
        // pagination out-of-range: kalau halaman sekarang melewati total_page
        // (mis. setelah hapus/mutasi/arsipkan item terakhir di halaman terakhir), mundur & fetch ulang
        if (this.pagination.page > this.pagination.total_page) {
          this.pagination.page = this.pagination.total_page;
          await this.fetchAdminPesertaDidik();
        }
      } catch (error) {
        console.error("Gagal mengambil data peserta didik:", error);
      } finally {
        this.adminLoading = false;
      }
    },

    // actions
    async fetchPesertaDidikMe() {
      this.siswaLoading = true;
      try {
        const response = await api.get("/pesertadidik/me");
        this.siswaData = response.data.data;
        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil biodata peserta didik:", error);
        throw error;
      } finally {
        this.siswaLoading = false;
      }
    },

    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAdminPesertaDidik();
    },

    toggleSort(column) {
      if (this.sortBy === column) {
        this.sortDir = this.sortDir === "asc" ? "desc" : "asc";
      } else {
        this.sortBy = column;
        this.sortDir = "asc";
      }
      this.pagination.page = 1;
      this.fetchAdminPesertaDidik();
    },

    async fetchPesertaDidikByNis(nis) {
      this.adminLoading = true;
      try {
        const response = await api.get(`/pesertadidik/${nis}`);
        this.adminDetail = response.data.data;
        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil peserta didik:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async createPesertaDidik(payload) {
      this.adminLoading = true;
      try {
        const response = await api.post("/pesertadidik", payload);
        this.pagination.page = 1;
        return response.data;
      } catch (error) {
        console.error("Gagal membuat data peserta didik:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async updatePesertaDidik(nis, payload) {
      this.adminLoading = true;
      try {
        const response = await api.put(`/pesertadidik/${nis}`, payload);
        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui data peserta didik:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async updateFotoPesertaDidik(nis, formData) {
      this.adminLoading = true;
      try {
        const response = await api.put(`/pesertadidik/${nis}/foto`, formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        return response.data;
      } catch (error) {
        console.error("Gagal memperbarui foto peserta didik:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async deletePesertaDidik(nis) {
      this.adminLoading = true;
      try {
        const response = await api.delete(`/pesertadidik/${nis}`);
        // refetch penuh supaya list & pagination.total selalu sinkron
        await this.fetchAdminPesertaDidik();
        return response.data;
      } catch (error) {
        console.error("Gagal menghapus data peserta didik:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    // ===== Mutasi =====
    async mutasiPesertaDidik(nis, payload) {
      // payload: { status_keluar, alasan_keluar, tanggal_keluar, tujuan, nomor_surat }
      this.adminLoading = true;
      try {
        const response = await api.post(`/pesertadidik/${nis}/mutasi`, payload);
        await this.fetchAdminPesertaDidik();
        return response.data;
      } catch (error) {
        console.error("Gagal memproses mutasi peserta didik:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async fetchAdminMutasi() {
      this.mutasiLoading = true;
      try {
        const response = await api.get("/pesertadidik/mutasi", {
          params: {
            search: this.mutasiFilters.q,
            page: this.mutasiPagination.page,
            limit: this.mutasiPagination.limit,
          },
        });
        this.adminMutasi = response.data.data || [];
        this.mutasiPagination.total = response.data.pagination?.total || 0;
        this.mutasiPagination.total_page = Math.max(
          1,
          Math.ceil(this.mutasiPagination.total / this.mutasiPagination.limit),
        );
      } catch (error) {
        console.error("Gagal mengambil data mutasi:", error);
      } finally {
        this.mutasiLoading = false;
      }
    },

    // ===== Arsip =====
    async arsipkanPesertaDidik(nis) {
      this.adminLoading = true;
      try {
        const response = await api.post(`/pesertadidik/${nis}/arsip`);
        await this.fetchAdminPesertaDidik();
        return response.data;
      } catch (error) {
        console.error("Gagal mengarsipkan peserta didik:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async fetchAdminArsip() {
      this.arsipLoading = true;
      try {
        const response = await api.get("/pesertadidik/arsip", {
          params: {
            search: this.arsipFilters.q,
            id_tahunpelajaran: this.arsipFilters.idTahunPelajaran || undefined,
            id_kelas: this.arsipFilters.idKelas || undefined,
            id_rombel: this.arsipFilters.idRombel || undefined,
            page: this.arsipPagination.page,
            limit: this.arsipPagination.limit,
          },
        });
        this.adminArsip = response.data.data || [];
        this.arsipPagination.total = response.data.pagination?.total || 0;
        this.arsipPagination.total_page = Math.max(
          1,
          Math.ceil(this.arsipPagination.total / this.arsipPagination.limit),
        );
      } catch (error) {
        console.error("Gagal mengambil data arsip:", error);
      } finally {
        this.arsipLoading = false;
      }
    },

    async fetchArsipById(id) {
      this.arsipLoading = true;
      try {
        const response = await api.get(`/pesertadidik/arsip/${id}`);
        this.arsipDetail = response.data.data;
        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil arsip:", error);
        throw error;
      } finally {
        this.arsipLoading = false;
      }
    },

    async deleteFotoPesertaDidik(nis) {
      this.adminLoading = true;
      try {
        const response = await api.delete(`/pesertadidik/${nis}/foto`);
        return response.data;
      } catch (error) {
        console.error("Gagal menghapus foto peserta didik:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async checkDuplicate(field, value, exclude = "") {
      try {
        const response = await api.get("/pesertadidik/check-duplicate", {
          params: { field, value, exclude },
        });
        return response.data.data; // { exists, nis?, nama_siswa? }
      } catch (error) {
        console.error("Gagal cek duplikat peserta didik:", error);
        throw error;
      }
    },

    async exportPesertaDidik() {
      try {
        const response = await api.get("/pesertadidik/export", {
          params: {
            search: this.filters.q,
            id_tahunpelajaran: this.filters.idTahunPelajaran || undefined,
            id_kelas: this.filters.idKelas || undefined,
            id_rombel: this.filters.idRombel || undefined,
          },
          responseType: "blob",
        });
        const disposition = response.headers["content-disposition"] || "";
        const match = disposition.match(/filename="?([^"]+)"?/);
        const filename = match ? match[1] : "daftar-peserta-didik.xlsx";
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", filename);
        document.body.appendChild(link);
        link.click();
        link.remove();
        window.URL.revokeObjectURL(url);
      } catch (error) {
        console.error("Gagal export peserta didik:", error);
        throw error;
      }
    },

    async importPesertaDidik(file) {
      this.importLoading = true;
      this.importErrors = [];
      try {
        const formData = new FormData();
        formData.append("file", file);
        const response = await api.post("/pesertadidik/import", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        // refetch supaya list & pagination langsung update dengan data baru
        await this.fetchAdminPesertaDidik();
        return response.data; // { success: true, message, data: { imported } }
      } catch (error) {
        // backend mengembalikan { success: false, message, errors: [...] } saat 400
        if (error.response?.data?.errors) {
          this.importErrors = error.response.data.errors;
        }
        console.error("Gagal import peserta didik:", error);
        throw error;
      } finally {
        this.importLoading = false;
      }
    },

    async downloadTemplate() {
      try {
        const response = await api.get("/pesertadidik/template", {
          responseType: "blob",
        });
        const disposition = response.headers["content-disposition"] || "";
        const match = disposition.match(/filename="?([^"]+)"?/);
        const filename = match
          ? match[1]
          : "template-import-peserta-didik.xlsx";
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", filename);
        document.body.appendChild(link);
        link.click();
        link.remove();
        window.URL.revokeObjectURL(url);
      } catch (error) {
        console.error("Gagal mengunduh template peserta didik:", error);
        throw error;
      }
    },

    async exportMutasi() {
      try {
        const response = await api.get("/pesertadidik/mutasi/export", {
          responseType: "blob",
        });
        const disposition = response.headers["content-disposition"] || "";
        const match = disposition.match(/filename="?([^"]+)"?/);
        const filename = match ? match[1] : "daftar-mutasi-peserta-didik.xlsx";
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", filename);
        document.body.appendChild(link);
        link.click();
        link.remove();
        window.URL.revokeObjectURL(url);
      } catch (error) {
        console.error("Gagal export mutasi:", error);
        throw error;
      }
    },

    async exportArsip() {
      try {
        const response = await api.get("/pesertadidik/arsip/export", {
          params: {
            search: this.arsipFilters.q,
            id_tahunpelajaran: this.arsipFilters.idTahunPelajaran || undefined,
            id_kelas: this.arsipFilters.idKelas || undefined,
            id_rombel: this.arsipFilters.idRombel || undefined,
          },
          responseType: "blob",
        });
        const disposition = response.headers["content-disposition"] || "";
        const match = disposition.match(/filename="?([^"]+)"?/);
        const filename = match ? match[1] : "daftar-arsip-peserta-didik.xlsx";
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", filename);
        document.body.appendChild(link);
        link.click();
        link.remove();
        window.URL.revokeObjectURL(url);
      } catch (error) {
        console.error("Gagal export arsip:", error);
        throw error;
      }
    },

    async fetchMutasiByNis(nis) {
      this.mutasiLoading = true;
      try {
        const response = await api.get(`/pesertadidik/mutasi/${nis}`);
        this.mutasiDetail = response.data.data;
        return response.data.data;
      } catch (error) {
        console.error("Gagal mengambil detil mutasi:", error);
        throw error;
      } finally {
        this.mutasiLoading = false;
      }
    },

    async bulkDeletePesertaDidik(nisList) {
      this.adminLoading = true;
      try {
        const response = await api.post("/pesertadidik/bulk/delete", {
          nis: nisList,
        });
        await this.fetchAdminPesertaDidik();
        return response.data;
      } catch (error) {
        console.error("Gagal menghapus massal peserta didik:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async bulkMutasiPesertaDidik(nisList, payload) {
      this.adminLoading = true;
      try {
        const response = await api.post("/pesertadidik/bulk/mutasi", {
          nis: nisList,
          ...payload,
        });
        await this.fetchAdminPesertaDidik();
        return response.data;
      } catch (error) {
        console.error("Gagal memutasi massal peserta didik:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },

    async bulkArsipkanPesertaDidik(nisList) {
      this.adminLoading = true;
      try {
        const response = await api.post("/pesertadidik/bulk/arsip", {
          nis: nisList,
        });
        await this.fetchAdminPesertaDidik();
        return response.data;
      } catch (error) {
        console.error("Gagal mengarsipkan massal peserta didik:", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },
  },
});
