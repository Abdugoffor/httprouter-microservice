package controllers

import (
	"encoding/json"
	"net/http"
	"shop_service/request"
	"shop_service/services"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ProductController struct {
	service services.ProductService
}

func NewProductController(s services.ProductService) *ProductController {
	return &ProductController{service: s}
}

// CREATE
func (c *ProductController) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var req request.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err := c.service.Create(req.Name, req.Price)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("created"))
}

// GET ALL
func (c *ProductController) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	data, _ := c.service.GetAll()
	JSON(w, http.StatusOK, data)
	// json.NewEncoder(w).Encode(data)
}

// GET BY ID
func (c *ProductController) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, _ := strconv.Atoi(ps.ByName("id"))

	data, err := c.service.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}

	JSON(w, http.StatusOK, data)
	// json.NewEncoder(w).Encode(data)
}

// UPDATE
func (c *ProductController) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, _ := strconv.Atoi(ps.ByName("id"))

	var req request.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err := c.service.Update(id, req.Name, req.Price)
	if err != nil {
		http.Error(w, "update failed", 400)
		return
	}

	w.Write([]byte("updated"))
}

// DELETE
func (c *ProductController) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id, _ := strconv.Atoi(ps.ByName("id"))

	err := c.service.Delete(id)
	if err != nil {
		http.Error(w, "delete failed", 400)
		return
	}

	w.Write([]byte("deleted"))
}
