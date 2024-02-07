package api

// TODO: REST endpoints ...
import (
	"encoding/json"
	"io"
	"net/http"
	"signing-service-challenge/domain"
	"signing-service-challenge/dto"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) CreateDevice(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}
	reqBody, _ := io.ReadAll(request.Body)
	var deviceRequest dto.CreateSignatureDeviceRequest
	err := json.Unmarshal(reqBody, &deviceRequest)
	//handle invalid json
	if err != nil {
		WriteErrorResponse(response, http.StatusUnprocessableEntity, []string{err.Error()})
		return
	}
	valid, err := dto.ValidateCreateSignatureDeviceRequest(deviceRequest)
	if !valid {
		WriteErrorResponse(response, http.StatusBadRequest, []string{err.Error()})
		return
	}
	id := uuid.NewString()
	newDevice, err := s.signatureDeviceService.CreateSignatureDevice(id, deviceRequest.Algorithm, deviceRequest.Label)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
		return
	}
	WriteAPIResponse(response, http.StatusCreated, newDevice)
}

func (s *Server) Sign(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}
	reqBody, _ := io.ReadAll(request.Body)
	var signRequest dto.SignatureRequest
	err := json.Unmarshal(reqBody, &signRequest)
	if err != nil {
		WriteErrorResponse(response, http.StatusUnprocessableEntity, []string{err.Error()})
		return
	}
	valid, err := dto.ValidateSignRequest(signRequest)
	if !valid {
		WriteErrorResponse(response, http.StatusBadRequest, []string{err.Error()})
		return
	}
	signedData, err := s.signatureDeviceService.SignTransaction(signRequest.Id, signRequest.Data)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
		return
	}
	// save new signature in DB
	newSignature := domain.Signature{
		Id:        uuid.NewString(),
		Signature: signedData.Signature,
		Data:      signedData.SignedData,
		SignedBy:  signRequest.Id,
	}
	s.signatureService.Save(newSignature)
	WriteAPIResponse(response, http.StatusAccepted, signedData)
}

func (s *Server) Verify(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}
	reqBody, _ := io.ReadAll(request.Body)
	var verifyRequest dto.VerificationRequest
	err := json.Unmarshal(reqBody, &verifyRequest)
	if err != nil {
		WriteErrorResponse(response, http.StatusUnprocessableEntity, []string{err.Error()})
		return
	}
	valid, err := dto.ValidateVerifyRequest(verifyRequest)
	if !valid {
		WriteErrorResponse(response, http.StatusBadRequest, []string{err.Error()})
		return
	}
	verification, err := s.signatureDeviceService.Verify(
		verifyRequest.DeviceId,
		verifyRequest.Signature,
		verifyRequest.Data,
	)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
		return
	}
	WriteAPIResponse(response, http.StatusAccepted, verification)
}

func (s *Server) GetAllDevices(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}
	result, err := s.signatureDeviceService.GetAll()
	if err != nil {
		WriteAPIResponse(response, http.StatusNotFound, []string{err.Error()})
		return
	}
	WriteAPIResponse(response, http.StatusOK, result)
}

func (s *Server) GetDevice(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}
	vars := mux.Vars(request)
	result, err := s.signatureDeviceService.GetById(vars["id"])
	if err != nil {
		WriteErrorResponse(response, http.StatusNotFound, []string{err.Error()})
		return
	}
	WriteAPIResponse(response, http.StatusOK, result)

}
