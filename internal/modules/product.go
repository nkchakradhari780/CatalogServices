package modules


type Product struct {
    ID         int   `json:"id,omitempty" `        
    Name       string   `json:"name" validate:"required"`                
    Price      int      `json:"price" validate:"required"`               
    Stock      int      `json:"stock" validate:"required"`              
    CategoryID string   `json:"category_id" validate:"required"`  
    Quantity    int     `json:"quantity" validate:"required"`       
    Brand      string   `json:"brand" validate:"required"`               
    Images     []string `json:"images,omitempty"`    
}
