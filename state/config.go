package state

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/teacat/noire"
)

type fmtDecoder func(interface{}) (*noire.Color, error)

type PaletteConfig struct {
	Name   string
	Colors []paletteColor `toml:"color"`
}

type paletteColor struct {
	RGB []int     `toml:"rgb"`
	HSV []float64 `toml:"hsv"`
	HEX string    `toml:"hex"`
}

func (s *State) save() error {
	log.Infof("saving state [%s]", s.path)

	config := PaletteConfig{Name: s.Name()}
	for _, ss := range s.sstates {
		config.Colors = append(config.Colors, ss.PColor())
	}

	f, err := os.OpenFile(s.path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(config)
}

func (s *State) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			// palette does not exist yet, will be created on save
			return nil
		}
		return fmt.Errorf("failed to load palette: %s", err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var config PaletteConfig
	if _, err := toml.Decode(string(b), &config); err != nil {
		return err
	}

	s.name = config.Name

	for n, pc := range config.Colors {
		if n > subStateCount {
			break
		}
		nc, err := pc.readColor()
		if err != nil {
			return fmt.Errorf("[color%d] %s", n, err)
		}
		s.sstates[n].SetNColor(nc)
		log.Debugf("loaded substate [%d] from %s", n, s.path)
	}

	log.Infof("loaded state [%s]", s.path)
	return nil
}

func (pc *paletteColor) readColor() (*noire.Color, error) {
	switch {
	case len(pc.RGB) != 0:
		if err := pc.validRGB(); err != nil {
			return nil, err
		}
		return noire.NewRGB(float64(pc.RGB[0]), float64(pc.RGB[1]), float64(pc.RGB[2])), nil
	case len(pc.HSV) != 0:
		if err := pc.validHSV(); err != nil {
			return nil, err
		}
		return noire.NewHSV(pc.HSV[0], pc.HSV[1], pc.HSV[2]), nil
	case len(pc.HEX) != 0:
		return nil, fmt.Errorf("missing definition")
	default:
		return nil, fmt.Errorf("missing definition")
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

func (pc *paletteColor) validHSV() error {
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
