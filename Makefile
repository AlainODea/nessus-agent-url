NAME        = nessus-agent-url
VERSION     = 0.0.1
PLATFORMS 	:= darwin/amd64 darwin/arm64 linux/386 linux/amd64 linux/arm64 windows/386 windows/amd64
package 	= main.go
binary 		= build/${NAME}

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

compile: clean ## Build binaries for all platforms defined in $PLATFORMS
ifndef NAME
	$(error NAME is undefined)
endif
ifndef VERSION
	$(error VERSION is undefined)
endif
	@ echo "Building binaries for version $(VERSION):"
	@ for platform in $(PLATFORMS); do 										  \
		platform_split=($${platform//\// }); 								  \
		GOOS=$${platform_split[0]}; 										  \
		GOARCH=$${platform_split[1]}; 										  \
		output_name=$(binary)'_'$(VERSION)'_'$$GOOS'-'$$GOARCH;				  \
		if [ $$GOOS = "windows" ]; then 									  \
			output_name+='.exe'; 											  \
		fi;	 																  \
		echo "- Build for $$platform -> $$output_name";						  \
		env GOOS=$$GOOS GOARCH=$$GOARCH go build -ldflags "-X main.Version=$${VERSION}" -o $$output_name $(package); \
	done
	@ echo "Done"

clean: ## Cleanup build artifacts and remove build/
	rm -rf build/
