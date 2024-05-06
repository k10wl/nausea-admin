function onImageChange(e, selector) {
  if (!(e.currentTarget instanceof HTMLInputElement)) {
    throw new Error("not an input");
  }
  const target = document.querySelector(selector);
  if (!(target instanceof HTMLImageElement)) {
    throw new Error("target is not an image");
  }
  const reader = new FileReader();
  reader.onload = (e) => {
    target.src = e.target.result;
  };
  reader.readAsDataURL(e.currentTarget.files[0]);
}

function dublicateClick(selector) {
  const target = document.querySelector(selector);
  if (!target) {
    throw new Error("failed to find selector");
  }
  target.click();
}

function resetFileInput(selector) {
  const target = document.querySelector(selector);
  if (!(target instanceof HTMLInputElement)) {
    throw new Error("not an input");
  }
  target.files = new DataTransfer().files;
}
