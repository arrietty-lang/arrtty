package tokenize

type Position struct {
	LineNo int
	Lat    int
	Wat    int
}

func NewPosition(ln, lat, wat int) *Position {
	return &Position{
		LineNo: ln,
		Lat:    lat,
		Wat:    wat,
	}
}
