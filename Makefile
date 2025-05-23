APP_NAME = joplin-fuse
VERSION := $(shell git describe --dirty --tags | sed 's/^v//' )
ARCH := amd64
BUILD_DIR := build
PKG_DIR := $(BUILD_DIR)/$(APP_NAME)_$(VERSION)
BIN_PATH := $(PKG_DIR)/usr/bin
DEB_FILE := $(APP_NAME)_$(VERSION)_$(ARCH).deb
AUTHOR := cretingame <you@example.com>
DESCRIPTION := Joplin Fuse is a Go-based tool that mounts your Joplin notes \
	       into a filesystem using FUSE (Filesystem in Userspace). \
	       This allows you to browse, read, and interact with your Joplin \
	       notes as if they were regular files on your system.

EXEC = $(APP_NAME)
MOUNT_POINT = /run/user/$(shell id -u)/$(EXEC)
SOURCES = $(shell find . -name "*.go" -not -path "./vendor/*")

build: $(EXEC)

package: $(BUILD_DIR)/$(DEB_FILE)

$(EXEC): $(SOURCES)
	go build

$(MOUNT_POINT):
	install --group=$(shell id -g) --owner=$(shell id -g) --directory $(MOUNT_POINT)

run: $(MOUNT_POINT) $(CURDIR)/$(EXEC)
	$(CURDIR)/$(EXEC) $(MOUNT_POINT)


$(BUILD_DIR)/$(DEB_FILE): build
	@echo "Packaging .deb file..."
	mkdir -p $(BIN_PATH)
	mv $(APP_NAME) $(BIN_PATH)/

	mkdir -p $(PKG_DIR)/DEBIAN
	echo "Package: $(APP_NAME)" > $(PKG_DIR)/DEBIAN/control
	echo "Version: $(VERSION)" >> $(PKG_DIR)/DEBIAN/control
	echo "Section: utils" >> $(PKG_DIR)/DEBIAN/control
	echo "Priority: optional" >> $(PKG_DIR)/DEBIAN/control
	echo "Architecture: $(ARCH)" >> $(PKG_DIR)/DEBIAN/control
	echo "Maintainer: $(AUTHOR)" >> $(PKG_DIR)/DEBIAN/control
	echo "Description: $(DESCRIPTION)" >> $(PKG_DIR)/DEBIAN/control

	dpkg-deb --build $(PKG_DIR)
	mv $(PKG_DIR).deb $(BUILD_DIR)/$(DEB_FILE)

clean:
	rm -f $(EXEC)
	-rmdir $(MOUNT_POINT)
	rm -rf $(BUILD_DIR)


.PHONY: clean run build package
