package main

import cvpgo "github.com/networkop/cvpgo/client"

type CvpClient struct {
	Client *cvpgo.CvpClient
}

type CvpInfo struct {
	CvpAddress string
	CvpUser    string
	CvpPwd     string
}

func getCvpClient(c *CvpInfo) (*CvpClient, error) {

	// authenticating with CVP
	cvp := cvpgo.New(c.CvpAddress, c.CvpUser, c.CvpPwd)

	// client declarations
	client := CvpClient{
		Client: &cvp,
	}

	return &client, nil
}
