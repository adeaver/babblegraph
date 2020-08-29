package main

import (
	"babblegraph/crawler/web"
	"fmt"
)

func main() {
	testURL := "https://cnnespanol.cnn.com/2020/08/29/la-lucha-de-europa-contra-el-covid-19-pasa-de-los-hospitales-a-las-calles/"
	data, err := web.GetPageDataForURL(testURL)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(data.BodyText)
}
