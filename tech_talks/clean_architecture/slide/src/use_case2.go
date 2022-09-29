
package usecase
import "entity"
//Hey use case there is out there some user saver that will save the //user
type UserSaver interface {
    Save(user entity.User) error
}
//UserRegisterUseCase will store his dependencies
type UserRegisterUseCase struct { 
    userSaver UserSaver
}
func newUserRegisterUseCase(userSaver UserSaver) UserRegisterUseCase { 
    return &{userSaver}
}
func (u UserRegisterUseCase) RegisterUser(user entity.User) error { 
    user.Validate() 
    err := u.userSaver.Save(user) 
    if (err != nil) { 
       return err 
    }
}