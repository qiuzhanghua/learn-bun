package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
)

type Book struct {
	bun.BaseModel `bun:"table:books"`
	ID            int64 `bun:"id,pk,autoincrement"`
	AuthorID      int64
	Author        *Author `bun:"rel:belongs-to,join:author_id=id"`
}

type Author struct {
	bun.BaseModel `bun:"table:authors"`
	ID            int64 `bun:"id,pk,autoincrement"`
	Book          *Book `bun:"rel:has-one"`
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
	err = InsertBookAndUser()
	fmt.Println(err)

	book := &Book{}
	err = db.NewSelect().Model(book). // model
						Relation("Author").                 // relation
						Where("`book`.`id` = 2").           // left join
						Where("`author`.`id` IS NOT NULL"). // simulate inner join
						Scan(ctx, book)
	fmt.Println(err)
	fmt.Println(book)
}

func InsertBookAndUser() error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Commit()

	// Insert a new Author
	author := &Author{}

	answer, err := tx.NewInsert().Model(author).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return errors.New("insert author: " + err.Error())
	}
	fmt.Println(author)
	book := &Book{}
	book.AuthorID = author.ID
	answer, err = tx.NewInsert().Model(book).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return errors.New("insert book: " + err.Error())
	}
	fmt.Println(answer)

	// Or update an existing Author
	book.AuthorID = author.ID
	answer, err = tx.NewUpdate().Model(book).Column("author_id").WherePK().Exec(ctx)
	if err != nil {
		tx.Rollback()
		return errors.New("update book: " + err.Error())
	}
	fmt.Println(answer)
	return nil
}
