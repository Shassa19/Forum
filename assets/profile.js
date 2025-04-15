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

function openAvatarPopup() {
  document.getElementById("avatar-popup").classList.remove("hidden");
}

function selectAvatar(avatar) {
  fetch("/update-avatar", {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body: `avatar=${avatar}`
  }).then(() => {
    document.getElementById("profile-avatar").src = `static/avatars/${avatar}`;
    document.getElementById("avatar-popup").classList.add("hidden");
  });
}

document.addEventListener("DOMContentLoaded", () => {
  fetch("/user-info")
    .then(res => res.json())
    .then(data => {
      document.getElementById("profile-username").textContent = data.username;
      document.getElementById("profile-avatar").src = `static/avatars/${data.avatar}`;
    });
});

function logout() {
  const formData = new FormData();
  formData.append("username", currentUser);

  const csrf = document.cookie
    .split("; ")
    .find(row => row.startsWith("csrf_token="))
    ?.split("=")
    .slice(1)
    .join("=") || "";

  fetch("/logout", {
    method: "POST",
    body: formData,
    headers: { "X-CSRF-Token": csrf }
  })
  .then(res => {
    if (res.ok) {
      alert("Déconnexion réussie !");
      window.location.href = "/index";
    } else {
      alert("Erreur lors de la déconnexion.");
    }
  });
}
