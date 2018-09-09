package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostToTitle(t *testing.T) {
	for _, tt := range []struct {
		name  string
		post  Post
		title string
	}{
		{
			name: "only Name",
			post: Post{
				Name: "title",
			},
			title: "title",
		},
		{
			name: "with Tag",
			post: Post{
				Name: "title",
				Tags: []string{
					"tag1",
				},
			},
			title: "title #tag1",
		},
		{
			name: "with multi Tags",
			post: Post{
				Name: "title",
				Tags: []string{
					"tag1",
					"tag2",
				},
			},
			title: "title #tag1 #tag2",
		},
		{
			name: "with Category",
			post: Post{
				Name:     "title",
				Category: "cate/gory",
			},
			title: "cate/gory/title",
		},
		{
			name: "with WIP flag",
			post: Post{
				Name: "title",
				WIP:  true,
			},
			title: "title [WIP]",
		},
		{
			name: "with id",
			post: Post{
				Name:   "title",
				Number: 1,
			},
			title: "title [id:1]",
		},
		{
			name: "with revision number",
			post: Post{
				Name:           "title",
				RevisionNumber: 1,
			},
			title: "title [rev:1]",
		},
		{
			name: "with spaces",
			post: Post{
				Name: "title  title",
			},
			title: "title  title",
		},
		{
			name: "with All metas",
			post: Post{
				Name:     "title  title",
				Category: "cate/gory",
				Tags: []string{
					"tag1",
					"tag2",
				},
				Number:         1,
				RevisionNumber: 1,
				WIP:            true,
			},
			title: "cate/gory/title  title #tag1 #tag2 [id:1] [rev:1] [WIP]",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.title, tt.post.ToTitle())
		})
	}
}

func TestPostParseTitle(t *testing.T) {
	for _, tt := range []struct {
		name  string
		post  Post
		title string
	}{
		{
			name: "only Name",
			post: Post{
				Name: "title",
			},
			title: "title",
		},
		{
			name: "with Tag",
			post: Post{
				Name: "title",
				Tags: []string{
					"tag1",
				},
			},
			title: "title #tag1",
		},
		{
			name: "with multi Tags",
			post: Post{
				Name: "title",
				Tags: []string{
					"tag1",
					"tag2",
				},
			},
			title: "title #tag1 #tag2",
		},
		{
			name: "with Category",
			post: Post{
				Name:     "title",
				Category: "cate/gory",
			},
			title: "cate/gory/title",
		},
		{
			name: "with WIP flag",
			post: Post{
				Name: "title",
				WIP:  true,
			},
			title: "title [WIP]",
		},
		{
			name: "with id",
			post: Post{
				Name:   "title",
				Number: 1,
			},
			title: "title [id:1]",
		},
		{
			name: "with revision number",
			post: Post{
				Name:           "title",
				RevisionNumber: 1,
			},
			title: "title [rev:1]",
		},
		{
			name: "with spaces",
			post: Post{
				Name: "title  title",
			},
			title: "title  title",
		},
		{
			name: "with All metas",
			post: Post{
				Name:     "title  title",
				Category: "cate/gory",
				Tags: []string{
					"tag1",
					"tag2",
				},
				RevisionNumber: 1,
				Number:         1,
				WIP:            true,
			},
			title: "cate/gory/title  title #tag1 #tag2 [id:1] [rev:1] [WIP]",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			post := Post{}
			err := post.ParseTitle(tt.title)
			assert.NoError(t, err)
			assert.Equal(t, tt.post, post)
		})
	}
}
