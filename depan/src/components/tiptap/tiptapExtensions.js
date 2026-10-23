import StarterKit from "@tiptap/starter-kit";
import Underline from "@tiptap/extension-underline";
import Link from "@tiptap/extension-link";
import CustomImage from "./CustomImage";
import TextAlign from "@tiptap/extension-text-align";
import Placeholder from "@tiptap/extension-placeholder";

import {
  Table,
  TableRow,
  TableHeader,
  TableCell,
} from "@tiptap/extension-table";

export const tiptapExtensions = [
  StarterKit.configure({
    heading: {
      levels: [1, 2, 3],
    },

    link: false,
    underline: false,
  }),

  Underline,

  Link.configure({
    openOnClick: false,
    autolink: true,
    linkOnPaste: true,
    HTMLAttributes: {
      class: "text-blue-600 underline",
      rel: "noopener noreferrer",
      target: "_blank",
    },
  }),

  CustomImage.configure({
    resize: {
      enabled: true,
      directions: ["top", "bottom", "left", "right"],
      minWidth: 50,
      minHeight: 50,
      alwaysPreserveAspectRatio: true,
    },
  }),

  TextAlign.configure({
    types: ["heading", "paragraph"],
  }),

  Placeholder.configure({
    placeholder: ({ node }) => {
      if (node.type.name === "paragraph") {
        return "Ketik disini...";
      }

      return "";
    },
  }),

  Table.configure({
    resizable: true,
  }),

  TableRow,
  TableHeader,
  TableCell,
];
