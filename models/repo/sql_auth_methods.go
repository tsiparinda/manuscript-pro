package repo

import (
	"brucheion/models"
	"database/sql"
	"strings"
)

//functions for working with user from auth module

func (repo *SqlRepository) GetGroupById(id int) (result models.Group) {
	//id, username, email, password_hash, is_verified, verification_code
	err := repo.Commands.GetGroupById.QueryRow(id).Scan(&result.Id, &result.Name)
	switch {
	case err == sql.ErrNoRows:
		repo.Logger.Panicf("GetGroupById: Group with id=%v not found!", id)
		return
	case err != nil:
		repo.Logger.Panicf("GetGroupById: %v", err.Error())
		return
	default:
		repo.Logger.Debugf("GetGroupById Group : %v", result)
		return
	}
}
func (repo *SqlRepository) GetGroupByName(n string) (result models.Group) {
	//id, username, email, password_hash, is_verified, verification_code
	repo.Debugf("GetGroupByName in: @%v@", n)
	err := repo.Commands.GetGroupByName.QueryRow(n).Scan(&result.Id, &result.Name)
	switch {
	case err == sql.ErrNoRows:
		repo.Logger.Panicf("GetGroupByName: Group with id=%v not found!", n)
		return
	case err != nil:
		repo.Logger.Panicf("GetGroupByName: %v", err.Error())
		return
	default:
		repo.Logger.Debugf("GetGroupByName Group : %v", result)
		return
	}
}

func (repo *SqlRepository) AddUserToGroup(creds *models.Credentials, id int) error {
	repo.Logger.Debugf("AddUserToGroup input: %v", creds)
	_, err := repo.Commands.AddUserToGroup.Exec(creds.Id, id)
	if err != nil {
		repo.Logger.Panicf("Cannot exec AddUserToGroup command: %v", err.Error())
		return err
	}

	group := repo.GetGroupById(id)
	creds.Roles = append(creds.Roles, group)

	return nil
}

func (repo *SqlRepository) AddNewUser(creds models.Credentials) (models.Credentials, error) {
	lastInsertId := 0
	err := repo.Commands.AddNewUser.QueryRow(creds.Username, strings.ToLower(creds.Email), creds.VerificationCode,
		creds.Password).Scan(&lastInsertId)
	if err != nil {
		repo.Logger.Panicf("Cannot exec AddNewUser command: %v", err.Error())
		return models.Credentials{}, err
	}
	repo.Logger.Debugf("New user added with ID: %v", lastInsertId)
	creds.Id = lastInsertId
	// // add new user to AllUsers group
	// if err := repo.AddUserToGroup(&creds, 2); err != nil {
	// 	return models.Credentials{}, err
	// }
	// // add new user to ToolsUsers group
	// if err := repo.AddUserToGroup(&creds, 3); err != nil {
	// 	return models.Credentials{}, err
	// }
	return creds, nil
}

func (repo *SqlRepository) AddUserAdmin(creds models.Credentials) (models.Credentials, error) {
	_, err := repo.Commands.AddUserAdmin.Exec(creds.Id, creds.Username, strings.ToLower(creds.Email), creds.Password)
	if err != nil {
		repo.Logger.Panicf("Cannot exec AddUserAdmin command: %v", err.Error())
		return models.Credentials{}, err
	}
	// add new user to AllUsers group
	if err := repo.AddUserToGroup(&creds, 1); err != nil {
		return models.Credentials{}, err
	}
	// add new user to ToolsUsers group
	if err := repo.AddUserToGroup(&creds, 2); err != nil {
		return models.Credentials{}, err
	}
	repo.Logger.Debug("Admin user added to DB")
	return creds, nil
}

func (repo *SqlRepository) UpdateUser(creds *models.Credentials) error {
	_, err := repo.Commands.UpdateUser.Exec(creds.Id, creds.Username, creds.Email, creds.IsVerified, creds.VerificationCode, creds.Password)
	if err != nil {
		repo.Logger.Panicf("UpdateUser Cannot exec UpdateUser command: %v", err.Error())
		return err
	}
	repo.Logger.Debugf("UpdateUser User with ID %v updated", creds.GetID())
	return nil
}

func (repo *SqlRepository) UpdateUserGroups(creds *models.Credentials) error {
	repo.Logger.Debugf("UpdateUserGroups in:", creds)
	// in creds.Roles we have []Group slice of allowed user's groups

	contains := func(g []models.Group, id int) bool {
		for i := range g {
			if g[i].Id == id {
				return true
			}
		}
		return false
	}
	olduser := repo.GetUserByID(creds.Id)
	// delete users groups, which found, but not allowed
	for i := range olduser.Roles {
		if !contains(creds.Roles, olduser.Roles[i].Id) {
			_, err := repo.Commands.DeleteUserGroup.Exec(creds.Id, olduser.Roles[i].Id)
			if err != nil {
				repo.Logger.Panicf("UpdateUserGroups: Cannot exec DeleteUserGroup command: %v", err.Error())
				return err
			}
		}
	}

	// insert new allowed groups
	for i := range creds.Roles {
		if !contains(olduser.Roles, creds.Roles[i].Id) {
			_, err := repo.Commands.AddUserToGroup.Exec(creds.Id, creds.Roles[i].Id)
			if err != nil {
				repo.Logger.Panicf("UpdateUserGroups: Cannot exec AddUserToGroup command: %v", err.Error())
				return err
			}
		}
	}

	repo.Logger.Debugf("UpdateUserGroups User groups with ID %v updated", creds.Id)
	return nil
}

func (repo *SqlRepository) UpdateGroup(group *models.Group) error {
	_, err := repo.Commands.UpdateGroup.Exec(group.Id, group.Name)
	if err != nil {
		repo.Logger.Panicf("UpdateGroup error: %v", err.Error())
		return err
	}
	repo.Logger.Debugf("UpdateGroup Group with ID %v updated", group.Id)
	return nil
}

func (repo *SqlRepository) GetUserByVerificationCode(code string) (models.Credentials, error) {
	var user models.Credentials
	err := repo.Commands.GetUserByVerificationCode.QueryRow(code).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.IsVerified, &user.VerificationCode)
	switch {
	case err == sql.ErrNoRows:
		return models.Credentials{}, nil
	case err != nil:
		repo.Logger.Panicf("An error arise by exec GetUserByVerificationCode command: %v", err.Error())
		return models.Credentials{}, err
	default:
		repo.Logger.Debugf("New user found by verification code: %v", user)
		return user, nil
	}
}

func (repo *SqlRepository) GetUserByEmail(email string) (models.Credentials, error) {
	var user models.Credentials
	//id, username, email, password_hash, is_verified, verification_code
	err := repo.Commands.GetUserByEmail.QueryRow(email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.IsVerified, &user.VerificationCode)
	switch {
	case err == sql.ErrNoRows:
		return models.Credentials{}, nil
	case err != nil:
		return models.Credentials{}, err
	default:
		// if err = repo.GetUserGroups(&user); err != nil {
		// 	repo.Logger.Panicf("Error in GetUserGroups: %v", err.Error())
		// 	return models.Credentials{}, err
		// }
		repo.Logger.Debugf("GetUserByEmail User found by email: %v", user)
		return user, nil
	}
}

func (repo *SqlRepository) GetUserByID(userid int) (result models.Credentials) {
	rows, err := repo.Commands.GetUserByID.QueryContext(repo.Context, userid)
	if err == nil {
		if userline, err := scanUser(rows); err != nil {
			repo.Logger.Panicf("Cannot scan data: %v", err.Error())
			return
		} else {
			if err = repo.GetUserGroups(&userline); err != nil {
				repo.Logger.Panicf("Error in GetUserGroups: %v", err.Error())
				return
			}
			result = userline
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec GetUser command: %v", err.Error())
	}
	return
}

func (repo *SqlRepository) GetUserByName(username string) (result models.Credentials) {

	rows, err := repo.Commands.GetUserByName.QueryContext(repo.Context, username)
	if err == nil {
		if userline, err := scanUser(rows); err != nil {
			repo.Logger.Panicf("Cannot scan data: %v", err.Error())
			return
		} else {
			if err = repo.GetUserGroups(&userline); err != nil {
				repo.Logger.Panicf("Error in GetUserGroups: %v", err.Error())
				return
			}
			result = userline
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec GetUserByName command: %v", err)
	}
	repo.Logger.Debugf("User:", result)
	return
}

func (repo *SqlRepository) GetUserGroups(creds *models.Credentials) error {

	//	usergroups := make([]models.Group, 0, 10)
	rows, err := repo.Commands.GetUserGroups.QueryContext(repo.Context, creds.Id)
	if err == nil {
		for rows.Next() {
			g := models.Group{}
			if err := rows.Scan(&g.Id, &g.Name); err != nil {
				repo.Logger.Debugf("GetUserGroups User has not attach to any group")
				return nil
			}
			//usergroups = append(usergroups, g)
			creds.Roles = append(creds.Roles, g)
		}
		//	repo.Logger.Debugf("GetUserGroups creds", creds)
		return nil
	} else {
		repo.Logger.Debugf("GetUserGroups err", err.Error())
	}
	return err
}

func (repo *SqlRepository) GetCredentials() []models.Credentials {
	var creds []models.Credentials
	//id, username, email, password_hash, is_verified, verification_code
	rows, err := repo.Commands.GetCredentials.QueryContext(repo.Context)
	if err == nil {
		for rows.Next() {
			var cred models.Credentials
			if err := rows.Scan(&cred.Id, &cred.Username, &cred.Email, &cred.IsVerified, &cred.Description); err != nil {
				repo.Logger.Debugf("GetCredentials errors", err.Error())
				return []models.Credentials{}
			}
			repo.GetUserGroups(&cred)
			creds = append(creds, cred)
		}
		return creds
	} else {
		repo.Logger.Debugf("GetCredentials err", err.Error())
	}
	return []models.Credentials{}
}

func (repo *SqlRepository) GetGroups() []models.Group {
	var groups []models.Group
	rows, err := repo.Commands.GetGroups.QueryContext(repo.Context)
	if err == nil {
		var group models.Group
		for rows.Next() {
			if err := rows.Scan(&group.Id, &group.Name); err != nil {
				repo.Logger.Debugf("GetGroups errors", err.Error())
				return groups
			}
			groups = append(groups, group)
		}
		return groups
	} else {
		repo.Logger.Debugf("GetGroups err", err.Error())
	}
	return groups
}
