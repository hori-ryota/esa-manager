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
	"strings"

	"github.com/hori-ryota/esa-go/esa"
	"github.com/hori-ryota/esa-manager/domain"
	"github.com/hori-ryota/zaperr"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var pushOpt = struct {
	apiToken string
	teamName string
	dir      string
}{}

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push to esa",
	Long:  `push to esa`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		esaClient := esa.NewClient(pushOpt.apiToken, pushOpt.teamName)

		for _, filename := range args {
			title := strings.TrimSuffix(filepath.ToSlash(filename), ".md")

			post := domain.Post{}
			err := post.ParseTitle(title)
			if err != nil {
				return zaperr.AppendFields(
					zaperr.Wrap(err, "failed to parse title"),
					zap.String("filename", filename),
					zap.String("title", title),
				)
			}

			srcFilepath := filepath.Join(pushOpt.dir, filename)

			body, err := ioutil.ReadFile(srcFilepath)
			if err != nil {
				return zaperr.AppendFields(
					zaperr.Wrap(err, "failed to read file"),
					zap.String("dir", pushOpt.dir),
					zap.String("filename", filename),
				)
			}

			var after esa.Post

			if post.Number == 0 {
				// new post
				b := string(body)
				resp, err := esaClient.CreatePost(
					ctx,
					esa.CreatePostParam{
						Name:     post.Name,
						BodyMD:   &b,
						Tags:     &post.Tags,
						Category: &post.Category,
						WIP:      &post.WIP,
					},
				)
				if err != nil {
					return zaperr.AppendFields(
						zaperr.Wrap(err, "failed to create post"),
						zap.String("filename", filename),
					)
				}
				after = *resp
				logger.Info("post created", zap.Any("created", after))
			} else {
				// update post
				b := string(body)
				resp, err := esaClient.UpdatePost(
					ctx,
					post.Number,
					esa.UpdatePostParam{
						Name:     &post.Name,
						BodyMD:   &b,
						Tags:     &post.Tags,
						Category: &post.Category,
						WIP:      &post.WIP,
					},
				)
				if err != nil {
					return zaperr.AppendFields(
						zaperr.Wrap(err, "failed to update post"),
						zap.String("filename", filename),
					)
				}
				logger.Info(
					"post updated",
					zap.Any("updated", resp.Post),
					zap.Bool("overlapped", resp.Overlapped),
				)
				after = resp.Post
			}

			f := filepath.Join(pushOpt.dir, filepath.FromSlash(domain.Post{
				Name:           after.Name,
				Number:         after.Number,
				Tags:           after.Tags,
				Category:       after.Category,
				FullName:       after.FullName,
				WIP:            after.WIP,
				BodyMD:         after.BodyMD,
				RevisionNumber: after.RevisionNumber,
			}.ToTitle()+".md"))
			if f == srcFilepath {
				continue
			}
			err = os.Rename(srcFilepath, f)
			if err != nil {
				return zaperr.AppendFields(
					zaperr.Wrap(err, "failed to rename file"),
					zap.String("srcFilepath", srcFilepath),
					zap.String("dstFilepath", f),
				)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringVarP(&pushOpt.apiToken, "apiToken", "a", "", "token for esa api (required)")
	err := pushCmd.MarkFlagRequired("apiToken")
	if err != nil {
		panic(err)
	}
	pushCmd.Flags().StringVarP(&pushOpt.teamName, "teamName", "t", "", "team name for esa api (required)")
	err = pushCmd.MarkFlagRequired("teamName")
	if err != nil {
		panic(err)
	}

	pushCmd.Flags().StringVarP(&pushOpt.dir, "dir", "d", ".", "dir for pushed files")
}
