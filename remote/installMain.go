package remote
import (
	"fmt"
	"os"
	"path/filepath"
)

// Attempts to install a modZip installation method (assuming type and platform checks have passed) to the given directory.
// source is a pointer because it's optional.
func tryExecuteModZipInstallationMethod(url string, source *string, newDir string) error {
	err := os.MkdirAll("installing", os.ModePerm)
	if err != nil {
		return fmt.Errorf("Unable to make temp directory: %s", err.Error())
	}
	defer os.RemoveAll("installing")

	file, err := download(url)
	if err != nil {
		return fmt.Errorf("Unable to download: %s", err.Error())
	}
	defer os.Remove(file.Name())

	dir, err := extract(file)
	if err != nil {
		return fmt.Errorf("Unable to extract: %s", err.Error())
	}
	defer os.RemoveAll(dir)
	
	dirSrc := dir
	if source != nil {
		dirSrc = filepath.Join(dirSrc, *source)
	}
	
	return copyDir(newDir, dirSrc)
}

// Attempts to install a ccmod to the given directory.
// Notice that this could be done by extraction, but that isn't necessary anymore.
func tryExecutePackedModInstallationMethod(url string, target string) error {
	err := os.MkdirAll("installing", os.ModePerm)
	if err != nil {
		return fmt.Errorf("Unable to make temp directory: %s", err.Error())
	}
	defer os.RemoveAll("installing")

	file, err := download(url)
	if err != nil {
		return fmt.Errorf("Unable to download: %s", err.Error())
	}
	defer os.Remove(file.Name())

	return copyFile(target, file.Name())
}

// Attempts to install the given installation method to the given directory.
func tryExecuteInstallationMethod(method ccModDBInstallationMethod, target string) error {
	if method.Platform != nil {
		platform := *method.Platform
		if platform != whatPlatformAreWe() {
			return fmt.Errorf("Installation method requires platform %s", platform)
		}
	}
	
	if method.Type == "modZip" {
		return tryExecuteModZipInstallationMethod(method.URL, method.Source, target)
	}
	if method.Type == "ccmod" {
		return tryExecutePackedModInstallationMethod(method.URL, target + ".ccmod")
	}
	return fmt.Errorf("Unable to interpret installation method of type %s", method.Type)
}
