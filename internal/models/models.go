package models

type Link struct {
	ID   string
	URL  string `firestore:"url" json:"url"`
	Text string `firestore:"text" json:"text"`
}

type Contacts struct {
	Email string `firestore:"email" json:"email,omitempty"`
	Links []Link `firestore:"links" json:"links,omitempty"`
}

type About struct {
	Bio string `firestore:"bio"`
}
