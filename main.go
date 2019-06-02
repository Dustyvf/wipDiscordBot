package main

import (
	"github.com/andersfylling/disgord"
	"strings"
)

// wip, command struct, used for command info, may revamp later to include more info on stuff
type cmd struct {
	execute   func(args []string, s disgord.Session, d *disgord.MessageCreate) error // function for command
	name      string                                                                 // name of command, currently irrelevent, only used to check if the command exists
	longdesc  string                                                                 // unimplemented long description to be used for help system
	shortdesc string                                                                 // unimplemented short description, to be used for tldr help
	man       []string                                                               // unimplemented multipage manual for help system, probably only to be used for complex commands
	parent    string                                                                 // unimplemented, to be used for defining the parent command for aliases
}

// main map of commands
var cmdmap = make(map[string]cmd)

// to be moved into its own file, but basically add commands into command map
func init() {
	cmdmap["ping"] = cmd{
		name: "currently this is irrelevent as long as something is in here",
		execute: func(args []string, session disgord.Session, data *disgord.MessageCreate) error {
			_, err := data.Message.Reply(session, "Pong!")
			return err
		},
	}
}

// thingy for parsing the command, trimming off the prefix (split into separate function?) and forking all the args into a string array
func commandfork(in, prefix string) (string, []string) {
	input := []rune(in)

	var msg []rune = input[len([]rune(prefix)):len(input)]

	x := strings.Index(string(msg), " ")

	if x <= 1 {
		var y []string
		return string(msg), y
	}
	cmds := string(msg[:x])
	arg := string(msg[(x + 1):])
	args := strings.Split(arg, " ")

	return cmds, args

}

// check if input contains prefix, supports multiple prefixes,
// returns whether or not it does have a prefix and due to supporting multiple prefixes in an array also returns the prefix
func prefixCheck(in string, prefix []string) (bool, string) {
	input := strings.ToLower(in)
	for i := 0; i < len(prefix); i++ {
		if strings.HasPrefix(input, prefix[i]) {
			return true, prefix[i]
		}
	}
	return false, ""
}

// to be implemented
func commandHandler() {

}

// hasPrefix is a disgord middleware to ensure that any incoming
// new messages contains at least one of the prefixes. Otherwise the
// handler will not be executed.
func hasPrefix(evt interface{}) interface{} {
	msgEvt := evt.(*disgord.MessageCreate)
	if a, _ := prefixCheck(msgEvt.Message.Content, gconf.Config.Prefix); a {
		return evt
	}

	// if nil is returned, this sends a signal to disgord
	// that the middleware denied the message content
	return nil
}

// main function
func main() {
	discord := disgord.New(&disgord.Config{
		BotToken: gconf.Auth.Token,
		Logger:   disgord.DefaultLogger(false),
	})
	defer func() {
		if err := discord.StayConnectedUntilInterrupted(); err != nil {
			panic(err)
		}
	}()

	discord.On(disgord.EvtMessageCreate, hasPrefix, func(session disgord.Session, data *disgord.MessageCreate) {
		_, pre := prefixCheck(data.Message.Content, gconf.Config.Prefix)
		command, args := commandfork(data.Message.Content, pre)
		session.Logger().Info(command, " | ", args, " | ", len(args))

		if cmdmap[command].name != "" {
			return // stop
		}

		if err := cmdmap[command].execute(args, session, data); err != nil {
			session.Logger().Error(err)
		}
	})
}
