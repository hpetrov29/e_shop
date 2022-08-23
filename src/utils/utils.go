package utils

import "os"

//The GetEnv function gets the value of the environment variable stored under the first argument key-name; if not defined it's assigned a default value (the second function argument).
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
