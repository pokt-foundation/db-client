make: gen_client gen_reader gen_writer

gen_client:
	mockery --name=IDBClient --filename=mock_client.go --recursive --inpackage
gen_reader:
	mockery --name=IDBReader --filename=mock_reader.go --recursive --inpackage
gen_writer:
	mockery --name=IDBWriter --filename=mock_writer.go --recursive --inpackage


# These targets spin up and shut down the E2E test env in docker.
test_env_up:
	docker-compose -f ./testdata/docker-compose.test.yml up -d --remove-orphans --build
	@echo "⏳ Waiting for test DB to be ready ..."
	until pg_isready -h localhost -p 5432 -U postgres -d postgres >/dev/null 2>&1; do sleep 0.01; done
	@echo "🚀 Test environment is up ..."
test_env_down:
	docker-compose -f ./testdata/docker-compose.test.yml down --remove-orphans -v
	@echo "✅ Test environment is down."


run_client_tests:
	-go test ./... -run E2E_PocketHTTPDBTestSuite -count=1;

# This target runs all tests, which includes spinning up the Docker test env.
test: test_env_up run_client_tests test_env_down

# This target install pre-commit to the repo and should be run only once, after cloning the repo for the first time.
init-pre-commit:
	wget https://github.com/pre-commit/pre-commit/releases/download/v2.20.0/pre-commit-2.20.0.pyz;
	python3 pre-commit-2.20.0.pyz install;
	python3 pre-commit-2.20.0.pyz autoupdate;
	go install golang.org/x/tools/cmd/goimports@v0.6.0;
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.0;
	go install -v github.com/go-critic/go-critic/cmd/gocritic@v0.6.5;
	python3 pre-commit-2.20.0.pyz run --all-files;
	rm pre-commit-2.20.0.pyz;
