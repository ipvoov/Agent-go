package md

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino/schema"
)

func Split(indexer *milvus.Indexer) {
	ctx := context.Background()

	// 初始化分割器
	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":    "h1",
			"##":   "h2",
			"###":  "h3",
			"####": "h4",
		},
		TrimHeaders: false,
	})
	if err != nil {
		panic(err)
	}

	// 准备要分割的文档
	content, err := os.OpenFile("./document/恋爱常见问题和回答 - 单身篇.md", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	defer content.Close()
	bs, err := os.ReadFile("./document/恋爱常见问题和回答 - 单身篇.md")
	if err != nil {
		panic(err)
	}
	docs := []*schema.Document{
		{
			ID:      "document",
			Content: string(bs),
		},
	}

	// 执行分割
	results, err := splitter.Transform(ctx, docs)
	if err != nil {
		panic(err)
	}

	for i, doc := range results {
		doc.ID = doc.ID + "_" + strconv.Itoa(i)
	}

	// 处理分割结果
	for i, doc := range results {
		println("片段", i+1, ":", doc.Content)
		println("标题层级：")
		for k, v := range doc.MetaData {
			if k == "h1" || k == "h2" || k == "h3" || k == "h4" {
				println("  ", k, ":", v)
			}
		}
	}

	ids, err := indexer.Store(ctx, results)
	if err != nil {
		return
	}

	fmt.Println(ids)
}
