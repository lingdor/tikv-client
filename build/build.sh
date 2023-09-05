cd ..
go mod download
cd build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tikv ../cmd/tikv.go
tar czvf tikv-amd64-linux.tar.gz tikv
rm tikv
echo "linux amd64 build success"
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o tikv ../cmd/tikv.go
tar czvf tikv-arm64-linux.tar.gz tikv
rm tikv
echo "linux arm64 build success"

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o tikv.exe ../cmd/tikv.go
zip tikv-amd64-windows.zip tikv.exe
rm tikv.exe
echo "windows amd64 build success"
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o tikv.exe ../cmd/tikv.go
zip tikv-arm64-windows.zip tikv.exe
rm tikv.exe
echo "windows arm64 build success"