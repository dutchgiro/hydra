package config

import (
	"io/ioutil"
	"os"
)

func WithTempFile(content string, fn func(string)) {
	f, _ := ioutil.TempFile("", "")
	f.WriteString(content)
	f.Close()
	defer os.Remove(f.Name())
	fn(f.Name())
}
