// +build=linux
package kernel

import "os"

func randomWrite() error {
	file, err := os.CreateTemp(".", "random")
	if err != nil {
		return err
	}

	file.WriteAt()
}
