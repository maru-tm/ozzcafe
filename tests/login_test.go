package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func TestLogin(t *testing.T) {
	// Устанавливаем адрес драйвера Selenium
	remoteWebDriver := "http://localhost:4444/wd/hub"

	// Создаем экземпляр WebDriver
	wd, err := selenium.NewRemote(
		selenium.Capabilities{
			"browserName": "chrome",
		},
		remoteWebDriver,
	)
	if err != nil {
		t.Fatalf("ошибка создания веб-драйвера: %v", err)
	}
	defer wd.Quit()

	// Открываем страницу логина
	err = wd.Get("http://localhost:8080/login")
	if err != nil {
		t.Fatalf("не удалось открыть страницу логина: %v", err)
	}

	// Находим элементы формы
	emailField, err := wd.FindElement(selenium.ByCSSSelector, "#email")
	if err != nil {
		t.Fatalf("не удалось найти поле email: %v", err)
	}

	passwordField, err := wd.FindElement(selenium.ByCSSSelector, "#password")
	if err != nil {
		t.Fatalf("не удалось найти поле пароля: %v", err)
	}

	loginButton, err := wd.FindElement(selenium.ByCSSSelector, "button[type='submit']")
	if err != nil {
		t.Fatalf("не удалось найти кнопку логина: %v", err)
	}

	// Вводим данные для входа
	emailField.SendKeys("aruzhanduyssenova@gmail.com")
	passwordField.SendKeys("asd")

	// Нажимаем на кнопку логина
	loginButton.Click()

	// Ожидаем редирект на страницу с токеном
	var currentURL string
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		var err error
		currentURL, err = wd.CurrentURL()
		if err != nil {
			return false, err
		}

		// Проверяем, что URL содержит правильный путь
		if !strings.HasPrefix(currentURL, "http://localhost:8080/admin?token=") {
			return false, nil
		}
		return true, nil
	}, 10*time.Second)

	// Проверка редиректа
	if err != nil {
		t.Fatalf("редирект не удался. Текущий URL: %v", currentURL)
	}

	t.Logf("Успешно редиректировано на: %s", currentURL)
}
