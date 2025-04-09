document.addEventListener("DOMContentLoaded", () => {
  const popup = document.getElementById("popup-overlay");
  const postForm = document.getElementById("postForm");
  const cancelBtn = document.getElementById("cancelBtn");
  const popupUser = document.getElementById("popup-user");

  let currentUser = "";

  // Récupération du pseudo connecté
  fetch("/me")
    .then(res => res.ok ? res.text() : Promise.reject("Non connecté"))
    .then(username => {
      console.log("Utilisateur connecté :", username);
      currentUser = username;
    })
    .catch(err => console.error("Erreur récupération pseudo :", err));

  // Ouvre la popup uniquement au clic
  window.openPostPopup = () => {
    console.log("openPostPopup appelée");
    if (!currentUser) {
      alert("Erreur : utilisateur non connecté.");
      return;
    }

    postForm.reset(); // Réinitialise le formulaire
    popupUser.textContent = currentUser;
    popup.classList.remove("hidden");
  };

  // Ferme la popup
  cancelBtn.onclick = () => {
    console.log("Annuler cliqué");
    popup.classList.add("hidden");
  };

  // Envoi du formulaire
  postForm.onsubmit = async (e) => {
    e.preventDefault();

    const formData = new FormData(postForm);
    formData.append("username", currentUser);

    const csrf = document.cookie
      .split("; ")
      .find(row => row.startsWith("csrf_token="))
      ?.split("=")[1];

    const res = await fetch("/createPost", {
      method: "POST",
      body: formData,
      headers: {
        "X-CSRF-Token": csrf || ""
      }
    });

    if (res.ok) {
      alert("Post publié !");
      popup.classList.add("hidden");
      postForm.reset();
    } else {
      alert("Erreur lors de la publication.");
    }
  };
});
