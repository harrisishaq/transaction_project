package model

import (
	"context"
	"encoding/json"
)

type UserContext struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	IsAdmin  bool
}

type UserCtx string

var userCtx UserCtx

func SetUserContext(ctx context.Context, user interface{}) context.Context {
	return context.WithValue(ctx, userCtx, user)
}

func GetUserContext(ctx context.Context) (result *UserContext) {
	var user = ctx.Value(userCtx)
	userByte, _ := json.Marshal(user)
	result = new(UserContext)
	if err := json.Unmarshal(userByte, result); err != nil {
		result = nil
		return
	}
	return result
}
