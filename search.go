package qanswer

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/ngaut/log"
	"github.com/silenceper/qanswer/util"
)

//SearchResult 搜索总数跟频率
type SearchResult struct {
	Sum  int32
	Freq int32
}

//GetSearchResult 返回各个搜索引擎返回的结果
func GetSearchResult(question string, answers []string) map[string][]*SearchResult {
	if question == "" {
		return nil
	}
	res := make(map[string][]*SearchResult)
	res["百度"] = baiduSearch(question, answers)
	return res
}

func baiduSearch(question string, answers []string) (result []*SearchResult) {
	resultMap := make(map[string]*SearchResult, len(answers))

	//搜索题目
	searchURL := fmt.Sprintf("http://www.baidu.com/s?wd=%s", url.QueryEscape(question))
	questionBody, err := util.HTTPGet(searchURL, 5)
	if err != nil {
		log.Errorf("search question:%s error", question)
		return
	}

	var wg sync.WaitGroup

	for k, answer := range answers {
		answer = plainAnswer(answer)
		answers[k] = answer

		wg.Add(1)
		go func(answer string) {
			defer wg.Done()
			searchResult := new(SearchResult)
			//题目搜索结果中包含的答案的数量
			searchResult.Freq = int32(strings.Count(string(questionBody), answer))

			//题目+结果搜索的总数
			keyword := fmt.Sprintf("%s %s", question, answer)
			searchURL := fmt.Sprintf("http://www.baidu.com/s?wd=%s", url.QueryEscape(keyword))
			body, err := util.HTTPGet(searchURL, 5)
			if err != nil {
				log.Errorf("search %s error", answer)
			} else {
				countRe, _ := regexp.Compile(`百度为您找到相关结果约([\d\,]+)`)
				result := countRe.FindAllStringSubmatch(string(body), -1)
				if len(result) > 0 {
					sum := result[0][1]
					sum = strings.Replace(sum, ",", "", -1)
					searchResult.Sum = util.MustInt32(sum)
				}
			}
			resultMap[answer] = searchResult
		}(answer)
	}
	wg.Wait()

	//将map转为slice 方便顺序输出
	for _, answer := range answers {
		result = append(result, resultMap[answer])
	}
	return result
}

//plainAnswer 去除答案中的 《》等字符
func plainAnswer(answer string) string {
	answer = strings.TrimPrefix(answer, "《")
	answer = strings.TrimSuffix(answer, "》")
	return answer
}
