package controller

import (
	"errors"
	
	"gorm.io/gorm"
	status "google.golang.org/grpc/status"
	codes "google.golang.org/grpc/codes"
	"github.com/google/uuid"
	
	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	log "github.com/Minerva-System/minerva-go/pkg/log"
	repository "github.com/Minerva-System/minerva-go/internal/svc/products/repository"
	model "github.com/Minerva-System/minerva-go/internal/model"
)

const PAGESIZE = 100

func ProductsIndex(db *gorm.DB, companyId string, page int64) (*rpc.ProductList, error) {
	if page < 0 {
		log.Error("Invalid page index: %d", page)
		return nil, status.Error(codes.OutOfRange, "Invalid product list page")
	}

	parsedCompanyId, err := uuid.Parse(companyId)
	if err != nil {
		log.Error("Company UUID is invalid: \"%s\"", companyId)
		return nil, status.Error(codes.InvalidArgument, "Index parameter has an invalid company UUID")
	}

	list, err := repository.ListProducts(db, parsedCompanyId, int(page) * PAGESIZE, PAGESIZE)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	return model.Product{}.ListToMessage(list), nil
}

func GetProduct(db *gorm.DB, companyId string, id string) (*rpc.Product, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		log.Error("UUID is invalid: \"%s\"", id)
		return nil, status.Error(codes.InvalidArgument, "Index parameter has an invalid UUID")
	}

	parsedCompanyId, err := uuid.Parse(companyId)
	if err != nil {
		log.Error("Company UUID is invalid: \"%s\"", companyId)
		return nil, status.Error(codes.InvalidArgument, "Index parameter has an invalid company UUID")
	}
	
	prd, err := repository.GetProduct(db, parsedCompanyId, parsedId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("Product %s not found", id)
			return nil, status.Error(codes.NotFound, "Product not found")
		}
		
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	msg := prd.ToMessage()
	return &msg, nil
}

func CreateProduct(db *gorm.DB, data *rpc.Product) (*rpc.Product, error) {
	if data.Id != nil {
		log.Error("New product must not have a predefined id")
		return nil, status.Error(codes.InvalidArgument, "New product must not have a predefined id")
	}

	m, err := model.Product{}.FromMessage(data)
	if err != nil {
		log.Error("Error mapping product data: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Error mapping product data: %v", err)
	}

	log.Debug("Creating product...")
	created, err := repository.CreateProduct(db, m)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	msg := created.ToMessage()
	return &msg, nil
}

func DeleteProduct(db *gorm.DB, companyId string, id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		log.Error("UUID is invalid: \"%s\"", id)
		return status.Error(codes.InvalidArgument, "Index parameter has an invalid UUID")
	}

	parsedCompanyId, err := uuid.Parse(companyId)
	if err != nil {
		log.Error("Company UUID is invalid: \"%s\"", companyId)
		return status.Error(codes.InvalidArgument, "Index parameter has an invalid company UUID")
	}

	exists, err := repository.ExistsProduct(db, parsedCompanyId, parsedId)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	if !exists {
		log.Error("Product %s not found", id)
		return status.Error(codes.NotFound, "Product not found")
	}
	
	err = repository.DeleteProduct(db, parsedCompanyId, parsedId)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}
	return nil
}

func UpdateProduct(db *gorm.DB, data *rpc.Product) (*rpc.Product, error) {
	if data.Id == nil {
		log.Error("Product id is missing")
		return nil, status.Error(codes.InvalidArgument, "Product id is missing")
	}

	if len(data.Unit) > 2 {
		log.Error("Unit has more than 2 characters")
		return nil, status.Error(codes.InvalidArgument, "Unit has more than 2 characters")
	}

	d, err := model.Product{}.FromMessage(data)
	if err != nil {
		log.Error("Error mapping product data: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Error mapping product data: %v", err)
	}

	result, err := repository.UpdateProduct(db, d)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	msg := result.ToMessage()
	return &msg, nil
}

