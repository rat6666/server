package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/infrastructure"
	"server/infrastructure/model"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CarsHandler ...
type CarsHandler interface {
	Index() http.HandlerFunc
	GetAllCars() http.HandlerFunc
	InsertOneCar() http.HandlerFunc
	DeleteOneCar() http.HandlerFunc
	GetOneCar() http.HandlerFunc
	UpdateOneCar() http.HandlerFunc
}

type carsHandler struct {
	carsRepository infrastructure.CarsRepository
}

//NewCarsHandler ...
func NewCarsHandler(carsRepository infrastructure.CarsRepository) CarsHandler {
	return &carsHandler{
		carsRepository: carsRepository,
	}
}

func (c *carsHandler) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "./templates/index.html")
	}
}

func (c *carsHandler) InsertOneCar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		result := &model.Car{IDMongo: primitive.NewObjectID()}

		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if result.ID == "" || result.Model == "" || result.Date == "" {
			err := errors.New("All fields must be filled")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := c.carsRepository.InsertOneCar(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	}
}

func (c *carsHandler) DeleteOneCar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		carID, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		deleteResult, err := c.carsRepository.DeleteOneCar(carID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Car with id %v was succefuly deleted", deleteResult)))
	}
}

func (c *carsHandler) GetOneCar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		carID, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := c.carsRepository.GetOneCar(carID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

func (c *carsHandler) UpdateOneCar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var result *model.Car

		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if result.ID == "" || result.Model == "" || result.Date == "" {
			err := errors.New("All fields must be filled")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := c.carsRepository.UpdateOneCar(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func (c *carsHandler) GetAllCars() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		result, err := c.carsRepository.GetAllCars()
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}
