package nsecurity

import (
	"errors"
	"strconv"
	"time"

	"nfgo.ga/nfgo/nerrors"
	util "nfgo.ga/nfgo/nutil"
	"nfgo.ga/nfgo/nutil/jwt"
)

// AuthTicket -
type AuthTicket struct {
	ClientType string
	RequestID  string
	Token      string
	Subject    string
	Timestamp  string
	Signature  string
}

// VerifyToken -
func (a *AuthTicket) VerifyToken(jwtSecret string) error {

	// Check the jwt token
	claims, err := jwt.ValidateToken(jwtSecret, a.Token)
	if err != nil {
		return nerrors.ErrUnauthorized
	}

	// Check the subject in the token
	if claims.Subject != a.Subject {
		return errors.New("the ticket's subject is not equal with the subject int the token")
	}
	return nil
}

// VerifySignature -
func (a *AuthTicket) VerifySignature(signKey string) bool {
	expectSig := util.Sha256(signKey + a.Timestamp + a.Subject + a.RequestID)
	return expectSig == a.Signature
}

// VerifyTimeWindow - check IsoverTimeWindow clientTs milliseconds since January 1, 1970 UTC.
func (a *AuthTicket) VerifyTimeWindow(timeWindow time.Duration) error {
	clientTs, err := strconv.ParseInt(a.Timestamp, 10, 64)
	if err != nil {
		return err
	}

	clientTime := time.Unix(0, clientTs*int64(time.Millisecond))
	serverTime := time.Now()

	var duration time.Duration
	if serverTime.After(clientTime) {
		duration = serverTime.Sub(clientTime)
	} else {
		duration = clientTime.Sub(serverTime)
	}
	if duration > timeWindow {
		return nerrors.ErrUnauthorized
	}
	return nil
}
