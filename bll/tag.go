package bll

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

func UpdateTag(todo entity.Todo, oldTodoTitle string) error {
	if todo.Title == oldTodoTitle {
		return nil
	}

	tags := getTags(todo.Title)
	oldTags := getTags(oldTodoTitle)

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
			if !dal.ExistTag(userID, tagName) {
				if err := dal.InsertTag(entity.Tag{
					ID:     uuid.NewString(),
					Name:   tagName,
					UserID: userID,
					Todos:  make([]entity.Todo, 0),
				}); err != nil {
					return fmt.Errorf("fails to create tag: %w", err)
				}
			}

			if err := dal.InsertTagTodo(userID, todo.ID, tagName); err != nil {
				return fmt.Errorf("fails to update tag: %w", err)
			}
		} else {
			// Remove old tag
			if err := dal.DeleteTagTodo(userID, todo.ID, tagName); err != nil {
				return fmt.Errorf("fails to update tag: %w", err)
			}
		}
	}

	return nil
}

var tagRegex = regexp.MustCompile(`^#(?P<tag>\\S+) `)

func getTags(title string) map[string]bool {
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
