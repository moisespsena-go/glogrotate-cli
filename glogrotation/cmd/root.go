// Copyright © 2019 Moises P. Sena <moisespsena@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"github.com/anmitsu/go-shlex"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/moisespsena-go/srvreader"

	"github.com/apex/log"
	"github.com/moisespsena-go/glogrotation-cli"

	"github.com/moisespsena-go/glogrotation"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var name = filepath.Base(os.Args[0])

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   name,
	Short: "Starts file writer rotation reads IN and writes to OUT",

	Long: longHelp,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			flags            = cmd.Flags()
			userExec         = viper.GetString("exec")
			userExecERD      = viper.GetBool("stderr-redirection-disabled")
			out              = viper.GetString("out")
			in               = viper.GetString("in")
			silent           = viper.GetBool("silent")
			printConfig, _   = flags.GetBool("print")
			maxUdpBufferSize = viper.GetInt("udp-max-bs")
		)

		for _, v := range args {
			if v != "" {
				out = v
			}
		}

		if out == "" {
			return errors.New("OUT not defined")
		}

		var cfg = logrotate.Config{
			HistoryDir:   strings.ReplaceAll(viper.GetString("history-dir"), "OUT", out),
			HistoryPath:  viper.GetString("history-path"),
			MaxSize:      viper.GetString("max-size"),
			HistoryCount: viper.GetInt("history-count"),
			Duration:     viper.GetString("duration"),
			FileMode:     os.FileMode(viper.GetInt("file-mode")),
			DirMode:      os.FileMode(viper.GetInt("dir-mode")),
		}

		var opt logrotate.Options
		if opt, err = cfg.Options(); err != nil {
			return
		}

		if printConfig {
			fmt.Fprintln(os.Stdout, "user-exec: "+userExec)
			fmt.Fprintln(os.Stdout, "stderr-redirection-disabled: ", userExecERD)
			fmt.Fprintln(os.Stdout, "out: "+out)
			fmt.Fprintln(os.Stdout, "in: "+in)
			fmt.Fprintln(os.Stdout, cfg.Yaml())
			fmt.Fprintln(os.Stdout, "udp-max-bs:", maxUdpBufferSize)
			return
		}

		Rotator := logrotate.New(out, opt)
		defer Rotator.Close()

		rw := glogrotation_cli.NewChanRW()
		var closers []io.Closer
		var ec int

		if userExec != "" {
			var args []string
			if args, err = shlex.Split(userExec, true); err != nil {
				return fmt.Errorf("user-exec flag: %s", err.Error())
			}
			func() {
				cmd := exec.Command(args[0], args[1:]...)
				cmd.Env = os.Environ()
				prepareCmd(cmd)
				inw, err := cmd.StdinPipe()
				if err != nil {
					log.Fatalf("pipe stdin failed: %s", err)
				}
				cmd.Stdout = rw

				if userExecERD {
					cmd.Stderr = os.Stderr
				} else {
					cmd.Stderr = rw
				}

				if err = cmd.Start(); err != nil {
					return
				}

				go io.Copy(inw, os.Stdin)

				sigs := make(chan os.Signal, 1)
				signal.Notify(sigs)
				go func() {
					defer func() {
						log.Info("child done")
						rw.Close()
					}()
					if err := cmd.Wait(); err != nil {
						switch et := err.(type) {
						case *exec.ExitError:
							if !strings.HasPrefix(et.Error(), "signal: ") {
								log.Errorf("exit code: %s", et)
							}
							ec = et.ExitCode()
						default:
							log.Error(err.Error())
						}
					}
				}()
				go func() {
					for sig := range sigs {
						if sig == syscall.SIGINT {
							sig = syscall.SIGTERM
						}
						if killable(sig) {
							cmd.Process.Signal(sig)
						}
					}
				}()
			}()
			if err != nil {
				return
			}
		}

		inm := map[string]interface{}{}

	inloop:
		// +udp:localhost:5678
		for _, in := range strings.Split(in, "+") {
			in = strings.TrimSpace(in)

			if _, ok := inm[in]; ok {
				continue
			}
			inm[in] = true

			switch {
			case in == "" || in == "-":
				if userExec == "" {
					go func() {
						if _, err := io.Copy(rw, os.Stdin); err != nil {
							if err != io.EOF {
								log.Errorf("STDIN: %s", err.Error())
							}
						}
					}()
				}
				continue inloop
			case srvreader.IsProto(in, "udp"):
				udps := srvreader.NewUDPServer(in, int16(maxUdpBufferSize), rw)
				closers = append(closers, udps)
				go func() {
					if err := udps.ListenAndServe(); err != nil {
						log.Errorf("UDP Serve[%s]: %s", in, err)
					}
				}()
			case srvreader.IsProto(in, "http"):
				tcps := srvreader.NewHTTPServerReader(in, rw)
				closers = append(closers, tcps)
				go func() {
					if err := tcps.ListenAndServe(); err != nil {
						log.Errorf("HTTP Serve[%s]: %s", in, err)
					}
				}()
			case srvreader.IsProto(in, "tcp"):
				tcps := srvreader.NewTCPServerReader(in, rw)
				closers = append(closers, tcps)
				go func() {
					if err := tcps.ListenAndServe(); err != nil {
						log.Errorf("TCP Serve[%s]: %s", in, err)
					}
				}()
			}
		}

		inm = nil

		defer func() {
			for _, c := range closers {
				c.Close()
			}
		}()

		if silent {
			_, err = io.Copy(Rotator, rw)
		} else {
			_, err = io.Copy(os.Stdout, io.TeeReader(rw, Rotator))
		}

		if ec != 0 {
			os.Exit(ec)
		}
		return
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	flags := rootCmd.Flags()
	flags.StringP("exec", "e", "", "execute PROGRAM and set your STDOUT and STDERR "+
		"(uses --stderr-redirection-disabled to disable it) as rotator input. Redirects the main STDIN to PROGRAM STDIN")
	flags.BoolP("stderr-redirection-disabled", "E", false, "on execute PROGRAM, disables "+
		"PROGRAM STDERR to STDOUT redirection.")
	flags.StringP("out", "o", "", "the OUTPUT file")
	flags.StringP("in", "i", "-", "the INPUT file. `-` (hyphen char) is STDIN. See INPUT "+
		"section for details")
	flags.StringP("history-dir", "c", "OUT.history", "history root directory")
	flags.StringP("history-path", "p", "%Y/%M", "dynamic direcotry path inside ROOT DIR "+
		"using TIME FORMAT")
	flags.StringP("duration", "d", "M", "rotates every DURATION. Accepted values: "+
		"Y - yearly, M - monthly, W - weekly, D - daily, h - hourly, m - minutely")
	flags.StringP("max-size", "S", "50M", "Forces rotation if current log size is "+
		"greather then MAX_SIZE. Values in bytes. Examples: 100, 100K, 50M, 1G, 1T")
	flags.Int16("udp-max-bs", 2048, "max UDP server buffer size. It's int16 value")
	flags.IntP("history-count", "C", 0, "Max history log count")
	flags.IntP("dir-mode", "M", 0750, "directory perms")
	flags.Lookup("dir-mode").DefValue = "0750"
	flags.IntP("file-mode", "m", 0640, "file perms")
	flags.Lookup("file-mode").DefValue = "0640"
	flags.Bool("print", false, "print current config")
	flags.Bool("silent", false, "disable tee to STDOUT")

	for _, v := range []string{
		"exec", "stderr-redirection-disabled",
		"out", "in",
		"history-dir", "history-path", "duration",
		"max-size", "history-count", "dir-mode",
		"file-mode", "silent", "udp-max-bs",
	} {
		viper.BindPFlag(v, flags.Lookup(v))
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(name)

	viper.AddConfigPath(".")
	viper.SetEnvPrefix(strings.Trim(strings.ToUpper(name), "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
