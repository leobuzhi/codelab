package model

type LikePerson struct {
	ID      int `gorm:"AUTO_INCREMENT primary_key"`
	UserID  int `gorm:"not null"`
	LikerID int `gorm:"not null"`
	//  1 mean use_id attention like_id
	//  2 mean like_id attention use_id
	//  3 mean use_id attention like_id and like_id attention use_id
	RelationShip int `gorm:"not null"`
}

func (LikePerson) TableName() string {
	return "like_person"
}

type Friend struct {
	ID      int `gorm:"AUTO_INCREMENT primary_key"`
	Person1 int `gorm:"not null"`
	Person2 int `gorm:"not null"`
}

func (Friend) TableName() string {
	return "friend"
}
