import { defineStore } from "pinia";
import api from "../services/api";

// Nilai akhir dikirim sebagai angka, atau null kalau kosong.
function toNumberOrNull(v) {
  if (v === null || v === undefined) return null;
  if (typeof v === "string") {
    const s = v.trim().replace(",", ".");
    if (s === "") return null;
    const n = Number(s);
    return Number.isNaN(n) ? NaN : n;
  }
  return v;
}

function snapshot(list) {
  return list.map((s) => ({
    nis: s.nis,
    nilai_akhir: toNumberOrNull(s.nilai_akhir),
    capaian_pembelajaran: (s.capaian_pembelajaran || "").trim(),
  }));
}

export const useNilaiUjianStore = defineStore("nilaiujian", {
  state: () => ({
    daftar: [],
    original: "[]",
    loading: false,
    submitting: false,
    importing: false,
    templateLoading: false,

    // mapel yang bisa dinilai (aktif, tanpa sub-mapel)
    mapelOptions: [],
  }),

  actions: {
    async fetchMapel() {
      try {
        const response = await api.get("/mapel/dinilai");
        this.mapelOptions = response.data.data || [];
      } catch (error) {
        console.error("Gagal mengambil mapel:", error);
        this.mapelOptions = [];
      }
    },

    async fetchDaftar(params) {
      this.loading = true;
      try {
        const response = await api.get("/nilai-ujian/daftar", { params });
        this.daftar = (response.data.data || []).map((s) => ({
          ...s,
          nilai_akhir: s.nilai_akhir === null ? "" : s.nilai_akhir,
        }));
        this.original = JSON.stringify(snapshot(this.daftar));
      } catch (error) {
        this.daftar = [];
        this.original = "[]";
        throw error;
      } finally {
        this.loading = false;
      }
    },

    clear() {
      this.daftar = [];
      this.original = "[]";
    },

    hasChanges() {
      return JSON.stringify(snapshot(this.daftar)) !== this.original;
    },

    // pesan kesalahan pertama, atau "" kalau semua valid
    validate() {
      for (const s of this.daftar) {
        const n = toNumberOrNull(s.nilai_akhir);
        if (Number.isNaN(n) || (n !== null && (n < 0 || n > 100))) {
          return `Nilai akhir ${s.nama_siswa} harus angka 0 - 100.`;
        }
      }
      return "";
    },

    async submitBulk(payload) {
      this.submitting = true;
      try {
        const items = this.daftar.map((s) => ({
          nis: s.nis,
          nilai_akhir: toNumberOrNull(s.nilai_akhir),
          capaian_pembelajaran: s.capaian_pembelajaran || "",
        }));
        const response = await api.post("/nilai-ujian/bulk", {
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
        const response = await api.get("/nilai-ujian/template", {
          params,
          responseType: "blob",
        });
        const url = window.URL.createObjectURL(response.data);
        const a = document.createElement("a");
        a.href = url;
        a.download = `template-nilai-ujian-${params.jenis}.xlsx`;
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
        const response = await api.post("/nilai-ujian/import", fd, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        return response.data;
      } finally {
        this.importing = false;
      }
    },
  },
});
