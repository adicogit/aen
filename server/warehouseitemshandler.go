package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"aen.it/poolmanager/log"
	"aen.it/poolmanager/warehouse"
	"github.com/gorilla/mux"
)

const DefaultPageSize = 20

type paginatedItemIDList struct {
	ID         []string `json:"id"`
	Page       int      `json:"page"`
	PageSize   int      `json:"pageSize"`
	TotalItems int      `json:"totalItems"`
	TotalPages int      `json:"totalPages"`
	HasItems   bool     `json:"hasItems"`
}

type createWarehouseItemRequest struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	PublicPrice   int    `json:"publicPrice"`
	IncomingPrice int    `json:"incomingPrice"`
	Quantity      int    `json:"quantity"`
}

type modifyWarehouseItemRequest struct {
	Name          string `json:"name"`
	PublicPrice   int    `json:"publicPrice"`
	IncomingPrice int    `json:"incomingPrice"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Message string `json:"message"`
}

func (server *Server) getWarehouseItems(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering getWarehouseItems")
	w.Header().Set("Content-Type", "application/json")

	// Parse pagination parameters from query string
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// Default values
	page := 1
	pageSize := DefaultPageSize

	// Parse page parameter
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else {
			log.Log.Warn("Invalid page parameter, using default", "page", pageStr)
		}
	}

	// Parse pageSize parameter
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		} else {
			log.Log.Warn("Invalid pageSize parameter, using default", "pageSize", pageSizeStr)
		}
	}

	// Get all items
	allItems := server.billiardManager.GetItemIDs()
	totalItems := len(allItems)

	// Calculate pagination
	totalPages := (totalItems + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	// Ensure page is within valid range
	if page > totalPages {
		page = totalPages
	}

	// Calculate start and end indices
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	// Ensure indices are within bounds
	if startIndex < 0 {
		startIndex = 0
	}
	if endIndex > totalItems {
		endIndex = totalItems
	}

	// Extract paginated items
	var paginatedItems []string
	if startIndex < totalItems {
		paginatedItems = allItems[startIndex:endIndex]
	} else {
		paginatedItems = []string{}
	}

	result := paginatedItemIDList{
		ID:         paginatedItems,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasItems:   totalItems > 0,
	}

	log.Log.Debug("Exiting getWarehouseItems", "page", page, "pageSize", pageSize, "totalItems", totalItems, "totalPages", totalPages)
	json.NewEncoder(w).Encode(result)
}

func (server *Server) getWarehouseItem(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering getWarehouseItem")
	w.Header().Set("Content-Type", "application/json")

	// Get item ID from URL path
	vars := mux.Vars(r)
	itemID := vars["itemID"]

	if itemID == "" {
		log.Log.Error("Item ID is required")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Item ID is required"})
		return
	}

	// Get the item from warehouse
	item, err := server.billiardManager.GetItem(itemID)
	if err != nil {
		log.Log.Error("Failed to get item", "error", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	// Get the quantity
	quantity := server.billiardManager.GetItemQuantity(itemID)

	// Create response with item details and quantity
	response := struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		PublicPrice   int    `json:"publicPrice"`
		IncomingPrice int    `json:"incomingPrice"`
		Quantity      int    `json:"quantity"`
	}{
		ID:            item.ID,
		Name:          item.Name,
		PublicPrice:   item.PublicPrice,
		IncomingPrice: item.IncomingPrice,
		Quantity:      quantity,
	}

	log.Log.Debug("Exiting getWarehouseItem", "itemID", itemID)
	json.NewEncoder(w).Encode(response)
}

func (server *Server) createWarehouseItems(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering createWarehouseItems")
	w.Header().Set("Content-Type", "application/json")

	var req createWarehouseItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Log.Error("Failed to decode request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	// Validate required fields
	if req.ID == "" {
		log.Log.Error("Item ID is required")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Item ID is required"})
		return
	}

	if req.Name == "" {
		log.Log.Error("Item name is required")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Item name is required"})
		return
	}

	if req.Quantity <= 0 {
		log.Log.Error("Quantity must be greater than 0")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Quantity must be greater than 0"})
		return
	}

	if req.PublicPrice < 0 {
		log.Log.Error("PublicPrice must be non-negative")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "PublicPrice must be non-negative"})
		return
	}

	if req.IncomingPrice < 0 {
		log.Log.Error("IncomingPrice must be non-negative")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "IncomingPrice must be non-negative"})
		return
	}

	// Create the item
	item := warehouse.Item{
		ID:            req.ID,
		Name:          req.Name,
		PublicPrice:   req.PublicPrice,
		IncomingPrice: req.IncomingPrice,
	}

	// Add items to warehouse
	server.billiardManager.AddItems(item, req.Quantity)

	log.Log.Debug("Exiting createWarehouseItems", "itemID", req.ID, "quantity", req.Quantity)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(successResponse{Message: "Item created successfully"})
}

func (server *Server) modifyWarehouseItems(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering modifyWarehouseItems")
	w.Header().Set("Content-Type", "application/json")

	// Get item ID from URL path
	vars := mux.Vars(r)
	itemID := vars["itemID"]

	if itemID == "" {
		log.Log.Error("Item ID is required")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Item ID is required"})
		return
	}

	var req modifyWarehouseItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Log.Error("Failed to decode request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	// Validate required fields
	if req.Name == "" {
		log.Log.Error("Item name is required")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Item name is required"})
		return
	}

	if req.PublicPrice < 0 {
		log.Log.Error("PublicPrice must be non-negative")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "PublicPrice must be non-negative"})
		return
	}

	if req.IncomingPrice < 0 {
		log.Log.Error("IncomingPrice must be non-negative")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "IncomingPrice must be non-negative"})
		return
	}

	// Update the item
	err := server.billiardManager.UpdateItem(itemID, req.Name, req.PublicPrice, req.IncomingPrice)
	if err != nil {
		log.Log.Error("Failed to update item", "error", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	log.Log.Debug("Exiting modifyWarehouseItems", "itemID", itemID)
	json.NewEncoder(w).Encode(successResponse{Message: "Item updated successfully"})
}

func (server *Server) deleteWarehouseItems(w http.ResponseWriter, r *http.Request) {
	log.Log.Debug("Entering deleteWarehouseItems")
	w.Header().Set("Content-Type", "application/json")

	// Get item ID from URL path
	vars := mux.Vars(r)
	itemID := vars["itemID"]

	if itemID == "" {
		log.Log.Error("Item ID is required")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Item ID is required"})
		return
	}

	// Delete the item
	err := server.billiardManager.DeleteItem(itemID)
	if err != nil {
		log.Log.Error("Failed to delete item", "error", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	log.Log.Debug("Exiting deleteWarehouseItems", "itemID", itemID)
	json.NewEncoder(w).Encode(successResponse{Message: "Item deleted successfully"})
}
