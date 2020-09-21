package uuidgen

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	epoch int64 = 1526285084373

	numWorkerBits = 16

	numSequenceBits = 12

	//MaxWorkID 最大机器ID
	MaxWorkID = -1 ^ (-1 << numWorkerBits)

	//MaxSequence 最大序列号
	MaxSequence = -1 ^ (-1 << numSequenceBits)
)

//SnowFlake SnowFlake
type SnowFlake struct {
	lastTimestamp uint64
	sequence      uint16
	workerID      uint16
	lock          sync.Mutex
}

func (sf *SnowFlake) pack() uint64 {
	uuid := (sf.lastTimestamp << (numWorkerBits + numSequenceBits)) | (uint64(sf.workerID) << numSequenceBits) | (uint64(sf.sequence))
	return uuid
}

//New returns a new snowflake node that can be used to generate snowflake
func New(workerID uint16) (*SnowFlake, error) {
	if workerID < 0 || workerID > MaxWorkID {
		return nil, errors.New("invalid worker Id")
	}
	return &SnowFlake{workerID: workerID}, nil
}

//GetK8sWorkID K8S内部WorkID
// 请确保K8S的CIDR为/16，否则可能会出现重复ID
func GetK8sWorkID() (uint16, error) {
	netInterface, err := net.InterfaceByName("eth0")
	if err != nil {
		return 0, err
	}
	if netInterface.Flags&net.FlagUp == 0 {
		return 0, fmt.Errorf("eth0 interface is down") // interface down
	}
	addrs, err := netInterface.Addrs()
	if err != nil {
		return 0, err
	}

	var workID uint16
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip == nil || ip.IsLoopback() {
			continue
		}

		ip = ip.To4()
		if ip == nil {
			continue // not an ipv4 address
		}

		workID = uint16(uint(ip[2])<<8 | uint(ip[3]))
		break
	}
	return workID, nil
}

//Generate  Next creates and returns a unique snowflake ID
func (sf *SnowFlake) Generate() (uint64, error) {
	sf.lock.Lock()
	defer sf.lock.Unlock()

	ts := timestamp()
	if ts == sf.lastTimestamp {
		sf.sequence = (sf.sequence + 1) & MaxSequence
		if sf.sequence == 0 {
			ts = sf.waitNextMilli(ts)
		}
	} else {
		sf.sequence = 0
	}

	if ts < sf.lastTimestamp {
		return 0, errors.New("invalid system clock")
	}

	sf.lastTimestamp = ts
	return sf.pack(), nil
}

// waitNextMilli if that microsecond is full
// wait for the next microsecond
func (sf *SnowFlake) waitNextMilli(ts uint64) uint64 {
	for ts == sf.lastTimestamp {
		time.Sleep(100 * time.Microsecond)
		ts = timestamp()
	}
	return ts
}

// timestamp
func timestamp() uint64 {
	return uint64(time.Now().UnixNano()/int64(1000000) - epoch)
}

//NetIP 获取指定网卡的IP地址
func eth0IP(name string) uint16 {

}
