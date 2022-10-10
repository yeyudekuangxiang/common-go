init:.
	rm -rf ~/.goctl/1.4.0
	mkdir -p ~/.goctl
	ln -s ${PWD}/deploy/goctl/1.4.0 ~/.goctl/1.4.0