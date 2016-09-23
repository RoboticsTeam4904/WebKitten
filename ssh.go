package main

import (
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func NewSession(address, port, user, password, key string) (*ssh.Session, error) {
	authMethods := []ssh.AuthMethod{}
	if key != "" {
		signer, signerErr := KeySigner(key)
		if signerErr != nil {
			Error.Println(signerErr)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}
	if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: authMethods,
	}
	sshClient, sshClientErr := ssh.Dial("tcp", address+":"+port, sshConfig)
	if sshClientErr != nil {
		panic("Could not dial: " + sshClientErr.Error())
	}
	sshSession, sshSessionErr := sshClient.NewSession()
	if sshSessionErr != nil {
		panic("Could not establish session: " + sshSessionErr.Error())
	}
	defer func() {
		closeErr := sshClient.Close()
		Error.Println(closeErr)
	}()
	return sshSession, nil
}

func KeySigner(keyPath string) (ssh.Signer, error) {
	key, keyErr := ioutil.ReadFile(keyPath)
	if keyErr != nil {
		return nil, keyErr
	}
	signer, signerErr := ssh.ParsePrivateKey(key)
	if signerErr != nil {
		return nil, signerErr
	}
	return signer, nil
}
