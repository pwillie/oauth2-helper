# Oauth2 command line helper

This is a very simple command line utility to aid in issuing an oauth2 access token.

Will open a web browser for user authentication whilst listening for the callback at localhost:5556/auth/callback.  On callback the access token will be written to stdout.

Example:
```
go run main.go --issuer https://example.idp.com/ --client-id exampleClientID --scope example/scope1 --scope example/scope2
```