package public
import "github.com/Masterminds/semver"

// PackageType represents a specific kind of package.
type PackageType int

// PackageTypeFirst is the first PackageType.
const PackageTypeFirst PackageType = 0

// PackageTypeBase represents a hidden base package that should not be visible from UI (but may be visible in CLI).
const PackageTypeBase PackageType = 0
// PackageTypeMod represents a mod.
const PackageTypeMod PackageType = 1
// PackageTypeTool represents a tool.
const PackageTypeTool PackageType = 2

// PackageTypeAfterLast is the last PackageType.
const PackageTypeAfterLast PackageType = 3

func (pt PackageType) String() string {
	switch pt {
		case PackageTypeBase:
			return "Base"
		case PackageTypeMod:
			return "Mod"
		case PackageTypeTool:
			return "Tool"
	}
	return string(pt)
}

// StringPlural is similar to String, but the resulting string is a plural
func (pt PackageType) StringPlural() string {
	if pt == PackageTypeBase {
		return "Base Packages"
	}
	return pt.String() + "s"
}

// PackageMetadata contains the metadata of a package.
type PackageMetadata struct {
	// The name of this package; identifies it hopefully uniquely. Must be the same as the key in various maps.
	Name string
	// The type of this package, 
	Type PackageType
	// The description of this package, human-readable.
	Description string
	// The version of this package, in SemVer format.
	Version *semver.Version
}

// Package represents a package, local or remote.
type Package interface {
	// Returns the metadata for this package.
	Metadata() PackageMetadata
}

// RemotePackage represents a package on a remote server.
type RemotePackage interface {
	Package
	// Installs the package. This does not check dependencies. This may affect other packages.
	Install(target *GameInstance) error
}

// LocalPackage represents a package installed in a GameInstance. It is invalidated when the GameInstance is modified.
type LocalPackage interface {
	Package
	// Attempts to remove the package. This does not check dependencies. This may affect other packages.
	Remove() error
	// Acquires a dependency map (mapping package IDs to version constraints). DO BE ALERTED! This may be made part of PackageMetadata when the time comes!
	Dependencies() map[string]string
}