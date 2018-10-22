GO_CMD=go
REPO_PATH=crawler.club/et
CMD_PATH=${REPO_PATH}/cmd/et
GIT_SHA=`git rev-parse --short HEAD || echo "GitNotFound"`
GO_LDFLAGS=-ldflags "-X ${REPO_PATH}/version.GitSHA=${GIT_SHA}"
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get

BUILD_DIR=build

all: linux darwin windows

linux:
	GOOS=linux $(GO_BUILD) $(GO_LDFLAGS) -o "$(BUILD_DIR)/linux/et" "${CMD_PATH}"
darwin:
	GOOS=darwin $(GO_BUILD) $(GO_LDFLAGS) -o "$(BUILD_DIR)/darwin/et" "${CMD_PATH}"
windows:
	GOOS=windows $(GO_BUILD) $(GO_LDFLAGS) -o "$(BUILD_DIR)/windows/et.exe" "${CMD_PATH}"

clean:
	$(GO_CLEAN)
	rm -fr $(BUILD_DIR)
