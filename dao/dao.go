package dao

var (
	// OAuth2Client ...
	OAuth2Client OAuth2Clienter = &oauth2Client{}
	// Admin ...
	Admin Adminer = &admin{}
	// User ...
	User Userer = &user{}
)
