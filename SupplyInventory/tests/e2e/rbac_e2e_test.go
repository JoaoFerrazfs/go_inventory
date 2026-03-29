//go:build integration
package e2e

import (
       "bufio"
       "net/http"
       "os"
       "os/exec"
       "strings"
       "testing"
       "time"

       httpExpect "github.com/gavv/httpexpect/v2"
       "go_inventory/SupplyInventory/tests/e2e/testsetup"
)

// Para rodar: go test -tags=integration -v ./SupplyInventory/tests/e2e
// Certifique-se que o servidor está rodando em http://localhost:8000
// e o banco está limpo para o teste.

func TestRBAC_E2E(t *testing.T) {
       // Setup e seed do admin user via helper reutilizável
       cfg := testsetup.GetTestDBConfig()
       db, err := testsetup.ConnectWithRetry(cfg, 15, 2*time.Second)
       if err != nil {
	       t.Fatalf("Erro ao conectar no banco para seed: %v", err)
       }
       defer db.Close()
       adminEmail := "admin_e2e@test.com"
       adminPass := "SenhaForte123"
       if err := testsetup.SeedAdminUser(db, "Admin E2E", adminEmail, adminPass, "admin"); err != nil {
	       t.Fatalf("Erro ao inserir admin no banco: %v", err)
       }

       // 1. Sobe o servidor Go apontando para .env.test
       appCmd := exec.Command("go", "run", "main.go")
       appCmd.Dir = "/app"
       appCmd.Env = append(os.Environ(), "ENV_FILE=.env.test")
       appStdout, _ := appCmd.StdoutPipe()
       appStderr, _ := appCmd.StderrPipe()
       var stdoutLines, stderrLines []string
       if err := appCmd.Start(); err != nil {
	       t.Fatalf("Falha ao subir servidor Go: %v", err)
       }
       defer func() {
	       _ = appCmd.Process.Kill()
	       t.Logf("STDOUT final do backend:\n%s", strings.Join(stdoutLines, "\n"))
	       t.Logf("STDERR final do backend:\n%s", strings.Join(stderrLines, "\n"))
       }()

       // Aguarda o servidor subir (procura por "Listening" no stdout)
       ready := make(chan struct{})
       go func() {
	       scanner := bufio.NewScanner(appStdout)
	       for scanner.Scan() {
		       line := scanner.Text()
		       stdoutLines = append(stdoutLines, line)
		       if strings.Contains(line, "Starting server...") {
			       close(ready)
			       return
		       }
	       }
       }()
       go func() {
	       scanner := bufio.NewScanner(appStderr)
	       for scanner.Scan() {
		       line := scanner.Text()
		       stderrLines = append(stderrLines, line)
		       if strings.Contains(line, "Starting server...") {
			       close(ready)
			       return
		       }
	       }
       }()
       select {
       case <-ready:
       case <-time.After(90 * time.Second):
	       t.Logf("STDOUT do backend:\n%s", strings.Join(stdoutLines, "\n"))
	       t.Logf("STDERR do backend:\n%s", strings.Join(stderrLines, "\n"))
	       t.Fatal("Timeout ao subir servidor Go para e2e")
       }

       baseURL := "http://127.0.0.1:8082"
       // Aguarda o backend realmente aceitar conexões HTTP
       var httpReady bool
       var healthResp *http.Response
       var healthErr error
       for i := 0; i < 20; i++ {
	       healthResp, healthErr = http.Get(baseURL + "/healthz")
	       if healthErr == nil && healthResp.StatusCode < 500 {
		       httpReady = true
		       break
	       }
	       time.Sleep(500 * time.Millisecond)
       }
       if !httpReady {
	       t.Fatalf("Servidor Go não respondeu em %s após aguardar", baseURL)
       }
       e := httpExpect.New(t, baseURL)

       // 1. Criar usuário admin
       // Tenta criar usuário admin, ignora erro 422 (já existe)
       var adminResp, userResp *httpExpect.Response
       adminResp = e.POST("/api/v1/users/create").
	       WithJSON(map[string]interface{}{
		       "name": "Admin E2E",
		       "email": adminEmail,
		       "password": adminPass,
	       }).
	       Expect()
       if adminResp.Raw().StatusCode != http.StatusCreated && adminResp.Raw().StatusCode != http.StatusUnprocessableEntity {
	       t.Fatalf("Falha ao criar admin: status %d", adminResp.Raw().StatusCode)
       }

       // 2. Promover para admin diretamente via SQL
       _, err = db.Exec(`UPDATE user_entities SET role = 'admin' WHERE email = ?`, adminEmail)
       if err != nil {
	       t.Fatalf("Falha ao promover admin: %v", err)
       }
       time.Sleep(1 * time.Second) // Aguarda propagação

       // 3. Criar usuário comum
       userEmail := "user_e2e@test.com"
       userPass := "SenhaForte123"

       // Tenta criar usuário comum, ignora erro 422 (já existe)
       userResp = e.POST("/api/v1/users/create").
	       WithJSON(map[string]interface{}{
		       "name": "User E2E",
		       "email": userEmail,
		       "password": userPass,
	       }).
	       Expect()
       if userResp.Raw().StatusCode != http.StatusCreated && userResp.Raw().StatusCode != http.StatusUnprocessableEntity {
	       t.Fatalf("Falha ao criar user: status %d", userResp.Raw().StatusCode)
       }

       // 4. Login admin
       adminToken := e.POST("/api/v1/auth/login").
	       WithJSON(map[string]interface{}{
		       "email": adminEmail,
		       "password": adminPass,
	       }).
	       Expect().
	       Status(http.StatusOK).
	       JSON().Object().Value("token").String().Raw()

       // 5. Login user
       userToken := e.POST("/api/v1/auth/login").
	       WithJSON(map[string]interface{}{
		       "email": userEmail,
		       "password": userPass,
	       }).
	       Expect().
	       Status(http.StatusOK).
	       JSON().Object().Value("token").String().Raw()

       // 6. Admin acessa rota protegida
       e.GET("/api/v1/admin/racks").
	       WithHeader("Authorization", "Bearer "+adminToken).
	       WithQuery("page", "1").
	       Expect().
	       Status(http.StatusOK)

       // 7. User NÃO acessa rota protegida
       e.GET("/api/v1/admin/racks").
	       WithHeader("Authorization", "Bearer "+userToken).
	       WithQuery("page", "1").
	       Expect().
	       Status(http.StatusForbidden)
}
