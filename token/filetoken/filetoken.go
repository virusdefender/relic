//
// Copyright (c) SAS Institute Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package filetoken

import (
	"crypto"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/sassoftware/relic/config"
	"github.com/sassoftware/relic/lib/certloader"
	"github.com/sassoftware/relic/lib/passprompt"
	"github.com/sassoftware/relic/token"
)

type fileToken struct {
	config    *config.Config
	tokenConf *config.TokenConfig
	prompt    passprompt.PasswordGetter
	mu        sync.Mutex
}

type fileKey struct {
	keyConf *config.KeyConfig
	signer  crypto.Signer
}

func Open(conf *config.Config, tokenName string, prompt passprompt.PasswordGetter) (token.Token, error) {
	tconf, err := conf.GetToken(tokenName)
	if err != nil {
		return nil, err
	}
	return &fileToken{
		config:    conf,
		tokenConf: tconf,
		prompt:    prompt,
	}, nil
}

func (tok *fileToken) Ping() error {
	return nil
}

func (tok *fileToken) Close() error {
	return nil
}

func (tok *fileToken) Config() *config.TokenConfig {
	return tok.tokenConf
}

func (tok *fileToken) ListKeys(opts token.ListOptions) error {
	return errors.New("not implemented for tokens of type \"file\"")
}

func (tok *fileToken) GetKey(keyName string) (token.Key, error) {
	keyConf, err := tok.config.GetKey(keyName)
	if err != nil {
		return nil, err
	}
	if keyConf.KeyFile == "" {
		return nil, fmt.Errorf("key \"%s\" needs a KeyFile setting", keyName)
	}
	blob, err := ioutil.ReadFile(keyConf.KeyFile)
	if err != nil {
		return nil, err
	}
	/* TODO: keyring support
	loginFunc := func(pin string) (bool, error) {
	keyringUser := fmt.Sprintf("%s.%s", tok.tokenConf.Name(), keyName)
	if savedPass != "" {
		tok.mu.Lock()
		defer tok.mu.Unlock()
		if err := token.Login(tok.tokenConf, tok.prompt, loginFunc, keyringUser, ""); err != nil {
			return nil, err
		}
	}
	*/
	privateKey, err := certloader.ParseAnyPrivateKey(blob, tok.prompt)
	if err != nil {
		return nil, err
	}
	return &fileKey{
		keyConf: keyConf,
		signer:  privateKey.(crypto.Signer),
	}, nil
}

func (key *fileKey) Public() crypto.PublicKey {
	return key.signer.Public()
}

func (key *fileKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	return key.signer.Sign(rand, digest, opts)
}

func (key *fileKey) Config() *config.KeyConfig {
	return key.keyConf
}

func (key *fileKey) GetID() []byte {
	return nil
}

func (tok *fileToken) Import(keyName string, privKey crypto.PrivateKey) (token.Key, error) {
	return nil, errors.New("function not implemented for tokens of type \"file\"")
}

func (tok *fileToken) ImportCertificate(cert *x509.Certificate, labelBase string) error {
	return errors.New("function not implemented for tokens of type \"file\"")
}

func (tok *fileToken) Generate(keyName string, keyType token.KeyType, bits uint) (token.Key, error) {
	// TODO - probably useful
	return nil, errors.New("function not implemented for tokens of type \"file\"")
}

func (key *fileKey) ImportCertificate(cert *x509.Certificate) error {
	return errors.New("function not implemented for tokens of type \"file\"")
}
