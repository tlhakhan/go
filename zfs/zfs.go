package zfs

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

// Zpool type
type Zpool struct {
	Name string
}

type Database []Dataset

type Dataset struct {
	Name   string `json:"dataset"`
	Type   string `json:"type"`
	Origin string `json:"origin"`
	GUID   string `json:"guid"`
}

// Zpool factory function returns a Zpool type to operate against and pass around
func New(zpoolName string) (*Zpool, error) {

	// script to validate zpoolName
	script := []byte(`zpool list -Ho name $0`)

	if _, err := executeScript(script, zpoolName); err != nil {
		return nil, errors.Errorf("zpool %q does not exist", zpoolName)
	}

	return &Zpool{zpoolName}, nil
}

func (z *Zpool) Exists(datasetName string) bool {
	db, err := z.GetDatabase()
	if err != nil {
		// not good
		panic(err)
	}

	found := sort.Search(len(db), func(i int) bool { return db[i].Name == datasetName })
	if found < len(db) && db[found].Name == datasetName {
		return true
	}
	return false
}

func (z *Zpool) GetDatabase() (Database, error) {

	script := []byte(`zfs list -Hpro name,type,origin,guid -s createtxg -t all $0`)

	out, err := executeScript(script, z.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get dataset list for %q.", z.Name)
	}

	// initialize slice
	db := make(Database, 0)

	// parse output
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		db = append(db, Dataset{fields[0], fields[1], fields[2], fields[3]})
	}

	return db, nil
}

// Create the filesystem
func (z *Zpool) CreateFilesystem(filesystemName string) error {

	if !strings.HasPrefix(filesystemName, z.Name) {
		return errors.Errorf("unable to create filesystem %q, prefix does not match zpool name %q.", filesystemName, z.Name)
	}

	script := []byte(`zfs create -p $0`)
	_, err := executeScript(script, filesystemName)

	if err != nil {
		return errors.Wrapf(err, "unable to create filesystem %q.", filesystemName)
	}

	return nil
}
