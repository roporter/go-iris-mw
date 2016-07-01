package logger

import (
	"strconv"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/config"
	"github.com/kataras/iris/logger"
	"github.com/roporter/go-loreley"
)

type loggerMiddleware struct {
	*logger.Logger
	config Config
}

// Serve serves the middleware
func (l *loggerMiddleware) Serve(ctx *iris.Context) {
	//all except latency to string
	var date, status, ip, method, path string
	var latency time.Duration
	var startTime, endTime time.Time
	path = ctx.PathString()
	method = ctx.MethodString()

	startTime = time.Now()

	ctx.Next()
	//no time.Since in order to format it well after
	endTime = time.Now()
	date = endTime.Format("15:04:05.999999 - 02/01/2006")
	latency = endTime.Sub(startTime)

	if l.config.Status {
		status = strconv.Itoa(ctx.Response.StatusCode())
	}

	if l.config.IP {
		ip = ctx.RemoteAddr()
	}

	if !l.config.Method {
		method = ""
	}

	if !l.config.Path {
		path = ""
	}

	getText, _ := loreley.CompileAndExecuteToString(
		`{bold}{fg 15}{bg 40} GET {from "" 0}{reset}`,
		nil,
		nil,
	)

	//finally print the logs
	if(method == "GET") {
		l.printf("%s %s | %v | %4v | %s | %s \n", getText, date, status, latency, ip, path)
	} else {
		l.printf("%s %v %4v %s %s %s \n", date, status, latency, ip, method, path)
	}

}

func (l *loggerMiddleware) printf(format string, a ...interface{}) {
	if l.config.EnableColors {
		l.Logger.Otherf(format, a...)
	} else {
		l.Logger.Printf(format, a...)
	}
}

// New returns the logger middleware
// receives two parameters, both of them optionals
// first is the logger, which normally you set to the 'iris.Logger'
// if logger is nil then the middlewares makes one with the default configs.
// second is optional configs(logger.Config)
func New(theLogger *logger.Logger, cfg ...Config) iris.HandlerFunc {
	if theLogger == nil {
		theLogger = logger.New(config.DefaultLogger())
	}
	c := DefaultConfig().Merge(cfg)
	l := &loggerMiddleware{Logger: theLogger, config: c}

	return l.Serve
}
