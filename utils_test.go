package main

import (
	"fmt"
	"testing"
)

func TestValidate(t *testing.T) {
	cmdreg := []string{
		"^git pull \\w+ \\w+$",
		"^git fetch \\w+ \\w+ && git reset --hard \\w+/\\w+$",
		"^git pull \\w+ \\w+ && supervisorctl restart \\w+$",
		"^git fetch \\w+ \\w+ && git reset --hard \\w+/\\w+ && supervisorctl restart \\w+$",
	}

	rightcmds := []string{
		"git pull origin master",
		"git pull origin master && supervisorctl restart webhooks",
		"git fetch origin master && git reset --hard origin/master && supervisorctl restart webhooks",
	}
	errcmds := []string{
		"ls -la",
		"cat /etc/passwd",
	}

	for _, cmd := range rightcmds {
		fmt.Println(cmd)
		if !Validate(cmdreg, cmd) {
			t.Error(fmt.Sprintf("validate error reg: %s,cmd: %s", cmdreg, cmd))
		}
	}

	for _, cmd := range errcmds {
		if Validate(cmdreg, cmd) {
			t.Error(fmt.Sprintf("validate error reg: %s,cmd: %s", cmdreg, cmd))
		}
	}

	pathregs := []string{
		"^/www/wwwroot/[^\\/]*$",
		"^/app/web/test",
	}
	rightpaths := []string{
		"/www/wwwroot/www",
		"/www/wwwroot/web.stestd",
	}

	errpaths := []string{
		"/www/wwwroot/../",
		"/home/www/wwwroot/web.stestd",
		"/www/wwwroot/ww../../../../",
	}
	for _, cmd := range rightpaths {
		fmt.Println(cmd)
		if !Validate(pathregs, cmd) {
			t.Error(fmt.Sprintf("validate error reg: %s,cmd: %s", pathregs, cmd))
		}
	}

	for _, cmd := range errpaths {
		fmt.Println(cmd)
		if Validate(pathregs, cmd) {
			t.Error(fmt.Sprintf("validate error reg: %s,cmd: %s", pathregs, cmd))
		}
	}

}
