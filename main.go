package main

import (
	"bytes"
	"html/template" //this package is used to generate HTML output that is safe against code injection
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	news "github.com/harshit-456/newsapp/news"
	"github.com/joho/godotenv"
	//"google.golang.org/genproto/googleapis/cloud/aiplatform/v1/schema/predict/params"

	//news "newsapp/news"
)

var tpl=template.Must(template.ParseFiles("index.html"))//this variable points to index html file in root of our project
//struct to represent query made by user
type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *news.Results
}
func(s *Search) IsLastPage() bool{
	return s.NextPage>=s.TotalPages
}
func (s *Search) CurrentPage() int{
	if s.NextPage==1{
		return s.NextPage
	}
	return s.NextPage-1
}
func (s *Search) PreviousPage() int{
	return s.CurrentPage()-1
}

func indexHandler( w http.ResponseWriter,_ *http.Request){
// w is used to send responses to an http request
//r represents http request recieved from client
//	w.Write([]byte("<h1> Hello World </h1>"));
//tpl.Execute(w,nil)//tpl must write output to w and nil data ispassed to tpl,here we are executimg template directly on responsewriter
buf:=&bytes.Buffer{}
err:=tpl.Execute(buf,nil)
if(err!=nil){
	http.Error(w,err.Error(),http.StatusInternalServerError)
	return
}
buf.WriteTo(w)

}
//searchHandler returns an anonymous fucntioon which sstisfies http.HandleFunc type.
//the second arg of HandleFun(patter,func) should have signature of form(w http.ResponseWriter,r *http.Request)
func searchHandler(newsapi *news.Client)http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){

	u,err:=url.Parse(r.URL.String())
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	params:=u.Query()
	searchQuery:=params.Get("q")//q represents users query
	page:=params.Get("page")
	if page==""{
		page="1"
	}
	results,err:=news.FetchEverything(searchQuery,page,newsapi)
	//the above code extracts q and page paramters from request url
	//println writes on server that is here
	//fprintln writes it on client as a response
	//fmt.Println("Search Query is",searchQuery)
	//fmt.Println("Page is",page)
	if(err!=nil){
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	nextpage,err:=strconv.Atoi(page)
	if(err!=nil){
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	//initall the value of nextpage is 1
	search:=&Search{
		Query:searchQuery,
		NextPage: nextpage,
		TotalPages:int(math.Ceil(float64(results.TotalResults) / float64(newsapi.PageSize))),
		Results: results,
	}
	//below we increment value of nextpage,so nextpage now becomes 2(if initially it is 1),so it points at next page
	if check:=!search.IsLastPage(); check{
		search.NextPage++;
	}
buf:=&bytes.Buffer{}
//template is first executed into an empty buffer so that we can check for erros after that buffers is written to response writer
err=tpl.Execute(buf,search)
if(err!=nil){
	http.Error(w,err.Error(),http.StatusInternalServerError)
	return
}
buf.WriteTo(w)
}
}

func main(){
	err:=godotenv.Load()//load method reads env file and sets variables in environment,this is helpful for stroing credentials
	if(err!=nil){
		log.Println("Erro loading env files")
	}

	port:=os.Getenv("PORT")
	if port==""{
		port="3000"
	}
apikey:=os.Getenv("NEWS_API_KEY")
if apikey==""{
	log.Fatal("Env apikey must be set")

}
//var myClient *http.Client
myClient := &http.Client{Timeout: 10 * time.Second}//request timesout after 10sec
	newsapi := news.NewsClient(myClient, apikey, 20)
//fileserver as the name tells serves u with a file in specific folder,this file can be used to handle the HTTP request as required,
fs:=http.FileServer(http.Dir("beauty"))//http file server returns http.Handle type that sserves http request with contents of file "beauty"
//fs points at beauty directory
//string value passed to http.dir can be anywhere on the native(your computers) file system
//even if beauty directory is kept in susnhine folder it will find it
	mux:=http.NewServeMux()//creates an http request mutliplexer
	//Essentially a request mux matches path of the url of incoming HTTP requests against list of registered patterns and calls the associated handler   whenever match occurs
	//below we register path and their associated handlers
	mux.Handle("/beauty/",http.StripPrefix("/beauty/",fs))	//strip prefix strips off /beauty/ part and forward the modififed request to handler returned by http.FileServer,for eg if reuest is /beauty/style.css then it sends style.css to fs and fs will look for style.css inside beauty ,(fs is already pointing the beauty directory) and serve this file to http reuqest /beauty/style.css
		mux.HandleFunc("/search",searchHandler(newsapi))
	mux.HandleFunc("/",indexHandler)//here we register pattern and the handling fucntion to mux.
//ServeMux is a struct and Hzndle and Handleucn are func assoicated to this struct
    http.ListenAndServe(":"+port,mux)//starts server on this port

}
