package usecase

import (
	"api/internal/app/alfie/repository"
	services2 "api/internal/app/alfie/services"
)

type useCase struct {
	Repository          repository.Repository
	MailVerifierService services2.MailVerifierService
	ValidationsService  services2.ValidationsService
	MediaCloudService   services2.MediaCloudService
}

func New(
	repository repository.Repository,
	mailVerifierService services2.MailVerifierService,
	validationsService services2.ValidationsService,
	mediaCloudService services2.MediaCloudService,
) UseCase {
	return &useCase{
		Repository:          repository,
		MailVerifierService: mailVerifierService,
		ValidationsService:  validationsService,
		MediaCloudService:   mediaCloudService,
	}
}
