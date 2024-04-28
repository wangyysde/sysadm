package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"sysadm/utils"
)

type User struct {
	FirstName string  `json:"firstName" db:"firstName"`
	LastName  string  `json:"lastName" db:"lastName"`
	Scores    []Score `json:"scores" db:"scores"`
	Score     Score   `json:"score" db:"score"`
}

type Score struct {
	Subject string
	Score   int
	Desc    *Desc
}

type NewScore struct {
	Item  string
	Score int
}

type Desc struct {
	Item1 string
	Item2 string
}

var users = []User{
	{FirstName: "larry", LastName: "wang", Score: Score{}},
	{FirstName: "jenny", LastName: "huang", Score: Score{}},
}

var user1 = User{FirstName: "aaaaaaa", LastName: "AAAAAAAA"}
var user2 = User{FirstName: "bbbbbbb", LastName: "BBBBBBBB"}
var userSlice = []*User{&user1, &user2}

func (u User) GetLastName(firstName string) string {
	return firstName
}

func main() {
	score2 := Score{Subject: "english", Score: 60, Desc: &Desc{Item1: "aaaaaa", Item2: "AAAAAAA"}}
	user2.Score = score2

	fName, e := GetFeildValueByName(userSlice[0], "Scores")
	if e != nil {
		fmt.Printf("get first name error: %s\n", e)
	} else {
		fmt.Printf("got first name is: %s\n", fName.(string))
	}

	s, e := GetFeildValueByName(userSlice[1], "Score")
	if e != nil {
		fmt.Printf("get scroe error: %s\n", e)
	} else {
		sScore := s.(Score)
		fmt.Printf("got score is: %+v\n", sScore)
	}

	e = setFieldValue(userSlice[0], "CCCCCCCCCCCCC", "FirstName")
	if e != nil {
		fmt.Printf("struct error: %s\n", e)
	} else {
		fmt.Printf("new struct data: %+v\n", userSlice[0])
	}

	// score1 := Score{Subject: "chinese", Score: 80}
	score2 = Score{Subject: "english", Score: 60, Desc: &Desc{Item1: "aaaaaa", Item2: "AAAAAAA"}}
	//	score3 := NewScore{Item: "huaxue", Score: 60}
	e = setFieldValue(userSlice[0], score2, "Score")
	if e != nil {
		fmt.Printf("struct error: %s\n", e)
	} else {
		fmt.Printf("new struct data: %+v\n", userSlice[0])
	}
	e = setFieldValue(userSlice[1], &score2, "Score")
	if e != nil {
		fmt.Printf("pointer error: %s\n", e)
	} else {
		fmt.Printf("new pointer data: %+v\n", userSlice[1])
	}
	fmt.Printf("user: %+v\n", userSlice)
	fmt.Printf("score: %+v\n", userSlice[1].Score.Desc)

	os.Exit(1)

}

func Unmarshal(data map[string]interface{}, dst any) error {
	dT := reflect.TypeOf(dst)
	if dT.Kind() != reflect.Pointer || dT.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("the type of dst is not a pointer or the destination where dst point to is not a struct")
	}

	dTElem := dT.Elem()
	dV := reflect.ValueOf(dst).Elem()

	for i := 0; i < dTElem.NumField(); i++ {
		field := dTElem.Field(i)
		if !field.IsExported() {
			continue
		}

		tag, okTag := field.Tag.Lookup("db")
		if !okTag || tag == "" {
			continue
		}

		v, _ := data[tag]
		switch fieldType := field.Type.Kind(); fieldType {
		case reflect.Bool:
			value := utils.Interface2Bool(v)
			dV.Field(i).SetBool(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var value = int64(0)
			vStr := utils.Interface2String(v)
			if vStr != "" {
				tmpValue, e := strconv.ParseInt(vStr, 10, 64)
				if e == nil {
					value = tmpValue
				} else {
					return fmt.Errorf("can not umarshal feild %s for %s", tag, e)
				}
			}

			dV.Field(i).SetInt(value)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var value uint64 = 0
			valueStr := utils.Interface2String(v)
			if valueStr != "" {
				tmpValue, e := strconv.ParseUint(valueStr, 10, 64)
				if e == nil {
					value = tmpValue
				} else {
					return fmt.Errorf("can not umarshal feild %s for %s", tag, e)
				}
			}
			dV.Field(i).SetUint(value)
		case reflect.Float32, reflect.Float64:
			var value float64 = 0
			vStr := utils.Interface2String(v)
			if vStr != "" {
				tmpValue, e := strconv.ParseFloat(vStr, 10)
				if e == nil {
					value = tmpValue
				} else {
					return fmt.Errorf("can not umarshal feild %s for %s", tag, e)
				}
			}
			dV.Field(i).SetFloat(value)
		case reflect.String:
			value := utils.Interface2String(v)
			dV.Field(i).SetString(value)
		default:
			continue
		}
	}

	return nil
}
