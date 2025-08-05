package parser_sql

type Word struct {
	LineNum int
	Value   string
}

type Words []*Word

type SQL struct {
	LastLineNum int
	Statements  []Statement
}

type Statement struct {
	LineNum   int
	QyeryType int
	Query     []*Token
}

type Token struct {
	LineNum int
	Value   string
}
