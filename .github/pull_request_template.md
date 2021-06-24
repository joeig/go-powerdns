* [ ] Does your submission pass all tests against mocks? `make test`
* [ ] Does your submission pass all tests against a real instance running PowerDNS Authoritative Server 4.1, 4.2 and 4.3? `docker-compose -f docker-compose-v4.3.yml up && make test-without-mocks`
* [ ] Do your tests meet or exceed the coverage results of the master branch? `make coverage && go tool cover -html=c.out`
* [ ] Does your submission pass the format guidelines? `make check-fmt`

Furthermore, the pipeline performs additional checks against your submission:

* [ ] Does your submission comply with common best practices? `golangci-lint run && staticcheck ./...`
