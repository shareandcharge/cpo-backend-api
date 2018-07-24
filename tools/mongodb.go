package tools

//This class connects to the Charge & Fuel Mobile Application to get a list of the drivers

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Driver struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Email   string        `bson:"email" json:"email"`
	About   interface{}   `bson:"about" json:"about"`
	Address string        `bson:"address" json:"address"`
	Balance float64 `json:"balance"`
	Token string `json:"token"`
}

// returns a list of all the users from the mobile app
func ReturnAllDrivers() ([]Driver, error) {
	url := "mongodb://sc:Thahm6uudaiweifa@18.184.196.227:27017/fsdata"

	var Session *mgo.Session

	var err error
	Session, err = mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	collection := Session.DB("fsdata").C("users")

	var drivers []Driver
	err = collection.Find(nil).All(&drivers)
	if err != nil {
		return nil, err
	}

	Session.Close()

	return drivers, nil

}
