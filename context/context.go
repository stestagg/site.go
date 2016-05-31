package context

import (
	"github.com/spf13/cast"
	"github.com/stestagg/site.go/log"
	)


type Context struct {
	next *Context
	Values map[string]interface{}
	Source string
}

func NewRootContext() Context {
	ctx := Context{}
	ctx.Values = make(map[string]interface{})
	ctx.Source = "root"
	return ctx
}

func (c *Context) NewEmptyContext(Source string) Context {
	ctx := Context{}
	ctx.next = c;
	ctx.Values = make(map[string]interface{})
	ctx.Source = Source
	return ctx
}

func (c *Context) NewForDir(name string) Context {
	ctx := c.NewEmptyContext(name)
	loadDirContext(&ctx, name)
	return ctx
}


// -----

func (c *Context) get(key string) interface{}{
	val, present := c.Values[key]
	if present {
		return val
	}
	if (c.next != nil) {
		return c.next.get(key)
	}
	return nil
}

func (c *Context) getAll(key string) []interface{} {
	cur_ctx := c
	out := make([]interface{}, 0)
	for cur_ctx != nil {
		val := cur_ctx.get(key)
		if val == nil {
			log.Debug("Masking(null) value for for key %s in %s", key, cur_ctx.Source)
			return out
		}
		out = append(out, val)
		cur_ctx = cur_ctx.next
	}
	return out
}

func (c *Context) Set(key string, value interface{}) {
	c.Values[key] = value
}

func (c *Context) GetString(key string) string {
	return cast.ToString(c.get(key))
}

func (c *Context) GetMap(key string) map[string]interface{} {
	return nil
}

func (c *Context) GetArray(key string) []interface{} {
	//values := c.getAll(key)
	out := make([]interface{}, 0)
	//for _, val := range(values) {

	//}
	return out
}