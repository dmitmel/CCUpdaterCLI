package ccmod

// .ccmod file metadata extraction
// This is a separate package because remote/ccmod.go will need this to peek at the files too.
// It also keeps the code nicely separated.
import (
	"encoding/json"
	"archive/zip"
	"github.com/CCDirectLink/CCUpdaterCLI"
	"fmt"
)

func GetMetadata(filename string) (ccmodupdater.PackageMetadata, error) {
	reader, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	for _, file := range reader.File {
		if file.Name == "package.json" {
			metadataFileHandle, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer metadataFileHandle.Close()
			metadata := make(ccmodupdater.PackageMetadata)
			err = json.NewDecoder(metadataFileHandle).Decode(&metadata)
			if err != nil {
				return nil, err
			}
			return metadata, nil
		}
	}
	return nil, fmt.Errorf("Unable to find package.json in packed (.ccmod) mod: %s", filename)
}

