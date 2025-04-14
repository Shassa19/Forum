let currentUser = "";

document.addEventListener("DOMContentLoaded", () => {
  // 🔄 Affiche le message de succès après redirection
  const params = new URLSearchParams(window.location.search);
  if (params.get("success") === "1") {
    const msg = document.getElementById("message");
    msg.textContent = "Connexion validée !";
    msg.classList.remove("hidden");
    setTimeout(() => msg.classList.add("hidden"), 3000);
  }

  // 🔄 Détection utilisateur connecté
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

  // 📥 Chargement des posts
  const postContainer = document.querySelector(".subject-list");

  fetch("/posts")
    .then(res => res.json())
    .then(posts => {
      postContainer.innerHTML = "";
      if (posts.length === 0) {
        postContainer.innerHTML = "<p>Aucun post pour l’instant.</p>";
        return;
      }

      posts.forEach(post => {
        const card = document.createElement("div");
        card.className = "subject-card";
        card.innerHTML = `
        <a href="/post?id=${post.id}" class="post-link">
          <h3>${post.title}</h3>
          <p>${post.content}</p>
          <span>${post.username} — ${new Date(post.date).toLocaleDateString()}</span>
        </a>
      `;

        postContainer.appendChild(card);
      });
    })
    .catch(err => {
      console.error("Erreur lors du chargement des posts :", err);
      postContainer.innerHTML = "<p>Erreur de chargement des posts.</p>";
    });

  // 🚪 Déconnexion
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

document.getElementById("btn-profile").onclick = () => {
  window.location.href = "../profil";
};
