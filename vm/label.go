package vm

type Label struct {
	Id     string
	Define bool
}

func (l *Label) String() string {
	if l.Define {
		return l.Id + ":"
	}
	return l.Id
}
