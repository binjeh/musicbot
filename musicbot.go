package main

import (
  
  "log"
  "time"
  "flag"
  "github.com/bwmarrin/dgvoice"
  "github.com/bwmarrin/discordgo"
  "github.com/spf13/viper"
  "github.com/fsnotify/fsnotify"
)

var bot map[string]string

// loadConfig load the config file
func loadConfig (filename string) (error){
  // Read the config.toml file
  viper.SetConfigType("toml")
  viper.SetConfigFile(filename)
  log.Println("Opening", filename)
  //viper.AddConfigPath(".")
  err := viper.ReadInConfig()
  if err != nil {
    log.Println("Config file (bot.toml) not found.\nExit ...\n")
    return err
  }
  bot = viper.GetStringMapString("bot")
  log.Println("URL:", bot["url"])
  return nil
}

// connectionOn open the url stream connection
func connectionOn(filename string) {
  // load the config file
  err := loadConfig(filename)
  if err != nil {
    log.Println(err)
    return
  }

  discord, err := discordgo.New("Bot " + bot["token"])
  if err != nil {
    log.Println("error creating Discord session,", err)
    return
  }
  defer discord.Close()
  
  log.Println("Bot is Opening.")

  //discord.AddHandler(messageCreate)

  // Open Websocket
  err = discord.Open()
  if err != nil {
    log.Println("Error Open():", err)
    return
  }

  log.Println("Bot user test.")

  _, err = discord.User("@me")
  if err != nil {
    // Login unsuccessful
    log.Println(err)
    return
  } // Login successful

  log.Println("Bot is now running. Press CTRL-C to exit.")

  // Set Status
  discord.UpdateStatus(0, bot["status"])
  
  // Join to voice channel
  vs, err := discord.ChannelVoiceJoin(bot["guild_id"], bot["channel_id"], false, false)
  if err != nil {
    log.Println(err)
    return
  }
  defer vs.Close()

  // Play to URL
  dgvoice.PlayAudioFile(vs, bot["url"])
}

func main() {
  file_name := flag.String("f", "bot.toml", "Set path for the config file.")
  flag.Parse()
  
  // Hot reload
  viper.WatchConfig()
  viper.OnConfigChange(func (e fsnotify.Event) {
    log.Println("The config file changed:", e.Name)
    dgvoice.KillPlayer()
  })
  
  for {
    connectionOn(*file_name)
    // Restart the bot if the file is changed or the connection fail
    log.Println("The bot is restarting ...")
    time.Sleep(5000 * time.Millisecond)
  }

  return
}


