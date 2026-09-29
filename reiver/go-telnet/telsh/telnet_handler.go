package telsh

import (
	"honey/reiver/go-oi"
	"honey/reiver/go-telnet"

	"bytes"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	defaultExitCommandName = "exit"
	defaultPrompt          = "\r\n> "
	defaultWelcomeMessage  = "\r\nWelcome!\r\n"
	defaultExitMessage     = "\r\nGoodbye!\r\n"
	defaultLogin           = ""
	defaultPassword        = ""
	defaultLoginMessage    = "Login: "
	defaultPasswordMessage = "Password: "
	defaultErrorMessage    = "Login incorrect. Try again.\r\n"
)

type ShellHandler struct {
	muxtex       sync.RWMutex
	producers    map[string]Producer
	elseProducer Producer

	ExitCommandName string
	Prompt          string
	WelcomeMessage  string
	Login           string
	Password        string
	LoginMessage    string
	PasswordMessage string
	ExitMessage     string
	ErrorMessage    string
}

func NewShellHandler() *ShellHandler {
	producers := map[string]Producer{}

	telnetHandler := ShellHandler{
		producers: producers,

		Prompt:          defaultPrompt,
		ExitCommandName: defaultExitCommandName,
		WelcomeMessage:  defaultWelcomeMessage,
		Login:           defaultLogin,
		Password:        defaultPassword,
		LoginMessage:    defaultLoginMessage,
		PasswordMessage: defaultPasswordMessage,
		ExitMessage:     defaultExitMessage,
		ErrorMessage:    defaultErrorMessage,
	}

	return &telnetHandler
}

func (telnetHandler *ShellHandler) Register(name string, producer Producer) error {

	telnetHandler.muxtex.Lock()
	telnetHandler.producers[name] = producer
	telnetHandler.muxtex.Unlock()

	return nil
}

func (telnetHandler *ShellHandler) MustRegister(name string, producer Producer) *ShellHandler {
	if err := telnetHandler.Register(name, producer); nil != err {
		panic(err)
	}

	return telnetHandler
}

func (telnetHandler *ShellHandler) RegisterHandlerFunc(name string, handlerFunc HandlerFunc) error {

	produce := func(ctx telnet.Context, name string, args ...string) Handler {
		return PromoteHandlerFunc(handlerFunc, args...)
	}

	producer := ProducerFunc(produce)

	return telnetHandler.Register(name, producer)
}

func (telnetHandler *ShellHandler) MustRegisterHandlerFunc(name string, handlerFunc HandlerFunc) *ShellHandler {
	if err := telnetHandler.RegisterHandlerFunc(name, handlerFunc); nil != err {
		panic(err)
	}

	return telnetHandler
}

func (telnetHandler *ShellHandler) RegisterElse(producer Producer) error {

	telnetHandler.muxtex.Lock()
	telnetHandler.elseProducer = producer
	telnetHandler.muxtex.Unlock()

	return nil
}

func (telnetHandler *ShellHandler) MustRegisterElse(producer Producer) *ShellHandler {
	if err := telnetHandler.RegisterElse(producer); nil != err {
		panic(err)
	}

	return telnetHandler
}

func (telnetHandler *ShellHandler) ServeTELNET(c net.Conn, ctx telnet.Context, writer telnet.Writer, reader telnet.Reader) {

	logger := ctx.Logger()
	if nil == logger {
		logger = internalDiscardLogger{}
	}

	colonSpaceCommandNotFoundEL := []byte(": command not found\r\n")

	var prompt bytes.Buffer
	var authprompt bytes.Buffer
	var exitCommandName string
	var welcomeMessage string
	var Login string
	var Password string
	var LoginMessage string
	var PasswordMessage string
	var exitMessage string
	var ErrorMessage string

	var tmpLogin string
	var tmpPassword string

	prompt.WriteString(telnetHandler.Prompt)
	promptBytes := prompt.Bytes()

	authprompt.WriteString("")
	authpromptBytes := authprompt.Bytes()

	exitCommandName = telnetHandler.ExitCommandName
	welcomeMessage = telnetHandler.WelcomeMessage
	Login = telnetHandler.Login
	Password = telnetHandler.Password
	LoginMessage = telnetHandler.LoginMessage
	PasswordMessage = telnetHandler.PasswordMessage
	exitMessage = telnetHandler.ExitMessage
	ErrorMessage = telnetHandler.ErrorMessage

	if _, err := oi.LongWriteString(writer, welcomeMessage); nil != err {
		logger.Errorf("Problem long writing welcome message: %v", err)
		return
	}
	logger.Debugf("Wrote welcome message: %q.", welcomeMessage)

	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	p := buffer[:]
	var line bytes.Buffer
	if (Login != "") && (Password != "") {
		if _, err := oi.LongWriteString(writer, LoginMessage); nil != err {
			logger.Errorf("Problem long writing welcome message: %v", err)
			return
		}
		logger.Debugf("Wrote welcome message: %q.", PasswordMessage)

		for {
			// Read 1 byte.
			n, err := reader.Read(p)
			if n <= 0 && nil == err {
				continue
			} else if n <= 0 && nil != err {
				break
			}

			line.WriteByte(p[0])
			if '\n' == p[0] {
				lineString := line.String()

				if "\r\n" == lineString {
					line.Reset()
					if _, err := oi.LongWrite(writer, authpromptBytes); nil != err {
						return
					}
					continue
				}
				fields := strings.Fields(lineString)
				if len(fields) <= 0 {
					line.Reset()
					if _, err := oi.LongWrite(writer, authpromptBytes); nil != err {
						return
					}
					continue
				}

				tmpLogin = fields[0]
				break
			}
		}
		line.Reset()
		if _, err := oi.LongWriteString(writer, PasswordMessage); nil != err {
			logger.Errorf("Problem long writing welcome message: %v", err)
			return
		}

		for {
			n, err := reader.Read(p)
			if n <= 0 && nil == err {
				continue
			} else if n <= 0 && nil != err {
				break
			}

			line.WriteByte(p[0])
			if '\n' == p[0] {
				lineString := line.String()

				if "\r\n" == lineString {
					line.Reset()
					if _, err := oi.LongWrite(writer, authpromptBytes); nil != err {
						return
					}
					continue
				}
				fields := strings.Fields(lineString)
				if len(fields) <= 0 {
					line.Reset()
					if _, err := oi.LongWrite(writer, authpromptBytes); nil != err {
						return
					}
					continue
				}

				tmpPassword = fields[0]
				break
			}
		}

		var producer Producer

		telnetHandler.muxtex.RLock()
		var ok bool
		producer, ok = telnetHandler.producers["auth"]
		telnetHandler.muxtex.RUnlock()

		if !ok {
			telnetHandler.muxtex.RLock()
			producer = telnetHandler.elseProducer
			telnetHandler.muxtex.RUnlock()
		}

		if nil != producer {
			ip := strings.Split(c.RemoteAddr().String(), ":")[0]
			const layout = "02/01/2006_15:04:05"
			td := time.Now().Format(layout)
			fields := strings.Fields(td + " " + ip + " " + tmpLogin + " " + tmpPassword)
			handler := producer.Produce(ctx, "auth", fields...)
			if nil != handler {
				if stdoutPipe, err := handler.StdoutPipe(); nil != err {
					connect(ctx, writer, stdoutPipe, nil)
				}
			}

		}
		if (tmpPassword != Password) || (tmpLogin != Login) {
			oi.LongWriteString(writer, ErrorMessage)
			return
		}
	}
	if _, err := oi.LongWrite(writer, promptBytes); nil != err {
		logger.Errorf("Problem long writing prompt: %v", err)
		return
	}
	logger.Debugf("Wrote prompt: %q.", promptBytes)

	line.Reset()
	for {
		// Read 1 byte.
		n, err := reader.Read(p)
		if n <= 0 && nil == err {
			continue
		} else if n <= 0 && nil != err {
			break
		}

		line.WriteByte(p[0])
		//logger.Tracef("Received: %q (%d).", p[0], p[0])

		if '\n' == p[0] {
			lineString := line.String()

			if "\r\n" == lineString {
				line.Reset()
				if _, err := oi.LongWrite(writer, promptBytes); nil != err {
					return
				}
				continue
			}

			//@TODO: support piping.
			fields := strings.Fields(lineString)
			logger.Debugf("Have %d tokens.", len(fields))
			logger.Tracef("Tokens: %v", fields)
			if len(fields) <= 0 {
				line.Reset()
				if _, err := oi.LongWrite(writer, promptBytes); nil != err {
					return
				}
				continue
			}

			field0 := fields[0]

			if exitCommandName == field0 {
				oi.LongWriteString(writer, exitMessage)
				return
			}

			var producer Producer

			telnetHandler.muxtex.RLock()
			var ok bool
			producer, ok = telnetHandler.producers[field0]
			telnetHandler.muxtex.RUnlock()

			if !ok {
				telnetHandler.muxtex.RLock()
				producer = telnetHandler.elseProducer
				telnetHandler.muxtex.RUnlock()
			}

			if nil == producer {
				oi.LongWrite(writer, []byte(field0))
				oi.LongWrite(writer, colonSpaceCommandNotFoundEL)
				line.Reset()
				if _, err := oi.LongWrite(writer, promptBytes); nil != err {
					return
				}
				continue
			}

			handler := producer.Produce(ctx, field0, fields[1:]...)
			if nil == handler {
				oi.LongWrite(writer, []byte(field0))
				oi.LongWrite(writer, colonSpaceCommandNotFoundEL)
				line.Reset()
				oi.LongWrite(writer, promptBytes)
				continue
			}
			var wg sync.WaitGroup

			if stdoutPipe, err := handler.StdoutPipe(); nil != err {
				//@TODO:
			} else if nil == stdoutPipe {
				//@TODO:
			} else {
				wg.Add(1)
				connect(ctx, writer, stdoutPipe, &wg)
			}

			if stderrPipe, err := handler.StderrPipe(); nil != err {
				//@TODO:
			} else if nil == stderrPipe {
				//@TODO:
			} else {
				wg.Add(1)
				connect(ctx, writer, stderrPipe, &wg)
			}

			if err := handler.Run(); nil != err {
				//@TODO:
			}
			// Wait for the output pumps to flush every byte to the socket
			// before emitting the prompt, so the prompt never races ahead of
			// the command output (e.g. the trailing newline of "help").
			wg.Wait()
			line.Reset()
			if _, err := oi.LongWrite(writer, promptBytes); nil != err {
				return
			}
		}

		//@TODO: Are there any special errors we should be dealing with separately?
		if nil != err {
			break
		}
	}

	oi.LongWriteString(writer, exitMessage)
	return
}

// connect pumps reader into writer in a background goroutine. When wg is
// non-nil it is marked done once the pump has flushed every byte, so the
// caller can wait for the output to be fully written before emitting the
// next prompt (otherwise the prompt can race ahead of the trailing bytes).
func connect(ctx telnet.Context, writer io.Writer, reader io.Reader, wg *sync.WaitGroup) {

	logger := ctx.Logger()

	go func(logger telnet.Logger) {
		if nil != wg {
			defer wg.Done()
		}

		var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
		p := buffer[:]

		for {
			// Read 1 byte.
			n, err := reader.Read(p)
			if n <= 0 && nil == err {
				continue
			} else if n <= 0 && nil != err {
				break
			}

			//logger.Tracef("Sending: %q.", p)
			//@TODO: Should we be checking for errors?
			oi.LongWrite(writer, p)
			//logger.Tracef("Sent: %q.", p)
		}
	}(logger)
}
