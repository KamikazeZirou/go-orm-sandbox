package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	models "orm_sandbox/sqlboiler/my_models"
)

func main() {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/sandbox_db?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	ctx := context.Background()

	// 普通にSelect + Relationがあるデータも取得
	issues, err := models.Issues(qm.Load("Labels")).All(ctx, db)
	if err != nil {
		log.Fatalln(err)
	}

	for _, issue := range issues {
		fmt.Println(issue)
		for _, label := range issue.R.Labels {
			fmt.Println(label.Name)
		}
	}

	// 条件を指定してSelect + Relationがあるデータも取得
	issues, err = models.Issues(models.IssueWhere.ID.EQ(1), qm.Load("Labels")).All(ctx, db)
	if err != nil {
		log.Fatalln(err)
	}

	for _, issue := range issues {
		fmt.Println(issue)
		for _, label := range issue.R.Labels {
			fmt.Println(label.Name)
		}
	}
}
