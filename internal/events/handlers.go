package events

import (
	"fmt"
	"log"
	"net/http"

	"github.com/freitasmatheusrn/social-fit/internal/events/dtos"
	"github.com/freitasmatheusrn/social-fit/internal/events/eventpgs"
	"github.com/freitasmatheusrn/social-fit/internal/types"

	"github.com/freitasmatheusrn/social-fit/pkg/renderer"
	"github.com/freitasmatheusrn/social-fit/pkg/rest"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service ServiceInterface
}

func NewHandler(s ServiceInterface) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) Form(c echo.Context) error {
	return renderer.Render(c, eventpgs.EventForm(nil, eventpgs.FormFields{}), 200)
}

func (h *Handler) Create(c echo.Context) error {
	var r dtos.CreateRequest
	err := c.Bind(&r)
	if err != nil {
		return rest.NewUnprocessableEntity("erro ao processar os dados")
	}
	currentUser := types.GetCurrentUser(c)
	r.UserID = currentUser.ID

	eventID, apiErr := h.Service.Create(c.Request().Context(), r)
	if apiErr != nil {
		log.Println(apiErr)
		data := eventpgs.FormFields{
			Title:         r.Title,
			PosterUrl:     r.PosterURL,
			Description:   r.Description,
			Status:        string(r.Status),
			Start_date:    r.StartDate,
			End_date:      r.EndDate,
			City:          r.City,
			Street:        r.Street,
			Neighbourhood: r.Neighbourhood,
			MaxCapacity:   int(r.MaxCapacity),
		}
		return renderer.Respond(c, eventpgs.Form(apiErr, data), r, apiErr)
	}
	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/dashboard/events/%s", eventID))
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Show(c echo.Context) error {
	eventID := c.Param("event_id")
	currentUser := types.GetCurrentUser(c)
	event, err := h.Service.GetEventWithDetails(c.Request().Context(), eventID, currentUser.ID)
	if err != nil {
		return err
	}

	return renderer.Render(c, eventpgs.Show(event), 200)
}
