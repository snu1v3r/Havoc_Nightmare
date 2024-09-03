package profile

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Profile struct {
	config map[string]any
}

func NewProfile() *Profile {
	var profile *Profile

	profile = new(Profile)
	profile.config = make(map[string]any)

	return profile
}

func (p *Profile) Parse(path string) error {
	var (
		file []byte
		err  error
	)

	file, err = os.ReadFile(path)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(file, &p.config)
	if err != nil {
		return err
	}

	err = p.sanityCheck()
	if err != nil {
		return err
	}

	return nil
}

func (p *Profile) Config() map[string]any {
	return p.config
}
