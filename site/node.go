package site

import (
	"path"

	"github.com/spf13/cast"
	"github.com/stestagg/site.go/context"
	"github.com/stestagg/site.go/log"
)

type Node struct{

	Path string
	FileName string
	Context *context.Context
	Pipeline []interface{}
}

func getPipelineForFile(fileName string, ctx *context.Context) ([]interface{}){
	for _, candidate := range(ctx.GetArray("site.handlers")) {
		candidate_map, err := cast.ToStringMapE(candidate)
		if err != nil{
			log.Error("site.handlers must be a mapping with string keys, not: %s", candidate)
			continue
		}
		pattern, found := candidate_map["pattern"]
		if !found {
			log.Error("Handler %s is missing %s", candidate_map, "pattern")
			continue
		}
		pattern_str, cast_ok := pattern.(string)
		if !cast_ok {
			log.Error("Handler %s %s must be a string", candidate_map, "pattern")
			continue
		}
		log.Debug("Testing node %s against handler: %s", fileName, pattern_str)
		matched, _ := path.Match(pattern_str, fileName)
		if matched{
			pipeline_val, found := candidate_map["pipeline"]
			if !found {
				log.Error("Handler %s is missing %s", candidate_map, "pipeline")
				continue
			}
			pipeline, cast_ok := pipeline_val.([]interface{})
			if !cast_ok {
				log.Error("Handler %s %s is not a list", candidate_map, "pipeline")
				continue
			}
			log.Debug("Found pipeline %s for %s", pipeline, fileName)
			return pipeline
		}
	}
	return nil
}


func NewNode(filePath string, fileName string, ctx *context.Context) Node {
	return Node{
		Path: filePath,
		FileName: fileName,
		Context: ctx,
		Pipeline: getPipelineForFile(fileName, ctx),
	}
}


func (n *Node) HasPipeline() bool {
	return n.Pipeline != nil
}