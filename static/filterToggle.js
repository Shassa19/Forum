document.addEventListener("DOMContentLoaded", () => {
    const filters = document.getElementById("filters");
    const burgerCheckbox = document.getElementById("hi");
  
    burgerCheckbox.addEventListener("change", () => {
      filters.classList.toggle("closed");
      filters.classList.toggle("open");
    });
  });
  
  