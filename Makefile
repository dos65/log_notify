TARGETDIR="build"
RESDIR="res"
NAME="lognotify"

.PHONY: deps test all clean new_target install

all: deps test new_target
	@go build -o ${TARGETDIR}/${NAME}
	@cp res/* ${TARGETDIR}

deps: 
	@echo "* Installing dependencies"
	@go get -t -v .
	@echo "* ... Done"

new_target:
	@rm -rf ${TARGETDIR}
	@mkdir -p ${TARGETDIR}
	
test: deps
	@echo "* Running tests"
	@go test -v
	@echo "* ... Done"

install:
	@mkdir /opt/lognotify
	@cp build/* /opt/lognotify
	@ln -s /opt/lognotify/lognotify /usr/local/bin/lognotify

uninstall:
	@rm -rf /opt/lognotify
	@rm -f /usr/local/bin/lognotify

clean:
	@rm -rf ${TARGETDIR} 
	@find . -perm +100 -type f -delete
