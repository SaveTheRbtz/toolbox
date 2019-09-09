package api_auth_impl

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/watermint/toolbox/infra/api/api_auth"
	"github.com/watermint/toolbox/infra/api/api_context"
	"github.com/watermint/toolbox/infra/api/api_context_impl"
	"github.com/watermint/toolbox/infra/app"
	app2 "github.com/watermint/toolbox/legacy/app"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"os"
)

type EcCachedAuth struct {
	peerName string
	tokens   map[string]string
	ec       *app2.ExecContext
	auth     api_auth.Console
}

func (z *EcCachedAuth) Auth(tokenType string) (ctx api_context.Context, err error) {
	if tok, e := z.tokens[tokenType]; e {
		tc := api_auth.TokenContainer{
			Token:     tok,
			TokenType: tokenType,
		}
		return api_context_impl.NewLegacy(z.ec, tc), nil
	}
	if ctx, err = z.auth.Auth(tokenType); err != nil {
		return nil, err
	} else {
		z.updateCache(tokenType, ctx)
		return ctx, nil
	}
}

func (z *EcCachedAuth) init() {
	z.tokens = make(map[string]string)

	if z.loadFile() == nil {
		return // return on success
	}
}

func (z *EcCachedAuth) cacheFile(kind string) string {
	px := sha256.Sum224([]byte(z.peerName))
	pn := fmt.Sprintf("%x.%s", px, kind)
	return z.ec.FileOnSecretsPath(pn)
}

func (z *EcCachedAuth) compatibleCachedFile() string {
	return z.cacheFile("tokens")
}
func (z *EcCachedAuth) secureCachedFile() string {
	return z.cacheFile("t")
}
func (z *EcCachedAuth) loadBytes(tb []byte) error {
	err := json.Unmarshal(tb, &z.tokens)
	if err != nil {
		z.ec.Log().Debug("unable to unmarshal tokens file", zap.Error(err))
		return err
	}
	return nil
}

func (z *EcCachedAuth) loadFile() error {
	if ex, err := z.loadCompatibleFile(); err == nil {
		return nil
	} else if !ex && app.BuilderKey != "" {
		_, err := z.loadSecureFile()
		return err
	}
	return nil
}

func (z *EcCachedAuth) loadCompatibleFile() (exists bool, err error) {
	tf := z.compatibleCachedFile()
	_, err = os.Stat(tf)
	if os.IsNotExist(err) {
		//z.ec.Log().Debug("token file not found", zap.String("path", tf))
		return false, err
	}
	z.ec.Log().Debug("Loading token file", zap.String("file", tf))
	tb, err := ioutil.ReadFile(tf)
	if err != nil {
		z.ec.Log().Debug("unable to read tokens file", zap.String("path", tf), zap.Error(err))
		return false, err
	}
	return true, z.loadBytes(tb)
}

func (z *EcCachedAuth) loadSecureFile() (exists bool, err error) {
	if app.BuilderKey == "" {
		z.ec.Log().Debug("Use compatible token file in dev mode")
		return false, errors.New("dev mode")
	}
	tf := z.secureCachedFile()
	z.ec.Log().Debug("Loading token file", zap.String("file", tf))
	_, err = os.Stat(tf)
	if os.IsNotExist(err) {
		//z.ec.Log().Debug("token file not found", zap.String("path", tf))
		return false, err
	}
	tb, err := ioutil.ReadFile(tf)
	if err != nil {
		z.ec.Log().Debug("unable to read tokens file", zap.String("path", tf), zap.Error(err))
		return false, err
	}

	key := []byte(app.BuilderKey + app.Name)
	key32 := sha256.Sum224([]byte(key))
	kb := make([]byte, 32)
	copy(kb[:], key32[:])

	bk, err := aes.NewCipher([]byte(kb))
	if err != nil {
		return false, err
	}
	gcm, err := cipher.NewGCM(bk)
	if err != nil {
		return false, err
	}
	ns := gcm.NonceSize()
	nonce, ct := tb[:ns], tb[ns:]
	v, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return false, err
	}
	return true, z.loadBytes(v)
}

func (z *EcCachedAuth) updateCompatible(tb []byte) error {
	tf := z.compatibleCachedFile()
	err := ioutil.WriteFile(tf, tb, 0600)
	if err != nil {
		z.ec.Log().Debug("unable to write tokens into file", zap.Error(err))
		return err
	}
	return nil
}

func (z *EcCachedAuth) updateSecure(tb []byte) error {
	key := []byte(app.BuilderKey + app.Name)
	key32 := sha256.Sum224([]byte(key))
	kb := make([]byte, 32)
	copy(kb[:], key32[:])
	bk, err := aes.NewCipher([]byte(kb))
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(bk)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	sealed := gcm.Seal(nonce, nonce, tb, nil)

	tf := z.secureCachedFile()
	err = ioutil.WriteFile(tf, sealed, 0600)
	if err != nil {
		z.ec.Log().Debug("unable to write tokens into file", zap.Error(err))
		return err
	}
	return nil
}

func (z *EcCachedAuth) updateCache(tokenType string, ctx api_context.Context) {
	// Do not store tokens into file
	if z.ec.NoCacheToken() {
		return
	}

	switch tc := ctx.(type) {
	case api_auth.TokenContext:
		z.tokens[tokenType] = tc.Token().Token
		tb, err := json.Marshal(z.tokens)
		if err != nil {
			z.ec.Log().Debug("unable to marshal tokens", zap.Error(err))
			return
		}
		if app.BuilderKey == "" {
			z.updateCompatible(tb)
		} else {
			z.updateSecure(tb)
		}
	}
}