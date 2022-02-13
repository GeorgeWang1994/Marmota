package db

import (
	"fmt"
	"log"
	"marmota/pkg/common/model"
	"time"
)

// QueryHosts 查询所有的主机列表
func QueryHosts() (map[string]int, error) {
	m := make(map[string]int)

	sql := "select id, hostname from host"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return m, err
	}

	defer rows.Close()
	for rows.Next() {
		var (
			id       int
			hostname string
		)

		err = rows.Scan(&id, &hostname)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		m[hostname] = id
	}

	return m, nil
}

// QueryMonitoredHosts 查询所有监控的主机
func QueryMonitoredHosts() (map[int]*model.Host, error) {
	hosts := make(map[int]*model.Host)
	now := time.Now().Unix()
	sql := fmt.Sprintf("select id, hostname from host where maintain_begin > %d or maintain_end < %d", now, now)
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return hosts, err
	}

	defer rows.Close()
	for rows.Next() {
		t := model.Host{}
		err = rows.Scan(&t.Id, &t.Name)
		if err != nil {
			log.Println("WARN:", err)
			continue
		}
		hosts[t.Id] = &t
	}

	return hosts, nil
}
