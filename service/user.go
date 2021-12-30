package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/mine/just-projecting/domain/model"
	"github.com/mine/just-projecting/domain/repository"
	log "github.com/mine/just-projecting/pkg/logger"
)

type UserService interface {
	Find() ([]*model.User, error)
	FindOne() (*model.User, error)
	CountUsersUsecase(userID int64) (result string, err error)
	ListMerchantOmzet(c *fiber.Ctx, date string, username string) (res []model.MerchantOmzet, err error)
	ListOutletOmzet(c *fiber.Ctx, date string, username string) (res []model.Outlet, err error)
	LoginUsecase(c *fiber.Ctx, form model.Login) (token string, user model.Auth, err error)
	GenerateToken(c *fiber.Ctx, form model.Login) (token string, err error)
}

type userService struct {
	// should be private
	UserRepository repository.UserRepository
	// cache, config, db transaction etc inject here
}

func NewUserService() *userService {
	return &userService{}
}

func (s *userService) SetUserRepository(userRepo repository.UserRepository) *userService {
	s.UserRepository = userRepo
	return s
}

func (s *userService) Validate() *userService {
	if s.UserRepository == nil {
		panic("handler need user repository")
	}
	return s
}

func (s *userService) Find() ([]*model.User, error) {
	/**
	 * business logic here
	 */
	return s.UserRepository.Find()
}

func (s *userService) FindOne() (*model.User, error) {
	return s.UserRepository.FindOne()
}

func (s *userService) CountUsersUsecase(userID int64) (result string, err error) {
	result, err = s.UserRepository.CountUserRepository(userID)
	if err != nil {
		panic("failed to call Repo")
	}

	return result, err
}

func (s *userService) ListMerchantOmzet(c *fiber.Ctx, date string, username string) (res []model.MerchantOmzet, err error) {
	l := log.GetLoggerContext(c.Context(), "service", "ListMerchantOmzet")
	l.Infof("Start Execute ListMerchantOmzet")

	listMerchant, err := s.UserRepository.GetListMerchant(username)
	if err != nil {
		l.Errorf("Error when do GetListMerchant")
	}

	res, err = s.UserRepository.ListMerchantOmzet(date, username)
	if err != nil {
		l.Errorf("Error when do ListMerchantOmzet")
	}

	mapOfListMerchant := make(map[string]bool)
	for _, v := range res {
		key := fmt.Sprintf("%d", v.MerchantID)
		mapOfListMerchant[key] = true
	}

	for _, v := range listMerchant {
		key := fmt.Sprintf("%d", v.MerchantID)
		_, found := mapOfListMerchant[key]
		if !found {
			res = append(res, model.MerchantOmzet{
				MerchantID:   v.MerchantID,
				MerchantName: v.MerchantName,
				Omzet:        0,
				CreatedAt:    date,
			})
		}
	}

	return res, nil
}

func (s *userService) ListOutletOmzet(c *fiber.Ctx, date string, username string) (res []model.Outlet, err error) {
	l := log.GetLoggerContext(c.Context(), "service", "ListOutletOmzet")
	l.Infof("Start Execute ListOutletOmzet")

	outletRes, err := s.UserRepository.GetListOutlet(username)
	if err != nil {
		l.Errorf("Error when do GetListOutlet")
	}

	res, err = s.UserRepository.ListOutletOmzet(date, username)
	if err != nil {
		l.Errorf("Error when do ListOutletOmzet")
	}

	if len(res) <= 0 {
		for _, v := range outletRes {
			res = append(res, model.Outlet{
				MerchantID:   v.MerchantID,
				OutletID:     v.OutletID,
				MerchantName: v.MerchantName,
				OutletName:   v.OutletName,
				Omzet:        0,
				CreatedAt:    date,
			})
		}
	} else if len(res) < len(outletRes) && len(res) != 0 {
		for _, v := range outletRes {
			for j := 0; j < len(res); j++ {
				if v.MerchantID != res[j].MerchantID && v.OutletID != res[j].OutletID && j == len(res)-1 {
					res = append(res, model.Outlet{
						MerchantID:   v.MerchantID,
						OutletID:     v.OutletID,
						MerchantName: v.MerchantName,
						OutletName:   v.OutletName,
						Omzet:        0,
						CreatedAt:    date,
					})
				} else if v.MerchantID == res[j].MerchantID && v.OutletID != res[j].OutletID && j == len(res)-1 {
					res = append(res, model.Outlet{
						MerchantID:   v.MerchantID,
						OutletID:     v.OutletID,
						MerchantName: v.MerchantName,
						OutletName:   v.OutletName,
						Omzet:        0,
						CreatedAt:    date,
					})
				} else if v.MerchantID == res[j].MerchantID && v.OutletID == res[j].OutletID {
					break
				}
			}
		}
	}

	return res, nil
}

func (s *userService) LoginUsecase(c *fiber.Ctx, form model.Login) (token string, user model.Auth, err error) {
	l := log.GetLoggerContext(c.Context(), "service", "LoginUsecase")

	data := []byte(form.Password)
	b := md5.Sum(data)

	userHashPassword := hex.EncodeToString(b[:])

	//from DB
	userData, err := s.UserRepository.GetDataByUsername(form.Username)

	if userData.Password == userHashPassword {
		token, err = s.GenerateToken(c, form)
		if err != nil {
			l.Errorf("Invalid secret Key")
		}
		return token, userData, nil
	} else {
		l.Errorf("Invalid password")
	}

	return
}

func (s *userService) GenerateToken(c *fiber.Ctx, form model.Login) (token string, err error) {
	atClaims := jwt.MapClaims{}
	atExp := time.Now().Add(time.Minute * time.Duration(30*1000*60)).Unix()

	atClaims["user_name"] = form.Username
	atClaims["exp"] = atExp

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	token, err = at.SignedString([]byte("e00cf2*T&*#@YRdskjnds5USHDbsjasbd89*^#%*&!jbkj"))
	if err != nil {
		return
	}

	return token, nil
}
