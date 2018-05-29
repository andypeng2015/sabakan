package main

func newConfig() *config {
	return &config{
		ListenHTTP:  "0.0.0.0:10080",
		URLPort:     "10080",
		EtcdServers: []string{"http://localhost:2379"},
		EtcdPrefix:  "/sabakan",
		EtcdTimeout: "2s",
		DHCPBind:    "0.0.0.0:10067",
		IPXEPath:    "/usr/lib/ipxe/ipxe.efi",
	}
}

type config struct {
	ListenHTTP  string   `yaml:"http"`
	URLPort     string   `yaml:"url-port"`
	EtcdServers []string `yaml:"etcd-servers"`
	EtcdPrefix  string   `yaml:"etcd-prefix"`
	EtcdTimeout string   `yaml:"etcd-timeout"`
	DHCPBind    string   `yaml:"dhcp-bind"`
	IPXEPath    string   `yaml:"ipxe-efi-path"`
}