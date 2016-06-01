package site

import (
	"path"
	"sync"
	"time"

	"github.com/spf13/afero"

	"github.com/stestagg/site.go/context"
	"github.com/stestagg/site.go/log"
)

type finder (chan<- Node)

func matchesPattern(patterns []string, test_name string) bool {
	for _, pattern := range(patterns) {
		match, _ := path.Match(pattern, test_name)
		if match { return true }
	}
	return false
}

func visitDir(name string, context *context.Context, found finder, wg *sync.WaitGroup) {
	log.Debug("Visiting dir %s", name)
	defer wg.Done()
	entries, err := afero.ReadDir(Site.fs, name)
	if err != nil { log.Panic("Error %s reading directory: %s", err, name) }
	dirContext := context.NewForDirWithEntries(Site.fs, name, entries)
	ignores := dirContext.GetStringArray("site.ignore_patterns")

	for _, entry := range entries {
		if matchesPattern(ignores, entry.Name()) {
			log.Debug("File %s matches ignores, skipping", entry.Name())
			continue
		}
		child_name := path.Join(name, entry.Name())
		if entry.IsDir() {
			wg.Add(1)
			visitDir(child_name, &dirContext, found, wg)
		} else {
			wg.Add(1)
			visitFile(child_name, entry.Name(), &dirContext, found, wg)
		}
	}
}

func visitFile(path string, name string, context *context.Context, found finder, wg *sync.WaitGroup) {
	log.Debug("Visiting file %s", path)
	ctx := context.NewEmptyContext(path)
	defer wg.Done()
	node := NewNode(path, name, &ctx)
	if node.HasPipeline() {
		found <- node
	} else {
		log.Info("Skipping %s as no pipeline could be found", name)
	}
}


func (s site) DiscoverFiles() (<-chan Node) {
	var wg sync.WaitGroup
	found := make(chan Node)
	wg.Add(1)
	startTime := time.Now()
	go visitDir(Site.Root, &Site.Context, found, &wg)
	go func() {
		wg.Wait()
		log.Info("Discovery complete in %s μs", int(time.Since(startTime) / time.Microsecond))
		close(found)
	}()

	return found
}
