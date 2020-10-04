package model

import (
	"GoTuringCoffee/internal/service/lib"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func initRepositoryManager() (*RepositoryManager, error) {
	dbConf := MongoDBConfig{
		URL:      "mongodb+srv://turingcoffee:test12345@cluster0.m5idb.gcp.mongodb.net/testturingcoffee?retryWrites=true&w=majority",
		Database: "testturingcoffee",
		//URL:      "mongodb://turing-coffee:test12345@ds343718.mlab.com:43718/test-turing-coffee",
		//Database: "test-turing-coffee",
	}
	dbConf.Collections.Cookbook = "cookbooks"

	ctx := context.TODO()
	repoManager, err := NewRepositoryManager(ctx, &dbConf)
	if err != nil {
		return nil, err
	}

	return repoManager, nil
}

func TestListCookbooks(t *testing.T) {
	repoManager, err := initRepositoryManager()
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, repoManager)
	ctx := context.TODO()

	// Create cookbook
	newCookbook := lib.Cookbook{
		Name:        "new-cookbook",
		Description: "new cookbook",
		Tags:        []string{},
		Notes:       []string{},
		Processes: []lib.Process{
			lib.Process{
				ID:        "1",
				Name:      "Go Home",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Impl:      &lib.Circle{},
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	createdCookbook, err := repoManager.Cookbook.Create(ctx, &newCookbook)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdCookbook.ID)

	// List all cookbooks
	cookbooks, err := repoManager.Cookbook.List(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, cookbooks)
	assert.Equal(t, len(cookbooks), 1)

	// Get the cookbook
	result, err := repoManager.Cookbook.Get(ctx, createdCookbook.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, cookbooks)
	assert.Equal(t, createdCookbook.ID, result.ID)

	// Delete the cookbook
	err = repoManager.Cookbook.Delete(ctx, result)
	assert.NoError(t, err)

	// List all cookbooks
	cookbooks2, err := repoManager.Cookbook.List(ctx)
	assert.NoError(t, err)
	assert.Empty(t, cookbooks2)
	assert.Equal(t, len(cookbooks2), 0)
}
