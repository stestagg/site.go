package context

import (
	"os"
	"path"
	"gopkg.in/yaml.v2"
	"github.com/spf13/afero"
//	"github.com/spf13/cast"

	"github.com/stestagg/site.go/log"
	)

func (c *Context) NewForDirWithEntries(fs afero.Fs, dir_path string, entries []os.FileInfo) Context {
	ctx := c.NewEmptyContext(dir_path)
	file_patterns := ctx.GetStringArray("site.context_pattern")
	log.Debug("Looking for context files matching: %s", file_patterns)
	if len(file_patterns) == 0{
		return ctx
	}
	for _, entry := range entries {
		file_name := entry.Name()
		for _, pattern := range(file_patterns) {
			match, _ := path.Match(pattern, file_name)
			if  match {
				file_path := path.Join(dir_path, file_name)
				contents, err := afero.ReadFile(fs, file_path)
				if err != nil {
					log.Panic("Error reading file %s: %s", file_path, err)
				}
				loadYamlContext(&ctx, file_path, contents)
				break
			}
		}
	}
	return ctx
}


func loadYamlContext(ctx *Context, path string, data []byte) {
	contents := make(map[string]interface{})
	err := yaml.Unmarshal(data, contents)
	if err != nil {
		log.Error("Could not parse %s as yaml: %s", path, err)
		return
	}
	log.Debug("Found context in %s: %s", ctx.Source, contents)
	ctx.Values = contents
}
