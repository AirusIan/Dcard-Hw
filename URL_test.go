package main

import (
	"log"
	"testing"

	"database/sql"
	"encoding/json"

	_ "github.com/lib/pq"
)

func TestGetHead(t *testing.T) {
	hd := GetHead("somekey")
	if hd != "abce" {
		t.Errorf("Wrong Head Key")
	}
}
func TestGetPage(t *testing.T) {
	p := GetPage("abce")
	if p.Articles != "Here is Page with key abce" {
		t.Errorf("Wrong Article")
	}
	if p.NextPagekey != "wasd" {
		t.Errorf("Wrong NextPageKey")
	}
}
func TestSetHead(t *testing.T) {
	db, err := sql.Open("postgres", "user=postgres password=1234 dbname=test sslmode=disable") //建立DB連線
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 插入資料以便進行測試
	_, err = db.Exec("INSERT INTO head(list_key, nextpagekey) VALUES($1, $2)", "test-list-key", "test-page-key")
	if err != nil {
		panic(err)
	}
	hd := GetHead("test-list-key")
	if hd != "test-page-key" {
		t.Errorf("Setting Head Error")
	}
}
func TestSetPage(t *testing.T) {
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
	_, err = db.Exec("INSERT INTO users(pgkey,info) VALUES($1,$2)", "test-page-key", data)
	if err != nil {
		panic(err)
	}
	var Pagekey string
	err = db.QueryRow("SELECT pgkey FROM users WHERE pgkey=$1", "test-page-key").Scan(&Pagekey)
	if err != nil {
		log.Fatal(err)
	}
	if Pagekey != "test-page-key" {
		t.Errorf("Setting Page Error")
	}

}
func TestSetArticle(t *testing.T) {
	db, err := sql.Open("postgres", "user=postgres password=1234 dbname=test sslmode=disable") //建立DB連線
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//更新指定Page的nextPagekey
	_, err = db.Exec(`UPDATE users SET info = jsonb_set("info", '{"article"}', to_jsonb($1::text), true)WHERE pgkey = $2`, "Here is Page with key test-page-key", "test-page-key")
	if err != nil {
		panic(err)
	}
	p := GetPage("test-page-key")
	if p.Articles != "Here is Page with key test-page-key" {
		t.Errorf("Setting Article Error")
	}
}
func TestSetNextKey(t *testing.T) {
	db, err := sql.Open("postgres", "user=postgres password=1234 dbname=test sslmode=disable") //建立DB連線
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(`UPDATE users SET info = jsonb_set("info", '{"nextPagekey"}', to_jsonb($1::text), true)WHERE pgkey = $2`, "change-key", "test-page-key")
	if err != nil {
		panic(err)
	}
	p := GetPage("test-page-key")
	if p.NextPagekey != "change-key" {
		t.Errorf("Setting NextPageKey Error")
	}

}
func TestSetCreateTime(t *testing.T) {

}
func TestClearList(t *testing.T) {

}
