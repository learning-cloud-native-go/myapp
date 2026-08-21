package form

type BookForm struct {
	Title         string `json:"title" validate:"required,max=255"`
	PublishedDate string `json:"published_date" validate:"required,datetime=2006-01-02"`
	Description   string `json:"description"`
	ImageURL      string `json:"image_url" validate:"url"`
	Status        string `json:"status" validate:"required,oneof=pending verified"`
}
