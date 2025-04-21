/* Récupère l'identifiant du post depuis l'URL */
/* Récupère le post depuis le backend */
/* Injecte le contenu dans la page HTML */
/* Affiche une erreur si la récupération échoue */

document.addEventListener("DOMContentLoaded", async () => {
  const params = new URLSearchParams(window.location.search);
  const postId = params.get("id");

  if (!postId) {
    document.querySelector(".post-card").innerHTML = "<p>Post introuvable.</p>";
    return;
  }

  try {
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
    document.querySelector(".post-meta").textContent = `${formattedDate}`;

    // 🔹 Affichage des likes/dislikes
    const postActions = document.createElement("div");
    postActions.className = "post-actions";
    postActions.innerHTML = `
      <button id="btn-like">👍</button>
      <span id="like-count">${post.likes || 0}</span>
      <button id="btn-dislike">👎</button>
      <span id="dislike-count">${post.dislikes || 0}</span>
    `;
    const contentBlock = document.querySelector(".post-content");
    contentBlock.style.position = "relative";

    postActions.style.position = "absolute";
    postActions.style.right = "0";
    postActions.style.bottom = "-40px";

    contentBlock.appendChild(postActions);

    const csrf = document.cookie
      .split("; ")
      .find(row => row.startsWith("csrf_token="))
      ?.split("=")
      .slice(1)
      .join("=");

    const btnLike = document.getElementById("btn-like");
    const btnDislike = document.getElementById("btn-dislike");
    const likeCount = document.getElementById("like-count");
    const dislikeCount = document.getElementById("dislike-count");

    btnLike?.addEventListener("click", () => sendReaction(1));
    btnDislike?.addEventListener("click", () => sendReaction(-1));

    async function sendReaction(value) {
      const action = value === 1 ? "like" : "dislike";
      const form = new URLSearchParams();
      form.append("post_id", postId);
      form.append("action", action); // <- ici au lieu de value
    
      const res = await fetch("/api/like", {
        method: "POST",
        body: form,
        headers: {
          "X-CSRF-Token": csrf,
          "Content-Type": "application/x-www-form-urlencoded"
        },
        credentials: "include"
      });
    
      if (res.ok) {
        const data = await res.json();
        if (data.success) location.reload();
      } else {
        alert("Action impossible. Êtes-vous connecté ?");
      }
    }    

    // 🔹 Chargement des commentaires
    const container = document.querySelector(".comments-section");
    fetch(`/comments?id=${postId}`)
      .then(res => res.json())
      .then(comments => {
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

    // 🔹 Envoi d’un commentaire
    const commentForm = document.querySelector(".comment-form");
    if (commentForm) {
      commentForm.addEventListener("submit", async e => {
        e.preventDefault();
        const textarea = commentForm.querySelector("textarea");
        const content = textarea.value.trim();
        if (!content) return;

        const resUser = await fetch("/me");
        const username = await resUser.text();

        const formData = new FormData();
        formData.append("post_id", postId);
        formData.append("content", content);
        formData.append("username", username);

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
    }

  } catch (err) {
    console.error("Erreur:", err);
    document.querySelector(".post-card").innerHTML = "<p>Erreur lors du chargement du post.</p>";
  }
});
