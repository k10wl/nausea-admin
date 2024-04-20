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
