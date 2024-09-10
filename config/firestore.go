package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/storage"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App
var storageClient *storage.Client

func FirestoreInit() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("./ServiceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	firebaseApp = app

	//intialize the firebase storage client
	client, err := firebaseApp.Storage(ctx)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	storageClient = client
}
