<template>
  <DashboardPageHeader />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Mata Pelajaran"
        description="Isilah data mata pelajaran dibawah ini."
      >
        <BookOpenIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="space-y-6 p-6">
        <!-- Nama -->
        <FormField label="Nama Mata Pelajaran">
          <input
            v-model="form.nama_mapel"
            type="text"
            :class="inputClass"
            placeholder="Contoh: Bahasa Indonesia"
            maxlength="100"
          />
        </FormField>

        <!-- Induk -->
        <FormField label="Induk">
          <select v-model="form.id_induk" :class="inputClass">
            <option :value="null">Mapel utama (tanpa induk)</option>
            <option
              v-for="item in indukOptions"
              :key="item.id_mapel"
              :value="item.id_mapel"
            >
              {{ item.nama_mapel }}
            </option>
          </select>

          <p class="mt-1.5 text-xs text-gray-500 dark:text-slate-400">
            Pilih induk jika ini sub-mapel, misalnya "Seni Rupa" di bawah "Seni
            Budaya &amp; Prakarya". Biarkan "Mapel utama" untuk mapel biasa.
          </p>
        </FormField>

        <!-- Urutan -->
        <FormField label="Urutan">
          <input
            v-model.number="form.urutan"
            type="number"
            min="0"
            step="1"
            :class="inputClass"
            placeholder="0"
          />

          <p class="mt-1.5 text-xs text-gray-500 dark:text-slate-400">
            Angka lebih kecil tampil lebih dulu pada daftar dan Buku Induk.
            Sub-mapel diurutkan di dalam induknya.
          </p>
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
            :to="{ name: 'dashboard-mapel' }"
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
import { computed, onMounted, reactive, ref } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";

import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import FormField from "@/components/form/FormField.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";

import { useToast } from "@/composables/useToast";
import { useMapelStore } from "@/stores/mapel";

import { ArrowPathIcon, BookOpenIcon } from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const mapelStore = useMapelStore();
const toast = useToast();

const isEdit = !!route.params.id;

const loading = ref(false);

const form = reactive({
  id_mapel: null,
  nama_mapel: "",
  id_induk: null,
  urutan: 0,
});

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";

const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";

// mapel yang sedang diubah tidak boleh jadi induk bagi dirinya sendiri
const indukOptions = computed(() =>
  mapelStore.induk.filter((item) => item.id_mapel !== form.id_mapel),
);

// =====================================
// Load Data (mode edit)
// =====================================

async function loadData() {
  if (!isEdit) return;

  loading.value = true;

  try {
    const data = await mapelStore.fetchById(route.params.id);

    Object.assign(form, {
      id_mapel: data.id_mapel,
      nama_mapel: data.nama_mapel || "",
      id_induk: data.id_induk ?? null,
      urutan: data.urutan ?? 0,
    });
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat data mata pelajaran.",
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

  const nama = form.nama_mapel.trim();
  const urutan = form.urutan === "" ? 0 : Number(form.urutan);

  if (!nama) {
    toast.error("Nama mata pelajaran wajib diisi.");
    return;
  }

  if (nama.length > 100) {
    toast.error("Nama mata pelajaran maksimal 100 karakter.");
    return;
  }

  if (!Number.isInteger(urutan) || urutan < 0) {
    toast.error("Urutan harus berupa angka bulat 0 atau lebih.");
    return;
  }

  loading.value = true;

  try {
    const payload = {
      nama_mapel: nama,
      id_induk: form.id_induk,
      urutan,
    };

    if (isEdit) {
      const response = await mapelStore.update(form.id_mapel, payload);

      if (response.changed) {
        toast.success(
          response.message || "Mata pelajaran berhasil diperbarui.",
        );
      } else {
        toast.warning(response.message || "Tidak ada perubahan data.");
      }
    } else {
      await mapelStore.create(payload);
      toast.success("Mata pelajaran berhasil ditambahkan.");
    }

    router.push({ name: "dashboard-mapel" });
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menyimpan mata pelajaran.",
    );
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  // dropdown induk harus siap sebelum data edit diisi ke form
  await mapelStore.fetchInduk();
  await loadData();
});
</script>
