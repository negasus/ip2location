package ip2location

import (
	"github.com/ip2location/ip2location-go"
)

var db *ip2location.DB

func Init(database string) error {
	var err error
	db, err = ip2location.OpenDB(database)
	if err != nil {
		return err
	}
	return nil
}

func Stop() {
	db.Close()
}

type ParseFunc func(string) (ip2location.IP2Locationrecord, error)

func Parse(ip string) (ip2location.IP2Locationrecord, error) {
	return db.Get_all(ip)
}
