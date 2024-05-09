package helpers

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"fmt"
	
)

var HmacSampleSecret = "EHUIQHRUEWRUEWBR"

func GenerateToken(ObjectId  primitive.ObjectID)(string,error){

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id": ObjectId,
		"exp":time.Now().Add(time.Hour).Unix(),
	})
	
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(HmacSampleSecret))
	if err!=nil{
		return "",err
	}

	return tokenString,err
}

func TokenValidation(tokenString string)(jwt.MapClaims,error){

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(HmacSampleSecret), nil
	})

	if err != nil {
		return nil,err
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			err = fmt.Errorf("token has expired")
			return claims,err
		}
		return claims,nil
		
	} else {
		return nil,fmt.Errorf("error while validatating the token")
	}
}