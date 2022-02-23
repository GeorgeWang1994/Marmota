package db

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"marmota/pkg/common/model"
	"time"

	"github.com/astaxie/beego/orm"
)

const timeLayout = "2006-01-02 15:04:05"

// 插入事件到数据库中
func insertEvent(q orm.Ormer, eve *model.Event) (res sql.Result, err error) {
	var status int
	if status = 0; eve.Status == "OK" {
		status = 1
	}
	sqltemplete := `INSERT INTO events (
		event_caseId,
		step,
		cond,
		status,
		timestamp
	) VALUES(?,?,?,?,?)`
	res, err = q.Raw(
		sqltemplete,
		eve.Id,
		eve.CurrentStep,
		fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
		status,
		time.Unix(eve.EventTime, 0).Format(timeLayout),
	).Exec()

	if err != nil {
		log.Errorf("insert event to db fail, error:%v", err)
	} else {
		lastid, _ := res.LastInsertId()
		log.Debug("insert event to db succ, last_insert_id:", lastid)
	}
	return
}

func InsertEvent(eve *model.Event) {
	q := orm.NewOrm()
	//insert case
	_, err := insertEvent(q, eve)
	if err != nil {
		log.Error(err)
	}
}
