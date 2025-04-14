/* Masquer ou afficher les sections sans recharger la page. */
/* Naviguer facilement entre choix → formulaire et formulaire → retour. */

function showForm(formId) {
    document.getElementById("choix").classList.add("hidden");
    document.getElementById(formId).classList.remove("hidden");
  }
  
  function back() {
    document.getElementById("login").classList.add("hidden");
    document.getElementById("register").classList.add("hidden");
    document.getElementById("choix").classList.remove("hidden");
  }
  