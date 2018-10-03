ifeq ($(strip $(GOPATH)),)
$(error "make sure you defined the GOPATH environment variable")
endif

GOBUILD := go build
GOCLEAN := go clean
GOTEST  := go test
GOGET   := go get
TARG    := $(notdir $(CURDIR))
TARGET  := $(TARG).exe
TARGET_UNIX := $(TARG)_unix
BUILD =

ifeq ($(OS),Windows_NT)
	BUILD += $(GOBUILD) -o $(TARGET) -v
else
	BUILD += CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(TARGET_UNIX) -v
endif

all: build
build:
	$(GOGET) -u github.com/bwmarrin/discordgo
	$(BUILD)
