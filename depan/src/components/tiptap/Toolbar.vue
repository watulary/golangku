<template>
  <div
    class="flex flex-wrap items-center gap-2 rounded-t-xl border border-gray-200 bg-gray-50 p-2 dark:border-slate-700 dark:bg-slate-900"
  >
    <!-- Undo / Redo -->
    <div class="flex items-center gap-1">
      <ToolbarButton
        title="Undo"
        :disabled="!editor.can().chain().focus().undo().run()"
        @click="editor.chain().focus().undo().run()"
      >
        <ArrowUturnLeftIcon class="h-5 w-5" />
      </ToolbarButton>

      <ToolbarButton
        title="Redo"
        :disabled="!editor.can().chain().focus().redo().run()"
        @click="editor.chain().focus().redo().run()"
      >
        <ArrowUturnRightIcon class="h-5 w-5" />
      </ToolbarButton>
    </div>

    <div class="h-6 w-px bg-gray-300 dark:bg-slate-700" />

    <!-- Heading -->
    <div class="flex items-center gap-1">
      <ToolbarButton
        title="Heading 1"
        :active="editor.isActive('heading', { level: 1 })"
        @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
      >
        H1
      </ToolbarButton>

      <ToolbarButton
        title="Heading 2"
        :active="editor.isActive('heading', { level: 2 })"
        @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
      >
        H2
      </ToolbarButton>

      <ToolbarButton
        title="Heading 3"
        :active="editor.isActive('heading', { level: 3 })"
        @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
      >
        H3
      </ToolbarButton>
    </div>

    <div class="h-6 w-px bg-gray-300 dark:bg-slate-700" />

    <!-- Text -->
    <div class="flex items-center gap-1">
      <ToolbarButton
        title="Bold"
        :active="editor.isActive('bold')"
        @click="editor.chain().focus().toggleBold().run()"
      >
        <strong>B</strong>
      </ToolbarButton>

      <ToolbarButton
        title="Italic"
        :active="editor.isActive('italic')"
        @click="editor.chain().focus().toggleItalic().run()"
      >
        <em>I</em>
      </ToolbarButton>

      <ToolbarButton
        title="Underline"
        :active="editor.isActive('underline')"
        @click="editor.chain().focus().toggleUnderline().run()"
      >
        <span class="underline">U</span>
      </ToolbarButton>

      <ToolbarButton
        title="Strike"
        :active="editor.isActive('strike')"
        @click="editor.chain().focus().toggleStrike().run()"
      >
        <span class="line-through">S</span>
      </ToolbarButton>
    </div>

    <div class="h-6 w-px bg-gray-300 dark:bg-slate-700" />

    <!-- Align -->
    <div class="flex items-center gap-1">
      <ToolbarButton
        title="Align Left"
        :active="editor.isActive({ textAlign: 'left' })"
        @click="editor.chain().focus().setTextAlign('left').run()"
      >
        <Bars3BottomLeftIcon class="h-5 w-5" />
      </ToolbarButton>

      <ToolbarButton
        title="Align Center"
        :active="editor.isActive({ textAlign: 'center' })"
        @click="editor.chain().focus().setTextAlign('center').run()"
      >
        <Bars3CenterLeftIcon class="h-5 w-5" />
      </ToolbarButton>

      <ToolbarButton
        title="Align Right"
        :active="editor.isActive({ textAlign: 'right' })"
        @click="editor.chain().focus().setTextAlign('right').run()"
      >
        <Bars3BottomRightIcon class="h-5 w-5" />
      </ToolbarButton>

      <ToolbarButton
        title="Justify"
        :active="editor.isActive({ textAlign: 'justify' })"
        @click="editor.chain().focus().setTextAlign('justify').run()"
      >
        <Bars4Icon class="h-5 w-5" />
      </ToolbarButton>
    </div>

    <div class="h-6 w-px bg-gray-300 dark:bg-slate-700" />

    <!-- List -->
    <div class="flex items-center gap-1">
      <ToolbarButton
        title="Bullet List"
        :active="editor.isActive('bulletList')"
        @click="editor.chain().focus().toggleBulletList().run()"
      >
        <ListBulletIcon class="h-5 w-5" />
      </ToolbarButton>

      <ToolbarButton
        title="Ordered List"
        :active="editor.isActive('orderedList')"
        @click="editor.chain().focus().toggleOrderedList().run()"
      >
        <NumberedListIcon class="h-5 w-5" />
      </ToolbarButton>
    </div>

    <div class="h-6 w-px bg-gray-300 dark:bg-slate-700" />

    <!-- Insert -->
    <div class="flex items-center gap-1">
      <!-- Blockquote -->
      <ToolbarButton
        title="Blockquote"
        :active="editor.isActive('blockquote')"
        @click="editor.chain().focus().toggleBlockquote().run()"
      >
        <ChatBubbleLeftRightIcon class="h-5 w-5" />
      </ToolbarButton>

      <!-- Code Block -->
      <ToolbarButton
        title="Code Block"
        :active="editor.isActive('codeBlock')"
        @click="editor.chain().focus().toggleCodeBlock().run()"
      >
        <CodeBracketIcon class="h-5 w-5" />
      </ToolbarButton>

      <!-- Horizontal Rule -->
      <ToolbarButton
        title="Horizontal Rule"
        @click="editor.chain().focus().setHorizontalRule().run()"
      >
        <MinusIcon class="h-5 w-5" />
      </ToolbarButton>
    </div>

    <div class="h-6 w-px bg-gray-300 dark:bg-slate-700" />

    <!-- Link -->
    <div class="flex items-center gap-1">
      <ToolbarButton
        title="Link"
        :active="editor.isActive('link')"
        @click="setLink"
      >
        <LinkIcon class="h-5 w-5" />
      </ToolbarButton>

      <ToolbarButton
        title="Hapus Link"
        @click="editor.chain().focus().unsetLink().run()"
      >
        <LinkSlashIcon class="h-5 w-5" />
      </ToolbarButton>
    </div>

    <div class="h-6 w-px bg-gray-300 dark:bg-slate-700" />

    <!-- Table -->
    <DashboardDropdown>
      <template #trigger>
        <ToolbarButton title="Table">
          <TableCellsIcon class="h-5 w-5" />
        </ToolbarButton>
      </template>

      <button
        type="button"
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="insertTable"
      >
        Sisipkan Tabel
      </button>

      <div class="h-px bg-gray-200 dark:bg-slate-700"></div>

      <button
        type="button"
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="addRowBefore"
      >
        Tambah Baris Atas
      </button>

      <button
        type="button"
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="addRowAfter"
      >
        Tambah Baris Bawah
      </button>

      <button
        type="button"
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="deleteRow"
      >
        Hapus Baris
      </button>

      <div class="h-px bg-gray-200 dark:bg-slate-700"></div>

      <button
        type="button"
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="addColumnBefore"
      >
        Tambah Kolom Kiri
      </button>

      <button
        type="button"
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="addColumnAfter"
      >
        Tambah Kolom Kanan
      </button>

      <button
        type="button"
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="deleteColumn"
      >
        Hapus Kolom
      </button>

      <div class="h-px bg-gray-200 dark:bg-slate-700"></div>

      <div class="h-px bg-gray-200 dark:bg-slate-700"></div>

      <button
        type="button"
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="mergeCells"
      >
        Gabung Sel
      </button>

      <button
        type="button"
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 transition-colors hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="splitCell"
      >
        Pisah Sel
      </button>

      <button
        type="button"
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="toggleHeaderRow"
      >
        Header Row
      </button>

      <div class="h-px bg-gray-200 dark:bg-slate-700"></div>

      <button
        type="button"
        class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-slate-200 dark:hover:bg-slate-800"
        @click="deleteTable"
      >
        Hapus Tabel
      </button>
    </DashboardDropdown>

    <!-- Image -->
    <div class="flex items-center gap-1">
      <ImageUploadButton :editor="editor" />
    </div>

    <!-- Image Align -->
    <div class="flex items-center gap-1">
      <ToolbarButton
        title="Align Gambar Kiri"
        :active="editor.isActive('image', { align: 'left' })"
        @click="setImageAlign('left')"
      >
        <Bars3BottomLeftIcon class="h-5 w-5" />
      </ToolbarButton>

      <ToolbarButton
        title="Align Gambar Tengah"
        :active="editor.isActive('image', { align: 'center' })"
        @click="setImageAlign('center')"
      >
        <Bars3CenterLeftIcon class="h-5 w-5" />
      </ToolbarButton>

      <ToolbarButton
        title="Align Gambar Kanan"
        :active="editor.isActive('image', { align: 'right' })"
        @click="setImageAlign('right')"
      >
        <Bars3BottomRightIcon class="h-5 w-5" />
      </ToolbarButton>
    </div>
  </div>
</template>

<script setup>
import ToolbarButton from "./ToolbarButton.vue";
import TableDropdown from "./TableDropdown.vue";
import ImageUploadButton from "./ImageUploadButton.vue";
import DashboardDropdown from "@/components/dashboard/DashboardDropdown.vue";

import {
  ArrowUturnLeftIcon,
  ArrowUturnRightIcon,
  Bars3BottomLeftIcon,
  Bars3CenterLeftIcon,
  Bars3BottomRightIcon,
  Bars4Icon,
  ListBulletIcon,
  NumberedListIcon,
  ChatBubbleLeftRightIcon,
  CodeBracketIcon,
  LinkIcon,
  LinkSlashIcon,
  TableCellsIcon,
  PlusIcon,
  MinusIcon,
  TrashIcon,
} from "@heroicons/vue/24/outline";

const props = defineProps({
  editor: {
    type: Object,
    required: true,
  },
});

function setLink() {
  console.log("Link diklik");
  const url = window.prompt(
    "Masukkan URL",
    props.editor.getAttributes("link").href || "",
  );

  if (url === null) return;

  if (url === "") {
    props.editor.chain().focus().unsetLink().run();
    return;
  }

  props.editor
    .chain()
    .focus()
    .extendMarkRange("link")
    .setLink({ href: url })
    .run();
}

function setImageAlign(align) {
  props.editor.chain().focus().setImageAlign(align).run();
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

function mergeCells() {
  props.editor.chain().focus().mergeCells().run();
}

function splitCell() {
  props.editor.chain().focus().splitCell().run();
}
</script>
