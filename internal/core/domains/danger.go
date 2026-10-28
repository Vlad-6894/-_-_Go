package domains

import (
	"fmt"
	"regexp"
	core_errors "security/internal/core/errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Danger struct {
	ID          uuid.UUID
	Version     int
	Name        string
	MD5         string
	DangerType  string
	DangerLevel string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type DangerPut struct {
	Name        Put[string]
	MD5         Put[string]
	DangerType  Put[string]
	DangerLevel Put[string]
}

var (
	md5Regex    = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	RightLevels = []string{"LOW", "MEDIUM", "HIGH"}
)

func NewDangerPut(
	name Put[string],
	md5 Put[string],
	dangerType Put[string],
	dangerLevel Put[string],
) DangerPut {
	return DangerPut{
		Name:        name,
		MD5:         md5,
		DangerType:  dangerType,
		DangerLevel: dangerLevel,
	}
}

func NewDanger(
	id uuid.UUID,
	version int,
	name string,
	md5 string,
	dangerType string,
	dangerLevel string,
	createdAt time.Time,
	updatedAt *time.Time,
) Danger {
	return Danger{
		ID:          id,
		Version:     version,
		Name:        name,
		MD5:         md5,
		DangerType:  dangerType,
		DangerLevel: dangerLevel,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func NewUnitializedDanger(
	name string,
	md5 string,
	dangerType string,
	dangerLevel string,
) Danger {
	return Danger{
		ID:          uuid.New(),
		Version:     -1,
		Name:        name,
		MD5:         md5,
		DangerType:  dangerType,
		DangerLevel: dangerLevel,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}
}

func (d Danger) Validate() error {
	if len([]rune(d.Name)) < 1 || len([]rune(d.Name)) > 1000 {
		return fmt.Errorf("name must be less 1000 and more 1 char: %w", core_errors.ErrInvalidArgument)
	}

	cleaned := strings.TrimSpace(d.MD5)

	cleaned = strings.ToLower(cleaned)

	if !md5Regex.MatchString(cleaned) {
		return fmt.Errorf("md5 must be 32 chars and have A-F or 0-9: %w", core_errors.ErrInvalidArgument)
	}

	if len([]rune(d.DangerType)) < 1 || len([]rune(d.DangerType)) > 1000 {
		return fmt.Errorf("danger type must be less 1000 and more 1 char: %w", core_errors.ErrInvalidArgument)
	}

	if is_contains := IsContains(d.DangerLevel); !is_contains {
		return fmt.Errorf("danger level must be LOW, MEDIUM or HIGH: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func IsContains(level string) bool {
	for _, v := range RightLevels {
		if v == level {
			return true
		}
	}

	return false
}

func (u *Danger) ApplyPut(patch DangerPut) error {

	tmp := *u

	if patch.Name.Set {
		tmp.Name = *&patch.Name.Value
	}

	if patch.MD5.Set {
		tmp.MD5 = patch.MD5.Value
	}

	if patch.DangerType.Set {
		tmp.DangerType = patch.DangerType.Value
	}

	if patch.DangerLevel.Set {
		tmp.DangerLevel = patch.DangerLevel.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched danger error: %w", err)
	}

	*u = tmp

	return nil
}
