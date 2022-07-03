package manager

import (
	"github.com/DipandaAser/tg-bot-storage/pkg/manager"
)

var defaultManager *manager.Manager

//InitManager initializes the manager. This function must be called before GetDefaultManager
func InitManager(tokens ...string) error {
	newManager, err := manager.NewManager(tokens...)
	if err != nil {
		return err
	}

	defaultManager = newManager
	return nil
}

//GetDefaultManager returns a manager. This function must be called after InitManager
func GetDefaultManager() *manager.Manager {
	return defaultManager
}
