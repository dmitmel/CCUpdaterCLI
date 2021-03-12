package commands

import (
	"fmt"

	"github.com/CCDirectLink/CCUpdaterCLI/cmd/internal"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/local"
	"github.com/CCDirectLink/CCUpdaterCLI/remote"
)

//InstallFrom installs mods from ccmod files
func InstallFrom(context *internal.Context, args []string, online bool) (*internal.Stats, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("cmd: No mods installed since no mods were specified")
	}

	stats := &internal.Stats{}

	remotePackages := make(map[string]ccmodupdater.RemotePackage)

	// Add online to remote packages
	if online {
		oc, err := context.Upgrade()
		if err != nil {
			stats.AddWarning(fmt.Sprintf("cmd: Could not go online: %s", err))
		} else {
			for k, v := range oc.RemotePackages() {
				remotePackages[k] = v
			}
		}
	}

	// Setup transaction and ccmod packages
	tx := make(ccmodupdater.PackageTX)
	installedMods := context.Game().Packages()

	for _, path := range args {
		ccm, err := remote.NewPackedModRemotePackage(path)
		if err != nil {
			return nil, fmt.Errorf("cmd: In packed mod %s: %s", path, err)
		}
		name := ccm.Metadata().Name()
		if _, modExists := installedMods[name]; modExists {
			stats.AddWarning(fmt.Sprintf("cmd: Could not install '%s' because it was already installed", name))
			continue
		}
		// Ok, add it (overriding any online copy of the mod)
		remotePackages[name] = ccm
		tx[name] = ccmodupdater.PackageTXOperationInstall
	}

	err := context.ExecuteWithRemotePackages(tx, stats, remotePackages)
	for _, warning := range local.CheckLocal(context.Game(), remotePackages) {
		stats.AddWarning(warning)
	}
	return stats, err
}
