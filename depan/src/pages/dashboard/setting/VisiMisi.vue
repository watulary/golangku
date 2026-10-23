<template>
  <DashboardPageHeader />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <SectionHeading
        title="Visi & Misi Sekolah"
        description="Judul dan keterangan visi misi yang tampil pada halaman website."
      >
        <DocumentTextIcon class="h-5 w-5" />
      </SectionHeading>

      <div class="space-y-6 p-6">
        <FormField label="Judul">
          <input
            v-model="form.judul"
            type="text"
            :class="inputClass"
            placeholder="Ketik disini..."
          />
        </FormField>

        <FormField label="Isi">
          <TiptapEditor v-model="form.keterangan" />
        </FormField>
      </div>
    </section>

    <div
      class="flex flex-col gap-4 rounded-2xl bg-white px-5 py-4 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5 md:flex-row md:items-center md:justify-between"
    >
      <button
        type="submit"
        :disabled="loading"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-blue-600 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:bg-gray-400 dark:disabled:bg-slate-700"
      >
        <ArrowPathIcon v-if="loading" class="h-5 w-5 animate-spin" />
        {{ loading ? "Menyimpan..." : "Simpan" }}
      </button>

      <div class="space-y-2 md:text-right">
        <div class="flex items-center gap-2 md:justify-end">
          <UserCircleIcon class="h-5 w-5 text-gray-500 dark:text-slate-400" />
          <span class="text-xs text-gray-600 dark:text-slate-400">
            Diperbarui oleh: {{ lastUpdatedBy }}
          </span>
        </div>

        <div class="flex items-center gap-2 md:justify-end">
          <ClockIcon class="h-5 w-5 text-gray-500 dark:text-slate-400" />
          <span class="text-xs text-gray-600 dark:text-slate-400">
            {{ lastUpdatedAt }}
          </span>
        </div>
      </div>
    </div>
  </form>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import FormField from "@/components/form/FormField.vue";
import SectionHeading from "@/components/form/SectionHeading.vue";
import TiptapEditor from "@/components/tiptap/TiptapEditor.vue";
import { useToast } from "@/composables/useToast";
import { useVisiMisiStore } from "@/stores/visimisi";
import { formatDateTime } from "@/utils/date";
import {
  ArrowPathIcon,
  ClockIcon,
  DocumentTextIcon,
  UserCircleIcon,
} from "@heroicons/vue/24/outline";

const visiMisiStore = useVisiMisiStore();
const toast = useToast();

const loading = ref(false);

const form = reactive({
  id: null,
  judul: "",
  keterangan: "",
});

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";

const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";

async function loadData() {
  loading.value = true;

  try {
    await visiMisiStore.getDashboardVisiMisi();

    Object.assign(form, {
      id: visiMisiStore.dashboardVisiMisi.id,
      judul: visiMisiStore.dashboardVisiMisi.judul || "",
      keterangan: visiMisiStore.dashboardVisiMisi.keterangan || "",
    });
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat data visi misi.",
    );
  } finally {
    loading.value = false;
  }
}

async function submit() {
  if (loading.value) return;

  loading.value = true;

  try {
    const response = await visiMisiStore.updateDashboardVisiMisi({
      judul: form.judul,
      keterangan: form.keterangan,
    });

    if (response.changed) {
      toast.success(response.message || "Visi Misi berhasil diperbarui.");
    } else {
      toast.warning(response.message || "Tidak ada perubahan data.");
    }
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menyimpan visi misi.");
  } finally {
    loading.value = false;
  }
}

const lastUpdatedBy = computed(() => {
  return visiMisiStore.dashboardVisiMisi?.nama_pengguna || "Belum diketahui";
});

const lastUpdatedAt = computed(() => {
  if (!visiMisiStore.dashboardVisiMisi?.waktu) {
    return "Belum pernah diperbarui";
  }

  return formatDateTime(visiMisiStore.dashboardVisiMisi.waktu);
});

onMounted(loadData);
</script>
