package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func main() {

	// genkey()
	gentoken()

}

func gentoken() {
	privatePEM, err := ioutil.ReadFile("/home/yoseph/workspace/github.com/yosephbuitrago/garagesale/private.pem")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		log.Fatal(err)
	}

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
	claims := struct {
		jwt.StandardClaims
		Roles []string
	}{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "service project",
			Subject:   "123123123",
			ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Roles: []string{"ADMIN"},
	}

	// This will generate a JWT with the claims embedded in them. The database
	// with need to be configured with the information found in the public key
	// file to validate these claims. Dgraph does not support key rotate at
	// this time.

	method := jwt.GetSigningMethod("RS256")
	token := jwt.NewWithClaims(method, claims)
	token.Header["kid"] = "72635462-2j34bj2b-234bj23"
	str, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// token, err := a.GenerateToken(keyID, claims)
	// if err != nil {
	// 	return errors.Wrap(err, "generating token")
	// }

	fmt.Printf("-----BEGIN TOKEN-----\n%s\n-----END TOKEN-----\n", str)
	// return nil

}

func genkey() {
	// Generate a new private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	// Create a file for the private key information in PEM form.
	privateFile, err := os.Create("private.pem")
	if err != nil {
		log.Fatal(err)
	}

	defer privateFile.Close()

	// Construct a PEM block for the private key.
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Write the private key to the private key file.
	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		log.Fatal(err)
	}

	// Marshal the public key from the private key to PKIX.
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	// Create a file for the public key information in PEM form.
	publicFile, err := os.Create("public.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer privateFile.Close()

	// Construct a PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// Write the public key to the private key file.
	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		log.Fatal(err)
	}

	fmt.Println("private and public key files generated")
}
