package collections

import "go.mongodb.org/mongo-driver/mongo"

//getting database collections
func GetCollection(client *mongo.Client, collectionName string, dbName string) *mongo.Collection {
    collection := client.Database(dbName).Collection(collectionName)
    return collection
}