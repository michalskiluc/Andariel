package spider

import (
	"Andariel/mongo"
	"gopkg.in/mgo.v2"
	"github.com/robfig/cron"
	"github.com/andygrunwald/go-trending"

	"log"
	"time"
	"strconv"
	"gopkg.in/mgo.v2/bson"
)

type RequestServiceProvider struct {
}

var RequestService *RequestServiceProvider

var TrendingCollection *mgo.Collection


func PrepareTren() {
	TrendingCollection = mongo.GithubSession.DB(mongo.MDGitName).C("trending")
	idIndex := mgo.Index{
		Key: 		[]string{"name"},
		Unique: 	true,
		DropDups: 	true,
		Background: 	true,
		Sparse: 	true,
	}

	if err := TrendingCollection.EnsureIndex(idIndex); err != nil {
		panic(err)
	}
}

type Trending struct {
	CreateTime	string	 		`json:"create_time"`
	Repos 		[]trending.Project	`json:"repos"`
}

func (this *RequestServiceProvider)CronJob() {
	c := cron.New()
	c.AddFunc("@daily", RequestService.GetTrendingToday)
	c.Start()
	select {}
}

// 获取一天的 trending
func (this *RequestServiceProvider) GetTrendingToday() {

	trend := trending.NewTrending()
	// trending.TimeTodya 可以换成TimeWeek or TimeMonth 来获取本周或本月的 trending .
	result, err := trend.GetProjects(trending.TimeToday, "go")

	if err != nil {
		log.Print(err)
	}
	t := time.Now().Format("20060102")
	i := Trending{
		CreateTime: 	t,
		Repos: 		result,
	}

	err = TrendingCollection.Insert(i)

	if err != nil {
		log.Print(err)
	}
}

func (this *RequestServiceProvider) GetTrendingWeek() {

	trend := trending.NewTrending()
	result, err := trend.GetProjects(trending.TimeWeek, "")

	if err != nil {
		log.Print(err)
	}
	_, w := time.Now().ISOWeek()
	t := strconv.Itoa(w)
	i := Trending{
		CreateTime: 	t,
		Repos: 		result,
	}

	err = TrendingCollection.Insert(i)

	if err != nil {
		log.Print(err)
	}
}

func (this *RequestServiceProvider) GetTrendingMonth() {

	trend := trending.NewTrending()
	result, err := trend.GetProjects(trending.TimeMonth, "swift")

	if err != nil {
		log.Print(err)
	}
	m := time.Now().Month().String()
	i := Trending{
		CreateTime: 	m,
		Repos: 		result,
	}

	err = TrendingCollection.Insert(i)

	if err != nil {
		log.Print(err)
	}
}

//从数据库获取 trending
func (this *RequestServiceProvider) GetTrendingFromMD(tm string) ([]trending.Project, error) {
	var m Trending

	err := TrendingCollection.Find(bson.M{"createtime":tm}).One(&m)

	if err != nil {
		log.Print(err)
	}
	return m.Repos, err
}
