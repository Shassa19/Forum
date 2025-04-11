document.addEventListener("DOMContentLoaded", async () => {
    const params = new URLSearchParams(window.location.search);
    const postId = params.get("id");
  
    if (!postId) {
      document.querySelector(".post-box").innerHTML = "<p>Post introuvable</p>";
      return;
    }
  
    try {
      const res = await fetch(`/getPost?id=${postId}`);
      const post = await res.json();
  
      document.getElementById("post-title").textContent = post.title;
      document.getElementById("post-content").textContent = post.content;
      document.getElementById("post-author").textContent = post.username;
      document.getElementById("post-date").textContent = new Date(post.date).toLocaleDateString();
      document.getElementById("like-count").textContent = post.likes || 0;
  
      // ➕ tu peux aussi charger les commentaires ici
    } catch (err) {
      console.error("Erreur récupération post:", err);
      document.querySelector(".post-box").innerHTML = "<p>Erreur de chargement du post</p>";
    }
  });
  