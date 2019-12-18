package restapi

import (
	"log"
	"net/http"
	"os"
	"server/infrastructure"
	"server/infrastructure/client"
	"server/infrastructure/collection"
	"server/restapi/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Server ...
type Server struct {
	carsHandler handlers.CarsHandler
	router      *mux.Router
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

//NewServer ...
func NewServer(carsHandler handlers.CarsHandler) *Server {
	return &Server{
		carsHandler: carsHandler,
		router:      mux.NewRouter(),
	}
}

//ConfigureAndRun ...
func ConfigureAndRun() {
	serverAdres := getEnv("BIND_ADDR", ":8080")

	mongoclient := client.NewClient

	mongocollection := collection.NewCollection

	carsRepository := infrastructure.NewCarsRepository(mongoclient, mongocollection)

	carsHandler := handlers.NewCarsHandler(carsRepository)

	s := NewServer(carsHandler)

	s.router.HandleFunc("/index", s.carsHandler.Index()).Methods("GET")
	s.router.HandleFunc("/", s.carsHandler.GetAllCars()).Methods("GET")
	s.router.HandleFunc("/", s.carsHandler.InsertOneCar()).Methods("POST")
	s.router.HandleFunc("/{id}", s.carsHandler.DeleteOneCar()).Methods("DELETE")
	s.router.HandleFunc("/{id}", s.carsHandler.GetOneCar()).Methods("GET")
	s.router.HandleFunc("/", s.carsHandler.UpdateOneCar()).Methods("PUT")

	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))

	log.Printf("listening at %s\n", serverAdres)
	log.Fatal(http.ListenAndServe(serverAdres, s.router))
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
