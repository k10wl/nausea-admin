const folderContentsEl = document.getElementById("folder-contents");

const SHOW_DELETED_CONTENT_CLASS = "show-deleted-content";

class ShowDeleted extends HTMLElement {
  constructor() {
    super();
    const shadow = this.attachShadow({ mode: "open" });
    const checkbox = document.createElement("input");
    const label = document.createElement("label");
    label.innerText = "Show deleted";
    checkbox.type = "checkbox";
    checkbox.checked = $.get("showDeleted");
    toggleDeletedClass(checkbox.checked);
    checkbox.addEventListener("change", function () {
      toggleDeletedClass(checkbox.checked);
      $.set("showDeleted", checkbox.checked);
    });
    label.appendChild(checkbox);
    shadow.appendChild(label);
  }
}
customElements.define("show-deleted-checkbox", ShowDeleted);

function toggleDeletedClass(show) {
  if (show) {
    folderContentsEl.classList.add(SHOW_DELETED_CONTENT_CLASS);
  } else {
    folderContentsEl.classList.remove(SHOW_DELETED_CONTENT_CLASS);
  }
}
