package repository

type ShortLinkRepository interface {
	SaveMappedUrls(mappedUrls map[string]string) error
	GetOrigUrlByShortCode(shortCode string) (string, error)
	GetMappedUrls() (map[string]string, error)
}
