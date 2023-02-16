make: gen_client gen_reader gen_writer

gen_client:
	mockery --name=IDBClient --filename=mock_client.go --recursive --inpackage
gen_reader:
	mockery --name=IDBReader --filename=mock_reader.go --recursive --inpackage
gen_writer:
	mockery --name=IDBWriter --filename=mock_writer.go --recursive --inpackage

test: test_env_up run_client_tests test_env_down
test_env_up:
	docker-compose -f ./testdata/docker-compose.test.yml up -d --remove-orphans --build;
	sleep 2;
test_env_down:
	docker-compose -f ./testdata/docker-compose.test.yml down --remove-orphans -v
run_client_tests:
	-go test ./... -run Test_RunDBClientTestSuite -count=1 -v;

init-pre-commit:
	wget https://github.com/pre-commit/pre-commit/releases/download/v2.20.0/pre-commit-2.20.0.pyz;
	python3 pre-commit-2.20.0.pyz install;
	python3 pre-commit-2.20.0.pyz autoupdate;
	go install golang.org/x/tools/cmd/goimports@latest;
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest;
	go install -v github.com/go-critic/go-critic/cmd/gocritic@latest;
	python3 pre-commit-2.20.0.pyz run --all-files;
	rm pre-commit-2.20.0.pyz;
