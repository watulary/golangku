<template>
  <DashboardPageHeader />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Tahun Pelajaran"
        description="Isilah tahun pelajaran dibawah ini."
      >
        <CalendarDaysIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="space-y-6 p-6">
        <!-- Nama -->
        <FormField label="Tahun Pelajaran">
          <input
            v-model="form.nama_tahunpelajaran"
            type="text"
            :class="inputClass"
            placeholder="2025/2026"
            maxlength="9"
          />

          <p class="mt-1.5 text-xs text-gray-500 dark:text-slate-400">
            Gunakan format
            <code class="rounded bg-gray-100 px-1 py-0.5 dark:bg-slate-800">
              YYYY/YYYY
            </code>
          </p>
        </FormField>

        <!-- Semester -->
        <FormField label="Semester">
          <select v-model="form.semester" :class="inputClass">
            <option value="Ganjil">Ganjil</option>
            <option value="Genap">Genap</option>
          </select>

          <p class="mt-1.5 text-xs text-gray-500 dark:text-slate-400">
            Semester yang sedang berjalan pada tahun pelajaran ini.
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
            :to="{ name: 'dashboard-tahunpelajaran' }"
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
import { useTahunPelajaranStore } from "@/stores/tahunpelajaran";

import { ArrowPathIcon, CalendarDaysIcon } from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const tahunPelajaranStore = useTahunPelajaranStore();
const toast = useToast();

const isEdit = !!route.params.id;

const loading = ref(false);

const tahunAjaranRegex = /^\d{4}\/\d{4}$/;

const form = reactive({
  id_tahunpelajaran: null,
  nama_tahunpelajaran: "",
  semester: "Ganjil",
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
    const data = await tahunPelajaranStore.fetchById(route.params.id);

    Object.assign(form, {
      id_tahunpelajaran: data.id_tahunpelajaran,
      nama_tahunpelajaran: data.nama_tahunpelajaran || "",
      semester: data.semester || "Ganjil",
    });
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat data tahun pelajaran.",
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

  const nama = form.nama_tahunpelajaran.trim();

  if (!nama) {
    toast.error("Tahun pelajaran wajib diisi.");
    return;
  }

  if (!tahunAjaranRegex.test(nama)) {
    toast.error(
      "Format tahun pelajaran tidak valid. Gunakan format YYYY/YYYY.",
    );
    return;
  }

  loading.value = true;

  try {
    const payload = {
      nama_tahunpelajaran: nama,
      semester: form.semester,
    };

    if (isEdit) {
      const response = await tahunPelajaranStore.update(
        form.id_tahunpelajaran,
        payload,
      );

      if (response.changed) {
        toast.success(
          response.message || "Tahun pelajaran berhasil diperbarui.",
        );
      } else {
        toast.warning(response.message || "Tidak ada perubahan data.");
      }
    } else {
      await tahunPelajaranStore.create(payload);
      toast.success("Tahun pelajaran berhasil ditambahkan.");
    }

    router.push({ name: "dashboard-tahunpelajaran" });
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menyimpan tahun pelajaran.",
    );
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  await loadData();
});
</script>
