package repository

import "api/app/models"

func (r *repository) CreateUser(user *models.User) (err error) {
	err = r.Db.Model(&models.User{}).Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateUser(user *models.User) (err error) {
	err = r.Db.Model(&models.User{ID: user.ID}).Select("*").Updates(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetUserByEmail(email string) (user *models.User, err error) {
	user = &models.User{}
	err = r.Db.Model(&models.User{}).Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
