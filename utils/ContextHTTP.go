package utils

import (
	"net/http"
	"strconv"
)

type Values struct {
	M map[string]string
}

func (v Values) Get(key string) string {
	return v.M[key]
}

func GetUserIDFromRequest (r *http.Request) uint {
	id64, _ := strconv.ParseUint(r.Context().Value("context").(Values).Get("user_id"),10, 64)
	return uint(id64)
}
