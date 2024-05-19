package db

import (
    "context"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    MessagesCollection *mongo.Collection
    Client             *mongo.Client
)

func ConnectMongoDB(logger *log.Logger) error {
    uri := os.Getenv("MONGODB_URI")
    dbName := os.Getenv("DB_NAME")
    username := os.Getenv("DB_USERNAME")
    password := os.Getenv("DB_PASSWORD")

    clientOptions := options.Client().ApplyURI(uri).SetAuth(options.Credential{
        Username: username,
        Password: password,
    })

    client, err := mongo.NewClient(clientOptions)
    if err != nil {
        logger.Fatal(err)
        return err
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    err = client.Connect(ctx)
    if err != nil {
        logger.Fatal(err)
        return err
    }

    Client = client
    MessagesCollection = client.Database(dbName).Collection("messages")
    return nil
}

func DisconnectMongoDB(logger *log.Logger) {
    if Client != nil {
        if err := Client.Disconnect(context.Background()); err != nil {
            logger.Println("Error disconnecting from MongoDB:", err)
        }
    }
}
