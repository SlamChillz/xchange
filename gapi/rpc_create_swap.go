package gapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"database/sql"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/pb"
	"github.com/slamchillz/xchange/utils"
	"github.com/slamchillz/xchange/common"
)

func (server *Server) CreateSwap(ctx context.Context, req *pb.CoinSwapRequest) (*pb.CoinSwapResponse, error) {
	payload, err := server.AuthenticateUser(ctx)
	if err != nil {
		return nil, authenticationError(err)
	}
	btcUsdRate := 34403.000000
	arg := db.GetPendingNetworkTransactionParams{
		CustomerID: payload.CustomerID,
		Network: req.Network,
		TransactionStatus: "PENDING",
	}
	count, err := server.storage.GetPendingNetworkTransaction(context.Background(), arg)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}
	if count > 0 {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprintf("You have a pending %s transaction on %s network", req.CoinName, req.Network))
	}
	address, err := common.GetSwapAddress(server.config, server.storage, req.GetCoinName(), req.GetNetwork(), arg.CustomerID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println(err)
			return nil, status.Errorf(codes.Internal, "Internal server error")
		}
	}
	currentUsdtRate := utils.RandomCoinswapRate()
	ngnEquivalent := float64(req.GetCoinAmountToSwap()) * currentUsdtRate
	if strings.ToUpper(req.GetCoinName()) == "BTC" {
		ngnEquivalent = ngnEquivalent * btcUsdRate
	}
	coinAmountToSwap := fmt.Sprintf("%f", req.GetCoinAmountToSwap())
	swapDetails, err := server.storage.CreateSwap(context.Background(), db.CreateSwapParams{
		CoinName: req.CoinName,
		CoinAmountToSwap: coinAmountToSwap,
		Network: req.GetNetwork(),
		PhoneNumber: req.GetPhoneNumber(),
		CoinAddress: address,
		TransactionRef: utils.RandomString(15),
		TransactionStatus: "PENDING",
		CurrentUsdtNgnRate: fmt.Sprintf("%f", currentUsdtRate),
		CustomerID: arg.CustomerID,
		NgnEquivalent: fmt.Sprintf("%f", ngnEquivalent),
		BankAccName: req.GetBankName(),
		BankAccNumber: req.GetBankAccNumber(),
		BankCode: req.GetBankCode(),
	})
	if err != nil {
		log.Panicln(err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return &pb.CoinSwapResponse{
		Id: swapDetails.ID,
		CoinName: swapDetails.CoinName,
		CoinAmountToSwap: req.CoinAmountToSwap,
		Network: swapDetails.Network,
		PhoneNumber: swapDetails.PhoneNumber,
		CoinAddress: swapDetails.CoinAddress,
		TransactionRef: swapDetails.TransactionRef,
		TransactionStatus: swapDetails.TransactionStatus,
		CurrentUsdtNgnRate: swapDetails.CurrentUsdtNgnRate,
		CustomerId: swapDetails.CustomerID,
		NgnEquivalent: swapDetails.NgnEquivalent,
		BankName: swapDetails.BankAccName,
		BankAccNumber: swapDetails.BankAccNumber,
		BankCode: swapDetails.BankCode,
	}, nil
}
