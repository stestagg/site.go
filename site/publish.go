package site

type site struct {
	fs afero.Fs
	Root string
	Context context.Context
}

var Site site;
