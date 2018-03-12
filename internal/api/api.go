package api

import (
  "fmt"
  "internal/blockchain"
  "encoding/json"
  "net/http"
  "log"
  "github.com/gorilla/mux"
)

//Alternative response objects
type GenericResponse struct {
  Message string `json:"message"`
}

type Node struct {
  Addr string `json:"address"`
}

type ResolveResponse struct {
  Message string `json:"message"`
  Chain []blockchain.Block `json:"chain"`
}

//JSONResponse returns JSON response body given an http.ResponseWriter and struct.
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(status)
  json.NewEncoder(w).Encode(data)
}

func StartNode(port *string) {
  router := mux.NewRouter()
  // instantiate shared Blockchain for this node
  BC := blockchain.NewBlockchain()

  // Service API routes
  router.HandleFunc("/chain", func(w http.ResponseWriter, r *http.Request) {
    getChain(w, r, BC)
    }).Methods("GET")

  router.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request) {
    mineBlock(w, r, BC)
    }).Methods("GET")

  router.HandleFunc("/transactions/new", func(w http.ResponseWriter, r *http.Request) {
    postTransaction(w, r, BC)
    }).Methods("POST")

  // Service Node Management routes
  router.HandleFunc("/nodes/register", func(w http.ResponseWriter, r *http.Request) {
    registerNode(w, r, BC)
  }).Methods("POST")

  router.HandleFunc("/nodes/resolve", func(w http.ResponseWriter, r *http.Request) {
    attainConsensus(w, r, BC)
  }).Methods("GET")

  // Start Server
  log.Printf("Starting Basic Coin Node on port %v ....", *port)
  http.Handle("/", router)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), router))
}

//Handler for GET /chain
func getChain(w http.ResponseWriter, r *http.Request, BC *blockchain.Blockchain) {
  JSONResponse(w, http.StatusOK, BC)
}

//Handler for GET /mine. JSON response includes new block added to blockchain
func mineBlock(w http.ResponseWriter, r *http.Request, BC *blockchain.Blockchain) {
  if len(BC.CurrentTX) == 0 {
    noCurrentTX := GenericResponse{Message: "No pending transactions."}
    JSONResponse(w, http.StatusOK, noCurrentTX)
  } else {
    BC.Mine()
    nb := BC.LastBlock()
    log.Printf("##### New Block %v created. #####", *nb.BlockHash)
    JSONResponse(w, http.StatusCreated, nb)
  }
}

//Handler for POST /transaction/new
func postTransaction(w http.ResponseWriter, r *http.Request, BC *blockchain.Blockchain) {
  var _iTx blockchain.Transaction
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&_iTx)
  if err != nil {
      panic(err)
  }
  defer r.Body.Close()

  BC.NewTransaction(_iTx.Sender, _iTx.Recipient, _iTx.Amount)
  newTX := BC.CurrentTX[len(BC.CurrentTX)-1]
  log.Printf("New Transaction %v created.", newTX.ID)
  JSONResponse(w, http.StatusCreated, newTX)
}

//Handler for POST /node/register
func registerNode(w http.ResponseWriter, r *http.Request, BC *blockchain.Blockchain) {
  var _n Node
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&_n)
  if err != nil {
      panic(err)
  }
  defer r.Body.Close()

  err = BC.NewNode(_n.Addr)
  if err != nil {
    e := GenericResponse{Message: err.Error()}
    JSONResponse(w, http.StatusBadRequest, e)
  } else {
    log.Printf("New Node %v registered.", _n)
    JSONResponse(w, http.StatusCreated, _n)
  }
}

//Handler for GET /node/resolve
func attainConsensus(w http.ResponseWriter, r *http.Request, BC *blockchain.Blockchain) {
  var resolveResp ResolveResponse
  replaced := BC.ResolveConflicts()
  if replaced {
    resolveResp = ResolveResponse{Message: "Conflicts resolved. Chain replaced.",
      Chain: BC.Chain,
    }
  } else {
    resolveResp = ResolveResponse{Message: "Existing chain persists.",
      Chain: BC.Chain,
    }
  }

  JSONResponse(w, http.StatusOK, resolveResp)
}
