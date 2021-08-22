package seed

import (
	"log"
	"github.com/jinzhu/gorm"
	"github.com/haikalvidya/goApiBlog/api/models"
)


var users = []models.User{
	models.User{
		Nickname: "HeiHei",
		Email:    "heihei@gmail.com",
		Password: "h31h31",
	},
	models.User{
		Nickname: "haiHai",
		Email:    "haiHai@gmail.com",
		Password: "h41H41",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "Golang API Tutorial",
		Content: "Hello world, Golang API tutorial",
	},
	models.Post{
		Title:   "Python API Web",
		Content: "Hello world, api web with flask python",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorId = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}