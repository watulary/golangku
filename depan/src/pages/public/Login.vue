<template>
  <div
    class="relative min-h-screen overflow-hidden bg-slate-50 flex items-center justify-center px-4 py-10"
  >
    <div class="relative w-full max-w-md">
      <!-- Kembali ke Beranda -->
      <div class="mb-4 flex justify-center">
        <RouterLink
          to="/"
          class="inline-flex items-center gap-1.5 text-sm font-medium text-slate-500 hover:text-[#041983] transition-colors duration-150"
        >
          <ArrowLeftIcon class="h-4 w-4" />
          Kembali ke Beranda
        </RouterLink>
      </div>

      <div
        class="overflow-hidden rounded-2xl bg-white shadow-xl shadow-slate-900/5 ring-1 ring-slate-900/5"
      >
        <div
          class="h-1 w-full bg-linear-to-r from-[#041983] to-[#f2c808]"
        ></div>

        <div class="p-7 sm:p-8">
          <!-- Logo -->
          <div class="flex justify-center">
            <div
              class="flex h-20 w-20 items-center justify-center rounded-2xl bg-[#041983]/5"
            >
              <img
                v-if="settingStore.setting?.logo_url"
                :src="settingStore.setting.logo_url"
                :alt="settingStore.setting?.nama_sekolah"
                class="h-14 w-14 object-contain"
              />
            </div>
          </div>

          <!-- Nama Sekolah -->
          <div class="mt-5 text-center">
            <h1
              class="text-xl sm:text-2xl font-bold tracking-tight text-[#041983]"
            >
              {{ settingStore.setting?.nama_sekolah }}
            </h1>
            <p class="mt-1.5 text-sm text-slate-500">
              Silahkan masuk untuk melanjutkan.
            </p>
          </div>

          <!-- Error -->
          <Transition
            enter-active-class="transition duration-150 ease-out"
            enter-from-class="opacity-0 -translate-y-1"
            enter-to-class="opacity-100 translate-y-0"
          >
            <div
              v-if="authStore.error"
              class="mt-6 flex items-start gap-2 rounded-xl border border-red-100 bg-red-50 px-4 py-3 text-sm text-red-700"
            >
              <ExclamationCircleIcon class="h-5 w-5 shrink-0 text-red-500" />
              <span>{{ authStore.error }}</span>
            </div>
          </Transition>

          <!-- Form -->
          <form class="mt-6 space-y-4" @submit.prevent="submit">
            <!-- Email -->
            <div>
              <label class="mb-1.5 block text-xs font-medium text-slate-500"
                >Email</label
              >
              <div class="relative">
                <EnvelopeIcon
                  class="pointer-events-none absolute left-3.5 top-1/2 h-5 w-5 -translate-y-1/2 text-slate-400"
                />
                <input
                  ref="emailInput"
                  v-model="email"
                  type="email"
                  autocomplete="username"
                  placeholder="nama@email.com"
                  class="w-full rounded-xl border border-slate-200 bg-slate-50/60 pl-11 pr-4 py-3 text-sm text-slate-700 placeholder:text-slate-400 transition duration-150 focus:border-[#041983] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#041983]/10"
                />
              </div>
            </div>

            <!-- Password -->
            <div>
              <label class="mb-1.5 block text-xs font-medium text-slate-500"
                >Password</label
              >
              <div class="relative">
                <LockClosedIcon
                  class="pointer-events-none absolute left-3.5 top-1/2 h-5 w-5 -translate-y-1/2 text-slate-400"
                />
                <input
                  v-model="password"
                  :type="showPassword ? 'text' : 'password'"
                  autocomplete="current-password"
                  placeholder="Kata sandi"
                  class="w-full rounded-xl border border-slate-200 bg-slate-50/60 pl-11 pr-12 py-3 text-sm text-slate-700 placeholder:text-slate-400 transition duration-150 focus:border-[#041983] focus:bg-white focus:outline-none focus:ring-4 focus:ring-[#041983]/10"
                />
                <button
                  type="button"
                  tabindex="-1"
                  class="absolute inset-y-0 right-3 flex items-center text-slate-400 hover:text-[#041983] transition-colors duration-150"
                  @click="showPassword = !showPassword"
                >
                  <EyeIcon v-if="!showPassword" class="w-5 h-5" />
                  <EyeSlashIcon v-else class="w-5 h-5" />
                </button>
              </div>
            </div>

            <!-- Tombol -->
            <button
              type="submit"
              :disabled="authStore.loading || !email.trim() || !password"
              class="flex w-full items-center justify-center gap-2 rounded-xl bg-[#041983] py-3 font-semibold text-white shadow-sm shadow-[#041983]/30 transition-all duration-200 hover:shadow-md hover:shadow-[#041983]/40 hover:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-60 disabled:hover:translate-y-0 disabled:hover:shadow-sm"
            >
              <svg
                v-if="authStore.loading"
                class="h-4 w-4 animate-spin text-[#f2c808]"
                viewBox="0 0 24 24"
                fill="none"
              >
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                ></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
                ></path>
              </svg>
              <ArrowRightOnRectangleIcon
                v-else
                class="h-5 w-5 text-[#f2c808]"
              />
              <span v-if="authStore.loading">Proses masuk...</span>
              <span v-else>Masuk</span>
            </button>
          </form>
        </div>
      </div>

      <div
        class="mt-6 flex flex-wrap items-center justify-center gap-1.5 text-xs text-slate-400"
      >
        <span>© {{ year }} With</span>
        <HeartIcon class="w-4 h-4 text-red-500 animate-pulse shrink-0" />
        <span>Tim Tendik</span>
        <span class="text-slate-300">•</span>
        <span>{{ schoolName }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from "vue";
import { useRouter } from "vue-router";

import {
  EyeIcon,
  EyeSlashIcon,
  ArrowRightOnRectangleIcon,
  ArrowLeftIcon,
  EnvelopeIcon,
  LockClosedIcon,
  ExclamationCircleIcon,
} from "@heroicons/vue/24/outline";
import { HeartIcon } from "@heroicons/vue/24/solid";

import { useAuthStore } from "../../stores/auth";
import { useSettingStore } from "../../stores/setting";

const router = useRouter();

const authStore = useAuthStore();
const settingStore = useSettingStore();

const schoolName = computed(() => {
  return settingStore.setting?.nama_sekolah || "Sekolah";
});

const email = ref("");
const emailInput = ref(null);
const password = ref("");
const showPassword = ref(false);

onMounted(async () => {
  await settingStore.fetchSetting();

  // Jika sudah login langsung kembali ke home
  if (authStore.isAuthenticated) {
    await authStore.fetchMe();
    router.replace("/dashboard");
    return;
  }
  await nextTick();

  emailInput.value?.focus();
});

const submit = async () => {
  if (!email.value || !password.value) {
    authStore.error = "Email dan password wajib diisi.";
    return;
  }

  const result = await authStore.login(email.value, password.value);

  if (result.success) {
    // sementara ke home
    router.replace("/dashboard");
  }
};

const year = new Date().getFullYear();
</script>
