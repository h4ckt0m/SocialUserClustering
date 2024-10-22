package models

type Twitteruser struct {
	Id 			string
	Name 		string
	Username 	string
	Description string
	Protected 	bool
	Followers 	[]Twitteruser
	Following 	[]Twitteruser
	Tweets []Tweet
}

type TwitterDAO struct {
	Id 			string
	Name 		string
	Username 	string
	Description string
	Protected 	bool
}

type FollowDAO struct {
	Id 			string
	IdFollowing string
}


type TwitterGetTweetResponse struct {
	Data []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
		Author_id string `json:"author_id"`
	} `json:"data"`
}

//TODO
type Tweet struct{
	Id 			string
	Text 		string
	Author 		string
}
	
type TwitterGetUserResponse struct {
	Data struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
		Description string `json:"description"`
	} `json:"data"`
}


	
type TwitterGetUsersResponse struct {
	Data []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
		Description string `json:"description"`
	} `json:"data"`
	Meta struct {
		Results_Count int `json:"result_count"`
		Next_Token string `json:"next_token"`
	} `json:"meta"`
}