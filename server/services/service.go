package services

import "net/http"

type Service interface {
	Handle() http.HandlerFunc
}
