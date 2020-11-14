
all:build manifest

build:
	go build -v -o bin/mynvim

manifest:
	$(CURDIR)/bin/mynvim -manifest mynvim  -location $(CURDIR)/plugin/plugin.vim
