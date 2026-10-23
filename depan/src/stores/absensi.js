import { defineStore } from "pinia";
import api from "../services/api";

export const useAbsensiStore = defineStore("absensi", {
  state: () => ({
    daftarSiswa: [],
    originalSiswa: [],
    siswaList: [],
    siswaLoading: false,
    loading: false,
    submitting: false,
    exportLoading: false,
    ringkasan: { total: 0, hadir: 0, sakit: 0, izin: 0, alpa: 0 },
    rekapSiswa: [],
    rekapJumlahHari: 31,
    rekapLoading: false,
    rekapExportLoading: false,
    filters: {
      tanggal: "",
      id_kelas: null,
      id_rombel: null,
    },
  }),

  actions: {
    async fetchDaftar(tanggal, idKelas, idRombel) {
      this.loading = true;
      try {
        const res = await api.get("/absensi/daftar", {
          params: { tanggal, id_kelas: idKelas, id_rombel: idRombel },
        });
        this.daftarSiswa = res.data.data || [];
        this.originalSiswa = JSON.parse(JSON.stringify(this.daftarSiswa));
        this.filters.tanggal = tanggal;
        this.filters.id_kelas = idKelas;
        this.filters.id_rombel = idRombel;
      } catch (err) {
        console.error("Gagal mengambil daftar absensi:", err);
        this.daftarSiswa = [];
        this.originalSiswa = [];
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async submitBulk(tanggal, idKelas, idRombel, items) {
      this.submitting = true;
      try {
        await api.post("/absensi/bulk", {
          tanggal,
          id_kelas: idKelas,
          id_rombel: idRombel,
          items,
        });
        this.originalSiswa = JSON.parse(JSON.stringify(this.daftarSiswa));
      } catch (err) {
        console.error("Gagal menyimpan absensi:", err);
        throw err;
      } finally {
        this.submitting = false;
      }
    },

    async exportAbsensi(tanggal, idKelas, idRombel, status = "") {
      this.exportLoading = true;
      try {
        const response = await api.get("/absensi/export", {
          params: { tanggal, id_kelas: idKelas, id_rombel: idRombel, status },
          responseType: "blob",
        });
        const disposition = response.headers["content-disposition"] || "";
        const match = disposition.match(/filename="?([^"]+)"?/);
        const filename = match ? match[1] : `laporan-absensi-${tanggal}.xlsx`;
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", filename);
        document.body.appendChild(link);
        link.click();
        link.remove();
        window.URL.revokeObjectURL(url);
      } catch (error) {
        console.error("Gagal export absensi:", error);
        throw error;
      } finally {
        this.exportLoading = false;
      }
    },

    async fetchRekapBulanan(bulan, tahun, idKelas, idRombel) {
      this.rekapLoading = true;
      try {
        const res = await api.get("/absensi/rekap-bulanan", {
          params: { bulan, tahun, id_kelas: idKelas, id_rombel: idRombel },
        });
        this.rekapSiswa = res.data.data.siswa || [];
        this.rekapJumlahHari = res.data.data.jumlah_hari || 31;
      } catch (err) {
        console.error("Gagal mengambil rekap bulanan absensi:", err);
        this.rekapSiswa = [];
        throw err;
      } finally {
        this.rekapLoading = false;
      }
    },

    async exportRekapBulanan(bulan, tahun, idKelas, idRombel) {
      this.rekapExportLoading = true;
      try {
        const response = await api.get("/absensi/rekap-bulanan/export", {
          params: { bulan, tahun, id_kelas: idKelas, id_rombel: idRombel },
          responseType: "blob",
        });
        const disposition = response.headers["content-disposition"] || "";
        const match = disposition.match(/filename="?([^"]+)"?/);
        const filename = match
          ? match[1]
          : `rekap-absensi-${bulan}-${tahun}.xlsx`;
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", filename);
        document.body.appendChild(link);
        link.click();
        link.remove();
        window.URL.revokeObjectURL(url);
      } catch (error) {
        console.error("Gagal export rekap bulanan absensi:", error);
        throw error;
      } finally {
        this.rekapExportLoading = false;
      }
    },

    async fetchAbsensiSaya(bulan, tahun) {
      this.siswaLoading = true;
      try {
        const res = await api.get("/absensi/saya", {
          params: { bulan, tahun },
        });
        this.siswaList = res.data.data.list || [];
      } catch (err) {
        console.error("Gagal mengambil absensi saya:", err);
        this.siswaList = [];
        throw err;
      } finally {
        this.siswaLoading = false;
      }
    },

    hasChanges() {
      // siswa dengan id === 0 berarti belum pernah ada record di DB untuk tanggal ini —
      // itu tetap harus disimpan meski statusnya masih default "Hadir" (belum diubah admin)
      const adaYangBelumTersimpan = this.daftarSiswa.some(
        (s) => !s.id || s.id === 0,
      );
      if (adaYangBelumTersimpan) return true;

      if (this.daftarSiswa.length !== this.originalSiswa.length) return true;
      return this.daftarSiswa.some((s) => {
        const orig = this.originalSiswa.find((o) => o.nis === s.nis);
        if (!orig) return true;
        return s.status !== orig.status || s.keterangan !== orig.keterangan;
      });
    },

    updateStatusLokal(nis, status) {
      const siswa = this.daftarSiswa.find((s) => s.nis === nis);
      if (siswa) siswa.status = status;
    },

    updateKeteranganLokal(nis, keterangan) {
      const siswa = this.daftarSiswa.find((s) => s.nis === nis);
      if (siswa) siswa.keterangan = keterangan;
    },

    clearDaftar() {
      this.daftarSiswa = [];
      this.originalSiswa = [];
    },

    setRingkasan(items) {
      this.ringkasan.total = items.length;
      this.ringkasan.hadir = items.filter((i) => i.status === "Hadir").length;
      this.ringkasan.sakit = items.filter((i) => i.status === "Sakit").length;
      this.ringkasan.izin = items.filter((i) => i.status === "Izin").length;
      this.ringkasan.alpa = items.filter((i) => i.status === "Alpa").length;
    },
  },
});
