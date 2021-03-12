package remote
import (
	"fmt"
	"path/filepath"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"github.com/CCDirectLink/CCUpdaterCLI/ccmod"
)

type packedModRemotePackage struct {
	metadata ccmodupdater.PackageMetadata
	location string
}

// NewPackedModRemotePackage checks the basic validity of a .ccmod file (i.e. that it has metadata) and then returns a ccmodupdater.RemotePackage for viewing or installing it.
func NewPackedModRemotePackage(location string) (ccmodupdater.RemotePackage, error) {
	metadata, err := ccmod.GetMetadata(location)
	if err != nil {
		return nil, err
	}
	if err = metadata.Verify(); err != nil {
		return nil, err
	}
	return packedModRemotePackage {
		metadata: metadata,
		location: location,
	}, nil
}

// Metadata implements RemotePackage.Metadata
func (mrp packedModRemotePackage) Metadata() ccmodupdater.PackageMetadata {
	return mrp.metadata
}

// Install implements RemotePackage.Install
func (mrp packedModRemotePackage) Install(game *ccmodupdater.GameInstance, log func(text string)) error {
	if mrp.metadata.Type() != ccmodupdater.PackageTypeMod {
		return fmt.Errorf("Unable to handle package type %s as a .ccmod", mrp.metadata.Type())
	}
	// Just copy it into the mods directory
	return copyFile(filepath.Join(game.Base(), "assets/mods", mrp.metadata.Name() + ".ccmod"), mrp.location)
}

