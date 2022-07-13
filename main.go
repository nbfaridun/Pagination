package main

import (
	"fmt"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Sessions struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	Browser    string    `json:"browser"`
	IpAddress  string    `json:"ip_address"`
	OS         string    `json:"os"`
	Device     string    `json:"device"`
	MerchantID int       `json:"merchant_id"`
	Phone      string    `json:"phone"`
	Session    string    `json:"session"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func getDB() *gorm.DB {
	dsn := "host=localhost user=faridun password=123123123 dbname=newdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Could not connect to the database")
	}

	return db
}

func getPagesAndSessions(pageSize int) {
	var session []Sessions
	getSession := getDB().Where("merchant_id = ? and device not similar to ?", 512, "%(Laptop)%").Order("created_at desc")
	count := getDB().Where("merchant_id = ? and device not similar to ?", 512, "%(Laptop)%").Find(&session).RowsAffected

	var paginationSteps int

	if int(count)%pageSize != 0 {
		paginationSteps = int(count)/pageSize + 1
	} else {
		paginationSteps = int(count) / pageSize
	}

	if int(count) < pageSize {
		return
	}

	fmt.Println(paginationSteps)

	p := paginator.New(&paginator.Config{
		Limit: pageSize,
		Order: paginator.DESC,
	})

	var pageSessions []Sessions

	result, pageCursor, err := p.Paginate(getSession, &pageSessions)

	if err != nil {
		panic(err.Error())
	}

	if result.Error != nil {
		panic(result.Error.Error())
	}

	fmt.Println(pageSessions)
	fmt.Println(pageCursor)

	for i := 2; i <= paginationSteps; i++ {

		p = paginator.New(&paginator.Config{
			After: *pageCursor.After,
			Limit: pageSize,
			Order: paginator.DESC,
		})

		result, pageCursor, err = p.Paginate(getSession, &pageSessions)

		if err != nil {
			panic(err.Error())
		}

		if result.Error != nil {
			panic(result.Error.Error())
		}

		if err != nil {
			panic(err.Error())
		}

		if result.Error != nil {
			panic(result.Error.Error())
		}
		fmt.Println(pageSessions)
		fmt.Println(pageCursor)

	}
}

func getPaginated(page int, pageSize int) []Sessions {
	var sessions []Sessions

	offset := (page - 1) * pageSize
	limit := pageSize

	getDB().Raw("select * from sessions offset ? limit ?", offset, limit).Scan(&sessions)

	return sessions
}

func useCursorToGetPage(pageSize int, myCursor *paginator.Cursor) ([]Sessions, *paginator.Cursor, error) {

	getSession := getDB().Where("merchant_id = ? and device not similar to ?", 512, "%(Laptop)%").Order("created_at desc")
	p := paginator.New(&paginator.Config{})
	if myCursor.After != nil {
		p = paginator.New(&paginator.Config{
			After: *myCursor.After,
			Limit: pageSize,
			Order: paginator.DESC,
		})
	} else {
		p = paginator.New(&paginator.Config{
			Limit: pageSize,
			Order: paginator.DESC,
		})
	}

	var pageSessions []Sessions

	result, pageCursor, err := p.Paginate(getSession, &pageSessions)

	if result.Error != nil {
		panic(result.Error.Error())
	}

	return pageSessions, &pageCursor, err
}

func main() {

	var pageCursor paginator.Cursor // пустой курсор для использования функции.

	err := getDB().AutoMigrate(&Sessions{})
	if err != nil {
		return
	}

	//getPageSessions(6)

	page, cursor, err := useCursorToGetPage(2, &pageCursor) // первая страница
	if err != nil {
		return
	}
	fmt.Println(page)

	page, cursor, err = useCursorToGetPage(2, cursor) // вторая страница
	if err != nil {
		return
	}
	fmt.Println(page)

	page, cursor, err = useCursorToGetPage(2, cursor) // третяя страница
	if err != nil {
		return
	}
	fmt.Println(page)
}
