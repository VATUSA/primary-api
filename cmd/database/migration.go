package main

import (
	"fmt"
	"github.com/VATUSA/primary-api/pkg/config"
	"github.com/VATUSA/primary-api/pkg/constants"
	"github.com/VATUSA/primary-api/pkg/database"
	"github.com/VATUSA/primary-api/pkg/database/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"os"
	"sync"
	"time"
)

type Facility struct {
	ID     string `gorm:"column:id" json:"id"`
	Name   string `gorm:"column:name" json:"name"`
	URL    string `gorm:"column:url" json:"url"`
	ATM    uint   `gorm:"column:atm" json:"atm"`
	DATM   uint   `gorm:"column:datm" json:"datm"`
	TA     uint   `gorm:"column:ta" json:"ta"`
	EC     uint   `gorm:"column:ec" json:"ec"`
	FE     uint   `gorm:"column:fe" json:"fe"`
	WM     uint   `gorm:"column:wm" json:"wm"`
	APIKey string `gorm:"column:apikey" json:"apikey"`
}

type User struct {
	CID          uint                 `gorm:"column:cid" json:"cid"`
	FirstName    string               `gorm:"column:fname" json:"fname"`
	LastName     string               `gorm:"column:lname" json:"lname"`
	Email        string               `gorm:"column:email" json:"email"`
	Rating       int                  `gorm:"column:rating" json:"rating"`
	DiscordID    string               `gorm:"column:discord_id" json:"discord_id"`
	HomeFacility constants.FacilityID `gorm:"column:facility" json:"facility"`
	JoinDate     time.Time            `gorm:"column:facility_join" json:"facility_join"`
	LastCertSync time.Time            `gorm:"column:last_cert_sync" json:"last_cert_sync"`
	LastLogin    time.Time            `gorm:"column:lastactivity" json:"lastactivity"`
	CreatedAt    time.Time            `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time            `gorm:"column:updated_at" json:"updated_at"`

	TransferOverride bool `gorm:"column:flag_xferOverride" json:"transfer_override"`
	NoStaffRole      bool `gorm:"column:flag_preventStaffAssign" json:"no_staff_role"`
}

type Visiting struct {
	CID       uint      `gorm:"column:cid" json:"cid"`
	Facility  string    `gorm:"column:facility" json:"facility"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

type ActionLog struct {
	From      uint      `gorm:"column:from" json:"from"`
	To        uint      `gorm:"column:to" json:"to"`
	Log       string    `gorm:"column:log" json:"log"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type RatingChange struct {
	CID       uint      `gorm:"column:cid" json:"cid"`
	OldRating int       `gorm:"column:from" json:"old_rating"`
	NewRating int       `gorm:"column:to" json:"new_rating"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	CreatedBy uint      `gorm:"column:created_by" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type Role struct {
	CID       uint      `gorm:"column:cid" json:"cid"`
	Facility  string    `gorm:"column:facility" json:"facility"`
	Role      string    `gorm:"column:role" json:"role"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func main() {
	_ = godotenv.Load(".env")

	dbCfg := config.NewDBConfig()
	database.DB = database.Connect(dbCfg)

	oldDb := &config.DBConfig{
		Host:        os.Getenv("OLD_DB_HOST"),
		Port:        os.Getenv("OLD_DB_PORT"),
		User:        os.Getenv("OLD_DB_USER"),
		Password:    os.Getenv("OLD_DB_PASSWORD"),
		Database:    os.Getenv("OLD_DB_DATABASE"),
		LoggerLevel: os.Getenv("OLD_DB_LOGGER_LEVEL"),
	}
	oldDbConn := database.Connect(oldDb)

	fmt.Println("Beginning Migration...")

	// Migrate old database to new database
	// 1. Migrate facilities
	MigrateFacilities(oldDbConn)
	// 2. Migrate users
	MigrateUsers(oldDbConn)

	fmt.Println("Migration Complete.")
}

func MigrateFacilities(oldDbConn *gorm.DB) {
	fmt.Sprintln("Migrating facilities")

	var facilities []Facility
	oldDbConn.Table("facilities").Find(&facilities)

	for _, facility := range facilities {
		newFacility := &models.Facility{
			ID:     constants.FacilityID(facility.ID),
			Name:   facility.Name,
			URL:    facility.URL,
			APIKey: facility.APIKey,
		}

		if err := newFacility.Create(); err != nil {
			fmt.Printf("Error creating facility: %s\n", err.Error())
			continue
		}
	}

	fmt.Sprintln("Done migrating facilities")
}

func MigrateUsers(oldDbConn *gorm.DB) {
	var users []User
	fmt.Println("Getting users from old database")
	oldDbConn.Table("controllers").Find(&users)
	fmt.Printf("Got %d users from old database\n", len(users))

	var wg sync.WaitGroup
	userChannel := make(chan User)

	// Worker function to migrate a single user
	worker := func() {
		defer wg.Done()
		for user := range userChannel {
			// 1. Create new user in new database
			newUser := &models.User{
				CID:              user.CID,
				FirstName:        user.FirstName,
				LastName:         user.LastName,
				Email:            user.Email,
				PreferredOIs:     fmt.Sprintf("%s%s", user.FirstName[:1], user.LastName[:1]),
				PilotRating:      0,
				ControllerRating: constants.ATCRating(user.Rating),
				DiscordID:        user.DiscordID,
				LastLogin:        user.LastLogin,
				LastCertSync:     time.Now(),
				CreatedAt:        user.CreatedAt,
				UpdatedAt:        user.UpdatedAt,
			}

			if err := newUser.Create(); err != nil {
				fmt.Printf("Error creating user: %s\n", err.Error())
				continue
			}

			// 2. Migrate user flags
			userFlag := &models.UserFlag{
				CID:                  user.CID,
				NoStaffRole:          user.NoStaffRole,
				UsedTransferOverride: user.TransferOverride,
			}

			if err := userFlag.Create(); err != nil {
				fmt.Printf("Error creating user flag: %s\n", err.Error())
				continue
			}

			// 3. Migrate user roster
			roster := &models.Roster{
				CID:       user.CID,
				Facility:  constants.FacilityID(user.HomeFacility),
				CreatedAt: user.JoinDate,
				OIs:       newUser.PreferredOIs,
				Home:      true,
				Visiting:  false,
				Status:    "Active",
				DeletedAt: nil,
			}

			if err := roster.Create(); err != nil {
				fmt.Printf("Error creating roster: %s\n", err.Error())
				continue
			}

			// 4. Migrate user visiting
			var visiting []Visiting
			oldDbConn.Table("visits").Where("cid = ?", user.CID).Find(&visiting)

			for _, visit := range visiting {
				roster := &models.Roster{
					CID:       visit.CID,
					Facility:  constants.FacilityID(visit.Facility),
					CreatedAt: visit.CreatedAt,
					Home:      false,
					Visiting:  true,
					Status:    "Active",
					DeletedAt: nil,
				}

				if err := roster.Create(); err != nil {
					fmt.Printf("Error creating visiting roster: %s\n", err.Error())
					continue
				}
			}

			// 5. Migrate user action logs
			var actionLogs []ActionLog
			oldDbConn.Table("action_log").Where("`to` = ?", user.CID).Find(&actionLogs)

			for _, actionLog := range actionLogs {
				newActionLog := &models.ActionLogEntry{
					CID:       actionLog.To,
					Entry:     actionLog.Log,
					CreatedAt: actionLog.CreatedAt,
					CreatedBy: "",
					UpdatedAt: actionLog.UpdatedAt,
					UpdatedBy: "",
				}

				if actionLog.From == 0 {
					newActionLog.CreatedBy = "System"
				} else {
					newActionLog.CreatedBy = fmt.Sprintf("%d", actionLog.From)
				}

				if err := newActionLog.Create(); err != nil {
					fmt.Printf("Error creating action log: %s\n", err.Error())
					continue
				}
			}

			// 6. Migrate user rating changes
			var ratingChanges []RatingChange
			oldDbConn.Table("promotions").Where("cid = ?", user.CID).Find(&ratingChanges)
			for _, ratingChange := range ratingChanges {
				rc := &models.RatingChange{
					CID:          ratingChange.CID,
					OldRating:    constants.ATCRating(ratingChange.OldRating),
					NewRating:    constants.ATCRating(ratingChange.NewRating),
					CreatedAt:    ratingChange.CreatedAt,
					CreatedByCID: ratingChange.CreatedBy,
					UpdatedAt:    ratingChange.UpdatedAt,
				}

				if err := rc.Create(); err != nil {
					fmt.Printf("Error creating rating change: %s\n", err.Error())
					continue
				}
			}

			// 7. Migrate user roles
			var roles []Role
			oldDbConn.Table("roles").Where("cid = ?", user.CID).Find(&roles)
			for _, role := range roles {
				if !constants.FacilityID(role.Facility).IsValidFacility() {
					fmt.Printf("Facility %s is not a valid facility", role.Facility)
					continue
				}
				newRole := &models.UserRole{
					CID:        role.CID,
					FacilityID: constants.FacilityID(role.Facility),
					CreatedAt:  role.CreatedAt,
				}

				if role.Role == "EC" {
					newRole.RoleID = constants.AssistantEventCoordinator
				} else if role.Role == "FE" {
					newRole.RoleID = constants.AssistantFacilityEngineer
				} else if role.Role == "WM" {
					newRole.RoleID = constants.AssistantWebMasterRole
				} else if role.Role == "TA" {
					newRole.RoleID = constants.TrainingAdministratorRole
				} else if role.Role == "DATM" {
					newRole.RoleID = constants.DeputyAirTrafficManagerRole
				} else if role.Role == "ATM" {
					newRole.RoleID = constants.AirTrafficManagerRole
				} else if role.Role == "INS" {
					newRole.RoleID = constants.InstructorRole
				} else if role.Role == "MTR" {
					newRole.RoleID = constants.MentorRole
				} else {
					if len(role.Role) == 3 && role.Role[0:2] == "US" {
						newRole.RoleID = constants.RoleID(fmt.Sprintf("USA%s", role.Role[2:3]))
					} else {
						if constants.RoleID(role.Role).IsValidRole() {
							newRole.RoleID = constants.RoleID(role.Role)
						} else {
							fmt.Printf("Role %s is not a valid role\n", role.Role)
							continue
						}
					}
				}

				rost, err := models.GetRosterByFacilityAndCID(newRole.FacilityID, newRole.CID)
				if err != nil {
					newRole.RosterID = 0
				} else {
					newRole.RosterID = rost.ID
				}

				if err := newRole.Create(); err != nil {
					fmt.Printf("Error creating role: %s\n", err.Error())
					continue
				}
			}
		}
	}

	// Spin up worker goroutines
	for i := 0; i < 25; i++ {
		wg.Add(1)
		go worker()
	}

	// Feed users to worker goroutines through channel
	for _, user := range users {
		userChannel <- user
	}
	close(userChannel)

	wg.Wait()

	fmt.Println("Done migrating users")
}
