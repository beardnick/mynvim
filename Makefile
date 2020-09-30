
all:build manifest

build:
	go build -v -o bin/win-container .

manifest:
	$(CURDIR)/bin/win-container -manifest win-container -location $(CURDIR)/plugin/plugin.vim
