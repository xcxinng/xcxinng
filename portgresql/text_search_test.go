package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"testing"

	"github.com/go-ego/gse"
	"github.com/go-ego/gse/hmm/idf"
	"github.com/go-pg/pg/v10"
	_ "github.com/lib/pq"
)

var (
	// subjects     = []string{"研究生", "企业家", "建筑师", "律师", "记者", "厨师", "运动员", "教练", "导演", "歌手"}
	verbs        = []string{"探索", "推广", "实践", "改进", "优化", "创新", "传播", "应用", "普及", "革新"}
	objects      = []string{"新能源", "智慧城市", "数字货币", "元宇宙", "量子计算", "基因工程", "太空科技", "虚拟现实", "智能制造", "生物科技"}
	places       = []string{"在苏州", "在厦门", "在青岛", "在大连", "在天津", "在长沙", "在济南", "在郑州", "在昆明", "在福州"}
	modifiers    = []string{"持续地", "专注地", "积极地", "科学地", "有序地", "稳步地", "智慧地", "务实地", "创造性地", "前瞻性地"}
	englishWords = []string{"Web3", "NFT", "Metaverse", "DeFi", "DAO", "Zero Trust", "Edge Computing", "Digital Twin", "SaaS", "Quantum"}
	// Define some sample Chinese words or phrases
	words = []string{"探索", "未来的", "理解", "创新", "挑战", "崛起", "指南", "影响", "进步", "艺术", "科学", "神秘"}
	// Define some sample Chinese subjects
	subjects = []string{"技术", "人工智能", "区块链", "量子计算", "气候变化", "太空探索", "医疗", "教育", "金融", "网络安全", "机器人", "数据科学"}
)

// 生成随机句子
func generateRandomSentence() string {
	subject1 := subjects[rand.Intn(len(subjects))]
	verb1 := verbs[rand.Intn(len(verbs))]
	object1 := objects[rand.Intn(len(objects))]
	place1 := places[rand.Intn(len(places))]
	modifier1 := modifiers[rand.Intn(len(modifiers))]
	englishWord := englishWords[rand.Intn(len(englishWords))]
	randomNumber := rand.Intn(1000)

	return fmt.Sprintf("%s%s%s%s%s，(%s-%d)",
		subject1, place1, modifier1, verb1, object1, englishWord, randomNumber)
}

func generateRandomChineseTitle() string {
	// Seed the random number generator
	// rand.Seed(time.Now().UnixNano())

	// Randomly select a word and a subject
	word := words[rand.Intn(len(words))]
	subject := subjects[rand.Intn(len(subjects))]

	// Combine them to form a title
	title := fmt.Sprintf("%s%s", word, subject)
	return title
}

// 批量插入文档
func batchInsertDocuments(count int, batchSize int) error {
	for i := 0; i < count; i += batchSize {
		err := db.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
			for j := 0; j < batchSize && (i+j) < count; j++ {
				content := generateRandomSentence()
				title := generateRandomChineseTitle()
				titleTokens := generateTsvector(title)
				contentTokens := generateTsvector(content)
				_, err := tx.Exec(`INSERT INTO documents (title, content, title_tokens, content_tokens) VALUES (?, ?, ?, ?)`,
					title, content, titleTokens, contentTokens)
				if err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			return err
		}

		if i%10000 == 0 {
			log.Printf("Inserted %d documents", i)
		}
	}
	return nil
}

var testText = "大学生在广州大数据研究所工作，专门研究大数据"

func TestTextSearch(t *testing.T) {
	content := `大学生在广州很屌的研究所工作，专门研究大数据，专业是计算机科学`
	content = gse.FilterSymbol(content)
	t.Log(seg.CutSearch(content, true))
}

func TestBatchInsertDocuments(t *testing.T) {
	initDB()
	batchInsertDocuments(200000, 1000)
}

func TestInsertOne(t *testing.T) {
	title := "大学生搞研究"
	content := "大学生在广州搞研究，专门研究大数据，口号是：我爱你中国"
	titleTokens := generateTsvector(title)
	contentTokens := generateTsvector(content)
	_, err := db.Exec(`INSERT INTO documents (title, content, title_tokens, content_tokens) VALUES (?, ?, ?, ?)`,
		title, content, titleTokens, contentTokens)
	if err != nil {
		t.Error(err)
	}
}

func TestSearchDocumentsKeyWords(t *testing.T) {
	// query := "我爱你中国"
	query := "科技与人工智能"
	segments := seg.CutStop(query)
	query = strings.Join(segments, "|") // 且关系用 "&"

	s := `SELECT doc.id,
       doc.title,
       doc.content,
       ts_rank(doc.tsvector_title_content, query) AS score
FROM documents doc,
     to_tsquery('simple', '%s') query
WHERE doc.tsvector_title_content @@ query
ORDER BY score DESC
LIMIT 10;`
	statement := fmt.Sprintf(s, query)
	var documents []Document
	_, err := db.Query(&documents, statement)
	if err != nil {
		t.Error(err)
	}
	printJSONPretty(documents)
}

func printJSONPretty(v interface{}) {
	b, err := json.MarshalIndent(v, "", "    ") // 使用4个空格作为缩进
	if err != nil {
		fmt.Printf("error marshaling json: %v\n", err)
		return
	}
	fmt.Println(string(b))
}

func TestAnalyze(t *testing.T) {
	content := `《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的的科幻片.`
	content = gse.FilterSymbol(content)

	// 2. 使用 TextRankWithPOS 而不是 TextRank，可以指定要包含的词性
	var tr idf.TextRanker
	tr.WithGse(seg)
	// 这里扩大词性范围，包含更多类型的词
	allowPOS := []string{"n", "nr", "nz", "ns", "nt", "nw", "vn", "v", "vd", "vg", "a", "ad", "an"}
	results := tr.TextRankWithPOS(content, 0, allowPOS) // topK=0 表示返回所有结果
	t.Log(results)
}

func TestWordFrequencies(t *testing.T) {
	words := []string{"的", "是", "在", "有", "和", "与", "计算机"}
	for _, word := range words {
		freq, pos, exists := seg.Find(word)
		t.Logf("Word: %s, Freq: %.0f, Pos: %s, Exists: %v", word, freq, pos, exists)
	}
}

var text = "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的的科幻片."

func TestCustomDict(t *testing.T) {
	r := seg.CutSearch(testText)
	t.Log(r)
	// [大学 生在 广州 大 数据 研究 所 研究所 工作 ， 专门 研究 大 数据]

	seg.AddToken("大数据", 10000, "n")
	r = seg.CutSearch(testText)
	t.Log(r)
	// [大学 生在 广州 大数据 研究 所 研究所 工作 ， 专门 研究 大数据]

	// var te idf.TagExtracter
	// te.WithGse(seg)
	// err := te.LoadIdf()
	// fmt.Println("load idf: ", err)

	// var tr idf.TextRanker
	// tr.WithGse(seg)
	// results := tr.TextRank(text, 5)
	// fmt.Println("results: ", results)
}

func Test_generateTsvector(t *testing.T) {
	// a := generateTsvector(testText)
	// t.Log(a)

	a := generateTsvector(testText)
	t.Log(a)

	a = generateTsvector(text)
	t.Log(a)
}
