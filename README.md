moov-io/auth
===

[![GoDoc](https://godoc.org/github.com/moov-io/auth?status.svg)](https://godoc.org/github.com/moov-io/auth)
[![Build Status](https://travis-ci.com/moov-io/auth.svg?branch=master)](https://travis-ci.com/moov-io/auth)
[![Coverage Status](https://coveralls.io/repos/github/moov-io/auth/badge.svg?branch=master)](https://coveralls.io/github/moov-io/auth?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/auth)](https://goreportcard.com/report/github.com/moov-io/auth)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moov-io/auth/master/LICENSE)

*project is under active development and is not production ready*

This repository holds the authentication service for [moov.io](https://github.com/moov-io). If you find a problem (security or otherwise), please contact us at [`security@moov.io`](mailto:security@moov.io).

The auth project supports various auth methods:
- REST authentication and user sign-up
- OAuth2 exchange (linked to an authenticated user)

Docs: [docs.moov.io](https://docs.moov.io/en/latest/) | [api docs](https://editor.swagger.io/?url=https://raw.githubusercontent.com/moov-io/auth/master/openapi.yaml)

## Project Status

This project is currently pre-production and could change without much notice, however we are looking for community feedback so please try out our code or give us feedback!

## Getting Started / Install

You can download [our docker image `moov/auth`](https://hub.docker.com/r/moov/auth/) from Docker Hub or use this repository. No configuration is required to serve on `localhost:8080`.

Metrics are served at `localhost:9090/metrics` in Prometheus format.

### Configuration

The follow are environment variables can be configured:

**Required**
- `DOMAIN`: Domain to set on cookies.

**Optional**
- `OAUTH2_CLIENTS_DB_PATH`: File path to our oauth2 clients database.
- `OAUTH2_TOKENS_DB_PATH`: File path to our oauth2 tokens database.
- `SQLITE_DB_PATH`: File path to our sqlite database
- `TLS_CERT` and `TLS_KEY`: File paths to TLS certificate and keyfile (in PEM encoding)

### Endpoints

| Method | Path | Description |
|---|---|---|
| GET | /ping | Always returns "PONG". Useful for readness check |
| POST | /users/create | Create a new user. (Signup) |
| GET | /users/login | Verify if a Cookie is valid for a user. |
| POST | /users/login | Login with an email and password.  |
| DELETE | /users/login | Invalidat a user's active cookies. |
| GET | /oauth2/authorize | Verify a Bearer OAuth2 token. |
| [GET&]POST | /oauth2/token | Create a new OAuth2 token. |
| POST | /oauth2/token/create | Create a new OAuth2 client credential set. |

### metrics

| Name | Help Text |
|---|---|
| auth_successes | Count of successful authorizations |
| auth_failures | Count of failed authorizations |
| auth_inactivations | Count of inactivated auths (i.e. user logout) |
| http_errors | Count of how many 5xx errors we send out |
| oauth2_client_generations | Count of auth tokens created |
| oauth2_token_generations | Count of auth tokens created |
| sqlite_connections | How many sqlite connections and what status they're in. |

## Getting Help

 channel | info
 ------- | -------
 [Project Documentation](https://docs.moov.io/en/latest/) | Our project documentation available online.
 Google Group [moov-users](https://groups.google.com/forum/#!forum/moov-users)| The Moov users Google group is for contributors other people contributing to the Moov project. You can join them without a google account by sending an email to [moov-users+subscribe@googlegroups.com](mailto:moov-users+subscribe@googlegroups.com). After receiving the join-request message, you can simply reply to that to confirm the subscription.
Twitter [@moov_io](https://twitter.com/moov_io)	| You can follow Moov.IO's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io) | If you are able to reproduce an problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](http://moov-io.slack.com/) | Join our slack channel to have an interactive discussion about the development of the project. [Request an invite to the slack channel](https://join.slack.com/t/moov-io/shared_invite/enQtNDE5NzIwNTYxODEwLTRkYTcyZDI5ZTlkZWRjMzlhMWVhMGZlOTZiOTk4MmM3MmRhZDY4OTJiMDVjOTE2MGEyNWYzYzY1MGMyMThiZjg)

## Supported and Tested Platforms

- 64-bit Linux (Ubuntu, Debian), macOS, and Windows

## Contributing

TODO(adam): contrib and CoC docs
Yes please! Please review our [Contributing guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) to get started!

Note: This project uses Go Modules, which requires Go 1.11 or higher, but we ship the vendor directory in our repository.

## License

Apache License 2.0 See [LICENSE](LICENSE) for details.
