# HorseyBot
HorseyBot is a service for manage media files via tracking them with keywords, users from group chat, set ping timer and maybe even more

## Table of Content

 - [Usage](#usage)
 - [Features](#features)
 - [Functionality](#functionality)
 - [Technologies](#technologies)

---

## Usage

### Application launch

Create an `.env` file in the root directory with the following variables:

```env
TG_TOKEN=your_telegram_bot_token
POSTGRES_USER=postgres_user_name
POSTGRES_PASSWORD=postgres_password
POSTGRES_DB=postgres_table_name
```

Then run ```docker-compose up -d``` command.

---

## Features

 - Connect the media with keywords for receive them with only one word (or choose from many of them)
 - Control the drop chance of media
 - Set ping timers for notify your friends from group chat
 - Manage all your connections and group chats in private messages

---

## Functionality 

### Commands

| Command          | Description                                                                      |
|------------------|----------------------------------------------------------------------------------|
| `/start`         | Register user and start using the bot                                            |
| `Помощь`         | Display all available commands and usage guide                                   |
| `Добавить связь` | Create a new connection for a group chat                                         |
| `Убрать связь`   | Disable a created connection                                                     |
| `Показать связи` | Shows all connections for a selected group chat. Filter them by users or actions |




### Database Schema

The service uses Postgres to store administrator, user, chat and subscriptions information

 - `administrators` - Stores userID and userName about administrators
 - `users` - Stores chatIDs to which group chats the user belongs along with the Telegram bot
 - `chats` - Keeps information about which chats bot a member of
 - `subscriptions` - Information of active connections

---

## Technologies

### Telegram Bot Service
 - **Language:** Go
 - **Telegram API:** telebot
 - **Database:** Postgres
 - **Query Builder:** Squirrel
 - **Logger:** slog
 - **Configuration:** yaml + .env files
 - **Container:** Docker
