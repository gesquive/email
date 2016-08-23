package cmd

import (
	"bufio"
	"fmt"
	"net/mail"
	"os"
	"os/user"

	cli "github.com/gesquive/cli-log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

var cfgFile string
var displayVersion string

var logDebug bool
var showVersion bool

// RootCmd is our only command
var RootCmd = &cobra.Command{
	Use:   "email [flags] <message>",
	Short: "Send an email from the command line",
	Long: `Send an email from the command line.

If a flag is tagged with 'multi', multiple versions of the flag are accepted`,
	ValidArgs: []string{"MESSAGE"},
	Run:       run,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	displayVersion = version
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.config/email.yml)")
	RootCmd.PersistentFlags().BoolVarP(&logDebug, "debug", "D", false,
		"Write debug messages to console")
	RootCmd.PersistentFlags().MarkHidden("debug")
	RootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false,
		"Show the version and exit")

	RootCmd.PersistentFlags().BoolP("strict-parsing", "e", false,
		"Fail to send the email when any email address is malformed")
	RootCmd.PersistentFlags().StringSliceP("to", "t", nil,
		"Destination addresses (multi)")
	RootCmd.MarkFlagRequired("to")
	RootCmd.PersistentFlags().StringP("from", "f", "",
		"From address on email (default $USER@$HOST)")
	RootCmd.PersistentFlags().StringP("reply-to", "r", "",
		"Reply to address")
	RootCmd.PersistentFlags().StringSliceP("cc", "c", nil,
		"Carbon copy adresses (multi)")
	RootCmd.PersistentFlags().StringSliceP("bcc", "b", nil,
		"Blind carbon copy addresses (multi)")
	RootCmd.PersistentFlags().StringP("subject", "s", "",
		"Subject of email")
	RootCmd.PersistentFlags().StringP("message", "m", "",
		"Plain text content of email")
	RootCmd.PersistentFlags().StringP("html-message", "H", "",
		"Alternate HTML content of email")
	RootCmd.PersistentFlags().StringSliceP("attachment", "a", nil,
		"File to attach to email (multi)")
	//TODO: remove this workaround
	//	when https://github.com/spf13/viper/issues/200 is fixed
	RootCmd.PersistentFlags().Lookup("to").DefValue = ""
	RootCmd.PersistentFlags().Lookup("cc").DefValue = ""
	RootCmd.PersistentFlags().Lookup("bcc").DefValue = ""
	RootCmd.PersistentFlags().Lookup("attachment").DefValue = ""

	RootCmd.PersistentFlags().StringP("smtp-server", "x", "localhost",
		"The SMTP server to send email through")
	RootCmd.PersistentFlags().Uint32P("smtp-port", "o", 25,
		"The port to use for the SMTP server")
	RootCmd.PersistentFlags().StringP("smtp-username", "u", "",
		"Authenticate the SMTP server with this user")
	RootCmd.PersistentFlags().StringP("smtp-password", "p", "",
		"Authenticate the SMTP server with this password")

	viper.SetEnvPrefix("email")
	viper.AutomaticEnv()
	viper.BindEnv("strict_parsing")
	viper.BindEnv("to")
	viper.BindEnv("from")
	viper.BindEnv("reply_to")
	viper.BindEnv("cc")
	viper.BindEnv("bcc")
	viper.BindEnv("subject")
	viper.BindEnv("message")
	viper.BindEnv("html_message")
	viper.BindEnv("attachment")
	viper.BindEnv("smtp_server")
	viper.BindEnv("smtp_port")
	viper.BindEnv("smtp_username")
	viper.BindEnv("smtp_password")

	viper.BindPFlag("strict-parsing", RootCmd.PersistentFlags().Lookup("strict-parsing"))
	viper.BindPFlag("email.to", RootCmd.PersistentFlags().Lookup("to"))
	viper.BindPFlag("email.from", RootCmd.PersistentFlags().Lookup("from"))
	viper.BindPFlag("email.reply-to", RootCmd.PersistentFlags().Lookup("reply-to"))
	viper.BindPFlag("email.cc", RootCmd.PersistentFlags().Lookup("cc"))
	viper.BindPFlag("email.bcc", RootCmd.PersistentFlags().Lookup("bcc"))
	viper.BindPFlag("email.subject", RootCmd.PersistentFlags().Lookup("subject"))
	viper.BindPFlag("email.message", RootCmd.PersistentFlags().Lookup("message"))
	viper.BindPFlag("email.html", RootCmd.PersistentFlags().Lookup("html-message"))
	viper.BindPFlag("email.attachments", RootCmd.PersistentFlags().Lookup("attachment"))
	viper.BindPFlag("smtp.server", RootCmd.PersistentFlags().Lookup("smtp-server"))
	viper.BindPFlag("smtp.port", RootCmd.PersistentFlags().Lookup("smtp-port"))
	viper.BindPFlag("smtp.username", RootCmd.PersistentFlags().Lookup("smtp-username"))
	viper.BindPFlag("smtp.password", RootCmd.PersistentFlags().Lookup("smtp-password"))

	viper.SetDefault("strict-parsing", false)
	viper.SetDefault("smtp.server", "localhost")
	viper.SetDefault("smtp.port", 25)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/email.yml") // adding home directory as first search path
	viper.AddConfigPath("/etc/email")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if !showVersion {
			fmt.Println("Error opening config: ", err)
		}
	}
}

func run(cmd *cobra.Command, args []string) {
	if logDebug {
		cli.SetLogLevel(cli.LevelDebug)
	}
	if showVersion {
		cli.Info(displayVersion)
		os.Exit(0)
	}
	cli.Debug("Running with debug turned on")

	message := viper.GetString("email.message")
	pipeInput := getPipedInput()
	if len(pipeInput) > 0 {
		message = pipeInput
	}
	sendMessage(message)
}

func getPipedInput() string {
	// Detect if user is piping in text
	fileInput, err := os.Stdin.Stat()
	if err != nil {
		cli.Error(err.Error())
		os.Exit(2)
	}

	var text string
	pipeFound := fileInput.Mode()&os.ModeNamedPipe != 0
	if pipeFound {
		cli.Debug("Pipe input detected, reading")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text += scanner.Text()
		}

		if err != nil {
			cli.Error("Failed to read piped data")
			cli.Error(err.Error())
		}
	}

	return text
}

func sendMessage(message string) {
	strict := viper.GetBool("strict-parsing")
	msg := gomail.NewMessage()
	fromAddress := viper.GetString("email.from")
	if len(fromAddress) == 0 {
		fromAddress = getDefaultEmailAddress()
	}
	cli.Debug("From: %s", fromAddress)
	msg.SetHeader("From", fromAddress)
	replyAddress, err := formatEmail(viper.GetString("email.reply-to"))
	if strict && err != nil {
		cli.Warn("%v", err)
		cli.Error("Will not send email")
		return
	} else if len(replyAddress) > 0 {
		cli.Debug("Setting reply-to: %s", replyAddress)
		msg.SetHeader("Reply-To", replyAddress)
	}
	toAddresses, err := formatEmailList(viper.GetStringSlice("email.to"), strict)
	if strict && err != nil {
		cli.Warn("%v", err)
		cli.Error("Will not send email")
		return
	} else if len(toAddresses) > 0 {
		cli.Debug("Adding cc: %v", toAddresses)
		msg.SetHeader("To", toAddresses...)
	}
	ccAddresses, err := formatEmailList(viper.GetStringSlice("email.cc"), strict)
	if strict && err != nil {
		cli.Warn("%v", err)
		cli.Error("Will not send email")
		return
	} else if len(ccAddresses) > 0 {
		cli.Debug("Adding cc: %v", ccAddresses)
		msg.SetHeader("Cc", ccAddresses...)
	}
	bccAddresses, err := formatEmailList(viper.GetStringSlice("email.bcc"), strict)
	if strict && err != nil {
		cli.Warn("%v", err)
		cli.Error("Will not send email")
		return
	} else if len(bccAddresses) > 0 {
		cli.Debug("Adding bcc: %v", bccAddresses)
		msg.SetHeader("Bcc", bccAddresses...)
	}
	msg.SetHeader("Subject", viper.GetString("email.subject"))
	msg.SetBody("text/plain", message)
	htmlMessage := viper.GetString("email.html")
	if len(htmlMessage) > 0 {
		if len(message) == 0 {
			msg.SetBody("text/html", htmlMessage)
		} else {
			msg.AddAlternative("text/html", htmlMessage)
		}
	}
	attachments := viper.GetStringSlice("email.attachments")
	if len(attachments) > 0 {
		for _, a := range attachments {
			msg.Attach(a)
		}
	}

	//TODO: Add an option for tls/ssl connections
	smtpHost := viper.GetString("smtp.server")
	smtpPort := viper.GetInt("smtp.port")
	username := viper.GetString("smtp.username")
	password := viper.GetString("smtp.password")
	var dialer *gomail.Dialer
	if len(username) > 0 || len(password) > 0 {
		cli.Debug("Connecting too %s:%s@%s:%d", username, password, smtpHost, smtpPort)
		dialer = gomail.NewDialer(smtpHost, smtpPort, username, password)
	} else {
		cli.Debug("Connecting too %s:%d", smtpHost, smtpPort)
		dialer = &gomail.Dialer{Host: smtpHost, Port: smtpPort}
	}

	if err := dialer.DialAndSend(msg); err != nil {
		cli.Error("An error occurred when sending email")
		cli.Fatalln(err)
	}
	msg.WriteTo(os.Stdout)
}

func formatEmailList(list []string, strictParsing bool) ([]string, error) {
	var formattedList []string
	for _, r := range list {
		formattedAddress, err := formatEmail(r)
		if err != nil {
			if strictParsing {
				return []string{},
					fmt.Errorf("Could not parse address '%s': %v", r, err)
			}
			cli.Warn("Could not parse address '%s'", r)
			cli.Warn("%v", err)
		}
		formattedList = append(formattedList, formattedAddress)
	}
	return formattedList, nil
}

func formatEmail(address string) (string, error) {
	email, err := mail.ParseAddress(address)
	if err != nil {
		return "", err
	}

	fAddress := ""
	if len(email.Name) > 0 {
		fAddress = fmt.Sprintf("\"%s\" <%s>", email.Name, email.Address)
	} else {
		fAddress = email.Address
	}
	return fAddress, nil
}

func getDefaultEmailAddress() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	userinfo, err := user.Current()
	username := "user"
	if err == nil {
		username = userinfo.Username
	}
	return fmt.Sprintf("%s@%s", username, hostname)
}
