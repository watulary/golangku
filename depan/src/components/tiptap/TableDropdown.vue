<template>
  <div class="relative" ref="dropdown">
    <ToolbarButton title="Table" @click="open = !open">
      <TableCellsIcon class="h-5 w-5" />
    </ToolbarButton>

    <div
      v-if="open"
      class="absolute left-0 top-11 z-50 w-60 overflow-hidden rounded-xl border border-gray-200 bg-white shadow-xl ring-1 ring-black/5 dark:border-slate-700 dark:bg-slate-900 dark:ring-white/10"
    >
      <button
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="run(insertTable)"
      >
        Sisipkan Tabel
      </button>

      <div class="h-px bg-gray-200 dark:bg-slate-700"></div>

      <button
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="run(addRowBefore)"
      >
        Tambah Baris Atas
      </button>

      <button
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="run(addRowAfter)"
      >
        Tambah Baris Bawah
      </button>

      <button
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="run(deleteRow)"
      >
        Hapus Baris
      </button>

      <div class="h-px bg-gray-200 dark:bg-slate-700"></div>

      <button
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="run(addColumnBefore)"
      >
        Tambah Kolom Kiri
      </button>

      <button
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="run(addColumnAfter)"
      >
        Tambah Kolom Kanan
      </button>

      <button
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="run(deleteColumn)"
      >
        Hapus Kolom
      </button>

      <div class="h-px bg-gray-200 dark:bg-slate-700"></div>

      <button
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="run(toggleHeaderRow)"
      >
        Header Row
      </button>

      <div class="h-px bg-gray-200 dark:bg-slate-700"></div>

      <button
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="run(deleteTable)"
      >
        Hapus Tabel
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from "vue";
import { TableCellsIcon } from "@heroicons/vue/24/outline";
import ToolbarButton from "./ToolbarButton.vue";

const props = defineProps({
  editor: {
    type: Object,
    required: true,
  },
});

const open = ref(false);
const dropdown = ref(null);

function run(callback) {
  callback();
  open.value = false;
}

function insertTable() {
  props.editor
    .chain()
    .focus()
    .insertTable({
      rows: 3,
      cols: 3,
      withHeaderRow: true,
    })
    .run();
}

function addRowBefore() {
  props.editor.chain().focus().addRowBefore().run();
}

function addRowAfter() {
  props.editor.chain().focus().addRowAfter().run();
}

function deleteRow() {
  props.editor.chain().focus().deleteRow().run();
}

function addColumnBefore() {
  props.editor.chain().focus().addColumnBefore().run();
}

function addColumnAfter() {
  props.editor.chain().focus().addColumnAfter().run();
}

function deleteColumn() {
  props.editor.chain().focus().deleteColumn().run();
}

function toggleHeaderRow() {
  props.editor.chain().focus().toggleHeaderRow().run();
}

function deleteTable() {
  props.editor.chain().focus().deleteTable().run();
}

function handleClickOutside(event) {
  if (dropdown.value && !dropdown.value.contains(event.target)) {
    open.value = false;
  }
}

function mergeCells() {
  props.editor.chain().focus().mergeCells().run();
}

function splitCell() {
  props.editor.chain().focus().splitCell().run();
}

onMounted(() => {
  document.addEventListener("click", handleClickOutside);
});

onBeforeUnmount(() => {
  document.removeEventListener("click", handleClickOutside);
});
</script>
