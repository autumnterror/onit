package service

import (
	"errors"
	"testing"
	"time"

	"github.com/autumnterror/onit/internal/domain"
	"github.com/autumnterror/onit/internal/repo/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupTest(t *testing.T) (*gomock.Controller, *mocks.MockProductRepository, *Service) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockProductRepository(ctrl)
	svc := NewService(mockRepo)
	return ctrl, mockRepo, svc
}

func TestService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		inputProduct   *domain.Product
		setupMocks     func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product)
		expectedError  error
		validateResult func(t *testing.T, product *domain.Product)
	}{
		{
			name: "Success - Valid product with all fields",
			inputProduct: &domain.Product{
				ID:          uuid.New().String(),
				Title:       "Ноутбук Apple MacBook Pro",
				Description: "Мощный ноутбук для разработки",
				Price:       150000,
				Photo:       "https://example.com/photo.jpg",
				CreatedAt:   time.Now().Unix(),
				UpdatedAt:   time.Now().Unix(),
			},
			setupMocks: func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product) {
				mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(p *domain.Product) error {
					assert.Equal(t, expectedProduct.ID, p.ID)
					assert.Equal(t, expectedProduct.Title, p.Title)
					assert.Equal(t, expectedProduct.Price, p.Price)
					return nil
				}).Times(1)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, product *domain.Product) {
				assert.NotEmpty(t, product.ID)
				assert.NotEmpty(t, product.Title)
				assert.NotEmpty(t, product.Description)
				assert.NotZero(t, product.Price)
			},
		},
		{
			name: "Success - Auto-fill empty title with default",
			inputProduct: &domain.Product{
				ID:          uuid.New().String(),
				Title:       "",
				Description: "Описание товара",
				Price:       50000,
				Photo:       "https://example.com/photo.jpg",
			},
			setupMocks: func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product) {
				mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(p *domain.Product) error {
					assert.Equal(t, "Без названия", p.Title)
					assert.Equal(t, expectedProduct.Description, p.Description)
					return nil
				}).Times(1)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, product *domain.Product) {
				assert.Equal(t, "Без названия", product.Title)
			},
		},
		{
			name: "Success - Auto-fill empty description",
			inputProduct: &domain.Product{
				ID:          uuid.New().String(),
				Title:       "Смартфон Samsung",
				Description: "",
				Price:       70000,
				Photo:       "https://example.com/phone.jpg",
			},
			setupMocks: func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product) {
				mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(p *domain.Product) error {
					assert.Equal(t, "Без описания", p.Description)
					return nil
				}).Times(1)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, product *domain.Product) {
				assert.Equal(t, "Без описания", product.Description)
			},
		},
		{
			name: "Success - Auto-fill empty photo with default URL",
			inputProduct: &domain.Product{
				ID:          uuid.New().String(),
				Title:       "Наушники Sony",
				Description: "Отличный звук",
				Price:       15000,
				Photo:       "",
			},
			setupMocks: func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product) {
				mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(p *domain.Product) error {
					assert.Contains(t, p.Photo, "encrypted-tbn0.gstatic.com")
					return nil
				}).Times(1)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, product *domain.Product) {
				assert.NotEmpty(t, product.Photo)
			},
		},
		{
			name: "Success - Auto-fill zero price with default 1,000,000",
			inputProduct: &domain.Product{
				ID:          uuid.New().String(),
				Title:       "Элитный товар",
				Description: "Дорогой товар",
				Price:       0,
				Photo:       "https://example.com/elite.jpg",
			},
			setupMocks: func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product) {
				mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(p *domain.Product) error {
					assert.Equal(t, 1000000, p.Price)
					return nil
				}).Times(1)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, product *domain.Product) {
				assert.Equal(t, 1000000, product.Price)
			},
		},
		{
			name: "Success - Auto-fill timestamps if zero",
			inputProduct: &domain.Product{
				ID:          uuid.New().String(),
				Title:       "Товар с автодатой",
				Description: "Тест",
				Price:       1000,
				Photo:       "photo.jpg",
				CreatedAt:   0,
				UpdatedAt:   0,
			},
			setupMocks: func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product) {
				mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(p *domain.Product) error {
					assert.NotZero(t, p.CreatedAt)
					assert.NotZero(t, p.UpdatedAt)
					return nil
				}).Times(1)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, product *domain.Product) {
				assert.NotZero(t, product.CreatedAt)
				assert.NotZero(t, product.UpdatedAt)
			},
		},
		{
			name: "Error - Invalid ID (not UUID)",
			inputProduct: &domain.Product{
				ID:          "invalid-id-format",
				Title:       "Ноутбук",
				Description: "Описание",
				Price:       50000,
			},
			setupMocks: func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product) {
				mockRepo.EXPECT().Create(gomock.Any()).Times(0)
			},
			expectedError: errors.New("неправильный id"),
			validateResult: func(t *testing.T, product *domain.Product) {
			},
		},
		{
			name: "Error - Empty ID",
			inputProduct: &domain.Product{
				ID:          "",
				Title:       "Товар без ID",
				Description: "Описание",
				Price:       30000,
			},
			setupMocks: func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product) {
				mockRepo.EXPECT().Create(gomock.Any()).Times(0)
			},
			expectedError: errors.New("неправильный id"),
			validateResult: func(t *testing.T, product *domain.Product) {
			},
		},
		{
			name: "Error - Repository returns error",
			inputProduct: &domain.Product{
				ID:          uuid.New().String(),
				Title:       "Тестовый товар",
				Description: "Описание",
				Price:       100000,
				Photo:       "photo.jpg",
			},
			setupMocks: func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product) {
				mockRepo.EXPECT().Create(gomock.Any()).Return(errors.New("database connection failed")).Times(1)
			},
			expectedError: ErrServer,
			validateResult: func(t *testing.T, product *domain.Product) {
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl, mockRepo, svc := setupTest(t)
			defer ctrl.Finish()

			expectedProduct := *tt.inputProduct

			tt.setupMocks(mockRepo, &expectedProduct)

			err := svc.Create(tt.inputProduct)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
				if tt.validateResult != nil {
					tt.validateResult(t, tt.inputProduct)
				}
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		inputProduct  *domain.Product
		setupMocks    func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product)
		expectedError error
	}{
		{
			name: "Success - Update valid product",
			inputProduct: &domain.Product{
				ID:          uuid.New().String(),
				Title:       "Обновленный ноутбук",
				Description: "Новое описание",
				Price:       120000,
				Photo:       "new-photo.jpg",
			},
			setupMocks: func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product) {
				mockRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(p *domain.Product) error {
					assert.Equal(t, expectedProduct.ID, p.ID)
					assert.Equal(t, "Обновленный ноутбук", p.Title)
					return nil
				}).Times(1)
			},
			expectedError: nil,
		},
		{
			name: "Success - Update with auto-fill empty fields",
			inputProduct: &domain.Product{
				ID:          uuid.New().String(),
				Title:       "",
				Description: "",
				Price:       0,
				Photo:       "",
			},
			setupMocks: func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product) {
				mockRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(p *domain.Product) error {
					assert.Equal(t, "Без названия", p.Title)
					assert.Equal(t, "Без описания", p.Description)
					assert.Equal(t, 1000000, p.Price)
					assert.Contains(t, p.Photo, "encrypted-tbn0.gstatic.com")
					assert.NotZero(t, p.UpdatedAt)
					return nil
				}).Times(1)
			},
			expectedError: nil,
		},
		{
			name: "Error - Invalid ID format",
			inputProduct: &domain.Product{
				ID:          "not-a-uuid",
				Title:       "Товар",
				Description: "Описание",
				Price:       50000,
			},
			setupMocks: func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product) {
				mockRepo.EXPECT().Update(gomock.Any()).Times(0)
			},
			expectedError: errors.New("неправильный id"),
		},
		{
			name: "Error - Repository returns error",
			inputProduct: &domain.Product{
				ID:          uuid.New().String(),
				Title:       "Товар для обновления",
				Description: "Описание",
				Price:       75000,
			},
			setupMocks: func(mockRepo *mocks.MockProductRepository, expectedProduct *domain.Product) {
				mockRepo.EXPECT().Update(gomock.Any()).Return(errors.New("update failed")).Times(1)
			},
			expectedError: ErrServer,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl, mockRepo, svc := setupTest(t)
			defer ctrl.Finish()

			expectedProduct := *tt.inputProduct
			tt.setupMocks(mockRepo, &expectedProduct)

			err := svc.Update(tt.inputProduct)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	validUUID := uuid.New().String()

	tests := []struct {
		name          string
		id            string
		setupMocks    func(mockRepo *mocks.MockProductRepository)
		expectedError error
	}{
		{
			name: "Success - Delete valid ID",
			id:   validUUID,
			setupMocks: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().Delete(validUUID).Return(nil).Times(1)
			},
			expectedError: nil,
		},
		{
			name: "Error - Invalid ID format",
			id:   "invalid-id",
			setupMocks: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().Delete(gomock.Any()).Times(0)
			},
			expectedError: errors.New("неправильный id"),
		},
		{
			name: "Error - Empty ID",
			id:   "",
			setupMocks: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().Delete(gomock.Any()).Times(0)
			},
			expectedError: errors.New("неправильный id"),
		},
		{
			name: "Error - Repository returns error",
			id:   validUUID,
			setupMocks: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().Delete(validUUID).Return(errors.New("delete failed")).Times(1)
			},
			expectedError: ErrServer,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl, mockRepo, svc := setupTest(t)
			defer ctrl.Finish()

			tt.setupMocks(mockRepo)

			err := svc.Delete(tt.id)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetById(t *testing.T) {
	t.Parallel()

	validUUID := uuid.New().String()
	expectedProduct := &domain.Product{
		ID:          validUUID,
		Title:       "Найденный товар",
		Description: "Описание товара",
		Price:       99900,
		Photo:       "photo-url.jpg",
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	tests := []struct {
		name          string
		id            string
		setupMocks    func(mockRepo *mocks.MockProductRepository)
		expectedError error
		expectedRes   *domain.Product
	}{
		{
			name: "Success - Get by valid ID",
			id:   validUUID,
			setupMocks: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().GetByID(validUUID).Return(expectedProduct, nil).Times(1)
			},
			expectedRes:   expectedProduct,
			expectedError: nil,
		},
		{
			name: "Error - Invalid ID format",
			id:   "not-uuid",
			setupMocks: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().GetByID(gomock.Any()).Times(0)
			},
			expectedRes:   nil,
			expectedError: errors.New("неправильный id"),
		},
		{
			name: "Error - Empty ID",
			id:   "",
			setupMocks: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().GetByID(gomock.Any()).Times(0)
			},
			expectedRes:   nil,
			expectedError: errors.New("неправильный id"),
		},
		{
			name: "Error - Product not found in repository",
			id:   uuid.New().String(),
			setupMocks: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().GetByID(gomock.Any()).Return(nil, errors.New("record not found")).Times(1)
			},
			expectedRes:   nil,
			expectedError: ErrServer,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl, mockRepo, svc := setupTest(t)
			defer ctrl.Finish()

			tt.setupMocks(mockRepo)

			result, err := svc.GetById(tt.id)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.expectedRes.ID, result.ID)
				assert.Equal(t, tt.expectedRes.Title, result.Title)
				assert.Equal(t, tt.expectedRes.Price, result.Price)
			}
		})
	}
}

func TestService_GetAll(t *testing.T) {
	t.Parallel()

	expectedProducts := []domain.Product{
		{
			ID:          uuid.New().String(),
			Title:       "Товар 1",
			Description: "Описание 1",
			Price:       1000,
			Photo:       "photo1.jpg",
		},
		{
			ID:          uuid.New().String(),
			Title:       "Товар 2",
			Description: "Описание 2",
			Price:       2000,
			Photo:       "photo2.jpg",
		},
	}

	tests := []struct {
		name          string
		setupMocks    func(mockRepo *mocks.MockProductRepository)
		expectedError error
		expectedCount int
	}{
		{
			name: "Success - Get all products",
			setupMocks: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().GetAll().Return(expectedProducts, nil).Times(1)
			},
			expectedError: nil,
			expectedCount: 2,
		},
		{
			name: "Success - Empty list",
			setupMocks: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().GetAll().Return([]domain.Product{}, nil).Times(1)
			},
			expectedError: nil,
			expectedCount: 0,
		},
		{
			name: "Error - Repository returns error",
			setupMocks: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().GetAll().Return(nil, errors.New("database error")).Times(1)
			},
			expectedError: ErrServer,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl, mockRepo, svc := setupTest(t)
			defer ctrl.Finish()

			tt.setupMocks(mockRepo)

			results, err := svc.GetAll()

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, results)
			} else {
				require.NoError(t, err)
				assert.Len(t, results, tt.expectedCount)
			}
		})
	}
}

func TestService_Validation_DoesNotOverrideExistingValues(t *testing.T) {
	t.Parallel()

	ctrl, mockRepo, svc := setupTest(t)
	defer ctrl.Finish()

	existingTimestamp := time.Now().Unix()
	product := &domain.Product{
		ID:          uuid.New().String(),
		Title:       "Существующее название",
		Description: "Существующее описание",
		Price:       75000,
		Photo:       "existing-photo.jpg",
		CreatedAt:   existingTimestamp,
		UpdatedAt:   existingTimestamp,
	}

	mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(p *domain.Product) error {
		assert.Equal(t, "Существующее название", p.Title)
		assert.Equal(t, "Существующее описание", p.Description)
		assert.Equal(t, 75000, p.Price)
		assert.Equal(t, "existing-photo.jpg", p.Photo)
		assert.Equal(t, existingTimestamp, p.CreatedAt)
		assert.Equal(t, existingTimestamp, p.UpdatedAt)
		return nil
	}).Times(1)

	err := svc.Create(product)
	assert.NoError(t, err)
}
