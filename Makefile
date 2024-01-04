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
	go build
	# TODO: bundle static assets and folder struct
	# vendor tailwind, htmx, & alpine js

templ :
	templ generate

clean :
	-rm *.exe dev.db
