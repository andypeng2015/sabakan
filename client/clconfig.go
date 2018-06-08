package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/coreos/container-linux-config-transpiler/config/types"
	yaml "gopkg.in/yaml.v2"
)

const baseFileDir = "./files"
const baseSystemdDir = "./systemd"
const baseNetworkdDir = "./networkd"

type systemd struct {
	Enabled bool   `yaml:"enabled"`
	Source  string `yaml:"source"`
}

type ignitionSource struct {
	Passwd   string    `yaml:"passwd"`
	Files    []string  `yaml:"files"`
	Systemd  []systemd `yaml:"systemd"`
	Networkd []string  `yaml:"networkd"`
	Include  string    `yaml:"include"`
}

func constructCLConfigTemplate(fname string) (io.Reader, error) {
	source, err := loadSource(fname)
	if err != nil {
		return nil, ErrorStatus(err)
	}

	var clConf types.Config
	err = constructCLConf(source, &clConf)
	if err != nil {
		return nil, ErrorStatus(err)
	}

	b, err := yaml.Marshal(clConf)
	if err != nil {
		return nil, ErrorStatus(err)
	}

	return bytes.NewReader(b), nil
}

func loadSource(fname string) (*ignitionSource, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var source ignitionSource
	err = yaml.Unmarshal(data, &source)
	if err != nil {
		return nil, err
	}
	return &source, nil
}

func constructCLConf(source *ignitionSource, clConf *types.Config) error {
	if source.Include != "" {
		include, err := loadSource(source.Include)
		if err != nil {
			return err
		}
		err = constructCLConf(include, clConf)
		if err != nil {
			return err
		}
	}
	if source.Passwd != "" {
		err := constructPasswd(source.Passwd, clConf)
		if err != nil {
			return err
		}
	}

	for _, file := range source.Files {
		err := constructFile(file, clConf)
		if err != nil {
			return err
		}
	}

	for _, s := range source.Systemd {
		err := constructSystemd(s, clConf)
		if err != nil {
			return err
		}
	}

	for _, n := range source.Networkd {
		err := constructNetworkd(n, clConf)
		if err != nil {
			return err
		}
	}
	return nil
}

func constructPasswd(passwd string, clConf *types.Config) error {
	pf, err := os.Open(passwd)
	if err != nil {
		return err
	}
	defer pf.Close()
	passData, err := ioutil.ReadAll(pf)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(passData, &clConf.Passwd)
}

func constructFile(inputFile string, clConf *types.Config) error {
	p := filepath.Join(baseFileDir, inputFile)
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	fi, err := os.Stat(p)
	if err != nil {
		return err
	}
	mode := int(fi.Mode())

	clConf.Storage.Files = append(clConf.Storage.Files, types.File{
		Path:       inputFile,
		Filesystem: "root",
		Mode:       &mode,
		Contents: types.FileContents{
			Inline: string(data),
		},
	})

	return nil
}

func constructSystemd(s systemd, clConf *types.Config) error {

	f, err := os.Open(filepath.Join(baseSystemdDir, s.Source))
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	clConf.Systemd.Units = append(clConf.Systemd.Units, types.SystemdUnit{
		Name:     s.Source,
		Enabled:  &s.Enabled,
		Contents: string(data),
	})

	return nil
}

func constructNetworkd(n string, clConf *types.Config) error {

	f, err := os.Open(filepath.Join(baseNetworkdDir, n))
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	clConf.Networkd.Units = append(clConf.Networkd.Units, types.NetworkdUnit{
		Name:     n,
		Contents: string(data),
	})

	return nil
}
