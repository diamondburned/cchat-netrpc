package client

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func readPlugins() ([]os.DirEntry, error) {
	c, err := os.UserConfigDir()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user config dir")
	}

	d, err := os.ReadDir(filepath.Join(c, "cchat-netrpc", "plugins"))
	if err != nil {
		// IsNotExist isn't a fatal error; it just means there's no plugins.
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to read plugins")
	}

	plugins := d[:0]

	for _, entry := range d {
		if entry.IsDir() {
			continue
		}

		plugins = append(plugins, entry)
	}

	return plugins, nil
}

func init() {
	// plugins, err := readPlugins()
	// if err != nil {
	// 	// Transfer the error over somehow.
	// 	services.RegisterSource(func() []error {
	// 		return []error{err}
	// 	})

	// 	return
	// }

	// services.RegisterService()
}
