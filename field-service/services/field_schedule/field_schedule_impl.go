package services

import (
	"context"
	"fmt"
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
	_, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = f.repository.GetFieldSchedule().Delete(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (f *FieldScheduleService) convertMonthName(inputDate string) string {
	date, err := time.Parse(time.DateOnly, inputDate)
	if err != nil {
		return ""
	}

	indonesiaMonth := map[string]string{
		"Jan": "Jan",
		"Feb": "Feb",
		"Mar": "Mar",
		"Apr": "Apr",
		"May": "Mei",
		"Jun": "Jun",
		"Jul": "Jul",
		"Aug": "Agu",
		"Sep": "Sep",
		"Oct": "Okt",
		"Nov": "Nov",
		"Dec": "Des",
	}

	formattedDate := date.Format("28 Sep")
	day := formattedDate[:3]
	month := formattedDate[3:]
	formattedDate = fmt.Sprintf("%s %s", day, indonesiaMonth[month])
	return formattedDate
}

// FindAllByIdAndDate implements IFieldScheduleService.
func (f *FieldScheduleService) FindAllFieldByIdAndDate(ctx context.Context, uuid string, date string) ([]dto.FieldScheduleForBookingReponse, error) {
	field, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	fieldSchedules, err := f.repository.GetFieldSchedule().FindAllByIdAndDate(ctx, int(field.ID), date)
	if err != nil {
		return nil, err
	}

	fieldScheduleResults := make([]dto.FieldScheduleForBookingReponse, 0, len(fieldSchedules))
	for _, fieldSchedule := range fieldSchedules {
		pricePerHour := float64(field.Field.PricePerHour)
		startTime, _ := time.Parse("12:21:01", fieldSchedule.Time.StartTime)
		endTime, _ := time.Parse("12:21:01", fieldSchedule.Time.EndTime)
		fieldScheduleResults = append(fieldScheduleResults, dto.FieldScheduleForBookingReponse{
			UUID:         fieldSchedule.UUID,
			PricePerHour: util.RupiahFormat(&pricePerHour),
			Date:         f.convertMonthName(fieldSchedule.Date.Format("2006-01-02")),
			Status:       fieldSchedule.Status.GetStatusString(),
			Time:         fmt.Sprintf("%s - %s", startTime.Format("15:04"), endTime.Format("15:04")),
		})
	}

	return fieldScheduleResults, nil
}

// FindAllWithPagination implements IFieldScheduleService.
func (f *FieldScheduleService) FindAllWithPagination(ctx context.Context, param *dto.FieldScheduleRequestParam) (*util.PaginationResult, error) {
	fieldSchedules, total, err := f.repository.GetFieldSchedule().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	fieldScheduleResults := make([]dto.FieldScheduleResponse, 0, len(fieldSchedules))
	for _, schedule := range fieldSchedules {
		fieldScheduleResults = append(fieldScheduleResults, dto.FieldScheduleResponse{
			UUID:      schedule.UUID,
			FieldName: schedule.Field.Name,
			Date:      schedule.Date.Format("2006-01-02"),
			Status:    schedule.Status,
			Time:      fmt.Sprintf("%s - %s", schedule.Time.StartTime, schedule.Time.EndTime),
			CreatedAt: schedule.CreatedAt,
			UpdateAt:  schedule.UpdatedAt,
		})
	}

	pagination := &util.PaginationParam{
		Count: total,
		Limit: param.Limit,
		Page:  param.Page,
		Data:  fieldScheduleResults,
	}

	response := util.GeneratePagination(*pagination)
	return &response, nil
}

// FindByUUID implements IFieldScheduleService.
func (f *FieldScheduleService) FindByUUID(ctx context.Context, uuid string) (*dto.FieldScheduleResponse, error) {
	fieldSchedule, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	response := dto.FieldScheduleResponse{
		UUID:      fieldSchedule.UUID,
		FieldName: fieldSchedule.Field.Name,
		Date:      fieldSchedule.Date.Format("2006-01-02"),
		Status:    fieldSchedule.Status.GetStatusString().GetStatusInt(),
		Time:      fmt.Sprintf("%s - %s", fieldSchedule.Time.StartTime, fieldSchedule.Time.EndTime),
		CreatedAt: fieldSchedule.CreatedAt,
		UpdateAt:  fieldSchedule.UpdatedAt,
	}

	return &response, nil
}

// GenereateScheduleForOneMonth implements IFieldScheduleService.
func (f *FieldScheduleService) GenereateScheduleForOneMonth(ctx context.Context, req dto.GenerateFieldScheduleForOneMonthRequest) error {
	field, err := f.repository.GetField().FindByUUID(ctx, req.FieldID)
	if err != nil {
		return err
	}

	times, err := f.repository.GetTime().FindAll(ctx)
	if err != nil {
		return err
	}

	numberOfDay := 30

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
