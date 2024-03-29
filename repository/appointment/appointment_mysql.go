package repository

import (
	dbCon "Halovet/driver"
	models "Halovet/models"
	"time"

	"database/sql"
	. "fmt"
	"log"
	"strconv"
)

var err error
var db *sql.DB

func init() {
	// KONEK KE DATABASE
	db, err = dbCon.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Insert(time_appointment string, doctor_name string, pet_owner_name string, pet_owner_id int, pet_type string, complaint string) (models.Appointment, bool) {
	var appointment models.Appointment

	sql_statement := "insert into appointment (appointment_time,doctor_name, pet_owner_name, pet_owner_id, pet_type, complaint_description) values (?,?,?,?,?,?)"
	row, err := db.Exec(sql_statement, time_appointment, doctor_name, pet_owner_name, pet_owner_id, pet_type, complaint)

	if err != nil {
		Println(err.Error())
		return appointment, false
	}

	id, err := row.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
	}
	appointment.AppointmentID = id
	appointment.Doctor_name = doctor_name
	appointment.Time_Appointment = time_appointment
	appointment.Pet_Owner_Name = pet_owner_name
	appointment.Pet_owner_id = pet_owner_id
	appointment.Pet_Type = pet_type
	appointment.Complaint = complaint
	appointment.CreatedAt = Sprintf(time.Now().Format("2006-01-02 15:04:05"))
	appointment.UpdatedAt = Sprintf(time.Now().Format("2006-01-02 15:04:05"))

	return appointment, true

}

func FindAppointmentbyUserID(id string) ([]models.Appointment, bool) {

	var Appointment models.Appointment
	var Appointments []models.Appointment
	var nullhandler string

	userid, err := strconv.Atoi(id)
	// Println(id)
	if err != nil {
		Println("format ID salah")
		return Appointments, false
	}

	sqlStatement := "select * from appointment where pet_owner_id = ?"
	results, err := db.Query(sqlStatement, userid)
	if err != nil {
		panic(err.Error())
		return Appointments, false
	}
	for results.Next() {
		err = results.Scan(&Appointment.AppointmentID,
			&Appointment.Time_Appointment,
			&Appointment.Doctor_name,
			&Appointment.Pet_Owner_Name,
			&Appointment.Pet_Type,
			&Appointment.Complaint,
			&Appointment.IsPaid,
			&Appointment.CreatedAt,
			&Appointment.UpdatedAt,
			&Appointment.Pet_owner_id,
			&Appointment.PhotoPath)
		if err != nil {
			panic(err.Error())
			return Appointments, false
		} else {
			if nullhandler == "0" {
				Appointment.PhotoPath = "-"
			} else {
				Appointment.PhotoPath = nullhandler
			}
			
		}
		Appointments = append(Appointments, Appointment)
	}
	return Appointments, true
}

func FindAllAppointment(limitstart string, limit string) ([]models.Appointment, int, error) {
	var Appointment models.Appointment
	var Appointments []models.Appointment

	realLimitStart, err := strconv.Atoi(limitstart)
	if err != nil {
		Println("format limit salah")
		return Appointments, 0, err
	}
	realLimit, err := strconv.Atoi(limit)
	if err != nil {
		Println("format limit salah")
		return Appointments, 0, err
	}

	sqlStatement := "select * from appointment order by created_at limit ?, ?"
	results, err := db.Query(sqlStatement, realLimitStart, realLimit)
	if err != nil {
		panic(err.Error())
		return Appointments, 0, err
	}

	var nullhandler string
	for results.Next() {
		err = results.Scan(&Appointment.AppointmentID,
			&Appointment.Time_Appointment,
			&Appointment.Doctor_name,
			&Appointment.Pet_Owner_Name,
			&Appointment.Pet_Type,
			&Appointment.Complaint,
			&Appointment.IsPaid,
			&Appointment.CreatedAt,
			&Appointment.UpdatedAt,
			&Appointment.Pet_owner_id,
			&nullhandler,
		)
		if err != nil {
			panic(err.Error())
		} else {
			if nullhandler == "0" {
				Appointment.PhotoPath = "-"
			} else {
				Appointment.PhotoPath = nullhandler
			}
			
		}
		Appointments = append(Appointments, Appointment)

	}

	var count int

	err = db.QueryRow("SELECT COUNT(*) FROM appointment").Scan(&count)

	if err != nil {
		log.Fatal(err)
		return Appointments, 0, err
	}

	return Appointments, count, nil
}

func FindbyID(id string) (models.Appointment, bool) {
	var appointment models.Appointment
	appointmentid, err := strconv.Atoi(id)
	if err != nil {
		Println("format ID salah")
	}
	var nullhandler string
	sql_statement := "select * from appointment where id = ?"
	err = db.QueryRow(sql_statement, appointmentid).
		Scan(&appointment.AppointmentID,
			&appointment.Time_Appointment,
			&appointment.Doctor_name,
			&appointment.Pet_Owner_Name,
			&appointment.Pet_Type,
			&appointment.Complaint,
			&appointment.IsPaid,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
			&appointment.Pet_owner_id,
			&nullhandler,
		)

	if err != nil {
		Println(err.Error())
		return appointment, false
	} else {
		// Appointments = append(Appointments, Appointment)
		if nullhandler == "0" {
			appointment.PhotoPath = "-"
		} else {
			appointment.PhotoPath = nullhandler
		}
		
	}

	return appointment, true
}

// func FindAll()

func Remove(id string) bool {

	appointmentid, err := strconv.Atoi(id)

	if err != nil {
		Println("format ID salah")
		return false
	}
	sql_statement := "delete from appointment where id = ?"
	_, err = db.Exec(sql_statement, appointmentid)

	if err != nil {
		Println(err.Error())
		return false
	}
	return true
}

func Update(id string, doctor_name string, pet_type string, complaint string, appointment_time string) bool {
	// var appointment models.Appointment
	appointmentid, err := strconv.Atoi(id)
	if err != nil {
		Println("format ID salah")
	}
	timeNow := Sprintf(time.Now().Format("2006-01-02 15:04:05"))
	sql_statement := "update appointment set doctor_name = ?, pet_type = ?, complaint_description = ?, appointment_time = ?, updated_at = ? where id = ?"
	_, err = db.Exec(sql_statement, doctor_name, pet_type, complaint, appointment_time, timeNow, appointmentid)

	if err != nil {
		Println(err.Error())
		return false
	}

	return true

}

func ValidatePayment(appointmentid int) bool {
	// appointmentid, err := strconv.Atoi(id)
	// if err != nil {
	// 	Println("format ID salah")
	// }

	sqlStatement := "update appointment set is_paid = 1 where id = ?"
	row, err := db.Exec(sqlStatement, appointmentid)

	if err != nil {
		Println(err.Error())
		return false
	}

	count, err := row.RowsAffected()
	if count == 0 {
		return false
	}
	return true
}

func CheckValidation(appointmentid int) bool {

	// appointmentid, err := strconv.Atoi(id)
	// if err != nil {
	// 	Println("format ID salah")
	// }

	var statusPembayaran int
	sqlStatement := "select is_paid appointment where id = ?"
	err = db.QueryRow(sqlStatement, appointmentid).Scan(&statusPembayaran)

	if statusPembayaran == 1 {
		return false
	}
	return true
}
