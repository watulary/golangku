import { defineStore } from "pinia";
import api from "../services/api";

function snapshot(list) {
  return list.map((s) => ({
    nis: s.nis,
    status: s.status || "",
  }));
}

export const useNaikKelasStore = defineStore("keputusannaikkelas", {
  state: () => ({
    daftar: [],
    tingkat: 0,
    original: "[]",
    loading: false,
    submitting: false,
    importing: false,
    templateLoading: false,
  }),

  actions: {
    async fetchDaftar(params) {
      this.loading = true;
      try {
        const response = await api.get("/keputusan-naikkelas/daftar", {
          params,
        });
        this.daftar = response.data.data || [];
        this.tingkat = response.data.tingkat || 0;
        this.original = JSON.stringify(snapshot(this.daftar));
      } catch (error) {
        this.daftar = [];
        this.tingkat = 0;
        this.original = "[]";
        throw error;
      } finally {
        this.loading = false;
      }
    },

    clear() {
      this.daftar = [];
      this.tingkat = 0;
      this.original = "[]";
    },

    hasChanges() {
      return JSON.stringify(snapshot(this.daftar)) !== this.original;
    },

    async submitBulk(payload) {
      this.submitting = true;
      try {
        const items = this.daftar.map((s) => ({
          nis: s.nis,
          status: s.status || "",
        }));
        const response = await api.post("/keputusan-naikkelas/bulk", {
          ...payload,
          items,
        });
        return response.data;
      } finally {
        this.submitting = false;
      }
    },

    async downloadTemplate(params) {
      this.templateLoading = true;
      try {
        const response = await api.get("/keputusan-naikkelas/template", {
          params,
          responseType: "blob",
        });
        const url = window.URL.createObjectURL(response.data);
        const a = document.createElement("a");
        a.href = url;
        a.download = `template-keputusan-naikkelas.xlsx`;
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
        const response = await api.post("/keputusan-naikkelas/import", fd, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        return response.data;
      } finally {
        this.importing = false;
      }
    },
  },
});
