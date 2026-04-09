package main

import (
	"fmt"

	"github.com/AlexandreJSimon/hexagonal-golang-project/pkg/jwt"
)

func main() {

	jwtT, _ := jwt.New(jwt.Config{SecretKey: "super_secret_key"})

	token, _ := jwtT.Generate("12345")

	fmt.Println(token)

	uid, _ := jwtT.Validate(token)

	fmt.Printf("%+v\n", uid)
}
