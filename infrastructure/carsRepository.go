package infrastructure

import (
	"context"
	"errors"
	"server/infrastructure/client"
	"server/infrastructure/collection"
	"server/infrastructure/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

//CarsRepository ...
type CarsRepository interface {
	GetAllCars() (*[]model.Car, error)
	InsertOneCar(car *model.Car) (*mongo.InsertOneResult, error)
	DeleteOneCar(carID primitive.ObjectID) (primitive.ObjectID, error)
	GetOneCar(carID primitive.ObjectID) (*model.Car, error)
	UpdateOneCar(cars *model.Car) (*mongo.UpdateResult, error)
}

type carsRepository struct {
	collectionFactoryFunc collection.FactoryFuncC
	clientFactoryFunc     client.FactoryFunc
}

// NewCarsRepository ...
func NewCarsRepository(clientFactoryFunc client.FactoryFunc, collectionFactoryFunc collection.FactoryFuncC) CarsRepository {
	return &carsRepository{
		collectionFactoryFunc: collectionFactoryFunc,
		clientFactoryFunc:     clientFactoryFunc,
	}
}

func (repo *carsRepository) GetAllCars() (*[]model.Car, error) {
	client, err := repo.clientFactoryFunc()
	if err != nil {
		return nil, err
	}

	collection := repo.collectionFactoryFunc(client)

	var allDocuments []model.Car

	cur, err := collection.Find(context.Background(), bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var elem model.Car
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		allDocuments = append(allDocuments, elem)
	}

	if err := client.Disconnect(context.Background()); err != nil {
		return nil, err
	}

	return &allDocuments, err
}

func (repo *carsRepository) InsertOneCar(car *model.Car) (*mongo.InsertOneResult, error) {
	client, err := repo.clientFactoryFunc()
	if err != nil {
		return nil, err
	}

	collection := repo.collectionFactoryFunc(client)

	filter := bson.M{
		"_id":   car.IDMongo,
		"ID":    car.ID,
		"Model": car.Model,
		"Date":  car.Date,
		// "Membership": bson.M{
		// 	"groupName": car.Membership.GroupName,
		// 	"groupID":   car.Membership.GroupID,
		// },
	}

	insertResult, err := collection.InsertOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err := client.Disconnect(context.Background()); err != nil {
		return nil, err
	}

	return insertResult, nil
}

func (repo *carsRepository) DeleteOneCar(carID primitive.ObjectID) (primitive.ObjectID, error) {
	client, err := repo.clientFactoryFunc()
	if err != nil {
		return carID, err
	}

	collection := repo.collectionFactoryFunc(client)

	deleteResult, err := collection.DeleteOne(context.Background(), bson.M{"_id": carID})
	if err != nil {
		return carID, errors.New("wrong data in request")
	}

	if deleteResult.DeletedCount == 0 {
		return carID, errors.New("object not found in DB")
	}

	if err := client.Disconnect(context.Background()); err != nil {
		return carID, err
	}

	return carID, nil
}

func (repo *carsRepository) GetOneCar(carID primitive.ObjectID) (*model.Car, error) {
	client, err := repo.clientFactoryFunc()
	if err != nil {
		return nil, err
	}

	collection := repo.collectionFactoryFunc(client)

	var elem model.Car

	if err := collection.FindOne(context.TODO(), bson.M{"_id": carID}).Decode(&elem); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Document not found")
		}
		return nil, err
	}

	if err := client.Disconnect(context.Background()); err != nil {
		return nil, err
	}

	return &elem, nil
}

func (repo *carsRepository) UpdateOneCar(cars *model.Car) (*mongo.UpdateResult, error) {
	client, err := repo.clientFactoryFunc()
	if err != nil {
		return nil, err
	}

	collection := repo.collectionFactoryFunc(client)

	filter := bson.M{"_id": bson.M{"$eq": cars.IDMongo}}
	update := bson.M{
		"$set": bson.M{
			"ID":    cars.ID,
			"Model": cars.Model,
			"Date":  cars.Date,
			// "Membership": bson.M{
			// 	"groupName": cars.Membership.GroupName,
			// 	"groupID":   cars.Membership.GroupID,
			// },
		},
	}

	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	if err := client.Disconnect(context.Background()); err != nil {
		return nil, err
	}

	return updateResult, nil
}
