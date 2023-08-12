package flax

import (
	"flag"
	"fmt"
	"os"
)

type Job interface {
	Run() error
	IsError() bool
}

type Flax struct {
	name    string
	ctx     *Context
	flagSet *flag.FlagSet
}

var allApp = []*Flax{}
var flagSet = flag.NewFlagSet("", flag.ContinueOnError)

func App(name string) *Flax {
	flax := &Flax{
		ctx: &Context{
			f:       func(c *Context) {},
			Err:     nil,
			Args:    make([]string, 0),
			typeMap: typeMap,
			rules:   make([]*Rule, 0),
			m:       map[string]struct{}{},
		},
		name:    name,
		flagSet: flagSet,
	}
	flax.ctx.flax = flax
	if name != "" {
		flax.flagSet = flag.NewFlagSet(name, flag.ContinueOnError)
	}
	allApp = append(allApp, flax)
	return flax
}

func (fx *Flax) Func(f func(c *Context)) *Flax {
	if fx.IsError() {
		return fx
	}
	fx.ctx.f = f
	return fx
}

func (f *Flax) Run() error {
	if f.IsError() {
		return f.ctx.Err
	}
	//为空则解析全局参数
	if f.flagSet == flagSet {
		flagSet.Parse(os.Args[1:])
		f.ctx.Args = flagSet.Args()
		//否则先检测本地参数，再检测全局参数
	} else {
		if len(os.Args) <= 1 {
			return fmt.Errorf("lack of args")
		}
		if os.Args[1] != f.name {
			return nil
		}
		f.flagSet.Parse(os.Args[2:])
		f.ctx.Args = f.flagSet.Args()
	}

	//若规则不满足则放弃执行
	if f.checkRules() {
		return nil
	}

	f.ctx.f(f.ctx)
	if f.IsError() {
		return f.ctx.Err
	}
	return nil
}

// 检查规则，若不满足则返回true
func (f *Flax) checkRules() bool {
	if len(f.ctx.rules) != 0 {
		if f.flagSet != nil {
			f.flagSet.Visit(func(fg *flag.Flag) {
				f.ctx.m[fg.Name] = struct{}{}
			})
		} else {
			flagSet.Visit(func(fg *flag.Flag) {
				f.ctx.m[fg.Name] = struct{}{}
			})
		}
	}
	for _, rule := range f.ctx.rules {
		//若规则不匹配则不执行
		if !rule.f(f.ctx) {
			return true
		}
	}
	return false
}

func (f *Flax) IsError() bool {
	return f.ctx.Err != nil
}

func Run() error {
	for _, v := range allApp {
		if v.IsError() {
			return v.ctx.Err
		}
		err := v.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
