package analyze

type StaticResources struct {
	Name    string
	Address SocketAddress
	Routes   []Route
	Clusters []Cluster
}

type Route struct {
	Prefix string
	Cluster string
}

type Cluster struct {
	Name string
	Host string
	Port int
}

type SocketAddress struct {
	Host string
	Port int
}
