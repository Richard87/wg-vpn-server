# WG VPN Server

- Easily add clients
- Secure

## TODO:

- [X] Client database
- [X] API server
    - [X] Clients
        - [X] List
        - [X] Create
        - [X] Modify name
        - [X] Delete
    - [ ] Config
        - [X] Get endpoint
        - [ ] Get next available IP
        - [ ] Use real endpoint
- [ ] UI
    - [X] Create clients
    - [X] List clients
    - [X] Generate QR Code
    - [X] Generate Config
    - [ ] Delete Client
    - [ ] Authentication
- [ ] WireGuard
    - [ ] Configure with clients from client.db (watch for changes)
    - [ ] Generate Private key if missing
    - [ ] Configure with Private key
    - [ ] Find next available IP
    - [ ] Validate IP on create
    - [ ] Use a userspace Wireguard implementation: https://github.com/cloudflare/boringtun
 