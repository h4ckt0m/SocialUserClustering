package main

import (
	"os"

	dao "./DAO"
	"./config"
	"./models"
	"./twitter"

	//Log library
	log "github.com/sirupsen/logrus"
)

var (
	cFile     = "config.json"
	daoObject = dao.TwitterDAO{}
)

func main() {
	//Load configuration
	conf, err := config.ParseConf(cFile)
	if err != nil {
		log.Fatalf("[SERVER] %s", err.Error())
	}

	//Create variables from configuration
	usernames := conf.Usernames
	apikey := conf.APIKEY
	daoObject.Server = conf.DBServer
	daoObject.Connect()

	if apikey == "" {
		log.Errorf("API_KEY cannot be empty")
		os.Exit(-1)
	}

	for _, username := range usernames {

		twitteruser := fillTwitterUserInfo(username, apikey)

		//Insert Into DB
		daoObject.InsertTwitterUser(twitteruser)

		var followersLevel2 []models.Twitteruser
		//var followersLevel3 []models.Twitteruser

		//Deep level 1
		for _, follower := range twitteruser.Followers {
			followerInfo := fillTwitterUserInfo(follower.Username, apikey)

			//Insert Into DB
			daoObject.InsertTwitterUser(followerInfo)
			followersLevel2 = append(followersLevel2, followerInfo.Followers...)
		}
	}
}

func fillTwitterUserInfo(username, apikey string) (twitteruser models.Twitteruser) {
	twitteruser.Username = username

	//Get ID from username
	twitteruser.Id = twitter.GetIDFromUsername(username, apikey)

	if twitteruser.Id == "" {
		log.Errorf("Couldnt find ID from username")
		os.Exit(-1)
	}

	//Get Description from ID
	twitteruser.Description = twitter.GetBioFromID(twitteruser.Id, apikey)

	//Get Followers from ID
	twitteruser.Followers = twitter.GetFollowersFromID(twitteruser.Id, apikey)

	//Get Following from ID
	twitteruser.Following = twitter.GetFollowingFromID(twitteruser.Id, apikey)

	//Get Tweets from ID
	twitteruser.Tweets = twitter.GetTweetsFromID(twitteruser.Id, apikey)

	return
}
