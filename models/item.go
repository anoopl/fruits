package models

type Item struct {
	ID    int    `json:id`
	Title string `json:title`
	Owner string `json:owner`
	Year  string `json: date`
}
