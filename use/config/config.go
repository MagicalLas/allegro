package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	enrolled map[string]BinaryDirectory
}

type BinaryDirectory struct {
	name string
	path string
}

func LoadFromFile(path string, newConfig bool) (*Config, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("fail to parse config file path. wrong path: %w", err)
	}

	data, err := readConfigFile(p, newConfig)
	if err != nil {
		return nil, fmt.Errorf("fail to read config file: %w", err)
	}

	c, err := parseConfigFile(data)
	if err != nil {
		return nil, fmt.Errorf("fail to parse config file: %w", err)
	}

	return c, nil
}

func readConfigFile(path string, newConfig bool) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if !newConfig {
			return nil, fmt.Errorf("fail to read config file from path: %w", err)
		}
		if newConfig {
			err = createDefaultConfigFile(path)
			if err != nil {
				return nil, fmt.Errorf("fail to create dafault config file: %w", err)
			}

			// retry
			data, err = os.ReadFile(path)
		}
	}

	return data, nil
}

func createDefaultConfigFile(path string) error {
	defaultConfig := &FileSchemaJSON{
		Enrolled: make([]EnrollSchemaJSON, 0),
	}
	data, err := json.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("fail to make default json config: %w", err)
	}

	err = os.WriteFile(path, data, 0777)
	if err != nil {
		return fmt.Errorf("fail to create default confnig file: %w", err)
	}

	return nil
}

type FileSchemaJSON struct {
	Enrolled []EnrollSchemaJSON `json:"enrolled"`
}

type EnrollSchemaJSON struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func parseConfigFile(data []byte) (*Config, error) {
	f := &FileSchemaJSON{}

	err := json.Unmarshal(data, f)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal json from file: %w", err)
	}

	err = validateConfig(f)
	if err != nil {
		return nil, fmt.Errorf("fail to pass validate config: %w", err)
	}

	c := &Config{
		enrolled: make(map[string]BinaryDirectory, 10),
	}

	for i := range f.Enrolled {
		name := f.Enrolled[i].Name
		path, err := filepath.Abs(f.Enrolled[i].Path)
		if err != nil { // double check
			return nil, fmt.Errorf("can't not convert path to absolute path: path=%v", name)
		}

		c.enrolled[name] = BinaryDirectory{
			name: name,
			path: path,
		}
	}

	return c, nil
}

func validateConfig(f *FileSchemaJSON) error {
	if f == nil {
		return fmt.Errorf("fail to validate config file: file is nil")
	}

	nameSet := map[string]struct{}{}

	for i := range f.Enrolled {
		name := f.Enrolled[i].Name
		path := f.Enrolled[i].Path

		_, exist := nameSet[name]
		if exist {
			return fmt.Errorf("detected duplicated name in validate: name=%v", name)
		}

		_, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("can't not convert path to absolute path: path=%v", name)
		}
	}

	return nil
}
