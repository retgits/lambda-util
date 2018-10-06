.PHONY: install test prep deps

#--- Variables ---
AWSREGION=
AWSBUCKET=
VALIDEMAIL=

#--- Help ---
help:
	@echo 
	@echo Makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo 

#--- Setup targets ---
prep: ## Make preparations to run the tests
	mkdir -p ./testdata/downloads
	echo "Hello World!" > ./testdata/helloworld.txt

#--- Test targets ---
test: prep ## Run all testcases
	export TESTDATADIR=`pwd`/testdata && export AWSREGION=${AWSREGION} && export AWSBUCKET=${AWSBUCKET} && export VALIDEMAIL=${VALIDEMAIL} && go test -coverprofile cover.out ./... 
	go tool cover -html=cover.out -o cover.html && open cover.html