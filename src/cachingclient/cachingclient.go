package cachingclient

import (
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

func GetSecretCached(secretid string) string {
	secretCache, _ := secretcache.New()
	result, _ := secretCache.GetSecretString(secretid)
	// Use secret to connect to secured resource.
	return result
}




