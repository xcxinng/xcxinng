package elasticsearch

import (
	"flag"
	"fmt"
	"log"
	"regexp"

	"github.com/go-ego/gse"
	"github.com/go-ego/gse/hmm/idf"
	"github.com/go-ego/gse/hmm/pos"
)

var (
	seg    gse.Segmenter
	posSeg pos.Segmenter

	text  = "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的的科幻片."
	text1 = flag.String("text", text, "要分词的文本")

	text2 = "西雅图地标建筑, Seattle Space Needle, 西雅图太空针. Sky tree."
)

func example() {
	// Loading the default dictionary
	// 这个词典应该是给NLP训练用的
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

	// addToken()

	cutExample()
	// cut()
	// //
	// cutPos()
	// segCut()

	// extAndRank(seg)
}

func cutExample() {
	// Cut就是依据 以字典构建出来的NLP模型 对输入的文本进行分词
	// n := gse.FilterLang(text, "zh")
	// n := gse.FilterSymbol(text)
	// res := seg.CutSearch(n, true)
	// res = seg.Trim(res)
	// fmt.Println(res)

	t := `学生在广州研究所工作，研究大数据`
	fmt.Println(seg.CutSearch(t))

	// fmt.Println(seg.Cut("consistency function used to traverse the tree during search", true))
	// t := "平安科技<span xx>这里是文,&quot; &quot本 \n \r \\\\\" <span>"
	// ct := gse.FilterHtml(t)
	// ct = gse.FilterSymbol(ct)
	// fmt.Println(seg.Trim(seg.Cut(ct)))
	// 分析是干嘛的？
	// a := seg.Analyze(res, "")
	// fmt.Println(a)
	// res = seg.Trim(res)

}

func addToken() {
	err := seg.AddToken("《复仇者联盟3：无限战争》", 100, "n")
	fmt.Println("add token: ", err)
	seg.AddToken("西雅图中心", 100)
	seg.AddToken("西雅图太空针", 100, "n")
	seg.AddToken("Space Needle", 100, "n")
	// seg.AddTokenForce("上海东方明珠广播电视塔", 100, "n")
	//
	seg.AddToken("太空针", 100)
	seg.ReAddToken("太空针", 100, "n")
	freq, pos, ok := seg.Find("太空针")
	fmt.Println("seg.Find: ", freq, pos, ok)

	// seg.CalcToken()
	err = seg.RemoveToken("太空针")
	fmt.Println("remove token: ", err)
}

// 使用 DAG 或 HMM 模式分词
func cut() {
	// "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的的科幻片."

	// use DAG and HMM
	hmm := seg.Cut(text, true)
	fmt.Println("cut use hmm: ", hmm)
	// cut use hmm:  [《复仇者联盟3：无限战争》 是 全片 使用 imax 摄影机 拍摄 制作 的 的 科幻片 .]

	cut := seg.Cut(text)
	fmt.Println("cut: ", cut)
	// cut:  [《 复仇者 联盟 3 ： 无限 战争 》 是 全片 使用 imax 摄影机 拍摄 制作 的 的 科幻片 .]

	hmm = seg.CutSearch(text, true)
	fmt.Println("cut search use hmm: ", hmm)
	//cut search use hmm:  [复仇 仇者 联盟 无限 战争 复仇者 《复仇者联盟3：无限战争》 是 全片 使用 imax 摄影 摄影机 拍摄 制作 的 的 科幻 科幻片 .]
	fmt.Println("analyze: ", seg.Analyze(hmm, text))

	cut = seg.CutSearch(text)
	fmt.Println("cut search: ", cut)
	// cut search:  [《 复仇 者 复仇者 联盟 3 ： 无限 战争 》 是 全片 使用 imax 摄影 机 摄影机 拍摄 制作 的 的 科幻 片 科幻片 .]

	cut = seg.CutAll(text)
	fmt.Println("cut all: ", cut)
	// cut all:  [《复仇者联盟3：无限战争》 复仇 复仇者 仇者 联盟 3 ： 无限 战争 》 是 全片 使用 i m a x 摄影 摄影机 拍摄 摄制 制作 的 的 科幻 科幻片 .]

	s := seg.CutStr(cut, ", ")
	fmt.Println("cut all to string: ", s)
	// cut all to string:  《复仇者联盟3：无限战争》, 复仇, 复仇者, 仇者, 联盟, 3, ：, 无限, 战争, 》, 是, 全片, 使用, i, m, a, x, 摄影, 摄影机, 拍摄, 摄制, 制作, 的, 的, 科幻, 科幻片, .

	analyzeAndTrim(cut)

	reg := regexp.MustCompile(`(\d+年|\d+月|\d+日|[\p{Latin}]+|[\p{Hangul}]+|\d+\.\d+|[a-zA-Z0-9]+)`)
	text1 := `헬로월드 헬로 서울, 2021年09月10日, 3.14`
	hmm = seg.CutDAG(text1, reg)
	fmt.Println("Cut with hmm and regexp: ", hmm, hmm[0], hmm[6])
}

// seg.Analyze 分析分词结果,主要做以下工作：
// 计算每个词在文本中的位置信息（开始和结束位置）
// 为每个词添加词性标注（如果有的话）
// 计算每个词的词频权重
// seg.Trim 去掉 stop word
func analyzeAndTrim(cut []string) {
	a := seg.Analyze(cut, "")
	fmt.Println("analyze the segment: ", a)
	// analyze the segment:

	cut = seg.Trim(cut)
	fmt.Println("cut all: ", cut)
	// cut all:  [复仇者联盟3无限战争 复仇 复仇者 仇者 联盟 3 无限 战争 是 全片 使用 i m a x 摄影 摄影机 拍摄 摄制 制作 的 的 科幻 科幻片]

	fmt.Println(seg.String(text2, true))
	// 西雅图/nr 地标/n 建筑/n ,/x  /x seattle/x  /x space needle/n ,/x  /x 西雅图太空针/n ./x  /x sky/x  /x tree/x ./x
	fmt.Println(seg.Slice(text2, true))
	// [西雅图 地标 建筑 ,   seattle   space needle ,   西雅图太空针 .   sky   tree .]
}

// 词性标注
/*
seg.Pos 功能可以标注多种词性，不仅仅是名词（Nouns）和动词（Verbs）。
在自然语言处理中，词性标注（POS Tagging）通常会识别一系列的词性类别，包括但不限于：

名词（Nouns, NN）：表示人、地点、事物或概念的名称。
动词（Verbs, VB）：表示动作、状态或发生的事件。
形容词（Adjectives, JJ）：用来修饰名词，描述其属性或特征。
副词（Adverbs, RB）：用来修饰动词、形容词或其他副词，描述程度、方式、时间或地点。
代词（Pronouns, PRP）：代替名词的词，如我、你、他、她、它、我们、他们等。
介词（Prepositions, IN）：表示名词或代词与其他词之间的关系，如在、到、从、通过等。
连词（Conjunctions, CC）：连接单词、短语或从句的词，如和、但、或等。
冠词（Articles, DT）：限定名词的特殊词，如 a、an、the。
数词（Numerals, CD）：表示数量或顺序的词，如 one、two、first、second。
感叹词（Interjections, UH）：表达情感或反应的词，如 哦、啊、嗯等。
符号（Symbols, .）：标点符号，如句号、逗号、问号等。
情态动词（Modals, MD）：表示能力、许可或请求的动词，如 can、could、may、might、will、would 等。
分词（Gerunds, VBG）：以 -ing 结尾的动词形式，用作名词。
动名词（Infinitives, VB）：以 to 开头的动词形式，用作名词
*/
func cutPos() {
	// "西雅图地标建筑, Seattle Space Needle, 西雅图太空针. Sky tree."

	po := seg.Pos(text2, true)
	fmt.Println("pos: ", po)
	// pos:  [{西雅图 nr} {地标 n} {建筑 n} {, x} {  x} {seattle x} {  x} {space needle n} {, x} {  x} {西雅图太空针 n} {. x} {  x} {sky x} {  x} {tree x} {. x}]

	po = seg.TrimWithPos(po, "zg")
	fmt.Println("trim pos: ", po)
	// trim pos:  [{西雅图 nr} {地标 n} {建筑 n} {, x} {  x} {seattle x} {  x} {space needle n} {, x} {  x} {西雅图太空针 n} {. x} {  x} {sky x} {  x} {tree x} {. x}]

	posSeg.WithGse(seg)
	po = posSeg.Cut(text, true)
	fmt.Println("pos: ", po)
	// pos:  [{《 x} {复仇 v} {者 k} {联盟 j} {3 x} {： x} {无限 v} {战争 n} {》 x} {是 v} {全片 n} {使用 v} {imax eng} {摄影 n} {机 n} {拍摄 v} {制作 vn} {的的 u} {科幻 n} {片 q} {. m}]

	po = posSeg.TrimWithPos(po, "zg")
	fmt.Println("trim pos: ", po)
	// trim pos:  [{《 x} {复仇 v} {者 k} {联盟 j} {3 x} {： x} {无限 v} {战争 n} {》 x} {是 v} {全片 n} {使用 v} {imax eng} {摄影 n} {机 n} {拍摄 v} {制作 vn} {的的 u} {科幻 n} {片 q} {. m}]
}

// 使用最短路径和动态规划分词
func segCut() {
	segments := seg.Segment([]byte(*text1))
	fmt.Println(gse.ToString(segments, true))
	// 《/x 复仇/v 者/k 复仇者/n 联盟/j 3/x ：/x 无限/v 战争/n 》/x 是/v 全片/n 使用/v imax/x 摄影/n 机/n 摄影机/n 拍摄/v 制作/vn 的/uj 的/uj 科幻/n 片/q 科幻片/n ./x

	segs := seg.Segment([]byte(text2))
	// log.Println(gse.ToString(segs, false))
	log.Println(gse.ToString(segs))
	// 西雅图/nr 地标/n 建筑/n ,/x  /x seattle/x  /x space needle/n ,/x  /x 西雅图太空针/n ./x  /x sky/x  /x tree/x ./x

	// 搜索模式主要用于给搜索引擎提供尽可能多的关键字
	// segs := seg.ModeSegment(text2, true)
	log.Println("搜索模式: ", gse.ToString(segs, true))
	// 搜索模式:  西雅图/nr 地标/n 建筑/n ,/x  /x seattle/x  /x space needle/n ,/x  /x 西雅图太空针/n ./x  /x sky/x  /x tree/x ./x

	log.Println("to slice", gse.ToSlice(segs, true))
	// to slice [西雅图 地标 建筑 ,   seattle   space needle ,   西雅图太空针 .   sky   tree .]
}

func extAndRank(segs gse.Segmenter) {
	var te idf.TagExtracter
	te.WithGse(segs)
	err := te.LoadIdf()
	fmt.Println("load idf: ", err)

	segments := te.ExtractTags(text, 5)
	fmt.Println("segments: ", len(segments), segments)
	// segments:  5 [{科幻片 1.6002581704125} {全片 1.449761569875} {摄影机 1.2764747747375} {拍摄 0.9690261695075} {制作 0.8246043033375}]

	var tr idf.TextRanker
	tr.WithGse(segs)

	results := tr.TextRank(text, 5)
	fmt.Println("results: ", results)
	// results:  [{机 1} {全片 0.9931964427972227} {摄影 0.984870660504368} {使用 0.9769826633059524} {是 0.8489363954683677}]
}
