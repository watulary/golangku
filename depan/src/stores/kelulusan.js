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

export const useKelulusanStore = defineStore("kelulusan", {
  state: () => ({
    // ==== Data Kelulusan (admin-only) ====
    adminKelulusan: [],
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
    sortBy: "nama_kelas",
    sortDir: "asc",
    importLoading: false,
    importErrors: [],

    // ==== Pengaturan Pengumuman + Isi SKL (kelulusan_waktu, single-row) ====
    waktu: null,
    waktuLoading: false,

    // ==== Galeri foto kelulusan (fotokelulusan, upload ZIP massal) ====
    fotoList: [],
    fotoLoading: false,
    fotoUploadLoading: false,

    // ==== Peserta Didik ====
    siswaData: null,
    siswaLoading: false,

    // ==== Cetak PDF ====
    cetakLoading: false,
  }),
  actions: {
    // ==== Data Kelulusan ====
    async fetchAdminKelulusan() {
      this.adminLoading = true;
      try {
        const res = await api.get("/kelulusan", {
          params: {
            search: this.filters.q,
            sort_by: this.sortBy,
            sort_dir: this.sortDir,
            limit: this.pagination.limit,
            page: this.pagination.page,
          },
        });
        this.adminKelulusan = res.data.data || [];
        this.pagination.total = res.data.pagination?.total || 0;
        this.pagination.total_page = Math.max(
          1,
          Math.ceil(this.pagination.total / this.pagination.limit),
        );
        if (this.pagination.page > this.pagination.total_page) {
          this.pagination.page = this.pagination.total_page;
          await this.fetchAdminKelulusan();
          return;
        }
      } catch (error) {
        console.error("Gagal mengambil data kelulusan", error);
      } finally {
        this.adminLoading = false;
      }
    },
    async fetchKelulusanById(nis) {
      this.adminLoading = true;
      try {
        const res = await api.get(`/kelulusan/${nis}`);
        this.adminDetail = res.data.data;
        return this.adminDetail;
      } catch (error) {
        console.error("Gagal mengambil detail kelulusan", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },
    async createKelulusan(payload) {
      this.adminLoading = true;
      try {
        await api.post("/kelulusan", payload);
      } catch (error) {
        console.error("Gagal menambah data kelulusan", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },
    async updateKelulusan(nis, payload) {
      this.adminLoading = true;
      try {
        const res = await api.put(`/kelulusan/${nis}`, payload);
        return res.data;
      } catch (error) {
        console.error("Gagal memperbarui data kelulusan", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },
    async deleteKelulusan(nis) {
      this.adminLoading = true;
      try {
        await api.delete(`/kelulusan/${nis}`);
        await this.fetchAdminKelulusan();
      } catch (error) {
        console.error("Gagal menghapus data kelulusan", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },
    // ==== Bulk actions ====
    async bulkDelete(nisList) {
      this.adminLoading = true;
      try {
        const res = await api.post("/kelulusan/bulk-delete", {
          nis: nisList,
        });
        await this.fetchAdminKelulusan();
        return res.data;
      } catch (error) {
        console.error("Gagal menghapus data terpilih", error);
        throw error;
      } finally {
        this.adminLoading = false;
      }
    },
    // ==== Import / Export / Template ====
    async importKelulusan(file) {
      this.importLoading = true;
      this.importErrors = [];
      try {
        const formData = new FormData();
        formData.append("file", file);
        const res = await api.post("/kelulusan/import", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        await this.fetchAdminKelulusan();
        return res.data;
      } catch (error) {
        if (error.response?.data?.errors) {
          this.importErrors = error.response.data.errors;
        }
        console.error("Gagal mengimpor data kelulusan", error);
        throw error;
      } finally {
        this.importLoading = false;
      }
    },
    async downloadTemplate() {
      const res = await api.get("/kelulusan/template", {
        responseType: "blob",
      });
      triggerDownload(res);
    },
    async exportKelulusan() {
      const res = await api.get("/kelulusan/export", {
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
      this.fetchAdminKelulusan();
    },
    changeLimit(limit) {
      this.pagination.limit = Number(limit);
      this.pagination.page = 1;
      this.fetchAdminKelulusan();
    },

    // ==== Pengaturan Pengumuman + Isi SKL (kelulusan_waktu) ====
    async fetchWaktu() {
      this.waktuLoading = true;
      try {
        const res = await api.get("/kelulusan/waktu");
        this.waktu = res.data.data;
        return this.waktu;
      } catch (error) {
        console.error("Gagal mengambil pengaturan pengumuman kelulusan", error);
        throw error;
      } finally {
        this.waktuLoading = false;
      }
    },
    // Form "Pengumuman" (judul/nomor/pengantar/.../tahun_pelajaran)
    async updateWaktu(payload) {
      this.waktuLoading = true;
      try {
        const res = await api.put("/kelulusan/waktu", payload);
        await this.fetchWaktu();
        return res.data;
      } catch (error) {
        console.error(
          "Gagal memperbarui pengaturan pengumuman kelulusan",
          error,
        );
        throw error;
      } finally {
        this.waktuLoading = false;
      }
    },
    // Form "Isi Surat SKL" (kop*_skl/isi*_skl/ttd_skl/...) — endpoint TERPISAH
    // dari updateWaktu, sama seperti simpan()/simpanskl() di controller lama.
    async updateSKL(payload) {
      this.waktuLoading = true;
      try {
        const res = await api.put("/kelulusan/skl", payload);
        await this.fetchWaktu();
        return res.data;
      } catch (error) {
        console.error("Gagal memperbarui isi surat SKL", error);
        throw error;
      } finally {
        this.waktuLoading = false;
      }
    },

    // ==== Galeri foto kelulusan (fotokelulusan) ====
    async fetchFoto() {
      this.fotoLoading = true;
      try {
        const res = await api.get("/kelulusan/foto");
        this.fotoList = res.data.data || [];
        return this.fotoList;
      } catch (error) {
        console.error("Gagal mengambil galeri foto kelulusan", error);
        throw error;
      } finally {
        this.fotoLoading = false;
      }
    },
    async uploadFotoZip(file) {
      this.fotoUploadLoading = true;
      try {
        const formData = new FormData();
        formData.append("gambar_zip", file);
        const res = await api.post("/kelulusan/foto/upload", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        });
        await this.fetchFoto();
        return res.data;
      } catch (error) {
        console.error("Gagal mengunggah zip foto kelulusan", error);
        throw error;
      } finally {
        this.fotoUploadLoading = false;
      }
    },
    async bulkDeleteFoto(idList) {
      this.fotoLoading = true;
      try {
        const res = await api.post("/kelulusan/foto/bulk-delete", {
          id: idList,
        });
        await this.fetchFoto();
        return res.data;
      } catch (error) {
        console.error("Gagal menghapus foto terpilih", error);
        throw error;
      } finally {
        this.fotoLoading = false;
      }
    },

    // ==== Peserta Didik ====
    // Backend return { success, data: { waktu, kelulusan } }. TERBUKTI dari
    // debugging: TIDAK ada interceptor axios yang unwrap res.data lebih
    // lanjut di proyek ini -- res.data di sini = seluruh body apa adanya
    // {success, data:{...}}. Konsumen (pesertadidik/Tampil.vue) HARUS akses
    // siswaData.data.waktu / siswaData.data.kelulusan (dua level), bukan
    // siswaData.waktu langsung.
    async fetchKelulusanPesertadidik() {
      this.siswaLoading = true;
      try {
        const res = await api.get("/kelulusan/me");
        this.siswaData = res.data;
        return this.siswaData;
      } catch (error) {
        console.error("Gagal mengambil data kelulusan peserta didik", error);
        throw error;
      } finally {
        this.siswaLoading = false;
      }
    },

    // ==== Cetak PDF ====
    // Backend return PDF langsung (GET, bukan render halaman lalu convert),
    // jadi cukup buka tab baru ke endpointnya — tidak perlu blob/download
    // manual kecuali browser memblokir popup.
    cetakUrl(nis) {
      return `${api.defaults.baseURL}/kelulusan/cetak/${nis}`;
    },
    async cetakSKL(nis) {
      this.cetakLoading = true;
      try {
        const res = await api.get(`/kelulusan/cetak/${nis}`, {
          responseType: "blob",
        });
        const url = window.URL.createObjectURL(
          new Blob([res.data], { type: "application/pdf" }),
        );
        window.open(url, "_blank");
      } catch (error) {
        // Network Error TANPA response sama sekali (error.response undefined)
        // adalah ciri khas request "dicegat" ekstensi download manager (IDM
        // dkk) sebelum axios sempat terima responsnya -- PDF-nya sebenarnya
        // tetap terunduh normal lewat IDM, cuma browser lapor gagal ke JS.
        // Kalau server BENAR-BENAR merespons error (403/404/500 dgn pesan),
        // error.response ADA -- itu baru genuinely gagal, lempar ke pemanggil.
        console.error("Gagal mencetak SKL", error);
        if (error.response) {
          throw error;
        }
        // Network Error murni: diamkan (jangan lempar), biar UI tidak
        // menampilkan toast "Gagal mencetak SKL" yang membingungkan padahal
        // file sudah terunduh lewat IDM.
      } finally {
        this.cetakLoading = false;
      }
    },
  },
});
