# Bajo

Bajo is a URL shortening service written in Go.

## Features

The fundamental features of this project are as follows:

- [x] Allow users to submit a URL of arbitrary length and receive a shortened URL.
- [x] Maintain a persistent, database backend mapping between submitted URLs and shortened URLs.
- [x] Redirect users to the submitted URL when the shortened URL is clicked.

With the fundamentals in place, some nice-to-have features are listed below:

- [x] Give users the ability to submit a custom shortened URL key.
- [ ] Provide click statistics for shortened URLs.

## Tests

Unit tests can be run within the container by executing the following commands:

```
docker-compose build
docker-compose run bajo ginkgo
```
