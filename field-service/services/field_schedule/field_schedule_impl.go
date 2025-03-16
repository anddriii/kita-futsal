package services

import (
	"context"
	"fmt"
	"time"

	"github.com/anddriii/kita-futsal/field-service/common/util"
	"github.com/anddriii/kita-futsal/field-service/constants"
	errFieldSchedule "github.com/anddriii/kita-futsal/field-service/constants/error/field_schedule"
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
			return errFieldSchedule.ErrFieldScheduleExist
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

// FindAllFieldByIdAndDate retrieves all field schedules by field UUID and date.
// It returns a list of available field schedules for booking on the given date.
func (f *FieldScheduleService) FindAllFieldByIdAndDate(ctx context.Context, uuid string, date string) ([]dto.FieldScheduleForBookingReponse, error) {
	// Retrieve field details using UUID
	field, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// Fetch all schedules related to the field ID and specific date
	fieldSchedules, err := f.repository.GetFieldSchedule().FindAllByIdAndDate(ctx, int(field.ID), date)
	if err != nil {
		return nil, err
	}

	// Prepare response slice
	fieldScheduleResults := make([]dto.FieldScheduleForBookingReponse, 0, len(fieldSchedules))
	for _, fieldSchedule := range fieldSchedules {
		pricePerHour := float64(field.Field.PricePerHour)
		startTime, _ := time.Parse("15:04:05", fieldSchedule.Time.StartTime)
		endTime, _ := time.Parse("15:04:05", fieldSchedule.Time.EndTime)

		// Append formatted data to response list
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

// FindAllWithPagination retrieves all field schedules with pagination.
// It returns a paginated response containing field schedules.
func (f *FieldScheduleService) FindAllWithPagination(ctx context.Context, param *dto.FieldScheduleRequestParam) (*util.PaginationResult, error) {
	// Fetch paginated schedules based on request parameters
	fieldSchedules, total, err := f.repository.GetFieldSchedule().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	// Prepare response slice
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

	// Generate pagination metadata
	pagination := &util.PaginationParam{
		Count: total,
		Limit: param.Limit,
		Page:  param.Page,
		Data:  fieldScheduleResults,
	}

	// Generate and return the paginated response
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

// GenereateScheduleForOneMonth membuat jadwal otomatis selama satu bulan untuk lapangan tertentu.
// Fungsi ini akan mengambil semua slot waktu yang tersedia dan membuat jadwal baru untuk setiap hari dalam 30 hari ke depan.
// Jika jadwal sudah ada untuk tanggal dan waktu tertentu, fungsi akan mengembalikan error untuk menghindari duplikasi.
func (f *FieldScheduleService) GenereateScheduleForOneMonth(ctx context.Context, req dto.GenerateFieldScheduleForOneMonthRequest) error {
	// Mengambil data lapangan berdasarkan FieldID yang diberikan dalam request.
	field, err := f.repository.GetField().FindByUUID(ctx, req.FieldID)
	if err != nil {
		return err // Jika lapangan tidak ditemukan, return error.
	}

	// Mengambil semua slot waktu yang tersedia dari database.
	times, err := f.repository.GetTime().FindAll(ctx)
	if err != nil {
		return err // Jika gagal mengambil data waktu, return error.
	}

	// Menentukan jumlah hari untuk pembuatan jadwal (30 hari ke depan).
	numberOfDay := 30

	// Membuat slice kosong untuk menyimpan jadwal yang akan dibuat.
	// Kapasitas awalnya dihitung berdasarkan jumlah slot waktu * jumlah hari.
	fieldSchedules := make([]models.FieldSchedule, 0, len(times)*numberOfDay)

	// Menentukan tanggal mulai pembuatan jadwal (besok dari hari ini).
	now := time.Now().Add(time.Duration(1) * 24 * time.Hour)

	// Looping selama 30 hari ke depan.
	for i := range numberOfDay {
		// Menentukan tanggal saat ini dalam iterasi.
		currentDate := now.AddDate(0, 0, i)

		// Looping untuk setiap slot waktu yang tersedia.
		for _, item := range times {
			// Mengecek apakah jadwal untuk tanggal dan waktu tertentu sudah ada di database.
			schedule, err := f.repository.GetFieldSchedule().FindByDateAndTimeId(ctx, currentDate.Format(time.DateOnly), int(item.ID), int(field.ID))
			if err != nil {
				return err // Jika terjadi error dalam pencarian, return error.
			}

			// Jika jadwal sudah ada untuk tanggal dan waktu tersebut, return error untuk menghindari duplikasi.
			if schedule != nil {
				return errFieldSchedule.ErrFieldScheduleExist
			}

			// Menambahkan jadwal baru ke dalam slice fieldSchedules.
			fieldSchedules = append(fieldSchedules, models.FieldSchedule{
				UUID:    uuid.New(),          // Generate UUID baru untuk jadwal.
				FieldId: field.ID,            // Menyimpan ID lapangan.
				TimeId:  item.ID,             // Menyimpan ID slot waktu.
				Date:    currentDate,         // Menyimpan tanggal jadwal.
				Status:  constants.Available, // Status jadwal diatur sebagai "tersedia".
			})
		}
	}

	// Melakukan batch insert untuk menyimpan semua jadwal baru ke dalam database.
	err = f.repository.GetFieldSchedule().Create(ctx, fieldSchedules)
	if err != nil {
		return err // Jika terjadi error saat menyimpan, return error.
	}

	return nil // Jika berhasil, return nil (tidak ada error).
}

func (f *FieldScheduleService) Update(ctx context.Context, uuid string, req *dto.UpdateFieldScheduleRequest) (*dto.FieldScheduleResponse, error) {
	// Mencari jadwal lapangan berdasarkan UUID yang diberikan
	fieldSchedule, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err // Jika tidak ditemukan, kembalikan error
	}

	// Mencari data waktu berdasarkan TimeID yang diberikan dalam request
	scheduleTime, err := f.repository.GetTime().FindByUUID(ctx, req.TimeID)
	if err != nil {
		return nil, err // Jika tidak ditemukan, kembalikan error
	}

	// Mengecek apakah sudah ada jadwal dengan tanggal dan waktu yang sama di lapangan yang sama
	isTimeExist, err := f.repository.GetFieldSchedule().FindByDateAndTimeId(ctx, req.Date, int(scheduleTime.ID), int(fieldSchedule.FieldId))
	if err != nil {
		return nil, err
	}

	// Jika tanggal diubah, cek apakah sudah ada jadwal di tanggal dan waktu yang sama
	if isTimeExist != nil && req.Date != fieldSchedule.Date.Format(time.DateOnly) {
		checkDate, err := f.repository.GetFieldSchedule().FindByDateAndTimeId(ctx, req.Date, int(scheduleTime.ID), int(fieldSchedule.FieldId))
		if err != nil {
			return nil, err
		}

		// Jika jadwal sudah ada di tanggal baru yang dipilih, kembalikan error
		if checkDate != nil {
			return nil, errFieldSchedule.ErrFieldScheduleExist
		}
	}

	// Parsing string tanggal ke format time.Time
	dateParsed, _ := time.Parse(time.DateOnly, req.Date)

	// Melakukan update jadwal di database
	fieldResult, err := f.repository.GetFieldSchedule().Update(ctx, uuid, &models.FieldSchedule{
		Date:   dateParsed,
		TimeId: scheduleTime.ID,
	})
	if err != nil {
		return nil, err
	}

	// Membentuk response DTO untuk dikirimkan ke client
	response := dto.FieldScheduleResponse{
		UUID:         fieldResult.UUID,
		FieldName:    fieldResult.Field.Name,
		Date:         fieldResult.Date.Format(time.DateOnly),
		PricePerHour: int(fieldSchedule.Status.GetStatusString().GetStatusInt()), // Mengonversi status ke harga per jam
		Time:         fmt.Sprintf("%s - %s", scheduleTime.StartTime, scheduleTime.EndTime),
		CreatedAt:    fieldResult.CreatedAt,
		UpdateAt:     fieldResult.UpdatedAt,
	}

	return &response, nil // Mengembalikan hasil update ke client
}

// UpdateStatus implements IFieldScheduleService.
func (f *FieldScheduleService) UpdateStatus(ctx context.Context, req dto.UpdateStatusFieldScheduleRequest) error {
	for _, item := range req.FieldScheduleIDs {
		_, err := f.repository.GetFieldSchedule().FindByUUID(ctx, item)
		if err != nil {
			return err
		}

		err = f.repository.GetFieldSchedule().UpdateStatus(ctx, constants.Booked, item)
		if err != nil {
			return err
		}

	}
	return nil
}

func NewFieldScheduleService(repository repositories.IRepoRegistry) IFieldScheduleService {
	return &FieldScheduleService{repository: repository}
}
