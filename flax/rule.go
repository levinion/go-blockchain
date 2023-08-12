package flax

type Rule struct {
	f func(c *Context) bool
}

func (fx *Flax) Rule(f func(c *Context) bool) *Flax {
	fx.ctx.rules = append(fx.ctx.rules, &Rule{f})
	return fx
}

func (fx *Flax) MustHave(names ...string) *Flax {
	fx.Rule(func(c *Context) bool {
		for _, name := range names {
			if !c.Exists(name) {
				return false
			}
		}
		return true
	})
	return fx
}

func (f *Flax) ExactArgs(num int) *Flax {
	f.Rule(func(c *Context) bool {
		return c.ExactArgs(num)
	})
	return f
}

func (f *Flax) RangeArgs(min int, max int) *Flax {
	f.Rule(func(c *Context) bool {
		return c.RangeArgs(min, max)
	})
	return f
}

func (c *Context) Exists(name string) bool {
	if _, ok := c.m[name]; ok {
		return true
	}
	return false
}

func (c *Context) ExactArgs(num int) bool {
	return c.Argc() == num
}

func (c *Context) RangeArgs(min int, max int) bool {
	return c.Argc() <= max && c.Argc() >= min
}
