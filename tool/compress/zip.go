package compress

import (
	"archive/zip"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
)

func Unzip(f string, dir string) error {
	return unzip(f, dir)
}
func unzip(f string, dir string) error {
	read, err := zip.OpenReader(f)
	if err != nil {
		errors.WithStack(err)
	}
	defer read.Close()
	return unzipFiles(read.File, dir)
}
func unzipFiles(files []*zip.File, dir string) error {
	fmt.Println(len(files))
	for _, f := range files {
		fmt.Printf("%s %+v\n", f.Name, f)
		if f.FileInfo().IsDir() {
			err := os.Mkdir(path.Join(dir, f.Name), f.Mode())
			if err != nil {
				errors.WithStack(err)
			}
		} else {
			r, err := f.Open()
			if err != nil {
				errors.WithStack(err)
			}
			data, err := ioutil.ReadAll(r)
			r.Close()
			if err != nil {
				errors.WithStack(err)
			}
			err = ioutil.WriteFile(path.Join(dir, f.Name), data, f.Mode())
			if err != nil {
				errors.WithStack(err)
			}
		}
	}
	return nil
}
