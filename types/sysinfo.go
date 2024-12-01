package types

type HostMemoryInfo struct {
	Total     uint64 `json:"total"`
	Used      uint64 `json:"used"`
	Available uint64 `json:"available"`
	Free      uint64 `json:"free"`
}

type SysInfoMessage struct {
	Key   string         `json:"key"`
	Value HostMemoryInfo `json:"value"`
}
