package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

var tags = make(map[string]entity.Tag)

func InsertTag(tag entity.Tag) error {
	tags[tag.ID] = tag
	return nil
}

func GetTag(userID, tagName string) (entity.Tag, error) {
	for _, tag := range tags {
		if tag.UserID == userID && tag.Name == tagName {
			return tag, nil
		}
	}

	return entity.Tag{}, utils.NewErrorWithNotFound("tag not found")
}

func GetTags(userID string) ([]entity.Tag, error) {
	vec := make([]entity.Tag, 0)
	for _, tag := range tags {
		if tag.UserID == userID {
			vec = append(vec, tag)
		}
	}

	return vec, nil
}

func InsertTagTodo(userID, todoID, tagName string) error {
	for _, tag := range tags {
		if tag.UserID == userID && tag.Name == tagName {
			tag.Todos = append(tag.Todos, entity.Todo{
				// TODO this is a known BUG, but can not be fixed
				// before database setup.
				ID: todoID,
			})
			tags[tag.ID] = tag
			return nil
		}
	}

	return utils.NewErrorWithNotFound("tag not found")
}

func DeleteTagTodo(userID, todoID, tagName string) error {
	for _, tag := range tags {
		if tag.UserID != userID || tag.Name != tagName {
			continue
		}

		for i, todo := range tag.Todos {
			if todo.ID != todoID {
				continue
			}

			if len(tag.Todos) == 1 {
				delete(tags, tag.ID)
				return nil
			}

			tag.Todos = append(tag.Todos[:i-1], tag.Todos[i+1:]...)
			tags[tag.ID] = tag
			return nil
		}

		break
	}

	return utils.NewErrorWithNotFound("tag not found")
}

func ExistTag(userID, tagName string) bool {
	for _, tag := range tags {
		if tag.UserID == userID && tag.Name == tagName {
			return true
		}
	}

	return false
}
