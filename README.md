# 🗣️ Forum Web - Golang & SQLite

Ce projet est un **forum web complet** développé en **Go**, avec une base **SQLite**, du HTML/CSS pour l'interface et du JavaScript pour le rendu dynamique.  
Il permet à un utilisateur de s’inscrire, se connecter, publier un post, voir ses publications, et consulter les posts des autres.

---

## 📌 Fonctionnalités principales

- ✅ Connexion et inscription sécurisées (avec protection CSRF + sessions)
- 📝 Création de posts via une popup
- 🧾 Lecture d’un post dans une page dédiée (`post.html`)
- 👤 Profil avec historique des topics postés et dernière réponse reçue
- ⏳ Date/heure automatique à la création d’un post
- ⚙️ Compatible Docker (via `Dockerfile`)
- 📦 Système modulaire (code Go découpé)

---
## Utilité des fichiers

📁 assets/
➡ Fichiers JavaScript du front

auth.js : gestion des événements d’authentification (login, register, logout)

filterToggle.js : active/désactive les filtres sur la page principale (ex. populaires, likés)

postDetails.js : charge les détails d’un post spécifique + commentaires

postLoader.js : récupère et affiche les posts sur la page principale

postPopup.js : gestion du popup de création de post

profile.js : interactions avec la page de profil (avatar, édition, commentaires récents)

📁 main/
➡ Backend en Go (logique serveur)

db.go : fonctions de requêtes SQL (posts, utilisateurs, commentaires, likes)

main.go : point d’entrée principal du serveur (handlers, routes, démarrage serveur)

session.go : gestion des sessions utilisateur et tokens CSRF

utils.go : fonctions utilitaires (validation, parsing, helpers)

📁 static/
➡ Contenu statique servi côté client

avatars/ : images d’avatar disponibles

fonts/ : polices custom du site

img/ : éventuelles images (fond, décor, etc.)

auth.css : style spécifique à la page d’authentification

post.css : style pour les pages de posts et commentaires

profil.css : style pour la page de profil utilisateur

style.css : style global (nav, header, footer, couleurs, polices)

🌐 Pages HTML
auth.html : page d’authentification (connexion, inscription)

index.html : page principale (accueil, posts visibles avec filtres)

post.html : page détaillée d’un post avec les commentaires

profil.html : profil de l’utilisateur connecté (topics publiés, avatar, édition)

⚙️ Autres fichiers
forum.db : base de données SQLite (locale)

go.mod / go.sum : dépendances du projet Go

dockerfile : configuration pour exécuter le projet dans un conteneur Docker

README.md : présentation rapide du projet

schema.sql : structure de la base de données (tables, champs, contraintes)

## 🧠 Arborescence du projet

Forum/
│
├── assets/
│   ├── auth.js
│   ├── filterToggle.js
│   ├── postDetails.js
│   ├── postLoader.js
│   ├── postPopup.js
│   └── profile.js
│
├── main/
│   ├── main.go
│   ├── db.go
│   ├── session.go
│   └── utils.go
│
├── static/
│   ├── auth.css
│   ├── post.css
│   ├── profil.css
│   └── style.css
│
├── auth.html
├── index.html
├── post.html
├── profil.html
│
├── forum.db
├── schema.sql
├── Dockerfile
├── go.mod
├── go.sum
└── README.md

---

## 🚀 Lancement de l'application

> Assurez-vous d’avoir Go installé.

```bash
cd main
go run /main/.

--- 

Ensuite, ouvrez dans votre navigateur :

http://localhost:8080/index


--- 

Projet développé par DEVERCHERE Geoffrey / DELARUE Louis / ASCIONE Romain
Dans un but pédagogique et personnel.