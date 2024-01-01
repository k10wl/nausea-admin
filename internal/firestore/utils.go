package firestore

import (
	"encoding/json"

	"nausea-admin/internal/models"

	"cloud.google.com/go/firestore"
)

func structToUpdate(s interface{}) (*[]firestore.Update, error) {
	var m map[string]interface{}
	tmp, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(tmp, &m)
	updates := []firestore.Update{}
	for key, value := range m {
		updates = append(updates, firestore.Update{Path: key, Value: value})
	}
	return &updates, nil
}

func linkToFirestoreLink(l models.Link) FirestoreLink {
	return FirestoreLink{
		URL:  l.URL,
		Text: l.Text,
	}
}
