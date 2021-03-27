module github.com/Richard87/wg-vpn-server

go 1.16

require (
	github.com/alexedwards/argon2id v0.0.0-20201228115903-cf543ebc1f7b
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gofiber/fiber/v2 v2.6.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/mdlayher/genetlink v1.0.0
	github.com/mdlayher/netlink v1.1.0
	go.etcd.io/bbolt v1.3.5
	golang.zx2c4.com/wireguard v0.0.20200121
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20200609130330-bd2cb7843e1b
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.6
)
