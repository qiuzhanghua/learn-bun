package main

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
)

type Book struct {
	ID       int64 `bun:"id,pk,autoincrement"`
	AuthorID int64
	Author   *Author `bun:"rel:belongs-to,join:author_id=id"`
}

type Author struct {
	ID   int64 `bun:"id,pk,autoincrement"`
	Book *Book `bun:"rel:has-one"`
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID    int64  `bun:"id,pk,autoincrement"`
	Name  string `bun:"name,notnull"`
	email string // unexported fields are ignored
}

type Order struct {
	ID int64 `bun:",pk"`
	// Order and Item in join:Order=Item are fields in OrderToItem model
	Items []Item `bun:"m2m:order_to_items,join:Order=Item"`
}

type Item struct {
	ID int64 `bun:",pk"`
}

type OrderToItem struct {
	OrderID int64  `bun:",pk"`
	Order   *Order `bun:"rel:belongs-to,join:order_id=id"`
	ItemID  int64  `bun:",pk"`
	Item    *Item  `bun:"rel:belongs-to,join:item_id=id"`
}

func main() {
	defer db.Close()

	ctx := context.Background()
	db.RegisterModel((*OrderToItem)(nil))

	err := db.ResetModel(ctx, (*User)(nil), (*Author)(nil), (*Book)(nil), (*Order)(nil), (*Item)(nil), (*OrderToItem)(nil))
	if err != nil {
		panic(err)
	}

	db.NewDropTable().Model((*OrderToItem)(nil)).IfExists().Exec(ctx)
	db.NewDropTable().Model((*Order)(nil)).IfExists().Exec(ctx)
	db.NewDropTable().Model((*Item)(nil)).IfExists().Exec(ctx)
	db.NewDropTable().Model((*User)(nil)).IfExists().Exec(ctx)
	db.NewDropTable().Model((*Book)(nil)).IfExists().Exec(ctx)
	db.NewDropTable().Model((*Author)(nil)).IfExists().Exec(ctx)

	db.NewCreateTable().Model((*Item)(nil)).IfNotExists().Exec(ctx)
	db.NewDropTable().Model((*Order)(nil)).IfExists().Exec(ctx)
	db.NewCreateTable().Model((*OrderToItem)(nil)).IfNotExists().Exec(ctx)
	db.NewCreateTable().Model((*User)(nil)).IfNotExists().Exec(ctx)
	db.NewCreateTable().Model((*Author)(nil)).IfNotExists().Exec(ctx)
	db.NewCreateTable().Model((*Book)(nil)).IfNotExists().Exec(ctx)

	var num int
	num, err = db.NewSelect().Model((*User)(nil)).Count(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(num)

}
