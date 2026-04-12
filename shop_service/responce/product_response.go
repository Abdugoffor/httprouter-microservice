package response

import "shop_service/models"

type ProductResponse struct {
	Data []models.Product `json:"data"`
}

type SingleProductResponse struct {
	Data models.Product `json:"data"`
}
