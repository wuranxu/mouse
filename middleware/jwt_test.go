package middleware

import "testing"

func TestJWT_ParseToken(t *testing.T) {
	token := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiQ3JlYXRlZEF0IjoiMjAyMy0wMi0yOFQwMDozMToxNCswODowMCIsIlVwZGF0ZWRBdCI6IjIwMjMtMDItMjhUMDA6MzE6MTcrMDg6MDAiLCJEZWxldGVkQXQiOm51bGwsIm5hbWUiOiJXT09EWSIsInVzZXJuYW1lIjoid29vZHkiLCJlbWFpbCI6IjYxOTQzNDE3NkBxcS5jb20iLCJwYXNzd29yZCI6Ind1cmFueHUiLCJsYXN0TG9naW5BdCI6IjIwMjMtMDItMjhUMDA6MzE6MzQrMDg6MDAiLCJyb2xlIjoxLCJ0b2tlbiI6IiJ9.iw5pRKK4roKUUz6_zNcyjXR3WYxrkvdziXR0fKNhnLo`
	parseToken, err := JWTUtil.ParseToken(token)
	if err != nil {
		t.Error(err)
	}
	t.Log(parseToken)
}
