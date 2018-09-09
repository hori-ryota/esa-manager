package domain

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/hori-ryota/zaperr"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Post is struct for post
type Post struct {
	Name           string
	Number         uint
	Tags           []string
	Category       string
	FullName       string
	WIP            bool
	BodyMD         string
	RevisionNumber int64
}

func (post Post) ToTitle() string {
	tags := post.Tags
	for i := range tags {
		tags[i] = "#" + tags[i]
	}
	absName := path.Join(post.Category, post.Name)

	marks := make([]string, 0, 3)
	if post.Number > 0 {
		marks = append(marks, fmt.Sprintf("[id:%d]", post.Number))
	}

	if post.RevisionNumber > 0 {
		marks = append(marks, fmt.Sprintf("[rev:%d]", post.RevisionNumber))
	}

	if post.WIP {
		marks = append(marks, "[WIP]")
	}

	ss := make([]string, 0, 1+len(tags))
	ss = append(ss, absName)
	ss = append(ss, tags...)
	ss = append(ss, marks...)
	return strings.Join(ss, " ")
}

var idReg = regexp.MustCompile(` \[id:([0-9]+)\]$`)
var revisionReg = regexp.MustCompile(` \[rev:([0-9]+)\]$`)

func (post *Post) ParseTitle(title string) error {

	wip := strings.HasSuffix(title, " [WIP]")
	if wip {
		title = strings.TrimSuffix(title, " [WIP]")
	}
	post.WIP = wip

	revMatch := revisionReg.FindStringSubmatch(title)
	if len(revMatch) > 0 {
		title = strings.TrimSuffix(title, revMatch[0])
		revnum, err := strconv.ParseInt(revMatch[1], 10, 64)
		if err != nil {
			return zaperr.AppendFields(
				errors.Wrap(err, "faield to parse revision number"),
				zap.String("number string", revMatch[1]),
				zap.String("revision mark string", revMatch[0]),
			)
		}
		post.RevisionNumber = revnum
	}

	idMatch := idReg.FindStringSubmatch(title)
	if len(idMatch) > 0 {
		title = strings.TrimSuffix(title, idMatch[0])
		idnum, err := strconv.ParseUint(idMatch[1], 10, 64)
		if err != nil {
			return zaperr.AppendFields(
				errors.Wrap(err, "faield to parse id number"),
				zap.String("number string", idMatch[1]),
				zap.String("id mark string", idMatch[0]),
			)
		}
		post.Number = uint(idnum)
	}

	ss := strings.Split(title, " #")
	if len(ss) >= 2 {
		post.Tags = ss[1:]
	}

	absName := ss[0]

	post.Name = path.Base(absName)

	dir := path.Dir(absName)
	if strings.HasPrefix(absName, dir) {
		post.Category = dir
	}

	return nil
}
