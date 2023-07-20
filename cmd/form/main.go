package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	Username string
	Password string
}

func main() {

	s := http.NewServeMux()
	s.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		// user := User{}
		// user.Username = r.FormValue("username")
		// user.Password = r.FormValue("password")

		// buf := new(bytes.Buffer)
		// buf.ReadFrom(r.Body)
		// json.Unmarshal(buf, v any)

		b, _ := io.ReadAll(r.Body)
		type aa struct {
			Name string `json:"name"`
		}
		a := aa{}

		json.Unmarshal(b, &a)
		fmt.Println(a)

		// json.Unmarshal(b, &a)

		w.Write(b)

		// b, err := json.Marshal(user)
		// if err != nil {
		// 	fmt.Fprint(w, err)
		// }

		// w.Write(b)

	})

	http.ListenAndServe("localhost:4000", s)

}
