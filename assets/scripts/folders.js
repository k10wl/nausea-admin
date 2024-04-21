const folderContentsEl = document.getElementById("folder-contents");

const SHOW_DELETED_CONTENT_CLASS = "show-deleted-content";

class ShowDeleted extends HTMLElement {
  constructor() {
    super();
    const shadow = this.attachShadow({ mode: "open" });
    const checkbox = document.createElement("input");
    const label = document.createElement("label");
    label.innerText = "Show hidden";
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

function toggleDeletedClass() {
  folderContentsEl.classList.toggle(SHOW_DELETED_CONTENT_CLASS);
}

class UpdatableInputFiles {
  /** @type {DataTransfer} @private */
  #dataTransfer;

  /** @type {HTMLInputElement} */
  input;

  /**
   * @param {Object} data
   * @param {HTMLInputElement} data.input
   * @param {(files: FileList, offset: number) => void} data.onChange - triggers upon input change
   */
  constructor({ input, onChange }) {
    this.input = input;
    this.#dataTransfer = new DataTransfer();
    input.addEventListener("change", () => {
      onChange(input.files, this.#dataTransfer.files.length);
      for (let i = 0; i < input.files.length; i++) {
        this.#dataTransfer.items.add(input.files.item(i));
      }
      this.#syncFiles();
    });
  }

  /** @param index {number} */
  remove(index) {
    if (index < 0 || index >= this.#dataTransfer.files.length) {
      throw new Error(
        `Can't access file at ${index} of ${this.#dataTransfer.files.length}`,
      );
    }
    this.#dataTransfer.items.remove(index);
    this.#syncFiles();
  }

  #syncFiles() {
    this.input.files = this.#dataTransfer.files;
  }

  clear() {
    this.#dataTransfer.items.clear();
    this.#syncFiles();
  }
}

class MediaPreview {
  /**
    * @param {Object} data
    * @param {HTMLElement} data.container
    * @param {HTMLTemplateElement} data.template
    * @param {(src: string) => HTMLElement} data.builder
    */
  constructor(data) {
    this.container = data.container;
    this.template = data.template;
  }
  /** @type {(files: FileList, offset?: number) => void} */
  updateContainer(fileList, offset = 0) {
    for (let i = 0; i < fileList.length; i++) {
      const file = fileList.item(i);
      this.container.appendChild(this.#createPreview(file, offset));
    }
  }

  /** @type {(files: File, offset?: number) => HTMLElement} */
  #createPreview(file, offset = 0) {
    const reader = new FileReader();
    const preview = this.template.content.firstElementChild.cloneNode(true);
    reader.onload = (fileReaderEvent) => {
      preview
        .querySelector("img")
        .setAttribute("src", fileReaderEvent.target.result);
      const button = preview.querySelector("button");
      button.setAttribute("index", i + offset);
      button.addEventListener("click", (e) => {
        updatableInputFiles.remove(e.target.getAttribute("index"));
        this.#updateFollowingParrents(
          e.target.parentElement.nextElementSibling,
        );
        e.target.parentElement.remove();
      });
    };
    reader.readAsDataURL(file);
    return preview
  }

  clear() {
    this.container.innerHTML = "";
  }

  #updateFollowingParrents(element) {
    const button = parent.querySelector("button");
    button.setAttribute("index", button.getAttribute("index") - 1);
    if (!element.nextElementSibling) {
      return;
    }
    this.#updateFollowingParrents(element.nextElementSibling);
  }
}

const input = document.getElementById("media-file-input");
const mediaPreview = new MediaPreview({
  container: document.getElementById("upload-preview-container"),
  template: document.getElementById("upload-preview-template"),
});
const updatableInputFiles = new UpdatableInputFiles({
  input: input,
  onChange: mediaPreview.updateContainer,
});

function cleanupCreateFolder() {
  document.getElementById("create-folder-error").innerHTML = "";
}

function openFileInput() {
  input.click();
}

function cleanupUpload() {
  document.getElementById("upload-media-error").innerHTML = "";
  updatableInputFiles.clear();
  mediaPreview.clear();
}

document.querySelectorAll("button[class*=\"rename-\"").forEach((button) => {
  button.addEventListener("click", (e) => {
    e.stopPropagation()
  })
})


const editFolder = document.getElementById("rename-folder-form")
function updateContent(name, id) {
  const base = editFolder.getAttribute('data-hx-base')
  editFolder.setAttribute('hx-patch', base+id)
  editFolder.setAttribute('hx-target', "#content-"+id)
  const textarea = editFolder.querySelector('textarea[name="name"]')
  textarea.textContent = name
  htmx.process(editFolder)
}
