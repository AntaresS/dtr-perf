package stress

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type Job struct {
	ClientName string
	Config     *Config
	Err        error
}

func StressPull(ctx context.Context, j *Job) error {

	image := getImageName(j)
	// parse the duration from Config file and time must be positive
	if strings.Contains(j.Config.Pull.Duration, "-") {
		return fmt.Errorf("duration must be positive")
	}
	logrus.Infof("run pulling test for %s", j.Config.Pull.Duration)
	duration, err := time.ParseDuration(j.Config.Pull.Duration)
	if err != nil {
		logrus.Errorf("time conversion failed: %s", err.Error())
		return err
	}
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	// continuously run the pulling test for the given duration
	timeout := time.After(duration)
	for {
		select {
		case <-ctx.Done():
			logrus.Errorf("context has been canceled %s", ctx.Err())
			return ctx.Err()
		case <-timeout:
			logrus.Infof("pulling test completed after %f minutes", duration.Minutes())
			return nil
		default:
			// hiccups caused by network or load can make test result inaccurate
			// thus add retries instead of returning failure immediately
			isSuccess := false
			for i := 0; i < 3; i++ {
				logrus.Infof("excuting pulling test at %s for the %v try", time.Now(), i)
				err = PullImage(ctx, j, cli, image)
				if err != nil {
					logrus.Errorf("pulling error: %s", err.Error())
					//return err
				} else {
					logrus.Info("finished pulling test with SUCCESS at %s", time.Now())
					time.Now().String()
					isSuccess = true
					break
				}
			}
			_, removeErr := cli.ImageRemove(ctx, image, types.ImageRemoveOptions{Force: true, PruneChildren: true})
			if removeErr != nil {
				log.Printf("remove image %s failed: %s ", image, removeErr)
			} else {
				log.Printf("remove image %s succeeded", image)
			}
			if !isSuccess {
				return err
			}
		}
	}
}

func PullImage(ctx context.Context, j *Job, cli client.APIClient, image string) error {

	return PullImageWithDockerClient(ctx, cli, j, image)
}

func PullImageWithDockerClient(ctx context.Context, cli client.APIClient, j *Job, image string) error {
	logrus.WithField("image", image).Info("pulling image")
	/*	options := types.ImagePullOptions{
		RegistryAuth: sharedutils.MakeRegistryAuth(j.Config.Username, j.Config.Password, j.Config.RefreshToken),
	}*/

	//progress, err := cli.ImagePull(ctx, image, options)
	progress, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer progress.Close()
	// pulling logic lifted from docker/testkit/provisioner/pull.go to avoid extra vendoring
	scanner := bufio.NewScanner(progress)
	for scanner.Scan() {
		line := scanner.Text()
		var msg JSONMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			// Consider this a hard fail so we don't get stuck if things get out of whack
			return fmt.Errorf("malformed progress line during pull of %s: %s - %s", image, line, err)
		}

		if msg.Error != nil {
			return fmt.Errorf("failed to load %s: %s", image, msg.Error.Message)
		}
		if msg.Progress != nil && msg.Progress.Total > 0 {
			fmt.Printf("\r%s %s layer %s %0.2f%%", image, msg.ID, msg.Status, float64(msg.Progress.Current)/float64(msg.Progress.Total)*100)
			if msg.Progress.Current == msg.Progress.Total {
				fmt.Println()
			}
		} else {
			logrus.Debugf("%s %s %s", image, msg.ID, msg.Status)
		}
	}
	return nil
}

func getImageName(j *Job) string {
	var image string
	if j.Config.DTRURL != "" {
		image = fmt.Sprintf("%s/%s/%s:%s", j.Config.DTRURL, j.Config.Pull.Namespace, j.Config.Pull.RepoName, j.Config.Pull.TagName)
	} else {
		if j.Config.Pull.Namespace != "" {
			image = fmt.Sprintf("%s/%s:%s", j.Config.Pull.Namespace, j.Config.Pull.RepoName, j.Config.Pull.TagName)
		} else {
			image = fmt.Sprintf("%s:%s", j.Config.Pull.RepoName, j.Config.Pull.TagName)
		}
	}
	return image
}

// Types lifted from docker/docker/pkg/jsonmessage to avoid TTY dependencies
// JSONError represents a JSON Error
type JSONError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// JSONProgress represents a JSON-encoded progress instance
type JSONProgress struct {
	//terminalFd uintptr
	Current int64 `json:"current,omitempty"`
	Total   int64 `json:"total,omitempty"`
	Start   int64 `json:"start,omitempty"`
}

// JSONMessage represents a JSON-encoded message regarding the status of a stream
type JSONMessage struct {
	Stream          string        `json:"stream,omitempty"`
	Status          string        `json:"status,omitempty"`
	Progress        *JSONProgress `json:"progressDetail,omitempty"`
	ProgressMessage string        `json:"progress,omitempty"` //deprecated
	ID              string        `json:"id,omitempty"`
	From            string        `json:"from,omitempty"`
	Time            int64         `json:"time,omitempty"`
	TimeNano        int64         `json:"timeNano,omitempty"`
	Error           *JSONError    `json:"errorDetail,omitempty"`
	ErrorMessage    string        `json:"error,omitempty"` //deprecated
	// Aux contains out-of-band data, such as digests for push signing.
	Aux *json.RawMessage `json:"aux,omitempty"`
}
