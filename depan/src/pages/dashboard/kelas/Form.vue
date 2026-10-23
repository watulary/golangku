<template>
  <DashboardPageHeader />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Kelas"
        description="Isilah nama kelas di bawah ini."
      >
        <RectangleStackIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="space-y-6 p-6">
        <!-- Nama -->
        <FormField label="Kelas">
          <input
            v-model="form.nama_kelas"
            type="text"
            :class="inputClass"
            placeholder="Contoh: 7"
            maxlength="50"
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
            :to="{ name: 'dashboard-kelas' }"
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
import { useKelasStore } from "@/stores/kelas";

import { ArrowPathIcon, RectangleStackIcon } from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const kelasStore = useKelasStore();
const toast = useToast();

const isEdit = !!route.params.id;

const loading = ref(false);

const form = reactive({
  id_kelas: null,
  nama_kelas: "",
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
    const data = await kelasStore.fetchById(route.params.id);

    Object.assign(form, {
      id_kelas: data.id_kelas,
      nama_kelas: data.nama_kelas || "",
    });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal memuat data kelas.");
  } finally {
    loading.value = false;
  }
}

// =====================================
// Submit
// =====================================

async function submit() {
  if (loading.value) return;

  const nama = form.nama_kelas.trim();

  if (!nama) {
    toast.error("Nama kelas wajib diisi.");
    return;
  }

  if (nama.length > 50) {
    toast.error("Nama kelas maksimal 50 karakter.");
    return;
  }

  loading.value = true;

  try {
    const payload = {
      nama_kelas: nama,
    };

    if (isEdit) {
      const response = await kelasStore.update(form.id_kelas, payload);

      if (response.changed) {
        toast.success(response.message || "Kelas berhasil diperbarui.");
      } else {
        toast.warning(response.message || "Tidak ada perubahan data.");
      }
    } else {
      await kelasStore.create(payload);
      toast.success("Kelas berhasil ditambahkan.");
    }

    router.push({ name: "dashboard-kelas" });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menyimpan kelas.");
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  await loadData();
});
</script>
