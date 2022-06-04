package service

import (
	"github.com/jackdallas/premiumizearr/internal/config"
)

//Service interface
type Service interface {
	New() (*config.Config, error)
	Start() error
	Stop() error
}
