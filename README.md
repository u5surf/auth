## auth

This repository holds the authentication service for [moov.io](https://github.com/moov-io). If you find a problem (security or otherwise), please contact us at [`security@moov.io`](mailto:security@moov.io).

The auth project supports various auth methods:
- REST authentication and user signup
- OAuth2 exchange (linked to an authenticated user)

### Getting Started

You can download [our docker image `moov/auth`](https://hub.docker.com/r/moov/auth/) from Docker Hub or use this repository. No configuration is required to serve on `localhost:8080`.

Metrics are served at `localhost:9090/metrics` in prometheus format.

### Configuration

The follow are environment variables can be configured:

**Required**
- `DOMAIN`: Domain to set on cookies.
**Optional**
- `OAUTH2_CLIENTS_DB_PATH`: Filepath to our oauth2 clients database.
- `OAUTH2_TOKENS_DB_PATH`: Filepath to our oauth2 tokens database.
- `SQLITE_DB_PATH`: Filepath to our sqlite database
- `TLS_CERT` and `TLS_KEY`: Filepaths to TLS certificate and keyfile (in PEM encoding)

### Endpoints

- GET /ping

- POST   /users/create
- GET    /users/login
- POST   /users/login
- DELETE /users/login

- GET    /oauth2/authorize
- POST   /oauth2/token  (NOTE: GET is supported with env var: ...)
- POST   /oauth2/token/create

### metrics

<dl>
    <dt>auth_successes</dt><dd>Count of successful authorizations</dd>
    <dt>auth_failures</dt><dd>Count of failed authorizations</dd>
    <dt>auth_inactivations</dt><dd>Count of inactivated auths (i.e. user logout)</dd>
    <dt>http_errors</dt><dd>Count of how many 5xx errors we send out</dd>
    <dt>oauth2_client_generations</dt><dd>Count of auth tokens created</dd>
    <dt>oauth2_token_generations</dt><dd>Count of auth tokens created</dd>
    <dt>sqlite_connections</dt><dd>How many sqlite connections and what status they're in.</dd>
</dl>
