package helpers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)


type PasswordM struct{
	Password 		string
	HashedPassword	string
}

func (p *PasswordM) CreatePasswordHash()(string,error){
	byteData,err := bcrypt.GenerateFromPassword([]byte(p.Password),0)
	if err!=nil{
		return "",err
	}

	return string(byteData),nil
}


func(p *PasswordM) ComparePassword()(bool,error){

	if err:= bcrypt.CompareHashAndPassword([]byte(p.HashedPassword),[]byte(p.Password));err!=nil{
		fmt.Println(err)
		return false,err
	}

	return true,nil
}	
