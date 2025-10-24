# qr-pastebin

## Requirements

- Authentication (Implemented - user logging in)
- Authorization (Implemented - only owners see their shares and only admins can delete shares)
- Load balancing (TODO: NGINX)
- DDoS prevention (TODO: Cloudflare)
- HTTPS (TODO: Cloudflare)
- Secure password input (Implemented - password input is hidden)
- Secure password storing (Implemented - passwords are being hashed in DB)
- Roles (Implemented - admins can delete any share)
- Access control (Implemented - users can access some shares using only password)

## Plan

- [x] Deploy SvelteKit frontend to Cloudflare
- [x] Create a dockerfile for Go backend
- [ ] Create Docker Compose for two backend instances
- [ ] Add Nginx to docker compose
- [ ] Host on DigitalOcean the docker compose file
- [ ] Point cloudflare's api IP to backend VM (DigitalOcean)
- [ ] Configure frontend and backend communication, env variables

## Frontend

SvelteKit web application is hosted on CloudFlare and can be reached via [qr-pastebin.pages.dev](https://qr-pastebin.pages.dev/).

After pushing new changes, CloudFlare automatically redeploys the app.

## Backend

The backend consists of a single docker container inside which the following services run:

- PostgreSQL database
- Go API (instance 1)

To launch, go into `qr-pastebin/api` and run:

```bash
docker compose up -d
```

To stop, run:

```bash
docker compose down
```

To also remove the persistent database data (for example, after updating the schema), include the `-v` flag like so:

```bash
docker compose down -v
```
