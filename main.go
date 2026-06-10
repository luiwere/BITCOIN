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

func showBlockchainInfo() error {
	var info struct {
		Chain      string  `json:"chain"`
		Blocks     int     `json:"blocks"`
		Difficulty float64 `json:"difficulty"`
	}

	err := rpc("getblockchaininfo", nil, "", &info)
	if err != nil {
		return err
	}

	fmt.Println("Chain:", info.Chain)
	fmt.Println("Blocks:", info.Blocks)
	fmt.Println("Difficulty:", info.Difficulty)

	return nil
}

func showWalletBalance(wallet string) error {
	_ = rpc("loadwallet", []any{wallet}, "", nil)

	var balance float64
	if err := rpc("getbalance", nil, wallet, &balance); err != nil {
		return err
	} 

	fmt.Printf("=== Wallet: %s ===\n", wallet)
	fmt.Printf("Balance: %v BTC\n", balance)
	return nil
}

func listTransactions(wallet string) error {

	_ = rpc ("loadwallet", []any{wallet}, "", nil)

	var txs []struct {
		Category      string  `json:"category"`
		Amount        float64 `json:"amount"`
		TxID          string  `json:"txid"`
		Confirmations int     `json:"confirmations"`
	}

	count := 10

	if err := rpc(
		"listtransactions",
		[]any{"*", count},
		wallet,
		&txs, 
		); err != nil {
		return err
	}

	fmt.Printf("=== Transactions for %s ===\n", wallet)

	for _, tx := range txs {
		dir := "OUT"

		switch tx.Category {
		case "receive", "generate", "immature":
			dir = "IN "
		}

		fmt.Printf(
			"%s %+.8f BTC | %d confs\n",
			dir,
			tx.Amount,
			tx.Confirmations,
		)
		fmt.Printf("\tTXID: %s\n", tx.TxID)
	}

	return nil
}

func decodeTransaction(txid string) error {
	var tx struct {
		Vin []struct {
			Coinbase	string	`json:"coinbase"`
			TxID		string	`json:"txid"`
			Vout		int		`json:"vout"`	
		} `json:"vin"`

		Vout []struct {
			Value			float64 `json:"value"`

			ScriptPubKey	struct {
				Address		string `json:"address"`
			} `json:"scriptPubKey"`
		} `json:"vout"`
	}

	if err := rpc("getrawtransaction", []any{txid, true}, "", &tx); err != nil {
		return err
	}

	fmt.Println("=== TRANSACTION ===")

	for _, vin := range tx.Vin {
		if vin.Coinbase != "" {
			fmt.Println("\tCOINBASE (mining reward)")
		} else {
			fmt.Printf("\tFrom: Tx: %s...\n", vin.TxID)
		}
	}

	for _, vout := range tx.Vout {
		address := "UNKNOWN"

		if len(vout.ScriptPubKey.Address) > 0 {
			address = vout.ScriptPubKey.Address
		}

		fmt.Printf("\t%.8f BTC -> %s\n", vout.Value, address)
	}

	return nil
}

func main() {
	err := showWalletBalance("alice")
	if err != nil {
		fmt.Println("RPC error:", err)
		return
	}

	err = listTransactions("alice")
	if err != nil {
		fmt.Println("RPC error:", err)
		return
	}

	err = showBlockchainInfo()
	if err != nil {
		fmt.Println("RPC error:", err)
		return
	}

	txid := "23c67dd878b5a37a022db235004b220fcf75b9ae089d01ea24221c725de95a54"

	err = decodeTransaction(txid)
	if err != nil {
		fmt.Println("RPC error:", err)
		return
	}
}

