package models

type NavBarStructure struct {
	StartButtons []Button
	EndButtons   []Button
}

type Button struct {
	ButtonName   string
	Url          string
	NotHxRequest bool
	Childs       []Button
}
