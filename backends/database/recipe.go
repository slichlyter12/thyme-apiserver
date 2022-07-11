package database

// Recipe that users can create
type Recipe struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Author      string            `json:"author"`
	Description string            `json:"description"`
	Cuisine     string            `json:"cuisine"`
	ImageName   string            `json:"imageName"`
	Ingredients map[string]string `json:"ingredients"`
	Steps       []string          `json:"steps"`
}
