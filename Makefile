DISTFILE=pkgtray

clean:
	@echo "Cleaning..."
	@rm -rf dist
	@echo "Done"

build:
	@echo "Building..."
	@mkdir -p dist
	@go build -o dist/${DISTFILE}
	@echo "Done"

install:
	@echo "Installing..."
	@install -DZs dist/${DISTFILE} -m 755 -t ${DESTDIR}/usr/bin
	@echo "Done"

uninstall:
	@echo "Uninstalling..."
	@rm -rf ${DESTDIR}/usr/bin/${DISTFILE}
	@echo "Done"

