package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	L "userTransactions/logging"
	"userTransactions/transactions"

	"github.com/IBM/sarama"
)

func (s *Service) createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		L.Logger.Errorf("Request method noot allowed: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Request method unallowed: %s", r.Method)))
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var userData createUserRequest
	err := d.Decode(&userData)
	if err != nil {
		L.Logger.Errorf("Failed decoding create user request: %+v. %s", r.Body, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed decoding create user request. %s", err.Error())))
		return
	}

	if userData.Email == "" {
		L.Logger.Info("Email is mandatory in request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email is mandatory in request"))
		return
	}

	id, err := s.repo.createUser(userData.Email)
	if err != nil {
		L.Logger.Info("Failed creating new user: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed creating new user. %s", err.Error())))
		return
	}

	// KAFKA notify
	newUserNotification := newUserKafkaMessage{
		Id:        id,
		Email:     userData.Email,
		CreatedAt: time.Now(),
	}

	msgJSON, err := json.Marshal(newUserNotification)
	if err != nil {
		L.Logger.Info("Error encoding newUserKafkaMessage struct to JSON: ", err)
	}

	msg := sarama.ProducerMessage{
		Topic: s.topic,
		Value: sarama.StringEncoder(msgJSON),
	}

	_, _, err = s.producer.SendMessage(&msg)
	if err != nil {
		L.Logger.Error("Failed sending kafka message", msg.Value, ": ", err)
	}

	L.Logger.Infof("New user %s created", userData.Email)
	w.Write([]byte("New user created"))
}

func (s *Service) getUserBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		L.Logger.Errorf("Request method not allowed: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Request method unallowed: %s", r.Method)))
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var userData getBalanceRequest
	err := d.Decode(&userData)
	if err != nil {
		L.Logger.Errorf("Failed decoding get balance request: %+v. %s", r.Body, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed decoding get balance request. %s", err.Error())))
		return
	}

	id, err := s.repo.getUserIDFromEmail(userData.Email)

	if err != nil {
		L.Logger.Error("Failed getting userID based on the given e-mail: ", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Failed getting userID based on the given e-mail: %s", err.Error())))
		return
	}

	/*
		NATS request to transactions service
		transctions service processing...
		NATS response from transactions service
	*/

	natsResp, err := s.nc.Request(s.subject+fmt.Sprintf("%d", id), nil, 10*time.Millisecond)
	if err != nil {
		L.Logger.Errorf("Failed getting balance for user %s: %s", userData.Email, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed getting balance for user %s: %s", userData.Email, err.Error())))
		return
	}

	if natsResp == nil {
		L.Logger.Error("NATS request timed out")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("NATS request timed out"))
		return
	}

	var balance transactions.BalanceNATSResponse

	err = json.Unmarshal(natsResp.Data, &balance)
	if err != nil {
		L.Logger.Error("Error decoding response: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error decoding NATS response: %v", err)))
		return
	}

	L.Logger.Infof("Balance retreived %+v", balance)

	resp := getBalanceResponse{
		Email:   userData.Email,
		Balance: balance.Balance,
	}

	msgJSON, err := json.Marshal(resp)
	if err != nil {
		L.Logger.Error("Error encoding getBalanceResponse struct to JSON: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint("Error encoding getBalanceResponse struct to JSON: ", err)))
		return
	}

	w.Write([]byte(msgJSON))
}
