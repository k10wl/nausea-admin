package models

type Media struct {
	ID
	URL string `firestore:"URL"`
	Timestamps
}

func NewMedia(URL string) (Media, error) {
	m := Media{URL: URL, Timestamps: NewTimestamps()}
	err := m.generateID()
	return m, err
}

func (m Media) AsContent() (MediaContent, error) {
	mc := MediaContent{
		ContentBase: ContentBase{RefID: m.ID.ID, Timestamps: NewTimestamps()},
		Name:        m.ID.ID,
		Description: "",
	}
	err := mc.generateID()
	return mc, err
}
