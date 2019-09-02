package dao

var (
	OAuth2Client OAuth2Clienter = &oauth2Client{}
	Admin        Adminer        = &admin{}
	User         Userer         = &user{}
)
