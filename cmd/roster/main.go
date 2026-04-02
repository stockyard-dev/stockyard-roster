package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-roster/internal/server";"github.com/stockyard-dev/stockyard-roster/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="8970"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./roster-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("roster: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Roster — team directory\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("roster: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
