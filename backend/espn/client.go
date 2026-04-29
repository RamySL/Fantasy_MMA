package espn

import (
   "io"
   "log"
   "net/http"
   "encoding/json"
   "fmt"
)

const base_url = "http://site.api.espn.com/apis/site/v2/sports/mma/ufc"

func Fetch() {
   resp, err := http.Get(base_url + "/scoreboard")
   if err != nil {
      log.Fatalln(err)
   }	

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }

   var decoded map[string]interface{}
   json.Unmarshal(body, &decoded)
   fmt.Println(decoded["season"])
}