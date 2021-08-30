package cert

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/gookit/gcli/v3"
	"github.com/ovrclk/akash/x/cert/types"
	cutils "github.com/ovrclk/akash/x/cert/utils"
	"github.com/ovrclk/akcmd/client"
	"github.com/ovrclk/akcmd/cmd/tx/cert/create"
)

func Cmd() *gcli.Command {
	cmd := &gcli.Command{
		Name: "cert",
		Desc: "Certificates transaction subcommands",
		Func: func(cmd *gcli.Command, args []string) error {
			cmd.ShowHelp()
			return nil
		},
		Subs: []*gcli.Command{create.Cmd(), revokeCMD()},
	}

	return cmd
}

var serial string

func revokeCMD() *gcli.Command {
	cmd := &gcli.Command{
		Name: "revoke",
		Desc: "revoke api certificate",
		Config: func(cmd *gcli.Command) {
			cmd.StrOpt(&serial, "serial", "", "", "revoke certificate by serial number")
			client.AddTxFlagsToCmd(cmd)
			cmd.Required(flags.FlagFrom)
		},
		Func: func(cmd *gcli.Command, args []string) error {
			cctx, err := client.GetClientTxContext()
			if err != nil {
				return err
			}

			if serial != "" {
				if _, valid := new(big.Int).SetString(serial, 10); !valid {
					return errors.New("invalid value in serial flag. expected integer")
				}
			} else {
				cpem, err := cutils.LoadPEMForAccount(cctx, cctx.Keyring)
				if err != nil {
					return err
				}

				blk, _ := pem.Decode(cpem.Cert)
				cert, err := x509.ParseCertificate(blk.Bytes)
				if err != nil {
					return err
				}

				serial = cert.SerialNumber.String()
			}

			msg := &types.MsgRevokeCertificate{
				ID: types.CertificateID{
					Owner:  cctx.FromAddress.String(),
					Serial: serial,
				},
			}

			return client.BroadcastTX(cctx, msg)
		},
	}

	return cmd
}
