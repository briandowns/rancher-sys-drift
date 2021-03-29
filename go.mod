module github.com/briandowns/rancher-sys-drift

go 1.13

replace k8s.io/client-go => k8s.io/client-go v0.18.0

require (
	github.com/elastic/go-sysinfo v1.6.0
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/rancher/lasso v0.0.0-20200905045615-7fcb07d6a20b
	github.com/rancher/wrangler v0.7.2
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.2
	golang.org/x/tools v0.0.0-20191017205301-920acffc3e65 // indirect
	google.golang.org/appengine v1.6.1 // indirect
	k8s.io/apiextensions-apiserver v0.18.0
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
)
