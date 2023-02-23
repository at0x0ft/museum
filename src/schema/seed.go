package schema

import (
    "bytes"
    "path/filepath"
    "io/ioutil"
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/jsonc"
)

const (
    SeedFilename = "seed.yml"
    DevContainerFileName string = "devcontainer.json"
    DockerComposeFileName string = "docker-compose.yml"
)

type Configs struct {
    VSCodeDevcontainer yaml.Node `yaml:"vscode_devcontainer"`
    DockerCompose yaml.Node `yaml:"docker_compose"`
}

type Seed struct {
    Version string `yaml:"version"`
    Variables yaml.Node `yaml:"variables"`
    Configs Configs `yaml:"configs"`
}

func LoadSeeds(pathList map[string]string) (map[string]*Seed, error) {
    result := make(map[string]*Seed)
    for name, path := range pathList {
        seed, err := LoadSeed(path)
        if err != nil {
            return nil, err
        }
        result[name] = seed
    }
    return result, nil
}

func LoadSeed(dirPath string) (*Seed, error) {
    filePath := filepath.Join(dirPath, SeedFilename)
    buf, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var data *Seed
    if err := yaml.Unmarshal(buf, &data); err != nil {
        return nil, err
    }
    return data, nil
}

func WriteDevcontainer(data *yaml.Node, dirPath string) error {
    jsoncContent, err := jsonc.Encode(data, 4)
    if err != nil {
        return err
    }

    filePath := filepath.Join(dirPath, DevContainerFileName)
    if err := ioutil.WriteFile(filePath, []byte(jsoncContent), 0644); err != nil {
        return err
    }
    return nil
}

func WriteDockerCompose(data *yaml.Node, dirPath string) error {
    filePath := filepath.Join(dirPath, DockerComposeFileName)
    var buf bytes.Buffer
    yamlEncoder := yaml.NewEncoder(&buf)
    defer yamlEncoder.Close()
    yamlEncoder.SetIndent(2)

    yamlEncoder.Encode(data)
    if err := ioutil.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
        return err
    }
    return nil
}
