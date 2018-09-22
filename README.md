## auth

This repository holds the authentication service for [moov.io](https://github.com/moov-io). If you find a problem (security or otherwise), please contact us at [`security@moov.io`](mailto:security@moov.io).

### runbook

// TODO(adam)

### configuration

The follow are environment variables which

- `DOMAIN`: Required

- `OAUTH2_DB_PATH`: TODO
- `SQLITE_DB_PATH`: TODO

- `TLS_CERT` and `TLS_KEY` TODO

### routes

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
    <dt>auth_token_generations</dt><dd>Count of auth tokens created</dd>
    <dt>sqlite_connections</dt><dd>How many sqlite connections and what status they're in.</dd>
</dl>
