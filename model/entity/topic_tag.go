package entity

type TopicTag struct {
	Id      int64 `json:"id"`
	TopicId int64 `json:"topicId"`
	TagId   int64 `json:"tagId"`
}

func (TopicTag) TableName() string {
	return "topic_tag"
}
