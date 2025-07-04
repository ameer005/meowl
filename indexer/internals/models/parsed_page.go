package models

type ParsedDocument struct {
	Url           string
	Title         string
	Description   string
	Tokens        []string
	internalLinks []string
	extrenalLinks []string
}

