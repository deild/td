package db

import (
	"fmt"
	"os"
	"path"

	"github.com/deild/td/helper"
)

// EnvDBPath environnement variable name for todo DB file
const EnvDBPath = "TODO_DB_PATH"

// DataStore structure
type DataStore struct {
	Path string
}

// NewDataStore search the path of database file
func NewDataStore() (*DataStore, error) {
	ds := new(DataStore)
	ds.Path = os.Getenv(EnvDBPath)

	if ds.Path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return ds, err
		}
		ds.Path = path.Join(cwd, ".todos")

	} else {
		dir, file := path.Split(ds.Path)
		if file == "" {
			ds.Path = path.Join(dir, ".todos")
		}

		fileInfo, err := os.Stat(ds.Path)

		if os.IsExist(err) {
			if fileInfo.IsDir() {
				ds.Path = path.Join(ds.Path, ".todos")
				err = os.Setenv(EnvDBPath, ds.Path)
				if err != nil {
					return ds, err
				}
			}
		}
	}

	return ds, nil
}

// Check if the database file exist
func (d *DataStore) Check() error {
	_, err := os.Stat(d.Path)
	if os.IsNotExist(err) {
		return fmt.Errorf("The database file \"%s\" doesn't exists", d.Path)
	}
	return nil
}

// Initialize the database file
func (d *DataStore) Initialize() error {
	var err error
	dir, _ := path.Split(d.Path)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s: One or more directories in this path doesn't exist", dir)
	}

	_, err = os.Stat(d.Path)
	if os.IsNotExist(err) {
		var w *os.File
		w, err = os.Create(d.Path)
		if err != nil {
			return err
		}
		defer helper.Check(w.Close)
		_, err = w.WriteString("[]")
		if err != nil {
			return err
		}
		return w.Sync()
	}
	return fmt.Errorf("%s: To-do file has been initialized before", d.Path)

}
