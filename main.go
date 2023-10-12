package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kilianp07/r-place-back/db"
	"github.com/kilianp07/r-place-back/routes"
)

func main() {

	r := routes.SetupRoutes() // Configurez les routes en utilisant le package routes

	if err := db.CreateDatabaseTables(); err != nil {
		// Gérer les erreurs de création de table
		fmt.Println(err)
	}
	// Démarrer la goroutine de diffusion (à partir de routes package)
	go routes.BroadcastColorPlacements()

	// Créez un routeur Gorilla Mux
	router := mux.NewRouter()

	// Utilisez le routeur configuré avec vos routes
	router.PathPrefix("/").Handler(r)

	// Créez un serveur HTTP
	server := &http.Server{
		Addr:    ":8080", // Port sur lequel le serveur écoutera
		Handler: router,
	}

	// Démarrer le serveur
	if err := server.ListenAndServe(); err != nil {
		// Gérer les erreurs de démarrage du serveur
		panic(err)
	}
}
