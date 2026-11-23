package model

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntityDir(t *testing.T) {
	s := Settings{RootDir: "/root/statements"}
	assert.Equal(t, "/root/statements/中國信託_CTBC", s.EntityDir(Entity{Path: "中國信託_CTBC"}))
	assert.Equal(t, "/abs/dir", s.EntityDir(Entity{Path: "/abs/dir"}))
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
