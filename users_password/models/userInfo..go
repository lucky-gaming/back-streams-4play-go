package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserInfo struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Login         string  `json:"login"`
    Email         *string `json:"email,omitempty"`
    Password      string  `json:"password"`
    CaptchaToken  string  `json:"captcha_token"`
    ActivateAccount *bool `json:"activate_account,omitempty"`
    AppSource     string  `json:"app_source"`
    TwoFactorCode *int    `json:"two_factor_code,omitempty"`
}
