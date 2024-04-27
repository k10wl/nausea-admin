/** @type {HTMLImageElement} */
const img = document.getElementById("about-image");
/** @type {HTMLInputElement} */
const fileInput = document.querySelector("input[type=\"file\"]")

/** @param {Event} e */
function onAboutImageChange(e) {
  if (!(e.currentTarget instanceof HTMLInputElement)) {
    throw new Error("not an input");
  }
  const reader = new FileReader();
  reader.onload = (e) => {
    img.src = e.target.result;
  };
  reader.readAsDataURL(e.currentTarget.files[0]);
}

img.onclick = (e) => {
  if (!(e.currentTarget instanceof HTMLImageElement)) {
    throw new Error("not an image");
  }
  fileInput.click()
}

function resetFileInput() {
  fileInput.files = new DataTransfer().files
}
