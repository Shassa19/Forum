/* Récupère l'identifiant du post depuis l'URL */
/* Récupère le post depuis le backend */
/* Injecte le contenu dans la page HTML */
/* Affiche une erreur si la récupération échoue */

// 🔹 Récupère l'identifiant du post depuis l'URL
document.addEventListener("DOMContentLoaded", async () => {
  const params = new URLSearchParams(window.location.search);
  const postId = params.get("id");

  if (!postId) {
    document.querySelector(".post-card").innerHTML = "<p>Post introuvable.</p>";
    return;
  }

  try {
    // 🔹 Récupération du post
    const res = await fetch(`/getPost?id=${postId}`);
    if (!res.ok) throw new Error("Post non trouvé");

    const post = await res.json();

    document.querySelector(".username").textContent = post.username;
    document.querySelector(".post-title").textContent = post.title;
    document.querySelector(".post-content").textContent = post.content;

    const date = new Date(post.date);
    const formattedDate = date.toLocaleDateString("fr-FR") + " - " + date.toLocaleTimeString("fr-FR", {
      hour: '2-digit',
      minute: '2-digit'
    });
    document.querySelector(".post-meta").textContent = `${formattedDate} - ${post.likes || 0} like(s)`;

    // 🔹 Affichage des commentaires
    fetch(`/comments?id=${postId}`)
      .then(res => res.json())
      .then(comments => {
        const container = document.querySelector(".comments-section");
        container.innerHTML = "";

        if (comments.length === 0) {
          container.innerHTML = "<p>Aucun commentaire pour l’instant.</p>";
          return;
        }

        comments.forEach(comment => {
          const el = document.createElement("div");
          el.className = "comment";
          el.innerHTML = `
            <div class="comment-avatar"></div>
            <div class="comment-body">
              <div class="comment-username">${comment.username}</div>
              <div class="comment-text">${comment.content}</div>
            </div>
          `;
          container.appendChild(el);
        });
      });

    // 🔹 Soumission du commentaire
    document.querySelector(".comment-form").addEventListener("submit", async e => {
      e.preventDefault();
      const textarea = e.target.querySelector("textarea");
      const content = textarea.value.trim();
      if (!content) return;

      // Récupération du username à la volée
      const resUser = await fetch("/me");
      const username = await resUser.text();

      const formData = new FormData();
      formData.append("post_id", postId);
      formData.append("content", content);
      formData.append("username", username);

      const csrf = document.cookie
        .split("; ")
        .find(row => row.startsWith("csrf_token="))
        ?.split("=")
        .slice(1)
        .join("="); // 🔥 garde tout même s'il y a des "="

      const res = await fetch("/add-comment", {
        method: "POST",
        body: formData,
        headers: {
          "X-CSRF-Token": csrf
        }
      });

      if (res.ok) {
        location.reload();
      } else {
        alert("Erreur lors de l’envoi du commentaire.");
      }
    });

  } catch (err) {
    console.error("Erreur:", err);
    document.querySelector(".post-card").innerHTML = "<p>Erreur lors du chargement du post.</p>";
  }
});

