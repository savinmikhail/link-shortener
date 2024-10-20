package repository

import (
	"encoding/json"
	"fmt"
	"os"
)

func SaveMappedUrls(mappedUrls map[string]string) error {
	jsonContent, err := json.Marshal(mappedUrls)
	if err != nil {
		return err
	}
	err = os.WriteFile("list.json", jsonContent, 0777)
	if err != nil {
		return err
	}

	return nil
}

func GetOrigUrlByShortCode(shortCode string) (string, error) {
	mappedUrls, err := GetMappedUrls()
	if err != nil {
		return "", err
	}
	origUrl, exists := mappedUrls[shortCode]
	if !exists {
		return "", fmt.Errorf("short code %s not found", shortCode)
	}
	return origUrl, nil
}

func GetMappedUrls() (map[string]string, error) {
	fileContent, err := os.ReadFile("list.json")
	if err != nil {
		return nil, err
	}
	mappedUrls := make(map[string]string)
	err = json.Unmarshal(fileContent, &mappedUrls)
	if err != nil {
		return nil, err
	}
	return mappedUrls, nil
}
