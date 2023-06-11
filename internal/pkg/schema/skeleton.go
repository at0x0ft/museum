package schema

import (
    "fmt"
    "os"
    "strings"
    "path/filepath"
    "gopkg.in/yaml.v3"
)

const (
    SkeletonFilename = "skeleton.yml"
    DockerFileDirectory = "./docker"
)

// TODO: rename const names.
const (
    COMMON_COLLECTION_NAME_KEY = "common"
    ARGUMENTS_KEY = "arguments"
    VSCODE_DEVCONTAINER_KEY = ARGUMENTS_KEY + ".vscode_devcontainer"
    DEVCONTAINER_PROJECT_NAME_KEY = VSCODE_DEVCONTAINER_KEY + ".project_name"
    DEVCONTAINER_ATTACH_SERVICE_KEY = VSCODE_DEVCONTAINER_KEY + ".attach_service"
    DEVCONTAINER_SOURCE_PATH_KEY = VSCODE_DEVCONTAINER_KEY + ".source_path"
    DOCKER_COMPOSE_KEY = ARGUMENTS_KEY + ".docker_compose"
    DOCKER_COMPOSE_PROJECT_PREFIX_KEY = DOCKER_COMPOSE_KEY + ".project_prefix"
    DOCKER_COMPOSE_FILES_KEY = DOCKER_COMPOSE_KEY + ".files"
    VSCODE_EXTENSION_VOLUMES_KEY = DOCKER_COMPOSE_KEY + ".vscode_extension_volumes"
    VSCODE_NORMAL_EXTENSION_VOLUME_NAME_KEY = VSCODE_EXTENSION_VOLUMES_KEY + ".normal"
    VSCODE_INSIDER_EXTENSION_VOLUME_NAME_KEY = VSCODE_EXTENSION_VOLUMES_KEY + ".insider"
)

// type Arguments struct {
//     VSCodeDevcontainer struct {
//         ProjectName string `yaml:"project_name"`
//         AttachService string `yaml:"attach_service"`
//         SourcePath string `yaml:"source_path"`
//     } `yaml:"vscode_devcontainer"`
//     DockerCompose struct {
//         ProjectPrefix string `yaml:"project_prefix"`
//         Files []string `yaml:"files"`
//         VSCodeExtensionVolumes struct {
//             Normal string `yaml:"normal"`
//             Insider string `yaml:"insider"`
//         } `yaml:"vscode_extension_volumes"`
//     } `yaml:"docker_compose"`
// }

type Arguments yaml.Node

type Collection struct {
    Name string `yaml:"name,omitempty"`
    Path string `yaml:"path"`
    // TODO: default value is false
    NoCompose bool `yaml:"no_compose"`
}

type Collections struct {
    Path string `yaml:"path"`
    List []Collection `yaml:"list"`
}

type Skeleton struct {
    Version string `yaml:"version"`
    Arguments Arguments `yaml:"arguments"`
    Collections Collections `yaml:"collections"`
}

func (self *Arguments) UnmarshalYAML(n *yaml.Node) error {
    self.Kind = n.Kind
    self.Style = n.Style
    self.Tag = n.Tag
    self.Value = n.Value
    self.Anchor = n.Anchor
    self.Alias = n.Alias
    self.Content = n.Content
    self.HeadComment = n.HeadComment
    self.LineComment = n.LineComment
    self.FootComment = n.FootComment
    self.Line = n.Line
    self.Column = n.Column
    return nil
}

func LoadSkeleton(dirPath string) (*Skeleton, error) {
    dirAbsPath, err := filepath.Abs(dirPath)
    if err != nil {
        return nil, err
    }
    fileAbsPath := filepath.Join(dirAbsPath, SkeletonFilename)
    buf, err := os.ReadFile(fileAbsPath)
    if err != nil {
        return nil, err
    }

    var data Skeleton
    if err := yaml.Unmarshal(buf, &data); err != nil {
        return nil, err
    }
    return data.getCanonical(dirAbsPath)
}

func resolvePath(baseAbsDir string, targetPath string) string {
    if filepath.IsAbs(targetPath) {
        return targetPath
    }
    return filepath.Join(baseAbsDir, targetPath)
}

func (self *Skeleton) serviceExists(serviceName string) bool {
    serviceNameSet := map[string]struct{}{}
    for _, collection := range self.Collections.List {
        serviceNameSet[collection.Name] = struct{}{}
    }
    _, ok := serviceNameSet[serviceName]
    return ok
}

func (self *Skeleton) validate() error {
    attachedService, err := self.Arguments.GetAttachServiceName()
    if err != nil {
        return err
    }
    if !self.serviceExists(attachedService) {
        if err != nil {
            return err
        }
        return fmt.Errorf(
            "[Error] '%s' specified value (= '%s') collection is not found!",
            DEVCONTAINER_ATTACH_SERVICE_KEY,
            attachedService,
        )
    }
    return nil
}

func (self *Skeleton) getCanonical(baseAbsDir string) (*Skeleton, error) {
    if !filepath.IsAbs(baseAbsDir) {
        err := fmt.Errorf("[Error] baseAbsDir = '%v' is not absolute path.", baseAbsDir)
        return nil, err
    }

    collections, err := self.Collections.getCanonical(baseAbsDir)
    if err != nil {
        return nil, err
    }

    result := &Skeleton{
        Version: self.Version,
        Arguments: self.Arguments,
        Collections: *collections,
    }
    if err := result.validate(); err != nil {
        return nil, err
    }
    return result, nil
}

func (self *Skeleton) getDockerComposeValueNode(arguments *yaml.Node) (*yaml.Node, error) {
    splitKey := strings.Split(DOCKER_COMPOSE_KEY, ".")
    dockerComposeKeyNodeValue := splitKey[len(splitKey) - 1]
    for index := 0; index < len(arguments.Content); index += 2 {
        if arguments.Content[index].Value == dockerComposeKeyNodeValue {
            return arguments.Content[index + 1], nil
        }
    }
    return nil, fmt.Errorf(
        "[Error] '%s' is not set in skeleton.yml!",
        DOCKER_COMPOSE_KEY,
    )
}

func (self *Skeleton) findVscodeExtensionVolumesKeyIndex(dockerComposeVolumeNode *yaml.Node) (int, error) {
    splitKey := strings.Split(VSCODE_EXTENSION_VOLUMES_KEY, ".")
    vscodeExtensionVolumes := splitKey[len(splitKey) - 1]
    result := -1
    for index := 0; index < len(dockerComposeVolumeNode.Content); index += 2 {
        if dockerComposeVolumeNode.Content[index].Value == vscodeExtensionVolumes {
            result = index
            break
        }
    }
    if result == -1 {
        return result, fmt.Errorf(
            "[Error] '%s' is not found in Skeleton struct!",
            VSCODE_EXTENSION_VOLUMES_KEY,
        )
    }
    return result, nil
}

func (self *Skeleton) filterValueNotSetContents(mapping *yaml.Node) *yaml.Node {
    var filteredContents []*yaml.Node
    for index := 0; index < len(mapping.Content); index += 2 {
        if mapping.Content[index + 1].Value != "" {
            filteredContents = append(
                filteredContents,
                mapping.Content[index],
                mapping.Content[index + 1],
            )
        }
    }
    mapping.Content = filteredContents
    return mapping
}

func (self *Skeleton) removeKeyValueFromContent(content []*yaml.Node, keyIndex int) ([]*yaml.Node, error) {
    if keyIndex < 0 || keyIndex + 1 >= len(content) {
        return nil, fmt.Errorf("[Error] given key index = %v is out of range!", keyIndex)
    }
    return append(content[:keyIndex], content[keyIndex + 2:]...), nil
}

func (self *Skeleton) filterOptionalArguments(arguments *yaml.Node) (*yaml.Node, error) {
    dockerComposeValueNode, err := self.getDockerComposeValueNode(arguments)
    if err != nil {
        return nil, err
    }
    vscodeExtensionVolumesKeyIndex, err := self.findVscodeExtensionVolumesKeyIndex(dockerComposeValueNode)
    if err != nil {
        return nil, err
    }

    vscodeExtensionVolumesValueNode := dockerComposeValueNode.Content[vscodeExtensionVolumesKeyIndex + 1]
    vscodeExtensionVolumesValueNode = self.filterValueNotSetContents(vscodeExtensionVolumesValueNode)
    if len(vscodeExtensionVolumesValueNode.Content) == 0 {
        filteredContent, err := self.removeKeyValueFromContent(
            dockerComposeValueNode.Content,
            vscodeExtensionVolumesKeyIndex,
        )
        if err != nil {
            return nil, err
        }
        dockerComposeValueNode.Content = filteredContent
    }

    return arguments, nil
}

func (self *Arguments) validateVSCodeDevContainer(key, value *yaml.Node) error {
    if key.Value != "vscode_devcontainer" {
        return nil
    }
    for index := 0; index < len(value.Content); index += 2 {
        childKey := value.Content[index]
        childValue := value.Content[index + 1]
        if childKey.Value == "project_name" && childValue.Value == "" {
            return fmt.Errorf(
                "[Error] '%s' is not specified!",
                DEVCONTAINER_PROJECT_NAME_KEY,
            )
        } else if childKey.Value == "attach_service" && childValue.Value == "" {
            return fmt.Errorf(
                "[Error] '%s' is not specified!",
                DEVCONTAINER_ATTACH_SERVICE_KEY,
            )
        } else if childKey.Value == "source_path" && childValue.Value == "" {
            return fmt.Errorf(
                "[Error] '%s' is not specified!",
                DEVCONTAINER_SOURCE_PATH_KEY,
            )
        }
    }
    return nil
}

func (self *Arguments) validateDockerCompose(key, value *yaml.Node) error {
    if key.Value != "docker_compose" {
        return nil
    }
    for index := 0; index < len(value.Content); index += 2 {
        childKey := value.Content[index]
        childValue := value.Content[index + 1]
        if childKey.Value == "project_prefix" && childValue.Value == "" {
            return fmt.Errorf(
                "[Error] '%s' is not specified!",
                DOCKER_COMPOSE_PROJECT_PREFIX_KEY,
            )
        } else if childKey.Value == "files" {
            for _, n := range childKey.Content {
                path := n.Value
                if path == "" {
                    return fmt.Errorf(
                        "[Error] 'path' is not specified in '%s'!",
                        DOCKER_COMPOSE_FILES_KEY,
                    )
                }
            }
        }
    }
    return nil
}

func (self *Arguments) validate() error {
    for index := 0; index < len(self.Content); index += 2 {
        key := self.Content[index]
        value := self.Content[index + 1]
        if err := self.validateVSCodeDevContainer(key, value); err != nil {
            return err
        }
        if err := self.validateDockerCompose(key, value); err != nil {
            return err
        }
    }
    return nil
}

func (self *Arguments) GetKey() *yaml.Node {
    return &yaml.Node{Kind:yaml.ScalarNode, Value:ARGUMENTS_KEY}
}

func (self *Arguments) GetAttachServiceName() (string, error) {
    for i := 0; i < len(self.Content); i += 2 {
        key := self.Content[i]
        value := self.Content[i + 1]
        if key.Value == "vscode_devcontainer" {
            for j := 0; j < len(value.Content); j += 2 {
                childKey := value.Content[j]
                childValue := value.Content[j + 1]
                if childKey.Value == "attach_service" {
                    return childValue.Value, nil
                }
            }
        }
    }
    return "", fmt.Errorf(
        "[Error] '%s' is not found!",
        DOCKER_COMPOSE_PROJECT_PREFIX_KEY,
    )
}

func (self *Arguments) getComposeProjectPrefix() (string, error) {
    for i := 0; i < len(self.Content); i += 2 {
        key := self.Content[i]
        value := self.Content[i + 1]
        if key.Value == "docker_compose" {
            for j := 0; j < len(value.Content); j += 2 {
                childKey := value.Content[j]
                childValue := value.Content[j + 1]
                if childKey.Value == "project_prefix" {
                    return childValue.Value, nil
                }
            }
        }
    }
    return "", fmt.Errorf(
        "[Error] '%s' is not found!",
        DOCKER_COMPOSE_PROJECT_PREFIX_KEY,
    )
}

func (self *Collections) validate() error {
    if self.Path == "" {
        return fmt.Errorf("[Error] 'collections.path' is not specified!")
    }
    return nil
}

func (self *Collections) getCanonical(baseAbsDir string) (*Collections, error) {
    if !filepath.IsAbs(baseAbsDir) {
        err := fmt.Errorf("[Error] baseAbsDir = '%v' is not absolute path.", baseAbsDir)
        return nil, err
    }
    if err := self.validate(); err != nil {
        return nil, err
    }

    var canonicalList []Collection
    absPath := resolvePath(baseAbsDir, self.Path)
    for _, collection := range self.List {
        canonicalCollection, err := collection.getCanonical(absPath)
        if err != nil {
            return nil, err
        }
        canonicalList = append(
            canonicalList,
            *canonicalCollection,
        )
    }
    return &Collections{Path: absPath, List: canonicalList}, nil
}

func (self *Collection) validate() error {
    if self.Name == "" && self.Path == "" {
        return fmt.Errorf("[Error] neither 'name' nor 'path' is specified in 'collections.list'!")
    } else if self.Name == COMMON_COLLECTION_NAME_KEY {
        return fmt.Errorf(
            "[Error] collection's 'name' = '%s' is not allowed!",
            COMMON_COLLECTION_NAME_KEY,
        )
    }
    return nil
}

func (self *Collection) getCanonical(baseAbsDir string) (*Collection, error) {
    if !filepath.IsAbs(baseAbsDir) {
        err := fmt.Errorf("[Error] baseAbsDir = '%v' is not absolute path.", baseAbsDir)
        return nil, err
    }
    if err := self.validate(); err != nil {
        return nil, err
    }

    result := *self
    if self.Name == "" {
        result.Name = filepath.Base(self.Path)
    } else if self.Path == "" {
        result.Path = filepath.Join(".", self.Name)
    }
    result.Path = resolvePath(baseAbsDir, self.Path)
    return &result, nil
}
