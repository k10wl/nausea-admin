package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

type Firestore struct {
	client *firestore.Client
	ctx    context.Context
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

func (f *Firestore) docAbout() *firestore.DocumentRef {
	return f.client.Doc("data/about")
}

func (f *Firestore) docContacts() *firestore.DocumentRef {
	return f.client.Doc("data/contacts")
}

func (f *Firestore) colLinks() *firestore.CollectionRef {
	return f.client.Collection("links")
}

func (f *Firestore) collectionFolders() *firestore.CollectionRef {
	return f.client.Collection("folders")
}
