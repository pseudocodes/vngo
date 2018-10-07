package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
	"vngo/core"
	"vngo/core/protocol"

	"github.com/apex/log"
	"github.com/spf13/viper"
)

type exitCode struct{ Code int }

func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(exitCode); ok {
			if exit.Code != 0 {
				fmt.Fprintln(os.Stderr, "TradeEngine failed at", time.Now().Format("January 2, 2006 at 3:04pm (MST)"))
			} else {
				fmt.Fprintln(os.Stderr, "Stopped TradeEngine at", time.Now().Format("January 2, 2006 at 3:04pm (MST)"))
			}

			os.Exit(exit.Code)
		}
		panic(e) // not an exitCode, bubble up
	}
}

func initConfig() {
	// The only command line arg is the config file
	configPath := flag.String("config-dir", ".", "Directory that contains the configuration file")
	flag.Parse()

	// Load the configuration from the file
	viper.SetConfigName("vngo")
	viper.AddConfigPath(*configPath)
	fmt.Fprintln(os.Stderr, "Reading configuration from", *configPath)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed reading configuration:", err.Error())
		panic(exitCode{1})
	}

	// setup viper to be able to read env variables with a configured prefix
	viper.SetDefault("general.env-var-prefix", "burrow")
	envPrefix := viper.GetString("general.env-var-prefix")
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Create the PID file to lock out other processes
	viper.SetDefault("general.pidfile", "burrow.pid")
	pidFile := viper.GetString("general.pidfile")
	if !core.CheckAndCreatePidFile(pidFile) {
		// Any error on checking or creating the PID file causes an immediate exit
		panic(exitCode{1})
	}
	defer core.RemovePidFile(pidFile)

	// Set up stderr/stdout to go to a separate log file, if enabled
	stdoutLogfile := viper.GetString("general.stdout-logfile")
	if stdoutLogfile != "" {
		core.OpenOutLog(stdoutLogfile)
	}
}

func start(app *protocol.ApplicationContext, exitChannel chan os.Signal) int {

	<-exitChannel
	log.Info("Shutdown triggered")

	app.Stop()
	return 0
}

func main() {
	defer handleExit()
	runtime.GOMAXPROCS(runtime.NumCPU())

	initConfig()
	// Register signal handlers for exiting
	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	// This triggers handleExit (after other defers), which will then call os.Exit properly
	panic(exitCode{start(nil, exitChannel)})
}
