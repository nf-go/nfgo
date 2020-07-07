package nsecurity

import (
	"testing"
)

func TestAuthTicketVerifySignature(t *testing.T) {
	signKey := "5f0fa825de41a7d3fd000002"
	ts := "1594800736455"
	subject := "1"
	requestID := "xxx"
	token := "eyJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2MjYzOTc2MDUuNjIyNTY0LCJzdWIiOiIxIn0.4S2w33nYU5rzOjPePP5t4YUnNVmRGJKlUxu8_ioogoU"
	sig := "d3dee5e6260b6a2ce59bc7a6a1a14d024756770cc7a6b5ec9a7ee4a9c5dc82d7"

	ticket := &AuthTicket{
		Timestamp: ts,
		Subject:   subject,
		RequestID: requestID,
		Signature: sig,
		Token:     token,
	}

	if !ticket.VerifySignature(signKey) {
		t.Fatal("error signature")
	}

	if err := ticket.VerifyToken("VoCra8#GEBAbRl*+vos9UF@??gi8Oy"); err != nil {
		t.Fatal(err)
	}

}
