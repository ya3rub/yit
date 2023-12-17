package internals

import "os"

var (
	// yit metadata dir
	YitMetadataDir = ".yit"
	// YitDefaultPermissions = os.FileMode(0744)
	YitDefaultDirPermissions = os.FileMode(0755)
	YitDefaultPermissions    = os.FileMode(0644)
	YitMetaDataDirContent    = []string{"objects", "refs"}
)
