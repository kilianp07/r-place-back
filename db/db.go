package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kilianp07/r-place-back/utils"
)

func CreateDatabaseTables() error {
	// Ouvrir une connexion à MySQL (assurez-vous que votre chaîne de connexion est correcte)
	db, err := sql.Open("mysql", "place:place@tcp(db:3306)/rplace")

	if err != nil {
		return err
	}
	defer db.Close()

	// Créer la base de données
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS rplace")
	if err != nil {
		return err
	}

	// Sélectionner la base de données
	_, err = db.Exec("USE `rplace`")
	if err != nil {
		return err
	}

	// Créer la table des placements de couleur
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS place (
            id INT AUTO_INCREMENT PRIMARY KEY,
            color VARCHAR(255),
            x INT,
            y INT
        )
    `)
	if err != nil {
		return err
	}

	return nil
}

func InsertPixel(pixel utils.Pixel) error {
	db, err := sql.Open("mysql", "place:place@tcp(db:3306)/rplace")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO place (color, x, y) VALUES (?, ?, ?)", pixel.Color, pixel.X, pixel.Y)
	if err != nil {
		return err
	}

	return nil
}
