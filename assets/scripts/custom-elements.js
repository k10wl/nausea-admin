customElements.define(
  "custom-dialog",
  class CustomDialog extends HTMLElement {
    slot = null;
    constructor() {
      super();
      this.attachShadow({ mode: "open" });
      const style = document.createElement("style");
      style.textContent = `
    dialog {
      position: fixed;
      top: 50%;
      left: 50%;
      translate: -50% -50%;
      padding: 0;
      margin: 0;
    }
    dialog > div {
      padding: var(--space-sm);
      margin: 0;
    }
    `;
      this.shadowRoot.appendChild(style);
      this.dialog = document.createElement("dialog");
      this.dialog.innerHTML = "<div><slot></slot></div>";
      this.shadowRoot.appendChild(this.dialog);
      this.slot = this.getSlot.bind(this)();
      this.onSuccessfullHTMXRequest = this.onSuccessfullHTMXRequest.bind(this);
    }

    connectedCallback() {
      this.htmxEvent();
      const trigger = this.getAttribute("trigger");
      const triggerEl = document.querySelector(trigger);
      if (!triggerEl) {
        throw new Error("Cannot find trigger for dialog");
      }
      triggerEl.addEventListener("click", () => {
        const onopen = this.getAttribute("onopen");
        if (onopen) {
          try {
            eval(onopen);
          } catch (e) {
            console.log(e);
          }
        }
        this.dialog.showModal();
      });
      this.dialog.addEventListener("click", (e) => {
        if (e.target === this.dialog) {
          this.closeAndReset();
        }
      });
      this.slot.addEventListener("slotchange", () => {
        const nodes = this.slot.assignedElements({});
        nodes.forEach((node) => {
          const elements = node.querySelectorAll("button[type='reset']");
          for (let i = 0; i < elements.length; i++) {
            elements
              .item(i)
              .addEventListener("click", () => this.closeAndReset());
          }
        });
      });
    }

    disconnectedCallback() {
      document.body.removeEventListener(
        "htmx:afterRequest",
        this.onSuccessfullHTMXRequest,
      );
    }

    htmxEvent() {
      document.body.addEventListener(
        "htmx:afterRequest",
        this.onSuccessfullHTMXRequest,
      );
    }

    onSuccessfullHTMXRequest(e) {
      if (200 <= e.detail.xhr.status && e.detail.xhr.status < 400) {
        this.closeAndReset();
      }
    }

    closeAndReset() {
      const nodes = this.slot.assignedNodes({ flatten: true });
      nodes.forEach((node) => {
        if (node.tagName === "FORM") {
          node.reset();
        }
      });
      this.shadowRoot.querySelector("dialog").close();
      const onclose = this.getAttribute("onclose");
      if (onclose) {
        try {
          eval(onclose);
        } catch (e) {
          console.log(e);
        }
      }
    }

    getSlot() {
      return this.shadowRoot.querySelector("slot");
    }
  },
);

customElements.define(
  "image-upload",
  class ImageUpload extends HTMLElement {
    constructor() {
      super();
      this.attachShadow({ mode: "open" });
      this.mount = this.mount.bind(this);
      this.clear = this.clear.bind(this);
      this.createForm = this.createForm.bind(this);
    }

    connectedCallback() {
      this.mount();
    }

    mount() {
      this.clear();
      this.createForm();
    }
    clear() {
      this.shadowRoot.innerHTML = "";
    }

    createForm() {
      const form = document.createElement("form");
      form.method = "dialog";
      const input = document.createElement("input");
      const imageContainer = document.createElement("div");
      imageContainer.id = "image-preview-container";
      input.type = "file";
      input.accept = "image/*";
      input.multiple = true;
      const dataTransfer = new DataTransfer()
      input.addEventListener("change", (e) => {
        const filesShift = [...dataTransfer.files, ...input.files]
        dataTransfer.items.clear()
        filesShift.forEach(file => dataTransfer.items.add(file))
        input.files = dataTransfer.files
        // kinda meh but works until there is no inner state on images
        imageContainer.innerHTML=''
        for (let i = 0; i < e.currentTarget.files.length; i++) {
          const reader = new FileReader();
          reader.onload = (e) => {
            const wrapper = document.createElement("div");
            wrapper.classList.add("wrapper");
            const button = document.createElement("button");
            button.onclick = () => {
              dataTransfer.items.remove(i)
              input.files = dataTransfer.files
              wrapper.remove();
            };
            button.textContent = "x";
            button.classList.add("remove");
            const image = document.createElement("img");
            image.src = e.target.result;
            wrapper.append(image, button);
            imageContainer.appendChild(wrapper);
          };
          reader.readAsDataURL(input.files.item(i));
        }
      });
      const button = document.createElement("button");
      button.innerText = "upload";
      button.classList.add("upload")
      const style = document.createElement("style");
      style.innerText = `
img {
  margin: auto;
  width: 100%;
  height: 100%;
  object-fit: contain;
}
#image-preview-container {
  margin-top: 1rem;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  grid-template-rows: repeat(auto-fill, 10rem);
  grid-gap: 0.25rem;
}
#image-preview-container:empty {
  display: none;
}
.wrapper {
  position: relative;
}
.remove {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
}
.upload {
  display: inline-block;
  margin-left: auto;
}
      `;
      form.append(input, button, imageContainer);
      this.shadowRoot.append(style, form);
    }
  },
);
