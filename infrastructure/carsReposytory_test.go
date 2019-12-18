package infrastructure

import (
	"errors"
	"server/infrastructure/client"
	"server/infrastructure/collection"
	"server/infrastructure/model"
	"server/mocks"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	idMongo, _ = primitive.ObjectIDFromHex("5dee1d9768c592130a6c56ed")

	car = &model.Car{}

	errFoo = errors.New("Some error")
)

func MockRepository(err error) (CarsRepository, *mocks.Client, *mocks.Collection) {
	clientMock := &mocks.Client{}
	collectionMock := &mocks.Collection{}
	clientMockFunc := client.NewClientFactoryFuncMock(clientMock, err)
	collectionMockFunc := collection.NewCollectionFactoryFuncMock(collectionMock)
	repo := NewCarsRepository(clientMockFunc, collectionMockFunc)
	return repo, clientMock, collectionMock
}
func TestCarsRepositoryGetAllCarError(t *testing.T) {
	t.Parallel()
	repo, clientMock, collectionMock := MockRepository(nil)
	clientMock.On("Database", nil).Return(nil)
	collectionMock.On("Find", mock.Anything, mock.Anything, mock.Anything).Return(nil, errFoo)
	_, err := repo.GetAllCars()
	assert.Equal(t, errFoo, err)
}
func TestCarsRepositoryGetAllCarClientError(t *testing.T) {
	t.Parallel()
	repo, _, _ := MockRepository(errFoo)
	_, err := repo.GetAllCars()
	assert.Equal(t, errFoo, err)
}
func TestCarsRepositoryInsertOneCar(t *testing.T) {
	t.Parallel()
	repo, clientMock, collectionMock := MockRepository(nil)
	clientMock.On("Disconnect", mock.Anything).Return(nil)
	collectionMock.On("InsertOne", mock.Anything, mock.Anything).Return(nil, nil)
	_, err := repo.InsertOneCar(car)
	assert.Equal(t, nil, err)
}
func TestCarsRepositoryInsertOneCarError(t *testing.T) {
	t.Parallel()
	repo, _, collectionMock := MockRepository(nil)
	collectionMock.On("InsertOne", mock.Anything, mock.Anything).Return(nil, errFoo)
	_, err := repo.InsertOneCar(car)
	assert.Equal(t, errFoo, err)
}
func TestCarsRepositoryInsertOneCarClientError(t *testing.T) {
	t.Parallel()
	repo, _, collectionMock := MockRepository(errFoo)
	collectionMock.On("InsertOne", mock.Anything, mock.Anything).Return(nil, nil)
	_, err := repo.InsertOneCar(car)
	assert.Equal(t, errFoo, err)
}
func TestCarsRepositoryInsertOneCarDisconnectError(t *testing.T) {
	t.Parallel()
	repo, clientMock, collectionMock := MockRepository(nil)
	clientMock.On("Disconnect", mock.Anything).Return(errFoo)
	collectionMock.On("InsertOne", mock.Anything, mock.Anything).Return(nil, nil)
	_, err := repo.InsertOneCar(car)
	assert.Equal(t, errFoo, err)
}
func TestCarsRepositoryDeleteOneCar(t *testing.T) {
	t.Parallel()
	repo, clientMock, collectionMock := MockRepository(nil)
	clientMock.On("Disconnect", mock.Anything).Return(nil)
	collectionMock.On("DeleteOne", mock.Anything, mock.Anything).Return(&mongo.DeleteResult{DeletedCount: 1}, nil)
	_, err := repo.DeleteOneCar(idMongo)
	assert.Equal(t, nil, err)
}
func TestCarsRepositoryDeleteOneCarError(t *testing.T) {
	t.Parallel()
	repo, _, collectionMock := MockRepository(nil)
	expectedError := errors.New("wrong data in request")
	collectionMock.On("DeleteOne", mock.Anything, mock.Anything).Return(&mongo.DeleteResult{DeletedCount: 1}, expectedError)
	_, err := repo.DeleteOneCar(idMongo)
	assert.Equal(t, expectedError, err)
}
func TestCarsRepositoryDeleteOneCarClientError(t *testing.T) {
	t.Parallel()
	repo, _, collectionMock := MockRepository(errFoo)
	collectionMock.On("DeleteOne", mock.Anything, mock.Anything).Return(&mongo.DeleteResult{DeletedCount: 1}, nil)
	_, err := repo.DeleteOneCar(idMongo)
	assert.Equal(t, errFoo, err)
}
func TestCarsRepositoryGetOneCar(t *testing.T) {
	t.Parallel()
	repo, clientMock, collectionMock := MockRepository(nil)
	clientMock.On("Disconnect", mock.Anything).Return(nil)
	expectedError := errors.New("Registry cannot be nil")
	collectionMock.On("FindOne", mock.Anything, mock.Anything).Return(&mongo.SingleResult{}, nil)
	_, err := repo.GetOneCar(idMongo)
	assert.Equal(t, expectedError, err)
}
func TestCarsRepositoryGetOneCarClientError(t *testing.T) {
	t.Parallel()
	repo, _, collectionMock := MockRepository(errFoo)
	collectionMock.On("FindOne", mock.Anything, mock.Anything).Return(nil, nil)
	_, err := repo.GetOneCar(idMongo)
	assert.Equal(t, errFoo, err)
}

func TestCarsRepositoryUpdateOneCar(t *testing.T) {
	t.Parallel()
	repo, clientMock, collectionMock := MockRepository(nil)
	clientMock.On("Disconnect", mock.Anything).Return(nil)
	collectionMock.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{MatchedCount: 1}, nil)
	_, err := repo.UpdateOneCar(car)
	assert.Equal(t, nil, err)
}
func TestCarsRepositoryUpdateOneCarError(t *testing.T) {
	t.Parallel()
	repo, _, collectionMock := MockRepository(nil)
	collectionMock.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, errFoo)
	_, err := repo.UpdateOneCar(car)
	assert.Equal(t, errFoo, err)
}
func TestCarsRepositoryUpdateOneCarClientError(t *testing.T) {
	t.Parallel()
	repo, _, collectionMock := MockRepository(errFoo)
	collectionMock.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, errFoo)
	_, err := repo.UpdateOneCar(car)
	assert.Equal(t, errFoo, err)
}
func TestCarsRepositoryUpdateOneCarDisconnectError(t *testing.T) {
	t.Parallel()
	repo, clientMock, collectionMock := MockRepository(nil)
	clientMock.On("Disconnect", mock.Anything).Return(errFoo)
	collectionMock.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{MatchedCount: 1}, nil)
	_, err := repo.UpdateOneCar(car)
	assert.Equal(t, errFoo, err)
}
