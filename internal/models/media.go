package models

type MediaSize struct {
	Width  int `firestore:"width"`
	Height int `firestore:"height"`
}

type Media struct {
	ID
	URL          string `firestore:"URL"`
	ThumbnailURL string `firestore:"thumbnailURL"`
	MediaSize
	Timestamps
}

func NewMedia(URL string, mediaSize MediaSize) (Media, error) {
	m := Media{URL: URL, Timestamps: NewTimestamps(), MediaSize: mediaSize}
	err := m.generateID()
	return m, err
}

func (m Media) AsContent(parentID string) (MediaContent, error) {
	mc := MediaContent{
		ContentBase:  ContentBase{RefID: m.ID.ID, Timestamps: NewTimestamps()},
		MediaSize:    m.MediaSize,
		Name:         m.ID.ID,
		URL:          m.URL,
		ThumbnailURL: m.ThumbnailURL,
		ParentID:     parentID,
		Description:  "",
	}
	err := mc.generateID()
	return mc, err
}
