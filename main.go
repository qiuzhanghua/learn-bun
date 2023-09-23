package main

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID    int64  `bun:"id,pk,autoincrement"`
	Name  string `bun:"name,notnull"`
	email string // unexported fields are ignored
}

func main() {
	defer db.Close()
	err := db.ResetModel(context.Background(), (*User)(nil))
	if err != nil {
		panic(err)
	}

	var num int
	num, err = db.NewSelect().Model((*User)(nil)).Count(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(num)
	//defer sqldb.Close()
}
