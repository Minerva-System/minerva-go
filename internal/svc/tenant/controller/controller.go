package controller

import (
	"errors"
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
	status "google.golang.org/grpc/status"
	codes "google.golang.org/grpc/codes"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"

	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	repository "github.com/Minerva-System/minerva-go/internal/svc/tenant/repository"
	log "github.com/Minerva-System/minerva-go/pkg/log"
	model "github.com/Minerva-System/minerva-go/internal/model"
	util "github.com/Minerva-System/minerva-go/pkg/util"
)

const PAGESIZE = 100

func CompanyIndex(db *gorm.DB, page int64) (*rpc.CompanyList, error) {
	if page < 0 {
		log.Error("Invalid page index: %d", page)
		return nil, status.Error(codes.OutOfRange, "Invalid company list page")
	}

	list, err := repository.ListCompanies(db, int(page) * PAGESIZE, PAGESIZE)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	return model.Company{}.ListToMessage(list), nil
}

func GetCompany(db *gorm.DB, id string) (*rpc.Company, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		log.Error("\"%s\" is an invalid UUID", id)
		return nil, status.Error(codes.InvalidArgument, "Index parameter is an invalid UUID")
	}

	company, err := repository.GetCompany(db, parsedId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("Company %s not found", id)
			return nil, status.Error(codes.NotFound, "Company not found")
		}

		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	msg := company.ToMessage()
	return &msg, nil
}

func GetCompanyBySlug(db *gorm.DB, slug string) (*rpc.Company, error) {
	slug, err := util.HygienizeSlug(slug)
	if err != nil {
		log.Error("Error while evaluating company slug: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Error while evaluating company slug: %v", err)
	}

	company, err := repository.GetCompanyBySlug(db, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("Company with slug \"%s\" not found", slug)
			return nil, status.Error(codes.NotFound, "Company not found")
		}

		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	msg := company.ToMessage()
	return &msg, nil
}

func GetCompanyExists(db *gorm.DB, id string) (*wrapperspb.BoolValue, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		log.Error("\"%s\" is an invalid UUID", id)
		return nil, status.Error(codes.InvalidArgument, "Index parameter is an invalid UUID")
	}

	exists, err := repository.ExistsCompany(db, parsedId)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	return wrapperspb.Bool(exists), nil
}

func CreateCompany(db *gorm.DB, data *rpc.Company) (*rpc.Company, error) {
	if data.Id != nil {
		log.Error("New company must not have a predefined id")
		return nil, status.Error(codes.InvalidArgument, "New company must not have a predefined id")
	}

	m, err := model.Company{}.FromMessage(data)
	if err != nil {
		log.Error("Error mapping company data: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Error mapping user data: %v", err)
	}

	log.Debug("Check if slug exists")
	exists, err := repository.ExistsCompanyBySlug(db, m.Slug)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	if exists {
		log.Error("A company with this slug already exists")
		return nil, status.Errorf(codes.AlreadyExists, "A company with this slug already exists")
	}

	log.Debug("Creating company")
	created, err := repository.CreateCompany(db, m)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	go func() {
		if json, err := json.Marshal(created); err == nil {
			log.Debug("Company created: %s", json)
		}
	}()

	msg := created.ToMessage()
	return &msg, nil
}

func DisableCompany(db *gorm.DB, id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		log.Error("\"%s\" is an invalid UUID", id)
		return status.Error(codes.InvalidArgument, "Index parameter is an invalid UUID")
	}

	// If the company was soft-deleted, then it isn't queryable, even when checking
	// for its existence. So if it doesn't exist from the GORM perspective, we
	// don't need to disable it again anyway
	exists, err := repository.ExistsCompany(db, parsedId)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	if !exists {
		log.Error("Company does not exist or is already disabled")
		return status.Errorf(codes.NotFound, "Company does not exist or is already disabled")
	}

	return repository.DisableCompany(db, parsedId)
}

func UpdateCompany(db *gorm.DB, data *rpc.Company) (*rpc.Company, error) {
	if data.Id == nil {
		log.Error("Company id is missing")
		return nil, status.Error(codes.InvalidArgument, "Company id is missing")
	}

	if data.Slug != "" {
		slug, err := util.HygienizeSlug(data.Slug)
		if err != nil {
			log.Error("Error while evaluating new company slug: %v", err)
			return nil, status.Errorf(codes.InvalidArgument, "Error while evaluating new company slug: %v", err)
		}

		data.Slug = slug

		exists, err := repository.ExistsCompanyBySlug(db, data.Slug)
		if err != nil {
			log.Error("Error accessing database: %v", err)
			return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
		}

		if exists {
			log.Error("Company slug \"%s\" is already taken", data.Slug)
			return nil, status.Errorf(codes.FailedPrecondition, "This slug has already been taken")
		}
	}

	d, err := model.Company{}.FromMessage(data)
	if err != nil {
		log.Error("Error mapping company data: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Error mapping company data: %v", err)
	}

	result, err := repository.UpdateCompany(db, d)
	if err != nil {
		log.Error("Error accessing database: %v", err)
		return nil, status.Errorf(codes.Internal, "Error accessing database: %v", err)
	}

	msg := result.ToMessage()
	return &msg, nil
}
