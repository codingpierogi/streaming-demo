package util

import (
	"encoding/json"
	"fmt"

	"github.com/codingpierogi/streaming-demo/types"
)

func SerializeMessage(sim types.SysInfoMessage) ([]byte, error) {
	serialized, err := json.Marshal(sim)

	if err != nil {
		return nil, fmt.Errorf("Failed to serialize message: %w", err)
	}

	return serialized, nil
}
