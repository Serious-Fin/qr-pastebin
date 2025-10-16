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

- [ ] Deploy SvelteKit frontend to Cloudflare (this gives HTTPS, get a free/cheap domain)
- [ ] Create a dockerfile for Go backend
- [ ] Create Docker Compose for two backend instances
- [ ] Add Nginx to docker compose
- [ ] Host on DigitalOcean the docker compose file
- [ ] Point cloudflare's api IP to backend VM (DigitalOcean)
- [ ] Configure frontend and backend communication, env variables
