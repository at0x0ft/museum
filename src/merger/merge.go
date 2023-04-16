package merger

import (
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/schema"
)

type seedMetadata struct {
    Name string
    Data *schema.Seed
}

func Merge(skeleton *schema.Skeleton) (*schema.Seed, error) {
    seedMetadataList, err := loadSeeds(skeleton)
    if err != nil {
        return nil, err
    }

    mergedVariables, err := mergeVariables(seedMetadataList)
    if err != nil {
        return nil, err
    }

    mergedConfigs, err := mergeConfigs(seedMetadataList)
    if err != nil {
        return nil, err
    }

    mergedSeed := &schema.Seed{
        Version: "0",
        Variables: *mergedVariables,
        Configs: *mergedConfigs,
    }
    return mergedSeed, nil
}

func loadSeeds(skeleton *schema.Skeleton) ([]seedMetadata, error) {
    var result []seedMetadata
    for _, collection := range skeleton.Collections {
        seed, err := schema.LoadSeed(collection.Path)
        if err != nil {
            return nil, err
        }
        newSeedMetadata := seedMetadata{
            Name: collection.Name,
            Data: seed,
        }
        result = append(result, newSeedMetadata)
    }
    return result, nil
}

func mergeVariables(seedMetadataList []seedMetadata) (*yaml.Node, error) {
    rootNodePath := ""
    appendedNodes := make(map[string]visitable)
    var err error
    for _, seedMetadata := range seedMetadataList {
        newRoot := createNewMappingNode()
        keyNode := &yaml.Node{
            Kind: yaml.ScalarNode,
            Tag: "!!str",
            Value: seedMetadata.Name,
        }
        newRoot.Content = append(
            newRoot.Content,
            keyNode,
            &seedMetadata.Data.Variables,
        )
        appendedNodes, err = appendTree(
            appendedNodes,
            newRoot,
            seedMetadata.Name,
        );
        if err != nil {
            return nil, err
        }
    }
    return appendedNodes[rootNodePath].getRaw(), nil
}

func createNewMappingNode() *yaml.Node {
    return &yaml.Node{
        Kind: yaml.MappingNode,
        Tag: "!!map",
    }
}

func mergeConfigs(seedMetadataList []seedMetadata) (*schema.Configs, error) {
    var mergedDevContainerRoot, mergedDockerComposeRoot *yaml.Node
    var err error
    mergedDevContainerRoot, err = mergeDevcontainerConfigs(seedMetadataList)
    if err != nil {
        return nil, err
    }
    mergedDockerComposeRoot, err = mergeDockerComposes(seedMetadataList)
    if err != nil {
        return nil, err
    }
    mergedConfigs := &schema.Configs{
        VSCodeDevcontainer: *mergedDevContainerRoot,
        DockerCompose: *mergedDockerComposeRoot,
    }
    return mergedConfigs, nil
}

func mergeDevcontainerConfigs(seedMetadataList []seedMetadata) (*yaml.Node, error) {
    rootNodePath := ""
    appendedConfigs := make(map[string]visitable)
    var err error
    for _, seedMetadata := range seedMetadataList {
        appendedConfigs, err = appendTree(
            appendedConfigs,
            &seedMetadata.Data.Configs.VSCodeDevcontainer,
            seedMetadata.Name,
        )
        if err != nil {
            return nil, err
        }
    }
    return appendedConfigs[rootNodePath].getRaw(), nil
}

func mergeDockerComposes(seedMetadataList []seedMetadata) (*yaml.Node, error) {
    rootNodePath := ""
    appendedConfigs := make(map[string]visitable)
    var err error
    for _, seedMetadata := range seedMetadataList {
        appendedConfigs, err = appendTree(
            appendedConfigs,
            &seedMetadata.Data.Configs.DockerCompose,
            seedMetadata.Name,
        )
        if err != nil {
            return nil, err
        }
    }
    return appendedConfigs[rootNodePath].getRaw(), nil
}

func appendTree(appendedNodes map[string]visitable, root *yaml.Node, collectionName string) (map[string]visitable, error) {
    r, err := visitableFactory("", root)
    if err != nil {
        return nil, err
    }
    _, err = r.visit(appendedNodes, collectionName)
    return appendedNodes, err
}
