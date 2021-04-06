package app

import (
	"encoding/json"
	"net/http"

	"github.com/ashishjuyal/banking-lib/logger"

	"github.com/slns/banking-auth/dto"
	"github.com/slns/banking-auth/service"
)

type AuthHandler struct {
	service service.AuthService
}

func (h AuthHandler) NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	// // writeResponse(w, http.StatusOK, "Handler not implemented...")
	// var registerRequest dto.LoginRequest
	// if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
	// 	logger.Error("Error while decoding login request: " + err.Error())
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// } else {
	// 	// Salt and hash the password using the bcrypt algorithm
	// 	// The second argument is the cost of hashing, which we arbitrarily set as 8 
	// 	// (this value can be more or less, depending on the computing power you wish to utilize)
	// 	hashedPassword, appErr := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), 8)
		
	// 	//_, appErr := h.service.Login(registerRequest)
	// 	if _, err = db.Query("insert into users values ($1, $2)", registerRequest.Username, string(hashedPassword)); appErr != nil {
	// 		// If there is any issue with inserting into the database, return a 500 error
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	} else {
	// 		writeResponse(w, http.StatusOK, registerRequest)
	// 	}
	// }
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		logger.Error("Error while decoding login request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := h.service.Login(loginRequest)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, *token)
		}
	}
}

/*
  Sample URL string
 http://localhost:8181/auth/verify?token=somevalidtokenstring&routeName=GetCustomer&customer_id=2000&account_id=95470
*/
func (h AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	urlParams := make(map[string]string)

	// converting from Query to map type
	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}

	if urlParams["token"] != "" {
		appErr := h.service.Verify(urlParams)
		if appErr != nil {
			writeResponse(w, appErr.Code, notAuthorizedResponse(appErr.Message))
		} else {
			writeResponse(w, http.StatusOK, authorizedResponse())
		}
	} else {
		writeResponse(w, http.StatusForbidden, notAuthorizedResponse("missing token"))
	}
}

func (h AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshRequest dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshRequest); err != nil {
		logger.Error("Error while decoding refresh token request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := h.service.Refresh(refreshRequest)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, *token)
		}
	}
}

func notAuthorizedResponse(msg string) map[string]interface{} {
	return map[string]interface{}{
		"isAuthorized": false,
		"message":      msg,
	}
}

func authorizedResponse() map[string]bool {
	return map[string]bool{"isAuthorized": true}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
