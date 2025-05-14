EXEC = hello-joplin
MOUNT_POINT = /run/user/$(shell id -u)/$(EXEC)

$(CURDIR)/$(EXEC):
	go build

$(MOUNT_POINT):
	install --group=$(shell id -g) --owner=$(shell id -g) --directory $(MOUNT_POINT)

run: $(MOUNT_POINT) $(CURDIR)/$(EXEC)
	$(CURDIR)/$(EXEC) $(MOUNT_POINT)

clean:
	rm -f $(CURDIR)/$(EXEC)
	rmdir $(MOUNT_POINT)


.PHONY: clean
