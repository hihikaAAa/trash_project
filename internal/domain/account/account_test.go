package account

import(
	//"errors"
	"testing"

	"github.com/google/uuid"
	//domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
)

func TestCreateAccountSuccess(t *testing.T){
	email, pHash, role := "hihika@gmail.com","1234", "WORKER"
	account, err := NewAccount(email, pHash, role)

	if err != nil{
		t.Fatalf("expected nil, got %v", err)
	}
	if account == nil{
		t.Fatal("expected account, got nil")
	}

	if account.Email != email{
		t.Fatalf("expected Email = %s, got %s", email, account.Email)
	}
	
	if account.PasswordHash != pHash{
		t.Fatalf("expected PasswordHash = %s, got %s", pHash, account.PasswordHash)
	}

	if account.ID == uuid.Nil{
		t.Fatal("expected ID, got nil")
	}

	if account.Role != WORKER{
		t.Fatalf("expected Role = %s, got %s", role, account.Role)
	}
}

func TestCreateAccountFail_BadEmail(t *testing.T) {
	tests := []string{
		"",
		"   ",
		"@gmail.com",
		"user@",
		"usergmail.com",
		"user@@gmail.com",
		"user@gmail",
		"user@gmail..com",
		"user..name@gmail.com",
		".user@gmail.com",
		"user.@gmail.com",
		"user@.gmail.com",
		"user@gmail.com.",
		"-user@gmail.com",
		"user@-gmail.com",  
		"user@gm_ail.com",
		"user@gm!ail.com",     
		"user name@gmail.com",
		"user\n@gmail.com",  
	}

	for _, email := range tests {
		account, err := NewAccount(email, "hash", "USER")
		if err == nil {
			t.Fatalf("expected error for email %q, got nil", email)
		}
		if account != nil {
			t.Fatalf("expected nil account for email %q", email)
		}
	}
}

func TestCreateAccountFail_BadRole(t *testing.T){
	account , err := NewAccount("hihika@gmail.com","hash","Person")

	if err == nil{
		t.Fatal("expected Error, got nil")
	}
	if account !=  nil{
		t.Fatal("expected account == nil")
	}
}



