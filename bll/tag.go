package bll

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

var tagRegex = regexp.MustCompile(`^#(?P<tag>\\S+) `)

func UpdateTag(userID, todoID, todoTitle, oldTodoTitle string) error {
	tags := make(map[string]bool)
	for {
		matches := tagRegex.FindStringSubmatch(todoTitle)
		if len(matches) == 0 {
			break
		}

		tags[matches[1]] = true
		todoTitle = strings.TrimLeft(todoTitle, " ")
	}

	oldTags := make(map[string]bool) // avoid duplicate tag
	for {
		matches := tagRegex.FindStringSubmatch(todoTitle)
		if len(matches) == 0 {
			break
		}

		oldTags[matches[1]] = true
		todoTitle = strings.TrimLeft(todoTitle, " ")
	}

	// diff
	for tagName := range oldTags {
		_, ok := tags[tagName]
		if ok {
			delete(tags, tagName)
		} else {
			tags[tagName] = false
		}
	}

	for tagName, op := range tags {
		if op {
			// Insert new tag
			if !dal.ExistTag(userID, tagName) {
				err := dal.InsertTag(entity.Tag{
					ID:     uuid.NewString(),
					Name:   tagName,
					UserID: userID,
					Todos:  make([]entity.Todo, 0, 1),
				})
				if err != nil {
					return fmt.Errorf("fails to create tag: %w", err)
				}
			}

			if err := dal.InsertTagTodo(userID, todoID, tagName); err != nil {
				return fmt.Errorf("fails to update tag: %w", err)
			}
		} else {
			// Remove old tag
			if err := dal.DeleteTagTodo(userID, todoID, tagName); err != nil {
				return fmt.Errorf("fails to update tag: %w", err)
			}
		}
	}

	return nil
}
