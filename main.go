package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/comhttp/explorer/app"
	daemon "github.com/leprosus/golang-daemon"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Get cmd line parameters
	port := flag.String("port", "", "Port")
	coin := flag.String("coin", "coin", "Coin")
	command := flag.String("command", "command", "Command")
	path := flag.String("path", "/var/db/jorm", "Path")
	loglevel := flag.String("loglevel", "info", "Logging level (debug, info, warn, error)")
	flag.Parse()

	//j.Log.SetLevel(parseLogLevel(*loglevel))
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Default level for this example is info, unless debug flag is present

	switch *loglevel {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	//log.Debug().Msg("This message appears only when log level set to Debug")
	//log.Info().Msg("This message appears when log level set to Debug or Info")

	log.Info().Msg("Service: explorer >>> Port: " + *port + "Coin: " + *coin + "JORM: ")

	err := daemon.Init(os.Args[0], map[string]interface{}{}, "./daemonized.pid")
	if err != nil {
		return
	}

	switch os.Args[1] {
	case "start":
		err = daemon.Start()
	case "stop":
		err = daemon.Stop()
	case "restart":
		err = daemon.Stop()
		err = daemon.Start()
	case "status":
		status := "stopped"
		if daemon.IsRun() {
			status = "started"
		}

		fmt.Printf("Application is %s\n", status)

		return
	case "":
	default:
		app.MainLoop(*path, *command, *coin)
		fmt.Println("JORM node is on: :" + *port)
	}
}
