package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type PlantController struct {
	plantService app.PlantService
}

func NewPlantController(ps app.PlantService) PlantController {
	return PlantController{
		plantService: ps,
	}
}

func (c PlantController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		plant, err := requests.Bind(r, requests.AddPlantRequest{}, domain.Plant{})
		if err != nil {
			log.Printf("PlantController -> Save: %s", err)
			BadRequest(w, err)
			return
		}

		plant.UserId = user.Id
		plant, err = c.plantService.Save(plant)
		if err != nil {
			log.Printf("PlantController -> Save: %s", err)
			InternalServerError(w, err)
			return
		}

		var plantDto resources.PlantDto
		Created(w, plantDto.DomainToDto(plant))
	}
}

func (c PlantController) GetForUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		plants, err := c.plantService.GetForUser(user.Id)
		if err != nil {
			log.Printf("PlantController -> GetForUser: %s", err)
			InternalServerError(w, err)
			return
		}

		var plantsDto resources.PlantsDto
		plantsDto = plantsDto.DomainToDtoCollection(plants)
		Success(w, plantsDto)
	}
}

func (c PlantController) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		plant := r.Context().Value(PlantKey).(domain.Plant)

		if user.Id != plant.UserId {
			err := errors.New("access denied")
			Forbidden(w, err)
			return
		}

		var plantDto resources.PlantDto
		Success(w, plantDto.DomainToDto(plant))
	}
}
