pkgname=pgsync
pkgver=1.0.0
pkgrel=1
pkgdesc="A beautiful, interactive command-line tool for migrating PostgreSQL databases"
arch=('x86_64' 'aarch64')
url="https://github.com/madss-bin/pgsync"
license=('MIT')
depends=('postgresql-libs')
makedepends=('go')
source=()
sha256sums=()

build() {
  cd "$startdir"
  export CGO_CPPFLAGS="${CPPFLAGS}"
  export CGO_CFLAGS="${CFLAGS}"
  export CGO_CXXFLAGS="${CXXFLAGS}"
  export CGO_LDFLAGS="${LDFLAGS}"
  export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
  
  go build -o pgsync .
}

package() {
  cd "$startdir"
  install -Dm755 pgsync "$pkgdir/usr/bin/pgsync"
  install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md"
}
