package schema

// import "github.com/at0x0ft/museum/internal/pkg/debug"    // 4debug
import (
    "bytes"
    _ "embed"
    "fmt"
    "os"
    "path/filepath"
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/internal/pkg/jsonc"
)

const (
    SeedFilename = "seed.yml"
    DevContainerFileName string = "devcontainer.json"
)

type Configs struct {
    VSCodeDevcontainer yaml.Node `yaml:"vscode_devcontainer"`
    DockerCompose yaml.Node `yaml:"docker_compose,omitempty"`
}

type Seed struct {
    Version string `yaml:"version"`
    Variables yaml.Node `yaml:"variables"`
    Configs Configs `yaml:"configs"`
}

//go:embed seed.common.yml
var commonSeedRawData []byte

func LoadSeed(dirPath string) (*Seed, error) {
    filePath := filepath.Join(dirPath, SeedFilename)
    buf, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var data *Seed
    if err := yaml.Unmarshal(buf, &data); err != nil {
        return nil, err
    }
    return data, nil
}

func GetCommonSeedData(variables *yaml.Node) (*Seed, error) {
    var data Seed
    if err := yaml.Unmarshal(commonSeedRawData, &data); err != nil {
        return nil, err
    }
    data.Variables = *variables
    return &data, nil
}

func (self *Seed) WriteToFile(dirPath string) error {
    filePath := filepath.Join(dirPath, SeedFilename)
    var buf bytes.Buffer
    yamlEncoder := yaml.NewEncoder(&buf)
    defer yamlEncoder.Close()
    yamlEncoder.SetIndent(2)

    yamlEncoder.Encode(&self)
    if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
        return err
    }
    return nil
}

func (self *Seed) WriteDevcontainer(dirPath string) error {
    jsoncContent, err := jsonc.Encode(&self.Configs.VSCodeDevcontainer, 4)
    if err != nil {
        return err
    }

    filePath := filepath.Join(dirPath, DevContainerFileName)
    if err := os.WriteFile(filePath, []byte(jsoncContent), 0644); err != nil {
        return err
    }
    return nil
}

func (self *Seed) WriteDockerCompose(dirPath string) error {
    filePath := filepath.Join(dirPath, DockerComposeFileName)
    var buf bytes.Buffer
    yamlEncoder := yaml.NewEncoder(&buf)
    defer yamlEncoder.Close()
    yamlEncoder.SetIndent(2)

    yamlEncoder.Encode(&self.Configs.DockerCompose)
    if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
        return err
    }
    return nil
}

func (self *Seed) findValueNodeInMapping(mapping *yaml.Node, keyName string) (*yaml.Node, error) {
    for index := 0; index < len(mapping.Content); index += 2 {
        keyNode, valueNode := mapping.Content[index], mapping.Content[index + 1]
        if keyNode.Value == keyName {
            return valueNode, nil
        }
    }
    return nil, fmt.Errorf("[Error] Cannot find '%s' key node in mapping!", keyName)
}

func (self *Seed) getCommonArgumentsRoot() (*yaml.Node, error) {
    commonServiceValueNode, err := self.findValueNodeInMapping(
        &self.Variables,
        COMMON_COLLECTION_NAME_KEY,
    )
    if err != nil {
        return nil, err
    }
    commonArgumentsValueNode, err := self.findValueNodeInMapping(
        commonServiceValueNode,
        ARGUMENTS_KEY,
    )
    if err != nil {
        return nil, err
    }
    return commonArgumentsValueNode, nil
}

func (self *Seed) GetComposeProjectPrefix() (string, error) {
    commonArgumentsRootNode, err := self.getCommonArgumentsRoot()
    if err != nil {
        return "", err
    }
    arguments := Arguments(*commonArgumentsRootNode)
    return (&arguments).getComposeProjectPrefix()
}

func (self *Seed) FilterDockerCompose() {
    var emptyContent []*yaml.Node
    self.Configs.DockerCompose.Content = emptyContent
}
