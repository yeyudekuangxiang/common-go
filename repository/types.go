package repository

type GetUserBy struct {
	OpenId string
}
type FindTopicLikeBy struct {
	TopicId int
	UserId  int
}
type GetTopicLikeListBy struct {
	TopicIds []int
}
