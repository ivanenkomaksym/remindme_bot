package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/ivanenkomaksym/remindme_bot/config"
	"github.com/ivanenkomaksym/remindme_bot/domain/usecases"
	"github.com/ivanenkomaksym/remindme_bot/keyboards"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TimezoneController serves a small HTML page to capture browser timezone and a callback to persist it
type TimezoneController struct {
	userUseCase usecases.UserUseCase
	bot         *tgbotapi.BotAPI
	config      config.Config
}

func NewTimezoneController(userUseCase usecases.UserUseCase, bot *tgbotapi.BotAPI, config config.Config) *TimezoneController {
	return &TimezoneController{userUseCase: userUseCase, bot: bot, config: config}
}

// ServePage serves minimal HTML+JS that grabs browser timezone and redirects to callback
func (c *TimezoneController) ServePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// Minimal HTML template
	tmpl := `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Set Timezone</title>
  <style>body{font-family:system-ui,-apple-system,Segoe UI,Roboto,Ubuntu,Cantarell,Noto Sans,sans-serif;padding:24px;display:flex;align-items:center;justify-content:center;min-height:100vh;background:#0b132b;color:#e0e6f0} .card{background:#1c2541;border-radius:12px;box-shadow:0 6px 30px rgba(0,0,0,.3);padding:24px;max-width:520px;width:100%} h1{font-size:20px;margin:0 0 8px} p{opacity:.9;margin:0 0 16px} .muted{opacity:.7;font-size:14px} .btn{display:inline-block;background:#3a86ff;color:#fff;text-decoration:none;padding:12px 16px;border-radius:8px;font-weight:600} .btn:focus{outline:2px solid rgba(58,134,255,.5);outline-offset:2px}</style>
  <script>
    function go(){
      try{
        var tz = encodeURIComponent(Intl.DateTimeFormat().resolvedOptions().timeZone);
        var url = 'set-timezone/callback?user_id=%s&timezone=' + tz;
        window.location.replace(url);
      }catch(e){
        document.getElementById('status').textContent = 'Failed to detect timezone: ' + e;
      }
    }
    window.addEventListener('load', go);
  </script>
  <noscript>
    <style>.js-only{display:none}</style>
  </noscript>
</head>
<body>
  <div class="card">
    <h1>Detecting your timezone…</h1>
    <p id="status" class="muted js-only">Please wait, redirecting…</p>
    <p class="muted">If nothing happens, ensure JavaScript is enabled.</p>
    <a class="btn" href="#" onclick="go();return false;">Try again</a>
  </div>
</body>
</html>`

	page := fmt.Sprintf(tmpl, c.config.Bot.PublicURL, template.HTMLEscapeString(userID))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(page))
}

// Callback persists the timezone and then redirects back to the bot deep link
func (c *TimezoneController) Callback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	user, err := c.userUseCase.GetUser(userID)
	if err != nil {
		return
	}

	s := keyboards.T(user.Language)

	timezone := r.URL.Query().Get("timezone")
	if timezone == "" {
		http.Error(w, "timezone is required", http.StatusBadRequest)
		return
	}

	_ = c.userUseCase.UpdateLocation(userID, timezone)

	// Redirect to bot; optional start parameter to return user to chat
	botUser := c.bot.Self.UserName
	if botUser == "" {
		botUser = "" // fallback: open Telegram app
	}

	confirmationText := fmt.Sprintf("%s *%s*.", s.MsgTimezoneSet, timezone)
	msg := tgbotapi.NewMessage(userID, confirmationText)
	msg.ParseMode = tgbotapi.ModeMarkdown

	// Attach the command to remove the keyboard from the user's screen.
	// The 'true' argument makes it selective, so it only disappears for this one user.
	msg.ReplyMarkup = nil

	// Send the message
	if _, err := c.bot.Send(msg); err != nil {
		log.Printf("Failed to send confirmation/removal message to user %d: %v", userID, err)
	}

	redirectURL := fmt.Sprintf("https://t.me/%s", botUser)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
