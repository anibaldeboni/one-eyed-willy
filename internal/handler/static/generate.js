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

  document.getElementById("generatePDF").addEventListener("click", async () => {
    clear();
    if (!htmlContent.value) {
      showError("Your HTML is empty.");
      return;
    }
    showSpinner(true);

    try {
      const response = await fetch("/generate", {
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
        showError("PDF generation failed. Please try again.");
      }
    } catch (error) {
      showError("PDF generation failed. Unexpected error");
      console.error(`An error occurred: ${error}`);
    }
    showSpinner(false);
  });
});
