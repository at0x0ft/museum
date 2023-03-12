package schema

import (
    "path/filepath"
    "io/ioutil"
    "gopkg.in/yaml.v3"
)

const (
    SkeletonFilename = "skeleton.yml"
    DockerFileDirectory = "./docker"
)

type Collection struct {
    Name string `yaml:"name"`
    Path string `yaml:"path"`
}

type Skeleton struct {
    Version string `yaml:"version"`
    CollectionsPath string `yaml:"collections_path"`
    Collections []Collection `yaml:"collections"`
}

func LoadSkeleton(dirPath string) (*Skeleton, error) {
    fileAbsPath, err := filepath.Abs(filepath.Join(dirPath, SkeletonFilename))
    if err != nil {
        return nil, err
    }

    buf, err := ioutil.ReadFile(fileAbsPath)
    if err != nil {
        return nil, err
    }

    var data *Skeleton
    if err := yaml.Unmarshal(buf, &data); err != nil {
        return nil, err
    }
    data.CollectionsPath = resolvePath(data.CollectionsPath, fileAbsPath)
    var collections []Collection
    for _, collection := range data.Collections {
        newCollectionConfig := Collection{
            Name: collection.Name,
            Path: resolvePath(collection.Path, data.CollectionsPath),
        }
        collections = append(collections, newCollectionConfig)
    }
    data.Collections = collections
    return data, nil
}

func resolvePath(targetPath, baseAbsPath string) string {
    if filepath.IsAbs(targetPath) {
        return targetPath
    }
    return filepath.Join(baseAbsPath, targetPath)
}
