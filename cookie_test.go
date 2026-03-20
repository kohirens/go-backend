package backend

import (
	"net/http"
	"os"
	"testing"

	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/www/storage"
)

func TestDecryptCookie(t *testing.T) {
	fixedDir := tmpDir + "/secrets"
	_ = os.Mkdir(fixedDir, os.ModePerm)
	_, e0 := fsio.CopyToDir(fixtureDir+"/secrets/test-01.json", fixedDir, string(os.PathSeparator))
	if e0 != nil {
		t.Fatal(e0)
		return
	}
	fixedStore, e1 := storage.NewLocalStorage(tmpDir)
	if e1 != nil {
		t.Fatalf("Error creating new local storage: %v", e1)
		return
	}
	fixtureApp := NewWithDefaults("test-01", fixedStore)
	fixtureApp.LoadGPG()

	cases := []struct {
		name    string
		cName   string
		cValue  string
		r       *http.Request
		a       App
		w       http.ResponseWriter
		want    string
		wantErr bool
	}{
		{
			"success",
			"test-01",
			"1234",
			&http.Request{
				Header: http.Header{
					//"Cookie": []string{
					//	//"test-01=LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3VjREVGIvQk4yUnVaZmtTQVFkQTRPc0pXL0kwUlU2L1RKQ0l3a1g2TnFIaTRtZDVYZWVwaVlDYzh2dVQ5d1V3CmlHbXVBWTQ3SWYvSWpSVFdZNVRuTWFPdjBHN3ZFR3N3T09zb2sxcXpMRnRrRVV2Q2xrSk9PL0pVNjRlK3p1SVQKMG9VQkszZS9ZOXNSUUZPdklaeHNCN0UyQ2trcnY3aDd0K1JFSXUwV2FKMDlCajAyWE9pZ2ZyNHp5VmQrMmg1dwoyb1E4NS9HL0ZlSFZWY0phRyttbm44MXQ1RTMrUEtYVGJ0UTFDdWQxcGNkVkMwTHZzVngvUnYxYm9iYXBleWwzCkZJbTNjelo5WlV2RzF4YXhZU0d4L0xWVzhtV1JHTjVoakgzL01TUk0zNXZTUG5PWlU2bTEKPXdFdzAKLS0tLS1FTkQgUEdQIE1FU1NBR0UtLS0tLQ==; Path=/; Secure",
					//},
				},
			},
			fixtureApp,
			&MockResponse{
				Headers: make(http.Header),
			},
			"1234",
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if e := EncryptCookie(c.cName, c.cValue, "/", c.w, c.a); (e != nil) != c.wantErr {
				t.Errorf("EncryptCookie error: %v", e)
				return
			}

			// grab the very cookie just encrypted and test decryption.
			c.r.Header.Set("Cookie", c.w.Header().Get("Set-Cookie"))

			got, err := DecryptCookie(c.cName, c.r, c.a)
			if (err != nil) != c.wantErr {
				t.Errorf("DecryptCookie() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			if string(got.Value) != c.want {
				t.Errorf("DecryptCookie() got = %v, want %v", got.Value, c.want)
				return
			}
		})
	}
}
