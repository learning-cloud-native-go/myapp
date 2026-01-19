package form

type BookForm struct {
	Title         string `json:"title" validate:"required,max=255"`
	Author        string `json:"author" validate:"required,alpha_space,max=255"`
	PublishedDate string `json:"published_date" validate:"required,datetime=2006-01-02"`
	ImageURL      string `json:"image_url" validate:"url"`
	Description   string `json:"description"`
}
