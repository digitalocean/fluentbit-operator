package copy

import (
	"fmt"
	"io"
	"os"
	"path"
)

// CopyFilesWithFilter copies all files directly inside the source directory
// into destination directory
// filterFn is used to dictate whether the path should be copied
func CopyFilesWithFilter(sourceDir, destDir string,
	filterFn func(info os.FileInfo) bool) error {

	sourceFi, err := os.Stat(sourceDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if !sourceFi.IsDir() {
		return fmt.Errorf("%s is not a directory", sourceDir)
	}

	destFi, err := os.Stat(destDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		err := os.Mkdir(destDir, 0755)
		if err != nil {
			return err
		}

		destFi, err = os.Stat(destDir)
		if err != nil {
			return err
		}
	}

	if !destFi.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}

	f, err := os.Open(sourceDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	fis, err := f.Readdir(-1)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, fi := range fis {
		if !filterFn(fi) {
			continue
		}

		fromPath := path.Join(sourceDir, fi.Name())
		fromFS, err := os.Open(fromPath)
		if err != nil {
			return err
		}
		defer fromFS.Close()

		toPath := path.Join(destDir, fi.Name())
		toFS, err := os.Create(toPath)
		if err != nil {
			return err
		}
		defer toFS.Close()

		_, err = io.Copy(toFS, fromFS)
		if err != nil {
			return err
		}
	}

	return nil
}
