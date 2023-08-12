package flax

import (
	"fmt"
	"strconv"
)

type Context struct {
	f       func(c *Context)
	rules   []*Rule
	Err     error
	Args    []string
	typeMap map[string]uint
	flax    *Flax
	m       map[string]struct{}
}

// return if not find
func (c *Context) Get(name string) any {
	if fg := c.flax.flagSet.Lookup(name); fg != nil {
		str := fg.Value.String()
		switch c.typeMap[name] {
		case String:
			return str
		case Int:
			i, _ := strconv.Atoi(str)
			return i
		case Int64:
			i, _ := strconv.Atoi(str)
			return int64(i)
		case Bool:
			i, _ := strconv.ParseBool(str)
			return i
		case Uint:
			i, _ := strconv.Atoi(str)
			return uint(i)
		case Uint64:
			i, _ := strconv.Atoi(str)
			return uint64(i)
		case Float64:
			i, _ := strconv.ParseFloat(str, 64)
			return i
		}
	}
	c.Err = fmt.Errorf("can not get flag value: %v", name)
	return nil
}

func (c *Context) Argc() int {
	return len(c.Args)
}

func (c *Context) Error(err error) {
	c.Err = err
}

func (c *Context) Errorf(format string, a ...any) {
	c.Err = fmt.Errorf(format, a...)
}
