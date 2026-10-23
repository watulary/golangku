import { defineStore } from "pinia";
import api from "../services/api";

function snapshot(list) {
  return list.map((s) => ({ nis: s.nis, anggota: !!s.anggota }));
}

export const useKeanggotaanEkskulStore = defineStore("keanggotaanekskul", {
  state: () => ({
    daftar: [],
    original: "",
    loading: false,
    submitting: false,
    importing: false,
    templateLoading: false,

    // ekskul aktif untuk dropdown
    ekskulOptions: [],
  }),

  actions: {
    async fetchEkskul() {
      try {
        const response = await api.get("/keanggotaan-ekskul/ekskul");
        this.ekskulOptions = response.data.data || [];
      } catch (error) {
        console.error("Gagal mengambil ekstrakurikuler:", error);
        this.ekskulOptions = [];
      }
    },

    async fetchDaftar(params) {
      this.loading = true;
      try {
        const response = await api.get("/keanggotaan-ekskul/daftar", {
          params,
        });
        this.daftar = response.data.data || [];
        this.original = JSON.stringify(snapshot(this.daftar));
      } catch (error) {
        this.daftar = [];
        this.original = "";
        throw error;
      } finally {
        this.loading = false;
      }
    },

    clear() {
      this.daftar = [];
      this.original = "";
    },

    hasChanges() {
      return JSON.stringify(snapshot(this.daftar)) !== this.original;
    },

    async submit(payload) {
      this.submitting = true;
      try {
        const nisAnggota = this.daftar
          .filter((s) => s.anggota)
          .map((s) => s.nis);
        const response = await api.post("/keanggotaan-ekskul/simpan", {
          ...payload,
          nis_anggota: nisAnggota,
        });
        return response.data;
      } finally {
        this.submitting = false;
      }
    },

    async downloadTemplate(params) {
      this.templateLoading = true;
      try {
        const response = await api.get("/keanggotaan-ekskul/template", {
          params,
          responseType: "blob",
        });
        const url = window.URL.createObjectURL(response.data);
        const a = document.createElement("a");
        a.href = url;
        a.download = `template-keanggotaan-ekskul-${params.semester}.xlsx`;
        document.body.appendChild(a);
        a.click();
        a.remove();
        window.URL.revokeObjectURL(url);
      } finally {
        this.templateLoading = false;
      }
    },

    async importExcel(file, params) {
      this.importing = true;
      try {
        const fd = new FormData();
        fd.append("file", file);
        Object.entries(params).forEach(([k, v]) => fd.append(k, v));
        const response = await api.post("/keanggotaan-ekskul/import", fd, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        return response.data;
      } finally {
        this.importing = false;
      }
    },
  },
});
