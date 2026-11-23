package model

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSettings(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	viper.Reset()
	viper.SetConfigType("json")
	require.NoError(t, viper.ReadConfig(strings.NewReader(`{
		"entities": [
			{"name": "ctbc", "root_dir": "~/Downloads/bank-estatements/中國信託_CTBC", "filename": "*.pdf", "password_env": "UPASSPDF_TW_ID"},
			{"name": "esun", "root_dir": "/abs/dir", "filename": "*.pdf", "password_env": "UPASSPDF_TW_ID"}
		]
	}`)))

	s, err := LoadSettings()
	require.NoError(t, err)
	require.Len(t, s.Entities, 2)
	assert.Equal(t, filepath.Join(home, "Downloads/bank-estatements/中國信託_CTBC"), s.Entities[0].RootDir)
	assert.Equal(t, "/abs/dir", s.Entities[1].RootDir)
}

func TestLoadSettingsErrors(t *testing.T) {
	viper.Reset()
	viper.SetConfigType("json")
	require.NoError(t, viper.ReadConfig(strings.NewReader(`{"entities": []}`)))
	_, err := LoadSettings()
	assert.Error(t, err)

	viper.Reset()
	viper.SetConfigType("json")
	require.NoError(t, viper.ReadConfig(strings.NewReader(`{"entities": [{"name": "ctbc"}]}`)))
	_, err = LoadSettings()
	assert.Error(t, err)
}

func TestSelect(t *testing.T) {
	s := Settings{Entities: []Entity{{Name: "ctbc"}, {Name: "esun"}}}
	got, err := s.Select("ESUN")
	require.NoError(t, err)
	assert.Equal(t, []Entity{{Name: "esun"}}, got)
	_, err = s.Select("nope")
	assert.Error(t, err)
}

func TestEntityPasswordFromEnv(t *testing.T) {
	viper.Reset()
	viper.AutomaticEnv()
	t.Setenv("UPASSPDF_TEST_PASSWORD", "A123456789")
	pw, ok := Entity{PasswordEnv: "UPASSPDF_TEST_PASSWORD"}.Password()
	assert.True(t, ok)
	assert.Equal(t, "A123456789", pw)
	_, ok = Entity{PasswordEnv: "UPASSPDF_UNSET_PASSWORD"}.Password()
	assert.False(t, ok)
}
