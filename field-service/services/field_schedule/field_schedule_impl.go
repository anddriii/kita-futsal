package services

import (
	"context"
	"time"

	"github.com/anddriii/kita-futsal/field-service/common/util"
	"github.com/anddriii/kita-futsal/field-service/constants"
	errFieldSchedule "github.com/anddriii/kita-futsal/field-service/constants/error/fieldSchedule"
	"github.com/anddriii/kita-futsal/field-service/domains/dto"
	"github.com/anddriii/kita-futsal/field-service/domains/models"
	"github.com/anddriii/kita-futsal/field-service/repositories"
	"github.com/google/uuid"
)

type FieldScheduleService struct {
	repository repositories.IRepoRegistry
}

// Create menambahkan jadwal lapangan baru berdasarkan permintaan pengguna.
func (f *FieldScheduleService) Create(ctx context.Context, req *dto.FieldScheduleRequest) error {
	// Mencari lapangan berdasarkan UUID yang diberikan dalam request
	field, err := f.repository.GetField().FindByUUID(ctx, req.FieldID)
	if err != nil {
		return err // Mengembalikan error jika lapangan tidak ditemukan
	}

	// Menyiapkan slice untuk menyimpan data jadwal lapangan yang akan dibuat
	fieldSchedules := make([]models.FieldSchedule, 0, len(req.TimeIDs))

	// Parsing tanggal dari request (format: YYYY-MM-DD)
	dateParsed, _ := time.Parse(time.DateOnly, req.Date)

	// Loop melalui setiap ID waktu yang diberikan dalam request
	for _, timeID := range req.TimeIDs {
		// Mencari data waktu berdasarkan UUID
		scheduleTime, err := f.repository.GetTime().FindByUUID(ctx, timeID)
		if err != nil {
			return err // Mengembalikan error jika data waktu tidak ditemukan
		}

		// Memeriksa apakah sudah ada jadwal untuk tanggal dan waktu tertentu di lapangan ini
		schedule, err := f.repository.GetFieldSchedule().FindByDateAndTimeId(ctx, req.Date, int(scheduleTime.ID), int(field.ID))
		if err != nil {
			return err
		}

		// Jika jadwal sudah ada, mengembalikan error bahwa jadwal sudah dipesan
		if schedule != nil {
			return errFieldSchedule.ErrFieldScheduleIsExist
		}

		// Menambahkan jadwal baru ke dalam slice fieldSchedules
		fieldSchedules = append(fieldSchedules, models.FieldSchedule{
			UUID:    uuid.New(),          // Membuat UUID baru untuk jadwal
			FieldId: field.ID,            // Mengaitkan dengan ID lapangan
			TimeId:  scheduleTime.ID,     // Mengaitkan dengan ID waktu
			Date:    dateParsed,          // Menyimpan tanggal dalam format Date
			Status:  constants.Available, // Status awal sebagai tersedia
		})
	}

	// Menyimpan semua data jadwal yang telah dibuat ke database
	err = f.repository.GetFieldSchedule().Create(ctx, fieldSchedules)
	if err != nil {
		return err // Mengembalikan error jika gagal menyimpan ke database
	}

	return nil // Mengembalikan nil jika operasi berhasil tanpa error
}

// Delete implements IFieldScheduleService.
func (f *FieldScheduleService) Delete(ctx context.Context, uuid string) error {
	panic("unimplemented")
}

// FindAllByIdAndDate implements IFieldScheduleService.
func (f *FieldScheduleService) FindAllByIdAndDate(ctx context.Context, fieldId int, date string) (dto.FieldScheduleForBookingReponse, error) {
	panic("unimplemented")
}

// FindAllWithPagination implements IFieldScheduleService.
func (f *FieldScheduleService) FindAllWithPagination(ctc context.Context, req *dto.FieldScheduleRequestParam) (util.PaginationResult, error) {
	panic("unimplemented")
}

// FindByUUID implements IFieldScheduleService.
func (f *FieldScheduleService) FindByUUID(ctx context.Context, uuid string) (dto.FieldScheduleResponse, error) {
	panic("unimplemented")
}

// GenereateScheduleForOneMonth implements IFieldScheduleService.
func (f *FieldScheduleService) GenereateScheduleForOneMonth(ctx context.Context, req dto.GenerateFieldScheduleForOneMonthRequest) error {
	panic("unimplemented")
}

// Update implements IFieldScheduleService.
func (f *FieldScheduleService) Update(ctx context.Context, uuid string, req *dto.UpdateFieldRequest) (*dto.FieldScheduleResponse, error) {
	panic("unimplemented")
}

// UpdateStatus implements IFieldScheduleService.
func (f *FieldScheduleService) UpdateStatus(ctx context.Context, req dto.UpdateStatusFieldScheduleRequest) {
	panic("unimplemented")
}

func NewFieldScheduleService(repository repositories.IRepoRegistry) IFieldScheduleService {
	return &FieldScheduleService{repository: repository}
}
