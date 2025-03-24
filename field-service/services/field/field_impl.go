package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"time"

	"github.com/anddriii/kita-futsal/field-service/common/gcs"
	"github.com/anddriii/kita-futsal/field-service/common/util"
	errCons "github.com/anddriii/kita-futsal/field-service/constants/error"
	"github.com/anddriii/kita-futsal/field-service/domains/dto"
	"github.com/anddriii/kita-futsal/field-service/domains/models"
	"github.com/anddriii/kita-futsal/field-service/repositories"
	"github.com/google/uuid"
)

type FieldService struct {
	repository repositories.IRepoRegistry
	gcs        gcs.IGCSClient
}

func NewFieldService(repository repositories.IRepoRegistry, gcs gcs.IGCSClient) IFieldService {
	return &FieldService{
		repository: repository,
		gcs:        gcs,
	}
}

func (f *FieldService) GetAllWithPagination(ctx context.Context, param *dto.FieldRequestParam) (*util.PaginationResult, error) {
	fields, total, err := f.repository.GetField().FindALlWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	// Mengonversi data field ke format FieldResponse
	fieldResults := make([]dto.FieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldResults = append(fieldResults, dto.FieldResponse{
			UUID:         field.UUID,
			Code:         field.Code,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Image,
			CreatedAt:    field.CreatedAt,
			UpdateAt:     field.UpdatedAt,
		})
	}

	pagination := &util.PaginationParam{
		Count: total,
		Page:  param.Page,
		Limit: param.Limit,
		Data:  fieldResults,
	}

	response := util.GeneratePagination(*pagination)
	return &response, nil
}

func (f *FieldService) GetAllWithoutPagination(ctx context.Context) ([]dto.FieldResponse, error) {
	fields, err := f.repository.GetField().FindAllWithoutPagination(ctx)
	if err != nil {
		return nil, err
	}

	fieldResults := make([]dto.FieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldResults = append(fieldResults, dto.FieldResponse{
			UUID:         field.UUID,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Image,
		})
	}

	return fieldResults, nil
}

func (f *FieldService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldResponse, error) {
	field, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	pricePerHour := float64(field.PricePerHour)
	fieldResult := dto.FieldResponse{
		UUID:         field.UUID,
		Code:         field.Code,
		Name:         field.Name,
		PricePerHour: util.RupiahFormat(&pricePerHour),
		Images:       field.Image,
		CreatedAt:    field.CreatedAt,
		UpdateAt:     field.UpdatedAt,
	}

	return &fieldResult, nil
}

func (f *FieldService) validateUpload(images []multipart.FileHeader) error {
	if len(images) == 0 {
		return errCons.ErrInvalidUploadFile
	}

	for _, image := range images {
		if image.Size > 5*1024*1024 { // 5MB batas ukuran
			return errCons.ErrSizeTooBig
		}
	}

	return nil
}

// processAndUploadImage membuka file gambar, membaca isinya ke dalam buffer,
// dan mengunggahnya ke Google Cloud Storage.
//
// Parameter:
// - ctx: Context untuk operasi asinkron
// - image: File gambar yang akan diunggah
//
// Return:
// - URL dari file yang diunggah
// - Error jika terjadi kesalahan
func (f *FieldService) processAndUploadImage(ctx context.Context, image multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Membaca isi file ke dalam buffer
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	if err != nil {
		return "", err
	}

	// Membuat nama file unik berdasarkan timestamp
	filename := fmt.Sprintf("images/%s-%s-%s", time.Now().Format("20060102150405"), image.Filename, path.Ext(image.Filename))

	// Mengunggah file ke GCS
	url, err := f.gcs.UploadFile(ctx, filename, buffer.Bytes())
	if err != nil {
		return "", err
	}

	return url, nil
}

// uploadImage memvalidasi daftar file gambar, memproses, dan mengunggahnya ke GCS.
//
// Parameter:
// - ctx: Context untuk operasi asinkron
// - images: Daftar file gambar yang akan diunggah
//
// Return:
// - Daftar URL dari gambar yang berhasil diunggah
// - Error jika ada kegagalan dalam proses upload
func (f *FieldService) uploadImage(ctx context.Context, images []multipart.FileHeader) ([]string, error) {
	err := f.validateUpload(images)
	if err != nil {
		return nil, err
	}

	// Menampung URL hasil upload
	urls := make([]string, 0, len(images))
	for _, image := range images {
		url, err := f.processAndUploadImage(ctx, image)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (f *FieldService) Create(ctx context.Context, req *dto.FieldRequest) (*dto.FieldResponse, error) {
	// Path direktori utama untuk local

	photo, err := util.UploadImageLocal(req.Images)
	if err != nil {
		return nil, err
	}

	// upload image for GCPs
	// imageUrl, err := f.uploadImage(ctx, req.Images)
	// if err != nil {
	// 	return nil, err
	// }

	field, err := f.repository.GetField().Create(ctx, &models.Field{
		Code:         req.Code,
		Name:         req.Name,
		PricePerHour: req.PricePerHour,
		Image:        photo,
	})
	if err != nil {
		return nil, err
	}

	response := dto.FieldResponse{
		UUID:         field.UUID,
		Code:         field.Code,
		Name:         field.Name,
		PricePerHour: field.PricePerHour,
		Images:       field.Image,
		CreatedAt:    field.CreatedAt,
		UpdateAt:     field.UpdatedAt,
	}

	return &response, nil
}

func (f *FieldService) Update(ctx context.Context, uuidParam string, req *dto.UpdateFieldRequest) (*dto.FieldResponse, error) {
	field, err := f.repository.GetField().FindByUUID(ctx, uuidParam)
	if err != nil {
		return nil, err
	}

	var imageUrls []string
	if req.Images == nil {
		imageUrls = field.Image // Gunakan gambar lama jika tidak ada gambar baru
	} else {
		imageUrls, err = f.uploadImage(ctx, req.Images) // Upload gambar baru jika tersedia
		if err != nil {
			return nil, err
		}
	}

	fieldResult, err := f.repository.GetField().Update(ctx, uuidParam, &models.Field{
		Code:         req.Code,
		Name:         req.Name,
		PricePerHour: req.PricePerHour,
		Image:        imageUrls,
	})
	if err != nil {
		return nil, err
	}

	uuidParsed, _ := uuid.Parse(uuidParam)
	response := dto.FieldResponse{
		UUID:         uuidParsed,
		Code:         fieldResult.Code,
		Name:         fieldResult.Name,
		PricePerHour: fieldResult.PricePerHour,
		Images:       fieldResult.Image,
		CreatedAt:    fieldResult.CreatedAt,
		UpdateAt:     fieldResult.UpdatedAt,
	}

	return &response, nil
}

func (f *FieldService) Delete(ctx context.Context, uuid string) error {
	// Cek apakah field dengan UUID tersebut ada
	_, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	// Hapus field dari repository
	err = f.repository.GetField().Delete(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}
