<template>
  <DashboardPageHeader />
  <DashboardBulkAction
    :selected-count="selectedNis.length"
    @clear="selectedNis = []"
  >
    <template #actions>
      <button
        type="button"
        @click="openBulkMutasi"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-amber-500 px-4 py-2 text-sm font-medium text-white transition hover:bg-amber-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500/40"
      >
        <ArrowsRightLeftIcon class="h-5 w-5" />
        Mutasikan
      </button>
      <button
        type="button"
        @click="showArsipDialog = true"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-purple-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-purple-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-purple-500/40"
      >
        <ArchiveBoxIcon class="h-5 w-5" />
        Arsipkan
      </button>
      <button
        type="button"
        @click="cetakBulk"
        :disabled="bukuIndukStore.cetakLoading"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-sky-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-sky-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-500/40 disabled:cursor-not-allowed disabled:opacity-60"
      >
        <ArrowPathIcon
          v-if="bukuIndukStore.cetakLoading"
          class="h-5 w-5 animate-spin"
        />
        <PrinterIcon v-else class="h-5 w-5" />
        {{ bukuIndukStore.cetakLoading ? "Mencetak..." : "Buku Induk" }}
      </button>
      <button
        type="button"
        @click="showDeleteDialog = true"
        class="inline-flex items-center justify-center gap-2 rounded-xl bg-red-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-red-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500/40"
      >
        <TrashIcon class="h-5 w-5" />
        Hapus
      </button>
    </template>
  </DashboardBulkAction>
  <div class="space-y-6">
    <div
      class="relative overflow-hidden rounded-2xl bg-white ring-1 ring-slate-900/5 shadow-sm dark:bg-slate-900 dark:ring-white/5"
    >
      <div
        class="absolute inset-x-0 top-0 h-1 bg-linear-to-r from-[#041983] to-[#f2c808]"
      ></div>

      <!-- Header -->
      <div class="border-b border-gray-100 px-6 py-5 dark:border-white/5">
        <!-- Baris 1 -->
        <div
          class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between"
        >
          <div class="flex items-center gap-3">
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#041983]/8 text-[#041983] dark:bg-white/5 dark:text-indigo-300"
            >
              <UserGroupIcon class="h-5 w-5" />
            </div>
            <h2 class="text-xl font-semibold text-gray-700 dark:text-slate-200">
              Daftar Peserta Didik Aktif
            </h2>
          </div>
          <div class="relative w-64 max-w-55 sm:max-w-xs md:max-w-sm">
            <MagnifyingGlassIcon
              class="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
            />
            <input
              v-model="pdStore.filters.q"
              @input="searchData"
              type="text"
              placeholder="Pencarian..."
              class="w-full rounded-xl border border-gray-300 bg-white py-2.5 pl-10 pr-4 text-sm outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:placeholder:text-slate-500"
            />
          </div>
        </div>
        <!-- Baris 2 -->
        <div
          class="mt-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
        >
          <!-- kiri -->
          <div class="flex flex-wrap items-center gap-3">
            <div class="flex items-center gap-2">
              <span class="text-sm text-slate-500 dark:text-slate-400">
                Tampilkan
              </span>
              <select
                :value="pdStore.pagination.limit"
                @change="pdStore.changeLimit($event.target.value)"
                class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
              >
                <option :value="10">10</option>
                <option :value="25">25</option>
                <option :value="50">50</option>
                <option :value="100">100</option>
                <option :value="500">500</option>
                <option :value="1000">1000</option>
              </select>
            </div>
            <select
              v-model="kelasFilter"
              @change="changeKelasFilter"
              class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
            >
              <option :value="0">Semua Kelas</option>
              <option
                v-for="k in kelasOptions"
                :key="k.id_kelas"
                :value="k.id_kelas"
              >
                {{ k.nama_kelas }}
              </option>
            </select>
            <select
              v-model="rombelFilter"
              @change="changeRombelFilter"
              class="rounded-xl border border-slate-300 bg-white px-3 py-1.5 text-sm font-medium text-slate-700 shadow-sm transition hover:border-blue-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-slate-600 dark:bg-[#0B0F1A] dark:text-slate-200"
            >
              <option :value="0">Semua Rombel</option>
              <option
                v-for="r in rombelOptions"
                :key="r.id_rombel"
                :value="r.id_rombel"
              >
                {{ r.nama_rombel }}
              </option>
            </select>
          </div>
          <!-- kanan -->
          <div class="flex items-center gap-2">
            <RouterLink
              :to="{ name: 'dashboard-pesertadidik-tambah' }"
              title="Tambah Peserta Didik"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-blue-500 hover:bg-blue-50 hover:text-blue-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-blue-500 dark:hover:bg-blue-500/10 dark:hover:text-blue-400"
            >
              <PlusIcon class="h-5 w-5" />
            </RouterLink>
            <button
              type="button"
              title="Unggah Peserta Didik"
              @click="showUploadDialog = true"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-yellow-500 hover:bg-yellow-50 hover:text-yellow-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-yellow-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-yellow-500 dark:hover:bg-yellow-500/10 dark:hover:text-yellow-300"
            >
              <ArrowUpTrayIcon class="h-5 w-5" />
            </button>
            <button
              type="button"
              title="Unduh Peserta Didik"
              @click="handleExport"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-emerald-500 hover:bg-emerald-50 hover:text-emerald-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-emerald-500 dark:hover:bg-emerald-500/10 dark:hover:text-emerald-400"
            >
              <ArrowDownTrayIcon class="h-5 w-5" />
            </button>
            <RouterLink
              :to="{ name: 'dashboard-pesertadidik-mutasi' }"
              title="Daftar Mutasi"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-amber-500 hover:bg-amber-50 hover:text-amber-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-amber-500 dark:hover:bg-amber-500/10 dark:hover:text-amber-400"
            >
              <ArrowsRightLeftIcon class="h-5 w-5" />
            </RouterLink>
            <RouterLink
              :to="{ name: 'dashboard-pesertadidik-arsip' }"
              title="Daftar Arsip"
              class="rounded-xl border border-slate-300 bg-white p-2.5 text-slate-600 shadow-sm transition hover:border-purple-500 hover:bg-purple-50 hover:text-purple-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-purple-500/40 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-300 dark:hover:border-purple-500 dark:hover:bg-purple-500/10 dark:hover:text-purple-400"
            >
              <ArchiveBoxIcon class="h-5 w-5" />
            </RouterLink>
          </div>
        </div>
      </div>
      <!-- Empty -->
      <div
        v-if="!pdStore.adminLoading && pdStore.adminPesertaDidik.length === 0"
        class="py-16 text-center"
      >
        <div
          class="flex flex-col items-center gap-2 text-gray-400 dark:text-slate-500"
        >
          <MagnifyingGlassIcon v-if="isSearching" class="h-10 w-10" />
          <InboxIcon v-else class="h-10 w-10" />
          <p class="text-sm font-medium">
            {{
              isSearching
                ? "Data tidak ditemukan."
                : "Belum ada data peserta didik."
            }}
          </p>
        </div>
      </div>
      <!-- Table -->
      <div v-else>
        <div class="max-h-[70vh] overflow-auto">
          <table class="min-w-250 w-full">
            <thead class="sticky top-0 z-20 bg-gray-50/90 dark:bg-[#0B0F1A]">
              <tr>
                <th class="w-10 px-4 py-3">
                  <input
                    type="checkbox"
                    :checked="isAllSelected"
                    @change="toggleSelectAll"
                    class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Foto
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Rombel
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="pdStore.toggleSort('nis')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    NIS
                    <component :is="sortIcon('nis')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="pdStore.toggleSort('nisn')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    NISN
                    <component :is="sortIcon('nisn')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="pdStore.toggleSort('nama_siswa')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Nama Peserta Didik
                    <component :is="sortIcon('nama_siswa')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="pdStore.toggleSort('jenis_kelamin')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    L/P
                    <component
                      :is="sortIcon('jenis_kelamin')"
                      class="h-4 w-4"
                    />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="pdStore.toggleSort('agama_siswa')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Agama
                    <component :is="sortIcon('agama_siswa')" class="h-4 w-4" />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="pdStore.toggleSort('tempat_lahirsiswa')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Tempat Lahir
                    <component
                      :is="sortIcon('tempat_lahirsiswa')"
                      class="h-4 w-4"
                    />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  <button
                    type="button"
                    @click="pdStore.toggleSort('tanggal_lahirsiswa')"
                    class="inline-flex items-center gap-1 hover:text-gray-900 dark:hover:text-white"
                  >
                    Tanggal Lahir
                    <component
                      :is="sortIcon('tanggal_lahirsiswa')"
                      class="h-4 w-4"
                    />
                  </button>
                </th>
                <th
                  scope="col"
                  class="px-4 py-3 text-left text-[11px] font-semibold tracking-wider text-gray-500 dark:text-slate-400"
                >
                  Aksi
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-white/5">
              <!-- Skeleton saat memuat -->
              <template v-if="pdStore.adminLoading">
                <tr v-for="n in 6" :key="`skeleton-${n}`">
                  <td v-for="c in 11" :key="c" class="px-4 py-3">
                    <div
                      class="h-4 w-20 animate-pulse rounded bg-gray-100 dark:bg-white/5"
                    ></div>
                  </td>
                </tr>
              </template>

              <tr
                v-else
                v-for="item in pdStore.adminPesertaDidik"
                :key="item.nis"
                class="transition hover:bg-gray-50 dark:hover:bg-[#0B0F1A]/40"
              >
                <td class="px-4 py-2.5">
                  <input
                    type="checkbox"
                    :checked="selectedNis.includes(item.nis)"
                    @change="toggleSelect(item.nis)"
                    class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </td>
                <td class="px-4 py-2.5">
                  <button
                    v-if="item.foto_url"
                    type="button"
                    title="Klik untuk perbesar"
                    class="group block cursor-pointer"
                    @click="openLightbox(item)"
                  >
                    <img
                      :src="item.foto_url"
                      class="h-10 w-10 rounded-full object-cover transition duration-200 group-hover:opacity-80"
                    />
                  </button>
                  <div
                    v-else
                    class="flex h-10 w-10 items-center justify-center rounded-full bg-gray-100 dark:bg-slate-800"
                  >
                    <UserIcon class="h-5 w-5 text-gray-300" />
                  </div>
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.nama_rombel }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.nis }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.nisn }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm font-medium text-gray-800 dark:text-slate-200"
                >
                  <RouterLink
                    :to="{
                      name: 'dashboard-pesertadidik-edit',
                      params: { nis: item.nis },
                    }"
                    class="hover:text-[#041983] hover:underline dark:hover:text-indigo-300"
                  >
                    {{ item.nama_siswa }}
                  </RouterLink>
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.jenis_kelamin }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.agama_siswa }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ item.tempat_lahirsiswa }}
                </td>
                <td
                  class="px-4 py-2.5 text-sm text-gray-600 dark:text-slate-400"
                >
                  {{ formatDateNumeric(item.tanggal_lahirsiswa) }}
                </td>
                <td class="px-4 py-2.5">
                  <button
                    type="button"
                    title="Cetak Buku Induk"
                    @click="cetakSatu(item.nis)"
                    :disabled="bukuIndukStore.cetakLoading"
                    class="rounded-lg border border-slate-300 p-1.5 text-slate-500 transition hover:border-sky-500 hover:bg-sky-50 hover:text-sky-600 disabled:cursor-not-allowed disabled:opacity-50 dark:border-slate-700 dark:text-slate-400 dark:hover:border-sky-500 dark:hover:bg-sky-500/10"
                  >
                    <ArrowPathIcon
                      v-if="cetakingId === item.nis"
                      class="h-4 w-4 animate-spin"
                    />
                    <PrinterIcon v-else class="h-4 w-4" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <DashboardPagination
          :pagination="pdStore.pagination"
          @change-page="changePage"
        />
      </div>
    </div>
  </div>

  <!-- Dialog Delete -->
  <DashboardConfirmDialog
    v-model="showDeleteDialog"
    title="Hapus Peserta Didik"
    :message="`Anda akan menghapus ${selectedNis.length} data peserta didik.\n\nData yang telah dihapus tidak dapat dikembalikan.`"
    confirm-text="Hapus"
    confirm-color="red"
    :loading="loadingAction"
    @confirm="handleDelete"
  />

  <!-- Dialog Arsipkan -->
  <DashboardConfirmDialog
    v-model="showArsipDialog"
    title="Arsipkan Peserta Didik"
    :message="`${selectedNis.length} peserta didik akan dipindahkan ke data arsip.\n\nData tidak akan hilang dari daftar aktif.`"
    confirm-text="Arsipkan"
    confirm-color="purple"
    :loading="loadingAction"
    @confirm="handleArsip"
  />

  <!-- Modal Mutasi -->
  <div
    v-if="showMutasiDialog"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
  >
    <div
      class="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl ring-1 ring-slate-900/5 dark:bg-slate-900 dark:ring-white/10"
    >
      <h3 class="text-lg font-semibold text-gray-800 dark:text-slate-200">
        Mutasi Peserta Didik
      </h3>
      <p class="mt-1 text-sm text-gray-500 dark:text-slate-400">
        {{ selectedNis.length }} peserta didik dipilih
      </p>

      <form @submit.prevent="handleMutasi" class="mt-4 space-y-4">
        <div>
          <label class="text-sm font-medium text-gray-700 dark:text-slate-300">
            Status Keluar
          </label>
          <select
            v-model="mutasiForm.status_keluar"
            required
            class="mt-1 w-full rounded-xl border border-gray-300 px-3 py-2 text-sm outline-none transition focus:border-[#041983] focus:ring-2 focus:ring-[#041983]/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200"
          >
            <option value="">--Pilih--</option>
            <option value="Pindah">Pindah</option>
            <option value="Keluar">Keluar</option>
            <option value="Meninggal">Meninggal</option>
          </select>
        </div>
        <div>
          <label class="text-sm font-medium text-gray-700 dark:text-slate-300">
            Tanggal Keluar
          </label>
          <input
            v-model="mutasiForm.tanggal_keluar"
            type="date"
            required
            class="mt-1 w-full rounded-xl border border-gray-300 px-3 py-2 text-sm outline-none transition scheme-light focus:border-[#041983] focus:ring-2 focus:ring-[#041983]/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200 dark:scheme-dark"
          />
        </div>
        <div>
          <label class="text-sm font-medium text-gray-700 dark:text-slate-300">
            Alasan
          </label>
          <input
            v-model="mutasiForm.alasan_keluar"
            type="text"
            class="mt-1 w-full rounded-xl border border-gray-300 px-3 py-2 text-sm outline-none transition focus:border-[#041983] focus:ring-2 focus:ring-[#041983]/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200"
          />
        </div>
        <div>
          <label class="text-sm font-medium text-gray-700 dark:text-slate-300">
            Sekolah Tujuan
          </label>
          <input
            v-model="mutasiForm.tujuan"
            type="text"
            class="mt-1 w-full rounded-xl border border-gray-300 px-3 py-2 text-sm outline-none transition focus:border-[#041983] focus:ring-2 focus:ring-[#041983]/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200"
          />
        </div>
        <div>
          <label class="text-sm font-medium text-gray-700 dark:text-slate-300">
            Nomor Surat
          </label>
          <input
            v-model="mutasiForm.nomor_surat"
            type="text"
            class="mt-1 w-full rounded-xl border border-gray-300 px-3 py-2 text-sm outline-none transition focus:border-[#041983] focus:ring-2 focus:ring-[#041983]/20 dark:border-slate-700 dark:bg-[#0B0F1A] dark:text-slate-200"
          />
        </div>

        <div class="flex justify-end gap-2 pt-2">
          <button
            type="button"
            @click="showMutasiDialog = false"
            class="rounded-xl border border-gray-300 px-4 py-2 text-sm font-medium text-gray-600 transition hover:bg-gray-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#041983]/40 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
          >
            Batal
          </button>
          <button
            type="submit"
            :disabled="loadingAction"
            class="rounded-xl bg-amber-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-amber-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500/40 disabled:opacity-50"
          >
            Mutasikan
          </button>
        </div>
      </form>
    </div>
  </div>

  <DashboardUploadDialog
    v-model="showUploadDialog"
    title="Unggah Peserta Didik"
    subtitle="Unggah data peserta didik menggunakan Template (.xlsx)"
    :upload-action="pdStore.importPesertaDidik"
    @download-template="handleDownloadTemplate"
    @success="handleImportSuccess"
  />

  <VueEasyLightbox
    :visible="visibleLightbox"
    :imgs="lightboxImages"
    :index="0"
    @hide="visibleLightbox = false"
  />
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { RouterLink } from "vue-router";
import VueEasyLightbox from "vue-easy-lightbox";
import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader.vue";
import DashboardPagination from "@/components/dashboard/DashboardPagination.vue";
import DashboardConfirmDialog from "@/components/dashboard/DashboardConfirmDialog.vue";
import DashboardUploadDialog from "@/components/dashboard/DashboardUploadDialog.vue";
import DashboardBulkAction from "@/components/dashboard/DashboardBulkActionBar.vue";
import { useToast } from "@/composables/useToast";
import { usePesertaDidikStore } from "@/stores/pesertadidik";
import { useKelasStore } from "@/stores/kelas";
import { useRombelStore } from "@/stores/rombel";
import { useBukuIndukStore } from "@/stores/bukuinduk";
import { useTahunPelajaranStore } from "@/stores/tahunpelajaran";
import { formatDateNumeric } from "@/utils/date";
import {
  ArchiveBoxIcon,
  ArrowsRightLeftIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  ChevronUpDownIcon,
  MagnifyingGlassIcon,
  PlusIcon,
  TrashIcon,
  UserIcon,
  ArrowDownTrayIcon,
  ArrowUpTrayIcon,
  UserGroupIcon,
  InboxIcon,
  PrinterIcon,
  ArrowPathIcon,
} from "@heroicons/vue/24/outline";

const pdStore = usePesertaDidikStore();
const kelasStore = useKelasStore();
const rombelStore = useRombelStore();
const toast = useToast();

const showDeleteDialog = ref(false);
const showUploadDialog = ref(false);
const showArsipDialog = ref(false);
const showMutasiDialog = ref(false);
const loadingAction = ref(false);

const bukuIndukStore = useBukuIndukStore();
const tahunPelajaranStore = useTahunPelajaranStore();

const isSearching = computed(() => pdStore.filters.q.trim() !== "");

// ===== Lightbox foto (pola sama spt GaleriFoto.vue) =====
const visibleLightbox = ref(false);
const lightboxImages = ref([]);
function openLightbox(item) {
  if (!item.foto_url) return;
  lightboxImages.value = [item.foto_url];
  visibleLightbox.value = true;
}

const kelasOptions = computed(() =>
  (kelasStore.list || []).filter((k) => k.aktif === 1),
);
const rombelOptions = computed(() =>
  (rombelStore.list || []).filter((r) => r.aktif === 1),
);

const selectedNis = ref([]);

const isAllSelected = computed(
  () =>
    pdStore.adminPesertaDidik.length > 0 &&
    pdStore.adminPesertaDidik.every((item) =>
      selectedNis.value.includes(item.nis),
    ),
);

const resolvedTahunPelajaran = computed(() => {
  const filtered = tahunPelajaranStore.list.find(
    (t) => t.id_tahunpelajaran === pdStore.filters.idTahunPelajaran,
  );
  if (filtered) return filtered.nama_tahunpelajaran;
  const aktif = tahunPelajaranStore.list.find((t) => t.aktif === 1);
  return aktif ? aktif.nama_tahunpelajaran : "";
});

const cetakingId = ref(null);

async function cetakSatu(nis) {
  if (!resolvedTahunPelajaran.value) {
    toast.error(
      "Tidak ada tahun pelajaran aktif. Silakan pilih filter tahun pelajaran terlebih dahulu.",
    );
    return;
  }
  cetakingId.value = nis;
  try {
    await bukuIndukStore.cetakSatu("aktif", nis, resolvedTahunPelajaran.value);
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal mencetak buku induk.");
  } finally {
    cetakingId.value = null;
  }
}

async function cetakBulk() {
  if (!resolvedTahunPelajaran.value) {
    toast.error(
      "Tidak ada tahun pelajaran aktif. Silakan pilih filter tahun pelajaran terlebih dahulu.",
    );
    return;
  }
  try {
    await bukuIndukStore.cetakBulk(
      "aktif",
      selectedNis.value,
      resolvedTahunPelajaran.value,
    );
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal mencetak buku induk massal.",
    );
  }
}

function toggleSelect(nis) {
  const idx = selectedNis.value.indexOf(nis);
  if (idx === -1) {
    selectedNis.value.push(nis);
  } else {
    selectedNis.value.splice(idx, 1);
  }
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedNis.value = selectedNis.value.filter(
      (nis) => !pdStore.adminPesertaDidik.some((item) => item.nis === nis),
    );
  } else {
    const pageNis = pdStore.adminPesertaDidik.map((item) => item.nis);
    selectedNis.value = [...new Set([...selectedNis.value, ...pageNis])];
  }
}

function openBulkMutasi() {
  mutasiForm.status_keluar = "";
  mutasiForm.tanggal_keluar = "";
  mutasiForm.alasan_keluar = "";
  mutasiForm.tujuan = "";
  mutasiForm.nomor_surat = "";
  showMutasiDialog.value = true;
}

const kelasFilter = ref(pdStore.filters.idKelas);
const rombelFilter = ref(pdStore.filters.idRombel);

function changeKelasFilter() {
  pdStore.filters.idKelas = Number(kelasFilter.value);
  pdStore.pagination.page = 1;
  pdStore.fetchAdminPesertaDidik();
}
function changeRombelFilter() {
  pdStore.filters.idRombel = Number(rombelFilter.value);
  pdStore.pagination.page = 1;
  pdStore.fetchAdminPesertaDidik();
}

let searchTimeout = null;
function searchData() {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(async () => {
    pdStore.pagination.page = 1;
    await pdStore.fetchAdminPesertaDidik();
  }, 300);
}

async function changePage(page) {
  pdStore.pagination.page = page;
  await pdStore.fetchAdminPesertaDidik();
}

function sortIcon(column) {
  if (pdStore.sortBy !== column) return ChevronUpDownIcon;
  return pdStore.sortDir === "asc" ? ChevronUpIcon : ChevronDownIcon;
}

// ===== Mutasi =====
const mutasiForm = reactive({
  status_keluar: "",
  tanggal_keluar: "",
  alasan_keluar: "",
  tujuan: "",
  nomor_surat: "",
});

async function handleExport() {
  try {
    await pdStore.exportPesertaDidik();
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal export data peserta didik.",
    );
  }
}

async function handleDownloadTemplate() {
  try {
    await pdStore.downloadTemplate();
  } catch (error) {
    toast.error("Gagal mengunduh template import.");
  }
}

function handleImportSuccess(message) {
  toast.success(message || "Data peserta didik berhasil diimport.");
}

async function handleDelete() {
  loadingAction.value = true;
  try {
    const res = await pdStore.bulkDeletePesertaDidik(selectedNis.value);
    toast.success(res.message || "Data peserta didik berhasil dihapus.");
    selectedNis.value = [];
    showDeleteDialog.value = false;
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal menghapus data peserta didik.",
    );
  } finally {
    loadingAction.value = false;
  }
}

async function handleArsip() {
  loadingAction.value = true;
  try {
    const res = await pdStore.bulkArsipkanPesertaDidik(selectedNis.value);
    const success = res.data?.success ?? [];
    const failed = res.data?.failed ?? {};
    const failedCount = Object.keys(failed).length;

    selectedNis.value = [];
    showArsipDialog.value = false;

    if (failedCount > 0) {
      const sampleReason = Object.values(failed)[0];
      toast.warning(
        `${success.length} berhasil diarsipkan, ${failedCount} gagal. ${sampleReason}`,
      );
    } else {
      toast.success(res.message || "Peserta didik berhasil diarsipkan.");
    }
  } catch (error) {
    toast.error(
      error.response?.data?.message || "Gagal mengarsipkan peserta didik.",
    );
  } finally {
    loadingAction.value = false;
  }
}

async function handleMutasi() {
  loadingAction.value = true;
  try {
    const res = await pdStore.bulkMutasiPesertaDidik(selectedNis.value, {
      ...mutasiForm,
    });
    toast.success(res.message || "Peserta didik berhasil dimutasi.");
    selectedNis.value = [];
    showMutasiDialog.value = false;
  } catch (error) {
    toast.error(error.response?.data?.message || "Gagal memproses mutasi.");
  } finally {
    loadingAction.value = false;
  }
}

onMounted(async () => {
  kelasStore.pagination.limit = 100;
  rombelStore.pagination.limit = 500;
  kelasStore.fetchAll();
  rombelStore.fetchAll();
  tahunPelajaranStore.pagination.limit = 9999; // ambil semua, jangan kepotong pagination
  tahunPelajaranStore.fetchAll();
  pdStore.fetchAdminPesertaDidik();
});
</script>
