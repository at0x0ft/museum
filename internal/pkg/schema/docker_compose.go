package schema

// import "fmt"    // 4debug
import (
    "bytes"
    "os"
    "strings"
    "path/filepath"
    "gopkg.in/yaml.v3"
)

const (
    DockerComposeFileName = "docker-compose.yml"
    DockerComposeVolumeDelimiter = ":"
)

type DockerCompose struct {
    Version string `yaml:"version"`
    Services map[string]DockerComposeService `yaml:"services"`
    Volumes map[string]string `yaml:"volumes"`
}

type DockerComposeService struct {
    Build struct {
        Context string `yaml:"context"`
        Dockerfile string `yaml:"dockerfile"`
        Args []string `yaml:"args"`
        Target string `yaml:"target"`
    } `yaml:"build"`
    WorkingDir string `yaml:"working_dir"`
    Volumes []string `yaml:"volumes"`
    EntryPoint string `yaml:"entrypoint"`
    Command string `yaml:"command"`
    User string `yaml:"user"`
    // TODO: append later
    // TODO: x-* anchor path conversion should support
}

func ConvertDockerComposeYamlToStruct(root *yaml.Node) (*DockerCompose, error) {
    var buf bytes.Buffer
    yamlEncoder := yaml.NewEncoder(&buf)
    defer yamlEncoder.Close()
    yamlEncoder.SetIndent(2)
    yamlEncoder.Encode(&root)

    var data *DockerCompose
    if err := yaml.Unmarshal(buf.Bytes(), &data); err != nil {
        return nil, err
    }
    return data, nil
}

func (self *DockerCompose) isVolumeMounted(volume string) bool {
    headPart := strings.Split(volume, DockerComposeVolumeDelimiter)[0]
    _, ok := self.Volumes[headPart]
    return ok
}

func (self *DockerCompose) convertHostRelpath(volume string, dirPath string) (string, error) {
    if self.isVolumeMounted(volume) {
        return volume, nil
    }
    volumeParts := strings.Split(volume, DockerComposeVolumeDelimiter)
    hostPath := volumeParts[0]
    if filepath.IsAbs(hostPath) {
        return volume, nil
    }

    hostAbsPath, err := filepath.Abs(filepath.Join(dirPath, hostPath))
    if err != nil {
        return "", err
    }
    var convertedParts []string
    convertedParts = append(convertedParts, hostAbsPath)
    convertedParts = append(convertedParts, volumeParts[1:]...)
    return strings.Join(convertedParts, DockerComposeVolumeDelimiter), nil
}

func (self *DockerCompose) ConvertVolumeRelpathToAbs(dirPath string) (*DockerCompose, error) {
    result := *self
    for name, service := range self.Services {
        for idx, volume := range service.Volumes {
            convertedVol, err := self.convertHostRelpath(volume, dirPath)
            if err != nil {
                return nil, err
            }
            result.Services[name].Volumes[idx] = convertedVol
        }
    }
    return &result, nil
}

func (self *DockerCompose) Write(dirPath string) error {
    filePath := filepath.Join(dirPath, DockerComposeFileName)
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
