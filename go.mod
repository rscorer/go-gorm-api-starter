module api

go 1.15

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/gorilla/mux v1.8.0
	github.com/jinzhu/gorm v1.9.16
	github.com/joho/godotenv v1.4.0
	github.com/sirupsen/logrus v1.8.1
	go.uber.org/atomic v1.9.0 // indirect
)

replace github.com/containerd/containerd => github.com/containerd/containerd v1.6.1

replace github.com/opencontainers/image-spec => github.com/opencontainers/image-spec v1.0.2
