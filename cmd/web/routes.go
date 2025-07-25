package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)

	mux.Get("/", app.Home)
	mux.Get("/virtual-terminal", app.VirtualTerminal)
	mux.Post("/virtual-terminal-payment-succeeded", app.VirtualTerminalPaymentSucceeded)
	mux.Get("/virtual-terminal-receipt", app.VirtualTerminalReceipt)

	mux.Post("/payment-succeeded", app.PaymentSucceeded)
	mux.Get("/receipt", app.Receipt)
	mux.Get("/widget/{id}", app.ChargeOnce)

	mux.Get("/plans/bronze", app.BronzePlan)
	mux.Get("/receipt/bronze", app.BronzePlanReceipt)

	// auth routes
	mux.Get("/login", app.LoginPage)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
