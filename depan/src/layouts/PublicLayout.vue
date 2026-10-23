<template>
  <div class="flex flex-col min-h-screen relative bg-slate-50">
    <PopupPengumuman />

    <!-- HEADER -->
    <header class="fixed top-3 inset-x-0 z-50">
      <div class="max-w-7xl mx-auto px-3 sm:px-6 lg:px-8">
        <div
          class="max-w-7xl mx-auto bg-white/80 backdrop-blur-xl rounded-2xl shadow-lg shadow-slate-900/5 ring-1 ring-slate-900/5 px-3 sm:px-6 lg:px-8"
        >
          <div class="flex justify-between h-14 md:h-16 items-center gap-4">
            <!-- Logo dan Menu -->
            <div class="flex items-center min-w-0 space-x-8">
              <router-link
                to="/"
                class="flex items-center gap-2.5 shrink-0 group"
              >
                <img
                  v-if="logoUrl"
                  :src="logoUrl"
                  :alt="schoolName"
                  class="h-9 w-9 md:h-10 md:w-10 object-contain transition-transform duration-200 group-hover:scale-105"
                />
                <span
                  class="text-base md:text-lg font-bold tracking-tight text-[#041983] truncate"
                  :class="{
                    'hidden sm:block': schoolName.length > 20,
                  }"
                >
                  {{ schoolName }}
                </span>
              </router-link>
              <!-- Desktop Menu -->
              <nav class="hidden md:flex items-center gap-1 text-sm">
                <!-- Profil Sekolah -->
                <div
                  class="relative"
                  @mouseover="openMenu('profil')"
                  @mouseleave="delayedClose"
                >
                  <button
                    class="flex items-center gap-1 px-3 py-2 rounded-lg font-medium text-slate-600 hover:text-[#041983] hover:bg-slate-50 transition-colors duration-200"
                  >
                    <span>Profil</span>
                    <ChevronDownIcon
                      class="w-3.5 h-3.5 transition-transform duration-200"
                      :class="{ 'rotate-180': activeMenu === 'profil' }"
                    />
                  </button>
                  <Transition
                    enter-active-class="transition duration-150 ease-out"
                    enter-from-class="opacity-0 -translate-y-1"
                    enter-to-class="opacity-100 translate-y-0"
                    leave-active-class="transition duration-100 ease-in"
                    leave-from-class="opacity-100 translate-y-0"
                    leave-to-class="opacity-0 -translate-y-1"
                  >
                    <div
                      v-if="activeMenu === 'profil'"
                      class="absolute left-0 mt-2 w-56 rounded-2xl bg-white/95 backdrop-blur-xl ring-1 ring-slate-900/5 shadow-xl shadow-slate-900/10 p-2"
                    >
                      <router-link
                        to="/profil/sejarah"
                        class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                      >
                        Sejarah
                      </router-link>
                      <router-link
                        to="/profil/visi-misi"
                        class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                      >
                        Visi Misi
                      </router-link>
                      <router-link
                        to="/profil/struktur-akademik"
                        class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                      >
                        Struktur Akademik
                      </router-link>
                      <router-link
                        to="/profil/struktur-tatausaha"
                        class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                      >
                        Struktur Tata Usaha
                      </router-link>
                      <router-link
                        to="/profil/struktur-komite"
                        class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                      >
                        Struktur Komite
                      </router-link>
                    </div>
                  </Transition>
                </div>
                <!-- Kegiatan -->
                <div
                  class="relative"
                  @mouseover="openMenu('kegiatan')"
                  @mouseleave="delayedClose"
                >
                  <button
                    class="flex items-center gap-1 px-3 py-2 rounded-lg font-medium text-slate-600 hover:text-[#041983] hover:bg-slate-50 transition-colors duration-200"
                  >
                    <span>Kegiatan</span>
                    <ChevronDownIcon
                      class="w-3.5 h-3.5 transition-transform duration-200"
                      :class="{ 'rotate-180': activeMenu === 'kegiatan' }"
                    />
                  </button>
                  <Transition
                    enter-active-class="transition duration-150 ease-out"
                    enter-from-class="opacity-0 -translate-y-1"
                    enter-to-class="opacity-100 translate-y-0"
                    leave-active-class="transition duration-100 ease-in"
                    leave-from-class="opacity-100 translate-y-0"
                    leave-to-class="opacity-0 -translate-y-1"
                  >
                    <div
                      v-if="activeMenu === 'kegiatan'"
                      class="absolute left-0 mt-2 w-56 rounded-2xl bg-white/95 backdrop-blur-xl ring-1 ring-slate-900/5 shadow-xl shadow-slate-900/10 p-2"
                    >
                      <router-link
                        to="/berita"
                        class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                        >Berita</router-link
                      >
                      <router-link
                        to="/pengumuman"
                        class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                        >Pengumuman</router-link
                      >
                      <router-link
                        to="/foto"
                        class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                        >Galeri Foto</router-link
                      >
                      <router-link
                        to="/video"
                        class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                        >Video</router-link
                      >
                    </div>
                  </Transition>
                </div>
                <!-- PPID -->
                <div
                  class="relative"
                  @mouseover="openMenu('ppid')"
                  @mouseleave="delayedClose"
                >
                  <button
                    class="flex items-center gap-1 px-3 py-2 rounded-lg font-medium text-slate-600 hover:text-[#041983] hover:bg-slate-50 transition-colors duration-200"
                  >
                    <span>PPID</span>
                    <ChevronDownIcon
                      class="w-3.5 h-3.5 transition-transform duration-200"
                      :class="{ 'rotate-180': activeMenu === 'ppid' }"
                    />
                  </button>
                  <Transition
                    enter-active-class="transition duration-150 ease-out"
                    enter-from-class="opacity-0 -translate-y-1"
                    enter-to-class="opacity-100 translate-y-0"
                    leave-active-class="transition duration-100 ease-in"
                    leave-from-class="opacity-100 translate-y-0"
                    leave-to-class="opacity-0 -translate-y-1"
                  >
                    <div
                      v-if="activeMenu === 'ppid'"
                      class="absolute left-0 mt-2 w-60 rounded-2xl bg-white/95 backdrop-blur-xl ring-1 ring-slate-900/5 shadow-xl shadow-slate-900/10 p-2"
                    >
                      <!-- Informasi Publik -->
                      <div class="relative group">
                        <button
                          class="w-full flex items-center justify-between px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                        >
                          <span>Informasi Publik</span>
                          <ChevronRightIcon class="w-3.5 h-3.5" />
                        </button>
                        <!-- Submenu -->
                        <div
                          class="absolute top-0 left-full ml-2 w-60 rounded-2xl bg-white/95 backdrop-blur-xl ring-1 ring-slate-900/5 shadow-xl shadow-slate-900/10 p-2 opacity-0 invisible translate-x-1 group-hover:opacity-100 group-hover:visible group-hover:translate-x-0 transition-all duration-150 ease-out"
                        >
                          <router-link
                            to="/ppid/informasi-dokumen"
                            class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                          >
                            Dokumen Informasi
                          </router-link>
                          <router-link
                            to="/ppid/informasi-setiap-saat"
                            class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                          >
                            Informasi Setiap Saat
                          </router-link>
                          <router-link
                            to="/ppid/informasi-berkala"
                            class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                          >
                            Informasi Berkala
                          </router-link>
                          <router-link
                            to="/ppid/informasi-serta-merta"
                            class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                          >
                            Informasi Serta Merta
                          </router-link>
                        </div>
                      </div>
                      <!-- Kontak -->
                      <router-link
                        to="/ppid/informasi-kontak"
                        class="block px-3.5 py-2.5 rounded-xl text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                      >
                        Kontak
                      </router-link>
                    </div>
                  </Transition>
                </div>
                <router-link
                  to="/spmb"
                  class="px-3 py-2 rounded-lg font-medium text-slate-600 hover:text-[#041983] hover:bg-slate-50 transition-colors duration-200"
                  >SPMB</router-link
                >
                <router-link
                  to="/unduh"
                  class="px-3 py-2 rounded-lg font-medium text-slate-600 hover:text-[#041983] hover:bg-slate-50 transition-colors duration-200"
                  >Unduh</router-link
                >
                <a
                  href="https://lms.smpn254.sch.id"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="px-3 py-2 rounded-lg font-medium text-slate-600 hover:text-[#041983] hover:bg-slate-50 transition-colors duration-200"
                >
                  LMS
                </a>
                <a
                  href="https://smpn254jkt.perpustakaan.co.id/home.ks"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="px-3 py-2 rounded-lg font-medium text-slate-600 hover:text-[#041983] hover:bg-slate-50 transition-colors duration-200"
                >
                  Pustaka Digital
                </a>
                <router-link
                  to="/lokasi"
                  class="px-3 py-2 rounded-lg font-medium text-slate-600 hover:text-[#041983] hover:bg-slate-50 transition-colors duration-200"
                  >Lokasi</router-link
                >
              </nav>
            </div>
            <!-- Tombol Login (Desktop) -->
            <div class="hidden md:flex items-center shrink-0">
              <router-link
                to="/login"
                class="flex items-center gap-2 rounded-full bg-[#041983] pl-4 pr-4.5 py-2 text-sm font-semibold text-white shadow-sm shadow-[#041983]/30 transition-all duration-200 hover:shadow-md hover:shadow-[#041983]/40 hover:-translate-y-0.5"
              >
                <ArrowRightOnRectangleIcon class="h-4 w-4 text-[#f2c808]" />
                Masuk
              </router-link>
            </div>
            <!-- Mobile Menu Button -->
            <div class="md:hidden flex items-center shrink-0">
              <button
                @click="mobileOpen = !mobileOpen"
                class="p-2 rounded-xl text-slate-600 hover:bg-slate-100 active:scale-95 transition"
                :aria-expanded="mobileOpen"
                aria-label="Buka menu"
              >
                <component
                  :is="mobileOpen ? XMarkIcon : Bars3Icon"
                  class="h-6 w-6"
                />
              </button>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- Overlay -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="mobileOpen"
        class="fixed inset-0 bg-slate-900/30 backdrop-blur-[2px] z-55 md:hidden"
        @click="closeMobileMenu"
      ></div>
    </Transition>

    <!-- MOBILE MENU -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div v-if="mobileOpen" class="fixed top-20 inset-x-0 z-60 md:hidden">
        <div class="max-w-7xl mx-auto px-3 sm:px-6 lg:px-8">
          <div
            class="bg-white/95 backdrop-blur-xl rounded-2xl shadow-xl shadow-slate-900/10 ring-1 ring-slate-900/5 overflow-hidden max-h-[75vh] overflow-y-auto"
          >
            <div class="p-3 space-y-1">
              <!-- === PROFIL === -->
              <div>
                <button
                  @click="mobileProfilOpen = !mobileProfilOpen"
                  class="flex items-center justify-between w-full px-3.5 py-3 rounded-xl text-slate-700 hover:bg-slate-50 active:bg-slate-100 transition-colors duration-150"
                >
                  <div class="flex items-center gap-3">
                    <span
                      class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#041983]/8 text-[#041983]"
                    >
                      <BuildingLibraryIcon class="h-4.5 w-4.5" />
                    </span>
                    <span class="font-medium">Profil</span>
                  </div>
                  <component
                    :is="mobileProfilOpen ? ChevronUpIcon : ChevronDownIcon"
                    class="h-4 w-4 text-slate-400"
                  />
                </button>
                <div
                  v-if="mobileProfilOpen"
                  class="ml-4 mt-1 mb-1 space-y-0.5 border-l-2 border-[#f2c808]/60 pl-4"
                >
                  <router-link
                    to="/profil/sejarah"
                    @click="closeMobileMenu"
                    class="flex items-center gap-2 px-3 py-2.5 rounded-lg text-sm text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                  >
                    Sejarah
                  </router-link>
                  <router-link
                    to="/profil/visi-misi"
                    @click="closeMobileMenu"
                    class="flex items-center gap-2 px-3 py-2.5 rounded-lg text-sm text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                  >
                    Visi Misi
                  </router-link>
                  <router-link
                    to="/profil/struktur-akademik"
                    @click="closeMobileMenu"
                    class="flex items-center gap-2 px-3 py-2.5 rounded-lg text-sm text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                  >
                    Struktur Akademik
                  </router-link>
                  <router-link
                    to="/profil/struktur-tatausaha"
                    @click="closeMobileMenu"
                    class="flex items-center gap-2 px-3 py-2.5 rounded-lg text-sm text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                  >
                    Struktur Tata Usaha
                  </router-link>
                  <router-link
                    to="/profil/struktur-komite"
                    @click="closeMobileMenu"
                    class="flex items-center gap-2 px-3 py-2.5 rounded-lg text-sm text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                  >
                    Struktur Komite
                  </router-link>
                </div>
              </div>
              <!-- === KEGIATAN === -->
              <div>
                <button
                  @click="mobileKegiatanOpen = !mobileKegiatanOpen"
                  class="flex items-center justify-between w-full px-3.5 py-3 rounded-xl text-slate-700 hover:bg-slate-50 active:bg-slate-100 transition-colors duration-150"
                >
                  <div class="flex items-center gap-3">
                    <span
                      class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#041983]/8 text-[#041983]"
                    >
                      <FlagIcon class="h-4.5 w-4.5" />
                    </span>
                    <span class="font-medium">Kegiatan</span>
                  </div>
                  <component
                    :is="mobileKegiatanOpen ? ChevronUpIcon : ChevronDownIcon"
                    class="h-4 w-4 text-slate-400"
                  />
                </button>
                <div
                  v-if="mobileKegiatanOpen"
                  class="ml-4 mt-1 mb-1 space-y-0.5 border-l-2 border-[#f2c808]/60 pl-4"
                >
                  <router-link
                    to="/berita"
                    @click="closeMobileMenu"
                    class="flex items-center gap-2 px-3 py-2.5 rounded-lg text-sm text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                  >
                    Berita
                  </router-link>
                  <router-link
                    to="/pengumuman"
                    @click="closeMobileMenu"
                    class="flex items-center gap-2 px-3 py-2.5 rounded-lg text-sm text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                  >
                    Pengumuman
                  </router-link>
                  <router-link
                    to="/foto"
                    @click="closeMobileMenu"
                    class="flex items-center gap-2 px-3 py-2.5 rounded-lg text-sm text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                  >
                    Galeri Foto
                  </router-link>
                  <router-link
                    to="/video"
                    @click="closeMobileMenu"
                    class="flex items-center gap-2 px-3 py-2.5 rounded-lg text-sm text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                  >
                    Video
                  </router-link>
                </div>
              </div>
              <!-- === PPID === -->
              <div>
                <button
                  @click="mobilePPIDOpen = !mobilePPIDOpen"
                  class="flex items-center justify-between w-full px-3.5 py-3 rounded-xl text-slate-700 hover:bg-slate-50 active:bg-slate-100 transition-colors duration-150"
                >
                  <div class="flex items-center gap-3">
                    <span
                      class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#041983]/8 text-[#041983]"
                    >
                      <DocumentTextIcon class="h-4.5 w-4.5" />
                    </span>
                    <span class="font-medium">PPID</span>
                  </div>
                  <component
                    :is="mobilePPIDOpen ? ChevronUpIcon : ChevronDownIcon"
                    class="h-4 w-4 text-slate-400"
                  />
                </button>
                <div
                  v-if="mobilePPIDOpen"
                  class="ml-4 mt-1 mb-1 space-y-0.5 border-l-2 border-[#f2c808]/60 pl-4"
                >
                  <!-- Informasi Publik -->
                  <div>
                    <button
                      @click="
                        mobileInformasiPublikOpen = !mobileInformasiPublikOpen
                      "
                      class="flex items-center justify-between w-full px-3 py-2.5 rounded-lg text-sm text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                    >
                      <span>Informasi Publik</span>
                      <component
                        :is="
                          mobileInformasiPublikOpen
                            ? ChevronUpIcon
                            : ChevronDownIcon
                        "
                        class="h-3.5 w-3.5 text-slate-400"
                      />
                    </button>
                    <div
                      v-if="mobileInformasiPublikOpen"
                      class="ml-3 mt-0.5 space-y-0.5 border-l border-slate-200 pl-3"
                    >
                      <router-link
                        to="/ppid/informasi-dokumen"
                        @click="closeMobileMenu"
                        class="block px-3 py-2 rounded-lg text-sm text-slate-500 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                      >
                        Dokumen Informasi
                      </router-link>
                      <router-link
                        to="/ppid/informasi-setiap-saat"
                        @click="closeMobileMenu"
                        class="block px-3 py-2 rounded-lg text-sm text-slate-500 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                      >
                        Informasi Setiap Saat
                      </router-link>
                      <router-link
                        to="/ppid/informasi-berkala"
                        @click="closeMobileMenu"
                        class="block px-3 py-2 rounded-lg text-sm text-slate-500 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                      >
                        Informasi Berkala
                      </router-link>
                      <router-link
                        to="/ppid/informasi-serta-merta"
                        @click="closeMobileMenu"
                        class="block px-3 py-2 rounded-lg text-sm text-slate-500 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                      >
                        Informasi Serta Merta
                      </router-link>
                    </div>
                  </div>
                  <!-- Kontak -->
                  <router-link
                    to="/ppid/informasi-kontak"
                    @click="closeMobileMenu"
                    class="flex items-center gap-2 px-3 py-2.5 rounded-lg text-sm text-slate-600 hover:bg-slate-50 hover:text-[#041983] transition-colors duration-150"
                  >
                    Kontak
                  </router-link>
                </div>
              </div>
              <!-- MENU BIASA -->
              <div class="pt-1 mt-1 border-t border-slate-100 space-y-0.5">
                <router-link
                  to="/spmb"
                  @click="closeMobileMenu"
                  class="flex items-center gap-3 px-3.5 py-3 rounded-xl text-slate-700 hover:bg-slate-50 active:bg-slate-100 transition-colors duration-150"
                >
                  <span
                    class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#041983]/8 text-[#041983] shrink-0"
                  >
                    <UsersIcon class="h-4.5 w-4.5" />
                  </span>
                  <span class="font-medium">SPMB</span>
                </router-link>
                <router-link
                  to="/unduh"
                  @click="closeMobileMenu"
                  class="flex items-center gap-3 px-3.5 py-3 rounded-xl text-slate-700 hover:bg-slate-50 active:bg-slate-100 transition-colors duration-150"
                >
                  <span
                    class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#041983]/8 text-[#041983] shrink-0"
                  >
                    <ArrowDownTrayIcon class="h-4.5 w-4.5" />
                  </span>
                  <span class="font-medium">Unduh</span>
                </router-link>
                <a
                  href="https://lms.smpn254.sch.id"
                  target="_blank"
                  rel="noopener noreferrer"
                  @click="closeMobileMenu"
                  class="flex items-center gap-3 px-3.5 py-3 rounded-xl text-slate-700 hover:bg-slate-50 active:bg-slate-100 transition-colors duration-150"
                >
                  <span
                    class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#041983]/8 text-[#041983] shrink-0"
                  >
                    <Squares2X2Icon class="h-4.5 w-4.5" />
                  </span>
                  <span class="font-medium">LMS</span>
                </a>
                <a
                  href="https://smpn254jkt.perpustakaan.co.id/home.ks"
                  target="_blank"
                  rel="noopener noreferrer"
                  @click="closeMobileMenu"
                  class="flex items-center gap-3 px-3.5 py-3 rounded-xl text-slate-700 hover:bg-slate-50 active:bg-slate-100 transition-colors duration-150"
                >
                  <span
                    class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#041983]/8 text-[#041983] shrink-0"
                  >
                    <BookOpenIcon class="h-4.5 w-4.5" />
                  </span>
                  <span class="font-medium">Pustaka Digital</span>
                </a>
                <router-link
                  to="/lokasi"
                  @click="closeMobileMenu"
                  class="flex items-center gap-3 px-3.5 py-3 rounded-xl text-slate-700 hover:bg-slate-50 active:bg-slate-100 transition-colors duration-150"
                >
                  <span
                    class="flex h-8 w-8 items-center justify-center rounded-lg bg-[#041983]/8 text-[#041983] shrink-0"
                  >
                    <MapPinIcon class="h-4.5 w-4.5" />
                  </span>
                  <span class="font-medium">Lokasi</span>
                </router-link>
              </div>
              <div class="pt-1 mt-1 border-t border-slate-100">
                <router-link
                  to="/login"
                  @click="closeMobileMenu"
                  class="flex items-center justify-center gap-2 px-3.5 py-3 rounded-xl bg-[#041983] text-white font-semibold shadow-sm shadow-[#041983]/30 transition-transform duration-150 active:scale-[0.98]"
                >
                  <ArrowRightOnRectangleIcon
                    class="h-4.5 w-4.5 text-[#f2c808]"
                  />
                  <span>Masuk</span>
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <!-- MAIN -->
    <main class="flex-1 pt-18 w-full">
      <div class="max-w-7xl mx-auto py-4 px-3 sm:px-6 lg:px-8">
        <router-view />
      </div>
    </main>

    <!-- FOOTER -->
    <footer class="relative pt-2 pb-6">
      <div class="max-w-7xl mx-auto px-3 sm:px-6 lg:px-8">
        <div
          class="relative overflow-hidden rounded-2xl ring-1 ring-slate-900/5 shadow-sm bg-white"
        >
          <div
            class="h-1 w-full bg-linear-to-r from-[#041983] to-[#f2c808]"
          ></div>
          <div class="relative px-5 sm:px-8 py-8">
            <div
              class="grid grid-cols-1 gap-x-8 gap-y-9 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5"
            >
              <!-- Visi -->
              <div>
                <h3
                  class="inline-block border-b-2 border-[#f2c808] pb-1 text-xs font-semibold uppercase tracking-wider text-[#041983]"
                >
                  Visi
                </h3>
                <p class="mt-3 max-w-xs text-sm leading-relaxed text-slate-500">
                  {{ setting.visi_sekolah }}
                </p>
              </div>
              <!-- Misi -->
              <div>
                <h3
                  class="inline-block border-b-2 border-[#f2c808] pb-1 text-xs font-semibold uppercase tracking-wider text-[#041983]"
                >
                  Misi
                </h3>
                <p class="mt-3 max-w-xs text-sm leading-relaxed text-slate-500">
                  {{ setting.deskripsi_sekolah }}
                </p>
              </div>
              <!-- Tautan Cepat -->
              <div>
                <h3
                  class="inline-block border-b-2 border-[#f2c808] pb-1 text-xs font-semibold uppercase tracking-wider text-[#041983]"
                >
                  Tautan Cepat
                </h3>
                <ul class="mt-3 space-y-2.5 text-sm">
                  <li>
                    <router-link
                      to="/berita"
                      class="text-slate-500 hover:text-[#041983] transition-colors duration-150"
                      >Berita</router-link
                    >
                  </li>
                  <li>
                    <router-link
                      to="/pengumuman"
                      class="text-slate-500 hover:text-[#041983] transition-colors duration-150"
                      >Pengumuman</router-link
                    >
                  </li>
                  <li>
                    <router-link
                      to="/foto"
                      class="text-slate-500 hover:text-[#041983] transition-colors duration-150"
                      >Galeri Foto</router-link
                    >
                  </li>
                  <li>
                    <router-link
                      to="/video"
                      class="text-slate-500 hover:text-[#041983] transition-colors duration-150"
                      >Video</router-link
                    >
                  </li>
                  <li>
                    <router-link
                      to="/lokasi"
                      class="text-slate-500 hover:text-[#041983] transition-colors duration-150"
                      >Lokasi</router-link
                    >
                  </li>
                </ul>
              </div>
              <!-- Tautan Lainnya -->
              <div>
                <h3
                  class="inline-block border-b-2 border-[#f2c808] pb-1 text-xs font-semibold uppercase tracking-wider text-[#041983]"
                >
                  Tautan Penting
                </h3>
                <ul class="mt-3 space-y-2.5 text-sm">
                  <li>
                    <a
                      href="https://cbt.smpn254.sch.id/peserta/login"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="flex items-center gap-2 text-slate-500 hover:text-[#041983] transition-colors duration-150"
                    >
                      <span>CBT</span>
                      <span
                        class="inline-flex items-center gap-1 rounded-full bg-green-100 px-2 py-0.5 text-[10px] font-semibold text-green-700"
                      >
                        <span
                          class="h-1.5 w-1.5 rounded-full bg-green-500"
                        ></span>
                        Online
                      </span>
                    </a>
                  </li>
                  <li>
                    <a
                      href="https://smpn254jkt.perpustakaan.co.id/home.ks"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="text-slate-500 hover:text-[#041983] transition-colors duration-150"
                      >Pustaka Digital</a
                    >
                  </li>
                  <li>
                    <a
                      href="https://lms.smpn254.sch.id/login/index.php"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="text-slate-500 hover:text-[#041983] transition-colors duration-150"
                      >LMS</a
                    >
                  </li>
                </ul>
              </div>
              <!-- Tautan Medsos -->
              <div
                v-if="
                  setting.facebook ||
                  setting.instagram ||
                  setting.twitter ||
                  setting.youtube
                "
              >
                <h3
                  class="inline-block border-b-2 border-[#f2c808] pb-1 text-xs font-semibold uppercase tracking-wider text-[#041983]"
                >
                  Tautan Lainnya
                </h3>
                <ul class="mt-3 space-y-2.5 text-sm">
                  <li v-if="setting.instagram">
                    <a
                      :href="setting.instagram"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="text-slate-500 hover:text-[#041983] transition-colors duration-150"
                      >Instagram</a
                    >
                  </li>
                  <li v-if="setting.youtube">
                    <a
                      :href="setting.youtube"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="text-slate-500 hover:text-[#041983] transition-colors duration-150"
                      >YouTube</a
                    >
                  </li>
                  <li v-if="setting.twitter">
                    <a
                      :href="setting.twitter"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="text-slate-500 hover:text-[#041983] transition-colors duration-150"
                      >Twitter</a
                    >
                  </li>
                  <li v-if="setting.facebook">
                    <a
                      :href="setting.facebook"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="text-slate-500 hover:text-[#041983] transition-colors duration-150"
                      >Facebook</a
                    >
                  </li>
                </ul>
              </div>
              <!-- Kontak (tampil otomatis kalau field-nya ada di data setting) -->
              <div v-if="setting.alamat || setting.email || setting.no_telp">
                <h3
                  class="inline-block border-b-2 border-[#f2c808] pb-1 text-xs font-semibold uppercase tracking-wider text-[#041983]"
                >
                  Kontak
                </h3>
                <ul class="mt-3 space-y-2.5 text-sm text-slate-500">
                  <li v-if="setting.alamat" class="flex items-start gap-2">
                    <MapPinIcon
                      class="h-4 w-4 mt-0.5 shrink-0 text-[#f2c808]"
                    />
                    <span>{{ setting.alamat }}</span>
                  </li>
                  <li v-if="setting.no_telp" class="flex items-center gap-2">
                    <PhoneIcon class="h-4 w-4 shrink-0 text-[#f2c808]" />
                    <span>{{ setting.no_telp }}</span>
                  </li>
                  <li v-if="setting.email" class="flex items-center gap-2">
                    <EnvelopeIcon class="h-4 w-4 shrink-0 text-[#f2c808]" />
                    <span>{{ setting.email }}</span>
                  </li>
                </ul>
              </div>
            </div>
            <div
              class="mt-8 flex flex-wrap items-center justify-center gap-1.5 border-t border-slate-100 pt-6 text-xs text-slate-400"
            >
              <span>© {{ year }} With</span>
              <HeartIcon class="w-4 h-4 text-red-500 animate-pulse shrink-0" />
              <span>Tim Tendik</span>
              <span class="text-slate-300">•</span>
              <span>{{ schoolName }}</span>
            </div>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from "vue";
import { useSettingStore } from "../stores/setting";
import PopupPengumuman from "../pages/public/PopupPengumuman.vue";
import {
  ArrowRightOnRectangleIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  ChevronRightIcon,
  DocumentTextIcon,
  ArrowDownTrayIcon,
  MapPinIcon,
  FlagIcon,
  BuildingLibraryIcon,
  UsersIcon,
  EnvelopeIcon,
  PhoneIcon,
  BookOpenIcon,
  Squares2X2Icon,
  Bars3Icon,
  XMarkIcon,
} from "@heroicons/vue/24/outline";
import { HeartIcon } from "@heroicons/vue/24/solid";
const settingStore = useSettingStore();
onMounted(async () => {
  await settingStore.fetchSetting();
});
const schoolName = computed(() => {
  return settingStore.setting?.nama_sekolah || "Sekolah";
});
const logoUrl = computed(() => {
  return settingStore.setting?.logo_url || "";
});
const setting = computed(() => {
  return settingStore.setting || {};
});
// dropdown desktop
const activeMenu = ref(null);
const hideTimeout = ref(null);
function openMenu(menu) {
  clearTimeout(hideTimeout.value);
  activeMenu.value = menu;
}
function delayedClose() {
  clearTimeout(hideTimeout.value);
  hideTimeout.value = setTimeout(() => {
    activeMenu.value = null;
  }, 250);
}
const mobileOpen = ref(false);
const mobileProfilOpen = ref(false);
const mobileKegiatanOpen = ref(false);
const mobilePPIDOpen = ref(false);
const mobileInformasiPublikOpen = ref(false);
const closeMobileMenu = () => {
  mobileOpen.value = false;
  mobileProfilOpen.value = false;
  mobileKegiatanOpen.value = false;
  mobilePPIDOpen.value = false;
  mobileInformasiPublikOpen.value = false;
};
const year = new Date().getFullYear();
</script>
