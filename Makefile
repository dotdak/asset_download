.PHONY: 

windows_64:
	GOOS=windows GOARCH=amd64 go build -o out/asset_download_windows_64.exe

windows_32: 
	GOOS=windows GOARCH=386 go build -o out/asset_download_windows_32.exe

linux_64: 
	GOOS=linux GOARCH=amd64 go build -o out/asset_download_linux_64

linux_32:
	GOOS=linux GOARCH=386 go build -o out/asset_download_linux_32

darwin:
	GOOS=darwin GOARCH=amd64 go build -o out/asset_download_darwin

arm: 
	GOOS=linux GOARCH=arm64 go build -o out/asset_download_arm
