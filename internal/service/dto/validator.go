package dto

import "test/pkg/validator"

func ValidCreateUserDTO(userDTO CreateUserDTO) bool {
	return validator.ValidEmail(userDTO.Email) && validator.ValidPassword(userDTO.Password)
}

func ValidUpdateUserDTO(userDTO UpdateUserDTO) bool {
	return validator.ValidEmail(userDTO.Email) &&
		validator.ValidPassword(userDTO.Password)
}
