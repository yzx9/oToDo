package todo

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/yzx9/otodo/infrastructure/util"
)

type Tag struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string

	UserID int64

	Todos []int64
}

// Update tag, should be called with `go UpdateTagAsync()`
func UpdateTag(todo *Todo, oldTodoTitle string) error {
	// TODO[bug]: handle error
	if todo.Title == oldTodoTitle {
		return nil
	}

	tags := parseTags(todo.Title)
	oldTags := parseTags(oldTodoTitle)

	// diff, avoid duplicate tag in title
	for tagName := range oldTags {
		_, ok := tags[tagName]
		if ok {
			delete(tags, tagName)
		} else {
			tags[tagName] = false
		}
	}

	// TODO How to update shared user
	userID := todo.UserID
	for tagName, op := range tags {
		if op {
			// Insert new tag
			exist, err := TagRepository.Exist(userID, tagName)
			if err != nil {
				return util.NewErrorWithUnknown("unknown error: %w", err)
			}

			if !exist {
				tag := Tag{
					Name:   tagName,
					UserID: userID,
					Todos:  make([]int64, 0),
				}
				if err := TagRepository.Save(&tag); err != nil {
					return fmt.Errorf("fails to create tag: %w", err)
				}
			}

			if err := TagTodoRepository.Save(userID, todo.ID, tagName); err != nil {
				return fmt.Errorf("fails to update tag: %w", err)
			}
		} else {
			// Remove old tag
			if err := TagTodoRepository.Delete(userID, todo.ID, tagName); err != nil {
				return fmt.Errorf("fails to update tag: %w", err)
			}
		}
	}

	return nil
}

func UpdateTagAsync(todo *Todo, oldTodoTitle string) {
	if err := UpdateTag(todo, oldTodoTitle); err != nil {
		// TODO[bug]: handle error
		fmt.Println(err)
	}
}

var tagRegex = regexp.MustCompile(`^#(?P<tag>\\S{1,16}) `)

func parseTags(title string) map[string]bool {
	tags := make(map[string]bool)
	for {
		matches := tagRegex.FindStringSubmatch(title)
		if len(matches) == 0 {
			break
		}

		tags[matches[1]] = true
		title = strings.TrimLeft(title, " ")
	}
	return tags
}
