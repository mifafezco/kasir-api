package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"github.com/spf13/viper"
	"os"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
)


type Config struct {
	Port    string `mapstructure:"PORT"`
	DBConn 	string `mapstructure:"DB_CONN"`
}

// var category = []Categories{
// 	{ID: 1, Name: "SUV", Description: "Mobil SUV"},
// 	{ID: 2, Name: "Minibus", Description: "Mobil Minibus"},
// 	{ID: 3, Name: "Truck", Description: "Truck"},
// }



func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}
	
	config := Config{
 	Port: viper.GetString("PORT"),
	DBConn: viper.GetString("DB_CONN"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	categoriesRepo := repositories.NewCategoriesRepository(db)
	categoriesService := services.NewCategoriesService(categoriesRepo)
	categoriesHandler := handlers.NewCategoriesHandler(categoriesService)
	
	// Setup routes
	http.HandleFunc("/api/categories", categoriesHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoriesHandler.HandleCategoriesByID)

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})
	
	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("gagal running server", err)
	}

}