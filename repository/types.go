package repository

type GetUserBy struct {
	OpenId string
}
type FindTopicLikeBy struct {
	TopicId int
	UserId  int
}
type GetTopicLikeListBy struct {
	TopicIds []int64
	UserId   int
}

type GetTopicPageListBy struct {
	ID         int `json:"id"`
	TopicTagId int `json:"topicTagId"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`  //limit为0时不限制数量
	UserId     int `json:"userId"` // 用于查询用户对帖子是否点赞
}
