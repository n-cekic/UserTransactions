package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"
	L "userTransactions/logging"
)

func (s *Service) depositHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		L.Logger.Errorf("Request method not allowed: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Request method unallowed: %s", r.Method)))
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var deposit depositRequest
	err := d.Decode(&deposit)
	if err != nil {
		L.Logger.Errorf("Failed decoding deposit request. %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed decoding deposit request. %s", err.Error())))
		return
	}

	newBalance, err := s.repo.deposit(deposit.UserId, deposit.Amount)
	if err != nil {
		L.Logger.Error("Failed creating deposit: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed creating deposit: " + err.Error()))
		return
	}

	resp := depositResponse{Balance: newBalance}
	payload, err := json.Marshal(resp)
	if err != nil {
		L.Logger.Error("Failed creating response: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed creating response: " + err.Error()))
		return
	}
	w.Write(payload)
}

func (s *Service) transferHnadler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		L.Logger.Errorf("Request method not allowed: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Request method unallowed: %s", r.Method)))
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var transferData transferRequest
	err := d.Decode(&transferData)
	if err != nil {
		L.Logger.Errorf("Failed decoding transfer request: %+v. %s", r.Body, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed decoding transfer request. %s", err.Error())))
		return
	}

	err = s.repo.transfer(transferData.FromUserId, transferData.ToUserId, transferData.Amount)
	if err != nil {
		L.Logger.Error("Failed creating transfer: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed creating transfer: " + err.Error()))
		return
	}

	w.Write([]byte("Transfer created"))

}
