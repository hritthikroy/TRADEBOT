# 📱 How to Get Telegram Bot Token and Chat ID

## Step 1: Create a Telegram Bot (Get Bot Token)

### 1.1 Open Telegram
- Open Telegram app on your phone or desktop
- Or use web version: https://web.telegram.org

### 1.2 Find BotFather
- In the search bar, type: `@BotFather`
- Click on the official BotFather (verified with blue checkmark)
- Click "START" button

### 1.3 Create Your Bot
1. Send this command: `/newbot`
2. BotFather will ask: "Alright, a new bot. How are we going to call it?"
3. Type your bot name (e.g., "My Trading Signals Bot")
4. BotFather will ask for a username
5. Type a username ending in "bot" (e.g., "mytradingsignals_bot")

### 1.4 Get Your Bot Token
- BotFather will reply with a message containing your token
- It looks like this: `123456789:ABCdefGHIjklMNOpqrsTUVwxyz-1234567890`
- **COPY THIS TOKEN** - you'll need it!

Example message from BotFather:
```
Done! Congratulations on your new bot. You will find it at t.me/mytradingsignals_bot. 
You can now add a description, about section and profile picture for your bot.

Use this token to access the HTTP API:
123456789:ABCdefGHIjklMNOpqrsTUVwxyz-1234567890

Keep your token secure and store it safely, it can be used by anyone to control your bot.
```

---

## Step 2: Get Your Chat ID

You have 3 options depending on where you want to receive signals:

### Option A: Personal Chat (Easiest)

#### A.1 Start Chat with Your Bot
1. Click on the link BotFather gave you (e.g., t.me/mytradingsignals_bot)
2. Click "START" button
3. Send any message to your bot (e.g., "Hello")

#### A.2 Get Your Chat ID
1. Open your browser
2. Go to this URL (replace `YOUR_BOT_TOKEN` with your actual token):
   ```
   https://api.telegram.org/botYOUR_BOT_TOKEN/getUpdates
   ```
   
   Example:
   ```
   https://api.telegram.org/bot123456789:ABCdefGHIjklMNOpqrsTUVwxyz/getUpdates
   ```

3. You'll see JSON response like this:
   ```json
   {
     "ok": true,
     "result": [
       {
         "update_id": 123456789,
         "message": {
           "message_id": 1,
           "from": {
             "id": 987654321,
             "is_bot": false,
             "first_name": "Your Name"
           },
           "chat": {
             "id": 987654321,
             "first_name": "Your Name",
             "type": "private"
           },
           "date": 1234567890,
           "text": "Hello"
         }
       }
     ]
   }
   ```

4. Look for `"chat":{"id":` - that number is your Chat ID
5. In the example above, the Chat ID is: `987654321`
6. **COPY THIS NUMBER**

---

### Option B: Telegram Channel (For Broadcasting)

#### B.1 Create a Channel
1. In Telegram, click the menu (☰)
2. Select "New Channel"
3. Give it a name (e.g., "Trading Signals")
4. Choose "Private" or "Public"
5. Click "Create"

#### B.2 Add Your Bot as Admin
1. Open your channel
2. Click on channel name at top
3. Click "Administrators"
4. Click "Add Administrator"
5. Search for your bot username (e.g., @mytradingsignals_bot)
6. Add it and give it "Post Messages" permission
7. Click "Done"

#### B.3 Get Channel Chat ID
1. Send a message to your channel
2. Go to this URL in browser:
   ```
   https://api.telegram.org/botYOUR_BOT_TOKEN/getUpdates
   ```

3. Look for `"chat":{"id":` in the response
4. For channels, the ID will be negative and start with `-100`
5. Example: `-1001234567890`
6. **COPY THIS NUMBER** (including the minus sign)

---

### Option C: Telegram Group

#### C.1 Create a Group
1. In Telegram, click menu (☰)
2. Select "New Group"
3. Add your bot to the group
4. Give the group a name

#### C.2 Make Bot Admin
1. Click on group name
2. Click "Edit"
3. Click "Administrators"
4. Add your bot as admin

#### C.3 Get Group Chat ID
1. Send a message in the group
2. Go to: `https://api.telegram.org/botYOUR_BOT_TOKEN/getUpdates`
3. Look for `"chat":{"id":`
4. Group IDs are negative (e.g., `-123456789`)
5. **COPY THIS NUMBER**

---

## Step 3: Add to .env File

1. Open `backend/.env` file
2. Replace the placeholder values:

```env
# Telegram Bot Configuration
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz-1234567890
TELEGRAM_CHAT_ID=987654321
```

### Examples:

**For Personal Chat:**
```env
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
TELEGRAM_CHAT_ID=987654321
```

**For Channel:**
```env
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
TELEGRAM_CHAT_ID=-1001234567890
```

**For Group:**
```env
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
TELEGRAM_CHAT_ID=-123456789
```

---

## Step 4: Test Your Configuration

### 4.1 Test Manually
Open this URL in your browser (replace with your values):
```
https://api.telegram.org/botYOUR_BOT_TOKEN/sendMessage?chat_id=YOUR_CHAT_ID&text=Test%20Message
```

Example:
```
https://api.telegram.org/bot123456789:ABCdefGHIjklMNOpqrsTUVwxyz/sendMessage?chat_id=987654321&text=Test%20Message
```

If successful, you'll see:
```json
{
  "ok": true,
  "result": {
    "message_id": 2,
    "from": {
      "id": 123456789,
      "is_bot": true,
      "first_name": "My Trading Signals Bot",
      "username": "mytradingsignals_bot"
    },
    "chat": {
      "id": 987654321,
      "first_name": "Your Name",
      "type": "private"
    },
    "date": 1234567890,
    "text": "Test Message"
  }
}
```

And you'll receive "Test Message" in Telegram!

### 4.2 Restart Backend
```bash
cd backend
go run .
```

You should see:
```
✅ Telegram bot initialized
```

---

## 🎉 You're Done!

Now you can:
1. Go to your trading bot web interface
2. Click on "Live Signals" tab
3. Click "📱 Start Telegram Bot (24/7)"
4. Signals will be sent to your Telegram!

---

## 🔍 Troubleshooting

### "Telegram bot not configured"
- Check that both `TELEGRAM_BOT_TOKEN` and `TELEGRAM_CHAT_ID` are in `.env`
- Make sure there are no extra spaces
- Restart the backend server

### "Failed to send message"
- Verify your bot token is correct
- For channels: Make sure bot is admin
- For groups: Make sure bot is member
- Test with the manual URL method above

### Can't find Chat ID
- Make sure you sent a message first
- Try refreshing the getUpdates URL
- Check you're using the correct bot token in the URL

### Bot not responding
- Make sure you clicked "START" on the bot
- For channels: Bot must be admin with "Post Messages" permission
- For groups: Bot must be added as member

---

## 📝 Quick Reference

**Get Bot Token:**
1. Message @BotFather
2. Send `/newbot`
3. Follow instructions
4. Copy token

**Get Chat ID:**
1. Send message to bot/channel/group
2. Visit: `https://api.telegram.org/botTOKEN/getUpdates`
3. Find `"chat":{"id":`
4. Copy the number

**Test:**
```
https://api.telegram.org/botTOKEN/sendMessage?chat_id=CHAT_ID&text=Test
```

---

## 🔐 Security Tips

1. **Never share your bot token** - it's like a password
2. **Don't commit .env to git** - it's already in .gitignore
3. **Use private channels** - don't broadcast to public
4. **Regenerate token if leaked** - message @BotFather with `/revoke`

---

Need help? Check the full setup guide in `TELEGRAM_BOT_SETUP.md`
