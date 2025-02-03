package prettylogs

import (
	"os"

	"github.com/mrfoh/pretty-logs/internal/formatter"
	"github.com/spf13/cobra"
)

var (
	FormatTemplateFile string
	TimestampFormat    string
	LevelKey           string
	TimeKey            string
	PidKey             string
	NameKey            string
	ContextKey         string
	MsgKey             string
	ErrorObjectKeys    []string
	HostnameKey        string
	RequestKey         string
	ResponseKey        string
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prettylogs",
		Short: "Pretty print json log files",
		Long:  "Takes a json log from another process and pretty prints it to a specified output",
		Run: func(cmd *cobra.Command, args []string) {
			formatTemplateFile, _ := cmd.Flags().GetString("formatTemplateFile")
			outputFormat, _ := cmd.Flags().GetString("messageFormat")
			timestampFormat, _ := cmd.Flags().GetString("timestampFormat")
			levelKey, _ := cmd.Flags().GetString("levelKey")
			timeKey, _ := cmd.Flags().GetString("timeKey")
			hostnameKey, _ := cmd.Flags().GetString("hostnameKey")
			pidKey, _ := cmd.Flags().GetString("pidKey")
			nameKey, _ := cmd.Flags().GetString("nameKey")
			contextKey, _ := cmd.Flags().GetString("contextKey")
			msgKey, _ := cmd.Flags().GetString("msgKey")
			errorObjectKeys, _ := cmd.Flags().GetStringSlice("errorLikeObjectKeys")
			requestKey, _ := cmd.Flags().GetString("requestKey")

			// Create a new json formatter
			formatter, err := formatter.NewJsonFormatter(&formatter.FormatterOptions{
				FormatTemplateFile: formatTemplateFile,
				OutputFormat:       outputFormat,
				TimestampFormat:    timestampFormat,
				LevelKey:           levelKey,
				TimeKey:            timeKey,
				PidKey:             pidKey,
				NameKey:            nameKey,
				ContextKey:         contextKey,
				MsgKey:             msgKey,
				ErrorObjectKeys:    errorObjectKeys,
				HostnameKey:        hostnameKey,
				RequestKey:         requestKey,
			})
			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}

			// Process the log
			formatter.Process(os.Stdin)
		},
	}

	cmd.PersistentFlags().StringVarP(&FormatTemplateFile, "formatTemplateFile", "F", "", "File containing the format template for the logs")
	cmd.PersistentFlags().StringVarP(&TimestampFormat, "timestampFormat", "f", formatter.DEFAULT_TIMESTAMP_FORMAT, "Timestamp format for the logs")
	cmd.PersistentFlags().StringVarP(&LevelKey, "levelKey", "L", formatter.LEVEL_KEY, "Key for the log level in the log line")
	cmd.PersistentFlags().StringVarP(&TimeKey, "timeKey", "t", formatter.TIME_KEY, "Key for the log timestamp in the log line")
	cmd.PersistentFlags().StringVarP(&PidKey, "pidKey", "p", formatter.PID_KEY, "Key for the log pid of the application in the log line")
	cmd.PersistentFlags().StringVarP(&NameKey, "nameKey", "n", formatter.NAME_KEY, "Key for the log name of the application in the log line")
	cmd.PersistentFlags().StringVarP(&ContextKey, "contextKey", "c", formatter.CONTEXT_KEY, "Key for the log context in the log line")
	cmd.PersistentFlags().StringVarP(&MsgKey, "msgKey", "m", formatter.MESSAGE_KEY, "Key for the log message in the log line")
	cmd.PersistentFlags().StringSliceVarP(&ErrorObjectKeys, "errorLikeObjectKeys", "k", formatter.ERROR_LIKE_KEYS, "Keys that are considered to be error objects in the log line")
	cmd.PersistentFlags().StringVarP(&HostnameKey, "hostnameKey", "H", formatter.HOSTNAME_KEY, "Key for the log hostname in the log line")
	cmd.PersistentFlags().StringVarP(&RequestKey, "requestKey", "r", formatter.REQUEST_KEY, "Key for the log request in the log line")
	cmd.PersistentFlags().StringVarP(&ResponseKey, "responseKey", "R", formatter.RESPONSE_KEY, "Key for the log response in the log line")

	return cmd
}
