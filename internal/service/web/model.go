package web

***REMOVED***
	"context"

	"github.com/globalsign/mgo"
***REMOVED***

type model struct {
	dbConf  *MongoDBConfig
	session *mgo.Session
	c       *mgo.Collection
***REMOVED***

func NewModel(dbConf *MongoDBConfig***REMOVED*** {
	return &model{
		dbConf: dbConf,
***REMOVED***
***REMOVED***

func (m *model***REMOVED*** Connect(ctx context.Context***REMOVED*** (err error***REMOVED*** {
	if m.session, err = mgo.Dial(m.dbConf.Url***REMOVED***; err != nil {
		return
***REMOVED***
	m.c = m.session.DB("turingcoffee"***REMOVED***.C("cookbook"***REMOVED***
	return
***REMOVED***

func (m *model***REMOVED*** Disconnect(***REMOVED*** {
	m.session.Close(***REMOVED***
***REMOVED***
