package extension

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/FriendsOfShopware/shopware-cli/version"
)

const (
	TypePlatformApp    = "app"
	TypePlatformPlugin = "plugin"
	TypeShopwareBundle = "shopware-bundle"

	ComposerTypePlugin = "shopware-platform-plugin"
	ComposerTypeApp    = "shopware-app"
	ComposerTypeBundle = "shopware-bundle"
)

func GetExtensionByFolder(path string) (Extension, error) {
	if _, err := os.Stat(fmt.Sprintf("%s/plugin.xml", path)); err == nil {
		return nil, fmt.Errorf("shopware 5 is not supported. Please use https://github.com/FriendsOfShopware/FroshPluginUploader instead")
	}

	if _, err := os.Stat(fmt.Sprintf("%s/manifest.xml", path)); err == nil {
		return newApp(path)
	}

	if _, err := os.Stat(fmt.Sprintf("%s/composer.json", path)); err != nil {
		return nil, fmt.Errorf("unknown extension type")
	}

	var ext Extension

	ext, err := newPlatformPlugin(path)
	if err != nil {
		ext, err = newShopwareBundle(path)
	}

	return ext, err
}

func GetExtensionByZip(filePath string) (Extension, error) {
	dir, err := os.MkdirTemp("", "extension")
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	file, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return nil, err
	}

	err = Unzip(file, dir)

	if err != nil {
		return nil, err
	}

	fileName := file.File[0].Name

	if strings.Contains(fileName, "..") {
		return nil, fmt.Errorf("invalid zip file")
	}

	extName := strings.Split(fileName, "/")[0]
	return GetExtensionByFolder(fmt.Sprintf("%s/%s", dir, extName))
}

type extensionTranslated struct {
	German  string `json:"german"`
	English string `json:"english"`
}

type extensionMetadata struct {
	Label       extensionTranslated
	Description extensionTranslated
}

type Extension interface {
	GetName() (string, error)
	GetResourcesDir() string

	// GetRootDir Returns the root folder where the code is located plugin -> src, app ->
	GetRootDir() string
	GetVersion() (*version.Version, error)
	GetLicense() (string, error)
	GetShopwareVersionConstraint() (*version.Constraints, error)
	GetType() string
	GetPath() string
	GetChangelog() (*extensionTranslated, error)
	GetMetaData() *extensionMetadata
	GetExtensionConfig() *Config
	Validate(context.Context, *ValidationContext)
}
