package infrastructure
import "entity"
//UserRepository to manage users persistence
type UserRepository struct { 
    psqlClient psql.PsqlClient
}
func NewUserRepository(psqlClient psql.PsqlClient) UserRepository 
{
return &UserRepository{psqlClient}}
//Save user
func (selfPtr *UserRepository)Save(user entity.User) error { 
 ... 
}
//GetByEmail gets the user by email
func (selfPtr *UserRepository)GetByEmail(email string) error { 
 ... 
}