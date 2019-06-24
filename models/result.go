package models

//Result is a model to store data about scrapping
type Result struct {
	ErrorLinks        []Link
	VisitedLinks      []Link
	VisitedLinksCount int
	Duration          string
	MemoryUsage       string
}
