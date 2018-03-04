package api

import (
  "internal/blockchain"
  "encoding/json"
  "net/http"
  "log"
  "github.com/gorilla/mux"
)

func StartNode() {
  router := mux.NewRouter()
  //instantiate shared Blockchain for this node
  BC := blockchain.NewBlockchain()

  router.HandleFunc("/chain", func(w http.ResponseWriter, r *http.Request) {
    getChain(w, r, BC)
    }).Methods("GET")
  log.Print("Starting BasicCoin Node on port 8000 ....")
  http.Handle("/", router)
  log.Fatal(http.ListenAndServe(":8000", router))
}

func mine() {

}

//Handler for GET /chain
func getChain(w http.ResponseWriter, r *http.Request, BC *blockchain.Blockchain) {
  json.NewEncoder(w).Encode(BC)
}
