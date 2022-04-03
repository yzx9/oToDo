package entity

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

var NewID func() int64

func startUpIDGenerator() error {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return fmt.Errorf("fails to create id generator")
	}

	NewID = func() int64 {
		return node.Generate().Int64()
	}

	return nil
}
