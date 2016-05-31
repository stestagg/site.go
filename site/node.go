package site

import (
	"github.com/stestagg/site.go/context"
	"github.com/stestagg/site.go/log"
)

type Node struct{

	Path string
	FileName string
	Context *context.Context
	Handlers []interface{}
}

func NewNode(Path string, FileName string, ctx *context.Context) Node {
	handlers := make([]interface{}, 0)
	for _, candidate := range(ctx.GetArray("site.handlers")) {
		candidate_map, was_cast := candidate.(map[string]interface{})
		if !was_cast {
			log.Error("site.handlers must be a mapping with string keys, not: %s", candidate)
			break
		}
	}
	return Node{
		Path: Path,
		FileName: FileName,
		Context: ctx,
		Handlers: handlers,
	}
}


func (n *Node) IsHandled() bool {
	return n.Handlers != nil
}