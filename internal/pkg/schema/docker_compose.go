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

type DockerComposeRootVolume map[string]map[string]interface{}

type DockerCompose struct {
    Version string `yaml:"version"`
    Services map[string]DockerComposeService `yaml:"services"`
    Volumes DockerComposeRootVolume `yaml:"volumes,omitempty"`
}

type DockerComposeServiceBuild struct {
    Context string `yaml:"context,omitempty"`
    Dockerfile string `yaml:"dockerfile,omitempty"`
    Args []string `yaml:"args,omitempty"`
    Target string `yaml:"target,omitempty"`
}

type DockerComposeServiceVolume string

type DockerComposeService struct {
    Build DockerComposeServiceBuild `yaml:"build,omitempty"`
    WorkingDir string `yaml:"working_dir,omitempty"`
    Volumes []DockerComposeServiceVolume `yaml:"volumes,omitempty"`
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

func convertToHostAbsPathIfPathIsRel(path, hostDirPath string) (string, error) {
    if filepath.IsAbs(path) {
        return path, nil
    }
    absPath, err := filepath.Abs(filepath.Join(hostDirPath, path))
    if err != nil {
        return "", err
    }
    return absPath, nil
}

func (self *DockerComposeServiceBuild) convertRelPathToAbs(
    dirPath string,
) (*DockerComposeServiceBuild, error) {
    result := *self
    if self.Context != "" {
        absPath, err := convertToHostAbsPathIfPathIsRel(self.Context, dirPath)
        if err != nil {
            return nil, err
        }
        result.Context = absPath
    }
    if self.Dockerfile != "" {
        absPath, err := convertToHostAbsPathIfPathIsRel(self.Dockerfile, dirPath)
        if err != nil {
            return nil, err
        }
        result.Dockerfile = absPath
    }
    return &result, nil
}

func (self *DockerComposeServiceVolume) isVolumeMounted(
    headPart string,
    rootVolumes DockerComposeRootVolume,
) bool {
    _, ok := rootVolumes[headPart]
    return ok
}

func (self *DockerComposeServiceVolume) convertRelPathToAbs(
    rootVolumes DockerComposeRootVolume,
    dirPath string,
) (*DockerComposeServiceVolume, error) {
    volumeParts := strings.Split(string(*self), DockerComposeVolumeDelimiter)
    if self.isVolumeMounted(volumeParts[0], rootVolumes) {
        return self, nil
    }
    hostPath := volumeParts[0]

    hostAbsPath, err := convertToHostAbsPathIfPathIsRel(hostPath, dirPath)
    if err != nil {
        return nil, err
    }
    var convertedParts []string
    convertedParts = append(convertedParts, hostAbsPath)
    convertedParts = append(convertedParts, volumeParts[1:]...)
    result := DockerComposeServiceVolume(strings.Join(convertedParts, DockerComposeVolumeDelimiter))
    return &result, nil
}

func (self *DockerComposeService) convertRelPathToAbs(
    rootVolumes DockerComposeRootVolume,
    dirPath string,
) (*DockerComposeService, error) {
    result := *self
    convertedBuild, err := self.Build.convertRelPathToAbs(dirPath)
    if err != nil {
        return nil, err
    }
    result.Build = *convertedBuild

    for idx, volume := range self.Volumes {
        convertedVolume, err := volume.convertRelPathToAbs(rootVolumes, dirPath)
        if err != nil {
            return nil, err
        }
        result.Volumes[idx] = *convertedVolume
    }
    return &result, nil
}

func (self *DockerCompose) ConvertRelPathToAbs(dirPath string) (*DockerCompose, error) {
    result := *self
    for name, service := range self.Services {
        convertedService, err := service.convertRelPathToAbs(self.Volumes, dirPath)
        if err != nil {
            return nil, err
        }
        result.Services[name] = *convertedService
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
