# Word-of-Wisdom

## Solution notes

- :trident: clean architecture (handler->service->repository)
- :book: standard Go project layout (well, more or less :blush:)
- :cd: github CI/CD + docker compose + Makefile included
- :white_check_mark: tests with mocks included

## HOWTO

- run with `make run`
- test with `make test` 
- get challenge with `curl -X GET http://localhost:9090/challenge --output -`
- verify challenge with `curl -X POST http://localhost:9090/validate \        
  -d "challenge={challenge}" \
  -d "solution={solution}" --output -`
- solve through client and get quote with `curl -X GET http://localhost:9091/solve --output -`

## Possible improvements

- Implement a repository for quotes instead of directly returning them in the service.
- Cover the solver service with integration tests.