package models


import (
	"time"
	pq "github.com/lib/pq"
)


type IdResponse struct {
	BannerId int `json:"id"`
}

type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
    Url   string `json:"url"`
}

type Banner struct {
	TagIds 		[]int 	`json:"tag_ids"`
	FeatureId 	int 	`json:"feature_id"`
	Content 	Content `json:"content"`
	IsActive	bool	`json:"is_actitve"`
}

type UpdateContent struct {
	Title *string `json:"title"`
	Text  *string `json:"text"`
    Url   *string `json:"url"`
}

type UpdateBanner struct {
	TagIds 		*[]int 				`json:"tag_ids,omitempty"`
	FeatureId 	*int 	 		 	`json:"feature_id,omitempty"`
	Content 	UpdateContent 		`json:"content,omitempty"`
	IsActive	*bool	 			`json:"is_active,omitempty"`
}

type BannerDB struct {
    ID              int       		`db:"id"`
    TagsIDs         pq.Int64Array   `db:"tags_ids"`
    FeatureID       int       		`db:"feature_id"`
    Title           string    		`db:"title"`
    Text            string    		`db:"text"`
    URL             string    		`db:"url"`
    IsActive        bool      		`db:"is_active"`
    CreateDatetime  time.Time 		`db:"create_datetime"`
    UpdateDatetime  time.Time 		`db:"update_datetime"`
}
