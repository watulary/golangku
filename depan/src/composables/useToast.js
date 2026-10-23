import { readonly, ref } from "vue";

const toasts = ref([]);

function createToast(type, message, duration = 4000, onClose = null) {
  const id = Date.now().toString(36) + Math.random().toString(36).slice(2);

  const toast = {
    id,
    type,
    message,
    duration,
  };

  toasts.value.push(toast);

  if (duration > 0) {
    setTimeout(() => {
      remove(id);

      if (typeof onClose === "function") {
        onClose();
      }
    }, duration);
  }

  return id;
}

function remove(id) {
  const index = toasts.value.findIndex((toast) => toast.id === id);

  if (index !== -1) {
    toasts.value.splice(index, 1);
  }
}

export function useToast() {
  return {
    toasts: readonly(toasts),

    success(message, duration = 4000, onClose = null) {
      return createToast("success", message, duration, onClose);
    },

    error(message, duration = 5000) {
      return createToast("error", message, duration);
    },

    warning(message, duration = 4500) {
      return createToast("warning", message, duration);
    },

    info(message, duration = 4000) {
      return createToast("info", message, duration);
    },

    remove,
  };
}
