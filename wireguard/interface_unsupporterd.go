// +build !linux

package wireguard

import (
	"fmt"
	"runtime"
)

func configureInterface(networkCidr string, ifaceName string) error {
	return fmt.Errorf("unsupported on %s", runtime.GOOS)
}
