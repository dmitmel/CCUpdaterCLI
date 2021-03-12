package internal

import (
	"fmt"
	"os"
	
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/local"
	"github.com/CCDirectLink/CCUpdaterCLI/remote"
)

//ContextOptions contains basic options for contexts that need to be supplied by main.go
type ContextOptions struct {
	// Game: Path to game (optional, blank string = cwd)
	Game string
	// Force: Ignore dependency checks
	Force bool
	// Verbose: Output more output
	Verbose bool
}

// Convenience converter for use by API
func GamePtrOptConv(game *string) ContextOptions {
	gameStr := ""
	if game != nil {
		gameStr = *game
	}
	return ContextOptions {
		Game: gameStr,
		Force: false,
		Verbose: false,
	}
}

//Context contains the context details for this command.
type Context struct {
	game *ccmodupdater.GameInstance
	options ContextOptions
	upgraded *OnlineContext
}

//NewContext creates a new local context.
func NewContext(opts ContextOptions) (*Context, error) {
	if opts.Game == "" {
		gameDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("Unable to get working directory (as game directory): %s", err)
		}
		opts.Game = gameDir
	}
	game := ccmodupdater.NewGameInstance(opts.Game)
	plugins, err := local.AllLocalPackagePlugins(game)
	if err != nil {
		return nil, fmt.Errorf("Unable to prepare for checking local packages: %s", err)
	}
	game.LocalPlugins = plugins
	return &Context{
		game: game,
		options: opts,
		upgraded: nil,
	}, nil
}

//NewOnlineContext creates a new online context.
func NewOnlineContext(opts ContextOptions) (*OnlineContext, error) {
	ctx, err := NewContext(opts)
	if err != nil {
		return nil, err
	}
	rwc, err := ctx.Upgrade()
	if err != nil {
		return nil, err
	}
	return rwc, nil
}

func (ctx *Context) Game() *ccmodupdater.GameInstance {
	return ctx.game
}

func (ctx *Context) Options() ContextOptions {
	return ctx.options
}

//Upgrade upgrades the Context to an OnlineContext.
func (ctx *Context) Upgrade() (*OnlineContext, error) {
	if ctx.upgraded != nil {
		return ctx.upgraded, nil
	}
	packages, err := remote.GetRemotePackages()
	if err != nil {
		return nil, err
	}
	rwc := &OnlineContext{
		Context: *ctx,
		remote: packages,
	}
	ctx.upgraded = rwc
	rwc.upgraded = rwc
	return rwc, nil
}

// Execute executes a package transaction. If this has been upgraded to an OnlineContext, it can install packages, otherwise it can't. For .ccmod use see ExecuteWithRemotePackages.
func (ctx *Context) Execute(tx ccmodupdater.PackageTX, stats *Stats) error {
	upgraded := ctx.upgraded
	var rp map[string]ccmodupdater.RemotePackage
	if upgraded != nil {
		rp = upgraded.remote
	} else {
		rp = make(map[string]ccmodupdater.RemotePackage)
	}
	return ctx.ExecuteWithRemotePackages(tx, stats, rp)
}

// ExecuteWithRemotePackages executes a package transaction with a specific set of remote packages.
func (ctx *Context) ExecuteWithRemotePackages(tx ccmodupdater.PackageTX, stats *Stats, remotePackages map[string]ccmodupdater.RemotePackage) error {
	tc := ccmodupdater.PackageTXContext{
		LocalPackages: ctx.game.Packages(),
		RemotePackages: remotePackages,
	}
	if !ctx.options.Force {
		solutions, err := tc.Solve(tx)
		if err != nil {
			return err
		}
		if len(solutions) > 1 {
			return fmt.Errorf("Dependency issue; can solve this in multiple ways. (This shouldn't happen in the current system.) %v", solutions)
		}
		if len(solutions) == 0 {
			return fmt.Errorf("Internal error caused no solutions to be returned yet no error was returned.")
		}
		tx = solutions[0]
	}
	return tc.Perform(ctx.game, tx, func (pkg string, pre bool, remove bool, install bool) {
		if install && remove {
			if pre {
				fmt.Fprintf(os.Stderr, "updating %s\n", pkg)
			} else {
				stats.Updated++
			}
		} else if install {
			if pre {
				fmt.Fprintf(os.Stderr, "installing %s\n", pkg)
			} else {
				stats.Installed++
			}
		} else if remove {
			if pre {
				fmt.Fprintf(os.Stderr, "removing %s\n", pkg)
			} else {
				stats.Removed++
			}
		}
	}, func (text string) {
		fmt.Fprintln(os.Stderr, text)
	})
}

//OnlineContext contains the details for an online context.
type OnlineContext struct {
	Context
	remote map[string]ccmodupdater.RemotePackage
}

//RemotePackages returns all the remote packages.
func (rwc *OnlineContext) RemotePackages() map[string]ccmodupdater.RemotePackage {
	target := map[string]ccmodupdater.RemotePackage{}
	for k, v := range rwc.remote {
		target[k] = v
	}
	return target
}

//Stats contains the statistics about the installed mods
type Stats struct {
	Installed int `json:"installed"`
	Updated   int `json:"updated"`
	Removed   int `json:"removed"`

	Warnings []string `json:"warnings,omitempty"`
}

//AddWarning to the statistics
func (stats *Stats) AddWarning(warning string) {
	stats.Warnings = append(stats.Warnings, warning)
}
