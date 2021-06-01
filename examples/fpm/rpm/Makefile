.PHONY: all remove add fpm
remove:
	aliasd remove -n fpm

add:
	aliasd add -f ./resource.yaml

fpm:
	${HOME}/.aliasd/bin/fpm -s dir -t rpm -n "aliasd" -v 0.1 -p $(shell pwd) -C $(shell pwd)/rpmbuild ./ 

all: remove add fpm 