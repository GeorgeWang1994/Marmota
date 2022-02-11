package file

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func FileExist(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if stat.IsDir() {
		return false, errors.Errorf("%s is dir", path)
	}
	return true, nil
}

func DirExist(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if !stat.IsDir() {
		return false, errors.Errorf("%s is file", path)
	}
	return true, nil
}

func CheckExistInDir(filename, dirPath string) (bool, error) {
	if exist, err := DirExist(dirPath); exist && err == nil {
		if files, err := ioutil.ReadDir(dirPath); err == nil {
			for _, f := range files {
				if f.Name() == filename {
					return true, nil
				}
			}
		}
	}
	return false, errors.Errorf("not exist %s", filename)
}

func WriteYamlToFile(in interface{}, outPath string) error {
	data, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outPath, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func ReadYamlFromFile(path string, out interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}
	return nil
}

func LookupFile(basePath string, filename string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(basePath, filename))
	if len(matches) == 0 {
		return lookupInNearestDir(basePath, filename)
	}
	return matches[0], err
}

func lookupInNearestDir(basePath string, filename string) (string, error) {
	if basePath == "/" {
		return "", errors.New("file not found")
	}
	nearest := filepath.Clean(filepath.Join(filepath.Dir(basePath), ".."))
	return LookupFile(nearest, filename)
}

func CopyFile(src, dest string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return errors.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func CopyDir(src string, dest string) error {
	if srcInfo, err := os.Stat(src); err != nil {
		return err
	} else {
		if !srcInfo.IsDir() {
			return errors.New("srcPath不是一个正确的目录！")
		}
	}
	if destInfo, err := os.Stat(dest); err != nil {
		return err
	} else {
		if !destInfo.IsDir() {
			return errors.New("destInfo不是一个正确的目录！")
		}
	}

	err := filepath.Walk(src, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if path == src {
			return nil
		}
		p := strings.Replace(path, "\\", "/", -1)
		destNewPath := strings.Replace(p, src, dest, -1)
		if !f.IsDir() {
			if err = CopyFile(path, destNewPath); err != nil {
				return err
			}
		} else {
			_ = os.MkdirAll(destNewPath, os.ModePerm)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func FileContent(path string) (string, error) {
	if exists, err := FileExist(path); err != nil && exists {
		return "", err
	}
	if cc, err := ioutil.ReadFile(path); err != nil {
		return "", err
	} else {
		return string(cc), nil
	}
}
