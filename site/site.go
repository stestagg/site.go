package site

import (
	"os"

	"github.com/spf13/afero"

	"github.com/stestagg/site.go/context"
	)

type site struct {
	fs afero.Fs
	Root string
	Context context.Context
}

var Site site;


func (site) SetRoot(root string){
	Site.Root = root
}


func init() {
	Site.fs = afero.NewOsFs()
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	Site.Root =wd
}