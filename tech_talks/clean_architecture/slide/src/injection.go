package main
import (
  "usecase"
  "infrastructure" 
  "psqlClient"
)
func main() { 
    psqlClient := psqlClient.NewPsqlClient() 
    userRepository := infrastructure.NewUserRepository(psqlClient)
    useCase := usecase.newUserRegisterUseCase(userRepository)
    //Now we can inject usecase into a delivery layer like web or     
    //cli
}