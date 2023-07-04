package dao

import "gorm.io/gorm"

type Entity interface {
	getTableName() string
}

func GetMySQLDb(e Entity) *gorm.DB {
	return mySqlDB.Table(e.getTableName())
}

func Insert(e Entity) {
	GetMySQLDb(e).Save(e)
}

func SelectById(e Entity, id int64) {
	GetMySQLDb(e).Select("*").Where("id=?", id).Take(&e)
}

func SelectByName(e Entity, name string) {
	GetMySQLDb(e).Where("name=?", name).Take(&e)
}

func SelectRandom(e Entity) {
	GetMySQLDb(e).Select("*").Order("RAND()").Take(&e)
}
