package model

import "time"

const (
	DEPLOY_STEP_CLONE     = 1
	DEPLOY_STEP_BUILD     = 2
	DEPLOY_STEP_RUN       = 3
	DEPLOY_STEP_RUNING    = 4
	DEPLOY_STEP_END       = 5
	DEPLOY_STATUS_RUNING  = 1
	DEPLOY_STATUS_SUCCESS = 2
	DEPLOY_STATUS_FAILED  = -1
)

type DeployHistory struct {
	Model
	OriginId uint          `json:"origin_id" form:"origin_id"`
	Branch   string        `json:"branch" form:"branch"`
	Image    string        `json:"image" form:"image"`
	Step     int           `json:"step" form:"step"`
	Status   int           `json:"status" form:"status"`
	Duration time.Duration `json:"duration" form:"duration"`
}
