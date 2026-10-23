<template>
  <DashboardPageHeader />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Kategori Berita"
        description="Isilah kategori berita dan urutan nya di bawah ini."
      >
        <TagIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="space-y-6 p-6">
        <!-- Nama -->
        <FormField label="Nama Kategori">
          <input
            v-model="form.nama_kategoriberita"
            type="text"
            :class="inputClass"
            placeholder="Contoh: Akademik"
            maxlength="200"
          />
        </FormField>

        <!-- Urutan -->
        <FormField label="Urutan">
          <input
            v-model.number="form.urutan"
            type="number"
            min="0"
            :class="inputClass"
            placeholder="0"
          />
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
            :to="{ name: 'dashboard-kategoriberita' }"
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
import { useKategoriBeritaStore } from "@/stores/kategori-berita";

import { ArrowPathIcon, TagIcon } from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const kategoriBeritaStore = useKategoriBeritaStore();
const toast = useToast();

const isEdit = !!route.params.id;

const loading = ref(false);

const form = reactive({
  id_kategoriberita: null,
  nama_kategoriberita: "",
  urutan: 0,
});

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";

const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";

// =====================================
// Load Data (mode edit)
// =====================================

async function loadData() {
  if (!isEdit) return;

  loading.value = true;

  try {
    const data = await kategoriBeritaStore.fetchKategoriBeritaById(
      route.params.id,
    );

    Object.assign(form, {
      id_kategoriberita: data.id_kategoriberita,
      nama_kategoriberita: data.nama_kategoriberita || "",
      urutan: data.urutan || 0,
    });
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat data kategori berita.",
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

  const nama = form.nama_kategoriberita.trim();

  if (!nama) {
    toast.error("Nama kategori wajib diisi.");
    return;
  }

  if (nama.length > 200) {
    toast.error("Nama kategori maksimal 200 karakter.");
    return;
  }

  loading.value = true;

  try {
    const payload = {
      nama_kategoriberita: nama,
      urutan: form.urutan || 0,
    };

    if (isEdit) {
      const response = await kategoriBeritaStore.updateKategoriBerita(
        form.id_kategoriberita,
        payload,
      );

      toast.success(response.message || "Kategori berita berhasil diperbarui.");
    } else {
      await kategoriBeritaStore.createKategoriBerita(payload);
      toast.success("Kategori berita berhasil ditambahkan.");
    }

    router.push({ name: "dashboard-kategoriberita" });
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menyimpan kategori berita.",
    );
  } finally {
    loading.value = false;
  }
}

onMounted(loadData);
</script>
