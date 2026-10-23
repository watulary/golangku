import Image from "@tiptap/extension-image";

const CustomImage = Image.extend({
  addAttributes() {
    return {
      ...this.parent?.(),
      align: {
        default: "center", // atau "left"
        parseHTML: (element) => element.getAttribute("data-align") || "center",
        renderHTML: (attributes) => {
          return {
            "data-align": attributes.align,
          };
        },
      },
    };
  },

  addCommands() {
    return {
      ...this.parent?.(),
      setImageAlign:
        (align) =>
        ({ commands }) => {
          return commands.updateAttributes(this.name, { align });
        },
    };
  },
});

export default CustomImage;
