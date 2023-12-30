package firestore

import (
	"context"

	"nausea-admin/internal/models"

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

func (f *Firestore) docInfo() *firestore.DocumentRef {
	return f.client.Doc("about/info")
}

func (f *Firestore) GetInfo() (*models.Info, error) {
	var info models.Info
	doc, err := f.docInfo().Get(f.ctx)
	if err != nil {
		return nil, err
	}
	err = doc.DataTo(&info)
	return &info, err
}

func (f *Firestore) WriteInfo(info models.Info) error {
	_, err := f.docInfo().Set(f.ctx, info)
	return err
}
