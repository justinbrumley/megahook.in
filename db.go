package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

type TokenNamespace struct {
	Token     string `json:"token"`
	Namespace string `json:"namespace"`
}

const (
	dbName         = "megahook"
	collectionName = "token_namespace"
)

var connstring string = "mongodb://megaman:SuperSecurePassword123@localhost:27017/megahook"

var dbClient *mongo.Client

func initDB() error {
	if os.Getenv("DB_CONN_STRING") != "" {
		connstring = os.Getenv("DB_CONN_STRING")
	}

	var err error
	dbClient, err = mongo.NewClient(options.Client().ApplyURI(connstring))
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = dbClient.Connect(ctx)
	if err != nil {
		return err
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = dbClient.Ping(ctx, readpref.Primary())
	return err
}

func getCollection() *mongo.Collection {
	return dbClient.Database(dbName).Collection(collectionName)
}

// Looks up token in db and returns namespace if found
func getTokenNamespace(token string) (*TokenNamespace, error) {
	collection := getCollection()

	result := &TokenNamespace{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := collection.FindOne(ctx, bson.D{{"token", token}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return result, nil
}

// Search for namespace by... namespace
func lookupNamespace(namespace string) (*TokenNamespace, error) {
	collection := getCollection()

	result := &TokenNamespace{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := collection.FindOne(ctx, bson.D{{"namespace", namespace}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return result, nil
}

// Create a new namespace record in db
func createNamespace(token string, namespace string) (*TokenNamespace, error) {
	collection := getCollection()

	ns := &TokenNamespace{
		Token:     token,
		Namespace: namespace,
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	_, err := collection.InsertOne(ctx, bson.D{{"token", token}, {"namespace", namespace}})
	if err != nil {
		return nil, err
	}

	return ns, nil
}

// Removes namespace from db by token
func deleteNamespace(token string) error {
	collection := getCollection()

	opts := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	_, err := collection.DeleteOne(ctx, bson.D{{"token", token}}, opts)
	return err
}
