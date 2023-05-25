package schema

// import "github.com/at0x0ft/museum/internal/pkg/debug"    // 4debug
import (
    "fmt"
    // "bytes"
    "path/filepath"
    "io/ioutil"
)

// TODO: replace with viper or other .env file r/w tools...

const (
    ComposeConfigFilename = "compose_config"
    ComposeConfigLinkDstFilename = ".env"
    ComposeNameKey = "COMPOSE_NAME"
)

type DockerComposeConfig struct {
    ComposeName string
}

func CreateComposeConfig(composeName string) *DockerComposeConfig {
    return &DockerComposeConfig{
        ComposeName: composeName,
    }
}

func (self *DockerComposeConfig) Write(dirPath string) error {
    filePath := self.GetFilepath(dirPath)
    content := fmt.Sprintf("%s='%s'\n", ComposeNameKey, self.ComposeName)
    if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
        return err
    }
    return nil
}

func (self *DockerComposeConfig) GetFilepath(dirPath string) string {
    return filepath.Join(dirPath, ComposeConfigFilename)
}
