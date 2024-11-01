# By default, all binary and archive files will be stored here.
InstallPath=$(shell pwd)

# Project dir, do not change this path(unless you know what you're doing).
projectPath=$(shell pwd)

# 一般是项目名称或包名
TARGET=projectName

# Compilation related information.
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date +"%Y-%m-%d")
LDFLAGS="-X main.CommitId=$(COMMIT) -X main.Built=$(DATE)"
GOBUILD=$(GO) build -ldflags $(LDFLAGS) -v
# 使用 := 意味着一旦初始化后，之后不能被覆盖（不能重新被赋值）
# 而使用 = 则意味着像编程语言中的变量一样，可以多次被重新赋值
# 但无论是哪种方式声明的变量，都能通过调用make命令时所传递的命令行参数覆盖
RPMROOT := ~/rpmbuild
VERSION := 0.0.1
REALEASE := 1
ARCH := $(shell uname -m)
GOOS:=linux
GOARCH:=amd64

# Setup go env
export PATH := $(shell go env GOPATH)/bin:$(PATH)
export GO111MODULE := on

clean:
	$(GO) clean
	rm -f $(InstallPath)/consumer $(InstallPath)/producer consumer.tar.gz producer.tar.gz

# rpm构建
rpm:
    @rm -rf ${TARGET} *.rpm
    go build -ldflags ${LDFLAGS} -o ${TARGET}
    @mkdir -p ${RPMROOT}/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS}
    @rm -rf ${RPMROOT}/SOURCES/${TARGET}*
    @mkdir -p ${TARGET}-${VERSION}
    @cp -f ${TARGET} ${TARGET}-${VERSION}
    # 配置文件和service文件按实际情况进行修改
    @cp -f config/${TARGET}.conf ${TARGET}-${VERSION}
    @cp -f config/${TARGET}.service ${TARGET}-${VERSION}
    @tar cvzf ${TARGET}-${VERSION}.tar.gz ${TARGET}-${VERSION}
    @rm -rf ${TARGET}-${VERSION}
    @mv ${TARGET}-${VERSION}.tar.gz ${RPMROOT}/SOURCES/
    @cp -f config/${TARGET}.spec ${RPMROOT}/SOURCES/
    @sed -i "s/ExVERSION/${VERSION}/g" ${RPMROOT}/SOURCES/${TARGET}.spec
    @sed -i "s/REALEASE/${REALEASE}/g" ${RPMROOT}/SOURCES/${TARGET}.spec
    rpmbuild -bb ${RPMROOT}/SOURCES/${TARGET}.spec
    @mv ${RPMROOT}/RPMS/${ARCH}/${TARGET}-${VERSION}-* .
