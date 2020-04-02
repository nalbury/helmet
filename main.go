package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

var (
	basePath string = "/usr/local/bin"
)

func unTar(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err != nil && err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		switch {
		case header == nil:
			continue
		case header.Name == "darwin-amd64/helm":
			f, err := os.OpenFile(dest, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			f.Close()
			fmt.Println("Successfully copied helm excutable to", dest)
			return nil
		default:
			continue
		}
	}
}

func installVersion(v string) error {
	tarBall := "helm-v" + v + "-darwin-amd64.tar.gz"
	dir, err := ioutil.TempDir("", "helmenv")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	f, err := os.Create(dir + "/" + tarBall)
	if err != nil {
		return err
	}
	defer f.Close()

	url := "https://get.helm.sh/" + tarBall
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("Downloaded to", dir+"/"+tarBall)

	err = unTar(dir+"/"+tarBall, basePath+"/"+"helm-"+v)
	if err != nil {
		return err
	}
	fmt.Println("Successfully installed helm", v)
	return nil
}

func setVersion(v string) error {
	execPath := basePath + "/helm"
	installPath := execPath + "-" + v

	if _, err := os.Stat(execPath); err == nil {
		if err := os.Remove(execPath); err != nil {
			return err
		}
	}
	if _, err := os.Stat(installPath); err != nil {
		return fmt.Errorf("Helm version %s is not installed")
	}
	if err := os.Symlink(installPath, execPath); err != nil {
		return err
	}
	return nil
}

func listVersions() {
	basePathContents, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, file := range basePathContents {
		//https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
		re, err := regexp.Compile(`helm-(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if re.MatchString(file.Name()) {
			fmt.Println(file.Name())
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected at least 1 arg, got ", (len(os.Args) - 1))
		os.Exit(1)
	}
	switch arg := os.Args[1]; arg {
	case "help":
		usage := `Usage: helmet {ls,install} [VERSION]

Run with no subcmds, helmet symlinks the specified helm version to /usr/local/bin/helm.

SubCommands:
  - ls											list installed versions
	- install VERSION					install the specified version
`
		fmt.Println(usage)
	case "ls":
		listVersions()
	case "install":
		if len(os.Args) != 3 {
			fmt.Println("Expected 1 arg for install, got ", (len(os.Args) - 2))
			os.Exit(1)
		}
		v := os.Args[2]
		err := installVersion(v)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		fmt.Println("Setting helm version:", arg)
		if err := setVersion(arg); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
