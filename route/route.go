package route


type Route interface {
	Filter(url string) *Target
}
