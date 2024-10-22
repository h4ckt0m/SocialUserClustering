package twitter

import (
	//Parsing libraries/utils
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	//Log library
	log "github.com/sirupsen/logrus"

	"../models"
)

var (
	TwitterAPIGetUserByUsername = "https://api.twitter.com/2/users/by/username/$USERNAME"
	TwitterAPIGetBioByID        = "https://api.twitter.com/2/users/$ID"
	TwitterAPIGetFollowersByID  = "https://api.twitter.com/2/users/$ID/followers"
	TwitterAPIGetFollowingByID  = "https://api.twitter.com/2/users/$ID/following"
	TwitterAPIGetTweetsByID     = "https://api.twitter.com/2/users/$ID/tweets"
	BEARER_KEY                  = "Bearer $BEARER"
	seconds                     = 60 * 15 * time.Second
)

func GetIDFromUsername(username, apikey string) string {
	urlget := strings.Replace(TwitterAPIGetUserByUsername, "$USERNAME", username, -1)
	bodyResponse := getContentAPITwitter(urlget, apikey)

	var response models.TwitterGetUserResponse
	json.Unmarshal([]byte(bodyResponse), &response)

	return response.Data.ID
}

func GetBioFromID(id, apikey string) string {
	urlget := strings.Replace(TwitterAPIGetBioByID, "$ID", id, -1)
	urlget += "?user.fields=description"
	bodyResponse := getContentAPITwitter(urlget, apikey)

	var response models.TwitterGetUserResponse
	json.Unmarshal([]byte(bodyResponse), &response)

	return response.Data.Description
}

func GetFollowersFromID(id, apikey string) (followers []models.Twitteruser) {
	urlget := strings.Replace(TwitterAPIGetFollowersByID, "$ID", id, -1)
	urlget += "?user.fields=description&max_results=1000"
	var response models.TwitterGetUsersResponse
	next_token := ""

	for i := 0; i < 1000; i++ {
		bodyResponse := getContentAPITwitter(urlget+next_token, apikey)
		json.Unmarshal([]byte(bodyResponse), &response)

		for _, user := range response.Data {
			follower := models.Twitteruser{Id: user.ID, Name: user.Name, Username: user.Username, Description: user.Description}
			followers = append(followers, follower)
		}

		if strings.Compare(response.Meta.Next_Token, "&pagination_token="+response.Meta.Next_Token) == 0 {
			break
		}
		if response.Meta.Results_Count != 1000 {
			break
		}
		if response.Meta.Next_Token == "" {
			break
		}

		next_token = "&pagination_token=" + response.Meta.Next_Token
	}

	return
}

func GetFollowingFromID(id, apikey string) (following []models.Twitteruser) {
	urlget := strings.Replace(TwitterAPIGetFollowingByID, "$ID", id, -1)
	urlget += "?user.fields=description&max_results=1000"
	var response models.TwitterGetUsersResponse
	next_token := ""

	for i := 0; i < 1000; i++ {
		bodyResponse := getContentAPITwitter(urlget+next_token, apikey)
		json.Unmarshal([]byte(bodyResponse), &response)

		for _, user := range response.Data {
			follower := models.Twitteruser{Id: user.ID, Name: user.Name, Username: user.Username, Description: user.Description}
			following = append(following, follower)
		}

		if strings.Compare(response.Meta.Next_Token, "&pagination_token="+response.Meta.Next_Token) == 0 {
			break
		}
		if response.Meta.Results_Count != 1000 {
			break
		}
		if response.Meta.Next_Token == "" {
			break
		}

		next_token = "&pagination_token=" + response.Meta.Next_Token
	}

	return
}

func GetTweetsFromID(id, apikey string) (tweets []models.Tweet) {
	urlget := strings.Replace(TwitterAPIGetTweetsByID, "$ID", id, -1)
	urlget += "?exclude=retweets&max_results=10&tweet.fields=author_id"
	bodyResponse := getContentAPITwitter(urlget, apikey)

	var response models.TwitterGetTweetResponse
	json.Unmarshal([]byte(bodyResponse), &response)

	for _, tweet := range response.Data {
		tweets = append(tweets, models.Tweet{Id: tweet.ID, Text: tweet.Text, Author: tweet.Author_id})
	}

	return
}

//HELPERS//

func getContentAPITwitter(twitterurl, apikey string) (bodyResponse string) {
	//jsonParser *json.Decoder
	client := &http.Client{}
	urlGet := twitterurl

	req, err := http.NewRequest("GET", urlGet, nil)
	if err != nil {
		// This was an error, but not a timeout
		log.Errorf("[Twtitter API Miner] Not a valid url %v", err)
		//panic(err)
		return
	}

	authoritation := BEARER_KEY
	authoritation = strings.Replace(authoritation, "$BEARER", apikey, -1)

	//Set authentication as browser and create phpsession id
	req.Header.Set("Authorization", authoritation)
	req.Close = true

	resp, err := client.Do(req)
	if e, ok := err.(net.Error); ok && e.Timeout() {
		// This was a timeout
		log.Debugf("[SERVER] TIMEOUT: %s", urlGet)
		return
	} else if err != nil {
		// This was an error, but not a timeout
		log.Errorf("[Twtitter API Miner] Failed to issue GET request: %v", err)
		time.Sleep(1000 * time.Millisecond)
		resp, err = client.Do(req)
		if err != nil {
			log.Debugf("[Twtitter API Miner] Failed to issue GET 2nd request: %v", err)
			return
		}
		//panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// TODO: Is this really an error?
		log.Warningf("The server returned the status code %d for %s", resp.StatusCode, urlGet)
		// return nil, nil
		if resp.StatusCode == 429 {
			log.Warningf("Waiting %v for rate limit", seconds)
			time.Sleep(seconds)
		}
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	bodyResponse = string(b)
	return
}
