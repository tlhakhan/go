package zfs

import (
	"bufio"
	"bytes"
	"github.com/pkg/errors"
	"strings"
)

// Zpool type
type Zpool struct {
	Name string
}

// Filesystem type
type Filesystem struct {
	Name   string `json:"name"`
	Origin string `json:"origin"`
}

// Snapshot type
type Snapshot struct {
	Name string `json:"name"`
	GUID string `json:"guid"`
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

// Does dataset exist?
func (z *Zpool) Exists(datasetName string) bool {

	// guard against empty names
	if datasetName == "" {
		return false
	}

	// script to get the name of the dataset given
	script := []byte(`zfs get -Ho name name $0`)
	_, err := executeScript(script, datasetName)
	if err != nil {
		return false
	}
	return true
}

// GetSnapshots gets an array of snapshots for the given filesystem name.  Snapshots array is ordered by createtxg from oldest to newest.
func (z *Zpool) GetSnapshots(filesystemName string) ([]Snapshot, error) {

	script := []byte(`zfs list -Hpro name,guid -s createtxg -t snapshot -d1 $0`)

	out, err := executeScript(script, filesystemName)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get snapshots for filesystem %q.", filesystemName)
	}

	// initialize slice
	snapshots := make([]Snapshot, 0)

	// parse output
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		snapshots = append(snapshots, Snapshot{fields[0], fields[1]})
	}
	return snapshots, nil

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

// GetFilesystem ...
func (z *Zpool) GetFilesystem(filesystemName string) (Filesystem, error) {
	script := []byte(`
                        filesystem=$0
                        zfs list -Ho name,origin -s createtxg -t filesystem ${filesystem}
                    `)

	out, err := executeScript(script, filesystemName)
	if err != nil {
		return Filesystem{}, errors.Wrapf(err, "unable to get filesystem %q.", filesystemName)
	}

	fields := strings.Fields(string(out))
	return Filesystem{fields[0], fields[1]}, nil
}

// Get all filesystems on the zpool, the array will be ordered by createtxg from oldest to newest.
func (z *Zpool) GetFilesystems() ([]Filesystem, error) {
	script := []byte(`
                        zpool=$0
                        zfs list -Hro name,origin -s createtxg -t filesystem ${zpool}
                    `)

	out, err := executeScript(script, z.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get filesystems for zpool %q.", z.Name)
	}

	filesystems := make([]Filesystem, 0)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		filesystems = append(filesystems, Filesystem{fields[0], fields[1]})
	}
	return filesystems, nil
}

// Get the zpool status.
func (z *Zpool) GetStatus() ([]byte, error) {
	script := []byte(`
                        zpool=$0
                        zpool status $zpool
                    `)
	out, err := executeScript(script, z.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get zpool status on %q.", z.Name)
	}
	return out, nil
}
