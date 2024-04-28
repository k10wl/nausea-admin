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
    checkbox.checked = $nausea.get("showDeleted");
    toggleDeletedClass(checkbox.checked);
    checkbox.addEventListener("change", function () {
      toggleDeletedClass(checkbox.checked);
      $nausea.set("showDeleted", checkbox.checked);
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
    this.updateContainer = this.updateContainer.bind(this);
  }
  /** @type {(files: FileList, offset?: number) => void} */
  updateContainer(fileList, offset = 0) {
    for (let i = 0; i < fileList.length; i++) {
      const file = fileList.item(i);
      this.container.appendChild(this.#createPreview(file, i + offset));
    }
  }

  /** @type {(files: File, offset?: number) => HTMLElement} */
  #createPreview(file, i) {
    const reader = new FileReader();
    const preview = this.template.content.firstElementChild.cloneNode(true);
    reader.onload = (fileReaderEvent) => {
      preview
        .querySelector("img")
        .setAttribute("src", fileReaderEvent.target.result);
      const button = preview.querySelector("button");
      button.setAttribute("index", i);
      button.addEventListener("click", (e) => {
        updatableInputFiles.remove(e.target.getAttribute("index"));
        if (e.target.parentElement.nextElementSibling) {
          this.#updateFollowingParrents(
            e.target.parentElement.nextElementSibling,
          );
        }
        e.target.parentElement.remove();
      });
    };
    reader.readAsDataURL(file);
    return preview;
  }

  clear() {
    this.container.innerHTML = "";
  }

  #updateFollowingParrents(element) {
    const button = element.querySelector("button");
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

document.querySelectorAll('button[class*="rename-"').forEach((button) => {
  button.addEventListener("click", (e) => {
    e.stopPropagation();
  });
});

class EditFolder {
  element;

  constructor(element) {
    this.element = element;
  }

  /** @param {{name: string, id: string}} data  */
  open(data) {
    this.updateForm(data);
    htmx.process(this.element);
    this.updateInputs(data);
  }

  /** @param {{name: string, id: string}} data  */
  updateForm(data) {
    const base = this.element.getAttribute("data-hx-base");
    this.element.setAttribute("hx-patch", base + data.id);
    this.element.setAttribute("hx-target", "#content-" + data.id);
  }

  /** @param {{name: string, id: string}} data  */
  updateInputs(data) {
    /** @type HTMLTextAreaElement */
    const nameEl = this.element.querySelector('textarea[name="name"]');
    nameEl.value = data.name;
    nameEl.updateHeight()
  }
}

class EditMedia extends EditFolder {
  /** @param {{name: string, description: string, mediaId: string, id: string, folderId: string}} data  */
  updateInputs(data) {
    super.updateInputs(data);
    /** @type HTMLTextAreaElement */
    const descriptionEl = this.element.querySelector('textarea[name="description"]');
    descriptionEl.value = data.description;
    descriptionEl.updateHeight()
  }

  /** @param {{name: string, description: string, mediaId: string, id: string, folderId: string}} data  */
  updateForm(data) {
    const base = this.element.getAttribute("data-hx-base");
    this.element.setAttribute(
      "hx-patch",
      `${base}${data.folderId}/${data.mediaId}`,
    );
    this.element.setAttribute("hx-target", "#content-" + data.id);
  }
}

const editFolder = new EditFolder(
  document.getElementById("rename-folder-form"),
);
const editMedia = new EditMedia(document.getElementById("rename-media-form"));
