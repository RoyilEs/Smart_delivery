package config

type SiteInfo struct {
	CreatedAt string `yaml:"created_at" json:"created_at"`
	Title     string `yaml:"title" json:"title"`
	Email     string `yaml:"email" json:"email"`
	Name      string `yaml:"name" json:"name"`
}
