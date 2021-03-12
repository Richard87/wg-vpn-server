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
    - [X] Configure with clients from client.db (watch for changes) / generate wg0.conf
    - [X] Generate Private key if missing
    - [X] Configure with Private key
    - [X] Find next available IP
    - [ ] Validate IP on create
    - [ ] Create wgX device if missing (depending on platform)
        - Init device if not up
          ```shell
          ip -4 address add 10.0.0.1/24 dev wg0
          ip link set mtu 1432 up dev wg0
          sysctl -q net.ipv4.conf.all.src_valid_mark=1
          sysctl -q net.ipv4.ip_forward=1
          ```
        - Redo available ip-addresses, claim one for WireGuard Servre
    - [ ] Use a userspace WireGuard implementation: https://github.com/cloudflare/boringtun
 

## Ro run:
Install boringtun: `cargo install boringtun`
Start boringtun: `WG_SUDO=1 sudo ~/.cargo/bin/boringtun -v debug -f wg0`

Requires GO 1.16
Install wg-vpn-server: `go get github.com/richard87/wg-vpn-server`
Start `sudo wg-vpn-server -wg-endpoint vpn.example.com`
