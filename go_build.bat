go build -tags static -ldflags "-s -w -H=windowsgui" && move /Y img_viewer.exe .\build\
pause