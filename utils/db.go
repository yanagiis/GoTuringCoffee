package main

import (
	"fmt"

	"github.com/yanagiis/GoTuringCoffee/internal/service/web/model"
)

func main() {
	conf := model.MongoDBConfig{
		Url: "mongodb://<username>:<password>@ds021000.mlab.com:21000/turing-coffee",
	}

	cModel := model.NewCookbook(&conf)
	cookbooks, err := cModel.ListCookbooks()
	if err != nil {
		panic(err)
	}
	fmt.Println(cookbooks)
}
