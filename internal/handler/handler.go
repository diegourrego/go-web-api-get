package handler

//func (h *MyHandler) CreateProduct() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var requestBody RequestBodyProduct
//		// Toca convertir el Json a estructura
//		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
//			code := http.StatusBadRequest
//			body := &ResponseBodyProduct{
//				Message: "Bad Request. check your body data",
//				Data:    nil,
//				Error:   true,
//			}
//			w.Header().Set("Content-Type", "application/json")
//			w.WriteHeader(code)
//			json.NewEncoder(w).Encode(body)
//		}
//
//		// Serializar
//		p := Product{
//			ID:          len(h.data.Products) + 1,
//			Name:        requestBody.Name,
//			Quantity:    requestBody.Quantity,
//			CodeValue:   requestBody.CodeValue,
//			IsPublished: requestBody.IsPublished,
//			Expiration:  requestBody.Expiration,
//		}
//
//		h.data.Products = append(h.data.Products, p)
//
//		// Le respondemos al cliente
//		code := http.StatusCreated
//		body := ResponseBodyProduct{
//			Message: "Product created successfully",
//			Data: &struct {
//				ID          int     `json:"id"`
//				Name        string  `json:"name"`
//				Quantity    int     `json:"quantity"`
//				CodeValue   string  `json:"code_value"`
//				IsPublished bool    `json:"is_published"`
//				Expiration  string  `json:"expiration"`
//				Price       float64 `json:"price"`
//			}{ID: p.ID, Name: p.Name, Quantity: p.Quantity, CodeValue: p.CodeValue,
//				IsPublished: p.IsPublished, Expiration: p.Expiration, Price: p.Price},
//			Error: false,
//		}
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(code)
//		if err := json.NewEncoder(w).Encode(body); err != nil {
//			fmt.Println(err)
//			return
//		}
//
//	}
//}
