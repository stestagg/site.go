package log

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cast"

)

var Verbosity int = 2

var ArgColor func(a ...interface{}) string = color.New(color.FgHiBlue).SprintFunc()

var InfoColor *color.Color = color.New(color.FgGreen)
var DebugColor *color.Color = color.New(color.FgCyan)
var WarnColor *color.Color = color.New(color.FgHiYellow)
var ErrColor *color.Color = color.New(color.BgRed).Add(color.Bold).Add(color.FgWhite)

func shouldSkip(level int) bool {
	return level > Verbosity
}

func wrapArgs(args []interface{}) ([]interface{}){
	formattedArgs := make([]interface{}, 0)
	for i := range args{
		arg := args[i]
		argstr := cast.ToString(arg)
		formattedArgs = append(formattedArgs, ArgColor(argstr))
	}
	return formattedArgs
}

func doOutput(level int, tag_color *color.Color, tag string, msg string, args []interface{}) {
	if shouldSkip(level) { return }
	wrappedArgs := wrapArgs(args)
	msg = fmt.Sprintf("[%s] %s\n", tag_color.SprintFunc()(tag), msg)
	fmt.Fprintf(color.Output, msg, wrappedArgs...)
}

func Info(msg string, args ...interface{}) {
	doOutput(2, InfoColor, "Info", msg, args)
}

func Debug(msg string, args ...interface{}) {
	doOutput(3, DebugColor, "Debug", msg, args)
}

func Warn(msg string, args ...interface{}) {
	doOutput(1, WarnColor, "Warn", msg, args)
}

func Error(msg string, args ...interface{}) {
	doOutput(0, ErrColor, "Error", msg, args)
}

func Panic(msg string, args ...interface{}) {
	doOutput(0, ErrColor, "Error", msg, args)
	if Verbosity > 5 {
		panic("")
	}
	os.Exit(1)
}