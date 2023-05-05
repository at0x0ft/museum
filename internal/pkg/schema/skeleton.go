package schema

import (
    "bytes"
    "fmt"
    "path/filepath"
    "io/ioutil"
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
    DOCKER_COMPOSE_KEY = ARGUMENTS_KEY + ".docker_compose"
    DOCKER_COMPOSE_PROJECT_PREFIX_KEY = DOCKER_COMPOSE_KEY + ".project_prefix"
    DOCKER_COMPOSE_FILES_KEY = DOCKER_COMPOSE_KEY + ".files"
    VSCODE_EXTENSION_VOLUMES_KEY = DOCKER_COMPOSE_KEY + ".vscode_extension_volumes"
    VSCODE_NORMAL_EXTENSION_VOLUME_NAME_KEY = VSCODE_EXTENSION_VOLUMES_KEY + ".normal"
    VSCODE_INSIDER_EXTENSION_VOLUME_NAME_KEY = VSCODE_EXTENSION_VOLUMES_KEY + ".insider"
)

type Arguments struct {
    VSCodeDevcontainer struct {
        ProjectName string `yaml:"project_name"`
        AttachService string `yaml:"attach_service"`
    } `yaml:"vscode_devcontainer"`
    DockerCompose struct {
        ProjectPrefix string `yaml:"project_prefix"`
        Files []string `yaml:"files"`
        VSCodeExtensionVolumes struct {
            Normal string `yaml:"normal"`
            Insider string `yaml:"insider"`
        } `yaml:"vscode_extension_volumes"`
    } `yaml:"docker_compose"`
}

type Collection struct {
    Name string `yaml:"name"`
    Path string `yaml:"path"`
    User string `yaml:"user"`
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

func LoadSkeleton(dirPath string) (*Skeleton, error) {
    dirAbsPath, err := filepath.Abs(dirPath)
    if err != nil {
        return nil, err
    }
    fileAbsPath := filepath.Join(dirAbsPath, SkeletonFilename)
    buf, err := ioutil.ReadFile(fileAbsPath)
    if err != nil {
        return nil, err
    }
    var data *Skeleton
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

func (self *Skeleton) attachedServiceExists() bool {
    serviceNameSet := map[string]struct{}{}
    for _, collection := range self.Collections.List {
        serviceNameSet[collection.Name] = struct{}{}
    }
    _, ok := serviceNameSet[self.GetCommonAttachedCollectionName()]
    return ok
}

func (self *Skeleton) validate() error {
    if !self.attachedServiceExists() {
        return fmt.Errorf(
            "[Error] '%s' specified value (= '%s') collection is not found!",
            DEVCONTAINER_ATTACH_SERVICE_KEY,
            self.GetCommonAttachedCollectionName(),
        )
    }
    return nil
}

func (self *Skeleton) getCanonical(baseAbsDir string) (*Skeleton, error) {
    if !filepath.IsAbs(baseAbsDir) {
        err := fmt.Errorf("[Error] baseAbsDir = '%v' is not absolute path.", baseAbsDir)
        return nil, err
    }

    arguments, err := self.Arguments.getCanonical(baseAbsDir)
    if err != nil {
        return nil, err
    }
    collections, err := self.Collections.getCanonical(baseAbsDir)
    if err != nil {
        return nil, err
    }

    result := &Skeleton{
        Version: self.Version,
        Arguments: *arguments,
        Collections: *collections,
    }
    if err := result.validate(); err != nil {
        return nil, err
    }
    return result, nil
}

func (self *Skeleton) GetRawArguments() (*yaml.Node, *yaml.Node, error) {
    var buf bytes.Buffer
    yamlEncoder := yaml.NewEncoder(&buf)
    defer yamlEncoder.Close()
    yamlEncoder.SetIndent(2)
    yamlEncoder.Encode(&self)

    var data yaml.Node
    if err := yaml.Unmarshal(buf.Bytes(), &data); err != nil {
        return nil, nil, err
    }
    var keyNode, valueNode *yaml.Node
    for index := 0; index < len(data.Content[0].Content); index += 2 {
        if data.Content[0].Content[index].Value == ARGUMENTS_KEY {
            keyNode = data.Content[0].Content[index]
            valueNode = data.Content[0].Content[index + 1]
            break
        }
    }
    if keyNode == nil && valueNode == nil {
        return nil, nil, fmt.Errorf(
            "[Error] '%s' is not found in skeleton.yml!",
            ARGUMENTS_KEY,
        )
    }
    return keyNode, valueNode, nil
}

func (self *Skeleton) GetCommonAttachedCollectionName() string {
    return self.Arguments.VSCodeDevcontainer.AttachService
}

func (self *Arguments) validate() error {
    if self.VSCodeDevcontainer.ProjectName == "" {
        return fmt.Errorf(
            "[Error] '%s' is specified!",
            DEVCONTAINER_PROJECT_NAME_KEY,
        )
    } else if self.VSCodeDevcontainer.AttachService == "" {
        return fmt.Errorf(
            "[Error] '%s' is specified!",
            DEVCONTAINER_ATTACH_SERVICE_KEY,
        )
    } else if self.DockerCompose.ProjectPrefix == "" {
        return fmt.Errorf(
            "[Error] '%s' is specified!",
            DOCKER_COMPOSE_PROJECT_PREFIX_KEY,
        )
    }
    for _, path := range self.DockerCompose.Files {
        if path == "" {
            return fmt.Errorf(
                "[Error] 'path' is not specified in '%s'!",
                DOCKER_COMPOSE_FILES_KEY,
            )
        }
    }
    return nil
}

func (self *Arguments) getCanonical(baseAbsDir string) (*Arguments, error) {
    if !filepath.IsAbs(baseAbsDir) {
        err := fmt.Errorf("[Error] baseAbsDir = '%v' is not absolute path.", baseAbsDir)
        return nil, err
    }
    if err := self.validate(); err != nil {
        return nil, err
    }

    result := *self
    var dockerComposeAbsPaths []string
    for _, path := range self.DockerCompose.Files {
        dockerComposeAbsPaths = append(
            dockerComposeAbsPaths,
            resolvePath(baseAbsDir, path),
        )
    }
    result.DockerCompose.Files = dockerComposeAbsPaths
    return &result, nil
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
