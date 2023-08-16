document.addEventListener("DOMContentLoaded", () => {
  const pdfInput = document.getElementById("pdfInput");
  const pdfPassword = document.getElementById("pdfPassword");
  const pdfPasswordConfirm = document.getElementById("pdfPasswordConfirm");
  const encryptButton = document.getElementById("encryptButton");
  const errorContainer = document.getElementById("error-container");
  const pdfContainer = document.getElementById("pdfContainer");
  const fileLink = document.getElementById("file-link");
  const pdfViewer = document.getElementById("pdfViewer");
  const errorMessage = document.getElementById("error");
  const spinner = document.getElementById("spinner");
  const showPassword = document.getElementById("showPassword");

  const showError = (message) => {
    errorMessage.innerText = message;
    errorContainer.style.display = "flex";
  };
  const showSpinner = (show) => {
    spinner.style.display = show ? "flex" : "none";
  };
  const clear = () => {
    pdfContainer.style.display = "none";
    errorContainer.style.display = "none";
    showPassword.checked = false;
    showSpinner(false);
  };
  clear();

  const handleError = (errorCode) => {
    switch (errorCode) {
      case 0:
        showError("Passwords do not match.");
        break;
      case 1:
        showError("Please select a PDF file.");
        break;
      case 413:
        showError("The selected files are too large.");
        break;
      default:
        showError("PDF encrypting failed. Unexpected error.");
        break;
    }
    showSpinner(false);
  };

  const checkPassword = () => {
    return (
      pdfPassword.value !== pdfPasswordConfirm.value ||
      pdfPassword.value === "" ||
      pdfPasswordConfirm.value === ""
    );
  };

  showPassword.addEventListener("click", () => {
    if (pdfPassword.type === "password") {
      pdfPassword.type = "text";
      pdfPasswordConfirm.type = "text";
    } else {
      pdfPassword.type = "password";
      pdfPasswordConfirm.type = "password";
    }
  });

  encryptButton.addEventListener("click", async () => {
    clear();
    const selectedFiles = pdfInput.files;
    if (checkPassword()) {
      handleError(0);
      return;
    }
    if (selectedFiles.length === 0) {
      handleError(1);
      return;
    }
    showSpinner(true);

    const formData = new FormData();
    formData.append("file", selectedFiles[0]);
    formData.append("password", pdfPassword.value);

    try {
      const response = await fetch("/pdf/encrypt", {
        method: "POST",
        body: formData,
      });

      if (response.ok) {
        const blob = await response.blob();
        const blobUrl = URL.createObjectURL(blob);
        pdfViewer.setAttribute("src", blobUrl);
        fileLink.href = blobUrl;
        pdfContainer.style.display = "flex";
      } else {
        handleError(response.status);
      }
    } catch (error) {
      handleError(response.status);
      console.error(`An error occurred: ${error}`);
    }
    showSpinner(false);
  });
});
