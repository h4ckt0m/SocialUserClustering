# TERESA Pasteminer Makefile

# Config
BINARY="twitteruserscrapper"
VERSION=v0.1.1-beta
OS=linux
ARCH=amd64


# Build stuff
TARGET=all
BUILD_TIME=`date +%FT%T%z`
BUILD=`git rev-parse HEAD`
LDFLAGS=-ldflags="\
	-s \
	-w \
	-X main.Version=${VERSION} \
	-X main.Build=${BUILD} \
	-X main.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(TARGET):
	CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}.elf
	upx --ultra-brute -o ${BINARY} ${BINARY}.elf
	if [ -f ${BINARY}.elf ] ; then rm ${BINARY}.elf ; fi

.PHONY: noupx
noupx:
	CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build ${LDFLAGS} -o ${BINARY}
