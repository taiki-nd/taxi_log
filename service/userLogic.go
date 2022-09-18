package service

import (
	"fmt"
	"log"

	"github.com/taiki-nd/taxi_log/model"
)

func UserValidation(user *model.User) (bool, []error) {
	var errs []error

	// uuid
	if len(user.Uuid) == 0 {
		log.Println("uuid null error")
		errs = append(errs, fmt.Errorf("uuid_null_error"))
	}
	//nickname
	if len(user.Nickname) == 0 {
		log.Println("nickname null error")
		errs = append(errs, fmt.Errorf("nickname_null_error"))
	}
	if len(user.Nickname) < 3 || 30 < len(user.Nickname) {
		log.Println("nickname letter count error")
		errs = append(errs, fmt.Errorf("nickname_letter_count_error"))
	}
	if len(user.Prefecture) == 0 {
		log.Println("prefecture null error")
		errs = append(errs, fmt.Errorf("prefecture_null_error"))
	}
	if len(user.Company) == 0 {
		log.Println("company null error")
		errs = append(errs, fmt.Errorf("company_null_error"))
	}
	if len(user.StyleFlg) == 0 {
		log.Println("style_flg null error")
		errs = append(errs, fmt.Errorf("style_flg_null_error"))
	}
	if !(user.StyleFlg == "every_other_day" || user.StyleFlg == "day" || user.StyleFlg == "night" || user.StyleFlg == "other") {
		log.Println("specified word error(open_flg)")
		errs = append(errs, fmt.Errorf("specified_word_error(open_flg)"))
	}
	if user.CloseDay == 0 {
		log.Println("close_day error")
		errs = append(errs, fmt.Errorf("close_day_error"))
	}
	if user.CloseDay < 1 || 31 < user.CloseDay {
		log.Println("close_day date error")
		errs = append(errs, fmt.Errorf("close_day_date_error"))
	}

	// errの出力
	if len(errs) != 0 {
		return false, errs
	}

	return true, nil
}
