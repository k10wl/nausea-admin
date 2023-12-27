package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

type Firestore struct {
	client *firestore.Client
	ctx    context.Context
}

type Info struct {
	Bio string `firestore:"bio"`
}

type DB struct {
	firestore *Firestore
}

func NewFirestoreClient(projectID string) *Firestore {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		panic(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		panic(err)
	}
	return &Firestore{
		client: client,
		ctx:    ctx,
	}
}

func NewDB(projectID string) *DB {
	firestore := NewFirestoreClient(projectID)
	return &DB{
		firestore: firestore,
	}
}

func (f *Firestore) docInfo() *firestore.DocumentRef {
	return f.client.Doc("about/info")
}

func (db *DB) GetInfo() Info {
	var info Info
	doc, err := db.firestore.docInfo().Get(db.firestore.ctx)
	if err != nil {
		panic(err)
	}
	err = doc.DataTo(&info)
	fmt.Printf("info: %v\n", info)
	if err != nil {
		panic(err)
	}
	return info
}

func (db *DB) WriteInfo(info Info) error {
	_, err := db.firestore.docInfo().Set(db.firestore.ctx, info)
	return err
}
