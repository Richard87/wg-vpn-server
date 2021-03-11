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
    - [X] Config
        - [X] Get endpoint
        - [X] Get next available IP
        - [X] Use real endpoint
- [ ] UI
    - [X] Create clients
    - [X] List clients
    - [X] Generate QR Code
    - [X] Generate Config
    - [X] Delete Client
    - [ ] Confirm deletion
    - [ ] Authentication
- [ ] WireGuard
    - [ ] Configure with clients from client.db (watch for changes) / generate wg0.conf
    - [X] Generate Private key if missing
    - [X] Configure with Private key
    - [X] Find next available IP
    - [ ] Validate IP on create
    - [ ] Use a userspace Wireguard implementation: https://github.com/cloudflare/boringtun
 