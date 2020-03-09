package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type githubConfig struct {
	Token       string `json:"token"`
	SourceRepo  string `json:"repository"`
	AuthorName  string `json:"username"`
	AuthorEmail string `json:"email"`
}

var client *github.Client
var ctx = context.Background()

func getTree(ref *github.Reference, opts *githubConfig, sourceFiles string) (tree *github.Tree, err error) {
	content, err := ioutil.ReadFile(sourceFiles)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	entries := []*github.TreeEntry{
		{
			Path:    github.String(filepath.Base(sourceFiles)),
			Type:    github.String("blob"),
			Content: github.String(string(content)),
			Mode:    github.String("100644"),
		},
	}

	tree, _, err = client.Git.CreateTree(ctx, opts.AuthorName, opts.SourceRepo, *ref.Object.SHA, entries)
	return tree, err
}

func pushCommit(ref *github.Reference, tree *github.Tree, opts *githubConfig) (err error) {

	parent, _, err := client.Repositories.GetCommit(ctx, opts.AuthorName, opts.SourceRepo, *ref.Object.SHA)
	if err != nil {
		return err
	}
	parent.Commit.SHA = parent.SHA

	commitMessage := "file updated"
	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: &opts.AuthorName, Email: &opts.AuthorEmail}
	commit := &github.Commit{Author: author, Message: &commitMessage, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(ctx, opts.AuthorName, opts.SourceRepo, commit)
	if err != nil {
		return err
	}

	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, opts.AuthorName, opts.SourceRepo, ref, false)
	return err
}

func updateFile(sourceFiles string, opts *githubConfig) bool {

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: opts.Token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	ref, _, err := client.Git.GetRef(ctx, opts.AuthorName, opts.SourceRepo, "refs/heads/master")
	if err != nil {
		return false
	}
	if ref == nil {
		return false
	}

	tree, err := getTree(ref, opts, sourceFiles)
	if err != nil {
		return false
	}

	if err := pushCommit(ref, tree, opts); err != nil {
		return false
	}
	return true
}
