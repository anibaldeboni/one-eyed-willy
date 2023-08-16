document.addEventListener("DOMContentLoaded", () => {
  const pdfInput = document.getElementById("pdfInput");
  const mergeButton = document.getElementById("mergeButton");
  const errorContainer = document.getElementById("error-container");
  const pdfContainer = document.getElementById("pdfContainer");
  const fileLink = document.getElementById("file-link");
  const pdfViewer = document.getElementById("pdfViewer");
  const errorMessage = document.getElementById("error");
  const spinner = document.getElementById("spinner");

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
    showSpinner(false);
  };
  clear();

  const handleError = (errorCode) => {
    switch (errorCode) {
      case 1:
        showError("Please select at least two PDF files.");
        break;
      case 413:
        showError("The selected files are too large.");
        break;
      default:
        showError("PDF merging failed. Unexpected error.");
    }
    showSpinner(false);
  };

  mergeButton.addEventListener("click", async () => {
    clear();
    showSpinner(true);
    const selectedFiles = pdfInput.files;
    if (selectedFiles.length < 2) {
      handleError(1);
      return;
    }

    const formData = new FormData();
    for (let i = 0; i < selectedFiles.length; i++) {
      formData.append("files", selectedFiles[i]);
    }

    try {
      const response = await fetch("/pdf/merge", {
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
