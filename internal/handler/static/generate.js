document.addEventListener("DOMContentLoaded", () => {
  const errorContainer = document.getElementById("error-container");
  const pdfContainer = document.getElementById("pdfContainer");
  const pdfViewer = document.getElementById("pdfViewer");
  const fileLink = document.getElementById("file-link");
  const htmlContent = document.getElementById("htmlContent");
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
        showError("Your HTML is empty.");
        break;
      case 413:
        showError("The HTML provided is too large.");
        break;
      default:
        showError("PDF generation failed. Unexpected error.");
    }
    showSpinner(false);
  };

  document.getElementById("generatePDF").addEventListener("click", async () => {
    clear();
    if (!htmlContent.value) {
      handleError(1);
      return;
    }
    showSpinner(true);

    try {
      const response = await fetch("/pdf/generate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ html: btoa(htmlContent.value) }),
      });

      if (response.ok) {
        const blob = await response.blob();
        const blobUrl = URL.createObjectURL(blob);
        pdfViewer.setAttribute("src", blobUrl);
        fileLink.href = blobUrl;
        pdfContainer.style.display = "flex"; // Show the PDF container
      } else {
        handleError(response.status);
      }
    } catch (error) {
      showError(response.status);
      console.error(`An error occurred: ${error}`);
    }
    showSpinner(false);
  });
});
