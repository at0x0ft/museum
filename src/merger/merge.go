package merger

import "fmt"    // 4debug
// import "github.com/at0x0ft/museum/debug"    // 4debug
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

    fmt.Println(seedMetadataList)  // 4debug

    mergedVariables := mergeVariables(seedMetadataList)

    // TODO: merge Configs
    // TODO: replace Tag paths

    // r, err := visitableFactory("", root)
    // if err != nil {
    //     return nil, err
    // }
    // return r.visit(variables)
    mergedSeed := &schema.Seed{
        Version: "0",
        Variables: *mergedVariables,
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

func mergeVariables(seedMetadataList []seedMetadata) *yaml.Node {
    newNode := createNewMappingNode()
    for _, e := range seedMetadataList {
        keyNode := &yaml.Node{
            Kind: yaml.ScalarNode,
            Tag: "!!str",
            Value: e.Name,
        }
        newNode.Content = append(
            newNode.Content,
            keyNode,
            &e.Data.Variables,
        )
    }
    return newNode
}

func createNewMappingNode() *yaml.Node {
    return &yaml.Node{
        Kind: yaml.MappingNode,
        Tag: "!!map",
    }
}
