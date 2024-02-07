package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Handler for retrieving signatures from DB
func (s *Server) GetAllSignatures(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}
	result, err := s.signatureService.GetAll()
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
		return
	}
	WriteAPIResponse(response, http.StatusOK, result)
}

func (s *Server) GetSignature(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}
	vars := mux.Vars(request)
	result, err := s.signatureService.GetById(vars["id"])
	if err != nil {
		WriteErrorResponse(response, http.StatusNotFound, []string{err.Error()})
		return
	}
	WriteAPIResponse(response, http.StatusOK, result)

}
