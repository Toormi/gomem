package snapshot

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var SnapshotUrl string

func init() {
	var err error
	if SnapshotUrl == "" {
		SnapshotUrl, err = os.Getwd()
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

// Restore restore data from snapshot
func Restore() (offset int64, data []byte, err error) {
	latestName, offset, err := lastSnapshotName(SnapshotUrl)
	if err != nil {
		return
	}

	if latestName != "" {
		data, err = unzipFile(fmt.Sprintf("%s/%s", SnapshotUrl, latestName))
		if err != nil {
			return
		}
	}

	// if offset not 0 should +1 ,set to newest offset
	if offset != 0 {
		offset++
	}

	return
}

// Restore restore date from snapshot with request and report file name
func RestoreDoubleOffset(subUrl string) (requestOffset, reportOffset int64, data []byte, err error) {
	latestName, requestOffset, reportOffset, err := lastDoubleOffsetSnapshotName(fmt.Sprintf("%s%s", SnapshotUrl, subUrl))
	if err != nil {
		return
	}

	if latestName != "" {
		data, err = unzipFile(fmt.Sprintf("%s%s/%s", SnapshotUrl, subUrl, latestName))
		if err != nil {
			return
		}
	}

	return
}

// BackUpDoubleOffset back up data to snapshot with request and report offset
func BackUpDoubleOffset(subUrl string, requestOffset, reportOffset int64, data []byte) error {
	fileName := fmt.Sprintf("%s%s/%d-%d", SnapshotUrl, subUrl, requestOffset, reportOffset)

	if err := zipFile(fileName, data); err != nil {
		return err
	}

	return nil
}

// BackUp back up data to snapshot
func BackUp(offset int64, data []byte) error {
	fileName := fmt.Sprintf("%s/%d", SnapshotUrl, offset)

	if err := zipFile(fileName, data); err != nil {
		return err
	}

	return nil
}

func lastSnapshotName(wd string) (snapshotName string, offset int64, err error) {
	if !FileExist(wd) {
		return
	}

	dir, err := ioutil.ReadDir(wd)
	if err != nil {
		return
	}

	var offsets []int64
	for _, file := range dir {
		fileName := file.Name()
		if !file.IsDir() && filepath.Ext(fileName) == ".zip" {
			sli := strings.Split(fileName, ".")
			if len(sli) != 2 {
				err = errors.New("file name not correct")
				return
			}

			o, innerr := strconv.ParseInt(sli[0], 10, 64)
			if innerr != nil {
				err = innerr
				return
			}
			offsets = append(offsets, o)
		}
	}

	if len(offsets) > 0 {
		sort.Slice(offsets, func(i, j int) bool {
			return offsets[i] < offsets[j]
		})

		offset = offsets[len(offsets)-1]
		snapshotName = fmt.Sprintf("%d.zip", offset)
	}

	return
}

func lastDoubleOffsetSnapshotName(wd string) (snapshotName string, requestOffset, reportOffset int64, err error) {
	if !FileExist(wd) {
		return
	}

	dir, err := ioutil.ReadDir(wd)
	if err != nil {
		return
	}

	type OffsetPair struct {
		request int64
		report  int64
	}
	var offsetPairs []*OffsetPair
	for _, file := range dir {
		fileName := file.Name()
		if !file.IsDir() && filepath.Ext(fileName) == ".zip" {
			sli := strings.Split(fileName, ".")
			if len(sli) != 2 {
				err = errors.New("file name not correct")
				return
			}

			offsets := strings.Split(sli[0], "-")
			if len(offsets) != 2 {
				err = errors.New("file name not double offset correct")
				return
			}

			reqOffset, innerr := strconv.ParseInt(offsets[0], 10, 64)
			if innerr != nil {
				err = innerr
				return
			}
			rptOffset, innerr := strconv.ParseInt(offsets[1], 10, 64)
			if innerr != nil {
				err = innerr
				return
			}
			offsetPairs = append(offsetPairs, &OffsetPair{request: reqOffset, report: rptOffset})
		}
	}

	if len(offsetPairs) > 0 {
		sort.Slice(offsetPairs, func(i, j int) bool {
			return offsetPairs[i].report < offsetPairs[j].report
		})

		offsetPair := offsetPairs[len(offsetPairs)-1]
		requestOffset = offsetPair.request
		reportOffset = offsetPair.report
		snapshotName = fmt.Sprintf("%d-%d.zip", requestOffset, reportOffset)
	}

	return
}

func unzipFile(fileName string) ([]byte, error) {
	r, err := zip.OpenReader(fileName)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if len(r.File) != 1 {
		return nil, errors.New("zip file not correct")
	}

	rc, err := r.File[0].Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return ioutil.ReadAll(rc)
}

func zipFile(fileName string, data []byte) error {
	txtName := fmt.Sprintf("%s.text", fileName)
	tmpName := fmt.Sprintf("%s/tmp.zip", SnapshotUrl)
	snapshotName := fmt.Sprintf("%s.zip", fileName)
	err := ioutil.WriteFile(txtName, data, 0666)
	if err != nil {
		return err
	}

	err = Zip(txtName, tmpName)
	if err != nil {
		return err
	}

	err = os.Rename(tmpName, snapshotName)
	if err != nil {
		return err
	}

	err = os.Remove(txtName)
	if err != nil {
		return err
	}

	return nil
}

// srcFile could be a single file or a directory
func Zip(srcFile string, destZip string) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	err = filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		var file *os.File
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
		// header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err = os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

func FileExist(wd string) bool {
	_, err := os.Stat(wd)
	if err == nil {
		return true
	}

	return os.IsExist(err)
}
