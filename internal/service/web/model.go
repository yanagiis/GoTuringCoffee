package web

***REMOVED***
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
***REMOVED***

type Model struct {
	dbConf  *MongoDBConfig
	session *mgo.Session
	c       *mgo.Collection
***REMOVED***

func NewModel(dbConf *MongoDBConfig***REMOVED*** *Model {
	return &Model{
		dbConf: dbConf,
***REMOVED***
***REMOVED***

func (m *Model***REMOVED*** ListCookbooks(***REMOVED*** ([]Cookbook, error***REMOVED*** {
	var cookbooks []bson.D
	if err := m.Connect(***REMOVED***; err != nil {
		return nil, err
***REMOVED***
	if err := m.c.Find(nil***REMOVED***.All(&cookbooks***REMOVED***; err != nil {
		return nil, err
***REMOVED***
***REMOVED***

func (m *Model***REMOVED*** GetCookbook(id string***REMOVED*** (Cookbook, error***REMOVED*** {
***REMOVED***

func (m *Model***REMOVED*** UpdateCookbook(id string, cookbook *Cookbook***REMOVED*** {
***REMOVED***

func (m *Model***REMOVED*** DeleteCookbook(id string***REMOVED*** {

***REMOVED***

func (m *Model***REMOVED*** Connect(***REMOVED*** (err error***REMOVED*** {
	if m.session == nil {
		if m.session, err = mgo.Dial(m.dbConf.Url***REMOVED***; err != nil {
			return
	***REMOVED***
***REMOVED***
	if m.c == nil {
		m.c = m.session.DB("turingcoffee"***REMOVED***.C("cookbook"***REMOVED***
***REMOVED***
	return
***REMOVED***

func (m *Model***REMOVED*** Disconnect(***REMOVED*** {
	if m.session != nil {
		m.session.Close(***REMOVED***
		m.session = nil
***REMOVED***
***REMOVED***
