# This is how we want to name the binary output
BINARY=comicknife

BUILDFOLDER = build/bin


default: fmt clean
	go build -o ${BUILDFOLDER}/${BINARY} cmd/*.go
	@echo "Your binary is ready. Check "${BUILDFOLDER}/${BINARY}

fmt:
	@echo "Run gofmt"
	@echo "Run goimports"
	bash fmt.sh

clean:
	rm -rf frontend/dist/ && rm -rf build/

all:
	cd cmd; gox -output="../build/${BINARY}_{{.OS}}_{{.Arch}}"
