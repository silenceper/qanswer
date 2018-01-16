build:
	cd cmd && go  build -o ../qanswer
build_linux_amd64:
	 cd cmd && GOOS=linux GOARCH=amd64 go build -o ../qanswer_linux_amd64
build_linux_386:
	 cd cmd && GOOS=linux GOARCH=386 go  build -o ../qanswer_linux_386
build_linux_arm:
	 cd cmd && GOOS=linux GOARCH=arm go  build -o ../qanswer_linux_arm
build_windows_amd64:
	 cd cmd && GOOS=windows GOARCH=amd64 go  build -o ../qanswer_windows_amd64.exe
build_windows_386:
	 cd cmd && GOOS=windows GOARCH=386 go  build -o ../qanswer_windows_386.exe
build_darwin_amd64:
	 cd cmd && GOOS=darwin GOARCH=amd64 go  build -o ../qanswer_darwin_amd64.exe

build_all:build_linux_amd64 build_linux_386 build_linux_arm build_windows_amd64 build_windows_386 build_darwin_amd64
