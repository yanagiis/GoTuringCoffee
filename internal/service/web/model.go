package web

***REMOVED***
	"github.com/globalsign/mgo"
	"github.com/yanagiis/GoTuringCoffee/internal/service/lib"
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

func (m *Model***REMOVED*** ListCookbooks(***REMOVED*** ([]lib.Cookbook, error***REMOVED*** {
	var cookbooks []lib.Cookbook
	if err := m.Connect(***REMOVED***; err != nil {
		return nil, err
***REMOVED***
	if err := m.c.Find(nil***REMOVED***.All(&cookbooks***REMOVED***; err != nil {
		return nil, err
***REMOVED***
	return cookbooks, nil
***REMOVED***

func (m *Model***REMOVED*** GetCookbook(id string***REMOVED*** (*lib.Cookbook, error***REMOVED*** {
	var cookbook lib.Cookbook
	if err := m.Connect(***REMOVED***; err != nil {
		return nil, err
***REMOVED***
	if err := m.c.FindId(id***REMOVED***.One(&cookbook***REMOVED***; err != nil {
		return nil, err
***REMOVED***
	return &cookbook, nil
***REMOVED***

func (m *Model***REMOVED*** UpdateCookbook(id string, cookbook *lib.Cookbook***REMOVED*** error {
	if err := m.Connect(***REMOVED***; err != nil {
		return err
***REMOVED***
	return m.c.UpdateId(id, cookbook***REMOVED***
***REMOVED***

func (m *Model***REMOVED*** DeleteCookbook(id string***REMOVED*** error {
	if err := m.Connect(***REMOVED***; err != nil {
		return err
***REMOVED***
	return m.c.RemoveId(id***REMOVED***
***REMOVED***

func (m *Model***REMOVED*** Connect(***REMOVED*** (err error***REMOVED*** {
	if m.session == nil {
		if m.session, err = mgo.Dial(m.dbConf.Url***REMOVED***; err != nil {
			return
	***REMOVED***
***REMOVED***
	if m.c == nil {
		m.c = m.session.DB("turing-coffee"***REMOVED***.C("cookbooknew"***REMOVED***
***REMOVED***
	return
***REMOVED***

func (m *Model***REMOVED*** Disconnect(***REMOVED*** {
	if m.session != nil {
		m.session.Close(***REMOVED***
		m.session = nil
***REMOVED***
***REMOVED***
