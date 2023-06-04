package schema

import (
    "os"
    "encoding/json"
    "path/filepath"
)

const (
    DevContainerFileName string = "devcontainer.json"
)

type DevContainer struct {
    Name string `json:"name"`
    DockerComposeFile []string `json:dockerComposeFile`
    Service string `json:"service"`
    WorkspaceFolder string `json:"workspaceFolder"`
    RemoteEnv map[string]string `json:"remoteEnv"`
    RemoteUser string `json:"remoteUser"`
    // == deprecated ==
    // TODO: map[string](string|int)
    Settings map[string]string `json:"settings"`
    Extensions []string `json:"extensions"`
}

func LoadDevcontainer(dirPath string) (*DevContainer, error) {
    devcontainerPath := filepath.Join(dirPath, DevContainerFileName)
    buf, err := os.ReadFile(devcontainerPath)
    if err != nil {
        return nil, err
    }

    var data DevContainer
    if err := json.Unmarshal(buf, &data); err != nil {
        return nil, err
    }
    return &data, nil
}
