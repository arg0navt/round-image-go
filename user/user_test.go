package user

import (
	"flag"
	"fmt"
	"testing"
)

var pkgdir = flag.String("pkgdir", "", "dir of package containing embedded files")

func TestCreateUser(t *testing.T) {
	fmt.Print(pkgdir)
	t.Error()
}
