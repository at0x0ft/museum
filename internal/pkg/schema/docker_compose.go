package schema

// import "fmt"    // 4debug
import (
    "bytes"
    "os"
    "strings"
    "os/exec"
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
    Volumes map[string]map[string]interface{} `yaml:"volumes,omitempty"`
}

type DockerComposeService struct {
    Build struct {
        Context string `yaml:"context,omitempty"`
        Dockerfile string `yaml:"dockerfile,omitempty"`
        Args []string `yaml:"args,omitempty"`
        Target string `yaml:"target,omitempty"`
    } `yaml:"build,omitempty"`
    WorkingDir string `yaml:"working_dir,omitempty"`
    Volumes []string `yaml:"volumes,omitempty"`
    EntryPoint string `yaml:"entrypoint,omitempty"`
    Command string `yaml:"command,omitempty"`
    User string `yaml:"user,omitempty"`
    // TODO: append later
    // TODO: x-* anchor path conversion should support
}

func loadDockerComposeByteData(data []byte) (*DockerCompose, error) {
    var res DockerCompose
    if err := yaml.Unmarshal(data, &res); err != nil {
        return nil, err
    }
    return &res, nil
}

func LoadMultipleDockerComposes(dockerComposeFileList []string) (*DockerCompose, error) {
    var args []string
    for _, dockerComposeFile := range dockerComposeFileList {
        args = append(args, "-f", dockerComposeFile)
    }
    args = append(args, "config")
    cmd := exec.Command("docker-compose", args...)
    // TODO: split stdout & stderr
    out, err := cmd.CombinedOutput()
    if err != nil {
        return nil, err
    }
    return loadDockerComposeByteData(out)
}

func ConvertDockerComposeYamlToStruct(root *yaml.Node) (*DockerCompose, error) {
    var buf bytes.Buffer
    yamlEncoder := yaml.NewEncoder(&buf)
    defer yamlEncoder.Close()
    yamlEncoder.SetIndent(2)
    yamlEncoder.Encode(&root)
    return loadDockerComposeByteData(buf.Bytes())
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
