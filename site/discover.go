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

func visitDir(name string, context context.Context, found finder, wg *sync.WaitGroup) {
	log.Debug("Visiting dir %s", name)
	defer wg.Done()
	entries, err := afero.ReadDir(Site.fs, name)
	if err != nil { log.Panic("Error %s reading directory: %s", err, name) }

	for _, entry := range entries {
		child_name := path.Join(name, entry.Name())
		if entry.IsDir() {
			dirContext := context.NewForDir(name)
			wg.Add(1)
			visitDir(child_name, dirContext, found, wg)
		} else {
			wg.Add(1)
			visitFile(child_name, context, found, wg)
		}
	}
}

func visitFile(name string, context context.Context, found finder, wg *sync.WaitGroup) {
	log.Debug("Visiting file %s", name)
	ctx := context.NewEmptyContext(name)
	defer wg.Done()
	node := Node{
		Path: name,
		Context: &ctx,
	}
	found <- node
}


func (s site) DiscoverFiles() (<-chan Node) {
	var wg sync.WaitGroup
	found := make(chan Node)
	wg.Add(1)
	startTime := time.Now()
	go visitDir(Site.Root, Site.Context, found, &wg)
	go func() {
		wg.Wait()
		log.Info("Discovery complete in %s ms", int(time.Since(startTime) / time.Millisecond))
		close(found)
	}()

	return found
}
