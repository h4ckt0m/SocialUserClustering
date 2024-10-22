package dao

import(
	//DB 
	mgo "gopkg.in/mgo.v2"

	//Log library
	log "github.com/sirupsen/logrus"

	//Time libraries
	"time"
)

//Object to create initial connection to DB
type TwitterDAO struct {
	Server   string
	Database string
}

//Variable definition to access DB
var db *mgo.Database

//Collections constant names
const (
    //DBs
    twitter_DB = "twitter"

    //Collections
    users_collection = "users"
    follows_collection = "follows"
    tweets_collection = "tweets"
)


//Initial connection to database
func (m *TwitterDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}

	session.SetMode(mgo.Monotonic, true)
	session.SetSocketTimeout(1 * time.Hour)
	db = session.DB(twitter_DB)

	//TODO ENSURE INDEXES
	collections, err := db.CollectionNames()
	if err != nil {
		log.Errorf("Error obtaining collection names: %v", err.Error())
	}

	ensureIndexes(collections,users_collection,[]string{"id"})
	ensureIndexes(collections,follows_collection,[]string{"id","idfollowing"})
	ensureIndexes(collections,tweets_collection,[]string{"id"})
}

func (m *TwitterDAO) ReloadSession() error {

	session, err := mgo.Dial(m.Server)
	if err != nil {
		return err
	}

	session.SetMode(mgo.Monotonic, true)
	session.SetSocketTimeout(1 * time.Hour)
	db = session.DB(m.Database)

	return nil
}

/////////////////////////

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// INTERNAL METHODS //

//Ensure unique indexes
func ensureIndexes(collections []string, collectionName string, indexUnique []string){
	//Check if collection already exists
	/*exist := false
	for _,collection := range collections{
		if collection == collectionName {
			exist = true
			break
		}
	}
*/
	//If doesnt, create index unique
	//if !exist {
		idx := mgo.Index{
			Key: indexUnique,
			DropDups: true,
			Unique: true,
		}

		db.C(collectionName).EnsureIndex(idx)
	//}
}