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

    // Remplissage des données dans le HTML
    document.querySelector(".username").textContent = post.username;
    document.querySelector(".post-title").textContent = post.title;
    document.querySelector(".post-content").textContent = post.content;

    const date = new Date(post.date);
    const formattedDate = date.toLocaleDateString() + " - " + date.toLocaleTimeString();
    document.querySelector(".post-meta").textContent = `${formattedDate} - ${post.likes || 0} like(s)`;

    // Commentaires à venir...
  } catch (err) {
    console.error("Erreur:", err);
    document.querySelector(".post-card").innerHTML = "<p>Erreur lors du chargement du post.</p>";
  }
});
