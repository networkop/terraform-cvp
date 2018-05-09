package main

import cvpgo "github.com/fredhsu/cvpgo/client"

type CvpClient struct {
	Client    *cvpgo.CvpClient
	Container string
}

type CvpInfo struct {
	CvpAddress   string
	CvpUser      string
	CvpPwd       string
	CvpContainer string
}

func getCvpClient(c *CvpInfo) (*CvpClient, error) {

	// authenticating with CVP
	cvp := cvpgo.New(c.CvpAddress, c.CvpUser, c.CvpPwd)

	// client declarations
	client := CvpClient{
		Client:    &cvp,
		Container: c.CvpContainer,
	}

	return &client, nil
}
