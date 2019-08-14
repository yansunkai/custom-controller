module github.com/yansunkai/custom-controller

go 1.12

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/apiserver v0.0.0
	k8s.io/client-go v0.0.0
)

replace (
	github.com/golang/glog => github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20181025213731-e84da0312774
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190312203227-4b39c73a6495
	golang.org/x/image => github.com/golang/image v0.0.0-20190227222117-0694c2d4d067
	golang.org/x/lint => github.com/golang/lint v0.0.0-20181217174547-8f45f776aaf1
	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20190312151609-d3739f865fa6
	golang.org/x/net => github.com/golang/net v0.0.0-20190206173232-65e2d4e15006
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190402181905-9f3314589c9a
	golang.org/x/sync => github.com/golang/sync v0.0.0-20181108010431-42b317875d0f
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190209173611-3b5209105503
	golang.org/x/text => github.com/golang/text v0.3.1-0.20181227161524-e6919f6577db
	golang.org/x/time => github.com/golang/time v0.0.0-20161028155119-f51c12702a4d
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190313210603-aa82965741a9
	k8s.io/api => ./k8s.io/api
	k8s.io/apimachinery => ./k8s.io/apimachinery
	k8s.io/apiserver => ./k8s.io/apiserver
	k8s.io/client-go => ./k8s.io/client-go
	k8s.io/component-base => ./k8s.io/component-base
)
