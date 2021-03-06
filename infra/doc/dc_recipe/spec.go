package dc_recipe

type DocColumn interface {
	ColName() string
	ColDesc() string
	ColExample() string
}

type DocRows interface {
	RowsName() string
	RowsDesc() string
	RowsCols() []DocColumn
	RowsHasExample() bool
}
