package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Int64(v int64) *int64 {
	return &v
}
func String(v string) *string {
	return &v
}

func Bool(v bool) *bool {
	return &v
}

func GetCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
