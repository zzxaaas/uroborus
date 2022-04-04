package model

type ProjectComment struct {
	Model
	OriginId  uint   `json:"origin_id" form:"origin_id"`
	TopicType int    `json:"topic_type" form:"topic_type"`
	Content   string `json:"content" form:"content"`
	FromUser  string `json:"from_user" form:"from_user"`
	ToUser    string `json:"to_user" form:"to_user"`
}
