package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-ego/gse"
	"github.com/go-pg/pg/v10"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "xianchaoxing"
	password = "123456"
	dbname   = "test"
)

var db *pg.DB

type dbLogger struct{}

func (d dbLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (d dbLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	if err != nil {
		return err
	}

	fmt.Printf("[DB] %s\n", query)
	return nil
}

func initDB() {
	var err error
	db = pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: dbname,
		Addr:     fmt.Sprintf("%s:%d", host, port),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AddQueryHook(dbLogger{})
}

func initTable() {
	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm;`)
	if err != nil {
		log.Fatalf("Failed to create extension: %v", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS documents (
            id SERIAL PRIMARY KEY,
            content TEXT,
            tsvector_col TSVECTOR
        );
    `)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	_, err = db.Exec(`
        CREATE INDEX IF NOT EXISTS tsvector_idx ON documents USING gin(tsvector_col);
    `)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}
}

var seg gse.Segmenter

func init() {
	err := seg.LoadDict("zh_s")
	if err != nil {
		log.Fatalf("Failed to load dictionary: %v", err)
	}

	err = seg.LoadStop("zh")
	if err != nil {
		log.Printf("Failed to load stop words: %v", err)
	}

	err = seg.LoadDict("/Users/xianchaoxing/go/pkg/mod/github.com/go-ego/gse@v0.80.3/data/dict/en/dict.txt")
	if err != nil {
		log.Fatalf("Failed to load dictionary: %v", err)
	}

	initDB()
}

func insertDocument(content string) error {
	// n := gse.FilterLang(content, "zh")
	seg.RemoveStop(content)
	// res := seg.CutSearch(content, true)
	// tsvector := fmt.Sprintf("'%s'", strings.Join(res, " "))
	// fmt.Println(content, tsvector)

	// 使用 CutAll 模式获取所有可能的分词结果
	segments := seg.CutAll(content)
	// 获取带权重的分词结果
	weights := seg.Analyze(segments, content)

	// 构建带权重的 tsvector
	var weightedTerms []string
	for _, weight := range weights {
		// 根据权重设置不同的优先级 (A > B > C > D)
		priority := "D"
		if weight.Freq > 5000 {
			priority = "A"
		} else if weight.Freq > 2500 {
			priority = "B"
		} else if weight.Freq > 1000 {
			priority = "C"
		}
		weightedTerms = append(weightedTerms,
			fmt.Sprintf("%s:%s", weight.Text, priority))
	}

	tsvector := fmt.Sprintf("'%s'", strings.Join(weightedTerms, " "))
	_, err := db.Exec(`INSERT INTO documents (content, tsvector_col) VALUES (?, to_tsvector(?))`, content, tsvector)
	if err != nil {
		return err
	}
	return nil
}

func generateTsvector(text string) string {
	// 使用 CutStop 模式 分词+去掉停用词
	segments := seg.CutStop(text)

	// 使用 CutSearch 模式获取所有可能的分词结果
	// segments := seg.CutSearch(text)
	// 分析词性
	// weights := seg.Analyze(segments, text)

	tsvector := fmt.Sprintf("'%s'", strings.Join(segments, " "))
	return tsvector
}

func searchDocuments(query string) ([]Document, error) {
	if query == "" {
		return nil, errors.New("query is empty")
	}
	var results []Document
	err := db.Model(&results).
		Column("content").
		Where(`tsvector_col @@ to_tsquery(?)`, query).
		Select()
	if err != nil {
		return nil, err
	}
	return results, nil
}

type Document struct {
	ID                   int     `pg:",pk" json:"id"`
	Title                string  `pg:"title" json:"title"`
	TitleTokens          string  `pg:"title_tokens" json:"title_tokens,omitempty"`
	Content              string  `pg:"content" json:"content"`
	ContentTokens        string  `pg:"content_tokens" json:"content_tokens,omitempty"`
	TsvectorTitleContent string  `pg:"tsvector_title_content" json:"tsvector_title_content,omitempty"`
	Score                float64 `pg:"score" json:"score,omitempty"`
}

type ContentParam struct {
	Content string `json:"content"`
}

// func main() {
// 	initDB()
// 	initTable()

// 	router := gin.Default()

// 	router.POST("/insert", func(c *gin.Context) {
// 		p := ContentParam{}
// 		if err := c.ShouldBindJSON(&p); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		if err := insertDocument(p.Content); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"message": "Document inserted successfully"})
// 	})

// 	router.GET("/search", func(c *gin.Context) {
// 		query := c.Query("query")
// 		results, err := searchDocuments(query)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"results": results})
// 	})

// 	router.Run(":8089") // listen and serve on 0.0.0.0:8080
// }
