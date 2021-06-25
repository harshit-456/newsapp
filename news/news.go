package news

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)


type Client struct{
	http  *http.Client
	key string 
	PageSize int
}
/*
httpClient field points to HTTP client used to make requests
key holds api key
pagesisze holds no of results to return per page(max of 100)

*/
type Article struct {
	Source struct {
		ID   interface{} `json:"id"`
		Name string      `json:"name"`
	} `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}
//Time is a struct to which Date() fucntion is associated,it returns year,month and day
//adding a function to article struct
func (a *Article) FormatPublishedDate() string {
	year, month, day := a.PublishedAt.Date()
	return fmt.Sprintf("%v %d, %d", month, day, year)
}
type Results struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}
func NewsClient(httpClient *http.Client,key string,pagesize int) *Client{
	if pagesize>100{
		pagesize=100
	}
	return &Client{httpClient,key,pagesize}//a new Client instance is returned
}
func FetchEverything(query,page string,c *Client)(*Results,error){
	endpoint:=fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%s&apiKey=%s&sortBy=publishedAt&language=en",url.QueryEscape(query),c.PageSize,page,c.key)//the resulting string will look like first arg of sprintf,after appending query,pagesize and key
	resp,err:=c.http.Get(endpoint)//making http request
	if err!=nil{
		return nil,err
	}
	defer resp.Body.Close()
	body,err:=ioutil.ReadAll(resp.Body)//response body is converrted to byte slice using ioutil.ReadAll and then decoded into result struct using json.unmarshall
	if err!=nil{
		return nil,err
	}
	if resp.StatusCode!=http.StatusOK{
		return nil,fmt.Errorf(string(body))
	}
	res:=&Results{}
	return res,json.Unmarshal(body,res)//unmarshall parses json encoding and stores them into a struct
	
}