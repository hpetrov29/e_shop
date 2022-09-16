package utils

import (
	"net/url"
	"os"
	"path"

	"github.com/google/uuid"
)

//The GetEnv function gets the value of the environment variable stored under the first argument key-name; if not defined it's assigned a default value (the second function argument).
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

//Parses the given URL string and returns the name of the object
func ObjNameFromUrl(imageUrl string) (string, error) {
	if imageUrl == "" {
		objId, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		return objId.String(), nil
	}

	urlPath, err := url.Parse(imageUrl)
	if err != nil {
		return "", err
	}

	return path.Base(urlPath.Path), nil
}
