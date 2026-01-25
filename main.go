package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pos-api/internal/http/handler"
	"pos-api/internal/repository_memory"
	"pos-api/internal/service"
)

func main() {
	// Product
	productRepo := repository_memory.NewProductRepo()
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)
	http.HandleFunc("GET /api/products", productHandler.GetProducts)
	http.HandleFunc("GET /api/products/", productHandler.GetProductByID)
	http.HandleFunc("POST /api/products", productHandler.CreateProduct)
	http.HandleFunc("PUT /api/products/", productHandler.UpdateProduct)
	http.HandleFunc("DELETE /api/products/", productHandler.DeleteProduct)

	// Category
	categoryRepo := repository_memory.NewCategoryRepo()
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	http.HandleFunc("GET /api/categories", categoryHandler.GetCategories)
	http.HandleFunc("GET /api/categories/", categoryHandler.GetCategoryByID)
	http.HandleFunc("POST /api/categories", categoryHandler.CreateCategory)
	http.HandleFunc("PUT /api/categories/", categoryHandler.UpdateCategory)
	http.HandleFunc("DELETE /api/categories/", categoryHandler.DeleteCategory)

	// Docs
	docsHandler := handler.NewDocsHandler("openapi.json")
	http.HandleFunc("GET /openapi.json", docsHandler.ServeSpec)
	http.HandleFunc("GET /docs", docsHandler.ServeDocs)
	http.HandleFunc("GET /docs/", docsHandler.RedirectDocs)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Starting server on :8081")

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
