package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	rpcURL      = "http://127.0.0.1:18443"
	rpcUser     = "luiwere"
	rpcPassword = "Montana"
)

type RPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type rpcResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func rpc(method string, params []any, wallet string, out any) error {
	url := rpcURL

	if wallet != "" {
		url += "/wallet/" + wallet
	}

	body, _ := json.Marshal(RPCRequest{
		JSONRPC: "1.0",
		ID:      "explorer",
		Method:  method,
		Params:  params,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.SetBasicAuth(rpcUser, rpcPassword)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var parsed rpcResponse
	json.NewDecoder(resp.Body).Decode(&parsed)

	if parsed.Error != nil {
		return fmt.Errorf("RPC error: %s", parsed.Error.Message)
	}

	return json.Unmarshal(parsed.Result, out)
}

func main() {
	// ---- Get balance ----
	var balance float64

	err := rpc("getbalance", []any{}, "bob", &balance)
	if err != nil {
		fmt.Println("RPC error:", err)
		return
	}

	fmt.Printf("Bob has %v BTC\n", balance)

	// ---- Get blockchain info ----
	var info struct {
		Chain  string `json:"chain"`
		Blocks int    `json:"blocks"`
	}

	err = rpc("getblockchaininfo", nil, "", &info)
	if err != nil {
		fmt.Println("RPC error:", err)
		return
	}

	fmt.Println("Chain:", info.Chain)
	fmt.Println("Blocks:", info.Blocks)
}
