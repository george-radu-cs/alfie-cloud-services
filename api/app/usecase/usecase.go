package usecase

import (
	"api/app/repository"
	"api/app/services"
)

type useCase struct {
	Repository          repository.Repository
	MailVerifierService services.MailVerifierService
	ValidationsService  services.ValidationsService
	MediaCloudService   services.MediaCloudService
}

func New(
	repository repository.Repository,
	mailVerifierService services.MailVerifierService,
	validationsService services.ValidationsService,
	mediaCloudService services.MediaCloudService,
) UseCase {
	return &useCase{
		Repository:          repository,
		MailVerifierService: mailVerifierService,
		ValidationsService:  validationsService,
		MediaCloudService:   mediaCloudService,
	}
}
