<template>
  <div
    class="rounded-xl border border-gray-200 bg-white dark:border-slate-700 dark:bg-slate-900"
  >
    <!-- Toolbar -->
    <Toolbar v-if="editor" :editor="editor" />

    <!-- Editor -->
    <EditorContent
      :editor="editor"
      class="tiptap"
      spellcheck="false"
      autocorrect="off"
      autocomplete="off"
      autocapitalize="off"
    />
  </div>
</template>

<script setup>
import { watch, onBeforeUnmount } from "vue";
import { EditorContent, useEditor } from "@tiptap/vue-3";

import Toolbar from "./Toolbar.vue";
import { tiptapExtensions } from "./tiptapExtensions";

import "./styles.css";

const props = defineProps({
  modelValue: {
    type: String,
    default: "",
  },

  editable: {
    type: Boolean,
    default: true,
  },
});

const emit = defineEmits(["update:modelValue"]);

const editor = useEditor({
  content: props.modelValue,

  extensions: tiptapExtensions,

  editable: props.editable,

  onUpdate({ editor }) {
    emit("update:modelValue", editor.getHTML());
  },
});

watch(
  () => props.modelValue,
  (value) => {
    if (!editor.value) return;

    if (value !== editor.value.getHTML()) {
      editor.value.commands.setContent(value, false);
    }
  },
);

watch(
  () => props.editable,
  (value) => {
    if (!editor.value) return;

    editor.value.setEditable(value);
  },
);

onBeforeUnmount(() => {
  editor.value?.destroy();
});
</script>
