package context

var ignore_patterns = []string{".*", "*.pyc", "sitego.yaml"}

var DefaultContext = map[string]interface{} {
	"site.context_pattern": "sitego.yaml",
	"site.ignore_patterns": ignore_patterns,
}