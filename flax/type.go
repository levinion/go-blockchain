package flax

var typeMap = map[string]uint{}

func (f *Flax) Var(t uint, name string, value any, usage string) *Flax {
	f.ctx.typeMap[name] = t
	switch t {
	case String:
		f.flagSet.String(name, value.(string), usage)
	case Bool:
		f.flagSet.Bool(name, value.(bool), usage)
	case Float64:
		f.flagSet.Float64(name, value.(float64), usage)
	case Int:
		f.flagSet.Int(name, value.(int), usage)
	case Int64:
		f.flagSet.Int64(name, value.(int64), usage)
	case Uint:
		f.flagSet.Uint(name, value.(uint), usage)
	case Uint64:
		f.flagSet.Uint64(name, value.(uint64), usage)
	}
	return f
}

func Var(t uint, name string, value any, usage string) {
	typeMap[name] = t
	switch t {
	case String:
		flagSet.String(name, value.(string), usage)
	case Bool:
		flagSet.Bool(name, value.(bool), usage)
	case Float64:
		flagSet.Float64(name, value.(float64), usage)
	case Int:
		flagSet.Int(name, value.(int), usage)
	case Int64:
		flagSet.Int64(name, value.(int64), usage)
	case Uint:
		flagSet.Uint(name, value.(uint), usage)
	case Uint64:
		flagSet.Uint64(name, value.(uint64), usage)
	}
}

const (
	String  = 24
	Bool    = 1
	Float64 = 14
	Int     = 2
	Int64   = 6
	Uint    = 7
	Uint64  = 11
)
