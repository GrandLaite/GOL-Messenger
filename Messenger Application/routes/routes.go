package routes

import (
	"gol_messenger/auth"
	"gol_messenger/messages"
	"gol_messenger/users"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(
	userHandler *users.UserHandler,
	messageHandler *messages.MessageHandler,
	authMiddleware *auth.AuthMiddleware,
) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/login", userHandler.LoginUserHandler).Methods("POST")
	router.HandleFunc("/register", userHandler.RegisterUserHandler).Methods("POST")

	router.Handle("/user", authMiddleware.Middleware(http.HandlerFunc(userHandler.GetUserHandler))).Methods("GET")
	router.Handle("/user/update", authMiddleware.Middleware(http.HandlerFunc(userHandler.UpdateUserHandler))).Methods("PUT")
	router.Handle("/user/delete", authMiddleware.Middleware(http.HandlerFunc(userHandler.DeleteUserHandler))).Methods("DELETE")

	
	router.Handle("/messages", authMiddleware.Middleware(http.HandlerFunc(messageHandler.ListMessagesHandler))).Methods("GET")
	router.Handle("/message", authMiddleware.Middleware(http.HandlerFunc(messageHandler.GetMessageHandler))).Methods("GET")
	router.Handle(
		"/message/create",
		authMiddleware.RoleMiddleware(auth.RoleUser)(
			authMiddleware.Middleware(http.HandlerFunc(messageHandler.CreateMessageHandler)),
		),
	).Methods("POST")
	router.Handle(
		"/message/update",
		authMiddleware.RoleMiddleware(auth.RolePremium)(
			authMiddleware.Middleware(http.HandlerFunc(messageHandler.UpdateMessageHandler)),
		),
	).Methods("PUT")
	router.Handle(
		"/message/delete",
		authMiddleware.RoleMiddleware(auth.RolePremium)(
			authMiddleware.Middleware(http.HandlerFunc(messageHandler.DeleteMessageHandler)),
		),
	).Methods("DELETE")
	router.Handle(
		"/message/like",
		authMiddleware.RoleMiddleware(auth.RoleUser)(
			authMiddleware.Middleware(http.HandlerFunc(messageHandler.LikeMessageHandler)),
		),
	).Methods("POST")
	router.Handle(
		"/message/superlike",
		authMiddleware.RoleMiddleware(auth.RolePremium)(
			authMiddleware.Middleware(http.HandlerFunc(messageHandler.SuperlikeMessageHandler)),
		),
	).Methods("POST")

	return router
}
