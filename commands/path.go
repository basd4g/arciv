package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func rootDir() string {
	// find arciv's root directory (exist .arciv)
	// ex . current dir is /hoge/fuga/wara
	// search /hoge/fuga/wara/.arciv , and next /hoge/fuga/.arciv , and next /hoge/.arciv , and next /.arciv
	currentDir, _ := os.Getwd()
	for dir := currentDir; strings.LastIndex(dir, "/") != -1; dir = dir[:strings.LastIndex(dir, "/")] {
		if f, err := os.Stat(dir + "/.arciv"); !os.IsNotExist(err) && f.IsDir() {
			return dir
		}
	}
	Exit(errors.New(".arciv is not found"), 1)
	return ""
}

func findPathsOfSelfRepo(includesDir bool) (relativePaths []string, err error) {
	root := SelfRepo().Path
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !includesDir && info.IsDir() {
			return nil
		}
		if len(root) >= len(path) {
			// exclude root directory
			return nil
		}
		relativePath := path[len(root)+1:]
		if info.IsDir() && relativePath == ".arciv" {
			return nil
		}
		if strings.HasPrefix(relativePath, ".arciv/") {
			return nil
		}
		// add relative path from root directory
		relativePaths = append(relativePaths, relativePath)
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return relativePaths, nil
}

func lsWithoutDir(dir string) (filenames []string, err error) {
	os.MkdirAll(dir, 0777)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}, err
	}
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}
	return filenames, nil
}

func loadLines(path string) ([]string, error) {
	if !Exists(path) {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return []string{}, err
		}
		file.Close()
	}
	var lines []string
	f, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return []string{}, err
	}
	return lines, nil
}

func writeLines(path string, lines []string) error {
	os.MkdirAll(filepath.Dir(path), 0777)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fw := bufio.NewWriter(file)
	defer fw.Flush()
	for _, line := range lines {
		fmt.Fprintln(fw, line)
	}
	return nil
}

func writeLineTail(path string, line string) error {
	os.MkdirAll(filepath.Dir(path), 0777)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Fprintln(file, line)
	return nil
}

func copyFile(from string, to string) error {
	w, err := os.Create(to)
	if err != nil {
		return err
	}
	defer w.Close()

	r, err := os.Open(from)
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = io.Copy(w, r)
	return err
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
