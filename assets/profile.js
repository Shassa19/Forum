let currentUser = "";

document.addEventListener("DOMContentLoaded", () => {
  // 🔹 Récupération de l'utilisateur connecté
  fetch("/me")
    .then(res => res.ok ? res.text() : Promise.reject("Non connecté"))
    .then(username => {
      currentUser = username;
      document.getElementById("profile-username").textContent = username;

      // 🔹 Récupération des topics de l'utilisateur
      fetch("/userPosts")
        .then(res => res.ok ? res.json() : Promise.reject("Erreur récupération topics"))
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

      // 🔹 Récupération des 3 derniers commentaires sur ses posts
      fetch("/last-comments")
        .then(res => res.ok ? res.json() : Promise.reject("Erreur récupération commentaires"))
        .then(comments => {
          const container = document.querySelector(".last-reply");
          const replyBlock = container.querySelector("p");

          if (comments.length === 0) {
            replyBlock.textContent = "Personne n'a répondu à tes posts visiblement...";
            return;
          }

          replyBlock.remove();

          comments.forEach(comment => {
            const link = document.createElement("a");
            link.href = `/post?id=${comment.post_id}`;
            link.className = "last-comment-link";
            link.textContent = `"${comment.content}" — par ${comment.author}`;

            const date = new Date(comment.date);
            const formatted = date.toLocaleDateString("fr-FR") + " à " + date.toLocaleTimeString("fr-FR", {
              hour: '2-digit',
              minute: '2-digit'
            });

            const span = document.createElement("span");
            span.textContent = ` (${formatted})`;

            const wrapper = document.createElement("p");
            wrapper.appendChild(link);
            wrapper.appendChild(span);
            container.appendChild(wrapper);
          });
        })
        .catch(err => {
          console.error("Erreur last-comments :", err);
          const container = document.querySelector(".last-reply p");
          container.textContent = "Erreur lors du chargement des dernières réponses.";
        });

      // 🔹 Infos utilisateur (avatar)
      fetch("/user-info")
        .then(res => res.ok ? res.json() : Promise.reject("Erreur récupération avatar"))
        .then(data => {
          document.getElementById("profile-username").textContent = data.username;
          document.getElementById("profile-avatar").src = `static/avatars/${data.avatar}`;
        });

    })
    .catch(err => {
      console.error("Erreur chargement profil :", err);
      document.querySelector(".profile-main").innerHTML = "<p>Erreur de chargement du profil.</p>";
    });
});

// 🔹 Fonctions utilitaires

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

// 🔹 Édition du profil
function openEditPopup() {
  document.getElementById("edit-popup").classList.remove("hidden");
}

function closeEditPopup() {
  document.getElementById("edit-popup").classList.add("hidden");
}

// 🔹 Envoi le nouveau pseudo + mdp coté back
document.getElementById("edit-form").addEventListener("submit", async (e) => {
  e.preventDefault();

  const newUsername = document.getElementById("new-username").value.trim();
  const newPassword = document.getElementById("new-password").value.trim();

  if (!newUsername || !newPassword) {
    alert("Veuillez remplir les deux champs.");
    return;
  }

  const formData = new FormData();
    formData.append("username", currentUser); // 🔥 essentiel pour Authorize
    formData.append("new_username", newUsername);
    formData.append("new_password", newPassword);


  const csrf = document.cookie
    .split("; ")
    .find(row => row.startsWith("csrf_token="))
    ?.split("=")
    .slice(1)
    .join("=") || "";

  const res = await fetch("/update-profile", {
    method: "POST",
    body: formData,
    headers: {
      "X-CSRF-Token": csrf
    }
  });

  if (res.ok) {
    alert("Profil mis à jour !");
    location.reload();
  } else {
    alert("Erreur lors de la mise à jour du profil.");
  }
});
