<template>
  <DashboardPageHeader />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Guru Wali"
        description="Isilah nama guru wali dan pilih guru yang bertanggung jawab di bawah ini."
      >
        <UserGroupIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="space-y-6 p-6">
        <!-- Nama Guru Wali -->
        <FormField label="Nama Wali">
          <input
            v-model="form.nama_guruwali"
            type="text"
            :class="inputClass"
            placeholder="Contoh: Wali 1"
          />
        </FormField>

        <!-- Guru (Pegawai) -->
        <FormField label="Guru">
          <div class="relative">
            <input
              v-model="pegawaiQuery"
              @input="onPegawaiSearch"
              @focus="showPegawaiDropdown = true"
              type="text"
              :class="inputClass"
              placeholder="Ketik nama atau NUPTK..."
              autocomplete="off"
            />

            <div
              v-if="showPegawaiDropdown && pegawaiQuery.trim() !== ''"
              class="absolute z-10 mt-1 max-h-56 w-full overflow-y-auto rounded-xl bg-white shadow-lg ring-1 ring-slate-900/5 dark:bg-[#0B0F1A] dark:ring-white/10"
            >
              <div
                v-if="guruwaliStore.pegawaiSearching"
                class="px-4 py-3 text-sm text-gray-500 dark:text-slate-400"
              >
                Mencari...
              </div>

              <div
                v-else-if="guruwaliStore.pegawaiOptions.length === 0"
                class="px-4 py-3 text-sm text-gray-500 dark:text-slate-400"
              >
                Guru tidak ditemukan.
              </div>

              <button
                v-for="p in guruwaliStore.pegawaiOptions"
                :key="p.nuptk"
                type="button"
                @click="selectPegawai(p)"
                class="flex w-full flex-col items-start px-4 py-2.5 text-left transition hover:bg-gray-50 dark:hover:bg-slate-800"
              >
                <span
                  class="text-sm font-medium text-gray-800 dark:text-slate-200"
                >
                  {{ p.nama_pegawai }}
                </span>
                <span class="text-xs text-gray-500 dark:text-slate-400">
                  {{ p.nuptk }}
                </span>
              </button>
            </div>
          </div>

          <p
            v-if="form.nuptk"
            class="mt-1.5 text-xs text-green-600 dark:text-green-400"
          >
            Terpilih: {{ form.nama_pegawai }} ({{ form.nuptk }})
          </p>
        </FormField>

        <!-- Status Aktif -->
        <FormField label="Status">
          <label class="inline-flex items-center gap-2 cursor-pointer">
            <input
              v-model="form.aktif"
              type="checkbox"
              class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500 dark:border-slate-600 dark:bg-[#0B0F1A]"
            />
            <span class="text-sm text-gray-700 dark:text-slate-300">
              Guru Wali Aktif
            </span>
          </label>
        </FormField>
      </div>

      <!-- Aksi -->
      <div
        class="flex flex-col gap-4 rounded-2xl bg-white px-5 py-4 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5 md:flex-row md:items-center md:justify-between"
      >
        <div class="flex gap-3">
          <button
            type="submit"
            :disabled="loading"
            class="inline-flex items-center justify-center gap-2 rounded-xl bg-blue-600 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:bg-gray-400 dark:disabled:bg-slate-700"
          >
            <ArrowPathIcon v-if="loading" class="h-5 w-5 animate-spin" />
            {{ loading ? "Menyimpan..." : "Simpan" }}
          </button>

          <RouterLink
            :to="{ name: 'dashboard-guruwali' }"
            class="inline-flex items-center justify-center rounded-xl border border-gray-300 px-5 py-2.5 text-sm font-medium text-gray-700 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
          >
            Batal
          </RouterLink>
        </div>
      </div>
    </section>
  </form>
</template>

<script setup>
import { onMounted, reactive, ref } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import FormField from "@/components/form/FormField.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";
import { useToast } from "@/composables/useToast";
import { useGuruWaliStore } from "@/stores/guruwali";
import { ArrowPathIcon, UserGroupIcon } from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const guruwaliStore = useGuruWaliStore();
const toast = useToast();

const isEdit = !!route.params.id;

const loading = ref(false);

const pegawaiQuery = ref("");
const showPegawaiDropdown = ref(false);

const form = reactive({
  id_guruwali: null,
  nama_guruwali: "",
  nuptk: "",
  nama_pegawai: "",
  aktif: true,
});

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";

const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500 disabled:cursor-not-allowed disabled:bg-gray-100 dark:disabled:bg-slate-800";

// =====================================
// Pegawai Search
// =====================================

let searchTimeout = null;

function onPegawaiSearch() {
  showPegawaiDropdown.value = true;
  form.nuptk = "";
  form.nama_pegawai = "";

  clearTimeout(searchTimeout);

  searchTimeout = setTimeout(() => {
    if (pegawaiQuery.value.trim() !== "") {
      guruwaliStore.searchPegawai(pegawaiQuery.value.trim());
    }
  }, 300);
}

function selectPegawai(pegawai) {
  form.nuptk = pegawai.nuptk;
  form.nama_pegawai = pegawai.nama_pegawai;
  pegawaiQuery.value = pegawai.nama_pegawai;
  showPegawaiDropdown.value = false;
}

// =====================================
// Load Data (mode edit)
// =====================================

async function loadData() {
  if (!isEdit) return;

  loading.value = true;

  try {
    const data = await guruwaliStore.fetchById(route.params.id);

    Object.assign(form, {
      id_guruwali: data.id_guruwali,
      nama_guruwali: data.nama_guruwali,
      nuptk: data.nuptk,
      nama_pegawai: data.nama_pegawai || "",
      aktif: !!data.aktif,
    });

    pegawaiQuery.value = data.nama_pegawai || "";
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat data guru wali.",
    );
  } finally {
    loading.value = false;
  }
}

// =====================================
// Submit
// =====================================

async function submit() {
  if (loading.value) return;

  const namaGuruWali = form.nama_guruwali.trim();

  if (!namaGuruWali) {
    toast.error("Nama guru wali wajib diisi.");
    return;
  }

  if (!form.nuptk) {
    toast.error("Guru wajib dipilih.");
    return;
  }

  loading.value = true;

  try {
    const payload = {
      nama_guruwali: namaGuruWali,
      nuptk: form.nuptk,
      aktif: form.aktif,
    };

    if (isEdit) {
      const response = await guruwaliStore.update(form.id_guruwali, payload);

      // "changed" dari backend berupa array nama kolom yang berubah — array
      // selalu truthy di JS meski kosong, jadi cek panjangnya, bukan sekadar
      // keberadaannya (beda dari pola ekskul.js yang mengasumsikan boolean).
      if (response.changed && response.changed.length > 0) {
        toast.success(response.message || "Guru wali berhasil diperbarui.");
      } else {
        toast.warning(response.message || "Tidak ada perubahan data.");
      }
    } else {
      await guruwaliStore.create(payload);
      toast.success("Guru wali berhasil ditambahkan.");
    }

    router.push({ name: "dashboard-guruwali" });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menyimpan guru wali.");
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  await loadData();
});
</script>
