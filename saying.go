package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// var (
// 	notion_token = "secret_pPSHpp3siLtjYsC4VfdOiQXQcvqlD1muEjHbzlr54e6"
// 	block_id     = "758ea3e7-1a27-4f1a-973e-5e2f0095c0ef"
// 	database_id  = "a352c8223a374aad9294117670dee0bf"
// )

var (
	notion_token string
	block_id     string
	database_id  string
)

func main() {
	rand.Seed(time.Now().Unix())

	flag.StringVar(&notion_token, "token", "", "请输入 Notion Token")
	flag.StringVar(&block_id, "bid", "", "请输入 Block id")
	flag.StringVar(&database_id, "did", "", "请输入 Database id")
	flag.Parse()

	s := QueryDatabase()
	UpdateNotionBlock(s)

}

// 构建 Block 参数
func BuildBlockData(btype string, msg string) []byte {

	type Content struct {
		Content string `json:"content"`
	}
	type Text struct {
		Text *Content `json:"text"`
	}
	type RichText struct {
		RichText []*Text `json:"rich_text"`
	}

	block := map[string]*RichText{
		btype: {
			RichText: []*Text{
				{
					Text: &Content{
						Content: msg,
					},
				},
			},
		},
	}

	data_json, err := json.Marshal(block)

	if err != nil {
		fmt.Println("JSON Marshal Error", err)
	}

	return data_json
}

// 构建 query 参数
func BuildQueryData() []byte {
	type Sorts struct {
		Property  string `json:"property"`
		Direction string `json:"direction"`
	}

	type Content struct {
		Contains string `json:"contains"`
	}

	type Filter struct {
		Property     string   `json:"property"`
		Multi_select *Content `json:"multi_select"`
	}

	type Query struct {
		PageSize int      `json:"page_size"`
		Sorts    []*Sorts `json:"sorts"`
		Filter   *Filter  `json:"filter"`
	}

	query_data := &Query{
		PageSize: 100,
		Sorts: []*Sorts{
			{
				Property:  "Created At",
				Direction: "descending",
			},
		},
		Filter: &Filter{
			Property: "Tags",
			Multi_select: &Content{
				Contains: "一言",
			},
		},
	}

	data_json, err := json.Marshal(query_data)

	if err != nil {
		fmt.Println("JSON Marshal Error", err)
	}

	return data_json
}

func sendreq(method string, url string, body []byte) []byte {
	req, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewBuffer(body))

	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+notion_token)
	req.Header.Set("Notion-Version", "2022-02-22")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Go/1.18")

	// uri, _ := nurl.Parse("http://127.0.0.1:8888")

	client := &http.Client{
		// Transport: &http.Transport{
		// 	Proxy: http.ProxyURL(uri),
		// },
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return content
}

func UpdateNotionBlock(msg string) {
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s", block_id)

	data := BuildBlockData("callout", msg)

	content := sendreq("PATCH", url, data)
	now := time.Now().Format("2006-01-02 15:04:05")

	if strings.Contains(string(content), "last_edited_time") {
		fmt.Println(now, "修改成功！")
	} else {
		fmt.Println(now, "修改失败！")
	}
}

func QueryDatabase() string {
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", database_id)
	body := BuildQueryData()
	content := sendreq("POST", url, body)
	sentence := ParseQuery(content)
	randSentence := sentence[rand.Intn(len(sentence))]
	return randSentence
}

func ParseQuery(content []byte) []string {
	type Content struct {
		Content string `json:"content"`
	}

	type Text struct {
		Text *Content `json:"text"`
	}
	type Name struct {
		Title []*Text `json:"title"`
	}

	type Properties struct {
		Name *Name `json:"Name"`
	}

	type Sentence struct {
		Properties *Properties `json:"properties"`
	}

	type Data struct {
		Results []*Sentence `json:"results"`
	}

	var result Data

	json.Unmarshal(content, &result)

	var ss = []string{}

	for _, item := range result.Results {
		content := item.Properties.Name.Title[0].Text.Content
		ss = append(ss, content)
	}

	return ss
}
