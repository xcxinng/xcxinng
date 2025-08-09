create table documents (
    id serial primary key,
    title varchar(255) not null,
    content text not null,
    title_tokens varchar(512), -- 冗余存储title字段中文分词后的结果
    content_tokens text, -- 冗余存储content字段中文分词后的结果
    tsvector_title_content tsvector, -- 通过触发器自动合并title和content同时计算权重，最后生成tsvector
    -- created_at timestamp not null default now()
    -- created_by varchar(128) not null ,
    -- updated_at timestamp not null,
    -- updated_by varchar(128) not null
);

-- 创建函数
CREATE OR REPLACE FUNCTION update_documents_tokens()
RETURNS TRIGGER AS $$
BEGIN
  -- 使用to_tsvector将文本转换为分词向量
  -- 'A' 和 'B' 分别代表不同的权重，可以根据需要调整
  NEW.tsvector_title_content := setweight(to_tsvector('simple', NEW.title_tokens), 'A') ||
                setweight(to_tsvector('simple', NEW.content_tokens), 'B');
  RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';


-- 创建触发器
CREATE TRIGGER trigger_update_documents_tokens
BEFORE INSERT OR UPDATE ON documents
FOR EACH ROW
EXECUTE FUNCTION update_documents_tokens();

-- 创建索引
CREATE INDEX idx_documents_tsvector_title_content ON documents USING GIN(tsvector_title_content);

-- 查询示例
SELECT doc.id,
       doc.title,
       doc.content,
       ts_rank(doc.tsvector_title_content, query) AS score
FROM documents doc,
     to_tsquery('simple', '探索') query
WHERE doc.tsvector_title_content @@ query
ORDER BY score DESC
LIMIT 10;


-- 或模式搜索
explain (ANALYZE, COSTS, VERBOSE, BUFFERS) SELECT id,title,content
FROM documents WHERE tsvector_title_content @@
to_tsquery('simple', '科技 | 人工智能 ') LIMIT 10;

-- 与模式搜索
-- 搜索同时含有 科技 或 人工智能 的文档并按照得分倒序排列
SELECT doc.id,doc.title,doc.content,
--计算tsvector匹配query得分
ts_rank(doc.tsvector_title_content, query) AS score
FROM documents doc,to_tsquery('simple', '科技 & 人工智能 ') query
WHERE doc.tsvector_title_content @@ query
ORDER BY score DESC LIMIT 10;

-- 短语模式搜索
SELECT doc.id, doc.title, doc.content,
       ts_rank(doc.tsvector_title_content, query) AS score
FROM documents doc,
     to_tsquery('simple', '厨师 <-> 苏州') query
WHERE doc.tsvector_title_content @@ query
ORDER BY score DESC
LIMIT 10;
