# Golang OAuth Example

## Authorization Code Grant

Prerequisites:

1. Register your application on Github: https://github.com/settings/applications/new:
    * Authorization callback URL: http://localhost:8080/oauth/redirect

## Client Credentials Grant

```
$ curl -X GET http://localhost:8081/protected
invalid access token
```

```
$ curl -X GET http://localhost:8081/credentials | jq
{
  "client_id": "ef1e4d63",
  "client_secret": "0147a508"
}
```

```
$ curl -X GET -H "Accept: application/json" "http://localhost:8081/token?grant_type=client_credentials&scope=all&client_id=ef1e4d63&client_secret=0147a508" | jq
```

```
$ curl -X GET http://localhost:8081/protected?access_token=GAMJVUZZOW-FZWEHUIBV4Q
$ curl -X GET -H "Authorization: Bearer GAMJVUZZOW-FZWEHUIBV4Q" http://localhost:8081/protected
Hi, I'm a protected data
```

