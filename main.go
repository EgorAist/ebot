package main

import "fmt"

var url = "https://lenta.ru/rss/last24"

func main()  {
	items, err := ParsePage(url)
	if err != nil{
		fmt.Println(err)
		return
	}

	err = ES(items)
	if err != nil{
		fmt.Println(err)
		return
	}

	PrintItems(items)
}
