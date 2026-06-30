package service

import (
	"encoding/json"

	"menu-recommend/internal/repository"
)

type AboutFeature struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	Bg          string `json:"bg"`
}

type AboutInfo struct {
	Name        string         `json:"name"`
	Subtitle    string         `json:"subtitle"`
	Description string         `json:"description"`
	Slogan      string         `json:"slogan"`
	Version     string         `json:"version"`
	Email       string         `json:"email"`
	Wechat      string         `json:"wechat"`
	TermsURL    string         `json:"terms_url"`
	PrivacyURL  string         `json:"privacy_url"`
	Features    []AboutFeature `json:"features"`
}

type AppConfigService struct {
	repo *repository.AppConfigRepo
}

func NewAppConfigService(repo *repository.AppConfigRepo) *AppConfigService {
	return &AppConfigService{repo: repo}
}

func (s *AppConfigService) GetAboutInfo() (*AboutInfo, error) {
	keys := []string{
		"about.name",
		"about.subtitle",
		"about.description",
		"about.slogan",
		"about.version",
		"about.email",
		"about.wechat",
		"about.terms_url",
		"about.privacy_url",
		"about.features",
	}
	configs, err := s.repo.FindActiveByKeys(keys)
	if err != nil {
		return nil, err
	}

	values := make(map[string]string, len(configs))
	for _, config := range configs {
		values[config.ConfigKey] = config.ConfigValue
	}

	info := &AboutInfo{
		Name:        values["about.name"],
		Subtitle:    values["about.subtitle"],
		Description: values["about.description"],
		Slogan:      values["about.slogan"],
		Version:     values["about.version"],
		Email:       values["about.email"],
		Wechat:      values["about.wechat"],
		TermsURL:    values["about.terms_url"],
		PrivacyURL:  values["about.privacy_url"],
		Features:    []AboutFeature{},
	}

	if raw := values["about.features"]; raw != "" {
		_ = json.Unmarshal([]byte(raw), &info.Features)
	}

	return info, nil
}
