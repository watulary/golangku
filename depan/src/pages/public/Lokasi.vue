<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="bg-white rounded-2xl shadow-lg p-6 md:p-8">
      <div class="self-start">
        <h2 class="text-2xl font-extrabold text-gray-600">Lokasi</h2>

        <div class="mt-2 h-0.5 w-48 flex overflow-hidden rounded-full">
          <div class="w-1/2 bg-[#041983]"></div>
          <div class="w-1/2 bg-[#f2c808]"></div>
        </div>
      </div>
    </div>

    <!-- Kontak + Map -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Informasi -->
      <div class="bg-white rounded-2xl shadow-lg p-6 md:p-6 lg:col-span-1">
        <div class="space-y-6">
          <div class="flex gap-4">
            <MapPinIcon class="w-5 h-5 text-[#f2c808] shrink-0 mt-1" />

            <div>
              <p class="text-[#041983] text-sm mt-1">
                {{ settingStore.setting?.alamat_sekolah }}
              </p>
            </div>
          </div>

          <div class="flex gap-4">
            <PhoneIcon class="w-5 h-5 text-[#f2c808] shrink-0" />

            <div>
              <p class="text-[#041983] text-sm mt-1">
                {{ settingStore.setting?.telp_sekolah }}
              </p>
            </div>
          </div>

          <div class="flex gap-4">
            <EnvelopeIcon class="w-5 h-5 text-[#f2c808] shrink-0" />

            <div>
              <p class="text-[#041983] text-sm mt-1">
                {{ settingStore.setting?.email_sekolah }}
              </p>
            </div>
          </div>

          <div class="flex gap-4">
            <GlobeAltIcon class="w-5 h-5 text-[#f2c808] shrink-0" />

            <div>
              <a
                :href="`${settingStore.setting?.web_sekolah}`"
                rel="noopener noreferrer"
                class="text-[#041983] text-sm mt-1"
              >
                {{ settingStore.setting?.web_sekolah }}
              </a>
            </div>
          </div>
        </div>
      </div>

      <!-- Peta -->
      <div class="bg-white rounded-2xl shadow-lg overflow-hidden lg:col-span-2">
        <iframe
          v-if="settingStore.setting?.peta_sekolah"
          :src="settingStore.setting.peta_sekolah"
          width="100%"
          height="500"
          style="border: 0"
          allowfullscreen
          loading="lazy"
          referrerpolicy="no-referrer-when-downgrade"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from "vue";
import { useSettingStore } from "../../stores/setting";
import {
  MapPinIcon,
  PhoneIcon,
  EnvelopeIcon,
  GlobeAltIcon,
} from "@heroicons/vue/24/outline";

const settingStore = useSettingStore();

onMounted(async () => {
  await settingStore.fetchSetting();
});
</script>
