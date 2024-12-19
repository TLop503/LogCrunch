package structs

type Heartbeat struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Seq       int    `json:"seq"`
}

type HbSeqErr struct {
	ExpSeq  int `json:"exp_seq"`
	RecvSeq int `json:"recv_seq"`
}
