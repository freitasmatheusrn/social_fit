package events

import (
	"log"
	"strings"
	"time"

	"github.com/freitasmatheusrn/social-fit/pkg/rest"
)

func (e Event) ValidateFields() *rest.ApiErr {
	var causes []rest.Causes
	if len(e.Title) < 10 {
		causes = append(causes, rest.Causes{Field: "title", Message: "Título precisa ter mais de 10 caracteres"})
	}
	if !e.Status.IsValid() {
		causes = append(causes, rest.Causes{Field: "status", Message: "Status inválido"})
	}
	if strings.TrimSpace(e.City) == "" {
		causes = append(causes, rest.Causes{Field: "city", Message: "Nome da cidade não pode ficar vazio"})
	}
	if strings.TrimSpace(e.Street) == "" {
		causes = append(causes, rest.Causes{Field: "street", Message: "Nome da rua não pode ficar vazio"})
	}
	if strings.TrimSpace(e.Neighbourhood) == "" {
		causes = append(causes, rest.Causes{Field: "neighbourhood", Message: "Nome do bairro não pode ficar vazio"})
	}
	if e.StartDate.After(e.EndDate) {
		causes = append(causes, rest.Causes{Field: "start_date", Message: "Data de início não pode ser após data do fim"})
	}
	if e.EndDate.Before(time.Now()) {

		log.Printf("hora de entrada: %v, hora de saída: %v", e.StartDate, e.EndDate)
		causes = append(causes, rest.Causes{Field: "end_date", Message: "Data de fim não pode ser no passado"})
	}
	if e.MaxCapacity < 0 {
		causes = append(causes, rest.Causes{Field: "max_capacity", Message: "Não pode ser menor que 0"})
	}
	if causes != nil {
		return rest.NewBadRequestValidationError("Campo(s) inválidos", causes)
	}
	return nil
}

func (s EventStatus) IsValid() bool {
	switch s {
	case EventStatusDraft, EventStatusPublished, EventStatusCancelled, EventStatusFinished:
		return true
	}
	return false
}
