package domain

import "gorm.io/gorm"

type Product struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	Price       int    `gorm:"not null"`
	Photo       string `gorm:"size:500"`
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type ProductHttp struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Photo       string `json:"photo"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

func (p *Product) ToHttp() *ProductHttp {
	return &ProductHttp{
		Id:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Price:       p.Price,
		Photo:       p.Photo,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
func (p *ProductHttp) ToDomain() *Product {
	return &Product{
		ID:          p.Id,
		Title:       p.Title,
		Description: p.Description,
		Price:       p.Price,
		Photo:       p.Photo,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ToHttpMany(prs []Product) []ProductHttp {
	var prsHTTP []ProductHttp
	for _, p := range prs {
		prsHTTP = append(prsHTTP, *p.ToHttp())
	}
	return prsHTTP
}
