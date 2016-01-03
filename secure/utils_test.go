package secure
import (
	"testing"
)

func TestIfCanLoadPage(t *testing.T){

	//arrange
	testCases := []struct{
		in string
		want string
	}{
		{"abc", "3a81oZNherrMQXNJriBBMRLm-k6JqX6iCp7u5ktV05ohkpkqJ0_BqDa6PCOj_uu9RU1EI2Q86A4qmslPpUyknw=="},
		{"xyz", "Sj7YFH43h2rcj3Yyjlq8wbRw5qz8GO_qATX5g2BJU6WOGDwaYIbpG6PoIdkm9f3rN3YcfKAyipY_XpKHBnW3KA=="},
	}

	for _,c := range testCases {
		//act
		got:= DoUserPasswdHash(c.in)

		//assert
		if got != c.want {
			t.Fatalf("Got %s, want %s!", got, c.want)
		}
	}
}

func TestIfCanGenerateKeyPair(t *testing.T){

	//arrange
	//act
	privateKey, err := GenerateRSAKeyPair()

	//assert
	if err != nil {
		t.Fail()
	}

	if privateKey==nil {
		t.Fail()
	}
}

func TestIfCanMarshallKeyPair(t *testing.T){

	//arrange
	privateKey, err := GenerateRSAKeyPair()

	//act
	priv_der, pub_pem, err := MarshalToX509KeyPair(privateKey)

	//assert
	if err != nil {
		t.Fatalf(err.Error())
	}

	if priv_der==nil || pub_pem==nil {
		t.Fail()
	}
}