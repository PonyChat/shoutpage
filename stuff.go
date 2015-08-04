package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/Xe/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string     `json:"user"`
	Password string     `json:"password"`
	Log      bool       `json:"log"`
	Networks []*Network `json:"networks"`
}

type Network struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	TLS      bool   `json:"tls"`
	Username string `json:"username"`
	Realname string `json:"realname"`
	Nick     string `json:"nick"`
	Join     string `json:"join"`
}

func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

func Crypt(password []byte) ([]byte, error) {
	defer clear(password)
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func MakeUser(name, password string) {
	ctext, err := Crypt([]byte(password))
	if err != nil {
		panic(err)
	}

	u := &User{
		Username: name,
		Password: string(ctext),
		Log:      true,
		Networks: []*Network{
			&Network{
				Name:     "PonyChat",
				Host:     "irc.ponychat.net",
				Port:     "6697",
				TLS:      true,
				Username: uuid.New()[:8],
				Realname: name + " via the PonyChat BNC",
				Nick:     name,
				Join:     "#bnc,#geek,#ponychat",
			},
		},
	}

	jsonOut, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		panic(err)
	}

	fname := os.Getenv("OUTPUT_PATH") + "/" + strings.ToLower(name) + ".json"
	log.Printf("Writing for %s to %s", name, fname)

	fout, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	_, err = fout.Write(jsonOut)
	if err != nil {
		panic(err)
	}
}

/*
func main() {
	ctext, err := Crypt([]byte("foobar"))

	if err != nil {
		log.Fatal(err)
	}

	u := &User{
		Username: "foo",
		Password: string(ctext),
		Log:      true,
		Networks: []*Network{
			&Network{
				Name:     "PonyChat",
				Host:     "irc.ponychat.net",
				Port:     "6697",
				TLS:      true,
				Username: "foo",
				Realname: "foo via the PonyChat BNC",
				Nick:     "foo",
				Join:     "#bnc,#geek,#ponychat",
			},
		},
	}

	fmt.Println(string(ctext))

	jsonOut, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonOut))
}
*/
