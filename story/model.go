package story

const (
	TopStories = iota
	NewStories
)

type Story struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
