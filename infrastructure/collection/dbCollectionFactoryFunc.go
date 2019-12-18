package collection

import (
	"os"
	"server/infrastructure/client"
)

var (
	databaseName   = getEnv("DATABASE_NAME", "vehicle_inspection")
	collectionName = getEnv("DATABASE_COLLECTION", "cars")
)

//FactoryFuncC ...
type FactoryFuncC func(c client.Client) Collection

//NewCollection ...
func NewCollection(c client.Client) Collection {
	collection := c.Database(databaseName).Collection(collectionName)
	return collection
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
