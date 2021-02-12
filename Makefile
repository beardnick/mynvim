
all:build manifest

build:
	go build -v -o bin/mynvim

manifest:
	$(CURDIR)/bin/mynvim -manifest mynvim  -location $(CURDIR)/plugin/mynvim.vim

install:all
	cp $(CURDIR)/plugin/* $(HOME)/.local/share/nvim/site/autoload/
	mkdir -p $(HOME)/.local/share/nvim/site/bin
	cp $(CURDIR)/bin/mynvim $(HOME)/.local/share/nvim/site/bin/
