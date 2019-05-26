package state

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/teacat/noire"
)

type fmtDecoder func(interface{}) (*noire.Color, error)

type PaletteConfig struct {
	Name   string
	Format string
	Colors []paletteColor `toml:"color"`
}

type paletteColor struct {
	RGB []int
	HSV []int
	HEX string
}

func (s *State) save(path string) error {
	log.Infof("saving state [%s]", path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(s.Config())
}

func (s *State) load(path string) error {
	var isDefault bool
	if path == "" {
		path = defaultPalettePath()
		isDefault = true
	}

	f, err := os.Open(path)
	if err != nil {
		if isDefault && os.IsNotExist(err) {

			return fmt.Errorf("failed to load palette: %s", err)
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		var config PaletteConfig
		if _, err := toml.Decode(b, &config); err != nil {
			return err
		}

		if config.Name == "" {
			config.Name = filepath.Base(path)
			config.Name = strings.ReplaceAll(config.Name, ".toml", "")
		}

		for n, pc := range config.Colors {
			if n > subStateCount {
				break
			}
			nc, err := pc.readColor()
			if err != nil {
				return fmt.Errorf("[color%d] %s", n, err)
			}
			s.sstates[n].SetNColor(nc)
			log.Debugf("loaded substate [%d] from %s", i, path)
		}

		log.Infof("loaded state [%s]", path)
		return nil
	}
}

func (pc *paletteColor) readColor() (*noire.Color, error) {
	switch {
	case len(pc.RGB) != 0:
		if err := pc.validRGB(); err != nil {
			return nil, fmt.Errorf("[color%d] %s", n, err)
		}
		return noire.NewRGB(float64(pc.RGB[0]), float64(pc.RGB[1]), float64(pc.RGB[2])), nil
	case len(pc.HSV) != 0:
		if err := pc.validHSV(); err != nil {
			return nil, fmt.Errorf("[color%d] %s", n, err)
		}
		return noire.NewHSV(float64(pc.HSV[0]), float64(pc.HSV[1]), float64(pc.HSV[2])), nil
	case len(pc.HEX) != 0:
		return nil, fmt.Errorf("[color%d] missing definition")
	default:
		return nil, fmt.Errorf("[color%d] missing definition")
	}
}

func (pc *paletteColor) validRGB() error {
	if len(pc.RGB) > 3 {
		return fmt.Errorf("malformed RGB (too many values)")
	}
	if len(pc.RGB) < 3 {
		return fmt.Errorf("malformed RGB (too few values)")
	}
	for _, x := range pc.RGB {
		if x > 255 || x < 0 {
			return fmt.Errorf("malformed RGB (values must be < 255)")
		}
	}
	return nil
}

func (pc *paletteColor) validHSV(a []int) error {
	if len(pc.HSV) > 3 {
		return fmt.Errorf("malformed HSV (too many values)")
	}
	if len(pc.HSV) < 3 {
		return fmt.Errorf("malformed HSV (too few values)")
	}
	if pc.HSV[0] < 0 || pc.HSV[0] > 359 {
		return fmt.Errorf("malformed HSV (hue out of 0-359 bounds)")
	}
	if pc.HSV[1] < 0 || pc.HSV[1] > 100 {
		return fmt.Errorf("malformed HSV (saturation out of 0-100 bounds)")
	}
	if pc.HSV[2] < 0 || pc.HSV[2] > 100 {
		return fmt.Errorf("malformed HSV (value out of 0-100 bounds)")
	}
	return nil
}
