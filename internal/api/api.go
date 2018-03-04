package api

import (
  "internal/blockchain"
  "encoding/json"
  "net/http"
  "log"
  "github.com/gorilla/mux"
)

//JSONResponse returns JSON response body given an http.ResponseWriter and struct.
func JSONResponse(w http.ResponseWriter, data interface{}) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(data)
}

func StartNode() {
  router := mux.NewRouter()
  // instantiate shared Blockchain for this node
  BC := blockchain.NewBlockchain()

  // Identify API routes
  router.HandleFunc("/chain", func(w http.ResponseWriter, r *http.Request) {
    getChain(w, r, BC)
    }).Methods("GET")

  router.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request) {
    mineBlock(w, r, BC)
    }).Methods("GET")

  router.HandleFunc("/transactions/new", func(w http.ResponseWriter, r *http.Request) {
    postTransaction(w, r, BC)
    }).Methods("POST")

  // Start Server
  log.Print("Starting BasicCoin Node on port 8000 ....")
  http.Handle("/", router)
  log.Fatal(http.ListenAndServe(":8000", router))
}

func mine() {

}

//Handler for GET /chain
func getChain(w http.ResponseWriter, r *http.Request, BC *blockchain.Blockchain) {
  JSONResponse(w, BC)
}

//Handler for GET /mine. JSON response includes new block added to blockchain
func mineBlock(w http.ResponseWriter, r *http.Request, BC *blockchain.Blockchain) {
  BC.Mine()
  nb := BC.LastBlock()
  log.Printf("##### New Block %v created. #####", *nb.BlockHash)
  JSONResponse(w, nb)
}

//Handler for POST /transaction/new
func postTransaction(w http.ResponseWriter, r *http.Request, BC *blockchain.Blockchain) {
  decoder := json.NewDecoder(r.Body)

  // parse request body as intermediate Transaction object
  var _iTx blockchain.Transaction
  err := decoder.Decode(&_iTx)
  if err != nil {
      panic(err)
  }
  defer r.Body.Close()
  BC.NewTransaction(_iTx.Sender, _iTx.Recipient, _iTx.Amount)
  newTX := BC.CurrentTX[len(BC.CurrentTX)-1]
  log.Printf("New Transaction %v created.", newTX.ID)
  JSONResponse(w, newTX)
}
