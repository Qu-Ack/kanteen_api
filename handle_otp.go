package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"

	"github.com/Qu-Ack/kanteen_api/internal/database"
	"github.com/google/uuid"
)

type MessageResponse struct {
	Status  string `json:"Status"`
	Details string `json:"Details"`
	Otp     string `json:"OTP"`
}

func GenerateOTP(length int) (string, error) {
	const charset = "0123456789"

	result := make([]byte, length)
	for i := range result {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}

func (apiconfig apiConfig) HandleCreateOTP(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Mobile string `json:"mobile"`
	}

	otp, err := GenerateOTP(4)

	if err != nil {
		log.Println("Error In HandleCreateOTP while generating otp", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	byte_body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error In HandleCreateOTP while reading body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	json_body := body{}

	err = json.Unmarshal(byte_body, &json_body)
	if err != nil {
		log.Println("Error In HandleCreateOTP while unmarsheling body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return

	}

	our_url := fmt.Sprintf("https://2factor.in/API/V1/1b095c83-7b14-11ef-8b17-0200cd936042/SMS/+91%v/%v/Standart", json_body.Mobile, otp)

	resp, err := http.Get(our_url)
	if err != nil {
		log.Println("Error In HandleCreateOTP while sending OTP to user", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	defer resp.Body.Close()
	response_body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error In HandleCreateOTP while reading body of message response", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	messageresp_body := MessageResponse{}

	err = json.Unmarshal(response_body, &messageresp_body)
	if err != nil {
		log.Println("Error In HandleCreateOTP while unmarsheling message response body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	if messageresp_body.Status == "failure" {
		log.Println("Error In HandleCreateOTP while sending OTP", messageresp_body)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	Id, err := uuid.NewUUID()
	if err != nil {
		log.Println("Error In HandleCreateOTP while generating UUID", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	_, err = apiconfig.DB.CreateOTP(r.Context(), database.CreateOTPParams{
		ID:     Id,
		Mobile: json_body.Mobile,
		Otp:    otp,
	})

	if err != nil {
		log.Println("Error In HandleCreateOTP while storing log session in DB")
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	WriteJSON(w, 201, map[string]uuid.UUID{
		"sesid": Id,
	})

}

func (apiconfig apiConfig) HandleVerifyOTP(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Sesid string `json:"sesid"`
		Otp   string `json:"otp"`
		Name  string `json:"name"`
	}

	byte_body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error In HandleVerifyOTP while reading body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	json_body := body{}
	err = json.Unmarshal(byte_body, &json_body)

	if err != nil {
		log.Println("Error In HandleVerifyOTP while unmarshelling body", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}

	uuid_id, err := uuid.Parse(json_body.Sesid)
	log.Println(json_body.Otp)
	log.Println(uuid_id)
	log.Println(json_body.Name)

	if err != nil {
		log.Println("Error In HandleVerifyOTP while parsing uuid", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}
	item, err := apiconfig.DB.GetOTP(r.Context(), uuid_id)
	if err != nil {
		log.Println("Error In HandleVerifyOTP while getting item from db", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return
	}
	log.Println(item.Otp)

	if item.Otp != json_body.Otp {
		log.Println("Error In HandleVerify OTP Wrong OTP by user")
		WriteJSONError(w, 401, "OTP verification failed")
		return
	}

	mobile, err := apiconfig.DB.GetOTP(r.Context(), uuid_id)

	user, err := apiconfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name:  json_body.Name,
		Phone: mobile.Mobile,
	})

	if err != nil {
		log.Println("Error In HandleGetOTP while getting things from db", err)
		WriteJSONError(w, 500, "Internal Server Error")
		return

	}

	WriteJSON(w, 200, map[string]string{"status": "success", "user_id": user.ID.UUID.String()})

	// we need to get the otp and the number that was stored in our server
	// we need to verify the match the otp sent by the client with the server's otp
	// then we need to send the response back to the user
}
