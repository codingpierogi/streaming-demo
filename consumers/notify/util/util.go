package util

import (
	"encoding/json"
	"fmt"

	"github.com/codingpierogi/streaming-demo/types"
)

func DeserializeMessage(data []byte) (types.NotifyMessage, error) {
	nm := types.NotifyMessage{}
	err := json.Unmarshal(data, &nm)

	if err != nil {
		return nm, fmt.Errorf("failed to deserialize message: %w", err)
	}

	return nm, nil
}
