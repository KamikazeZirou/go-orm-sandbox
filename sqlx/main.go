package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
	"time"
	"unicode"
)

type Issue struct {
	Id          int
	Title       string
	Description string
	Labels      []*Label
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Label struct {
	Id   int
	Name string
}

func main() {
	// "parseTime=true"は、mysqlのTIMESTAMPをgoのtime.Timeにマッピングするのに必要。
	db, err := sqlx.Connect("mysql", "root:password@tcp(localhost:3306)/sandbox_db?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Goのフィールド名 ->  DBのフィールド名に変換する関数を指定する。
	// DBのフィールド名はsnake_caseのときに使うと楽。
	// (何も指定しないとstrings.ToLower相当の動作)
	db.MapperFunc(ToLowerSnakeCase)

	issues, err := getIssues(db)
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range issues {
		fmt.Println(*v)
	}

	issue, err := getIssue(db, 1)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(*issue)
}

// 単純なSelect
func getIssues(db *sqlx.DB) (issues []*Issue, err error) {
	err = db.Select(&issues, "SELECT id, title, description, created_at, updated_at FROM issues")
	return
}

// 表結合 + 条件指定をしてSelect + Nestした配列取得
func getIssue(db *sqlx.DB, id int) (*Issue, error) {
	namedArgs := map[string]interface{}{
		"issueId": id,
	}
	query := `
		SELECT id, title, description, created_at, updated_at
		FROM issues
		WHERE issues.id = :issueId
`
	query, args, err := sqlx.Named(query, namedArgs)
	if err != nil {
		return nil, err
	}

	issue := &Issue{}
	err = db.Get(issue, query, args...)
	if err != nil {
		return nil, err
	}

	// Nestした配列は別途取得しないといけない
	labels, err := getLabels(db, issue.Id)
	if err != nil {
		return nil, err
	}
	issue.Labels = labels

	return issue, nil
}

func getLabels(db *sqlx.DB, issueId int) ([]*Label, error) {
	namedArgs := map[string]interface{}{
		"issueId": issueId,
	}
	query := `
			SELECT id, name
			FROM labels
				INNER JOIN issue_labels il on il.label_id = labels.id
			WHERE il.issue_id = :issueId
`
	query, args, err := sqlx.Named(query, namedArgs)
	if err != nil {
		return nil, err
	}

	var labels []*Label
	err = db.Select(&labels, query, args...)
	if err != nil {
		return nil, err
	}

	return labels, nil
}

func ToLowerSnakeCase(s string) string {
	s = strings.TrimSpace(s)
	buffer := make([]rune, 0, len(s)*2)

	for i, v := range s {
		if unicode.IsUpper(v) {
			if i != 0 {
				buffer = append(buffer, '_')
			}

			buffer = append(buffer, unicode.ToLower(v))
		} else {
			buffer = append(buffer, v)
		}
	}

	h := string(buffer)
	return h
	//return string(buffer)
}
