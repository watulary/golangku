<template>
  <DashboardPageHeader />

  <form class="mx-auto w-full space-y-6" @submit.prevent="submit">
    <!-- Informasi Sekolah -->
    <section :class="cardClass" class="relative">
      <div
        class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <CardHeading
        title="Informasi Sekolah"
        description="Kelola identitas utama website sekolah."
      >
        <BuildingOffice2Icon class="h-5 w-5" />
      </CardHeading>

      <div class="grid gap-6 p-6 lg:grid-cols-[13rem_minmax(0,1fr)]">
        <!-- Logo -->

        <div
          class="flex min-h-56 items-center justify-center rounded-2xl border border-dashed border-gray-300 bg-gray-50 p-5 dark:border-slate-700 dark:bg-[#0B0F1A]"
        >
          <div class="group relative">
            <!-- Logo -->
            <img
              v-if="form.logo_url"
              :src="form.logo_url"
              alt="Logo Sekolah"
              class="h-28 w-28 cursor-pointer rounded-2xl object-contain transition duration-200 group-hover:opacity-80"
              @click="openLightbox"
            />
            <!-- Overlay Zoom -->
            <div
              v-if="form.logo_url"
              class="pointer-events-none absolute inset-0 flex items-center justify-center rounded-2xl bg-black/0 opacity-0 transition duration-200 group-hover:bg-black/30 group-hover:opacity-100"
            >
              <MagnifyingGlassPlusIcon class="h-5 w-5 text-white" />
            </div>

            <BuildingOffice2Icon
              v-else
              class="h-20 w-20 text-gray-300 transition duration-200 group-hover:opacity-80 dark:text-slate-600"
            />

            <!-- Overlay Kompresi -->
            <div
              v-if="processingLogo"
              class="pointer-events-none absolute inset-0 flex items-center justify-center rounded-2xl bg-black/50"
            >
              <ArrowPathIcon class="h-6 w-6 animate-spin text-white" />
            </div>

            <!-- Edit -->

            <button
              type="button"
              title="Ganti Logo"
              :disabled="processingLogo"
              @click="logoInput.click()"
              class="absolute -right-2 -bottom-2 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-slate-700 dark:hover:text-white"
            >
              <PencilSquareIcon class="h-4 w-4" />
            </button>

            <!-- Hapus -->

            <button
              v-if="form.logo_url && !isDefaultLogo"
              type="button"
              title="Hapus Logo"
              :disabled="processingLogo"
              @click="removeLogo"
              class="absolute -left-2 -bottom-2 flex h-8 w-8 scale-90 items-center justify-center rounded-lg border border-gray-200 bg-white text-gray-600 opacity-100 shadow-lg transition-all duration-200 md:opacity-0 md:group-hover:scale-100 md:group-hover:opacity-100 hover:bg-red-50 hover:text-red-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40 disabled:cursor-not-allowed disabled:opacity-60 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:bg-red-500/10 dark:hover:text-red-400"
            >
              <XMarkIcon class="h-4 w-4" />
            </button>

            <!-- Input -->

            <input
              ref="logoInput"
              type="file"
              accept=".jpg,.jpeg,.png,.webp"
              class="hidden"
              @change="handleLogoChange"
            />
          </div>
        </div>

        <!-- Identitas -->
        <div class="space-y-5">
          <FormField label="Nama Sekolah">
            <input
              v-model="form.nama_sekolah"
              type="text"
              :class="inputClass"
            />
          </FormField>

          <div class="grid gap-5 md:grid-cols-3">
            <FormField label="NPSN">
              <input v-model="form.npsn" type="text" :class="inputClass" />
            </FormField>

            <FormField label="NSSS">
              <input v-model="form.nsss" type="text" :class="inputClass" />
            </FormField>

            <FormField label="NSP">
              <input v-model="form.nsp" type="text" :class="inputClass" />
            </FormField>
          </div>
        </div>
      </div>
    </section>

    <!-- Lokasi dan Kontak -->
    <div class="grid gap-6 xl:grid-cols-12">
      <section :class="[cardClass, 'xl:col-span-8']">
        <CardHeading
          title="Lokasi & Kontak"
          description="Lengkapi alamat dan kontak sekolah."
        >
          <MapPinIcon class="h-5 w-5" />
        </CardHeading>

        <div class="grid gap-5 p-6 md:grid-cols-2">
          <FormField class="md:col-span-2" label="Alamat">
            <textarea
              v-model="form.alamat_sekolah"
              rows="3"
              :class="textareaClass"
            />
          </FormField>

          <FormField label="Kelurahan / Desa">
            <input
              v-model="form.kelurahan_sekolah"
              type="text"
              :class="inputClass"
            />
          </FormField>

          <FormField label="Kecamatan">
            <input
              v-model="form.kecamatan_sekolah"
              type="text"
              :class="inputClass"
            />
          </FormField>

          <FormField label="Kabupaten / Kota">
            <input
              v-model="form.kabupaten_sekolah"
              type="text"
              :class="inputClass"
            />
          </FormField>

          <FormField label="Provinsi">
            <input
              v-model="form.provinsi_sekolah"
              type="text"
              :class="inputClass"
            />
          </FormField>

          <FormField label="Kode Pos">
            <input
              v-model="form.kodepos_sekolah"
              type="text"
              inputmode="numeric"
              :class="inputClass"
            />
          </FormField>

          <FormField label="Telepon">
            <input v-model="form.telp_sekolah" type="tel" :class="inputClass" />
          </FormField>

          <FormField label="Email">
            <input
              v-model="form.email_sekolah"
              type="email"
              :class="inputClass"
            />
          </FormField>

          <FormField label="Website">
            <input
              v-model="form.web_sekolah"
              type="url"
              placeholder="https://..."
              :class="inputClass"
            />
          </FormField>
        </div>
      </section>

      <!-- Media Sosial -->
      <section :class="[cardClass, 'xl:col-span-4']">
        <CardHeading title="Media Sosial" description="Kanal resmi sekolah.">
          <GlobeAltIcon class="h-5 w-5" />
        </CardHeading>

        <div class="space-y-5 p-6">
          <FormField label="Facebook">
            <input
              v-model="form.facebook"
              type="url"
              placeholder="https://..."
              :class="inputClass"
            />
          </FormField>

          <FormField label="Instagram">
            <input
              v-model="form.instagram"
              type="url"
              placeholder="https://..."
              :class="inputClass"
            />
          </FormField>

          <FormField label="Twitter">
            <input
              v-model="form.twitter"
              type="url"
              placeholder="https://..."
              :class="inputClass"
            />
          </FormField>

          <FormField label="YouTube">
            <input
              v-model="form.youtube"
              type="url"
              placeholder="https://..."
              :class="inputClass"
            />
          </FormField>
        </div>
      </section>
    </div>

    <!-- Konten Website -->
    <section :class="cardClass">
      <CardHeading
        title="Konten Website"
        description="Konten singkat yang tampil pada halaman publik."
      >
        <DocumentTextIcon class="h-5 w-5" />
      </CardHeading>

      <div class="grid gap-5 p-6 md:grid-cols-2">
        <FormField label="Deskripsi Sekolah">
          <textarea
            v-model="form.deskripsi_sekolah"
            rows="5"
            :class="textareaClass"
          />
        </FormField>

        <FormField label="Visi Sekolah">
          <textarea
            v-model="form.visi_sekolah"
            rows="5"
            :class="textareaClass"
          />
        </FormField>

        <FormField class="md:col-span-2" label="Profil Sekolah">
          <textarea
            v-model="form.profil_sekolah"
            rows="5"
            :class="textareaClass"
          />
        </FormField>

        <FormField class="md:col-span-2" label="Footer Website">
          <textarea
            v-model="form.footer_sekolah"
            rows="3"
            :class="textareaClass"
          />
        </FormField>
      </div>
    </section>

    <!-- Fitur Publik -->
    <div class="grid gap-6 xl:grid-cols-2">
      <!-- Google Maps -->
      <section :class="cardClass">
        <CardHeading
          title="Google Maps"
          description="Lokasi sekolah pada halaman publik."
        >
          <MapIcon class="h-5 w-5" />
        </CardHeading>

        <div class="space-y-5 p-6">
          <FormField label="Google Maps Embed">
            <textarea
              v-model="form.peta_sekolah"
              rows="4"
              placeholder="https://www.google.com/maps/embed?... atau kode iframe"
              :class="textareaClass"
            />
          </FormField>

          <div
            v-if="mapEmbedUrl"
            class="overflow-hidden rounded-xl border border-gray-200 dark:border-white/10"
          >
            <iframe
              :src="mapEmbedUrl"
              title="Pratinjau lokasi sekolah"
              class="h-72 w-full border-0"
              loading="lazy"
              referrerpolicy="no-referrer-when-downgrade"
            />
          </div>

          <div
            v-else
            class="flex h-72 items-center justify-center rounded-xl border border-dashed border-gray-300 bg-gray-50 p-6 text-center text-sm text-gray-500 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-400"
          >
            Pratinjau peta akan muncul setelah URL embed Google Maps yang valid
            diisi.
          </div>
        </div>
      </section>

      <!-- Popup Pengumuman -->
      <section :class="cardClass">
        <CardHeading
          title="Popup Pengumuman"
          description="Pengumuman singkat pada halaman publik."
        >
          <MegaphoneIcon class="h-5 w-5" />
        </CardHeading>

        <div class="space-y-5 p-6">
          <label
            class="flex cursor-pointer items-start gap-3 rounded-xl border border-gray-200 bg-gray-50 p-4 dark:border-white/10 dark:bg-[#0B0F1A]"
          >
            <input
              v-model="popupEnabled"
              type="checkbox"
              class="mt-0.5 h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />

            <span>
              <span
                class="block text-sm font-medium text-gray-700 dark:text-slate-200"
              >
                Tampilkan popup pengumuman
              </span>

              <span
                class="mt-1 block text-xs leading-relaxed text-gray-500 dark:text-slate-400"
              >
                Aktifkan untuk menampilkan pengumuman pada halaman publik.
              </span>
            </span>
          </label>

          <FormField label="Judul Popup">
            <input v-model="form.heroarea2" type="text" :class="inputClass" />
          </FormField>

          <FormField label="Deskripsi Popup">
            <textarea
              v-model="form.heroarea3"
              rows="4"
              :class="textareaClass"
            />
          </FormField>

          <FormField label="Tautan Popup">
            <input
              v-model="form.heroarea4"
              type="url"
              placeholder="https://..."
              :class="inputClass"
            />
          </FormField>
        </div>
      </section>
    </div>

    <!-- Aksi -->
    <div
      class="flex flex-col gap-4 rounded-2xl bg-white px-5 py-4 ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5 md:flex-row md:items-center md:justify-between"
    >
      <!-- Tombol -->
      <button
        type="submit"
        :disabled="loading"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-blue-600 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-blue-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 disabled:cursor-not-allowed disabled:bg-gray-400 dark:disabled:bg-slate-700"
      >
        <ArrowPathIcon v-if="loading" class="h-5 w-5 animate-spin" />

        {{ loading ? "Menyimpan..." : "Simpan" }}
      </button>

      <!-- Informasi -->
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
  <VueEasyLightbox
    :visible="lightboxVisible"
    :imgs="lightboxImg"
    :index="0"
    @hide="lightboxVisible = false"
  />
</template>

<script setup>
import { computed, defineComponent, h, onMounted, reactive, ref } from "vue";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import { useToast } from "@/composables/useToast";
import { useSettingStore } from "@/stores/setting";
import { formatDateTime } from "@/utils/date";
import { compressImage } from "@/utils/image";
import VueEasyLightbox from "vue-easy-lightbox";
import {
  ArrowPathIcon,
  BuildingOffice2Icon,
  DocumentTextIcon,
  GlobeAltIcon,
  MapIcon,
  MapPinIcon,
  MegaphoneIcon,
  PhoneIcon,
  PencilSquareIcon,
  UserCircleIcon,
  ClockIcon,
  MagnifyingGlassPlusIcon,
  XMarkIcon,
} from "@heroicons/vue/24/outline";

const settingStore = useSettingStore();
const toast = useToast();
const loading = ref(false);
const selectedLogo = ref(null);
const hapusLogo = ref(false);
const logoBase = ref(""); // alamat folder logo di server, diambil dari logo_url
const processingLogo = ref(false);
const logoInput = ref(null);
const lightboxVisible = ref(false);
const lightboxImg = ref("");
const form = reactive({
  id_sekolah: null,
  nama_sekolah: "",
  alamat_sekolah: "",
  npsn: "",
  nsss: "",
  nsp: "",
  kelurahan_sekolah: "",
  kecamatan_sekolah: "",
  kabupaten_sekolah: "",
  provinsi_sekolah: "",
  kodepos_sekolah: "",
  telp_sekolah: "",
  email_sekolah: "",
  web_sekolah: "",
  deskripsi_sekolah: "",
  visi_sekolah: "",
  footer_sekolah: "",
  profil_sekolah: "",
  facebook: "",
  twitter: "",
  instagram: "",
  youtube: "",
  peta_sekolah: "",
  logo_sekolah: "",
  logo_url: "",
  heroarea1: "0",
  heroarea2: "",
  heroarea3: "",
  heroarea4: "",
});

const cardClass =
  "overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5";
const inputClass =
  "w-full rounded-xl border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-700 outline-none transition placeholder:text-gray-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500";
const textareaClass = `${inputClass} resize-y`;

// Harus sama dengan defaultLogo di handler.go
const DEFAULT_LOGO = "1636812404_b39deec2f158cb789567.png";

const isDefaultLogo = computed(() => form.logo_url.endsWith(DEFAULT_LOGO));

const popupEnabled = computed({
  get: () =>
    ["1", "true", "yes"].includes(String(form.heroarea1).toLowerCase()),
  set: (enabled) => {
    form.heroarea1 = enabled ? "1" : "0";
  },
});

const mapEmbedUrl = computed(() => {
  const value = form.peta_sekolah.trim();
  if (!value) return "";

  const iframeSource = value.match(/<iframe[^>]+src=["']([^"']+)["']/i)?.[1];
  const url = iframeSource || value;

  try {
    const parsed = new URL(url);
    return parsed.protocol === "https:" &&
      parsed.hostname.endsWith("google.com")
      ? parsed.href
      : "";
  } catch {
    return "";
  }
});

async function loadData() {
  loading.value = true;

  try {
    await settingStore.getDashboardSetting();
    Object.assign(form, settingStore.dashboardSetting);
    logoBase.value = (form.logo_url || "").replace(/[^/]+$/, "");
    selectedLogo.value = null;
    hapusLogo.value = false;
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal memuat data pengaturan.",
    );
  } finally {
    loading.value = false;
  }
}

// =====================================
// Upload Logo
// =====================================

async function handleLogoChange(event) {
  const file = event.target.files?.[0];

  event.target.value = "";

  if (!file) {
    return;
  }

  // =====================================
  // Validasi Format
  // =====================================

  const allowed = ["image/jpeg", "image/png", "image/webp"];

  if (!allowed.includes(file.type)) {
    toast.error("Logo harus berformat JPG, JPEG, PNG atau WEBP.");
    return;
  }

  // =====================================
  // Kompres Otomatis (maks 100 KB)
  // =====================================

  processingLogo.value = true;

  try {
    const compressed = await compressImage(file, 100);

    selectedLogo.value = compressed;
    hapusLogo.value = false;

    // Preview
    form.logo_url = URL.createObjectURL(compressed);
  } catch {
    toast.error("Gagal memproses logo.");
  } finally {
    processingLogo.value = false;
  }
}

// =====================================
// Hapus Logo (diterapkan saat tombol Simpan ditekan)
// =====================================

function removeLogo() {
  selectedLogo.value = null;
  hapusLogo.value = true;
  form.logo_url = logoBase.value ? logoBase.value + DEFAULT_LOGO : "";
  form.logo_sekolah = DEFAULT_LOGO;
}

async function submit() {
  if (loading.value) return;

  loading.value = true;

  try {
    // =====================================
    // Form Data
    // =====================================

    const data = new FormData();

    Object.entries(form).forEach(([key, value]) => {
      if (key !== "logo_url") {
        data.append(key, value ?? "");
      }
    });

    if (selectedLogo.value) {
      data.append("logo", selectedLogo.value);
    }

    if (hapusLogo.value) {
      data.append("hapus_logo", "1");
    }

    // =====================================
    // Submit
    // =====================================

    const response = await settingStore.updateDashboardSetting(data);

    // Logo sudah tersimpan/terhapus di server; jangan kirim ulang di simpan berikutnya.
    selectedLogo.value = null;
    hapusLogo.value = false;

    // Samakan tampilan dengan data server (mis. logo default setelah dihapus).
    if (response.data?.logo_url) {
      form.logo_url = response.data.logo_url;
      form.logo_sekolah = response.data.logo_sekolah || "";
      logoBase.value = form.logo_url.replace(/[^/]+$/, "");
    }

    // Muat ulang setting publik supaya favicon (dan logo di header) langsung
    // mengikuti perubahan tanpa refresh halaman.
    if (response.changed) {
      settingStore.loaded = false;
      await settingStore.fetchSetting();
    }

    if (response.changed) {
      toast.success(response.message || "Website berhasil diperbarui.");
    } else {
      toast.warning(response.message || "Tidak ada perubahan data.");
    }
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal menyimpan pengaturan.");
  } finally {
    loading.value = false;
  }
}

const FormField = defineComponent({
  name: "FormField",
  props: {
    label: { type: String, required: true },
    required: Boolean,
  },
  setup(props, { slots, attrs }) {
    return () =>
      h("div", { ...attrs }, [
        h(
          "label",
          {
            class:
              "mb-2 block text-sm font-medium text-gray-700 dark:text-slate-300",
          },
          [
            props.label,
            props.required
              ? h("span", { class: "ml-1 text-red-500" }, "*")
              : null,
          ],
        ),
        slots.default?.(),
      ]);
  },
});

const CardHeading = defineComponent({
  name: "CardHeading",
  props: {
    title: { type: String, required: true },
    description: { type: String, required: true },
  },
  setup(props, { slots }) {
    return () =>
      h(
        "div",
        {
          class:
            "flex items-start gap-3 border-b border-gray-100 px-6 py-5 dark:border-white/5",
        },
        [
          h(
            "span",
            {
              class:
                "mt-0.5 flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-blue-100 text-blue-600 dark:bg-blue-500/10 dark:text-blue-300",
            },
            slots.default?.(),
          ),
          h("div", [
            h(
              "h2",
              {
                class:
                  "text-lg font-semibold text-gray-700 dark:text-slate-200",
              },
              props.title,
            ),
            h(
              "p",
              { class: "mt-1 text-sm text-gray-500 dark:text-slate-400" },
              props.description,
            ),
          ]),
        ],
      );
  },
});

const lastUpdatedBy = computed(() => {
  return settingStore.dashboardSetting?.updated_by_nama || "Belum diketahui";
});

const lastUpdatedAt = computed(() => {
  if (!settingStore.dashboardSetting?.updated_at) {
    return "Belum pernah diperbarui";
  }

  return formatDateTime(settingStore.dashboardSetting.updated_at);
});

function openLightbox() {
  if (!form.logo_url) return;
  lightboxImg.value = form.logo_url;
  lightboxVisible.value = true;
}

onMounted(loadData);
</script>
