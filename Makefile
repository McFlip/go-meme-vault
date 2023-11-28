bdd-run : test.exe
ifdef spec
	cd bdd_tests; pnpx cypress run --spec cypress/e2e/${spec}/${spec}.feature
else
	cd bdd_tests; pnpx cypress run
endif

test.exe :
	go build -o test.exe

clean :
	rm *.exe