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

export const useKenaikanStore = defineStore("kenaikan", {
  state: () => ({
    // ==== Data Kenaikan (admin-only, tanpa halaman publik) ====
    adminKenaikan: [],
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
    sortBy: "nama_siswa",
    sortDir: "asc",
    importLoading: false,
    importErrors: [],

    // ==== Pengaturan Pengumuman (kenaikan_waktu, single-row) ====
    waktu: null,
    waktuLoading: false,
    //pesertadidik
    siswaData: null,
    siswaLoading: false,
  }),
  actions: {
    // ==== Data Kenaikan ====
    async fetchAdminKenaikan() {
      this.adminLoading = true;
      try {
        const res = await api.get("/kenaikan", {
          params: {
            search: this.filters.q,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
            limit: this.pagination.limit,
            page: this.pagination.page,
          },
        });
        this.adminKenaikan = res.data.data || [];
        this.pagination.total = res.data.pagination?.total || 0;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );
        if (this.pagination.page > this.pagination.total_page) {
          this.pagination.page = this.pagination.total_page;
          await this.fetchAdminKenaikan();
          return;
        }
      } catch (error) {
        console.error("Gagal mengambil data kenaikan", error);
      } finally {
        this.adminLoading = false;
      }
    },
    async fetchKenaikanById(nis) {
      this.adminLoading = true;
      try {
        const res = await api.get(`/kenaikan/${nis}`);
        this.adminDetail = res.data.data;
        return this.adminDetail;
      } catch (error) {
        console.error("Gagal mengambil detail kenaikan", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },
    async createKenaikan(payload) {
      this.adminLoading = true;
      try {
        await api.post("/kenaikan", payload);
      } catch (error) {
        console.error("Gagal menambah data kenaikan", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },
    async updateKenaikan(nis, payload) {
      this.adminLoading = true;
      try {
        const res = await api.put(`/kenaikan/${nis}`, payload);
        return res.data;
      } catch (error) {
        console.error("Gagal memperbarui data kenaikan", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },
    async deleteKenaikan(nis) {
      this.adminLoading = true;
      try {
        await api.delete(`/kenaikan/${nis}`);
        await this.fetchAdminKenaikan();
      } catch (error) {
        console.error("Gagal menghapus data kenaikan", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },
    // ==== Bulk actions ====
    async bulkDelete(nisList) {
      this.adminLoading = true;
      try {
        const res = await api.post("/kenaikan/bulk-delete", {
          nis: nisList,
        });
        await this.fetchAdminKenaikan();
        return res.data;
      } catch (error) {
        console.error("Gagal menghapus data terpilih", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },
    // ==== Import / Export / Template ====
    async importKenaikan(file) {
      this.importLoading = true;
      this.importErrors = [];
      try {
        const formData = new FormData();
        formData.append("file", file);
        const res = await api.post("/kenaikan/import", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        await this.fetchAdminKenaikan();
        return res.data;
      } catch (error) {
        if (error.response?.data?.errors) {
          this.importErrors = error.response.data.errors;
        }
        console.error("Gagal mengimpor data kenaikan", error);
        throw error;
      } finally {
        this.importLoading = false;
      }
    },
    async downloadTemplate() {
      const res = await api.get("/kenaikan/template", { responseType: "blob" });
      triggerDownload(res);
    },
    async exportKenaikan() {
      const res = await api.get("/kenaikan/export", {
        params: { search: this.filters.q },
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
      this.fetchAdminKenaikan();
    },
    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAdminKenaikan();
    },

    // ==== Pengaturan Pengumuman (kenaikan_waktu) ====
    async fetchWaktu() {
      this.waktuLoading = true;
      try {
        const res = await api.get("/kenaikan-waktu");
        this.waktu = res.data.data;
        return this.waktu;
      } catch (error) {
        console.error("Gagal mengambil pengaturan pengumuman kenaikan", error);
        throw error;
      } finally {
        this.waktuLoading = false;
      }
    },
    async updateWaktu(payload) {
      this.waktuLoading = true;
      try {
        const res = await api.put("/kenaikan-waktu", payload);
        await this.fetchWaktu();
        return res.data;
      } catch (error) {
        console.error(
          "Gagal memperbarui pengaturan pengumuman kenaikan",
          error,
        );
        throw error;
      } finally {
        this.waktuLoading = false;
      }
    },

    //pesertadidik
    async fetchKenaikanPesertadidik() {
      this.siswaLoading = true;
      try {
        const res = await api.get("/kenaikan/me");
        this.siswaData = res.data.data;
        return this.siswaData;
      } catch (error) {
        console.error("Gagal mengambil data kenaikan peserta didik", error);
        throw error;
      } finally {
        this.siswaLoading = false;
      }
    },
  },
});
