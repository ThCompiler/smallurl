package app

import (
	"io"
	"log"
	"net/http"
	"os"
	"smallurl/internal/app/config"
	v1 "smallurl/internal/app/delivery/http/v1"
	plogger "smallurl/internal/pkg/logger"
	sh "smallurl/internal/shortcut/delivery/http/v1/handlers"
	"smallurl/pkg/logger"

	sf "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"

	_ "smallurl/docs"
)

func prepareLogger(cfg config.LoggerInfo) (*logger.Logger, *os.File) {
	var logOut io.Writer

	var logFile *os.File

	var err error

	if cfg.Directory != "" {
		logFile, err = plogger.OpenLogDir(cfg.Directory)
		if err != nil {
			log.Fatalf("[App] Init - create logger error: %s", err) // nolint: revive // логгер инициализируется,
			// ошибку открытия лог файла больше нечем логировать
		}

		logOut = logFile
	} else {
		logOut = os.Stderr
		logFile = nil
	}

	l := logger.New(
		logger.Params{
			AppName:                  cfg.AppName,
			LogDir:                   cfg.Directory,
			Level:                    cfg.Level,
			UseStdAndFile:            cfg.UseStdAndFile,
			AddLowPriorityLevelToCmd: cfg.AllowShowLowLevel,
		},
		logOut,
	)

	return l, logFile
}

func prepareRoutes(shortcutHandlers *sh.ShortcutHandlers) v1.Routes {
	return v1.Routes{
		// "Swagger"
		v1.Route{
			Method:      http.MethodGet,
			Pattern:     "/swagger/*any",
			HandlerFunc: gs.WrapHandler(sf.Handler),
		},

		// "GenShort"
		v1.Route{
			Method:      http.MethodPost,
			Pattern:     "/shorten",
			HandlerFunc: shortcutHandlers.GenShort,
		},

		// "GetLong"
		v1.Route{
			Method:      http.MethodGet,
			Pattern:     "/:" + sh.ShortURLParam,
			HandlerFunc: shortcutHandlers.GetLong,
		},
	}
}
