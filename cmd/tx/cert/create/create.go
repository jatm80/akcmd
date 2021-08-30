package create

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/interact"
	"github.com/ovrclk/akash/x/cert/client/cli"
	"github.com/ovrclk/akash/x/cert/types"
	cutils "github.com/ovrclk/akash/x/cert/utils"
	"github.com/ovrclk/akcmd/client"
	"github.com/pkg/errors"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "create",
		Desc: "create/update api certificates",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{clientCMD(), serverCMD()},
	}

	return cmd
}

func clientCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name:   "client",
		Desc:   "create client api certificate",
		Config: config,
		Func: func(cmd *gcli.Command, args []string) error {
			return doCreateCmd(cmd, nil)
		},
	}

	return cmd
}

func serverCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name:   "server",
		Desc:   "create server api certificate",
		Config: config,
		Func:   doCreateCmd,
	}

	return cmd
}

var opts = struct {
	nbf       string
	naf       string
	rie       bool
	toGenesis bool
}{}

func config(cmd *gcli.Command) {
	client.AddTxFlagsToCmd(cmd)
	cmd.Required(flags.FlagFrom)

	cmd.StrOpt(&opts.nbf, "nbf", "", "",
		"certificate is not valid before this date. default current timestamp. RFC3339")
	cmd.StrOpt(&opts.naf, "naf", "", "",
		"certificate is not valid after this date. default 365d. days or RFC3339")
	cmd.BoolOpt(&opts.rie, "rie", "", false, "revoke current certificate if it already present on chain")

	// fixme shall we use gentx instead? 🤔
	cmd.BoolOpt(&opts.toGenesis, "to-genesis", "", false, "export certificate to genesis")
}

func handleCreate(cctx sdkclient.Context, pemFile string, domains []string) error {
	msg, err := createAuthPem(cctx, pemFile, domains)
	if err != nil {
		return err
	}

	if !opts.toGenesis {
		return client.BroadcastTX(cctx, msg)
	}

	return addCertToGenesis(cctx, types.GenesisCertificate{
		Owner: msg.Owner,
		Certificate: types.Certificate{
			State:  types.CertificateValid,
			Cert:   msg.Cert,
			Pubkey: msg.Pubkey,
		},
	})
}

func doCreateCmd(cmd *gcli.Command, domains []string) error {
	revokeIfExists := opts.rie

	cctx, err := client.GetClientTxContext()
	if err != nil {
		return err
	}

	fromAddress := cctx.GetFromAddress()

	pemFile := cctx.HomeDir + "/" + fromAddress.String() + ".pem"

	if _, err = os.Stat(pemFile); os.IsNotExist(err) {
		_ = cctx.PrintString(fmt.Sprintf("no certificate found for address %s. generating new...\n", fromAddress))

		return handleCreate(cctx, pemFile, domains)
	}

	cpem, err := cutils.LoadPEMForAccount(cctx, cctx.Keyring)
	if err != nil {
		return err
	}

	blk, _ := pem.Decode(cpem.Cert)
	x509cert, err := x509.ParseCertificate(blk.Bytes)
	if err != nil {
		return err
	}

	// if revoke-if-exists flag is true query is performed automatically
	// then certificate is being revoked (if valid) and file is removed without any prompts
	yes := revokeIfExists
	if !yes {
		yes = getConfirmation(fmt.Sprintf("certificate file for address %q already exist. check it on-chain status?", fromAddress))
	}

	if yes {
		params := &types.QueryCertificatesRequest{
			Filter: types.CertificateFilter{
				Owner:  cctx.FromAddress.String(),
				Serial: x509cert.SerialNumber.String(),
			},
		}

		res, err := types.NewQueryClient(cctx).Certificates(context.Background(), params)
		if err != nil {
			return err
		}

		removeFile := revokeIfExists

		if len(res.Certificates) == 0 {
			if !revokeIfExists {
				yes = getConfirmation("this certificate has not been found on chain. would you like to commit it?")

				if !yes {
					yes = getConfirmation("would you like to remove the file?")
					removeFile = yes
				} else {
					cpem, err := cutils.LoadPEMForAccount(cctx, cctx.Keyring)
					if err != nil {
						return err
					}

					msg := &types.MsgCreateCertificate{
						Owner:  fromAddress.String(),
						Cert:   cpem.Cert,
						Pubkey: cpem.Pub,
					}

					if err = msg.ValidateBasic(); err != nil {
						return err
					}

					return client.BroadcastTX(cctx, msg)
				}
			}
		} else {
			if res.Certificates[0].Certificate.IsState(types.CertificateValid) {
				msg := &types.MsgRevokeCertificate{
					ID: types.CertificateID{
						Owner:  cctx.FromAddress.String(),
						Serial: x509cert.SerialNumber.String(),
					},
				}

				err = client.BroadcastTX(cctx, msg)
				if err == nil {
					removeFile = true
				}
			}
		}

		if removeFile {
			if err = os.Remove(pemFile); err != nil {
				return err
			}

			_ = cctx.PrintString("generating new...\n")
			return handleCreate(cctx, pemFile, domains)
		}
	}

	return nil
}

func addCertToGenesis(cctx sdkclient.Context, cert types.GenesisCertificate) error {
	cdc := cctx.Codec

	serverCtx := client.GetServerContext()
	config := serverCtx.Config

	config.SetRoot(cctx.HomeDir)

	if err := cert.Validate(); err != nil {
		return errors.Errorf("failed to validate new genesis certificate: %v", err)
	}

	genFile := config.GenesisFile()
	appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
	if err != nil {
		return errors.Errorf("failed to unmarshal genesis state: %v", err)
	}

	certsGenState := types.GetGenesisStateFromAppState(cdc, appState)

	if certsGenState.Certificates.Contains(cert) {
		return errors.Errorf("cannot add already existing certificate")
	}
	certsGenState.Certificates = append(certsGenState.Certificates, cert)

	certsGenStateBz, err := cdc.MarshalJSON(certsGenState)
	if err != nil {
		return errors.Errorf("failed to marshal auth genesis state: %v", err)
	}

	appState[types.ModuleName] = certsGenStateBz

	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return errors.Errorf("failed to marshal application genesis state: %v", err)
	}

	genDoc.AppState = appStateJSON
	return genutil.ExportGenesisFile(genDoc, genFile)
}

func createAuthPem(cctx sdkclient.Context, pemFile string, domains []string) (*types.MsgCreateCertificate, error) {
	fromAddress := cctx.GetFromAddress()
	// note operation below needs more digging to ensure security. current implementation is more like example
	//      private key we generate has to be password protected
	//      from user prospective remembering/handling yet another password
	//      would be a subject of obliviousness. instead we utilize account's key
	//      to generate signature of it's address and use it as a password to encrypt
	//      private key.
	//      from security prospective this signature must never be exposed to prevent certificate leak.
	//      from other hand user must never obtain signature of it's own addresses in shell
	//      so yet again - to be discussed
	sig, _, err := cctx.Keyring.SignByAddress(fromAddress, fromAddress.Bytes())
	if err != nil {
		return nil, err
	}

	var priv *ecdsa.PrivateKey

	if priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader); err != nil {
		return nil, err
	}

	nbf := time.Now()
	naf := nbf.Add(time.Hour * 24 * 365)

	if val := opts.nbf; val != "" {
		nbf, err = time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
	}

	if val := opts.naf; val != "" {
		if strings.HasSuffix(val, "d") {
			days, err := strconv.ParseUint(strings.TrimSuffix(val, "d"), 10, 32)
			if err != nil {
				return nil, err
			}

			naf = nbf.Add(time.Hour * 24 * time.Duration(days))
		} else {
			naf, err = time.Parse(time.RFC3339, val)
			if err != nil {
				return nil, err
			}
		}
	}

	serialNumber := new(big.Int).SetInt64(time.Now().UTC().UnixNano())

	extKeyUsage := []x509.ExtKeyUsage{
		x509.ExtKeyUsageClientAuth,
	}

	if len(domains) > 0 {
		extKeyUsage = append(extKeyUsage, x509.ExtKeyUsageServerAuth)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: fromAddress.String(),
			ExtraNames: []pkix.AttributeTypeAndValue{
				{
					Type:  cli.AuthVersionOID,
					Value: "v0.0.1",
				},
			},
		},
		Issuer: pkix.Name{
			CommonName: fromAddress.String(),
		},
		NotBefore:             nbf,
		NotAfter:              naf,
		KeyUsage:              x509.KeyUsageDataEncipherment | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           extKeyUsage,
		BasicConstraintsValid: true,
	}

	var ips []net.IP

	for i := len(domains) - 1; i >= 0; i-- {
		if ip := net.ParseIP(domains[i]); ip != nil {
			ips = append(ips, ip)
			domains = append(domains[:i], domains[i+1:]...)
		}
	}

	if len(domains) != 0 || len(ips) != 0 {
		template.PermittedDNSDomainsCritical = true
		template.PermittedDNSDomains = domains
		template.DNSNames = domains
		template.IPAddresses = ips
	}

	var certDer []byte
	if certDer, err = x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv); err != nil {
		_ = cctx.PrintString(fmt.Sprintf("Failed to create certificate: %v\n", err))
		return nil, err
	}

	var keyDer []byte
	if keyDer, err = x509.MarshalPKCS8PrivateKey(priv); err != nil {
		return nil, err
	}

	var pubKeyDer []byte
	if pubKeyDer, err = x509.MarshalPKIXPublicKey(priv.Public()); err != nil {
		return nil, err
	}

	var blk *pem.Block
	// fixme #1182
	blk, err = x509.EncryptPEMBlock(rand.Reader, types.PemBlkTypeECPrivateKey, keyDer, sig, x509.PEMCipherAES256) // nolint: staticcheck
	if err != nil {
		_ = cctx.PrintString(fmt.Sprintf("failed to encrypt key file: %v\n", err))
		return nil, err
	}

	var pemOut *os.File
	if pemOut, err = os.OpenFile(pemFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return nil, err
	}

	defer func() {
		if err = pemOut.Close(); err != nil {
			_ = cctx.PrintString(fmt.Sprintf("failed to close key file: %v\n", err))
		} else {
			_ = os.Chmod(pemFile, 0400)
		}
	}()

	if err = pem.Encode(pemOut, &pem.Block{Type: types.PemBlkTypeCertificate, Bytes: certDer}); err != nil {
		_ = cctx.PrintString(fmt.Sprintf("failed to write certificate to pem file: %v\n", err))
		return nil, err
	}

	if err = pem.Encode(pemOut, blk); err != nil {
		_ = cctx.PrintString(fmt.Sprintf("failed to write key to pem file: %v\n", err))
		return nil, err
	}

	msg := &types.MsgCreateCertificate{
		Owner: fromAddress.String(),
		Cert: pem.EncodeToMemory(&pem.Block{
			Type:  types.PemBlkTypeCertificate,
			Bytes: certDer,
		}),
		Pubkey: pem.EncodeToMemory(&pem.Block{
			Type:  types.PemBlkTypeECPublicKey,
			Bytes: pubKeyDer,
		}),
	}

	if err = msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return msg, nil
}

func getConfirmation(prompt string) bool {
	return interact.Confirm(prompt)
}
