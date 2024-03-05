package kubernetes

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ApplyManifest(manifestContent string) error {
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = bytes.NewBufferString(manifestContent)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("failed to apply manifest: %v, output: %s", err, output)
	}

	fmt.Printf("Manifest applied successfully: %s\n", output)
	return nil
}
