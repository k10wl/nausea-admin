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

const dataTransfer = new DataTransfer();
/** @type HTMLInputElement */
const input = document.getElementById("media-file-input");

function cleanupCreateFolder() {
  document.getElementById('create-folder-error').innerHTML = ''
}

const uploadPreviewContainer = document.getElementById(
  "upload-preview-container",
);
/** @type HTMLTemplateElement */
const uploadPreviewTemplate = document.getElementById(
  "upload-preview-template",
);

let filesArray = [];

function removeFile(index) {
  if (index < 0 || index >= filesArray.length) {
    console.error(index, filesArray)
    alert(
      "error upon removing file, please start over again (smth wrong was removed)",
    );
    return;
  }
  filesArray.splice(index, 1);
  updateInputElement();
}

function updateInputElement() {
  const dataTransfer = new DataTransfer();
  for (let i = 0; i < filesArray.length; i++) {
    dataTransfer.items.add(filesArray[i]);
  }
  input.files = dataTransfer.files;
}

function updateFileDisplay(newFiles, offset) {
  for (let i = 0; i < newFiles.length; i++) {
    const file = newFiles[i];
    const reader = new FileReader();
    const preview =
      uploadPreviewTemplate.content.firstElementChild.cloneNode(true);
    uploadPreviewContainer.appendChild(preview);
    reader.onload = (fileReaderEvent) => {
      /** @type HTMLElement */
      preview
        .querySelector("img")
        .setAttribute("src", fileReaderEvent.target.result);
      const button = preview.querySelector("button");
      button.setAttribute("index", i + offset);
      button.addEventListener("click", (e) => {
        removeFile(e.target.getAttribute("index"));
        /** @type HTMLDivElement */
        const parent = e.target.parentElement;
        updateFollowingParrents(parent.nextElementSibling);
        parent.remove();
      });
    };
    reader.readAsDataURL(file);
  }
}

function updateFollowingParrents(parent) {
  if (!parent) {
    return;
  }
  /** @type HTMLButtonElement */
  const button = parent.querySelector("button");
  button.setAttribute("index", button.getAttribute("index") - 1);
  updateFollowingParrents(parent.nextElementSibling);
}

input.addEventListener("change", () => {
  const newFiles = [...input.files];
  updateFileDisplay(newFiles, filesArray.length);
  filesArray.push(...newFiles);
  updateInputElement();
});

function openFileInput() {
 input.click()
}

function cleanupUpload() {
  document.getElementById('upload-media-error').innerHTML = ''
  uploadPreviewContainer.innerHTML = ''
  filesArray = []
}
