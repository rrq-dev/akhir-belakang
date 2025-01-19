package routes

import (
	"akhir-belakang/config"
	"akhir-belakang/controller"
	"akhir-belakang/helper"
	"net/http"
)

func URL(w http.ResponseWriter, r *http.Request) {
	// set access control header
	if config.SetAccessControlHeaders(w, r){
		return
	}

	//ambil metode dan path dari request
	method := r.Method
	path := r.URL.Path

	switch {
	case method == "POST" && path == "/register":
		controller.Register(w, r)
	case method == "POST" && path == "/login":
		controller.Login(w, r)
	
	// crud untuk feedback
	case method == "GET" && path == "/datalokasi":
		helper.ValidateTokenMiddleware(http.HandlerFunc(controller.GetDataLocation)).ServeHTTP(w, r)
	
	case method == "POST" && path == "/createlokasi":
		helper.ValidateTokenMiddleware(helper.RoleMiddleware("admin")(http.HandlerFunc(controller.CreateDataLocations))).ServeHTTP(w, r)	
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}


}