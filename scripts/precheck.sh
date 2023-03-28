rm -rf log_osv.txt log_test.txt log_build.txt log_lint.txt vendor
go mod tidy && go mod vendor ; \
osv-scanner . >> log_osv.txt ; \
go test -v ./... >> log_test.txt ; \
go build -v -modcacherw ./... >> log_build.txt ; \
golangci-lint run --out-format=github-actions --tests=false --timeout=10m >> log_lint.txt