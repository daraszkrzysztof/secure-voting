package main
import (
	"testing"
	"github.com/joho/godotenv"
	"log"
	"github.com/stretchr/testify/assert"
	"os"
)


func TestIfCanCreateNewSecureVotingService(t *testing.T) {

	//arrange
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	expectedName := os.Getenv("ADMIN_LOGIN")
	expectedPasswd := os.Getenv("ADMIN_SECRET")

	//act
	svs := NewSecureVotingService();

	//assert
	if svs == nil {
		t.Fatalf("Couldn't create new secure voting service")
	}

	assert.NotNil(t, expectedName)
	assert.NotNil(t, expectedPasswd)
	assert.Equal(t, svs.sv.Admins[expectedName].Id, expectedName)
	assert.NotEqual(t, svs.sv.Admins[expectedName].PasswdHash, expectedPasswd)
	assert.NotEmpty(t, svs.sv.Admins[expectedName].PasswdHash)
}