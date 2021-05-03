---
title: HTTP parameter passing in Golang
date: 2018-08-20 17:14:29
tags: [HTTP Golang]
keywords:
description:
---

<div align="center">
{% asset_img main.png %}
</div>

#### Preface

In HTTP API design, how parameters are passed is quite important. Let's see the way how to do it in Golang.

<!-- more -->

#### The Way HTTP do it

There are mainly 3 ways to pass data when HTTP: 

*  encode in url parameters:
```
http://foo.com/login?name=joe&dob=19650203
```
*  In url path, RESTful way:
```
# get all articles under health tag and set date
GET http://foo.com/articles/{health}/{date}
```
*  In JSON body
```javascript
{
    "id" : "123",
    "name" : "frank"
}
```

#### How to pass in Golang

following illustration will be using *github.com/gorilla/mux* to build HTTP server: 

```golang
import "github.com/gorilla/mux"

func NewRouter() http.Handler {
    router := mux.NewRouter()

    router.HandleFunc("/", action.Index).Methods("GET")
    router.HandleFunc("/v1/trip_info/{date}/count", action.GetCabsPickupCountInfo).Methods("POST")
    router.HandleFunc("/v1/trip_info/update_cache", action.UpdateCache).Methods("PUT")
    return router
}
```

##### Pass in URL parameters

pass name and D.O.B like following

```
http://foo.com/login?name=joe&dob=19650203
```

parse data in parameters:  
```golang
func Login(res http.ResponseWriter, req *http.Request) {
    name := req.URL.Query().Get("name")
    dob := req.URL.Query().Get("dob")
}
```

##### Pass in URL path

pass in URL path, following is how to pass date:   

```
/v1/trip_info/{date}/count
```

```golang
func GetCabsPickupCountInfo(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    callID := vars["date"]
}
```

##### Pass in JSON Body

We can also pass parameters in HTTP body using JSON: 

in data.json file
```javascript
{
  "id": "1",
  "date" : "2018-08-11",
  "tags" : ["health", "sports"],
  "title" :  "article 1 changed by frank",
  "body"  :  "article 1 body"
}
```

Following command will invoke article_service/articles/create API, passing json in HTTP body
```
curl -XPOST http://article_service/articles/create -d@data.json
```

using following code to parse data from HTTP JSON body:  

*  Create a construct to unmarshal json
```golang
type Article struct {
    Id string           `json:"id"`
    Title string        `json:"title"`
    Date string         `json:"date"`
    Body string         `json:"body"`
    Tags []string       `json:"tags"`
}
```

*  Parse HTTP Body
```golang
const (
        maxBodyLength = 1000
)
func ArticleUpsert(res http.ResponseWriter, req *http.Request) {
        // in case coming request is too large
        req.Body = http.MaxBytesReader(res, req.Body, maxBodyLength)
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
                log.Error(err)
                return
        }
        var article Article
        err = json.Unmarshal(body, &article)
        if err != nil {
               ...
        }
        ...
}
```


#### Conclusion

In this blog we see 3 different ways to get data in HTTP by Golang; In a lot of browsers, url length may have limit

>  URLs over 2,000 characters will not work in the most popular web browsers. 

So if it is required to pass more data, pass in HTTP body may be a better choice, there is no limit by specification, but server or client have different limitations about it. 


#### Reference

[What is the maximum length of a URL in different browsers?](https://stackoverflow.com/questions/417142/what-is-the-maximum-length-of-a-url-in-different-browsers)  
[What is the size limit of a post request?](https://serverfault.com/questions/151090/is-there-a-maximum-size-for-content-of-an-http-post)  