# WG VPN Server
 - Easily add clients
 - Secure

## TODO:

 - API server
   - List clients
   - Create client
   - Modify client name
   - Delete client
 - Client database
 - WireGuard
   - Configure with clients from client.db (watch for changes)
   - Configure with Private key
 - Init
   - API Users: Use a commandline hashed password, or create one
     --user admin:$argon2i$v=19$m=16,t=2,p=1$S1p3Z0FTQTViZkh0MURTVA$jxPFAzQ3kSrbEPSibCQIrg
     (If no specified, generate a Admin user with random password)
   - Clients.db: Use a commandline specified database file, or create a default one
     ./var/clients.db
   - WG Private key: use ./var/wg.private
     (use -create-private-key-if-missing command line argument to create it if missing)
   - server:
     use ./var/server.key (or specified by --http-server-key) and ./var/server/crt (or specified by --http-server-crt)
     - use port 8000 (or specified by --http-port)
 