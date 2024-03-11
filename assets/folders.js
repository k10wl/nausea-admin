/** @type {HTMLDialogElement} */
const createFolderDialogEl = document.getElementById("create-folder-dialog")
/** @type {HTMLFormElement} */
const createFolderFormEl =document.getElementById("create-folder-form")

function openCreateFolderDialog() {
  createFolderDialogEl.showModal()
}

function closeCreateFolderDialog() {
  createFolderFormEl.reset()
  createFolderDialogEl.close()
}

function createFolderError() {

}

createFolderDialogEl.addEventListener("click", (e) => {
  if (e.target === createFolderDialogEl) {
    closeCreateFolderDialog()
  }
})
