package commands

import (
	"os"
	"time"

	"github.com/mazzus/elevator-simulator/backend/src/elevator"
	httpHandlers "github.com/mazzus/elevator-simulator/backend/src/handlers/http"
	standardHandler "github.com/mazzus/elevator-simulator/backend/src/handlers/standard"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
)

const mainUsage = "Usage"
const module = "Root"

var (
	clientPort   int
	webPort      int
	floors       int
	logLevel     string
	shh          bool
	speed        float64
	margin       float64
	updatePeriod int32
)

var log = logging.MustGetLogger("root")

func init() {
	f := Root.PersistentFlags()
	f.IntVarP(&clientPort, "client-port", "c", 15657, "Specifies the port used for the client.")
	f.IntVarP(&webPort, "web-port", "w", 3001, "Specifies the port used for the web frontend.")
	f.IntVarP(&floors, "floors", "f", 4, "Specifies the number of floors.")
	f.StringVarP(&logLevel, "log-level", "l", "INFO", "The minimum log level which will show.")
	f.BoolVar(&shh, "shh", false, "This is used to silence the friendly welcome message")
	f.Float64VarP(&speed, "speed", "s", 0.4, "This sets the velocity of the elevator")
	f.Float64VarP(&margin, "margin", "m", 0.05, "This determines the size of the margin in which the floor sensors will detect the elevator")
	f.Int32VarP(&updatePeriod, "update-period", "p", 5, "The update period of the elevator, in milliseconds")
}

// Root is the root command
var Root = &cobra.Command{
	Use:   "elevator-simulator",
	Short: "Graphical simulator fro the Sanntidslab",
	Long:  mainUsage,
	RunE: func(cmd *cobra.Command, args []string) error {

		initializeLogging(logLevel)

		if !shh {
			log.Info(
				`
Hi!
This is an early version of the simulator, feedback is greatly appreciated!
The simulator will probably have some breakdowns, to prevent you from having a mental breakdown of your own, submit an issue!
Link for filing issues: https://github.com/mazzus/graphical-elevator-simulator/issues

Happy coding!

Tired of this message? The argument --shh will silence it ;)
			`)
		}

		var elev elevator.SafeElevator
		elev.Elevator = elevator.NewElevator(0, floors, speed, margin)

		go httpHandlers.HTTPServer(webPort, &elev)
		go standardHandler.Server(clientPort, &elev)

		updateTicker := time.NewTicker(time.Millisecond * time.Duration(updatePeriod))
		lastUpdateTime := time.Now()
		for {
			select {
			case <-updateTicker.C:
				now := time.Now()
				dTime := now.Sub(lastUpdateTime) // Using time.Now() because the time in the channel might be long ago because of the select
				elev.Update(dTime.Seconds())
				lastUpdateTime = now
			}
		}

	},
}

func initializeLogging(logLevel string) {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	format := logging.MustStringFormatter("%{color}%{level:.4s} %{time:15:04:05.000} [%{module}]--->%{color:reset} %{message} \n\n")

	backendFormatter := logging.NewBackendFormatter(backend, format)
	leveled := logging.AddModuleLevel(backendFormatter)
	level, err := logging.LogLevel(logLevel)
	if err != nil {
		panic(err)
	}
	leveled.SetLevel(level, "")

	logging.SetBackend(leveled)
}
