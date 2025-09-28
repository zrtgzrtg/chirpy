package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/zrtgzrtg/chirpy/internal/auth"
	"github.com/zrtgzrtg/chirpy/internal/database"
)

func handlerReady(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))

}
func handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", int(apiCfg.fileserverHits.Load()))))
}
func handlerReset(w http.ResponseWriter, req *http.Request) {
	apiCfg.reset(w, req)
	apiCfg.db.DeleteUsers(context.Background())
}
func handlerValidate(w http.ResponseWriter, rb requestBody) ([]byte, bool) {
	if utf8.RuneCountInString(rb.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return []byte{}, false
	} else {

		cleanBody := filterBadWords([]byte(rb.Body))
		return cleanBody, true
	}
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	respStruct, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(respStruct)
}
func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	errStruct := jsonError{msg}
	resp, _ := json.Marshal(&errStruct)
	w.Write(resp)
}
func filterBadWords(body []byte) []byte {
	badWords := map[string]string{
		"kerfuffle": "****",
		"sharbert":  "****",
		"fornax":    "****",
	}
	stringBody := string(body)
	stringList := strings.Split(stringBody, " ")

	for i, str := range stringList {
		lowerStr := strings.ToLower(str)
		val, ok := badWords[lowerStr]
		if ok {
			stringList[i] = val
		}
	}
	return []byte(strings.Join(stringList, " "))

}

func handlerUser(w http.ResponseWriter, req *http.Request) {
	//if apiCfg.platform != "dev" {
	//respondWithError(w, http.StatusForbidden, "Not on the right platform")
	//return
	//}
	var usr User
	err := json.NewDecoder(req.Body).Decode(&usr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Decode error")
		return
	}
	if usr.Email == "" {
		respondWithError(w, http.StatusBadRequest, "No email field in request")
		return
	}
	hashPass, err := auth.HashPassword(usr.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "password hashing went wrong")
		return
	}
	crtParams := database.CreateUserParams{
		Email:          usr.Email,
		HashedPassword: hashPass,
	}
	retUser, err := apiCfg.db.CreateUser(context.Background(), crtParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	myUsr := UserResponse{
		ID:        retUser.ID,
		CreatedAt: retUser.CreatedAt,
		UpdatedAt: retUser.UpdatedAt,
		Email:     retUser.Email,
	}
	jsonResp, err := json.Marshal(&myUsr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "decode error retUser")
		return
	}
	w.WriteHeader(201)
	w.Write(jsonResp)
}

func mapToUser(dbUsr database.User) User {
	return User{
		ID:        dbUsr.ID,
		CreatedAt: dbUsr.CreatedAt,
		UpdatedAt: dbUsr.UpdatedAt,
		Email:     dbUsr.Email,
	}
}
func mapToChirp(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
}

func handlerPostChirp(w http.ResponseWriter, req *http.Request) {
	var reqParams requestChirpBody
	err := json.NewDecoder(req.Body).Decode(&reqParams)
	if err != nil {
		log.Printf("Decode error: %T: %v", err, err)
		respondWithError(w, http.StatusInternalServerError, "decoding req.body went wrong")
		return
	}
	cleanBody, ok := handlerValidate(w, requestBody{reqParams.Body})
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "validate wrong")
	}
	uid, err := uuid.Parse(reqParams.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "uuid was wrong")
		return
	}
	dbusr, err := apiCfg.db.GetUser(context.Background(), uid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "db request for usr went wrong")
		return
	}
	usr := mapToUser(dbusr)
	createChirpPar := database.CreateChirpParams{
		Body:   string(cleanBody),
		UserID: usr.ID,
	}
	chirpResp, err := apiCfg.db.CreateChirp(context.Background(), createChirpPar)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "creating chirp in db went wrong!")
		return
	}
	chirp := mapToChirp(chirpResp)
	jsonChirp, err := json.Marshal(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "chirp json marshal went wrong")
		return
	}
	w.WriteHeader(201)
	w.Write(jsonChirp)
}
func handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := apiCfg.db.GetChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "chirp request in db went wrong")
		return
	}
	chirps := []Chirp{}
	for _, chirp := range dbChirps {
		chirps = append(chirps, mapToChirp(chirp))
	}
	chirpsStr, err := json.Marshal(&chirps)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "marshalling chirps went wrong")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(chirpsStr)
}

func handlerGetChirp(w http.ResponseWriter, req *http.Request) {

	id := req.PathValue("chirpID")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "no id found. Should probably not get called. This should be equal to the get all chirps endpoint")
		return
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "parsing uuid from request went wrong")
		return
	}
	dbChirp, err := apiCfg.db.GetChirpById(context.Background(), uid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	chirp := mapToChirp(dbChirp)
	chirpStr, err := json.Marshal(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "marshall for chirp went wrong")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(chirpStr)
}

func handlerLogin(w http.ResponseWriter, req *http.Request) {
	loginRequest := LoginRequest{}
	err := json.NewDecoder(req.Body).Decode(&loginRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "decoding body went wrong")
		return
	}
	usr, err := apiCfg.db.GetUserByEmail(context.Background(), loginRequest.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	_, err = auth.CheckPasswordHash(loginRequest.Password, usr.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	usrResp := UserResponse{
		ID:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email:     usr.Email,
	}
	w.WriteHeader(http.StatusOK)
	usrRespStr, err := json.Marshal(&usrResp)
	w.Write(usrRespStr)
}
