package database

import (
	"app-api/config"
	"app-api/internal/models"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var AllModels = []interface{}{
	&models.Note{},
	&models.Comment{},
}

func InitDatabase(config *config.DatabaseConfig) error {
	var err error
	// Use := only for the first declaration, not for assigning to the global DB variable
	fmt.Printf("Connecting to %s", config.GetDatabaseString())
	DB, err = gorm.Open(postgres.Open(config.GetDatabaseString()), &gorm.Config{})

	if err != nil {
		return err
	}

	posgresDB, err := DB.DB()

	if err != nil {
		return err
	}

	posgresDB.SetMaxIdleConns(10)
	posgresDB.SetMaxOpenConns(100)

	return nil
}

func Migrate(config *config.Config, logger *zap.Logger) error {
	if DB == nil {
		return errors.New("database is not initialized")
	}

	if err := DB.AutoMigrate(&models.Role{}); err != nil {
		return err
	}

	if err := seedRoles(DB); err != nil {
		return err
	}

	if DB.Migrator().HasTable(&models.User{}) && DB.Migrator().HasColumn(&models.User{}, "role_id") {
		if err := assignDefaultRole(DB); err != nil {
			return err
		}
	}

	if err := DB.AutoMigrate(&models.User{}); err != nil {
		fmt.Printf("Failed to migrate model: %T, error: %v\n", &models.User{}, err)
		return err
	}

	if err := assignDefaultRole(DB); err != nil {
		return err
	}

	for _, model := range AllModels {
		if err := DB.AutoMigrate(model); err != nil {
			fmt.Printf("Failed to migrate model: %T, error: %v\n", model, err)
			return err
		}
	}

	initAdmin(DB, config)

	return nil
}

func CloseDB() {
	if DB != nil {
		posgresDB, error := DB.DB()
		if error != nil {
			return
		}
		posgresDB.Close()
	}
}

func GetDB() *gorm.DB {
	return DB
}

func seedRoles(db *gorm.DB) error {
	defaultRoles := []string{"member", "admin"}

	for _, roleName := range defaultRoles {
		var role models.Role
		if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				role = models.Role{Name: roleName}
				if err := db.Create(&role).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

func assignDefaultRole(db *gorm.DB) error {
	var memberRole models.Role
	if err := db.Where("name = ?", "member").First(&memberRole).Error; err != nil {
		return err
	}

	return db.Model(&models.User{}).
		Where("role_id = 0 OR role_id IS NULL").
		Update("role_id", memberRole.ID).Error
}

func initAdmin(db *gorm.DB, config *config.Config) error {
	var adminUser models.User
	if err := db.Where("email = ?", config.DefaultAdmin.Email).First(&adminUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var adminRole models.Role
			if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
				return err
			}
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(config.DefaultAdmin.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			adminUser = models.User{
				Email:        config.DefaultAdmin.Email,
				Name:         config.DefaultAdmin.Name,
				PasswordHash: string(hashedPassword),
				RoleID:       adminRole.ID,
			}
			if err := db.Create(&adminUser).Error; err != nil {
				return err
			}

			// simulate some workload
			note := models.Note{
				Title:       "Welcome Note",
				Description: "This is a welcome note",
				Content:     "Welcome to the platform!",
				UserID:      adminUser.ID,
				Public:      true,
			}

			if err := db.Create(&note).Error; err != nil {
				return err
			}

			comment := models.Comment{
				NoteID: note.ID,
				UserID: adminUser.ID,
				Text:   config.Flag,
			}
			if err := db.Create(&comment).Error; err != nil {
				return err
			}

			logger := zap.L()
			logger.Info("Admin user created", zap.String("email", adminUser.Email))

		} else {
			return err
		}
	}
	return nil
}
