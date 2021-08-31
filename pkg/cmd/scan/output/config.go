package output

import "fmt"

type OutputConfig struct {
	Key     string
	Options map[string]string
}

func (o *OutputConfig) String() string {
	return fmt.Sprintf("%s://%s", o.Key, o.Options["path"])
}
