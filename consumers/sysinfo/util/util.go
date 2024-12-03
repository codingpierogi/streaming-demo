package util

import (
	"encoding/json"
	"fmt"

	"github.com/codingpierogi/streaming-demo/types"
)

func DeserializeMessage(data []byte) (types.SysInfoMessage, error) {
	sim := types.SysInfoMessage{}
	err := json.Unmarshal(data, &sim)

	if err != nil {
		return sim, fmt.Errorf("failed to deserialize message: %w", err)
	}

	return sim, nil
}
