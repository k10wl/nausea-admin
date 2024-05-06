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
      this.attachTriggers();
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

    attachTriggers() {
      const trigger = this.getAttribute("trigger");
      const triggerList = document.querySelectorAll(trigger);
      triggerList.forEach((triggerEl) => {
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
        this.attachTriggers();
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
  "custom-scalable",
  class CustomImage extends HTMLElement {
    constructor() {
      super();
      const root = this.attachShadow({ mode: "open" });
      const img = document.createElement("img");
      img.loading = "lazy";
      this.getAttributeNames().forEach((name) =>
        img.setAttribute(name, this.getAttribute(name)),
      );
      const container = document.createElement("div");
      const slot = document.createElement("slot");
      container.appendChild(slot);
      slot.addEventListener("slotchange", () => {
        const [el] = slot.assignedElements();
        let prevSrc = el.src;
        el.addEventListener("click", () => {
          if (container.classList.toggle("expanded")) {
            el.src = this.getAttribute("src");
            document.body.style.setProperty("overflow", "hidden", "important");
          } else {
            el.src = prevSrc;
            document.body.style.removeProperty("overflow");
          }
        });
      });
      const style = document.createElement("style");
      style.textContent = `
div.expanded {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.1);
  backdrop-filter: blur(5px);
  & ::slotted(IMG) {
    margin: auto;
    object-fit: contain !important;
  }
}
      `;
      container.append(img);
      root.append(style, container);
    }
  },
);

customElements.define(
  "custom-auto-resize",
  class AutoResize extends HTMLTextAreaElement {
    static observedAttributes = ["textContent"];

    constructor() {
      super();
      this.updateHeight = this.updateHeight.bind(this);
    }

    connectedCallback() {
      this.style.overflow = "hidden";
      this.addEventListener("input", this.updateHeight);
      this.updateHeight();
    }

    updateHeight() {
      // aight, this is weird, but without timeout it does not update properly
      // I have neither no idea why the fuck it does that nor time to figure out
      setTimeout(() => {
        this.style.height = "1px";
        this.style.height = `${this.scrollHeight}px`;
      });
    }
  },
  { extends: "textarea" },
);
