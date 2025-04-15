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