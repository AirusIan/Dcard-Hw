package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type srv struct{}
type Page struct {
	Articles    string `json:"article"`
	NextPagekey string `json:"nextPagekey"`
}

func (h srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("some-list-key") // get URL param with key "some-list-key"
	var keytoStr string
	keytoStr = fmt.Sprintf("%v", key) // 將拿到的list-key轉為string
	//fmt.Print(keytoStr)
	k := GetHead(keytoStr) // 藉由list-key得到第一個Page key
	//fmt.Println(k)
	p := GetPage(k)
	for p.Articles != "" { // 由GetPage一直拿到下一個Page 直到沒有下一個Page為止  並將Page寫到 Response中
		//fmt.Println(p)
		jsonbytes, _ := json.Marshal(p)
		w.Write(jsonbytes)
		if p.NextPagekey == "" {
			break
		}
		p = GetPage(p.NextPagekey)

	}
}
func GetHead(s string) string {
	db, err := sql.Open("postgres", "user=postgres password=1234 dbname=test sslmode=disable") //建立DB連線
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var nextPagekey string
	err = db.QueryRow("SELECT nextPagekey FROM head WHERE list_key=$1", s).Scan(&nextPagekey) //取得第一個pageKey
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Print(nextPagekey)
	return nextPagekey
}
func GetPage(p string) Page {
	page := Page{}
	var info []byte
	db, err := sql.Open("postgres", "user=postgres password=1234 dbname=test sslmode=disable") //建立DB連線
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.QueryRow("SELECT info FROM users WHERE pgkey=$1", p).Scan(&info)
	if err != nil {
		panic(err)
	}
	err2 := json.Unmarshal(info, &page)
	if err2 != nil {
		panic(err2)
	}
	//fmt.Print(page)
	return page

}
func SetHead(lt_key string, nPgkey string) {
	db, err := sql.Open("postgres", "user=postgres password=1234 dbname=test sslmode=disable") //建立DB連線
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//建立list-key, nextPagekey的值組
	_, err = db.Exec("INSERT INTO head(list_key, nextPagekey) VALUES($1, $2)", lt_key, nPgkey)
	if err != nil {
		panic(err)
	}
}
func SetNextkey(target string, next string) {
	db, err := sql.Open("postgres", "user=postgres password=1234 dbname=test sslmode=disable") //建立DB連線
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//更新指定Page的nextPagekey
	_, err = db.Exec(`UPDATE users SET info = jsonb_set("info", '{"nextPagekey"}', to_jsonb($1::text), true)WHERE pgkey = $2`, next, target)
	if err != nil {
		panic(err)
	}

	//fmt.Println("JSON value updated successfully.")
}
func SetArticle(target string, artic string) {
	db, err := sql.Open("postgres", "user=postgres password=1234 dbname=test sslmode=disable") //建立DB連線
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//更新指定Page的nextPagekey
	_, err = db.Exec(`UPDATE users SET info = jsonb_set("info", '{"article"}', to_jsonb($1::text), true)WHERE pgkey = $2`, artic, target)
	if err != nil {
		panic(err)
	}

	//fmt.Println("JSON value updated successfully.")
}
func SetPage(pagekey string) {
	db, err := sql.Open("postgres", "user=postgres password=1234 dbname=test sslmode=disable") //建立DB連線
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//建立list-key, nextPagekey的值組
	page := &Page{
		Articles:    "",
		NextPagekey: "",
	}
	data, err := json.Marshal(page)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("INSERT INTO users(pgkey,info) VALUES($1,$2)", pagekey, data)
	if err != nil {
		panic(err)
	}
}
func SetCreateTime() {

}
func ClearList() {

}
func main() {
	/*http.Handle("/", srv{})
	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)*/

	/*SetHead("anotherkey", "ghjk")
	SetNextkey("zxcv", "vvvv")*/
	SetPage("deleteKey")
	SetArticle("deleteKey", "Delete Page")

}
