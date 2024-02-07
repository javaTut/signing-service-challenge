package api

import (
	"encoding/json"
	"net/http"
	"signing-service-challenge/services"

	"github.com/gorilla/mux"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress          string
	signatureDeviceService services.SignatureDeviceService
	signatureService       services.SignatureService
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, deviceService services.SignatureDeviceService, signatureService services.SignatureService) *Server {
	return &Server{
		listenAddress: listenAddress,
		// TODO: add services / further dependencies here ...
		signatureDeviceService: deviceService,
		signatureService:       signatureService,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {

	// gorilla mux router used to handle path variables
	router := mux.NewRouter()

	router.HandleFunc("/api/v0/health", s.Health)
	router.HandleFunc("/api/v0/sign", s.Sign).Methods("POST")
	router.HandleFunc("/api/v0/verify", s.Verify).Methods("POST")
	router.HandleFunc("/api/v0/devices", s.CreateDevice).Methods("POST")
	router.HandleFunc("/api/v0/devices", s.GetAllDevices).Methods("GET")
	router.HandleFunc("/api/v0/devices/{id}", s.GetDevice).Methods("GET")
	router.HandleFunc("/api/v0/signatures", s.GetAllSignatures).Methods("GET")
	router.HandleFunc("/api/v0/signatures/{id}", s.GetSignature).Methods("GET")

	return http.ListenAndServe(s.listenAddress, router)
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}
