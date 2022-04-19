package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	renderChi "github.com/go-chi/render"
	"github.com/joleques/go-sqs-poc/src/sqs"
	renderPkg "github.com/unrolled/render"
	"net/http"
	"time"
)

var render *renderPkg.Render

type Message struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func (a *Message) Bind(r *http.Request) error {
	return nil
}

type Result struct {
	StatusCod int    `json:"statusCod"`
	Message   string `json:"message"`
}

func Start() {
	render = renderPkg.New()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(renderChi.SetContentType(renderChi.ContentTypeJSON))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Status Ok"))
	})

	r.Route("/producer", func(r chi.Router) {
		r.Post("/", SaveProducer)
	})

	http.ListenAndServe(":3000", r)
}

func SaveProducer(writer http.ResponseWriter, request *http.Request) {
	data := &Message{}
	if err := renderChi.Bind(request, data); err != nil {
		render.JSON(writer, 400, Result{StatusCod: 400, Message: "Invalid request"})
		return
	}
	messageSqs := sqs.MessageSQS{Id: data.Id, Message: data.Message}
	resultProducer, err := sqs.Send(messageSqs)
	if err != nil {
		render.JSON(writer, 400, Result{StatusCod: 400, Message: err.Error()})
		return
	}
	result := Result{StatusCod: 201, Message: resultProducer}
	render.JSON(writer, result.StatusCod, result)
}
