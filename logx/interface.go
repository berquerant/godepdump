package logx

import "io"

func Setup(w io.Writer, debug bool) {
	setup(w)
	if debug {
		SetLevel(Ldebug)
	} else {
		SetLevel(Linfo)
	}
}

func SetLevel(level Level) { levelVar.Set(level.intoSlog()) }

func Debug(msg string, v ...Attr) { getInstance().debug(msg, v...) }
func Info(msg string, v ...Attr)  { getInstance().info(msg, v...) }
func Error(msg string, v ...Attr) { getInstance().error(msg, v...) }
func Warn(msg string, v ...Attr)  { getInstance().warn(msg, v...) }
