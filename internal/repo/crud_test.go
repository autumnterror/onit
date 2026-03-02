package repo

import (
	"github.com/autumnterror/onit/internal/domain"
	"github.com/autumnterror/onit/internal/infra/psql"
	"github.com/autumnterror/utils_go/pkg/utils/uid"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestGood(t *testing.T) {
	db, err := psql.NewDB("postgres://postgres:postgrespw@localhost:5432/db?sslmode=disable")
	assert.NoError(t, err)

	assert.NoError(t, db.AutoMigrate(&domain.Product{}))

	id := uid.New()

	repo := NewProductRepository(db)
	product := &domain.Product{
		ID:          id,
		Title:       "Keyboard",
		Description: "Mechanical keyboard",
		Price:       12000,
		Photo:       "photo.jpg",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	assert.NoError(t, repo.Create(product))

	all, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(all))
	log.Println(all)

	newProduct := &domain.Product{
		ID:          id,
		Title:       "New",
		Description: "New",
		Price:       12321312,
		Photo:       "New",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	assert.NoError(t, repo.Update(newProduct))

	one, err := repo.GetByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, one)
	assert.Equal(t, newProduct, one)

	log.Println(one)

	assert.NoError(t, repo.Delete(id))

	all, err = repo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(all))
	log.Println(all)
}
