// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hori-ryota/esa-go/esa"
	"github.com/hori-ryota/esa-manager/domain"
	"github.com/hori-ryota/zaperr"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var pullOpt = struct {
	apiToken string
	teamName string
	q        string
	dir      string
}{}

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull from esa",
	Long:  `pull from esa`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		esaClient := esa.NewClient(pullOpt.apiToken, pullOpt.teamName)
		page := uint(1)
		perPage := uint(100)

		posts := make([]domain.Post, 0, 100)

		for {
			logger.Info("fetch posts", zap.Uint("page", page), zap.Uint("perPage", perPage))
			resp, err := esaClient.ListPosts(
				ctx,
				esa.ListPostsParam{
					Q: pullOpt.q,
				},
				page,
				perPage,
			)
			if err != nil {
				return zaperr.AppendFields(
					err,
					zap.Uint("page", page),
					zap.Uint("perPage", perPage),
				)
			}
			for _, p := range resp.Posts {
				posts = append(posts, domain.Post{
					Name:           p.Name,
					Number:         p.Number,
					Tags:           p.Tags,
					Category:       p.Category,
					FullName:       p.FullName,
					WIP:            p.WIP,
					BodyMD:         p.BodyMD,
					RevisionNumber: p.RevisionNumber,
				})
			}

			if resp.PageResp.NextPage == nil {
				break
			}
			page = *resp.PageResp.NextPage
		}

		for _, p := range posts {
			f := filepath.Join(pullOpt.dir, filepath.FromSlash(p.ToTitle())) + ".md"
			err := os.MkdirAll(filepath.Dir(f), 0755)
			if err != nil {
				return zaperr.AppendFields(
					zaperr.Wrap(err, "failed to create dir"),
					zap.String("dir", filepath.Dir(f)),
					zap.String("filepath", f),
				)
			}
			err = ioutil.WriteFile(f, []byte(p.BodyMD), 0600)
			if err != nil {
				return zaperr.AppendFields(
					zaperr.Wrap(err, "failed to write file"),
					zap.String("filepath", f),
				)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().StringVarP(&pullOpt.apiToken, "apiToken", "a", "", "token for esa api (required)")
	err := pullCmd.MarkFlagRequired("apiToken")
	if err != nil {
		panic(err)
	}
	pullCmd.Flags().StringVarP(&pullOpt.teamName, "teamName", "t", "", "team name for esa api (required)")
	err = pullCmd.MarkFlagRequired("teamName")
	if err != nil {
		panic(err)
	}

	pullCmd.Flags().StringVarP(&pullOpt.q, "query", "q", "", "query for search. see https://docs.esa.io/posts/104")
	pullCmd.Flags().StringVarP(&pullOpt.dir, "dir", "d", ".", "dir for pulled files")
}
