import { defineStore } from "pinia";
import api from "../services/api";

function snapshot(list) {
  return list.map((r) => ({
    nis: r.nis,
    prestasi: (r.prestasi || "").trim(),
    nilai: (r.nilai || "").trim(),
    keterangan: (r.keterangan || "").trim(),
  }));
}

let uid = 0;
function nextKey() {
  uid += 1;
  return `p-${uid}`;
}

export const usePrestasiStore = defineStore("prestasi", {
  state: () => ({
    daftar: [], // baris prestasi (satu baris = satu prestasi)
    siswa: [], // siswa rombel untuk dropdown
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
        const response = await api.get("/prestasi/daftar", { params });
        this.siswa = response.data.siswa || [];
        this.daftar = (response.data.data || []).map((r) => ({
          ...r,
          _key: nextKey(),
        }));
        this.original = JSON.stringify(snapshot(this.daftar));
      } catch (error) {
        this.daftar = [];
        this.siswa = [];
        this.original = "[]";
        throw error;
      } finally {
        this.loading = false;
      }
    },

    clear() {
      this.daftar = [];
      this.siswa = [];
      this.original = "[]";
    },

    tambahBaris(nis = "") {
      this.daftar.push({
        _key: nextKey(),
        id_prestasi: 0,
        nis,
        prestasi: "",
        nilai: "",
        keterangan: "",
      });
    },

    hapusBaris(key) {
      this.daftar = this.daftar.filter((r) => r._key !== key);
    },

    hasChanges() {
      return JSON.stringify(snapshot(this.daftar)) !== this.original;
    },

    // pesan kesalahan pertama, atau "" kalau semua valid
    validate() {
      const seen = new Set();
      for (const [i, r] of this.daftar.entries()) {
        const no = i + 1;
        const prestasi = (r.prestasi || "").trim();
        if (!r.nis) return `Baris ${no}: siswa belum dipilih.`;
        if (!prestasi) return `Baris ${no}: prestasi belum diisi.`;
        if (prestasi.length > 255)
          return `Baris ${no}: prestasi maksimal 255 karakter.`;
        if ((r.nilai || "").trim().length > 20)
          return `Baris ${no}: nilai maksimal 20 karakter.`;
        if ((r.keterangan || "").trim().length > 255)
          return `Baris ${no}: keterangan maksimal 255 karakter.`;
        const key = `${r.nis}|${prestasi.toLowerCase()}`;
        if (seen.has(key))
          return `Baris ${no}: prestasi yang sama terisi dua kali untuk siswa ini.`;
        seen.add(key);
      }
      return "";
    },

    async submitBulk(payload) {
      this.submitting = true;
      try {
        const items = this.daftar.map((r) => ({
          nis: r.nis,
          prestasi: r.prestasi || "",
          nilai: r.nilai || "",
          keterangan: r.keterangan || "",
        }));
        const response = await api.post("/prestasi/bulk", {
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
        const response = await api.get("/prestasi/template", {
          params,
          responseType: "blob",
        });
        const url = window.URL.createObjectURL(response.data);
        const a = document.createElement("a");
        a.href = url;
        a.download = `template-prestasi-${params.semester}.xlsx`;
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
        const response = await api.post("/prestasi/import", fd, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        return response.data;
      } finally {
        this.importing = false;
      }
    },
  },
});
