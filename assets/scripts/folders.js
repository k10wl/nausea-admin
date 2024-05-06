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

class SharedCustomDialog {
  /** @type {HTMLFormElement} */
  element;

  /** @param {HTMLFormElement} element */
  constructor(element) {
    this.element = element;
  }

  open(data) {
    this.updateForm(data);
    htmx.process(this.element);
    this.updateInputs(data);
  }
  updateForm() {}
  updateInputs() {}
}

class EditFolder extends SharedCustomDialog {
  /** @param {{name: string, id: string, fromInside?: boolean}} data  */
  updateForm(data) {
    const base = this.element.getAttribute("data-hx-base");
    let path = base + data.id;
    let target = "#content-" + data.id;
    if (data.fromInside) {
      this.element.setAttribute("hx-swap", "innerHTML");
      path += "?from-inside";
      target = "#folder-name";
    }
    this.element.setAttribute("hx-patch", path);
    this.element.setAttribute("hx-target", target);
  }

  /** @param {{name: string, id: string}} data  */
  updateInputs(data) {
    /** @type HTMLTextAreaElement */
    const nameEl = this.element.querySelector('textarea[name="name"]');
    nameEl.value = data.name;
    nameEl.updateHeight();
  }
}

class EditMedia extends EditFolder {
  /** @param {{name: string, description: string, mediaId: string, id: string, folderId: string}} data  */
  updateInputs(data) {
    super.updateInputs(data);
    /** @type HTMLTextAreaElement */
    const descriptionEl = this.element.querySelector(
      'textarea[name="description"]',
    );
    descriptionEl.value = data.description;
    descriptionEl.updateHeight();
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

class DeleteForever extends SharedCustomDialog {
  /** @param {{mediaId: string, id: string, folderId: string}} data  */
  updateForm(data) {
    const base = this.element.getAttribute("data-hx-base");
    this.element.querySelector("#delete-forever-error").innerHTML = "";
    let path = `${base}${data.folderId}`;
    if (data.mediaId) {
      path += `/${data.mediaId}`;
    }
    this.element.setAttribute("hx-delete", path);
    this.element.setAttribute("hx-target", "#content-" + data.id);
  }
}

const editFolder = new EditFolder(
  document.getElementById("rename-folder-form"),
);
const editMedia = new EditMedia(document.getElementById("rename-media-form"));
const deleteForever = new DeleteForever(
  document.getElementById("delete-forever-form"),
);

class DragReorder {
  /** @typedef {(e: DragEvent & {state: State}) => void} DragReoderEvent */

  /** @typedef {Object} DragReorderEvents
   * @property {DragReoderEvent} [onDragStart]
   * @property {DragReoderEvent} [onDragOver]
   * @property {DragReoderEvent} [onDragEnd]
   */

  /** @typedef {Object} State
   * @property {string} draggableQuerySelector
   * @property {HTMLElement | null} selected
   * @property {HTMLElement | null} over
   * @property {number} from
   * @property {number} to
   */

  /** @type {State} */
  #state;

  /** @type {DragReorderEvents} */
  #events;

  /** @param {{state: State, events: DragReorderEvents}} data */
  constructor(data) {
    this.#state = {
      selected: null,
      over: null,
      draggableQuerySelector: data.state.draggableQuerySelector,
      from: -1,
      to: -1,
    };
    this.#events = {
      onDragStart: data.events.onDragStart,
      onDragOver: data.events.onDragOver,
      onDragEnd: data.events.onDragEnd,
    };
  }

  /** @param {DragEvent} e */
  onDragStart(e) {
    e.dataTransfer.effectAllowed = "move";
    e.dataTransfer.setData("text/plain", null);
    this.#state.selected = e.target;
    this.#state.over = e.target;
    const nodeList = document.querySelectorAll(
      this.#state.draggableQuerySelector,
    );
    for (let i = 0; i < nodeList.length; i++) {
      if (this.#state.selected === nodeList.item(i)) {
        this.#state.from = i;
        break;
      }
    }
    this.#state.to = -1;
    e.state = this.#state;
    this.#events.onDragStart?.(e);
  }

  /** @param {DragEvent} e */
  onDragOver(e) {
    const over = this.#parentQuerySelector(
      e.target,
      this.#state.draggableQuerySelector,
    );
    if (!(over instanceof HTMLElement)) {
      return;
    }
    this.#state.over = over;
    e.state = this.#state;
    this.#events.onDragOver?.(e);
  }

  /** @param {DragEvent} e */
  onDragEnd(e) {
    const nodeList = document.querySelectorAll(
      this.#state.draggableQuerySelector,
    );
    for (let i = 0; i < nodeList.length; i++) {
      if (this.#state.selected === nodeList.item(i)) {
        this.#state.to = i;
        break;
      }
    }
    e.state = this.#state;
    this.#events.onDragEnd?.(e);
    this.#state.selected = null;
    this.#state.over = null;
  }

  /**
   * @param {HTMLElement} element
   * @param {QueuingStrategyInit} querySelector
   * @param {number} [maxDepth=100]
   * @returns {HTMLElement | null}
   */
  #parentQuerySelector(element, querySelector, maxDepth = 100) {
    if (maxDepth === 0) {
      return null;
    }
    if (element === null || element.matches(querySelector)) {
      return element;
    }
    return this.#parentQuerySelector(
      element.parentNode,
      querySelector,
      maxDepth - 1,
    );
  }
}

function isBefore(el1, el2) {
  let cur;
  if (el2.parentNode === el1.parentNode) {
    for (cur = el1.previousSibling; cur; cur = cur.previousSibling) {
      if (cur === el2) {
        return true;
      }
    }
  }
  return false;
}

const mediaDragReorder = new DragReorder({
  state: { draggableQuerySelector: "li[data-type='media']" },
  events: {
    onDragStart: (e) => e.target.classList.add("opacity-1/2"),
    onDragOver: (e) => {
      if (isBefore(e.state.selected, e.state.over)) {
        e.state.over.parentNode.insertBefore(e.state.selected, e.state.over);
        return;
      }
      e.state.over.parentNode.insertBefore(
        e.state.selected,
        e.state.over.nextSibling,
      );
    },
    onDragEnd: (e) => {
      e.target.classList.remove("opacity-1/2");
      if (e.state.from === e.state.to) {
        return;
      }
      htmx.ajax("POST", window.location.pathname + "/reorder-media", {
        values: { from: e.state.from, to: e.state.to },
        target: "none",
        swap: "none",
      });
    },
  },
});
