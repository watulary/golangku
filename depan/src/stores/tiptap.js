import { defineStore } from "pinia";
import api from "@/services/api";

export const useTiptapStore = defineStore("tiptap", {
  actions: {
    async uploadImage(file) {
      const formData = new FormData();

      formData.append("image", file);

      const response = await api.post("/tiptap/upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });

      return response.data;
    },
  },
});
