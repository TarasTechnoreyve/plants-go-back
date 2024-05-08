package controllers

import (
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
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
	}
}
