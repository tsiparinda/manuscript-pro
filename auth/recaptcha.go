package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const siteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"

type LoginRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	RecaptchaResponse string `json:"g-recaptcha-response"`
}

type SiteVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func CheckRecaptcha(secret, response, action string) (err error) {
	req, err := http.NewRequest(http.MethodPost, siteVerifyURL, nil)
	if err != nil {
		return
	}
	// Add necessary request parameters.
	q := req.URL.Query()
	q.Add("secret", secret)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()
	// Make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Decode response.
	var body SiteVerifyResponse
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return
	}

	// Check recaptcha verification success.
	if !body.Success {
		err = errors.New("unsuccessful recaptcha verify request")
		return
	}

	// NEXT CODE FOR reCaptcha v3!!!!!!!!!!!!!!!!!!!!!

	// Check response score.
	if body.Score < 0.5 {
		err = errors.New(fmt.Sprintf("lower received score than expected (%v)", body))
		// handler.Logger.Debugf("CheckRecaptcha ", err.Error())
		return
	}

	// Check response action.
	if body.Action != action {
		err = errors.New("mismatched recaptcha action")
		return
	}

	return nil
}
