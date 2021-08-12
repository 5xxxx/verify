package verify

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IsValidMobilePhoneNumber(mobilePhoneNumber string) bool {
	reg := `^(((13[0-9]{1})|(14[0-9]{1})|(15[0-9]{1})|(16[0-9]{1})|(17[0-9]{1})|(18[0-9]{1})|(19[0-9]{1}))+\d{8})$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(mobilePhoneNumber)
}

func IsPassword(password string) bool {
	reg := `[\S]{6,20}`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(password)
}

func IsAccountName(accountName string) bool {
	reg := `[\S]{6,20}`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(accountName)
}

func IsPayPassword(password string) bool {
	reg := `[0-9]{6}`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(password)
}

func IsVerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func IsValidUsername(username string) bool {
	match, _ := regexp.MatchString("/^[a-z]{1}([a-z0-9]){4,11}$", username)
	if !match {
		return false
	}
	return true
}

func IsObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}

var nameReg = regexp.MustCompile(`^\p{Han}[\p{Han}·]{1,15}$`)

func IsRealName(name string) bool {
	return nameReg.Match([]byte(name))
}

var weight = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
var valid_value = [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
var valid_province = map[string]string{
	"11": "北京市",
	"12": "天津市",
	"13": "河北省",
	"14": "山西省",
	"15": "内蒙古自治区",
	"21": "辽宁省",
	"22": "吉林省",
	"23": "黑龙江省",
	"31": "上海市",
	"32": "江苏省",
	"33": "浙江省",
	"34": "安徽省",
	"35": "福建省",
	"36": "山西省",
	"37": "山东省",
	"41": "河南省",
	"42": "湖北省",
	"43": "湖南省",
	"44": "广东省",
	"45": "广西壮族自治区",
	"46": "海南省",
	"50": "重庆市",
	"51": "四川省",
	"52": "贵州省",
	"53": "云南省",
	"54": "西藏自治区",
	"61": "陕西省",
	"62": "甘肃省",
	"63": "青海省",
	"64": "宁夏回族自治区",
	"65": "新疆维吾尔自治区",
	"71": "台湾省",
	"81": "香港特别行政区",
	"91": "澳门特别行政区",
}

func isValidCitizenNo18(citizenNo18 *[]byte) bool {
	nLen := len(*citizenNo18)
	if nLen != 18 {
		return false
	}
	nSum := 0
	for i := 0; i < nLen-1; i++ {
		n, _ := strconv.Atoi(string((*citizenNo18)[i]))
		nSum += n * weight[i]
	}
	mod := nSum % 11
	if valid_value[mod] == (*citizenNo18)[17] {
		return true
	}
	return false
}

func isLeapYear(nYear int) bool {
	if nYear <= 0 {
		return false
	}
	if (nYear%4 == 0 && nYear%100 != 0) || nYear%400 == 0 {
		return true
	}
	return false
}

func checkBirthdayValid(nYear, nMonth, nDay int) bool {
	if nYear < 1900 || nMonth <= 0 || nMonth > 12 || nDay <= 0 || nDay > 31 {
		return false
	}

	curYear, curMonth, curDay := time.Now().Date()
	if nYear == curYear {
		if nMonth > int(curMonth) {
			return false
		} else if nMonth == int(curMonth) && nDay > curDay {
			return false
		}
	}

	if 2 == nMonth {
		if isLeapYear(nYear) && nDay > 29 {
			return false
		} else if nDay > 28 {
			return false
		}
	} else if 4 == nMonth || 6 == nMonth || 9 == nMonth || 11 == nMonth {
		if nDay > 30 {
			return false
		}
	}

	return true
}

func checkProvinceValid(citizenNo []byte) bool {
	provinceCode := make([]byte, 0)
	provinceCode = append(provinceCode, citizenNo[:2]...)
	provinceStr := string(provinceCode)

	for i := range valid_province {
		if provinceStr == i {
			return true
		}
	}

	return false
}

func IsValidCitizenNo(idNo string) bool {
	if len(idNo) == 0 {
		return false
	}
	var citizenNo *[]byte
	_idNo := []byte(idNo)
	citizenNo = &_idNo
	if !isValidCitizenNo18(citizenNo) {
		return false
	}
	for i, v := range *citizenNo {
		n, _ := strconv.Atoi(string(v))
		if n >= 0 && n <= 9 {
			continue
		}
		if v == 'X' && i == 16 {
			continue
		}
		return false
	}
	if !checkProvinceValid(*citizenNo) {
		return false
	}
	nYear, _ := strconv.Atoi(string((*citizenNo)[6:10]))
	nMonth, _ := strconv.Atoi(string((*citizenNo)[10:12]))
	nDay, _ := strconv.Atoi(string((*citizenNo)[12:14]))
	if !checkBirthdayValid(nYear, nMonth, nDay) {
		return false
	}
	return true
}

func CitizenNoInfo(idNo string) (birthday time.Time, sex string, address string, err error) {
	if !IsValidCitizenNo(idNo) {
		err = errors.New("不合法的身份证号码。")
		return
	}
	citizenNo := []byte(idNo)
	birthday, err = time.Parse("2006-01-02", string(citizenNo[6:10])+"-"+string(citizenNo[10:12])+"-"+string(citizenNo[12:14]))
	if err != nil {
		return time.Time{}, "", "", err
	}
	genderMask, _ := strconv.Atoi(string(citizenNo[16]))
	if genderMask%2 == 0 {
		sex = "女"
	} else {
		sex = "男"
	}
	address = valid_province[string(citizenNo[:2])]
	return
}
