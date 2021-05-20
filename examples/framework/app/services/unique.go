package services

import "github.com/google/uuid"

type UniqueIDProvider struct {

}

func (p *UniqueIDProvider) Get() string {
	return uuid.New().String()
}
