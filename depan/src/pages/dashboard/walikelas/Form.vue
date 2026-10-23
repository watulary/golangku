<template>
  <DashboardPageHeader />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Wali Kelas"
        description="Pilih kelas, rombel, dan guru yang menjadi wali kelas."
      >
        <UserGroupIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="space-y-6 p-6">
        <!-- Kelas -->
        <FormField label="Kelas">
          <select
            v-model.number="form.id_kelas"
            @change="onKelasChange"
            :class="inputClass"
          >
            <option :value="0" disabled>--Pilih--</option>
            <option
              v-for="k in kelasOptions"
              :key="k.id_kelas"
              :value="k.id_kelas"
            >
              {{ k.nama_kelas }}
            </option>
          </select>
        </FormField>

        <!-- Rombel -->
        <FormField label="Rombel">
          <select
            v-model.number="form.id_rombel"
            :disabled="!form.id_kelas"
            :class="inputClass"
          >
            <option :value="0" disabled>
              {{ form.id_kelas ? "--Pilih--" : "--Pilih--" }}
            </option>
            <option
              v-for="r in rombelOptions"
              :key="r.id_rombel"
              :value="r.id_rombel"
            >
              {{ r.nama_rombel }}
            </option>
          </select>
        </FormField>

        <!-- Wali Kelas (Pegawai) -->
        <FormField label="Wali Kelas">
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
                v-if="waliKelasStore.pegawaiSearching"
                class="px-4 py-3 text-sm text-gray-500 dark:text-slate-400"
              >
                Mencari...
              </div>

              <div
                v-else-if="waliKelasStore.pegawaiOptions.length === 0"
                class="px-4 py-3 text-sm text-gray-500 dark:text-slate-400"
              >
                Pegawai tidak ditemukan.
              </div>

              <button
                v-for="p in waliKelasStore.pegawaiOptions"
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
            :to="{ name: 'dashboard-walikelas' }"
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
import { useWaliKelasStore } from "@/stores/walikelas";
import { useKelasStore } from "@/stores/kelas";
import { useRombelStore } from "@/stores/rombel";

import { ArrowPathIcon, UserGroupIcon } from "@heroicons/vue/24/outline";

const route = useRoute();
const router = useRouter();
const waliKelasStore = useWaliKelasStore();
const kelasStore = useKelasStore();
const rombelStore = useRombelStore();
const toast = useToast();

const isEdit = !!route.params.id;

const loading = ref(false);

const pegawaiQuery = ref("");
const showPegawaiDropdown = ref(false);

const form = reactive({
  id_walikelas: null,
  id_kelas: 0,
  id_rombel: 0,
  nuptk: "",
  nama_pegawai: "",
});

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";

const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500 disabled:cursor-not-allowed disabled:bg-gray-100 dark:disabled:bg-slate-800";

const kelasOptions = computed(() => {
  const list = kelasStore.list.filter((k) => k.aktif);

  if (form.id_kelas && !list.some((k) => k.id_kelas === form.id_kelas)) {
    const current = kelasStore.list.find((k) => k.id_kelas === form.id_kelas);
    if (current) list.push(current);
  }

  return list;
});

const rombelOptions = computed(() => {
  const list = rombelStore.list.filter(
    (r) => r.id_kelas === form.id_kelas && r.aktif,
  );

  if (form.id_rombel && !list.some((r) => r.id_rombel === form.id_rombel)) {
    const current = rombelStore.list.find(
      (r) => r.id_rombel === form.id_rombel,
    );
    if (current) list.push(current);
  }

  return list;
});

// =====================================
// Kelas & Rombel
// =====================================

async function onKelasChange() {
  form.id_rombel = 0;
}

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
      waliKelasStore.searchPegawai(pegawaiQuery.value.trim());
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
    const data = await waliKelasStore.fetchById(route.params.id);

    Object.assign(form, {
      id_walikelas: data.id_walikelas,
      id_kelas: data.id_kelas,
      id_rombel: data.id_rombel,
      nuptk: data.nuptk,
      nama_pegawai: data.nama_pegawai || "",
    });

    pegawaiQuery.value = data.nama_pegawai || "";
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat data wali kelas.",
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

  if (!form.id_kelas) {
    toast.error("Kelas wajib dipilih.");
    return;
  }

  if (!form.id_rombel) {
    toast.error("Rombel wajib dipilih.");
    return;
  }

  if (!form.nuptk) {
    toast.error("Wali kelas (pegawai) wajib dipilih.");
    return;
  }

  loading.value = true;

  try {
    const payload = {
      id_kelas: form.id_kelas,
      id_rombel: form.id_rombel,
      nuptk: form.nuptk,
    };

    if (isEdit) {
      const response = await waliKelasStore.update(form.id_walikelas, payload);

      if (response.changed) {
        toast.success(response.message || "Wali kelas berhasil diperbarui.");
      } else {
        toast.warning(response.message || "Tidak ada perubahan data.");
      }
    } else {
      await waliKelasStore.create(payload);
      toast.success("Wali kelas berhasil ditambahkan.");
    }

    router.push({ name: "dashboard-walikelas" });
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menyimpan wali kelas.");
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  await kelasStore.changeLimit(999);
  await rombelStore.changeLimit(999);
  await loadData();
});
</script>
