function showForm(formId) {
    document.getElementById("choix").classList.add("hidden");
    document.getElementById(formId).classList.remove("hidden");
  }
  
  function back() {
    document.getElementById("login").classList.add("hidden");
    document.getElementById("register").classList.add("hidden");
    document.getElementById("choix").classList.remove("hidden");
  }
  