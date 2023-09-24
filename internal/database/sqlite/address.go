package database_sqlite

import (
	"net"
	"strings"

	"github.com/stavros-k/go-dmarc-analyzer/internal/types"
)

type AddressModel struct {
	IP        string `gorm:"primaryKey"`
	Hostname  string
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdateAt  int64 `gorm:"autoUpdateTime"`
}

func (s *SqliteStorage) CreateAddress(address *types.Address) error {
	hostname, err := net.LookupAddr(address.IP)
	if err != nil {
		hostname = []string{""}
	}

	addr := &AddressModel{
		IP:       address.IP,
		Hostname: strings.Join(hostname, ","),
	}

	return s.db.Create(addr).Error
}
