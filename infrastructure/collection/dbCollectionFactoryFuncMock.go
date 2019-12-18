package collection

import "server/infrastructure/client"

//NewCollectionFactoryFuncMock ...
func NewCollectionFactoryFuncMock(collectionMock Collection) func(client.Client) Collection {
	return func(client.Client) Collection {
		return collectionMock
	}
}
