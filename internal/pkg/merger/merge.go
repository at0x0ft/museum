package merger

// import "fmt"    // 4debug
// import "github.com/at0x0ft/museum/internal/pkg/debug"   // 4debug
import (
    "gopkg.in/yaml.v3"
    "github.com/at0x0ft/museum/internal/pkg/schema"
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

func getCommonVariables(skeleton *schema.Skeleton) (*yaml.Node, error) {
    argumentsKey := skeleton.Arguments.GetKey()
    argumentsValue := yaml.Node(skeleton.Arguments)

    variablesNode := createNewMappingNode()
    variablesNode.Content = append(
        variablesNode.Content,
        argumentsKey,
        &argumentsValue,
    )
    return variablesNode, nil
}

func getCommonSeedMetadata(skeleton *schema.Skeleton) (*seedMetadata, error) {
    commonVariables, err := getCommonVariables(skeleton)
    if err != nil {
        return nil, err
    }

    commonSeedData, err := schema.GetCommonSeedData(commonVariables)
    if err != nil {
        return nil, err
    }

    return &seedMetadata{
        Name: schema.COMMON_COLLECTION_NAME_KEY,
        Data: commonSeedData,
    }, nil
}

func loadSeeds(skeleton *schema.Skeleton) ([]seedMetadata, error) {
    var result []seedMetadata
    commonSeedMetadata, err := getCommonSeedMetadata(skeleton)
    if err != nil {
        return nil, err
    }
    result = append(result, *commonSeedMetadata)

    var restSeeds []seedMetadata
    for _, collection := range skeleton.Collections.List {
        seed, err := schema.LoadSeed(collection.Path)
        if err != nil {
            return nil, err
        }

        if collection.NoCompose {
            seed.FilterDockerCompose()
        }
        newSeedMetadata := seedMetadata{
            Name: collection.Name,
            Data: seed,
        }

        commonAttachedCollectionName, err := skeleton.Arguments.GetAttachServiceName()
        if err != nil {
            return nil, err
        }
        if newSeedMetadata.Name == commonAttachedCollectionName {
            result = append(result, newSeedMetadata)
        } else {
            restSeeds = append(restSeeds, newSeedMetadata)
        }
    }
    result = append(result, restSeeds...)
    return result, nil
}

func appendTree(appendedNodes map[string]visitable, collectionName string, root *yaml.Node) (map[string]visitable, error) {
    r, err := visitableFactory("", root)
    if err != nil {
        return nil, err
    }
    _, err = r.visit(appendedNodes, collectionName)
    return appendedNodes, err
}

func (self *seedMetadata) getCollectionNameAddedVariables() *yaml.Node {
    root := createNewMappingNode()
    keyNode := &yaml.Node{
        Kind: yaml.ScalarNode,
        Tag: "!!str",
        Value: self.Name,
    }
    root.Content = append(
        root.Content,
        keyNode,
        &self.Data.Variables,
    )
    return root
}

func mergeVariables(seedMetadataList []seedMetadata) (*yaml.Node, error) {
    rootNodePath := ""
    appendedNodes := make(map[string]visitable)
    var err error
    for _, seedMetadata := range seedMetadataList {
        variables := seedMetadata.getCollectionNameAddedVariables()
        appendedNodes, err = appendTree(
            appendedNodes,
            seedMetadata.Name,
            variables,
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
            seedMetadata.Name,
            &seedMetadata.Data.Configs.VSCodeDevcontainer,
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
            seedMetadata.Name,
            &seedMetadata.Data.Configs.DockerCompose,
        )
        if err != nil {
            return nil, err
        }
    }
    return appendedConfigs[rootNodePath].getRaw(), nil
}
