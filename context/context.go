package context

import (
	"github.com/spf13/cast"

	"github.com/stestagg/site.go/log"
	)


type Context struct {
	next *Context
	Values map[string]interface{}
	Source string
	Final bool
}

func NewRootContext() Context {
	ctx := Context{}
	ctx.Values = DefaultContext
	ctx.Source = "root"
	ctx.Final = true
	return ctx
}

func (c *Context) NewEmptyContext(Source string) Context {
	ctx := Context{}
	ctx.next = c;
	ctx.Values = make(map[string]interface{})
	ctx.Source = Source
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
		val, present := cur_ctx.Values[key]
		if present {
			if val == nil {
				log.Debug("Masking(null) value for for key %s in %s", key, cur_ctx.Source)
				return out
			}
			out = append(out, val)
		}
		cur_ctx = cur_ctx.next
	}
	return out
}

func (c *Context) Set(key string, value interface{}) {
	if (c.Final) {
		log.Panic("Attempt to set key %s in read-only context %s", key, c.Source)
	}
	c.Values[key] = value
}

func (c *Context) GetString(key string) string {
	return cast.ToString(c.get(key))
}

func (c *Context) GetMap(key string) map[string]interface{} {
	return nil
}

func strSliceToSlice(val []string) []interface{} {
	out := make([]interface{}, len(val))
	for i, v := range(val) {
		out[i] = v
	}
	return out
}

func (c *Context) GetArray(key string) []interface{} {
	values := c.getAll(key)
	out := make([]interface{}, 0)
	for _, val := range(values) {
		switch v := val.(type) {
		case string:
			out = make([]interface{}, 1)
			out[0] = v
		case []interface{}:
			out = append(out, v...)
		case []string:
			out = append(out, strSliceToSlice(v)...)
		default:
			log.Panic("Value %s cant be interpreted as an array", val)
		}
	}
	return out
}

func (c *Context) GetStringArray(key string) []string {
	values := c.GetArray(key)
	out := make([]string, len(values))
	for i, v := range(values) {
		out[i] = cast.ToString(v)
	}
	return out
}