package vm3

import "fmt"

type LiteralKind int

const (
	KString LiteralKind = iota
	KInt
	KFloat
)

func (lk LiteralKind) String() string {
	switch lk {
	case KString:
		return "KString"
	case KInt:
		return "KInt"
	case KFloat:
		return "KFloat"
	default:
		return "illegal"
	}
}

type Literal struct {
	kind LiteralKind
	s    string
	i    int
	f    float64
}

func NewLiteral[T string | int | float64](v T) *Literal {
	switch any(v).(type) {
	case string:
		return &Literal{
			kind: KString,
			s:    any(v).(string),
		}
	case int:
		return &Literal{
			kind: KInt,
			i:    any(v).(int),
		}
	case float64:
		return &Literal{
			kind: KFloat,
			f:    any(v).(float64),
		}
	}
	return nil
}

func (l *Literal) String() string {
	var v any
	switch l.kind {
	case KString:
		v = l.s
	case KInt:
		v = l.i
	case KFloat:
		v = l.f
	}
	return fmt.Sprintf("Literal{ kind: %s, value: %v }", l.kind.String(), v)
}

func (l *Literal) GetKind() LiteralKind {
	return l.kind
}

func (l *Literal) GetString() string {
	return l.s
}
func (l *Literal) SetString(s string) {
	l.s = s
}

func (l *Literal) GetInt() int {
	return l.i
}
func (l *Literal) SetInt(i int) {
	l.i = i
}

func (l *Literal) GetFloat() float64 {
	return l.f
}
func (l *Literal) SetFloat(f float64) {
	l.f = f
}
