* [ ] Does your submission pass all tests against mocks? `make test`
* [ ] Do your mocks match the actual behaviour of PowerDNS Authoritative Server 4.4, 4.5 and 4.6? `docker-compose -f docker-compose-v4.6.yml up && make test-without-mocks`
* [ ] Do your tests meet or exceed the coverage results of the master branch? `make coverage && go tool cover -html=./coverage.out`
* [ ] Does your submission pass the format guidelines? `make check-fmt`

Furthermore, the pipeline performs additional checks against your submission:

* [ ] Does your submission comply with common best practices? `golangci-lint run && staticcheck ./...`
