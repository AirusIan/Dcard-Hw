package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

type Users struct {
	Pg_id int    `json:"pg_id"`
	Pgkey string `json:"pgkey"`
	Page  Page   `json:"info"`
}
type Page struct {
	Articles    string `json:"articles"`
	NextPageKey string `json:"nextPagekey,omitempty"`
}

func (p Page) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Page) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &p)
}
func GetHead(s string) {
	http.HandleFunc(s, func(rw http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("some-list-key") // get URL param with key "some-list-key"
		fmt.Fprint(rw, key)

	})

}

func main() {
	db, err := sql.Open("postgres", "user=postgres password=1234 dbname=test sslmode=disable") //建立DB連線
	if err != nil {
		panic(err)
	}
	defer db.Close()

	GetHead("/")
	http.ListenAndServe(":8080", nil)

	pgk := "abce"
	rows, err := db.Query("SELECT * FROM users WHERE pgkey=$1", pgk)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var pg_id int
		var pgkey string
		var info []byte
		// 從查詢結果中取出JSON數據
		err = rows.Scan(&pg_id, &pgkey, &info)
		var obj map[string]interface{}
		err = json.Unmarshal(info, &obj)
		if err != nil {
			panic(err)
		}
		//fmt.Print(pg_id)
		fmt.Print(pgkey)
		fmt.Println("Articles : ", obj["article"], "nextPagekey : ", obj["nextPagekey"])
	}

}
