module api

go 1.18

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/jinzhu/gorm v1.9.16
	github.com/joho/godotenv v1.4.0
	github.com/sirupsen/logrus v1.8.1
)

require (
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
)

replace github.com/containerd/containerd => github.com/containerd/containerd v1.6.4

replace github.com/opencontainers/image-spec => github.com/opencontainers/image-spec v1.0.2
