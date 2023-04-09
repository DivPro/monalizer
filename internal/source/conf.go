package source

type Conf struct {
	Headers map[string]string `yaml:"headers"`
	URLs    []string          `yaml:"urls"`
}
