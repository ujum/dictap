# Dictup - Save and learn new words

## APi

http://localhost:8080/swagger/index.html

## Usage

Add env vars for auth with Google Account:

- `DICTUP_SERVER_SECURITY_GOOGLEOAUTH2_CONFIG_CLIENTID=your client_id`
- `DICTUP_SERVER_SECURITY_GOOGLEOAUTH2_CONFIG_CLIENTSECRET=your client_secret`

```sh
$ docker-compose up
$ make migrate-up
```