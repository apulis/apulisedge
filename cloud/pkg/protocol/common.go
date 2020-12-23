// compatible with kubeedge protocol

package protocol

// ApulisHeader struct
type ApulisHeader struct {
	ClusterId int64 `json:"clusterId"`
	GroupId   int64 `json:"groupId"`
	UserId    int64 `json:"userId"`
}
