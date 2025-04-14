/* Récupère le pseudo de l'utilisateur connecté + le csrf_token (depuis les cookies) */
/* Ouvre la popup de création de post quand on clique sur "Créer un sujet" */
/* Ferme la popup si l'utilisateur clique sur "Annuler" */
/* Gère l’envoi sécurisé du post avec fetch en POST et token CSRF dans les headers */

document.addEventListener("DOMContentLoaded", () => {
  const popup = document.getElementById("popup-overlay");
  const postForm = document.getElementById("postForm");
  const cancelBtn = document.getElementById("cancelBtn");
  const popupUser = document.getElementById("popup-user");

  let currentUser = "";
  let csrfToken = "";

  //récup pseudo utilisateur connecté
  fetch("/me")
    .then(res => res.ok ? res.text() : Promise.reject("Non connecté"))
    .then(username => {
      console.log("Utilisateur connecté :", username);
      currentUser = username;
      popupUser.textContent = username;

      // Récupération correcte même avec '=' dans la valeur
    csrfToken = decodeURIComponent(
      document.cookie
        .split("; ")
        .find(row => row.startsWith("csrf_token="))
        ?.split("=")
        .slice(1)
        .join("=") || ""
    );
    })
    .catch(err => console.error("Erreur récupération pseudo :", err));

  //ouvre la popup au clic
  window.openPostPopup = () => {
    console.log("openPostPopup appelée");
    if (!currentUser) {
      alert("Erreur : utilisateur non connecté.");
      return;
    }
    postForm.reset();
    popupUser.textContent = currentUser;
    popup.classList.remove("hidden");
  };

  //ferme la popup
  cancelBtn.onclick = () => {
    console.log("Annuler cliqué");
    popup.classList.add("hidden");
  };

  //envoi du formulaire
  postForm.onsubmit = async (e) => {
    e.preventDefault();

    const formData = new FormData(postForm);
    formData.append("username", currentUser);

    const res = await fetch("/createPost", {
      method: "POST",
      body: formData,
      headers: {
        "X-CSRF-Token": csrfToken || ""
      }
    });

    if (res.ok) {
      alert("Post publié !");
      popup.classList.add("hidden");
      postForm.reset();
    } else {
      console.log("Status:", res.status);
      const msg = await res.text();
      console.error("Erreur lors de la publication :", msg);
      alert("Erreur lors de la publication.");
    }
  };
});
