package dao

import(
	//Internal Imports
	"../models"

	//DB 
	//"gopkg.in/mgo.v2/bson"

	"strings"

	//Log
	log "github.com/sirupsen/logrus"
)


func (m *TwitterDAO) InsertTwitterUser(twitteruser models.Twitteruser) {
	twitterDAO := models.TwitterDAO{Id:twitteruser.Id, Name:twitteruser.Name, Username:twitteruser.Username, Description:twitteruser.Description, Protected:twitteruser.Protected}
	err := db.C(users_collection).Insert(&twitterDAO)
	if err != nil && !strings.Contains(err.Error(),"duplicate key") {
		log.Errorf("[Twitter] Error inserting %v : %v",twitteruser.Id, err)
	}

	for _, follower := range twitteruser.Followers {
		m.InsertTwitterUser(follower)
		m.InsertTwitterFollow(follower.Id, twitteruser.Id)
	}

	for _, following := range twitteruser.Following {
		m.InsertTwitterUser(following)
		m.InsertTwitterFollow(twitteruser.Id, following.Id)
	}

	for _, tweet := range twitteruser.Tweets {
		m.InsertTwitterTweet(tweet)
	}

	return
}


func (m *TwitterDAO) InsertTwitterFollow(id, idfollowed string) {
	followDAO := models.FollowDAO{Id: id, IdFollowing: idfollowed}
	err := db.C(follows_collection).Insert(&followDAO)
	if err != nil && !strings.Contains(err.Error(),"duplicate key") {
		log.Errorf("[Twitter] Error inserting %v : %v",followDAO.Id, err)
	}
}

func (m *TwitterDAO) InsertTwitterTweet(tweet models.Tweet) {
	err := db.C(tweets_collection).Insert(&tweet)
	if err != nil && !strings.Contains(err.Error(),"duplicate key") {
		log.Errorf("[Twitter] Error inserting %v : %v",tweet.Id, err)
	}
	return
}