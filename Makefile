.PHONY: clean

bdd-run : clean test.exe
ifdef spec
	./test.exe & export SERVER_PID=$$!; cd bdd_tests; pnpx cypress run --spec cypress/e2e/${spec}/${spec}.feature; kill $${SERVER_PID}
else
	./test.exe & export SERVER_PID=$$!; cd bdd_tests; pnpx cypress run; kill $${SERVER_PID}
endif

bdd-open : clean test.exe
	./test.exe & export SERVER_PID=$$!; cd bdd_tests; pnpx cypress open; kill $${SERVER_PID}

test.exe : templ
	go build -o test.exe

dev : clean templ
	air

build : clean templ
	env GOOS=linux GOARCH=amd64 go build -o build/go-meme-vault
	env GOOS=windows GOARCH=amd64 go build -o build/go-meme-vault.exe
	# TODO: bundle static assets and folder struct
	# vendor tailwind, htmx, & alpine js
	rm dist/public/img/full/.gitkeep
	rm dist/public/img/tn/.gitkeep
	ln -s build/go-meme-vault dist/go-meme-vault
	tar czf release/linux.tar.gz dist
	unlink dist/go-meme-vault
	ln build/go-meme-vault.exe dist/go-meme-vault.exe
	zip -r release/windows.zip dist
	unlink dist/go-meme-vault.exe
	touch dist/public/img/full/.gitkeep
	touch dist/public/img/tn/.gitkeep

templ :
	templ generate

clean :
	-rm *.exe dev.db public/img/tn/*.jpg
	-rm dist/go-meme-vault
	-rm dist/go-meme-vault.exe
	-rm release/*
