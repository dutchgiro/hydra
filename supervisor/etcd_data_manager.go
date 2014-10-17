package supervisor

// import (
// 	"os"
// )

// type DataManager interface {
// 	RemoveAll() error
// }

// type EtcdDataManager struct {
// }

// func (e *EtcdDataManager) RemoveAll() error {
// 	if err := os.RemoveAll(e.EtcdService.(*etcd.Etcd).Config.DataDir); err != nil {
// 		return err
// 	}
// 	return nil
// }
