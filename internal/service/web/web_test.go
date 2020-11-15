package web

import (
	"context"

	repo "GoTuringCoffee/internal/service/web/model/repository"
)

func initRepositoryManager() (*repo.RepositoryManager, error) {
	dbConf := repo.MongoDBConfig{
		URL:      "mongodb+srv://turingcoffee:D56pNXo9bKlosH8W@cluster0.m5idb.gcp.mongodb.net/testturingcoffee?retryWrites=true&w=majority",
		Database: "testturingcoffee",
	}
	dbConf.Collections.Cookbook = "cookbooks"

	ctx := context.TODO()
	repoManager, err := repo.NewRepositoryManager(ctx, &dbConf)
	if err != nil {
		return nil, err
	}

	// Delete all cookbooks before testing
	repoManager.Cookbook.DeleteAll(ctx)

	return repoManager, nil
}
