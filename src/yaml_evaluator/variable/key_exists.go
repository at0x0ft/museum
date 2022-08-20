package variable

import (
    "fmt"
    "strings"
)

func KeyExists(key string, variables *map[string]map[string]string) (string, error) {
    keys := strings.Split(key, ".")
    if keyLength := len(keys); keyLength != 2 {
        return "", fmt.Errorf("Variable key error: key length (=%d) != 2.", keyLength)
    }

    if firstKeyVariables, ok1 := (*variables)[keys[0]]; ok1 {
        if variable, ok2 := firstKeyVariables[keys[1]]; ok2 {
            return variable, nil
        } else {
            return "", fmt.Errorf("Variable key error: second key = '%s' not found.", keys[1])
        }
    } else {
        return "", fmt.Errorf("Variable key error: first key = '%s' not found.", keys[0])
    }
}
