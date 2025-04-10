document.addEventListener("DOMContentLoaded", () => {
    const postContainer = document.querySelector(".subject-list");
  
    fetch("/posts")
      .then(res => res.json())
      .then(posts => {
        postContainer.innerHTML = ""; // Nettoie les anciens blocs (ex : exemples statiques)
  
        if (posts.length === 0) {
          postContainer.innerHTML = "<p>Aucun post pour l’instant.</p>";
          return;
        }
  
        posts.forEach(post => {
          const card = document.createElement("div");
          card.className = "subject-card";
          card.innerHTML = `
            <h3>${post.title}</h3>
            <p>${post.content}</p>
            <span>${post.username} — ${new Date(post.date).toLocaleDateString()}</span>
          `;
          postContainer.appendChild(card);
        });
      })
      .catch(err => {
        console.error("Erreur lors du chargement des posts :", err);
        postContainer.innerHTML = "<p>Erreur de chargement des posts.</p>";
      });
  });
  