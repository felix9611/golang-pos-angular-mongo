package dto

import "golang-fixedasset-mongo-backend/backend/models"

type GeneralUpdateResponse struct {
	MatchedCount  int    `json:"MatchedCount"`
	ModifiedCount int    `json:"ModifiedCount"`
	UpsertedCount int    `json:"UpsertedCount"`
	UpsertedID    string `json:"UpsertedID"`
}

type GeneralUpdateInactiveResponseBody struct {
	Data    GeneralUpdateResponse `json:"data"`
	Id      string                `json:"id"`
	Message string                `json:"message"`
}

type GeneralCreateResponseBody struct {
	Id string `json:"id"`
}

type DataResponse struct {
	Data interface{} `json:"data"`
}

type AuthDataResponse struct {
	Data models.SysUsers `json:"data"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

/*
{"MatchedCount":1,"ModifiedCount":1,"UpsertedCount":0,"UpsertedID":null}
*/
