package main

import (
	"fmt"

	"GoTuringCoffee/internal/service/web/model"
)

func main() {
	conf := model.MongoDBConfig{
		Url: "mongodb://turingcoffee:turingcoffeepassword@ds021000.mlab.com:21000/turing-coffee",
	}

	cModel := model.NewCookbookModel(&conf)
	cookbooks, err := cModel.ListCookbooks()
	if err != nil {
		panic(err)
	}
	fmt.Println(cookbooks)
}
