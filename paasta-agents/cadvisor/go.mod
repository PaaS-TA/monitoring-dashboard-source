module github.com/google/cadvisor

go 1.13

require (
	cloud.google.com/go v0.54.0
	github.com/aws/aws-sdk-go v1.15.11
	github.com/blang/semver v3.5.1+incompatible
	github.com/containerd/containerd v1.5.13
	github.com/containerd/typeurl v1.0.2
	github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0
	github.com/euank/go-kmsg-parser v2.0.0+incompatible
	github.com/gogo/protobuf v1.3.2
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/karrick/godirwalk v1.7.5
	github.com/mindprince/gonvml v0.0.0-20190828220739-9ebdce4bb989
	github.com/mistifyio/go-zfs v2.1.2-0.20190413222219-f784269be439+incompatible
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/runc v1.0.2
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.10.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad
	google.golang.org/grpc v1.33.2
	k8s.io/klog/v2 v2.4.0
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920
)
