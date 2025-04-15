/* Affiche un message éphémère après une connexion réussie */
/* Récupère l'utilisateur connecté et affiche les bons boutons */
/* Charge dynamiquement les posts depuis la base */
/* Génère les cartes HTML des posts (avec lien vers page post.html) */
/* Gère la déconnexion sécurisée via POST et CSRF */
/* Redirige l’utilisateur vers la page /profil */
let allPosts = [];
let currentUser = "";

document.addEventListener("DOMContentLoaded", () => {
  // ✅ Message de succès après redirection
  const params = new URLSearchParams(window.location.search);
  if (params.get("success") === "1") {
    const msg = document.getElementById("message");
    msg.textContent = "Connexion validée !";
    msg.classList.remove("hidden");
    setTimeout(() => msg.classList.add("hidden"), 3000);
  }

  // ✅ Détection de l'utilisateur connecté
  fetch("/me")
    .then(res => res.ok ? res.text() : null)
    .then(username => {
      if (username) {
        currentUser = username;
        document.getElementById("btn-auth").classList.add("hidden");
        document.getElementById("btn-profile").classList.remove("hidden");
        document.getElementById("btn-logout").classList.remove("hidden");
      }
    });

  // ✅ Chargement initial des posts
  fetch("/posts")
    .then(res => res.json())
    .then(posts => {
      allPosts = posts;
      displayPosts(posts);
    })
    .catch(err => {
      console.error("Erreur lors du chargement des posts :", err);
      document.querySelector(".subject-list").innerHTML = "<p>Erreur de chargement des posts.</p>";
    });

  // ✅ Recherche dynamique dans la barre de recherche
  document.querySelector(".search-bar input").addEventListener("input", e => {
    const query = e.target.value.toLowerCase();
    const filtered = allPosts.filter(post =>
      post.title.toLowerCase().includes(query) ||
      post.content.toLowerCase().includes(query) ||
      post.username.toLowerCase().includes(query)
    );
    displayPosts(filtered);
  });

  // ✅ Déconnexion sécurisée
  document.getElementById("btn-logout").onclick = async () => {
    const formData = new FormData();
    formData.append("username", currentUser);

    const csrf = document.cookie
      .split("; ")
      .find(row => row.startsWith("csrf_token="))
      ?.split("=")
      .slice(1)
      .join("=") || "";

    const res = await fetch("/logout", {
      method: "POST",
      body: formData,
      headers: { "X-CSRF-Token": csrf }
    });

    if (res.ok) {
      currentUser = "";
      document.getElementById("btn-auth").classList.remove("hidden");
      document.getElementById("btn-profile").classList.add("hidden");
      document.getElementById("btn-logout").classList.add("hidden");
      alert("Déconnexion réussie !");
    } else {
      alert("Erreur lors de la déconnexion.");
    }
  };
});

// ✅ Bouton redirection vers la page profil
document.getElementById("btn-profile").onclick = () => {
  window.location.href = "../profil";
};

// ✅ Fonction d'affichage des posts
function displayPosts(posts) {
  const postContainer = document.querySelector(".subject-list");
  postContainer.innerHTML = "";

  if (posts.length === 0) {
    postContainer.innerHTML = "<p>Aucun résultat trouvé.</p>";
    return;
  }

  posts.forEach(post => {
    const card = document.createElement("div");
    card.className = "subject-card";

    const date = new Date(post.date);
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = String(date.getFullYear()).slice(-2);
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');

    const formattedDate = `${day}/${month}/${year} à ${hours}:${minutes}`;

    card.innerHTML = `
      <a href="/post?id=${post.id}" class="post-link">
        <h3>${post.title}</h3>
        <p>${post.content}</p>
        <span>${post.username} — ${formattedDate}</span>
      </a>
    `;
    postContainer.appendChild(card);
  });
}
