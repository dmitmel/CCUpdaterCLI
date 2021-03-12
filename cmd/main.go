package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/api"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
	"github.com/CCDirectLink/CCUpdaterCLI/cmd/commands"
)

func assertContext(options internal.ContextOptions) *internal.Context {
	context, err := internal.NewContext(options)
	if err != nil {
		fmt.Printf("UNABLE TO FIND GAME in %s\n", err.Error())
		os.Exit(1)
	}
	return context
}
func assertOnlineContext(options internal.ContextOptions) *internal.OnlineContext {
	onlineContext, err := assertContext(options).Upgrade()
	if err != nil {
		fmt.Printf("UNABLE TO GO ONLINE in %s\n", err.Error())
		os.Exit(1)
	}
	return onlineContext
}

func main() {
	flagGame := flag.String("game", "", "if set it overrides the path of the game")
	flagForce := flag.Bool("force", false, "for commands that perform actions: ignores automatic dependency handling")

	flagPort := flag.Int("port", 9392, "the port which the api server listens on")
	flagHost := flag.String("host", "localhost", "the host which the api server listens on")

	flagVerbose := flag.Bool("v", false, "makes certain commands report more verbose output")
	flagAll := flag.Bool("all", false, "for list: indicates all kinds of packages should be shown")

	flag.Parse()

	options := internal.ContextOptions {
		Game: *flagGame,
		Force: *flagForce,
		Verbose: *flagVerbose,
	}

	if len(os.Args) == 1 {
		printHelp()
		return
	}

	op := flag.Arg(0)
	args := flag.Args()[1:]

	switch op {
	case "install",
		"i":
		printStatsAndError(commands.Install(assertOnlineContext(options), args))
	case "remove",
		"delete",
		"uninstall":
		printStatsAndError(commands.Uninstall(assertContext(options), args))
	case "update":
		printStatsAndError(commands.Update(assertOnlineContext(options), args))
	case "list":
		commands.List(assertOnlineContext(options), *flagAll)
	case "outdated":
		commands.Outdated(assertOnlineContext(options))
	case "api":
		api.StartAt(*flagHost, *flagPort)
	case "version":
		printVersion()
	case "help":
		printHelp()
	default:
		fmt.Printf("%s\n is not a command", op)
		printHelp()
		os.Exit(1)
	}
}

func printStatsAndError(stats *internal.Stats, err error) {
	if stats != nil && stats.Warnings != nil {
		for _, warning := range stats.Warnings {
			fmt.Printf("Warning in %s\n", warning)
		}
	}

	if err != nil {
		fmt.Printf("ERROR in %s\n", err.Error())
	}

	if stats != nil {
		fmt.Printf("Installed %d, updated %d, removed %d\n", stats.Installed, stats.Updated, stats.Removed)
	}
}
