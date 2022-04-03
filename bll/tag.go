package bll

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/infrastructure/util"
	"github.com/yzx9/otodo/model/entity"
)

// Update tag, should be called with `go UpdateTagAsync()`
func UpdateTag(todo *entity.Todo, oldTodoTitle string) error {
	// TODO[bug]: handle error
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
			exist, err := repository.ExistTag(userID, tagName)
			if err != nil {
				return util.NewErrorWithUnknown("unknown error: %w", err)
			}

			if !exist {
				tag := entity.Tag{
					Name:   tagName,
					UserID: userID,
					Todos:  make([]entity.Todo, 0),
				}
				if err := repository.InsertTag(&tag); err != nil {
					return fmt.Errorf("fails to create tag: %w", err)
				}
			}

			if err := repository.InsertTagTodo(userID, todo.ID, tagName); err != nil {
				return fmt.Errorf("fails to update tag: %w", err)
			}
		} else {
			// Remove old tag
			if err := repository.DeleteTagTodo(userID, todo.ID, tagName); err != nil {
				return fmt.Errorf("fails to update tag: %w", err)
			}
		}
	}

	return nil
}

func UpdateTagAsync(todo *entity.Todo, oldTodoTitle string) {
	if err := UpdateTag(todo, oldTodoTitle); err != nil {
		// TODO[bug]: handle error
		fmt.Println(err)
	}
}

var tagRegex = regexp.MustCompile(`^#(?P<tag>\\S{1,16}) `)

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
