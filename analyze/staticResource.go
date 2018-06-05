package analyze

type StaticResources struct {
	Name    string        `yaml:"name"`
	Address SocketAddress `yaml:"address"`
	Route SocketAddress `yaml:"route"`
}

type SocketAddress struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
