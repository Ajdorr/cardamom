.PHONY: all core web

all: core web

core:
	mkdir -p docker/core/build/
	cd ../core; make build
	mv ../core/cardamom-core docker/core/build/
	mv ../core/cardamom-core-migrate docker/core-migrate/build/
	
web:
	mkdir -p docker/web/build/
	cd ../web; yarn build
	rm -rf docker/web/build/
	mv ../web/build/ docker/web/
