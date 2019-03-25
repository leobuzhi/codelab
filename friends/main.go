package main

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/leobuzhi/codelab/friends/model"
)

var (
	db *gorm.DB
)

// NOTE(leobuzhi): In this design, we ensure userID1 < userID2 in like_person
// 1. if there appear reverse relation, there will appear a row lock conflict(use insert ... on duplicate)
// 2. we use ignore and bitwise or for idempotent,2 | 1 == 1 | 2 == 3 | 1 == 3 | 2 == 3
func attention(userID1, userID2 int) error {
	if userID1 == userID2 {
		return fmt.Errorf("userID1 and userID2 should not be equal, userID1: %d, userID2: %d", userID1, userID2)
	}
	var lp model.LikePerson
	tx := db.Begin()
	if userID1 < userID2 {
		// NOTE(leobuzhi): We use raw sql because gorm does not support on duplicate key update
		tx.Exec("insert into `like_person`(user_id, liker_id, relation_ship) values(?, ?, 1) on duplicate key update relation_ship=relation_ship | 1;", userID1, userID2)
		tx.Find(&lp, "user_id = ? AND liker_id = ?", userID1, userID2)
		if lp.RelationShip == 3 {
			tx.Exec("insert ignore into friend(person1, person2) values(?,?);", userID1, userID2)
		}
	} else {
		tx.Exec("insert into `like_person`(user_id, liker_id, relation_ship) values(?, ?, 2) on duplicate key update relation_ship=relation_ship | 2;", userID2, userID1)
		tx.Find(&lp, "user_id = ? AND liker_id = ?", userID2, userID1)
		if lp.RelationShip == 3 {
			tx.Exec("insert ignore into friend(person1, person2) values(?,?);", userID2, userID1)
		}
	}
	tx.Commit()
	return nil
}

func setupDB(dbname string) {
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("root:root@/%s?charset=utf8&parseTime=True&loc=Local", dbname))
	if err != nil {
		glog.Fatalf("init db faield, err: %v\n", err)
	}

	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.LikePerson{})
	db.Model(&model.LikePerson{}).AddUniqueIndex("uk_user_id_liker_id", "user_id", "liker_id")

	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.Friend{})
	db.Model(&model.Friend{}).AddUniqueIndex("uk_friend", "person1", "person2")
}

func teardown() {
	db.DropTable(&model.LikePerson{})
	db.DropTable(&model.Friend{})
	defer db.Close()
}

func main() {
	setupDB("testdb")
	fmt.Println(attention(1000, 2000))
	fmt.Println(attention(2000, 1000))
}
