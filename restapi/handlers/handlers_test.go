package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"server/infrastructure/model"
	"server/mocks"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	carID = "5dee1d9768c592130a6c56ed"

	err = errors.New("Some error")

	car = []byte(`{
	"id": "5dee1d9768c592130a6c56ed",
	"ID": "1155 FD-6",
	"Model": "GAZ",
	"Date": "21.01.09"
}`)

	badcar = []byte(`{
	"id": "5dee1d9768c592130a6c56ed",
	"Model": "GAZ",
	"Date": "21.01.09"
}`)

	idMongo, _ = primitive.ObjectIDFromHex("5dee1d9768c592130a6c56ed")

	mcar = &model.Car{
		IDMongo: idMongo,
		ID:      "1155 FD-6", //`bson:"id" json:"id,omitempty"`
		Model:   "GAZ",       //`bson:"model" json:"model,omitempty"`
		Date:    "21.01.09",  //`bson:"data" json:"data,omitempty"`
	}
)

func mockRepository() CarsHandler {
	carsRepoMock := &mocks.CarsRepository{}
	return NewCarsHandler(carsRepoMock)
}

func TestIndexHandler(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/index", rr.Body)
	handler := mockRepository().Index()
	io.WriteString(rr, "<html><body>Hello World!</body></html>")
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
}
func TestInsertOneCarError500(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("POST", "/car", bytes.NewBuffer(car))
	rr := httptest.NewRecorder()
	carsRepoMock.On("InsertOneCar", mock.Anything).Return(nil, err)
	userHandler := NewCarsHandler(carsRepoMock)
	handler := userHandler.InsertOneCar()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}
func TestInsertOneCarError400(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("POST", "/car", bytes.NewBuffer(badcar))
	rr := httptest.NewRecorder()
	carsRepoMock.On("InsertOneCar", mock.Anything).Return(nil, nil)
	userHandler := NewCarsHandler(carsRepoMock)
	handler := userHandler.InsertOneCar()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusBadRequest)
}
func TestInsertOneCarOK(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("POST", "/car", bytes.NewBuffer(car))
	rr := httptest.NewRecorder()
	carsRepoMock.On("InsertOneCar", mock.Anything).Return(nil, nil)
	userHandler := NewCarsHandler(carsRepoMock)
	handler := userHandler.InsertOneCar()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusCreated)
}
func TestDeleteOneCarOK(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("DELETE", "/", bytes.NewBuffer(car))
	rr := httptest.NewRecorder()
	vars := map[string]string{"id": carID}
	req = mux.SetURLVars(req, vars)
	carsRepoMock.On("DeleteOneCar", idMongo).Return(nil, nil)
	userHandler := NewCarsHandler(carsRepoMock)
	handler := userHandler.DeleteOneCar()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
}
func TestDeleteOneCarError404(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("DELETE", "/", bytes.NewBuffer(car))
	rr := httptest.NewRecorder()
	vars := map[string]string{"id": carID}
	req = mux.SetURLVars(req, vars)
	carsRepoMock.On("DeleteOneCar", idMongo).Return(nil, err)
	userHandler := NewCarsHandler(carsRepoMock)
	handler := userHandler.DeleteOneCar()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusNotFound)
}
func TestDeleteOneCarError400(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("DELETE", "/", nil)
	rr := httptest.NewRecorder()
	vars := map[string]string{"id": "111"}
	req = mux.SetURLVars(req, vars)
	carsRepoMock.On("DeleteOneCar", idMongo).Return(nil, nil)
	userHandler := NewCarsHandler(carsRepoMock)
	handler := userHandler.DeleteOneCar()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusBadRequest)
}
func TestGetOneCarOK(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("Get", "/"+carID, nil)
	rr := httptest.NewRecorder()
	carsRepoMock.On("GetOneCar", idMongo).Return(mcar, nil)
	userHandler := NewCarsHandler(carsRepoMock)
	vars := map[string]string{"id": carID}
	req = mux.SetURLVars(req, vars)
	handler := userHandler.GetOneCar()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
}
func TestGetOneCarError404(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("Get", "/"+carID, nil)
	rr := httptest.NewRecorder()
	carsRepoMock.On("GetOneCar", idMongo).Return(nil, err)
	userHandler := NewCarsHandler(carsRepoMock)
	vars := map[string]string{"id": carID}
	req = mux.SetURLVars(req, vars)
	handler := userHandler.GetOneCar()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusNotFound)
}
func TestGetOneCarError400(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("Get", "/"+carID, nil)
	rr := httptest.NewRecorder()
	carsRepoMock.On("GetOneCar", idMongo).Return(nil, err)
	userHandler := NewCarsHandler(carsRepoMock)
	vars := map[string]string{"id": "111"}
	req = mux.SetURLVars(req, vars)
	handler := userHandler.GetOneCar()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusBadRequest)
}
func TestUpdateOneCarOK(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("POST", "/car", bytes.NewBuffer(car))
	rr := httptest.NewRecorder()
	carsRepoMock.On("UpdateOneCar", mcar).Return(nil, nil)
	userHandler := NewCarsHandler(carsRepoMock)
	handler := userHandler.UpdateOneCar()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
}
func TestUpdateOneCarError400(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("POST", "/car", bytes.NewBuffer(badcar))
	rr := httptest.NewRecorder()
	carsRepoMock.On("UpdateOneCar", mcar).Return(nil, nil)
	userHandler := NewCarsHandler(carsRepoMock)
	handler := userHandler.UpdateOneCar()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusBadRequest)
}
func TestUpdateOneCarError500(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("POST", "/car", bytes.NewBuffer(car))
	rr := httptest.NewRecorder()
	carsRepoMock.On("UpdateOneCar", mcar).Return(nil, err)
	userHandler := NewCarsHandler(carsRepoMock)
	handler := userHandler.UpdateOneCar()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}
func TestGetAllCarsOK(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	carsRepoMock.On("GetAllCars").Return(nil, nil)
	userHandler := NewCarsHandler(carsRepoMock)
	handler := userHandler.GetAllCars()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
}
func TestGetAllCarsError404(t *testing.T) {
	t.Parallel()
	carsRepoMock := &mocks.CarsRepository{}
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	carsRepoMock.On("GetAllCars").Return(nil, err)
	userHandler := NewCarsHandler(carsRepoMock)
	handler := userHandler.GetAllCars()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusNotFound)
}
