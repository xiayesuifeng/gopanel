package caddy

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"time"
)

func ValidateConfig(config []byte) error {
	path, err := exec.LookPath("caddy")
	if err != nil {
		return err
	}

	tmpFile := fmt.Sprintf("%s/gopanel.%d.json", os.TempDir(), time.Now().Unix())
	if err := ioutil.WriteFile(tmpFile, config, 0600); err != nil {
		return err
	}

	cmd := exec.Command(path, "validate", "-config", tmpFile)
	out, err := cmd.CombinedOutput()

	os.Remove(tmpFile)

	if err != nil {
		str := string(out)
		exp, err := regexp.Compile("validate: decoding config: (.*)")
		if err == nil && exp.Match(out) {
			ssm := exp.FindStringSubmatch(str)
			if len(ssm) > 0 {
				return errors.New(ssm[len(ssm)-1])
			}
		}

		return errors.New(str)
	}

	return nil
}
