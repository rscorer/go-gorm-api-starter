module api

go 1.15

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/containerd/containerd v1.5.8 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/joho/godotenv v1.4.0
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/sirupsen/logrus v1.8.1
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/sys v0.0.0-20211124211545-fe61309f8881 // indirect
)

replace github.com/containerd/containerd => github.com/containerd/containerd v1.5.8

replace github.com/opencontainers/image-spec => github.com/opencontainers/image-spec v1.0.2
