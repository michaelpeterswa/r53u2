package settings

import "os"

func SetAWSEnvironment(awss AWSSettings) error {
	err := os.Setenv("AWS_ACCESS_KEY_ID", awss.AWSAccessKeyId)
	if err != nil {
		return err
	}

	err = os.Setenv("AWS_SECRET_ACCESS_KEY", awss.AWSAccessKeySecret)
	if err != nil {
		return err
	}

	err = os.Setenv("AWS_DEFAULT_REGION", awss.AWSDefaultRegion)
	if err != nil {
		return err
	}

	return nil
}
