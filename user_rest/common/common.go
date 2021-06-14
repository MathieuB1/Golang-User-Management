package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"user_rest/user_rest/context"
)

// Helper for JSON Serializer
func SerializeSender(users interface{}) ([]byte, error) {
	res, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Helper for Creating Hash
func CreateHash(key *string) string {
	hasher := md5.New()
	addSalt := context.GlobalCtx.APP_SALT + *key
	hasher.Write([]byte(addSalt))
	return hex.EncodeToString(hasher.Sum(nil))
}
