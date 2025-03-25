package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	InitDB()

	fmt.Println("Bonjour depuis mon application Go en Docker !")

	// Servir les fichiers statiques
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes d'authentification
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)

	log.Println("Serveur lancé sur :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Exécution du template avec gestion des erreurs
func RenderTemplate(w http.ResponseWriter, tmpl Template, data interface{}) {
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Println("Erreur d'exécution du template :", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
	}
}
