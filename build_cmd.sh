go build -ldflags "-X 'main.version=1.1.0' -X 'main.buildTime=$(date -u +'%Y-%m-%d %H:%M:%S UTC')'" 
