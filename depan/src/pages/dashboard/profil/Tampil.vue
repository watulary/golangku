<template>
  <DashboardPageHeader />

  <!-- Hero -->
  <div
    class="relative mb-8 overflow-hidden rounded-2xl bg-white p-6 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5 sm:p-8"
  >
    <div
      class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-[#041983] to-[#f2c808]"
    ></div>

    <h1
      class="text-2xl font-bold tracking-tight text-gray-700 dark:text-slate-200"
    >
      Informasi Akun
    </h1>
    <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
      Kelola informasi akun Anda.
    </p>
  </div>

  <div class="space-y-8">
    <!-- Profil -->

    <div class="grid gap-6 xl:grid-cols-3">
      <!-- Foto -->

      <div
        class="rounded-2xl bg-white p-6 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
      >
        <div class="mb-6 flex items-center gap-3">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/[0.08] text-[#041983] dark:bg-white/5 dark:text-indigo-300"
          >
            <CameraIcon class="h-5 w-5" />
          </div>
          <h2 class="text-lg font-semibold text-gray-700 dark:text-slate-200">
            Foto Profil
          </h2>
        </div>

        <div class="flex flex-col items-center">
          <div class="group relative">
            <img
              :src="authStore.user?.foto_url"
              alt="Foto Profil"
              class="h-32 w-32 cursor-pointer rounded-full border-4 border-gray-100 object-cover transition duration-200 group-hover:opacity-80 dark:border-white/10"
              @click="openLightbox"
            />

            <div
              class="pointer-events-none absolute inset-0 flex items-center justify-center rounded-full bg-black/40 opacity-0 transition duration-200 group-hover:opacity-100"
            >
              <MagnifyingGlassPlusIcon class="h-5 w-5 text-white" />
            </div>

            <div
              v-if="fotoUploading"
              class="absolute inset-0 flex items-center justify-center rounded-full bg-black/50"
            >
              <ArrowPathIcon class="h-6 w-6 animate-spin text-white" />
            </div>

            <button
              type="button"
              title="Ganti Foto"
              :disabled="fotoUploading"
              @click.stop="fileInputRef.click()"
              class="absolute -right-3 -bottom-3 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 disabled:cursor-not-allowed disabled:opacity-60 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700 dark:hover:text-white"
            >
              <PencilSquareIcon class="h-4 w-4" />
            </button>

            <button
              v-if="!isFotoDefault"
              type="button"
              title="Hapus Foto"
              :disabled="fotoUploading"
              @click.stop="showDeleteFotoDialog = true"
              class="absolute -left-3 -bottom-3 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 hover:bg-red-50 hover:text-red-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40 disabled:cursor-not-allowed disabled:opacity-60 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-red-500/10 dark:hover:text-red-400"
            >
              <XMarkIcon class="h-4 w-4" />
            </button>

            <input
              ref="fileInputRef"
              type="file"
              accept=".jpg,.jpeg,.png,.webp"
              class="hidden"
              @change="handleFotoSelected"
            />
          </div>

          <label
            class="mt-2 mb-2 block text-xs font-medium uppercase tracking-wide text-gray-500 dark:text-slate-400"
          >
            Status Password
          </label>

          <span
            v-if="authStore.user?.password_changed"
            class="inline-flex rounded-full bg-emerald-100 px-3 py-1 text-xs font-medium text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-300"
          >
            Diubah
          </span>

          <span
            v-else
            class="inline-flex rounded-full bg-amber-100 px-3 py-1 text-xs font-medium text-amber-700 dark:bg-amber-500/10 dark:text-amber-300"
          >
            Standar
          </span>
        </div>
      </div>

      <!-- Informasi -->

      <div
        class="xl:col-span-2 rounded-2xl bg-white p-6 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
      >
        <div class="mb-6 flex items-center gap-3">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/[0.08] text-[#041983] dark:bg-white/5 dark:text-indigo-300"
          >
            <IdentificationIcon class="h-5 w-5" />
          </div>
          <h2 class="text-lg font-semibold text-gray-700 dark:text-slate-200">
            Identitas Akun
          </h2>
        </div>

        <div class="grid gap-5 md:grid-cols-2">
          <div>
            <label
              class="mb-1 block text-xs font-medium uppercase tracking-wide text-gray-500 dark:text-slate-400"
            >
              Nama
            </label>

            <p class="text-sm text-gray-700 dark:text-slate-200">
              {{ authStore.user?.nama_pengguna }}
            </p>
          </div>

          <div>
            <label
              class="mb-1 block text-xs font-medium uppercase tracking-wide text-gray-500 dark:text-slate-400"
            >
              Email
            </label>

            <p class="text-sm text-gray-700 dark:text-slate-200">
              {{ authStore.user?.email }}
            </p>
          </div>

          <div>
            <label
              class="mb-1 block text-xs font-medium uppercase tracking-wide text-gray-500 dark:text-slate-400"
            >
              Level
            </label>

            <p class="text-sm text-gray-700 dark:text-slate-200">
              {{ authStore.displayRole }}
            </p>
          </div>

          <div>
            <label
              class="mb-1 block text-xs font-medium uppercase tracking-wide text-gray-500 dark:text-slate-400"
            >
              Login Terakhir
            </label>

            <p class="text-sm text-gray-700 dark:text-slate-200">
              {{ lastLogin }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Password -->

    <div
      class="rounded-2xl bg-white p-6 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
    >
      <div class="mb-6 flex items-center gap-3">
        <div
          class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/[0.08] text-[#041983] dark:bg-white/5 dark:text-indigo-300"
        >
          <KeyIcon class="h-5 w-5" />
        </div>
        <h2 class="text-lg font-semibold text-gray-700 dark:text-slate-200">
          Ubah Password
        </h2>
      </div>

      <div class="grid gap-5">
        <div>
          <label
            class="mb-2 block text-sm font-medium text-gray-700 dark:text-slate-200"
          >
            Password Lama
          </label>

          <div class="relative">
            <input
              v-model="form.current_password"
              :type="showPassword.current ? 'text' : 'password'"
              autocomplete="current-password"
              class="w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-100"
            />

            <button
              type="button"
              @click="showPassword.current = !showPassword.current"
              class="absolute inset-y-0 right-3 flex items-center text-gray-400 hover:text-gray-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:hover:text-gray-300"
            >
              <EyeSlashIcon v-if="showPassword.current" class="h-5 w-5" />

              <EyeIcon v-else class="h-5 w-5" />
            </button>
          </div>
        </div>

        <div>
          <label
            class="mb-2 block text-sm font-medium text-gray-700 dark:text-slate-200"
          >
            Password Baru
          </label>

          <div class="relative">
            <input
              v-model="form.new_password"
              :type="showPassword.new ? 'text' : 'password'"
              autocomplete="new-password"
              class="w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-100"
            />

            <button
              type="button"
              @click="showPassword.new = !showPassword.new"
              class="absolute inset-y-0 right-3 flex items-center text-gray-400 hover:text-gray-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:hover:text-gray-300"
            >
              <EyeSlashIcon v-if="showPassword.new" class="h-5 w-5" />

              <EyeIcon v-else class="h-5 w-5" />
            </button>
          </div>
        </div>

        <div
          class="rounded-xl border border-gray-100 bg-gray-50 p-4 dark:border-white/5 dark:bg-[#0B0F1A]"
        >
          <p class="mb-3 text-sm font-medium text-gray-700 dark:text-slate-200">
            Password harus memenuhi syarat berikut:
          </p>

          <ul class="space-y-2 text-sm">
            <li
              class="flex items-center gap-2"
              :class="ruleClass(passwordRules.minLength)"
            >
              <CheckCircleIcon
                v-if="passwordRules.minLength"
                class="h-5 w-5 shrink-0"
              />

              <MinusCircleIcon v-else class="h-5 w-5 shrink-0" />

              <span>Minimal 8 karakter</span>
            </li>

            <li
              class="flex items-center gap-2"
              :class="ruleClass(passwordRules.maxLength)"
            >
              <CheckCircleIcon
                v-if="passwordRules.maxLength"
                class="h-5 w-5 shrink-0"
              />
              <MinusCircleIcon v-else class="h-5 w-5 shrink-0" />
              <span>Maksimal 64 karakter</span>
            </li>

            <li
              class="flex items-center gap-2"
              :class="ruleClass(passwordRules.uppercase)"
            >
              <CheckCircleIcon
                v-if="passwordRules.uppercase"
                class="h-5 w-5 shrink-0"
              />
              <MinusCircleIcon v-else class="h-5 w-5 shrink-0" />
              <span>Mengandung huruf besar (A-Z)</span>
            </li>

            <li
              class="flex items-center gap-2"
              :class="ruleClass(passwordRules.lowercase)"
            >
              <CheckCircleIcon
                v-if="passwordRules.lowercase"
                class="h-5 w-5 shrink-0"
              />
              <MinusCircleIcon v-else class="h-5 w-5 shrink-0" />
              <span>Mengandung huruf kecil (a-z)</span>
            </li>

            <li
              class="flex items-center gap-2"
              :class="ruleClass(passwordRules.number)"
            >
              <CheckCircleIcon
                v-if="passwordRules.number"
                class="h-5 w-5 shrink-0"
              />
              <MinusCircleIcon v-else class="h-5 w-5 shrink-0" />
              <span>Mengandung angka (0-9)</span>
            </li>

            <li
              class="flex items-center gap-2"
              :class="ruleClass(passwordRules.special)"
            >
              <CheckCircleIcon
                v-if="passwordRules.special"
                class="h-5 w-5 shrink-0"
              />
              <MinusCircleIcon v-else class="h-5 w-5 shrink-0" />
              <span>Mengandung karakter spesial</span>
            </li>
          </ul>
        </div>

        <p
          v-if="sameAsCurrentPassword"
          class="mt-3 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-500/20 dark:bg-red-500/10 dark:text-red-300"
        >
          Password baru harus berbeda dari password lama.
        </p>

        <div>
          <label
            class="mb-2 block text-sm font-medium text-gray-700 dark:text-slate-200"
          >
            Konfirmasi Password
          </label>

          <div class="relative">
            <input
              v-model="form.confirm_password"
              :type="showPassword.confirm ? 'text' : 'password'"
              autocomplete="new-password"
              class="w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-100"
            />

            <button
              type="button"
              @click="showPassword.confirm = !showPassword.confirm"
              class="absolute inset-y-0 right-3 flex items-center text-gray-400 hover:text-gray-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:hover:text-gray-300"
            >
              <EyeSlashIcon v-if="showPassword.confirm" class="h-5 w-5" />

              <EyeIcon v-else class="h-5 w-5" />
            </button>
          </div>

          <p
            v-if="form.confirm_password"
            class="mt-3 rounded-lg border px-4 py-3 text-sm"
            :class="
              passwordMatched
                ? 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-500/20 dark:bg-emerald-500/10 dark:text-emerald-300'
                : 'border-red-200 bg-red-50 text-red-700 dark:border-red-500/20 dark:bg-red-500/10 dark:text-red-300'
            "
          >
            {{
              passwordMatched
                ? "Konfirmasi password sudah sesuai."
                : "Konfirmasi password tidak sama."
            }}
          </p>
        </div>

        <div class="pt-2">
          <button
            type="button"
            @click="submitPassword"
            :disabled="!canSubmit || loading"
            class="inline-flex items-center justify-center gap-2 rounded-xl px-5 py-2.5 text-sm font-medium text-white transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 disabled:cursor-not-allowed disabled:bg-gray-400 dark:disabled:bg-slate-700"
            :class="
              canSubmit && !loading
                ? 'bg-blue-600 hover:bg-blue-700'
                : 'bg-gray-400 dark:bg-slate-700'
            "
          >
            <ArrowPathIcon v-if="loading" class="h-5 w-5 animate-spin" />

            {{ loading ? "Menyimpan..." : "Simpan" }}
          </button>
        </div>
      </div>
    </div>
  </div>

  <VueEasyLightbox
    :visible="lightboxVisible"
    :imgs="lightboxImages"
    :index="0"
    @hide="lightboxVisible = false"
  />

  <DashboardConfirmDialog
    v-model="showDeleteFotoDialog"
    title="Hapus Foto Profil"
    message="Foto profil akan dihapus dan dikembalikan ke foto default. Lanjutkan?"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="fotoDeleting"
    @confirm="confirmDeleteFoto"
  />
</template>

<script setup>
import { computed, reactive, ref } from "vue";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";
import { formatDateTimeNumeric } from "@/utils/date";
import { validatePassword, isPasswordValid } from "@/utils/password";
import { compressImage } from "@/utils/image";
import { useToast } from "@/composables/useToast";
import { useAuthStore } from "@/stores/auth";
import { useRouter } from "vue-router";
import VueEasyLightbox from "vue-easy-lightbox";
import {
  EyeIcon,
  EyeSlashIcon,
  CheckCircleIcon,
  MinusCircleIcon,
  ArrowPathIcon,
  MagnifyingGlassPlusIcon,
  PencilSquareIcon,
  XMarkIcon,
  CameraIcon,
  IdentificationIcon,
  KeyIcon,
} from "@heroicons/vue/24/outline";

const authStore = useAuthStore();

const toast = useToast();

const router = useRouter();

const loading = ref(false);

const lastLogin = computed(() => {
  return formatDateTimeNumeric(authStore.user?.last_login);
});

const form = reactive({
  current_password: "",
  new_password: "",
  confirm_password: "",
});

const showPassword = reactive({
  current: false,
  new: false,
  confirm: false,
});

const lightboxVisible = ref(false);

const lightboxImages = computed(() =>
  authStore.user?.foto_url ? [authStore.user.foto_url] : [],
);

function openLightbox() {
  if (!authStore.user?.foto_url) {
    return;
  }

  lightboxVisible.value = true;
}

const FOTO_DEFAULT = "1636812404_b39deec2f158cb789567.jpg";

const fileInputRef = ref(null);

const fotoUploading = ref(false);

const showDeleteFotoDialog = ref(false);

const fotoDeleting = ref(false);

const isFotoDefault = computed(() => {
  const url = authStore.user?.foto_url;

  if (!url) {
    return true;
  }

  return (
    url.endsWith(FOTO_DEFAULT) ||
    url.endsWith("/L.jpg") ||
    url.endsWith("/P.jpg")
  );
});

async function handleFotoSelected(event) {
  const file = event.target.files?.[0];

  event.target.value = "";

  if (!file) {
    return;
  }

  const allowed = ["image/jpeg", "image/png", "image/webp"];

  if (!allowed.includes(file.type)) {
    toast.error("Foto harus berformat JPG, JPEG, PNG atau WEBP.");
    return;
  }

  fotoUploading.value = true;

  try {
    const compressed = await compressImage(file, 100);

    const formData = new FormData();
    formData.append("foto", compressed);

    await authStore.updateFotoProfil(formData);

    toast.success("Foto profil berhasil diubah.");
  } catch (error) {
    toast.error(
      error.response?.data?.message ||
        authStore.error ||
        "Gagal mengubah foto profil.",
    );
  } finally {
    fotoUploading.value = false;
  }
}

async function confirmDeleteFoto() {
  fotoDeleting.value = true;

  try {
    await authStore.deleteFotoProfil();

    toast.success("Foto profil berhasil dihapus.");
    showDeleteFotoDialog.value = false;
  } catch (error) {
    toast.error(
      error.response?.data?.message ||
        authStore.error ||
        "Gagal menghapus foto profil.",
    );
  } finally {
    fotoDeleting.value = false;
  }
}

const sameAsCurrentPassword = computed(() => {
  if (!form.current_password || !form.new_password) {
    return false;
  }

  return form.current_password === form.new_password;
});

const passwordMatched = computed(() => {
  if (!form.confirm_password) {
    return false;
  }

  return form.new_password === form.confirm_password;
});

const canSubmit = computed(() => {
  return (
    form.current_password.trim() !== "" &&
    form.new_password.trim() !== "" &&
    form.confirm_password.trim() !== "" &&
    passwordMatched.value &&
    !sameAsCurrentPassword.value &&
    isPasswordValid(form.new_password)
  );
});

async function submitPassword() {
  if (!canSubmit.value || loading.value) {
    return;
  }

  loading.value = true;

  try {
    const response = await authStore.changePassword({
      current_password: form.current_password,
      new_password: form.new_password,
      confirm_password: form.confirm_password,
    });

    form.current_password = "";
    form.new_password = "";
    form.confirm_password = "";

    toast.success(
      response.message || "Password berhasil diubah. Silakan login kembali.",
      2500,
      () => {
        authStore.logout();
        router.push("/login");
      },
    );
  } catch (error) {
    toast.error(
      error.response?.data?.message ||
        authStore.error ||
        "Gagal mengubah password.",
    );
  } finally {
    loading.value = false;
  }
}

const ruleClass = (passed) => {
  return passed
    ? "text-emerald-600 dark:text-emerald-400"
    : "text-gray-400 dark:text-slate-500";
};

const passwordRules = computed(() => validatePassword(form.new_password));
</script>
