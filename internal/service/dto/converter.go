package dto

import (
	"test/internal/domain"
	"test/pkg/api/params"
)

func ConvertCreateUserDTO(userDTO CreateUserDTO) domain.User {
	return domain.User{
		PasswordHash: userDTO.Password,
		Email:        userDTO.Email,
	}

}

func ConvertUpdateUserDTO(userDTO UpdateUserDTO) (domain.User, error) {
	oid, err := params.ParseIdToObjectID(userDTO.Id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:           oid,
		PasswordHash: userDTO.Password,
		Email:        userDTO.Email,
	}, nil
}
