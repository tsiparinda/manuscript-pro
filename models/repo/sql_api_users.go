package repo

import (
	"brucheion/models"
)

func (repo *SqlRepository) GetUsers(userid int) []models.User {
	var users []models.User
	//id, username, email, password_hash, is_verified, verification_code
	rows, err := repo.Commands.GetUsers.QueryContext(repo.Context, userid)
	if err == nil {
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Description, &user.Avatar); err != nil {
				repo.Logger.Debugf("GetUsers errors", err.Error())
				return []models.User{}
			}
			users = append(users, user)
		}
		return users
	} else {
		repo.Logger.Debugf("GetUsers err", err.Error())
	}
	return []models.User{}
}
