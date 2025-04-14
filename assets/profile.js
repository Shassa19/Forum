document.addEventListener("DOMContentLoaded", () => {
  // 🔹 Récupération de l'utilisateur connecté
  fetch("/me")
    .then(res => res.ok ? res.text() : Promise.reject("Non connecté"))
    .then(username => {
      document.getElementById("profile-username").textContent = username;

      // 🔹 Récupération des topics de l'utilisateur
      fetch("/userPosts")
        .then(res => res.json())
        .then(posts => {
          const container = document.querySelector(".topics-list");
          container.innerHTML = "<h3>Mes topics</h3>";

          if (posts.length === 0) {
            container.innerHTML += "<p>Aucun sujet publié.</p>";
            return;
          }

          posts.forEach(post => {
            const card = document.createElement("div");
            card.className = "topic-card";
            card.innerHTML = `<a href="/post?id=${post.id}">${post.title}</a>`;
            container.appendChild(card);
          });
        });
    })
    .catch(err => {
      document.querySelector(".profile-main").innerHTML = "<p>Erreur de chargement du profil.</p>";
      console.error(err);
    });
});
