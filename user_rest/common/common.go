package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

// Helper for Serialize
func SerializeAndSendResponse(w *http.ResponseWriter, users interface{}) {

	res, err := json.Marshal(users)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}

	(*w).Header().Set("Content-Type", "application/json")
	(*w).Write(res)
}

func CreateHash(key *string) string {
	hasher := md5.New()
	hasher.Write([]byte(*key))
	return hex.EncodeToString(hasher.Sum(nil))
}
