package database

import (
	"net"
	"strings"

	"github.com/stavros-k/go-dmarc-analyzer/internal/types"
)

type AddressModel struct {
	IP       string `json:"ip" gorm:"primaryKey"`
	Hostname string `json:"hostname"`
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
