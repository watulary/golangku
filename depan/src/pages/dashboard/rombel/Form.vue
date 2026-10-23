<template>
  <DashboardPageHeader />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Rombel"
        description="Pilih kelas dan buat nama rombel."
      >
        <UsersIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="space-y-6 p-6">
        <!-- Kelas -->
        <FormField label="Kelas">
          <select v-model.number="form.id_kelas" :class="inputClass">
            <option :value="0" disabled>--Pilih--</option>
            <option
              v-for="k in kelasStore.list"
              :key="k.id_kelas"
              :value="k.id_kelas"
            >
              {{ k.nama_kelas }}
            </option>
          </select>
        </FormField>

        <!-- Nama -->
        <FormField label="Rombel">
          <input
            v-model="form.nama_rombel"
            type="text"
            :class="inputClass"
            placeholder="Contoh: 7A"
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
            :to="{ name: 'dashboard-rombel' }"
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
import { useRombelStore } from "@/stores/rombel";
import { useKelasStore } from "@/stores/kelas";

import { ArrowPathIcon, UsersIcon } from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const rombelStore = useRombelStore();
const kelasStore = useKelasStore();
const toast = useToast();

const isEdit = !!route.params.id;

const loading = ref(false);

const form = reactive({
  id_rombel: null,
  id_kelas: 0,
  nama_rombel: "",
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
    const data = await rombelStore.fetchById(route.params.id);

    Object.assign(form, {
      id_rombel: data.id_rombel,
      id_kelas: data.id_kelas,
      nama_rombel: data.nama_rombel || "",
    });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal memuat data rombel.");
  } finally {
    loading.value = false;
  }
}

// =====================================
// Submit
// =====================================

async function submit() {
  if (loading.value) return;

  const nama = form.nama_rombel.trim();

  if (!nama) {
    toast.error("Nama rombel wajib diisi.");
    return;
  }

  if (nama.length > 50) {
    toast.error("Nama rombel maksimal 50 karakter.");
    return;
  }

  if (!form.id_kelas) {
    toast.error("Kelas wajib dipilih.");
    return;
  }

  loading.value = true;

  try {
    const payload = {
      id_kelas: form.id_kelas,
      nama_rombel: nama,
    };

    if (isEdit) {
      const response = await rombelStore.update(form.id_rombel, payload);

      if (response.changed) {
        toast.success(response.message || "Rombel berhasil diperbarui.");
      } else {
        toast.warning(response.message || "Tidak ada perubahan data.");
      }
    } else {
      await rombelStore.create(payload);
      toast.success("Rombel berhasil ditambahkan.");
    }

    router.push({ name: "dashboard-rombel" });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menyimpan rombel.");
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  await kelasStore.changeLimit(10);
  await loadData();
});
</script>
