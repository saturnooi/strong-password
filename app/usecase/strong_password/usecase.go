package strongpassword

import (
	"context"
	"database/sql"
	"encoding/json"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/moonrhythm/validator"
)

type StrongPasswordReq struct {
	InitPassword string `json:"init_password"`
}

type StrongPasswordRes struct {
	NumOfSteps int `json:"num_of_steps"`
}

func (p *StrongPasswordReq) Valid() error {
	v := validator.New()

	password := utf8.RuneCountInString(p.InitPassword)
	v.Must(password > 0 && password <= 40, "full name is required and must be less than 40 characters")
	return v.Error()
}

func StrongPasswordSteps(c *gin.Context, req StrongPasswordReq, db *sql.DB) (*StrongPasswordRes, error) {
	if err := req.Valid(); err != nil {
		return nil, err
	}

	step := calculateSteps(req.InitPassword)
	res := StrongPasswordRes{
		NumOfSteps: step,
	}
	log := StrongPasswordDto{
		Req: req,
		Res: res,
	}

	if err := insertLog(c.Request.Context(), log, db); err != nil {
		return nil, err
	}

	return &res, nil
}

func calculateSteps(password string) int {
	n := len(password)
	hasLower, hasUpper, hasDigit := false, false, false
	countRepeatedcharacters := 0
	change := 0

	for idx, ch := range password {
		if unicode.IsLower(ch) && hasLower == false {
			hasLower = true
		}
		if unicode.IsUpper(ch) && hasUpper == false {
			hasUpper = true
		}
		if unicode.IsDigit(ch) && hasDigit == false {
			hasDigit = true
		}

		if idx > 0 {
			if password[idx] == password[idx-1] {
				countRepeatedcharacters++
			} else {
				countRepeatedcharacters = 0
			}
			if countRepeatedcharacters >= 2 {
				change++
				countRepeatedcharacters = 0
			}
		}
	}

	missing := 0
	if !hasLower {
		missing++
	}
	if !hasUpper {
		missing++
	}
	if !hasDigit {
		missing++
	}

	if n < 6 {
		missing = max(missing, 6-n)
	}
	if n > 20 {
		missing = max(missing, n-20)
	}

	missing += change
	return missing
}

type StrongPasswordDto struct {
	Req StrongPasswordReq
	Res StrongPasswordRes
}

func insertLog(ctx context.Context, dto StrongPasswordDto, db *sql.DB) error {
	reqJSON, err := json.Marshal(dto.Req)
	if err != nil {
		return err
	}

	resJSON, err := json.Marshal(dto.Res)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		ctx,
		`INSERT INTO strong_password_log (req, res) VALUES ($1, $2)`,
		reqJSON, resJSON,
	)
	if err != nil {
		return err
	}
	return nil
}
