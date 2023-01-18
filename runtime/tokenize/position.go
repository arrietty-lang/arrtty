package tokenize

type Position struct {
	LineNo int
	Lat    int
	Wat    int
}

func (p *Position) Clone() *Position {
	return &Position{
		LineNo: p.LineNo,
		Lat:    p.Lat,
		Wat:    p.Wat,
	}
}

func NewPosition(ln, lat, wat int) *Position {
	return &Position{
		LineNo: ln,
		Lat:    lat,
		Wat:    wat,
	}
}
