# qr-pastebin

## Requirements

- Authentication (Implemented - user logging in)
- Authorization (Implemented - only owners see their shares and only admins can delete shares)
- Load balancing (Implemented - NGINX)
- HTTPS (Implemented - local certificates)
- Secure password input (Implemented - password input is hidden)
- Secure password storing (Implemented - passwords are being hashed in DB)
- Roles (Implemented - admins can delete any share)
- Access control (Implemented - users can access some shares using only password)

## Setup

Before running the project set up certificates for HTTPS in localhost. Example for mac:

```bash
brew install mkcert
mkcert -install
mkcert localhost 127.0.0.1 ::1
```

Copy the certificates to folders `web/certs/` and `api/certs/`

## Backend

Run the whole project by running `production.sh` script

The backend consists of a single docker container inside which the following services run:

- PostgreSQL database
- Go API (instance 1)

To launch, go into `qr-pastebin/api` and run:

```bash
docker compose up -d
docker compose up -d --build --scale api=2 (to launch two api instances)
```

To stop, run:

```bash
docker compose down
```

To also remove the persistent database data (for example, after updating the schema), include the `-v` flag like so:

```bash
docker compose down -v
```

## Exec'ing into DB from docker

Connect:
```bash
psql -U postgres -d qr_pastebin
```

List databases:
```bash
\t
```

List tables:
```bash
\dt
```

Table info:
```bash
\d users
```

Quit:
```bash
\q
```
